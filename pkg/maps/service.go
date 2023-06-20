package maps

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"googlemaps.github.io/maps"
)

const (
	mapSerivceLogger      = "maps.service"
	envGoogleMapsApiToken = "TRAVELREYS_GOOGLE_MAPS_APIKEY"
)

var (
	ErrInvalidSessionToken = errors.New("maps.service.invalidsessiontoken")
	ErrInvalidField        = errors.New("maps.service.invalidfield")
)

type Service interface {
	PlacesAutocomplete(context.Context, string, string, string, string) (AutocompletePredictionList, error)
	PlaceDetails(context.Context, string, []string, string, string) (Place, error)
	PlaceAtmosphere(ctx context.Context, placeID string, fields []string, sessiontoken, lang string) (PlaceAtmosphere, error)
	Directions(ctx context.Context, originPlaceID, destPlaceID string, modes []string) (RouteList, error)
	OptimizeRoute(ctx context.Context, originPlaceID, destPlaceID string, waypointsPlaceID []string) (RouteList, []int, error)
}

type service struct {
	c      *maps.Client
	logger *zap.Logger
}

func GetApiToken() string {
	return os.Getenv(envGoogleMapsApiToken)
}

func NewDefaulService(logger *zap.Logger) (Service, error) {
	return NewService(GetApiToken(), logger)
}

func NewService(apiToken string, logger *zap.Logger) (Service, error) {
	c, err := maps.NewClient(maps.WithAPIKey(apiToken))
	if err != nil {
		return nil, err
	}
	return &service{c, logger.Named(mapSerivceLogger)}, nil
}

func (svc *service) PlacesAutocomplete(ctx context.Context, query, types, sessiontoken, lang string) (AutocompletePredictionList, error) {
	stuuid, err := uuid.Parse(sessiontoken)
	if err != nil {
		return AutocompletePredictionList{}, ErrInvalidSessionToken
	}
	req := &maps.PlaceAutocompleteRequest{
		Input:        query,
		Types:        maps.AutocompletePlaceType(types),
		SessionToken: maps.PlaceAutocompleteSessionToken(stuuid),
		Language:     lang,
	}
	res, err := svc.c.PlaceAutocomplete(ctx, req)
	if err != nil {
		svc.logger.Error("PlacesAutocomplete",
			zap.String("query", query),
			zap.String("types", types),
			zap.String("sessiontoken", sessiontoken),
			zap.Error(err),
		)
		return AutocompletePredictionList{}, err
	}
	preds := AutocompletePredictionList{}
	for _, ap := range res.Predictions {
		preds = append(preds, AutocompletePrediction{ap})
	}
	return preds, nil
}

func (svc *service) stringsToPlaceDefaultsFieldMasks(fields []string) ([]maps.PlaceDetailsFieldMask, error) {
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

func (svc *service) placeDetails(ctx context.Context, placeID string, fields []string, sessiontoken, lang string) (maps.PlaceDetailsResult, error) {
	fieldMasks, err := svc.stringsToPlaceDefaultsFieldMasks(fields)
	if err != nil {
		return maps.PlaceDetailsResult{}, err
	}

	req := &maps.PlaceDetailsRequest{
		PlaceID:  placeID,
		Fields:   fieldMasks,
		Language: lang,
	}
	if sessiontoken != "" {
		stuuid, err := uuid.Parse(sessiontoken)
		if err != nil {
			svc.logger.Error("PlaceDetails",
				zap.String("placeID", placeID),
				zap.String("fields", strings.Join(fields, ",")),
				zap.String("sessiontoken", sessiontoken),
				zap.Error(err),
			)
			return maps.PlaceDetailsResult{}, err
		}
		req.SessionToken = maps.PlaceAutocompleteSessionToken(stuuid)
	}
	res, err := svc.c.PlaceDetails(ctx, req)
	return res, err
}

func (svc *service) PlaceDetails(ctx context.Context, placeID string, fields []string, sessiontoken, lang string) (Place, error) {
	result, err := svc.placeDetails(ctx, placeID, fields, sessiontoken, lang)
	if err != nil {
		return Place{}, err
	}
	return PlaceFromPlaceDetailsResult(result), nil
}

func (svc *service) PlaceAtmosphere(ctx context.Context, placeID string, fields []string, sessiontoken, lang string) (PlaceAtmosphere, error) {
	result, err := svc.placeDetails(ctx, placeID, fields, sessiontoken, lang)
	if err != nil {
		return PlaceAtmosphere{}, err
	}
	return PlaceAtmosphere{result}, err
}

func (svc *service) Directions(
	ctx context.Context,
	originPlaceID,
	destPlaceID string,
	modes []string,
) (RouteList, error) {
	routes := RouteList{}
	for _, mode := range modes {
		req := &maps.DirectionsRequest{
			Origin:      fmt.Sprintf("place_id:%s", originPlaceID),
			Destination: fmt.Sprintf("place_id:%s", destPlaceID),
			Mode:        maps.Mode(mode),
		}

		groutes, _, err := svc.c.Directions(ctx, req)
		if err != nil {
			svc.logger.Error("Directions",
				zap.String("originPlaceID", originPlaceID),
				zap.String("destPlaceID", destPlaceID),
				zap.String("modes", strings.Join(modes, ",")),
				zap.Error(err),
			)
			return RouteList{}, err
		}
		if len(groutes) <= 0 {
			continue
		}
		route, _ := GetShortestDurationGMapsRoute(groutes)
		routes = append(routes, RouteFromRouteResult(route, mode))
	}
	return routes, nil
}

func (svc *service) OptimizeRoute(ctx context.Context, originPlaceID, destPlaceID string, waypointsPlaceID []string) (RouteList, []int, error) {
	wpWithLabel := []string{}
	for _, wp := range waypointsPlaceID {
		wpWithLabel = append(wpWithLabel, fmt.Sprintf("place_id:%s", wp))
	}

	req := &maps.DirectionsRequest{
		Origin:      fmt.Sprintf("place_id:%s", originPlaceID),
		Destination: fmt.Sprintf("place_id:%s", destPlaceID),
		Mode:        maps.TravelModeWalking,
		Waypoints:   wpWithLabel,
		Optimize:    true,
	}

	routes := RouteList{}
	groutes, _, err := svc.c.Directions(ctx, req)
	if err != nil {
		svc.logger.Error("OptimizeRoute",
			zap.String("originPlaceID", originPlaceID),
			zap.String("destPlaceID", destPlaceID),
			zap.String("waypointsPlaceID", strings.Join(waypointsPlaceID, ",")),
			zap.Error(err),
		)
		return routes, nil, err
	}

	if len(groutes) <= 0 {
		return routes, []int{}, nil
	}

	for _, r := range groutes {
		routes = append(routes, RouteFromRouteResult(r, ""))
	}
	return routes, groutes[0].WaypointOrder, err
}
