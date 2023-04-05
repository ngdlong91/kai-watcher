package whale

import (
	"fmt"
	"github.com/ngdlong91/kai-watcher/kardia"
	"github.com/ngdlong91/kai-watcher/utils"
)

func newLevelOneAlert(tx *kardia.Transaction) string {
	return fmt.Sprintf(" 🚨<b>From </b> %s -> %s : %s KAI. TxHash: %s", tx.From, tx.To, utils.HumanizeCurrency(tx.Value), tx.Hash)
}

func newLevelTwoAlert(tx *kardia.Transaction) string {
	return fmt.Sprintf(" 🚨🚨<b>From </b> %s -> %s : %s KAI. TxHash: %s", tx.From, tx.To, utils.HumanizeCurrency(tx.Value), tx.Hash)
}

func newLevelThreeAlert(tx *kardia.Transaction) string {
	return fmt.Sprintf(" 🚨🚨🚨<b>From </b> %s -> %s : %s KAI. TxHash: %s", tx.From, tx.To, utils.HumanizeCurrency(tx.Value), tx.Hash)
}

func newLevelFourAlert(tx *kardia.Transaction) string {
	return fmt.Sprintf(" 🚨🚨🚨🚨<b>From </b> %s -> %s : %s KAI. TxHash: %s", tx.From, tx.To, utils.HumanizeCurrency(tx.Value), tx.Hash)
}
