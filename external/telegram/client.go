// Package telegram
package telegram

import (
	"fmt"
	tgAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type Client interface {
	Send(msg string) error
}

type client struct {
	lgr     *zap.Logger
	bot     *tgAPI.BotAPI
	groupID int64
}

func NewClient(cfg Config) (Client, error) {
	fmt.Println("Config", cfg)
	bot, err := tgAPI.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, err
	}

	c := &client{
		lgr:     cfg.Logger,
		bot:     bot,
		groupID: cfg.GroupID,
	}
	return c, nil
}

func (c *client) Send(msg string) error {
	fmt.Println("New message to group id", c.groupID)
	message := tgAPI.NewMessage(c.groupID, msg)
	message.ParseMode = tgAPI.ModeMarkdown
	message.DisableWebPagePreview = true

	if _, err := c.bot.Send(message); err != nil {
		return err
	}
	return nil
}
