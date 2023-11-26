package finance

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
)

// 1. https://cdn.jsdelivr.net/gh/fawazahmed0/currency-api@1/2023-05-22/currencies/sgd.json
// 2. https://exchangerate.host/#/#docs

var (
	FxApiEndpoints = []string{
		"http://api.exchangerate.host",
	}

	FxApiKey = os.Getenv("TRAVELREYS_EXCHANGE_RATE_KEY")
)

type Service interface {
	GetFxRates(ctx context.Context, base string) (ExchangeRates, error)
}

type service struct {
	store  Store
	logger *zap.Logger
}

func NewService(store Store, logger *zap.Logger) Service {
	return &service{store, logger}
}

func (svc service) GetFxRates(ctx context.Context, base string) (ExchangeRates, error) {
	rates, err := svc.store.ReadLatestFxRates(ctx, base)
	if err == nil {
		svc.logger.Info("GetFxRates", zap.String("cache hit", base))
		return rates, nil
	}

	if err != ErrRatesNotFound {
		svc.logger.Error("GetFxRates", zap.Error(err))
	}

	client := http.Client{}
	request, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/live?access_key=%s&source=%s",
			FxApiEndpoints[rand.Intn(len(FxApiEndpoints))],
			FxApiKey,
			base,
		),
		nil,
	)
	if err != nil {
		svc.logger.Error("GetFxRates", zap.Error(err))
		return ExchangeRates{}, err
	}

	resp, err := client.Do(request)
	if err != nil {
		svc.logger.Error("GetFxRates", zap.Error(err))
		return ExchangeRates{}, err
	}

	var exRateResp ExchangeRateHostResponse
	if err := json.NewDecoder(resp.Body).Decode(&exRateResp); err != nil {
		svc.logger.Error("GetFxRates", zap.Error(err))
		return ExchangeRates{}, err
	}

	rates = NewExchangeRatesFromExchangeRateHostResponse(exRateResp)
	svc.store.SaveFxRates(ctx, rates, 60*time.Minute)
	return rates, nil
}
