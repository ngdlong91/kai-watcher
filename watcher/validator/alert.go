// Package kai
package validator

import (
	"fmt"

	"github.com/ngdlong91/kai-watcher/types"
)

func missedBlockAlert(v types.Validator) string {
	return fmt.Sprintf("Validator: %s [%s]. Total missed block: [%d]", v.Name, v.OwnerAddress, v.MissedBlockCounter)
}
