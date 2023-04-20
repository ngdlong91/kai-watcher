// Package types
package types

type UnbondedRecord struct {
	Balances        string `json:"balance"`
	CompletionTimes string `json:"completionTime"`
}

type DelegatorInfo struct {
	Address          string  `json:"address"`
	Balance          float64 `json:"balance"`
	ValidatorRecords []ValidatorStakeRecord
	TotalStaked      float64 `json:"total_staked"`
	TotalAmount      float64 `json:"total_amount"`
}

type ValidatorStakeRecord struct {
	ValidatorName       string `json:"validator_name"`
	ValidatorAddress    string `json:"validator_address"`
	ValidatorSMCAddress string `json:"validator_smc_address"`
	Amount              string `json:"amount"`
}
