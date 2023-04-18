package types

type APIValidator struct {
	Address               string          `json:"address" bson:"address,omitempty"`
	SmcAddress            string          `json:"smcAddress" bson:"smcAddress,omitempty"`
	Status                uint8           `json:"status" bson:"status"`
	Role                  int             `json:"role" bson:"role"`
	Jailed                bool            `json:"jailed" bson:"jailed"`
	Name                  string          `json:"name,omitempty" bson:"name,omitempty"`
	VotingPowerPercentage string          `json:"votingPowerPercentage" bson:"votingPowerPercentage,omitempty"`
	StakedAmount          string          `json:"stakedAmount" bson:"stakedAmount,omitempty"`
	AccumulatedCommission string          `json:"accumulatedCommission" bson:"accumulatedCommission,omitempty"`
	UpdateTime            uint64          `json:"updateTime" bson:"updateTime,omitempty"`
	CommissionRate        string          `json:"commissionRate" bson:"commissionRate,omitempty"`
	TotalDelegators       int             `json:"totalDelegators" bson:"totalDelegators,omitempty"`
	MaxRate               string          `json:"maxRate" bson:"maxRate,omitempty"`
	MaxChangeRate         string          `json:"maxChangeRate" bson:"maxChangeRate,omitempty"`
	SigningInfo           *APISigningInfo `json:"signingInfo" bson:"signingInfo,omitempty"`
	Delegators            []*APIDelegator `json:"delegators,omitempty" bson:"delegators,omitempty"`
}

type APIDelegator struct {
	ValidatorSMCAddress string `json:"validatorSMCAddress" bson:"validatorSMCAddress,omitempty"`
	Address             string `json:"address" bson:"address,omitempty"`
	StakedAmount        string `json:"stakedAmount" bson:"stakedAmount,omitempty"`
	Reward              string `json:"reward" bson:"reward,omitempty"`
}

type APISigningInfo struct {
	StartHeight        uint64  `json:"startHeight"`
	IndexOffset        uint64  `json:"indexOffset"`
	Tombstoned         bool    `json:"tombstoned"`
	MissedBlockCounter uint64  `json:"missedBlockCounter"`
	IndicatorRate      float64 `json:"indicatorRate"`
	JailedUntil        uint64  `json:"jailedUntil"`
}
