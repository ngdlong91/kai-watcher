package validator

import (
	"go.uber.org/zap"
)

type Config struct {
	URL    string
	Logger *zap.Logger

	AlertToken      string
	MissedBlockStep uint64
}
