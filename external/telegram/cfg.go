// Package telegram
package telegram

import (
	"go.uber.org/zap"
)

type Config struct {
	Token  string
	Logger *zap.Logger
}
