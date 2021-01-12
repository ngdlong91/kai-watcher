// Package types
package types

type TokenStats struct {
	Name                     string  `json:"-" bson:"name"`
	TokenName                string  `json:"name" bson:"tokenName"`
	Symbol                   string  `json:"symbol" bson:"symbol"`
	Decimal                  int64   `json:"decimal" bson:"decimal"`
	ERC20CirculatingSupply   int64   `json:"erc20_circulating_supply" bson:"erc20CirculatingSupply"`
	MainnetCirculatingSupply int64   `json:"mainnet_circulating_supply" bson:"mainnetCirculatingSupply"`
	TotalSupply              int64   `json:"total_supply" bson:"totalSupply"`
	Price                    float64 `json:"price" bson:"price"`
	Volume24h                float64 `json:"volume_24h" bson:"volume24H"`
	Change1h                 float64 `json:"change_1h" bson:"change1H"`
	Change24h                float64 `json:"change_24h" bson:"change24H"`
	Change7d                 float64 `json:"change_7d" bson:"change7D"`
	MarketCap                float64 `json:"market_cap" bson:"marketCap"`
}

type SupplyInfo struct {
	ERC20CirculatingSupply int64 `json:"erc20CirculatingSupply"`
	MainnetGenesisAmount   int64 `json:"mainnetGenesisAmount"`
}

type TxStats struct {
	NumTxs uint64 `json:"numTxs"`
	Time   uint64 `json:"time"`
}
