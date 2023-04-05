package whale

import (
	"context"
	"github.com/ngdlong91/kai-watcher/external/telegram"
	"github.com/ngdlong91/kai-watcher/kardia"
	"github.com/ngdlong91/kai-watcher/utils"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type watcher struct {
	node               kardia.Node
	levelOneLimit      decimal.Decimal
	levelTwoLimit      decimal.Decimal
	levelThreeLimit    decimal.Decimal
	levelFourLimit     decimal.Decimal
	alert              telegram.Client
	currentBlockHeight uint64
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
	watcher := &watcher{
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

func (w *watcher) Run(ctx context.Context) error {
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
	lgr.Info("------------------------1")
	block, err := w.node.BlockByHeight(ctx, latestBlockNumber)
	if err != nil {
		return err
	}

	lgr.Info("Inspect block", zap.Uint64("NumTxs", block.NumTxs))
	for _, tx := range block.Txs {
		// Tx Value in decimal
		lgr.Info("Process tx with value", zap.String("Value", tx.Value), zap.String("hash", tx.Hash))
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
		lgr.Info("Send telegram with message", zap.String("Alert", alertMsg))
		if alertMsg != "" {
			if err := w.alert.Send(alertMsg); err != nil {
				return err
			}
		}
	}

	return nil
}
