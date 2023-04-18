package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
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

func main() {
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

	pool, err := pgxpool.Connect(ctx, gCfg.StorageURI)
	if err != nil {
		panic(err)
	}

	fetchDelegatorTask := tasks.FetchDelegators{
		Logger: logger,
		Pool:   pool,
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
