// Package staking
package staking

import (
	"context"

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

	block, err := w.node.BlockByHeight(ctx, latestBlockNumber)
	if err != nil {
		return err
	}

	b := types.Builder()
	for id, r := range block.Receipts {
		if r.ContractAddress != w.stakingAddress {
			continue
		}
		tx := block.Txs[id]
		if tx == nil {
			//todo: something really bad here
			continue
		}
		transaction := b.Transaction(tx)
		// Get decode data of smc call
		fc, err := w.node.DecodeInputData(tx.ContractAddress, tx.InputData)
		if err != nil {
			return err
		}

		v, err := w.node.Validator(ctx, tx.ContractAddress)
		if err != nil {
			return err
		}
		validator := b.Validator(v)
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
	}

	return nil
}
