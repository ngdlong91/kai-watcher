// Package staking
package staking

import (
	"context"
	"github.com/ngdlong91/kai-watcher/cfg"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/ngdlong91/kai-watcher/external/telegram"
	"github.com/ngdlong91/kai-watcher/kardia"

	"github.com/ngdlong91/kai-watcher/types"
)

type watcher struct {
	node               kardia.Node
	stakingAddress     string
	alert              telegram.Client
	currentBlockHeight uint64
	validators         []*kardia.Validator
	lastFetch          int64
	logger             *zap.Logger
}

func NewWatcher(cfg Config) (*watcher, error) {
	node, err := kardia.NewNode(cfg.URL, cfg.Logger)
	if err != nil {
		return nil, err
	}
	alertCfg := telegram.Config{
		Token:  cfg.AlertToken,
		Logger: cfg.Logger,
	}
	alert, err := telegram.NewClient(alertCfg)
	if err != nil {
		return nil, err
	}
	watcher := &watcher{
		node:           node,
		alert:          alert,
		stakingAddress: cfg.StakingAddress,
		logger:         cfg.Logger,
	}

	return watcher, nil
}

func WatchStakingSMC(ctx context.Context, cfg cfg.EnvConfig, interval time.Duration) {
	lgr := cfg.Logger.With(zap.String("Watcher", "Staking"))
	sCfg := Config{
		URL:            cfg.KardiaTrustedNodes[0],
		Logger:         lgr,
		AlertToken:     cfg.TelegramToken,
		StakingAddress: cfg.StakingAddress,
	}
	watcher, err := NewWatcher(sCfg)
	if err != nil {
		lgr.Error("cannot create watcher", zap.Error(err))
		panic(err)
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

func (w *watcher) Run(ctx context.Context) error {
	if w.lastFetch+600 < time.Now().Unix() {
		validators, err := w.node.Validators(ctx)
		if err != nil {
			return err
		}
		w.validators = validators
		w.lastFetch = time.Now().Unix()
		w.logger.Info("Validator info", zap.Any("va", validators))
	}

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

	if latestBlockNumber%10 == 0 {
		lgr.Info("Process block", zap.Uint64("Number", latestBlockNumber))
	}

	block, err := w.node.BlockByHeight(ctx, latestBlockNumber)
	if err != nil {
		return err
	}

	b := types.Builder()
	isSkip := true
	var validator *kardia.Validator
	for id, r := range block.Receipts {
		lgr.Info("REceipt", zap.Any("r", r))
		for _, l := range r.Logs {
			lgr.Info("L", zap.Any("l", l))
			for vid, v := range w.validators {
				lgr.Info("V", zap.Any("v", v))
				lgr.Info("VSMC", zap.Any("SMC", v.SMCAddress.String()))
				if strings.ToLower(l.Address) == strings.ToLower(v.SMCAddress.String()) {
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
			transaction := b.Transaction(tx)
			//// Get decode data of smc call
			fc, err := w.node.DecodeInputData(l.Address, l.Data)
			if err != nil {
				lgr.Error("Cannot decode input", zap.Error(err))
				return err
			}

			var alertMsg string
			switch fc.MethodName {
			case DelegateMethod:
				alertMsg = newDelegateAlert(transaction, validator)
			case UndelegatedMethod:
				alertMsg = newUndelegatedAlert(transaction, validator)
			}
			if fc.MethodName == DelegateMethod {

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
