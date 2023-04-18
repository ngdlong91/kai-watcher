// Package telegram
package telegram

import (
	"go.uber.org/zap"
)

type Config struct {
	Token   string
	GroupID int64
	Logger  *zap.Logger
}
