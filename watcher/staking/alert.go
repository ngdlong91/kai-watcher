// Package kai
package staking

import (
	"fmt"
	"github.com/ngdlong91/kai-watcher/kardia"

	"github.com/ngdlong91/kai-watcher/types"
)

func newUndelegatedAlert(tx types.Transaction, v *kardia.Validator) string {
	return fmt.Sprintf("New undelegated: Address [%s] to %s [%s]: %s KAI. Details: %s", tx.From, v.Name, v.SMCAddress.String(), tx.Value, tx.Hash)
}

func newDelegateAlert(tx types.Transaction, v *kardia.Validator) string {
	return fmt.Sprintf("New delegate: Address [%s] to %s [%s]: %s KAI. Details: %s", tx.From, v.Name, v.SMCAddress.String(), tx.Value, tx.Hash)
}
