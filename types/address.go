package types

type Address struct {
	Hash         string  `json:"hash" bson:"hash"`
	Name         string  `json:"name" bson:"name"`
	Balance      string  `json:"balance" bson:"balance"`
	BalanceFloat float64 `bson:"balanceFloat"`

	Rank uint64 `json:"rank" bson:"rank"`

	// SMC
	IsContract   bool   `json:"isContract" bson:"isContract"`
	OwnerAddress string `json:"ownerAddress" bson:"ownerAddress"`
}

func (o *Address) SetBalanceInFloat(b float64) {
	o.BalanceFloat = b
}

type SMC struct {
	Hash            string   `json:"hash" bson:"hash"`
	OwnerAddress    string   `json:"ownerAddress" bson:"ownerAddress"`
	ErcTypes        []string `json:"ercTypes" bson:"ercTypes"`
	Interfaces      []string `json:"interfaces" bson:"interfaces"`
	TxCount         int      `json:"txCount" bson:"txCount"`
	HolderCount     int      `json:"holderCount" bson:"holderCount"`
	InternalTxCount int      `json:"internalTxCount" bson:"internalTxCount"`
	TokenTxCount    int      `json:"tokenTxCount" bson:"tokenTxCount"`
}
