package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/ngdlong91/kai-watcher/watcher/staking"
	"github.com/ngdlong91/kai-watcher/watcher/validator"
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
	gCfg.Logger = logger

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

	go validator.WatchValidators(ctx, gCfg, 5*time.Second)
	go staking.WatchStakingSMC(ctx, gCfg, 5*time.Second)
	//go whale.WatchWhaleTransaction(ctx, gCfg, 5*time.Second)

	<-waitExit
	logger.Info("Stopped")
}
