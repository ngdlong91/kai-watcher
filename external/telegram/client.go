// Package telegram
package telegram

import (
	tgAPI "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
)

type Client interface {
	Send(msg string) error
}

type client struct {
	lgr *zap.Logger
	bot *tgAPI.BotAPI
}

func NewClient(cfg Config) (Client, error) {
	bot, err := tgAPI.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, err
	}

	c := &client{
		lgr: cfg.Logger,
		bot: bot,
	}
	return c, nil
}

func (c *client) Send(msg string) error {
	message := tgAPI.NewMessage(-467712071, msg)
	if _, err := c.bot.Send(message); err != nil {
		return err
	}
	return nil
}
