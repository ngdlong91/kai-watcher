package whale

import (
	"fmt"

	"github.com/ngdlong91/kai-watcher/types"
	"github.com/ngdlong91/kai-watcher/utils"
)

func newLevelOneAlert(tx *types.Transaction) string {
	return fmt.Sprintf(" 🚨\nFrom: [%s](https://explorer.kardiachain.io/address/%s) \nTo: [%s](https://explorer.kardiachain.io/address/%s) \nValue:  %s KAI \nTxHash: [%s](https://explorer.kardiachain.io/tx/%s)", tx.From, tx.From, tx.To, tx.To, utils.HumanizeCurrency(tx.Value), tx.Hash, tx.Hash)
}

func newLevelTwoAlert(tx *types.Transaction) string {
	return fmt.Sprintf(" 🚨🚨\nFrom: [%s](https://explorer.kardiachain.io/address/%s) \nTo: [%s](https://explorer.kardiachain.io/address/%s) \nValue:  %s KAI \nTxHash: [%s](https://explorer.kardiachain.io/tx/%s)", tx.From, tx.From, tx.To, tx.To, utils.HumanizeCurrency(tx.Value), tx.Hash, tx.Hash)
}

func newLevelThreeAlert(tx *types.Transaction) string {
	return fmt.Sprintf(" 🚨🚨🚨 \nFrom: [%s](https://explorer.kardiachain.io/address/%s) \nTo: [%s](https://explorer.kardiachain.io/address/%s) \nValue:  %s KAI \nTxHash: [%s](https://explorer.kardiachain.io/tx/%s)", tx.From, tx.From, tx.To, tx.To, utils.HumanizeCurrency(tx.Value), tx.Hash, tx.Hash)
}

func newLevelFourAlert(tx *types.Transaction) string {
	return fmt.Sprintf(" 🚨🚨🚨🚨 \nFrom: [%s](https://explorer.kardiachain.io/address/%s) \nTo: [%s](https://explorer.kardiachain.io/address/%s) \nValue:  %s KAI \nTxHash: [%s](https://explorer.kardiachain.io/tx/%s)", tx.From, tx.From, tx.To, tx.To, utils.HumanizeCurrency(tx.Value), tx.Hash, tx.Hash)
}
