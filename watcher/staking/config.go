package staking

import (
	"go.uber.org/zap"
)

type Config struct {
	URL            string
	Logger         *zap.Logger
	StakingAddress string

	AlertToken string
}

const (
	DelegateMethod    = "delegate"
	UndelegatedMethod = "undelegated"
)
