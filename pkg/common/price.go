package common

const (
	PriceSplitMethodSolo       = "solo"
	PriceSplitMethodAbsolute   = "absolute"
	PriceSplitMethodPercentage = "percentage"
)

type PriceSplitTarget struct {
	ID      string  `json:"id" bson:"id"`
	Value   float64 `json:"value" bson:"value"`
	Settled bool    `json:"settled" bson:"settled"`
}

type PriceSplitOptions struct {
	Method  string                      `json:"method" bson:"method"`
	Targets map[string]PriceSplitTarget `json:"targets" bson:"targets"`
}

type Price struct {
	Amount   float64 `json:"amount" bson:"amount"`
	Currency string  `json:"currency" bson:"currency"`
}

type PriceItem struct {
	Price        `bson:"inline"`
	SplitOptions PriceSplitOptions `json:"splitOptions" bson:"splitOptions"`
}
