package main

import (
	"context"
	"fmt"
	"github.com/ngdlong91/kai-watcher/external/alert"
	"github.com/ngdlong91/kai-watcher/kclient"
	"github.com/ngdlong91/kai-watcher/repo"
	"github.com/ngdlong91/kai-watcher/utils"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"github.com/ngdlong91/kai-watcher/cfg"
	"github.com/ngdlong91/kai-watcher/tasks"
)

var (
	gCfg   cfg.EnvConfig
	logger *zap.Logger
)

func preload() {

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

	ctx := context.Background()

	_, cancel := context.WithCancel(ctx)
	sigCh := make(chan os.Signal, 1)
	waitExit := make(chan bool)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for range sigCh {
			cancel()
			waitExit <- true
		}
	}()

	fmt.Printf("Global config: %+v \n", gCfg)

	pool, err := pgxpool.Connect(ctx, gCfg.StorageURI)
	if err != nil {
		fmt.Println("Err", err.Error())
		panic(err)
	}

	node, err := kclient.NewNode(gCfg.KardiaTrustedNodes[0], gCfg.Logger)
	if err != nil {
		panic(err)
	}

	delegatorDB := &repo.Delegator{
		Conn:   pool,
		Logger: logger,
	}
	fetchDelegatorTask := tasks.FetchDelegators{
		Logger:      logger,
		Node:        node,
		DelegatorDB: delegatorDB,
		Pool:        pool,
	}

	c := cron.New()
	_, err = c.AddFunc(gCfg.CronFetchDelegators, fetchDelegatorTask.Execute)
	if err != nil {
		panic(err.Error())
	}
	defer c.Stop()
	go func() {
		c.Start()
	}()

	<-waitExit
	fmt.Println("Exit")
}
