package finance

import (
	"fmt"
	"strconv"
	"time"
)

/*
	{
	    "success": true,
	    "terms": "https://exchangerate.host/terms",
	    "privacy": "https://exchangerate.host/privacy",
	    "timestamp": 1430068515,
	    "source": "USD",
	    "quotes": {
	        "USDAUD": 1.278384,
	        "USDCHF": 0.953975,
	        "USDEUR": 0.919677,
	        "USDGBP": 0.658443,
	        "USDPLN": 3.713873
	    }
	}
*/
type ExchangeRateHostResponse struct {
	Timestamp int                `json:"timestamp"`
	Source    string             `json:"source"`
	Quotes    map[string]float64 `json:"quotes"`
}

type ExchangeRates struct {
	Date  string             `json:"date"`
	Base  string             `json:"base"`
	Rates map[string]float64 `json:"rates"`
}

func NewExchangeRatesFromExchangeRateHostResponse(resp ExchangeRateHostResponse) ExchangeRates {
	i, _ := strconv.ParseInt(fmt.Sprintf("%d", resp.Timestamp), 10, 64)
	tm := time.Unix(i, 0)
	exRate := ExchangeRates{
		Date:  tm.Format("2006-01-02"),
		Base:  resp.Source,
		Rates: map[string]float64{},
	}

	for pair, val := range resp.Quotes {
		exRate.Rates[pair[3:]] = val
	}

	exRate.Rates[resp.Source] = 1
	return exRate
}
