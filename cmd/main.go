package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/ngdlong91/kai-watcher/watcher/staking"
	"github.com/ngdlong91/kai-watcher/watcher/validator"
	"github.com/ngdlong91/kai-watcher/watcher/whale"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ngdlong91/kai-watcher/cfg"
	"github.com/ngdlong91/kai-watcher/external/alert"
	"github.com/ngdlong91/kai-watcher/utils"
)

var (
	gCfg   cfg.EnvConfig
	logger *zap.Logger
)

func preload() {
	if err := godotenv.Load(); err != nil {
		panic(err.Error())
	}

	runtime.GOMAXPROCS(runtime.NumCPU())
	tempCfg, err := cfg.Load()
	if err != nil {
		panic(err)
	}
	gCfg = tempCfg

	alertCfg := alert.Config{
		DSN:         gCfg.SentryDSN,
		Environment: gCfg.ServerMode,
	}
	if err := alert.NewAlert(alertCfg); err != nil {
		panic(err)
	}

	lgrCfg := utils.LoggerConfig{
		ServerMode: gCfg.ServerMode,
		LogLevel:   gCfg.LogLevel,
	}
	tempLgr, err := utils.NewLogger(lgrCfg)
	if err != nil {
		panic("cannot init logger")
	}
	logger = tempLgr.With(zap.String("services", "watcher"))

}

func main() {
	preload()
	defer func() {
		if err := recover(); err != nil {
			logger.Error("cannot recover")
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal, 1)
	waitExit := make(chan bool)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for range sigCh {
			cancel()
			waitExit <- true
		}
	}()

	go WatchValidators(ctx, 5*time.Second)
	go WatchStakingSMC(ctx, 5*time.Second)
	go WatchWhaleTransaction(ctx, 5*time.Second)

	<-waitExit
	logger.Info("Stopped")
}

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
