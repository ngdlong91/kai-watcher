package types

const (
	defaultLimit = 50
	MaximumLimit = 100
)

type Pagination struct {
	Skip  int    `json:"skip"`
	Limit int    `json:"limit"`
	Total uint64 `json:"total"`
}

func (f *Pagination) Sanitize() {
	if f.Skip < 0 {
		f.Skip = 0
	}
	if f.Limit <= 0 {
		f.Limit = defaultLimit
	} else if f.Limit > MaximumLimit {
		f.Limit = MaximumLimit
	}
}
