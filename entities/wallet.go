package entities

type Wallet struct {
	ID        int64
	Address   string
	Name      string
	UpdatedAt int64
}

func (e *Wallet) Fields() ([]string, []interface{}) {
	return []string{"id", "address", "name", "updated_at"},
		[]interface{}{&e.ID, &e.Address, &e.Name, &e.UpdatedAt}
}

func (e *Wallet) TableName() string {
	return "wallets"
}
