package maps

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/travelreys/travelreys/pkg/common"
	"github.com/travelreys/travelreys/pkg/reqctx"
)

func errToHttpCode() func(err error) int {
	notFoundErrors := []error{}
	appErrors := []error{ErrInvalidField, ErrInvalidSessionToken}

	return func(err error) int {
		if common.ErrorContains(notFoundErrors, err) {
			return http.StatusNotFound
		}
		if common.ErrorContains(appErrors, err) {
			return http.StatusUnprocessableEntity
		}
		if errors.Is(err, ErrRBAC) {
			return http.StatusUnauthorized
		}
		return http.StatusInternalServerError
	}
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(common.Errorer); ok && e.Error() != nil {
		common.EncodeErrorFactory(errToHttpCode())(ctx, e.Error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func MakeHandler(svc Service) http.Handler {
	r := mux.NewRouter()

	opts := []kithttp.ServerOption{
		kithttp.ServerBefore(reqctx.ContextWithClientInfo),
		kithttp.ServerErrorEncoder(common.EncodeErrorFactory(errToHttpCode())),
	}

	placeAutocompleteHandler := kithttp.NewServer(NewPlacesAutocompleteEndpoint(svc), decodePlaceAutocompleteRequest, encodeResponse, opts...)
	placeDetailsHandler := kithttp.NewServer(NewPlaceDetailsEndpoint(svc), decodePlaceDetailsRequest, encodeResponse, opts...)
	placeAtmosphereHandler := kithttp.NewServer(NewPlaceAtmosphereEndpoint(svc), decodePlaceAtmosphereRequest, encodeResponse, opts...)
	directionsHandler := kithttp.NewServer(NewDirectionsEndpoint(svc), decodeDirectionRequest, encodeResponse, opts...)

	r.Handle("/api/v1/maps/place/autocomplete", placeAutocompleteHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/maps/place/details", placeDetailsHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/maps/place/atmosphere", placeAtmosphereHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/maps/place/directions", directionsHandler).Methods(http.MethodGet)

	return r
}

func decodePlaceAutocompleteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	q := r.URL.Query()
	return PlacesAutocompleteRequest{
		Query:        q.Get("query"),
		Sessiontoken: q.Get("sessiontoken"),
		Types:        q.Get("types"),
		Lang:         q.Get("language"),
	}, nil
}

func decodePlaceDetailsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	q := r.URL.Query()
	return PlaceDetailsRequest{
		PlaceID:      q.Get("placeID"),
		Sessiontoken: q.Get("sessiontoken"),
		Fields:       strings.Split(q.Get("fields"), ","),
		Lang:         q.Get("language"),
	}, nil
}

func decodePlaceAtmosphereRequest(_ context.Context, r *http.Request) (interface{}, error) {
	q := r.URL.Query()
	return PlaceAtmosphereRequest{
		PlaceID:      q.Get("placeID"),
		Sessiontoken: q.Get("sessiontoken"),
		Fields:       strings.Split(q.Get("fields"), ","),
		Lang:         q.Get("language"),
	}, nil
}

func decodeDirectionRequest(_ context.Context, r *http.Request) (interface{}, error) {
	q := r.URL.Query()
	return DirectionsRequest{
		OriginPlaceID: q.Get("originPlaceID"),
		DestPlaceID:   q.Get("destPlaceID"),
		Modes:         strings.Split(q.Get("modes"), ","),
	}, nil
}
