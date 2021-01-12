// Package types
package types

import (
	"github.com/kardiachain/go-kardia/lib/common"
)

type Validator struct {
	OwnerAddress       string `json:"ownerAddress" bson:"ownerAddress"`
	SmcAddress         string `json:"smcAddress" bson:"smcAddress"`
	Name               string `json:"name,omitempty" bson:"name"`
	IsJailed           bool   `json:"isJailed"`
	JailedUntil        uint64 `json:"jailedUntil" bson:"jailedUntil"`
	MissedBlockCounter uint64 `json:"missedBlockCounter" bson:"missedBlockCounter"`
}

type Delegator struct {
	Address      common.Address `json:"address" bson:"address"`
	Name         string         `json:"name,omitempty" bson:"name"`
	StakedAmount string         `json:"stakedAmount" bson:"stakedAmount"`
	Reward       string         `json:"reward" bson:"reward"`
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
	ProtocolVersion ProtocolVersion `json:"protocol_version" bson:"protocolVersion"`
	ID              string          `json:"id" bson:"id"`                  // authenticated identifier
	ListenAddr      string          `json:"listen_addr" bson:"listenAddr"` // accepting incoming
	Network         string          `json:"network" bson:"network"`        // network/chain ID
	Version         string          `json:"version" bson:"version"`        // major.minor.revision
	Moniker         string          `json:"moniker" bson:"moniker"`        // arbitrary moniker
	Peers           []*PeerInfo     `json:"peers,omitempty" bson:"peers"`  // peers details
}

type ValidatorsByDelegator struct {
	Name                  string         `json:"name"`
	Validator             common.Address `json:"validator"`
	ValidatorContractAddr common.Address `json:"validatorContractAddr"`
	ValidatorStatus       uint8          `json:"validatorStatus"`
	ValidatorRole         int            `json:"validatorRole"`
	StakedAmount          string         `json:"stakedAmount"`
	ClaimableRewards      string         `json:"claimableRewards"`
	UnbondedAmount        string         `json:"unbondedAmount"`
	WithdrawableAmount    string         `json:"withdrawableAmount"`
}
