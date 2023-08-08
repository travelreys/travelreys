package finance

const (
	PriceSplitMethodSolo       = "solo"
	PriceSplitMethodAbsolute   = "absolute"
	PriceSplitMethodPercentage = "percentage"
)

type PriceSplitTarget struct {
	ID      string  `json:"id" bson:"id" msgpack:"id"`
	Value   float64 `json:"value" bson:"value" msgpack:"value"`
	Settled bool    `json:"settled" bson:"settled" msgpack:"settled"`
}

type PriceSplitOptions struct {
	Method  string                      `json:"method" bson:"method" msgpack:"method"`
	Targets map[string]PriceSplitTarget `json:"targets" bson:"targets" msgpack:"targets"`
}

type Price struct {
	Amount   float64 `json:"amount" bson:"amount" msgpack:"amount"`
	Currency string  `json:"currency" bson:"currency" msgpack:"currency"`
}

type PriceItem struct {
	Price        `bson:"inline" msgpack:"inline"`
	SplitOptions PriceSplitOptions `json:"splitOptions" bson:"splitOptions" msgpack:"splitOptions"`
}
