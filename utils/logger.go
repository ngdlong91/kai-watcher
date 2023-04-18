// Package main
package utils

import (
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/ngdlong91/kai-watcher/cfg"
)

type LoggerConfig struct {
	ServerMode string
	LogLevel   string
}

func NewLogger(config LoggerConfig) (*zap.Logger, error) {
	logCfg := zap.NewProductionConfig()
	switch config.ServerMode {
	case cfg.ModeDev:
		logCfg = zap.NewDevelopmentConfig()
		logCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	case cfg.ModeProduction:
		logCfg = zap.NewProductionConfig()
	}

	switch config.LogLevel {
	case "info":
		logCfg.Level.SetLevel(zapcore.InfoLevel)
	case "debug":
		logCfg.Level.SetLevel(zapcore.DebugLevel)
	case "warn":
		logCfg.Level.SetLevel(zapcore.WarnLevel)
	default:
		logCfg.Level.SetLevel(zapcore.InfoLevel)
	}
	sentryOpts := zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.RegisterHooks(core, func(entry zapcore.Entry) error {
			if config.ServerMode != cfg.ModeProduction {
				return nil
			}
			var isCapture bool
			e := sentry.NewEvent()
			switch entry.Level {
			case zap.WarnLevel:
				isCapture = true
				e.Level = sentry.LevelWarning
			case zap.ErrorLevel:
				e.Level = sentry.LevelError
				isCapture = true
			default:
				isCapture = false
			}
			if isCapture {
				e.Message = entry.Message
				sentry.CaptureEvent(e)
			}

			return nil
		})
	})

	return logCfg.Build(sentryOpts)
}

func LoggerForPackage(lgr *zap.Logger, packageName string) *zap.Logger {
	return lgr.With(zap.String("package", packageName))
}

func LoggerForMethod(lgr *zap.Logger, method string) *zap.Logger {
	return lgr.With(zap.String("method", method))
}
