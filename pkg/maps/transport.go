package maps

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/tiinyplanet/tiinyplanet/pkg/common"
	"github.com/tiinyplanet/tiinyplanet/pkg/utils"
)

var (
	notFoundErrors     = []error{}
	appErrors          = []error{ErrInvalidField, ErrInvalidSessionToken}
	unauthorisedErrors = []error{}
)

var (
	encodeErrFn = utils.EncodeErrorFactory(utils.ErrorToHTTPCodeFactory(notFoundErrors, appErrors, unauthorisedErrors))

	opts = []kithttp.ServerOption{
		// kithttp.ServerBefore(reqctx.MakeContextFromHTTPRequest),
		kithttp.ServerErrorEncoder(encodeErrFn),
	}
)

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(common.Errorer); ok && e.Error() != nil {
		encodeErrFn(ctx, e.Error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func MakeHandler(svc Service) http.Handler {
	r := mux.NewRouter()
	placeAutocompleteHandler := kithttp.NewServer(NewPlacesAutocompleteEndpoint(svc), decodePlaceAutocompleteRequest, encodeResponse, opts...)
	placeDetailsHandler := kithttp.NewServer(NewPlaceDetailsEndpoint(svc), decodePlaceDetailsRequest, encodeResponse, opts...)

	r.Handle("/api/v1/maps/place/autocomplete", placeAutocompleteHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/maps/place/details", placeDetailsHandler).Methods(http.MethodGet)
	return r
}

func decodePlaceAutocompleteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	q := r.URL.Query()
	return PlacesAutocompleteRequest{
		Query:        q.Get("query"),
		Sessiontoken: q.Get("sessiontoken"),
		Types:        q.Get("types"),
	}, nil

}

func decodePlaceDetailsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	q := r.URL.Query()
	return PlaceDetailsRequest{
		PlaceID:      q.Get("placeID"),
		Sessiontoken: q.Get("sessiontoken"),
		Fields:       strings.Split(q.Get("fields"), ","),
	}, nil
}
