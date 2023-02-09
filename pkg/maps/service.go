package maps

import (
	"context"
	"errors"
	"os"

	"github.com/google/uuid"
	"googlemaps.github.io/maps"
)

var (
	ErrInvalidSessionToken = errors.New("ErrInvalidSessionToken")
	ErrInvalidField        = errors.New("ErrInvalidField")
)

type Service interface {
	PlacesAutocomplete(context.Context, string, string, string) (AutocompletePredictionList, error)
	PlaceDetails(context.Context, string, []string, string) (Place, error)
}

type service struct {
	apiToken string
	c        *maps.Client
}

func NewService() (Service, error) {
	apiToken := os.Getenv("TIINYPLANET_GMAPS_APIKEY")

	c, err := maps.NewClient(maps.WithAPIKey(apiToken))
	if err != nil {
		return nil, err
	}

	return &service{apiToken, c}, nil
}

func (svc *service) stringArrayToPlaceDetailsFieldMasksArray(fields []string) ([]maps.PlaceDetailsFieldMask, error) {
	list := []maps.PlaceDetailsFieldMask{}
	for _, str := range fields {
		field, err := maps.ParsePlaceDetailsFieldMask(str)
		if err != nil {
			return list, ErrInvalidField
		}
		list = append(list, field)
	}
	return list, nil
}

func (svc *service) PlacesAutocomplete(ctx context.Context, query, types, sessiontoken string) (AutocompletePredictionList, error) {
	stuuid, err := uuid.Parse(sessiontoken)
	if err != nil {
		return AutocompletePredictionList{}, ErrInvalidSessionToken
	}
	req := &maps.PlaceAutocompleteRequest{
		Input:        query,
		Types:        maps.AutocompletePlaceType(types),
		SessionToken: maps.PlaceAutocompleteSessionToken(stuuid),
	}
	res, err := svc.c.PlaceAutocomplete(ctx, req)
	if err != nil {
		return AutocompletePredictionList{}, err
	}
	preds := AutocompletePredictionList{}
	for _, ap := range res.Predictions {
		preds = append(preds, AutocompletePrediction{ap})
	}
	return preds, nil

}

func (svc *service) PlaceDetails(ctx context.Context, placeID string, fields []string, sessiontoken string) (Place, error) {
	fieldMasks, err := svc.stringArrayToPlaceDetailsFieldMasksArray(fields)
	if err != nil {
		return Place{}, err
	}

	req := &maps.PlaceDetailsRequest{
		PlaceID: placeID,
		Fields:  fieldMasks,
	}
	if sessiontoken != "" {
		stuuid, err := uuid.Parse(sessiontoken)
		if err != nil {
			return Place{}, err
		}
		req.SessionToken = maps.PlaceAutocompleteSessionToken(stuuid)
	}

	res, err := svc.c.PlaceDetails(ctx, req)
	return Place{res}, err

}
