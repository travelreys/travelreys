package flights

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/tiinyplanet/tiinyplanet/pkg/common"
	"github.com/tiinyplanet/tiinyplanet/pkg/utils"

	"github.com/gorilla/mux"
)

var (
	notFoundErrors = []error{}
	appErrors      = []error{ErrInvalidSearchRequest}
	authErrors     = []error{}
)

var (
	encodeErrFn = utils.EncodeErrorFactory(utils.ErrorToHTTPCodeFactory(notFoundErrors, appErrors, authErrors))

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
	searchFlightsHandler := kithttp.NewServer(NewSearchEndpoint(svc), decodeSearchRequest, encodeResponse, opts...)
	r.Handle("/api/v1/flights/search", searchFlightsHandler).Methods(http.MethodGet)
	return r
}

func decodeSearchRequest(_ context.Context, r *http.Request) (interface{}, error) {
	q := r.URL.Query()
	opts := FlightsSearchOptions{}

	if q.Has("returnDate") {
		opts["returnDate"] = q.Get("returnDate")
	}
	if q.Has("cabinClass") {
		opts["cabinClass"] = q.Get("cabinClass")
	}
	if q.Has("currency") {
		opts["currency"] = q.Get("currency")
	}
	if q.Has("duration") {
		opts["duration"] = q.Get("duration")
	}
	if q.Has("stops") {
		opts["stop"] = q.Get("stops")
	}

	numAdults, err := strconv.Atoi(q.Get("numAdults"))
	if err != nil {
		numAdults = 1
	}

	departDate, err := time.Parse("2006-01-02", q.Get("departDate"))
	if err != nil {
		departDate = time.Time{}
	}

	return SearchRequest{
		origIATA:   q.Get("origIATA"),
		destIATA:   q.Get("destIATA"),
		numAdults:  uint64(numAdults),
		departDate: departDate,
		opts:       opts,
	}, nil
}
