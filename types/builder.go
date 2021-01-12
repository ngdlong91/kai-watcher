// Package types
package types

import (
	"github.com/ngdlong91/kai-watcher/kardia"
	"github.com/ngdlong91/kai-watcher/utils"
)

type builder struct {
}

func Builder() *builder {
	return &builder{}
}

func (b *builder) Validator(v *kardia.Validator) Validator {
	validator := Validator{
		OwnerAddress:       v.Signer.Hex(),
		SmcAddress:         v.SMCAddress.Hex(),
		Name:               utils.Uint8ToString(v.Name),
		IsJailed:           v.Jailed,
		MissedBlockCounter: v.SigningInfo.MissedBlockCounter.Uint64(),
		JailedUntil:        v.SigningInfo.JailedUntil.Uint64(),
	}
	return validator
}

func (b *builder) Transaction(tx *kardia.Transaction) Transaction {
	transaction := Transaction{}
	return transaction
}

func (b *builder) FunctionCall(fc *kardia.FunctionCall) FunctionCall {
	funcCall := FunctionCall{}
	return funcCall
}
