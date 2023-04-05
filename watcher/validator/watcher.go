// Package validator
package validator

import (
	"context"
	"github.com/ngdlong91/kai-watcher/cfg"
	"time"

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

func WatchValidators(ctx context.Context, cfg cfg.EnvConfig, interval time.Duration) {
	lgr := cfg.Logger.With(zap.String("Watcher", "Validators"))
	vCfg := Config{
		URL:        cfg.KardiaTrustedNodes[0],
		Logger:     lgr,
		AlertToken: cfg.TelegramToken,
	}
	watcher, err := NewWatcher(vCfg)
	if err != nil {
		lgr.Error("cannot create watcher", zap.Error(err))
		panic(err)
	}
	lgr.Info("Start watch validator")
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			if err := watcher.Run(ctx); err != nil {
				lgr.Error("Run error", zap.Error(err))
				continue
			}
		}
	}
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
