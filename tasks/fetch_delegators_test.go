package tasks

import (
	"github.com/joho/godotenv"
	"github.com/ngdlong91/kai-watcher/repo"
	"testing"
)

func init() {
	if err := godotenv.Load("dev.env"); err != nil {
		panic(err)
	}
}

func TestCron_FetchDelegator(t *testing.T) {
	conn, node, lgr := SetupTestTaskEnv()
	delegatorRepo := &repo.Delegator{
		Conn:   conn,
		Logger: lgr,
	}
	task := &FetchDelegators{
		Logger:      lgr,
		Pool:        conn,
		node:        node,
		DelegatorDB: delegatorRepo,
	}

	task.Execute()
}
