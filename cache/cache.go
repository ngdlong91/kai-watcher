package cache

import (
	"github.com/go-redis/redis/v8"
	"os"
)

func SetupTestCache() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("CACHE_HOST"),
		DB:       0,
		Password: os.Getenv("CACHE_PASSWORD"),
	})
}
