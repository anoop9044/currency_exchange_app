package model

type ExchangeRate struct {
	Currency string  `json:"currency"`
	Rate     float64 `json:"rate"`
	UpdateTime *int64 `json:"updateTime,omitempty"`
}