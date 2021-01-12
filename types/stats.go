package types

import (
	"time"
)

type StakingStats struct {
	Name                       string `bson:"name"`
	TotalValidators            int    `json:"totalValidators" bson:"totalValidators"`
	TotalProposers             int    `json:"totalProposers" bson:"totalProposers"`
	TotalCandidates            int    `json:"totalCandidates" bson:"totalCandidates"`
	TotalDelegators            int    `json:"totalDelegators" bson:"totalDelegators"`
	TotalStakedAmount          string `json:"totalStakedAmount" bson:"totalStakedAmount"`
	TotalValidatorStakedAmount string `json:"totalValidatorStakedAmount" bson:"totalValidatorStakedAmount"`
	TotalDelegatorStakedAmount string `json:"totalDelegatorStakedAmount" bson:"totalDelegatorStakedAmount"`
}

type AddressStats struct {
	Name           string `bson:"name"`
	TotalAddresses uint64 `json:"totalAddresses" bson:"totalAddresses"`
	TotalContracts uint64 `json:"totalContracts" bson:"totalContracts"`
}

type TransactionStats struct {
	Name           string `bson:"name"`
	TotalTx        uint64 `json:"totalTx" bson:"totalTx"`
	ContractTx     uint64 `json:"contractInteract" bson:"contractInteract"`
	NormalTx       uint64 `json:"normalTx" bson:"normalTx"`
	AddressInvoked uint64 `json:"addressInvoked" bson:"addressInvoked"`
}

type DailyStats struct {
	Name      string    `bson:"name"`
	Timeline  time.Time `json:"timeline" bson:"timeline"`
	Txs       uint64    `json:"txs" bson:"txs"`
	Addresses uint64    `json:"addresses" bson:"addresses"`
	Contracts uint64    `json:"contracts" bson:"contracts"`
	Staking   uint64    `json:"staking" bson:"staking"`
}
