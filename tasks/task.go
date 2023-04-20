package tasks

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ngdlong91/kai-watcher/kclient"
	"github.com/ngdlong91/kai-watcher/repo"
	"go.uber.org/zap"
	"os"
)

// Task interface for implement schedule tasks
type Task interface {
	Execute()
}

func SetupTestTaskEnv() (*pgxpool.Pool, *kclient.Node, *zap.Logger) {
	conn, lgr := repo.SetupDBTestKit(context.Background())
	client, _ := kclient.NewNode(os.Getenv("KARDIA_URL"), lgr)
	return conn, client, lgr
}
