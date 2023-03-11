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
	ReadAndCreateIfNotExists(context.Context, string) (Moodboard, error)
	Update(context.Context, string, string) error
	AddPin(context.Context, string, string) (string, error)
	UpdatePin(context.Context, string, string, string) error
	DeletePin(context.Context, string, string) error
}

type service struct {
	store Store
}

func NewService(store Store) Service {
	return &service{store}
}

func (s service) ReadAndCreateIfNotExists(ctx context.Context, id string) (Moodboard, error) {
	mb, err := s.store.Read(ctx, id)
	if err == ErrMoodboardNotFound {
		mb = NewMoodboard(id)
		if err := s.store.Save(ctx, mb); err != nil {
			return Moodboard{}, err
		}
	}
	return mb, nil
}

func (s service) Update(ctx context.Context, id string, title string) error {
	_, err := s.store.Read(ctx, id)
	if err != nil {
		return err
	}
	return s.store.Update(ctx, id, UpdateBoardFilter{title})
}

func (s service) AddPin(ctx context.Context, id string, url string) (string, error) {
	_, err := s.store.Read(ctx, id)
	if err != nil {
		return "", err
	}

	intent := opengraph.Intent{
		Context:     ctx,
		HTTPClient:  &http.Client{Timeout: httpTimeout},
		Strict:      true,
		TrustedTags: []string{"meta", "title"},
	}
	ogp, err := opengraph.Fetch(url, intent)
	if err != nil {
		return "", err
	}
	pin := PinFromOGP(ogp)
	if err := s.store.SavePin(ctx, id, pin); err != nil {
		return "", err
	}
	return pin.ID, nil
}

func (s service) UpdatePin(ctx context.Context, id, pinID string, notes string) error {
	_, err := s.store.Read(ctx, id)
	if err != nil {
		return err
	}
	return s.store.UpdatePin(ctx, id, pinID, UpdatePinFilter{Notes: notes})
}

func (s service) DeletePin(ctx context.Context, id, pinID string) error {
	_, err := s.store.Read(ctx, id)
	if err != nil {
		return err
	}
	return s.store.DeletePin(ctx, id, pinID)
}
