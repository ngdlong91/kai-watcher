package kclient

import (
	"math/big"
	"time"

	"github.com/kardiachain/go-kardia/lib/common"
)

type FunctionCall struct {
	Function   string                 `json:"function"`
	MethodID   string                 `json:"methodID"`
	MethodName string                 `json:"methodName"`
	Arguments  map[string]interface{} `json:"arguments"`
}

type SMCCallArgs struct {
	From     string   `json:"from"`     // the sender of the 'transaction'
	To       *string  `json:"to"`       // the destination contract (nil for contract creation)
	Gas      uint64   `json:"gas"`      // if 0, the call executes with near-infinite gas
	GasPrice *big.Int `json:"gasPrice"` // HYDRO <-> gas exchange ratio
	Value    *big.Int `json:"value"`    // amount of HYDRO sent along with the call
	Data     string   `json:"data"`     // input data, usually an ABI-encoded contract method invocation
}

type RPCValidator struct {
	Name       [32]uint8
	Signer     common.Address
	SMCAddress common.Address
	Tokens     *big.Int
	Jailed     bool

	DelegationShares      *big.Int
	AccumulatedCommission *big.Int
	UbdEntryCount         *big.Int
	UpdateTime            *big.Int
	MinSelfDelegation     *big.Int
	Status                uint8

	UnbondingTime   *big.Int
	UnbondingHeight *big.Int
	Address         string `json:"address" bson:"address,omitempty"`
	SmcAddress      string `json:"smcAddress" bson:"smcAddress,omitempty"`
	Role            int    `json:"role" bson:"role"`

	VotingPowerPercentage string `json:"votingPowerPercentage" bson:"votingPowerPercentage,omitempty"`
	StakedAmount          string `json:"stakedAmount" bson:"stakedAmount,omitempty"`

	CommissionRate  string          `json:"commissionRate" bson:"commissionRate,omitempty"`
	TotalDelegators int             `json:"totalDelegators" bson:"totalDelegators,omitempty"`
	MaxRate         string          `json:"maxRate" bson:"maxRate,omitempty"`
	MaxChangeRate   string          `json:"maxChangeRate" bson:"maxChangeRate,omitempty"`
	SigningInfo     *SigningInfo    `json:"signingInfo" bson:"signingInfo,omitempty"`
	Delegators      []*RPCDelegator `json:"delegators,omitempty" bson:"delegators,omitempty"`
}

type SigningInfo struct {
	StartHeight        *big.Int `json:"startHeight"`
	IndexOffset        *big.Int `json:"indexOffset"`
	Tombstoned         bool     `json:"tombstoned"`
	MissedBlockCounter *big.Int `json:"missedBlockCounter"`
	IndicatorRate      float64  `json:"indicatorRate"`
	JailedUntil        *big.Int `json:"jailedUntil"`
}

type Commission struct {
	Rate          *big.Int
	MaxRate       *big.Int
	MaxChangeRate *big.Int
}

type RPCDelegator struct {
	Address      common.Address `json:"address"`
	StakedAmount *big.Int       `json:"stakedAmount"`
	Reward       *big.Int       `json:"reward"`
}

type Log struct {
	Address       string                 `json:"address,omitempty" bson:"address"`
	MethodName    string                 `json:"methodName,omitempty" bson:"methodName"`
	ArgumentsName string                 `json:"argumentsName,omitempty" bson:"argumentsName"`
	Arguments     map[string]interface{} `json:"arguments,omitempty" bson:"arguments"`
	Topics        []string               `json:"topics,omitempty" bson:"topics"`
	Data          string                 `json:"data,omitempty" bson:"data"`
	BlockHeight   uint64                 `json:"blockHeight,omitempty" bson:"blockHeight"`
	Time          time.Time              `json:"time" bson:"time"`
	TxHash        string                 `json:"transactionHash"  bson:"transactionHash"`
	TxIndex       uint                   `json:"transactionIndex,omitempty" bson:"transactionIndex"`
	BlockHash     string                 `json:"blockHash,omitempty" bson:"blockHash"`
	Index         uint                   `json:"logIndex,omitempty" bson:"logIndex"`
	Removed       bool                   `json:"removed,omitempty" bson:"removed"`
}
