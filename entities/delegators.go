package entities

type Delegator struct {
	ID           int    `json:"id"`
	Address      string `json:"address"`
	StakedAmount string `json:"staked_amount"`
	UpdatedAt    int64  `json:"updated_at"`
	Balance      string `json:"balance"`
}

func (e *Delegator) Fields() ([]string, []interface{}) {
	return []string{"id", "address", "staked_amount", "updated_at", "balance"},
		[]interface{}{&e.ID, &e.Address, &e.StakedAmount, &e.UpdatedAt, &e.Balance}
}

func (e *Delegator) TableName() string {
	return `delegators`
}
