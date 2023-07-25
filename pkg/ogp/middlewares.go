package ogp

import (
	"context"
	"errors"
	"net"

	"github.com/travelreys/travelreys/pkg/reqctx"
	"go.uber.org/zap"
)

var (
	ErrValidation = errors.New("ogp.ErrValidation")
	ErrRBAC       = errors.New("ogp.ErrRBAC")
)

type validationMiddleware struct {
	next   Service
	logger *zap.Logger
}

func SvcWithValidation(svc Service, logger *zap.Logger) Service {
	return &validationMiddleware{svc, logger.Named("ogp.validationMiddleware")}
}

func (mw validationMiddleware) Fetch(ctx context.Context, url string) (Opengraph, error) {
	if url == "" {
		return Opengraph{}, ErrValidation
	}
	ips, err := net.LookupIP(url)
	if err != nil {
		return Opengraph{}, ErrValidation
	}
	for _, ip := range ips {
		if ip.IsPrivate() {
			return Opengraph{}, ErrValidation
		}
	}
	return mw.next.Fetch(ctx, url)
}

type rbacMiddleware struct {
	next   Service
	logger *zap.Logger
}

func SvcWithRBACMw(svc Service, logger *zap.Logger) Service {
	return &rbacMiddleware{svc, logger.Named("ogp.rbacMiddleware")}
}

func (mw rbacMiddleware) Fetch(ctx context.Context, url string) (Opengraph, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return Opengraph{}, ErrRBAC
	}
	return mw.next.Fetch(ctx, url)
}
