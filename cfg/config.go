// Package cfg
package cfg

import (
	"os"
	"strconv"
	"strings"
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

	StorageDriver  string
	StorageURI     string
	StorageDB      string
	StorageMinConn int
	StorageMaxConn int
	StorageIsFlush bool

	// SentryConfiguration
	SentryDSN string

	// Telegram token
	TelegramToken string
}

func Load() (EnvConfig, error) {

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

	cfg := EnvConfig{
		ServerMode: os.Getenv("SERVER_MODE"),
		Port:       os.Getenv("PORT"),

		LogLevel: os.Getenv("LOG_LEVEL"),

		StakingAddress: os.Getenv("STAKING_ADDR"),

		KardiaURLs:         kardiaURLs,
		KardiaTrustedNodes: kardiaTrustedNodes,

		StorageDriver:  os.Getenv("STORAGE_DRIVER"),
		StorageURI:     os.Getenv("STORAGE_URI"),
		StorageDB:      os.Getenv("STORAGE_DB"),
		StorageMinConn: storageMinConn,
		StorageMaxConn: storageMaxConn,
		StorageIsFlush: storageIsFLush,

		SentryDSN:     os.Getenv("SENTRY_DNS"),
		TelegramToken: os.Getenv("TELEGRAM_TOKEN"),
	}

	return cfg, nil
}
