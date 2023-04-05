package main

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/ngdlong91/kai-watcher/watcher/staking"
)

func WatchStakingSMC(ctx context.Context, interval time.Duration) {
	lgr := logger.With(zap.String("Watcher", "Staking"))
	cfg := staking.Config{
		URL:            gCfg.KardiaTrustedNodes[0],
		Logger:         logger,
		AlertToken:     gCfg.TelegramToken,
		StakingAddress: gCfg.StakingAddress,
	}
	watcher, err := staking.NewWatcher(cfg)
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
