package cache

import (
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
)

func SetupTestCache() (*redis.Client, *zap.Logger) {
	godotenv.Load("dev.env")
	lgr, _ := zap.NewDevelopment()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("CACHE_HOST"),
		DB:       0,
		Password: os.Getenv("CACHE_PASSWORD"),
	})

	return redisClient, lgr
}
