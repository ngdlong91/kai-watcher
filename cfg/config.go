// Package cfg
package cfg

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

const (
	ModeDev        = "dev"
	ModeProduction = "prod"
)

type EnvConfig struct {
	ServerMode string
	Port       string

	LogLevel string

	StakingAddress string

	KardiaURLs         []string
	KardiaTrustedNodes []string

	StorageURI     string
	StorageMinConn int
	StorageMaxConn int
	StorageIsFlush bool

	CacheHost     string
	CachePassword string

	LevelOneLimit   string
	LevelTwoLimit   string
	LevelThreeLimit string
	LevelFourLimit  string
	ValidatorLimit  string

	// SentryConfiguration
	SentryDSN string

	// Telegram token
	TelegramToken string
	TelegramGroup int64

	CronFetchDelegators string

	Logger *zap.Logger
}

var serverMode = "dev"

func Load() (EnvConfig, error) {
	if serverMode == ModeProduction {
		godotenv.Load()
	} else {
		godotenv.Load("dev.env")
	}

	var (
		kardiaTrustedNodes []string
		kardiaURLs         []string
	)
	kardiaTrustedNodesStr := os.Getenv("KARDIA_TRUSTED_NODES")
	if kardiaTrustedNodesStr != "" {
		kardiaTrustedNodes = strings.Split(kardiaTrustedNodesStr, ",")
	} else {
		panic("missing trusted RPC URLs in config")
	}
	kardiaURLsStr := os.Getenv("KARDIA_URL")
	if kardiaURLsStr != "" {
		kardiaURLs = strings.Split(kardiaURLsStr, ",")
	} else {
		panic("missing RPC URLs in config")
	}

	storageMinConnStr := os.Getenv("STORAGE_MIN_CONN")
	storageMinConn, err := strconv.Atoi(storageMinConnStr)
	if err != nil {
		storageMinConn = 8
	}

	storageMaxConnStr := os.Getenv("STORAGE_MAX_CONN")
	storageMaxConn, err := strconv.Atoi(storageMaxConnStr)
	if err != nil {
		storageMaxConn = 32
	}

	storageIsFlushStr := os.Getenv("STORAGE_IS_FLUSH")
	storageIsFLush, err := strconv.ParseBool(storageIsFlushStr)
	if err != nil {
		storageIsFLush = false
	}

	telegramGroupIDStr := os.Getenv("TELEGRAM_GROUP")
	fmt.Println("Telegram group ID", telegramGroupIDStr)
	telegramGroup, err := strconv.ParseInt(telegramGroupIDStr, 10, 64)
	if err != nil {
		fmt.Println("err", err.Error())
		telegramGroup = 0
	}

	cfg := EnvConfig{
		ServerMode: os.Getenv("SERVER_MODE"),
		Port:       os.Getenv("PORT"),

		LogLevel: os.Getenv("LOG_LEVEL"),

		StakingAddress: os.Getenv("STAKING_ADDR"),

		KardiaURLs:         kardiaURLs,
		KardiaTrustedNodes: kardiaTrustedNodes,

		StorageURI:     os.Getenv("STORAGE_URI"),
		StorageMinConn: storageMinConn,
		StorageMaxConn: storageMaxConn,
		StorageIsFlush: storageIsFLush,

		CacheHost:     os.Getenv("CACHE_HOST"),
		CachePassword: os.Getenv("CACHE_PASSWORD"),

		LevelOneLimit:   os.Getenv("LEVEL_ONE_LIMIT"),
		LevelTwoLimit:   os.Getenv("LEVEL_TWO_LIMIT"),
		LevelThreeLimit: os.Getenv("LEVEL_THREE_LIMIT"),
		LevelFourLimit:  os.Getenv("LEVEL_FOUR_LIMIT"),

		ValidatorLimit: os.Getenv("VALIDATOR_LIMIT"),

		SentryDSN:     os.Getenv("SENTRY_DNS"),
		TelegramToken: os.Getenv("TELEGRAM_TOKEN"),
		TelegramGroup: telegramGroup,

		CronFetchDelegators: os.Getenv("CRON_FETCH_DELEGATORS"),
	}

	return cfg, nil
}
