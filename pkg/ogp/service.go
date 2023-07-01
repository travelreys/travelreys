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
		Strict:      false,
		TrustedTags: []string{"meta", "title"},
	}

	c := &http.Client{Timeout: httpTimeout}
	req, _ := http.NewRequest(http.MethodGet, queryUrl, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")

	resp, err := c.Do(req)
	if err != nil {
		return Opengraph{}, err
	}

	ogp := &opengraph.OpenGraph{Intent: intent}
	if err := ogp.Parse(resp.Body); err != nil {
		return Opengraph{}, err
	}

	return OpengraphFromRawGraph(ogp), nil
}
