package main

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/ngdlong91/kai-watcher/watcher/validator"
)

func WatchValidators(ctx context.Context, interval time.Duration) {
	lgr := logger.With(zap.String("Watcher", "Validators"))
	cfg := validator.Config{
		URL:        gCfg.KardiaTrustedNodes[0],
		Logger:     logger,
		AlertToken: gCfg.TelegramToken,
	}
	watcher, err := validator.NewWatcher(cfg)
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
