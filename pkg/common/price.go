package common

type PriceWithCurrency struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}
