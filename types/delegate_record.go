package types

type DelegateRecord struct {
	Hash string
	From string
	To   string
}

func (d DelegateRecord) TableName() string {
	return "delegate_records"
}

func (d DelegateRecord) Fields() ([]string, []interface{}) {
	return []string{
			"hash", "from", "to",
		}, []interface{}{
			&d.Hash, &d.From, &d.To,
		}
}
