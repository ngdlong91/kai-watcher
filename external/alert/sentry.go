// Package sentry
package alert

import (
	"github.com/getsentry/sentry-go"
)

type Config struct {
	DSN         string
	Environment string
}

func NewAlert(cfg Config) error {
	opts := sentry.ClientOptions{
		Dsn:         cfg.DSN,
		Environment: cfg.Environment,
	}
	if err := sentry.Init(opts); err != nil {
		return err
	}
	return nil
}
