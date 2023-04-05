package main

import (
	"context"
	"github.com/ngdlong91/kai-watcher/watcher/whale"
	"go.uber.org/zap"
	"time"
)

func WatchWhaleTransaction(ctx context.Context, interval time.Duration) {
	lgr := logger.With(zap.String("Watcher", "WhaleAlert"))
	cfg := whale.Config{
		URL:             gCfg.KardiaTrustedNodes[0],
		Logger:          logger,
		AlertToken:      gCfg.TelegramToken,
		LevelOneLimit:   gCfg.LevelOneLimit,
		LevelTwoLimit:   gCfg.LevelTwoLimit,
		LevelThreeLimit: gCfg.LevelThreeLimit,
		LevelFourLimit:  gCfg.LevelFourLimit,
	}
	watcher, err := whale.NewWatcher(cfg)
	if err != nil {
		lgr.Error("cannot create watcher", zap.Error(err))
		panic(err)
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
