package images

import (
	"context"
	"encoding/json"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/tiinyplanet/tiinyplanet/pkg/common"
	"github.com/tiinyplanet/tiinyplanet/pkg/utils"

	"github.com/gorilla/mux"
)

var (
	notFoundErrors     = []error{}
	appErrors          = []error{ErrEmptyQuery}
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
	searchImageHandler := kithttp.NewServer(NewSearchEndpoint(svc), decodeSearchRequest, encodeResponse, opts...)
	r.Handle("/api/v1/images/search", searchImageHandler).Methods(http.MethodGet)
	return r
}

func decodeSearchRequest(_ context.Context, r *http.Request) (interface{}, error) {
	q := r.URL.Query().Get("query")
	return SearchRequest{Query: q}, nil
}
