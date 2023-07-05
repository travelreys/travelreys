package finance

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
	"go.uber.org/zap"
)

const (
	ExchangeRatesPrefix = "finance:exchange-rates"
)

var (
	ErrRatesNotFound = errors.New("finance.ErrRatesNotFound")
)

type Store interface {
	ReadLatestFxRates(context.Context, string) (ExchangeRates, error)
	SaveFxRates(context.Context, ExchangeRates, time.Duration) error
}

type store struct {
	rdb    redis.UniversalClient
	logger *zap.Logger
}

func NewStore(rdb redis.UniversalClient, logger *zap.Logger) Store {
	return &store{rdb, logger}
}

func (s store) ReadLatestFxRates(ctx context.Context, base string) (ExchangeRates, error) {
	key := fmt.Sprintf("%s:%s", ExchangeRatesPrefix, base)
	cmd := s.rdb.Get(ctx, key)
	if errors.Is(cmd.Err(), redis.Nil) {
		return ExchangeRates{}, ErrRatesNotFound
	}

	var rates ExchangeRates
	if err := json.Unmarshal([]byte(cmd.Val()), &rates); err != nil {
		s.rdb.Del(ctx, key)
		return ExchangeRates{}, err
	}
	return rates, nil
}

func (s store) SaveFxRates(ctx context.Context, rates ExchangeRates, ttl time.Duration) error {
	key := fmt.Sprintf("%s:%s", ExchangeRatesPrefix, rates.Base)
	data, err := json.Marshal(rates)
	if err != nil {
		return err
	}
	cmd := s.rdb.Set(ctx, key, string(data), ttl)
	return cmd.Err()
}
