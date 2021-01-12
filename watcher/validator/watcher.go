// Package validator
package validator

import (
	"context"

	"go.uber.org/zap"

	"github.com/ngdlong91/kai-watcher/external/telegram"
	"github.com/ngdlong91/kai-watcher/kardia"
	"github.com/ngdlong91/kai-watcher/types"
)

const (
	defaultMissedBlockStep = 10
)

type watcher struct {
	node            kardia.Node
	alert           telegram.Client
	missedBlockStep uint64
	logger          *zap.Logger
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
	missedBlockStep := cfg.MissedBlockStep
	if missedBlockStep == 0 {
		missedBlockStep = defaultMissedBlockStep
	}

	return &watcher{
		node:            node,
		alert:           alert,
		missedBlockStep: missedBlockStep,
		logger:          cfg.Logger,
	}, nil
}

func (w *watcher) Run(ctx context.Context) error {
	validators, err := w.node.Validators(ctx)
	if err != nil {
		return err
	}
	b := types.Builder()
	for _, v := range validators {
		validator := b.Validator(v)
		if w.isAlert(validator) {
			alertMsg := missedBlockAlert(validator)
			if err := w.alert.Send(alertMsg); err != nil {
				return err
			}
		}
	}
	// Update list validator to storage

	return nil
}

func (w *watcher) isAlert(v types.Validator) bool {
	missedBlockCounter := v.MissedBlockCounter
	return missedBlockCounter > 0 && missedBlockCounter%w.missedBlockStep == 0
}
