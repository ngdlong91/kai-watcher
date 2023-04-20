// Package types
package types

import (
	"github.com/ngdlong91/kai-watcher/utils"
	"math/big"

	"github.com/kardiachain/go-kardia/lib/common"
)

const (
	// BloomByteLength represents the number of bytes used in a header log bloom.
	BloomByteLength = 256

	// BloomBitLength represents the number of bits used in a header log bloom.
	BloomBitLength = 8 * BloomByteLength
)

type Bloom [BloomByteLength]byte

type StakingStats struct {
	TotalValidators            int    `json:"totalValidators"`
	TotalProposers             int    `json:"totalProposers"`
	TotalCandidates            int    `json:"totalCandidates"`
	TotalDelegators            int    `json:"totalDelegators"`
	TotalStakedAmount          string `json:"totalStakedAmount"`
	TotalValidatorStakedAmount string `json:"totalValidatorStakedAmount"`
	TotalDelegatorStakedAmount string `json:"totalDelegatorStakedAmount"`
	MaxValidator               int    `json:"maxValidator"`
}

type ValidatorStats struct {
	TotalValidators            int    `json:"totalValidators"`
	TotalProposers             int    `json:"totalProposers"`
	TotalCandidates            int    `json:"totalCandidates"`
	TotalDelegators            int    `json:"totalDelegators"`
	TotalStakedAmount          string `json:"totalStakedAmount"`
	TotalValidatorStakedAmount string `json:"totalValidatorStakedAmount"`
	TotalDelegatorStakedAmount string `json:"totalDelegatorStakedAmount"`
}

type Validators struct {
	TotalValidators            int          `json:"totalValidators"`
	TotalProposers             int          `json:"totalProposers"`
	TotalCandidates            int          `json:"totalCandidates"`
	TotalDelegators            int          `json:"totalDelegators"`
	TotalStakedAmount          string       `json:"totalStakedAmount"`
	TotalValidatorStakedAmount string       `json:"totalValidatorStakedAmount"`
	TotalDelegatorStakedAmount string       `json:"totalDelegatorStakedAmount"`
	Validators                 []*Validator `json:"validators"`
}

type Validator struct {
	Name       string
	Signer     string
	SMCAddress string
	Tokens     *big.Int
	Jailed     bool

	DelegationShares      *big.Int
	AccumulatedCommission *big.Int
	UbdEntryCount         *big.Int
	UpdateTime            *big.Int
	MinSelfDelegation     *big.Int
	UnbondingTime         *big.Int
	UnbondingHeight       *big.Int

	Role                  int `json:"role" bson:"role"`
	Status                uint8
	VotingPowerPercentage string `json:"votingPowerPercentage" bson:"votingPowerPercentage,omitempty"`
	TotalDelegators       int    `json:"totalDelegators" bson:"totalDelegators,omitempty"`

	CommissionRate *big.Int `json:"commissionRate" bson:"commissionRate,omitempty"`
	MaxRate        *big.Int `json:"maxRate" bson:"maxRate,omitempty"`
	MaxChangeRate  *big.Int `json:"maxChangeRate" bson:"maxChangeRate,omitempty"`

	SigningInfo *SigningInfo `json:"signingInfo" bson:"signingInfo,omitempty"`
	Delegators  []*Delegator `json:"delegators,omitempty" bson:"delegators,omitempty"`
}

type RPCValidator struct {
	Name                  [32]uint8       `json:"name"`
	ValAddr               common.Address  `json:"validatorAddress"`
	ValStakingSmc         common.Address  `json:"valStakingSmc"`
	Tokens                *big.Int        `json:"tokens"`
	Jailed                bool            `json:"jailed"`
	DelegationShares      *big.Int        `json:"delegationShares"`
	AccumulatedCommission *big.Int        `json:"accumulatedCommission"`
	MinSelfDelegation     *big.Int        `json:"minSelfDelegation"`
	UbdEntryCount         *big.Int        `json:"ubdEntryCount"`
	UpdateTime            *big.Int        `json:"updateTime"`
	Status                uint8           `json:"status"`
	Role                  int             `json:"role"`
	UnbondingTime         *big.Int        `json:"unbondingTime"`
	UnbondingHeight       *big.Int        `json:"unbondingHeight"`
	CommissionRate        *big.Int        `json:"commissionRate,omitempty"`
	MaxRate               *big.Int        `json:"maxRate,omitempty"`
	MaxChangeRate         *big.Int        `json:"maxChangeRate,omitempty"`
	SigningInfo           *SigningInfo    `json:"signingInfo"`
	Delegators            []*RPCDelegator `json:"delegators,omitempty"`
}

func (r RPCValidator) ToValidator() Validator {
	return Validator{
		Name:                  utils.ValidatorName(r.Name),
		Signer:                r.ValAddr.String(),
		SMCAddress:            r.ValStakingSmc.String(),
		Tokens:                nil,
		Jailed:                false,
		DelegationShares:      nil,
		AccumulatedCommission: nil,
		UbdEntryCount:         nil,
		UpdateTime:            nil,
		MinSelfDelegation:     nil,
		UnbondingTime:         nil,
		UnbondingHeight:       nil,
		Role:                  0,
		Status:                0,
		VotingPowerPercentage: "",
		TotalDelegators:       0,
		CommissionRate:        nil,
		MaxRate:               nil,
		MaxChangeRate:         nil,
		SigningInfo:           nil,
		Delegators:            nil,
	}
}

type RPCDelegator struct {
	Address      common.Address `json:"address"`
	StakedAmount *big.Int       `json:"stakedAmount"`
	Reward       *big.Int       `json:"reward"`
}

type Delegator struct {
	ValidatorSMCAddress string `json:"validatorSMCAddress" bson:"validatorSMCAddress,omitempty"`
	Address             string `json:"address" bson:"address,omitempty"`
	StakedAmount        string `json:"stakedAmount" bson:"stakedAmount,omitempty"`
	Reward              string `json:"reward" bson:"reward,omitempty"`
}

type SlashEvents struct {
	Period   string `json:"period"`
	Fraction string `json:"fraction"`
	Height   string `json:"height"`
}

type SigningInfo struct {
	StartHeight        uint64  `json:"startHeight"`
	IndexOffset        uint64  `json:"indexOffset"`
	Tombstoned         bool    `json:"tombstoned"`
	MissedBlockCounter uint64  `json:"missedBlockCounter"`
	IndicatorRate      float64 `json:"indicatorRate"`
	JailedUntil        uint64  `json:"jailedUntil"`
}

type RPCPeerInfo struct {
	NodeInfo         *NodeInfo `json:"node_info"`
	IsOutbound       bool      `json:"is_outbound"`
	ConnectionStatus struct {
		Duration uint64 `json:"Duration"`
	} `json:"connection_status"`
	RemoteIP string `json:"remote_ip"`
}

type PeerInfo struct {
	Duration uint64 `json:"Duration,omitempty" bson:"duration"`
	Moniker  string `json:"moniker,omitempty" bson:"moniker"` // arbitrary moniker
	RemoteIP string `json:"remote_ip,omitempty" bson:"remoteIp"`
}

type ProtocolVersion struct {
	P2P   uint64 `json:"p2p"`
	Block uint64 `json:"block"`
	App   uint64 `json:"app"`
}

type DefaultNodeInfoOther struct {
	TxIndex    string `json:"tx_index" bson:"txIndex"`
	RPCAddress string `json:"rpc_address" bson:"rpcAddress"`
}

type NodeInfo struct {
	ProtocolVersion ProtocolVersion      `json:"protocol_version" bson:"protocolVersion"`
	ID              string               `json:"id" bson:"id"`                  // authenticated identifier
	ListenAddr      string               `json:"listen_addr" bson:"listenAddr"` // accepting incoming
	Network         string               `json:"network" bson:"network"`        // network/chain ID
	Version         string               `json:"version" bson:"version"`        // major.minor.revision
	Moniker         string               `json:"moniker" bson:"moniker"`        // arbitrary moniker
	Peers           map[string]*PeerInfo `json:"peers,omitempty" bson:"peers"`  // peers details
	Other           DefaultNodeInfoOther `json:"other" bson:"other"`            // other application specific data
}

type ValidatorsByDelegator struct {
	Name                    string            `json:"name"`
	Validator               common.Address    `json:"validator"`
	ValidatorContractAddr   common.Address    `json:"validatorContractAddr"`
	ValidatorStatus         uint8             `json:"validatorStatus"`
	ValidatorRole           int               `json:"validatorRole"`
	StakedAmount            string            `json:"stakedAmount"`
	ClaimableRewards        string            `json:"claimableRewards"`
	UnbondedRecords         []*UnbondedRecord `json:"unbondedRecords"`
	TotalWithdrawableAmount string            `json:"totalWithdrawableAmount"`
	TotalUnbondedAmount     string            `json:"totalUnbondedAmount"`
	UnbondedAmount          string            `json:"unbondedAmount"`
	WithdrawableAmount      string            `json:"withdrawableAmount"`
}
