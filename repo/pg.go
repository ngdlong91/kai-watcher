package repo

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"os"
)

//////////////////////////////
//////// FOR TESTING ////////
//////////////////////////////

func SetupDBTestKit(ctx context.Context) (*pgxpool.Pool, *zap.Logger) {
	lgr, _ := zap.NewDevelopment()
	pool, err := pgxpool.Connect(ctx, os.Getenv("DB_URL"))
	if err != nil {
		panic(err)
	}
	return pool, lgr
}
