// Package staking
package staking

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ngdlong91/kai-watcher/cache"
	"github.com/ngdlong91/kai-watcher/entities"
	"github.com/ngdlong91/kai-watcher/repo"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	"github.com/ngdlong91/kai-watcher/cfg"
	"github.com/ngdlong91/kai-watcher/external/telegram"
	"github.com/ngdlong91/kai-watcher/kclient"
	"github.com/ngdlong91/kai-watcher/types"
	"github.com/ngdlong91/kai-watcher/utils"
)

type Watcher struct {
	node               *kclient.Node
	alert              telegram.Client
	currentBlockHeight uint64
	validators         []*types.Validator
	lastFetch          int64
	limit              decimal.Decimal
	logger             *zap.Logger

	WalletCache interface {
		SetWallets(ctx context.Context, wallet string, name string) error
		WalletName(ctx context.Context, wallet string) (string, error)
	}

	WalletRepo interface {
		Insert(ctx context.Context, e *entities.Wallet) error
		Retrieve(ctx context.Context, wallet string) (string, error)
	}
}

func NewWatcher(cfg Config) (*Watcher, error) {
	node, err := kclient.NewNode(cfg.URL, cfg.Logger)
	if err != nil {
		return nil, err
	}
	alertCfg := telegram.Config{
		Token:   cfg.AlertToken,
		Logger:  cfg.Logger,
		GroupID: cfg.AlertTo,
	}
	alert, err := telegram.NewClient(alertCfg)
	if err != nil {
		return nil, err
	}
	limit, err := decimal.NewFromString(cfg.ValidatorLimit)
	if err != nil {
		return nil, err
	}
	watcher := &Watcher{
		node:   node,
		alert:  alert,
		limit:  limit,
		logger: cfg.Logger,
	}

	return watcher, nil
}

func WatchStakingSMC(ctx context.Context, cfg cfg.EnvConfig, interval time.Duration) {
	lgr := cfg.Logger.With(zap.String("Watcher", "Staking"))
	sCfg := Config{
		URL:            cfg.KardiaTrustedNodes[0],
		Logger:         lgr,
		AlertToken:     cfg.TelegramToken,
		AlertTo:        cfg.TelegramGroup,
		ValidatorLimit: cfg.ValidatorLimit,
	}

	watcher, err := NewWatcher(sCfg)
	if err != nil {
		lgr.Error("cannot create watcher", zap.Error(err))
		panic(err)
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.CacheHost,
		DB:       0,
		Password: cfg.CachePassword,
	})

	watcher.WalletCache = &cache.Wallet{
		Client: redisClient,
		Logger: lgr,
	}
	pool, err := pgxpool.Connect(ctx, cfg.StorageURI)
	if err != nil {
		panic(err)
	}
	watcher.WalletRepo = &repo.Wallet{
		Logger: lgr,
		Pool:   pool,
	}
	lgr.Info("Start staking watcher")
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			if err := watcher.Run(ctx); err != nil {
				continue
			}
		}
	}
}

func (w *Watcher) Run(ctx context.Context) error {
	if w.lastFetch+600 < time.Now().Unix() {
		validators, err := w.node.Validators(ctx)
		if err != nil {
			return err
		}
		w.validators = validators
		w.lastFetch = time.Now().Unix()
		w.logger.Info("Validator info", zap.Int("ValidatorSize", len(validators)))
		for _, v := range validators {
			w.logger.Info("VInfo", zap.Any("V", v))
		}
	}

	lgr := w.logger
	latestBlockNumber, err := w.node.LatestBlockNumber(ctx)
	if err != nil {
		return err
	}
	if w.currentBlockHeight == latestBlockNumber {
		lgr.Debug("Same latest block. Wait for next block")
		return nil
	}
	w.currentBlockHeight = latestBlockNumber

	if latestBlockNumber%10 == 0 {
		lgr.Info("Process block", zap.Uint64("Number", latestBlockNumber))
	}

	block, err := w.node.BlockByHeight(ctx, latestBlockNumber)
	if err != nil {
		return err
	}

	isSkip := true
	var validator *types.Validator
	for id, r := range block.Receipts {
		for _, l := range r.Logs {
			for vid, v := range w.validators {
				if strings.ToLower(l.Address) == strings.ToLower(v.SMCAddress) {
					isSkip = false
					validator = w.validators[vid]
					break
				}
			}
			if isSkip {
				continue
			}
			tx := block.Txs[id]
			if tx == nil {
				//todo: something really bad here
				lgr.Error("tx nil")
				continue
			}

			abi := w.node.ValidatorABI()
			//// Get decode data of smc call
			unpackedLog, err := kclient.UnpackLog(&l, abi)
			if err != nil {
				lgr.Error("Cannot decode input", zap.Error(err))
				return err
			}

			lgr.Info("UnpackedLog", zap.Any("L", unpackedLog))

			var alertMsg string
			switch unpackedLog.MethodName {
			case DelegateMethod:
				delegateAmount := utils.ToDecimal(fmt.Sprintf("%v", l.Arguments["_amount"]), 18)
				if delegateAmount.Cmp(w.limit) > 0 {
					alertMsg = newDelegateAlert(unpackedLog, validator)
				}

			case UndelegatedMethod:
				undelegateAmount := utils.ToDecimal(fmt.Sprintf("%v", l.Arguments["_amount"]), 18)
				if undelegateAmount.Cmp(w.limit) > 0 {
					alertMsg = newUndelegatedAlert(unpackedLog, validator)
				}

			}
			if alertMsg != "" {
				if err := w.alert.Send(alertMsg); err != nil {
					return err
				}
			}
			isSkip = true
		}

	}

	return nil
}
