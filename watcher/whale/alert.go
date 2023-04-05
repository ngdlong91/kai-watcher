package whale

import (
	"fmt"
	"github.com/ngdlong91/kai-watcher/kardia"
	"github.com/ngdlong91/kai-watcher/utils"
)

func newLevelOneAlert(tx *kardia.Transaction) string {
	return fmt.Sprintf(" 🚨\nFrom: [%s](https://explorer.kardiachain.io/address/%s) \nTo: [%s](https://explorer.kardiachain.io/address/%s) \nValue:  %s KAI \nTxHash: [%s](https://explorer.kardiachain.io/tx/%s)", tx.From, tx.From, tx.To, tx.To, utils.HumanizeCurrency(tx.Value), tx.Hash, tx.Hash)
}

func newLevelTwoAlert(tx *kardia.Transaction) string {
	return fmt.Sprintf(" 🚨🚨\nFrom: [%s](https://explorer.kardiachain.io/address/%s) \nTo: [%s](https://explorer.kardiachain.io/address/%s) \nValue:  %s KAI \nTxHash: [%s](https://explorer.kardiachain.io/tx/%s)", tx.From, tx.From, tx.To, tx.To, utils.HumanizeCurrency(tx.Value), tx.Hash, tx.Hash)
}

func newLevelThreeAlert(tx *kardia.Transaction) string {
	return fmt.Sprintf(" 🚨🚨🚨 \nFrom: [%s](https://explorer.kardiachain.io/address/%s) \nTo: [%s](https://explorer.kardiachain.io/address/%s) \nValue:  %s KAI \nTxHash: [%s](https://explorer.kardiachain.io/tx/%s)", tx.From, tx.From, tx.To, tx.To, utils.HumanizeCurrency(tx.Value), tx.Hash, tx.Hash)
}

func newLevelFourAlert(tx *kardia.Transaction) string {
	return fmt.Sprintf(" 🚨🚨🚨🚨 \nFrom: [%s](https://explorer.kardiachain.io/address/%s) \nTo: [%s](https://explorer.kardiachain.io/address/%s) \nValue:  %s KAI \nTxHash: [%s](https://explorer.kardiachain.io/tx/%s)", tx.From, tx.From, tx.To, tx.To, utils.HumanizeCurrency(tx.Value), tx.Hash, tx.Hash)
}
