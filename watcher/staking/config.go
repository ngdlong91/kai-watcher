package staking

import (
	"go.uber.org/zap"
)

type Config struct {
	URL            string
	Logger         *zap.Logger
	ValidatorLimit string

	AlertToken string
}

const (
	DelegateMethod    = "Delegate"
	UndelegatedMethod = "Undelegate"
)
