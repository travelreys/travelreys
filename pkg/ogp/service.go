package ogp

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
	Fetch(context.Context, string) (Opengraph, error)
}

type service struct{}

func NewService() Service {
	return service{}
}

func (svc service) Fetch(ctx context.Context, queryUrl string) (Opengraph, error) {
	intent := opengraph.Intent{
		Context:     ctx,
		HTTPClient:  &http.Client{Timeout: httpTimeout},
		Strict:      true,
		TrustedTags: []string{"meta", "title"},
	}

	graph, err := opengraph.Fetch(queryUrl, intent)
	if err != nil {
		return Opengraph{}, err
	}

	return OpengraphFromRawGraph(graph), nil
}
