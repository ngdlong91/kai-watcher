// Package telegram
package telegram

import (
	tgAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	//endpoint := fmt.Sprintf("%s&parse_mode=MarkdownV2", tgAPI.APIEndpoint)
	//bot, err := tgAPI.NewBotAPIWithAPIEndpoint(cfg.Token, endpoint)
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
	//finalMsg := tgAPI.EscapeText(tgAPI.ModeHTML, msg)
	message := tgAPI.NewMessage(-873461799, msg)
	message.ParseMode = tgAPI.ModeMarkdown
	message.DisableWebPagePreview = true

	if _, err := c.bot.Send(message); err != nil {
		return err
	}
	return nil
}
