package moodboard

import (
	"context"
	"net/http"
	"time"

	"github.com/otiai10/opengraph/v2"
)

const (
	httpTimeout = 5 * time.Second
)

type Service interface {
	AddBasePin(context.Context, string) (string, error)
	UpdatePin(context.Context, string, string) error
	DeletePin(context.Context, string) error
}

type service struct {
	store Store
}

func NewService(store Store) Service {
	return &service{store}
}

func (s service) AddBasePin(ctx context.Context, url string) (string, error) {
	intent := opengraph.Intent{
		Context:     ctx,
		HTTPClient:  &http.Client{Timeout: httpTimeout},
		Strict:      true,
		TrustedTags: []string{"meta", "title"},
	}
	ogp, err := opengraph.Fetch("https://ogp.me", intent)
	if err != nil {
		return "", err
	}
	return "", nil

}

func (s service) UpdatePin(ctx context.Context, id string, notes string) error {
	return nil
}

func (s service) DeletePin(ctx context.Context, id string) error {
	return nil
}
