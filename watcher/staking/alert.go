// Package kai
package staking

import (
	"fmt"
	"github.com/ngdlong91/kai-watcher/kardia"
	"github.com/ngdlong91/kai-watcher/utils"
)

func newUndelegatedAlert(l *kardia.Log, v *kardia.Validator) string {
	return fmt.Sprintf("🥊🥊🥊🥊Undelegated \nAddress: [%s](https://explorer.kardiachain.io/address/%s) \nTo: [%s](https://explorer.kardiachain.io/address/%s) \nValue: %s KAI \nTxHash: [%s](https://explorer.kardiachain.io/tx/%s)",
		fmt.Sprintf("%v", l.Arguments["_delAddr"]), fmt.Sprintf("%v", l.Arguments["_delAddr"]), utils.ValidatorName(v.Name), v.SMCAddress.String(), utils.HumanizeCurrency(fmt.Sprintf("%v", l.Arguments["_amount"])), l.TxHash, l.TxHash)
}

func newDelegateAlert(l *kardia.Log, v *kardia.Validator) string {
	return fmt.Sprintf("💰💰💰💰 Delegate \nAddress: [%s](https://explorer.kardiachain.io/address/%s) \nTo: [%s](https://explorer.kardiachain.io/address/%s) \nValue: %s KAI \n TxHash: [%s](https://explorer.kardiachain.io/tx/%s)",
		fmt.Sprintf("%v", l.Arguments["_delAddr"]), fmt.Sprintf("%v", l.Arguments["_delAddr"]), utils.ValidatorName(v.Name), v.SMCAddress.String(), utils.HumanizeCurrency(fmt.Sprintf("%v", l.Arguments["_amount"])), l.TxHash, l.TxHash)
}
