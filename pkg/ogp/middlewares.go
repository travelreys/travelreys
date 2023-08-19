package ogp

import (
	"context"
	"errors"
	"net"
	"net/url"

	"github.com/travelreys/travelreys/pkg/common"
	"github.com/travelreys/travelreys/pkg/reqctx"
	"go.uber.org/zap"
)

var (
	ErrRBAC = errors.New("ogp.ErrRBAC")
)

type validationMiddleware struct {
	next   Service
	logger *zap.Logger
}

func SvcWithValidation(svc Service, logger *zap.Logger) Service {
	return &validationMiddleware{svc, logger.Named("ogp.validationMiddleware")}
}

func (mw validationMiddleware) Fetch(ctx context.Context, queryURL string) (Opengraph, error) {
	if queryURL == "" {
		return Opengraph{}, common.ErrValidation
	}
	u, err := url.Parse(queryURL)
	if err != nil {
		return Opengraph{}, common.ErrValidation
	}

	ips, err := net.LookupIP(u.Host)
	if err != nil {
		return Opengraph{}, common.ErrValidation
	}
	for _, ip := range ips {
		if ip.IsPrivate() {
			return Opengraph{}, common.ErrValidation
		}
	}
	return mw.next.Fetch(ctx, queryURL)
}

type rbacMiddleware struct {
	next   Service
	logger *zap.Logger
}

func SvcWithRBACMw(svc Service, logger *zap.Logger) Service {
	return &rbacMiddleware{svc, logger.Named("ogp.rbacMiddleware")}
}

func (mw rbacMiddleware) Fetch(ctx context.Context, queryURL string) (Opengraph, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return Opengraph{}, ErrRBAC
	}
	return mw.next.Fetch(ctx, queryURL)
}
