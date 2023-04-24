package whale

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ngdlong91/kai-watcher/cache"
	"github.com/ngdlong91/kai-watcher/entities"
	"github.com/ngdlong91/kai-watcher/kclient"
	"github.com/ngdlong91/kai-watcher/repo"
	"github.com/ngdlong91/kai-watcher/types"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	"github.com/ngdlong91/kai-watcher/cfg"
	"github.com/ngdlong91/kai-watcher/external/telegram"
	"github.com/ngdlong91/kai-watcher/utils"
)

type Watcher struct {
	node               *kclient.Node
	levelOneLimit      decimal.Decimal
	levelTwoLimit      decimal.Decimal
	levelThreeLimit    decimal.Decimal
	levelFourLimit     decimal.Decimal
	alert              telegram.Client
	currentBlockHeight uint64
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

func WatchWhaleTransaction(ctx context.Context, cfg cfg.EnvConfig, interval time.Duration) {
	lgr := cfg.Logger.With(zap.String("Watcher", "WhaleAlert"))
	wCfg := Config{
		URL:             cfg.KardiaTrustedNodes[0],
		Logger:          lgr,
		AlertToken:      cfg.TelegramToken,
		AlertTo:         cfg.TelegramGroup,
		LevelOneLimit:   cfg.LevelOneLimit,
		LevelTwoLimit:   cfg.LevelTwoLimit,
		LevelThreeLimit: cfg.LevelThreeLimit,
		LevelFourLimit:  cfg.LevelFourLimit,
	}
	watcher, err := NewWatcher(wCfg)
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

	lgr.Info("Start whale alert watcher")
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			lgr.Info("Tick")
			if err := watcher.Run(ctx); err != nil {
				lgr.Error("watcher err", zap.Error(err))
				continue
			}
		}
	}
}

func NewWatcher(cfg Config) (*Watcher, error) {
	node, err := kclient.NewNode(cfg.URL, cfg.Logger)
	if err != nil {
		return nil, err
	}
	alertCfg := telegram.Config{
		Token:   cfg.AlertToken,
		GroupID: cfg.AlertTo,
		Logger:  cfg.Logger,
	}

	cfg.Logger.Info("Watcher", zap.Any("Config", cfg))
	alert, err := telegram.NewClient(alertCfg)
	if err != nil {
		return nil, err
	}
	levelOneLimit, err := decimal.NewFromString(cfg.LevelOneLimit)
	if err != nil {
		return nil, err
	}

	levelTwoLimit, err := decimal.NewFromString(cfg.LevelTwoLimit)
	if err != nil {
		return nil, err
	}

	levelThreeLimit, err := decimal.NewFromString(cfg.LevelThreeLimit)
	if err != nil {
		return nil, err
	}

	levelFourLimit, err := decimal.NewFromString(cfg.LevelFourLimit)
	if err != nil {
		return nil, err
	}
	watcher := &Watcher{
		node:            node,
		alert:           alert,
		levelOneLimit:   levelOneLimit,
		levelTwoLimit:   levelTwoLimit,
		levelThreeLimit: levelThreeLimit,
		levelFourLimit:  levelFourLimit,
		logger:          cfg.Logger,
	}

	return watcher, nil
}

func (w *Watcher) Run(ctx context.Context) error {
	//todo: change get latest block number to subscribe newHeads event
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
	block, err := w.node.BlockByHeight(ctx, latestBlockNumber)
	if err != nil {
		return err
	}

	for _, tx := range block.Txs {
		// Tx Value in decimal
		w.getNames(ctx, tx)
		txValue := utils.ToDecimal(tx.Value, 18)
		var alertMsg string
		if txValue.Cmp(w.levelFourLimit) >= 0 {
			alertMsg = newLevelFourAlert(tx)
		} else if txValue.Cmp(w.levelThreeLimit) >= 0 {
			alertMsg = newLevelThreeAlert(tx)
		} else if txValue.Cmp(w.levelTwoLimit) >= 0 {
			alertMsg = newLevelTwoAlert(tx)
		} else if txValue.Cmp(w.levelOneLimit) >= 0 {
			alertMsg = newLevelOneAlert(tx)
		}
		if alertMsg != "" {
			lgr.Info("Process tx with value", zap.String("Value", tx.Value), zap.String("hash", tx.Hash))
			if err := w.alert.Send(alertMsg); err != nil {
				return err
			}
		}
	}

	return nil
}

func (w *Watcher) getNames(ctx context.Context, tx *types.Transaction) {
	from := strings.ToLower(tx.From)
	to := strings.ToLower(tx.To)
	var fromName, toName string
	var err error
	fromName, err = w.WalletCache.WalletName(ctx, from)
	if err != nil || fromName == "" {
		fromName, err = w.WalletRepo.Retrieve(ctx, from)
		if err != nil {
			w.logger.Error("cannot get wallet name from db", zap.Error(err))
			fromName = ""
		} else {
			if err := w.WalletCache.SetWallets(ctx, from, fromName); err != nil {
				w.logger.Error("cannot set wallet", zap.Error(err))
				fromName = ""
			}
		}
	}
	tx.FromName = fromName

	toName, err = w.WalletCache.WalletName(ctx, to)
	if err != nil || toName == "" {
		toName, err = w.WalletRepo.Retrieve(ctx, to)
		if err != nil {
			w.logger.Error("cannot get wallet name from db", zap.Error(err))
			toName = ""
		} else {
			if err := w.WalletCache.SetWallets(ctx, to, toName); err != nil {
				w.logger.Error("cannot set wallet", zap.Error(err))
				fromName = ""
			}
		}
	}
	tx.ToName = toName

	fmt.Printf("TX: %+v \n", tx)
}
