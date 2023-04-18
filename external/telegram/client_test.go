package telegram

import (
	"fmt"
	tgAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ngdlong91/kai-watcher/utils"
	"go.uber.org/zap"
	"testing"
)

func TestClient(t *testing.T) {
	lgr, _ := zap.NewDevelopment()
	bot, _ := tgAPI.NewBotAPI("_")
	c := &client{
		bot: bot,
		lgr: lgr,
	}

	msg := fmt.Sprintf(" 🚨**text** %s >>> %s : %s KAI. TxHash: %s", "FromSender", "ToRecv", utils.HumanizeCurrency("1000"), "hasttx")

	c.Send(msg)
}
