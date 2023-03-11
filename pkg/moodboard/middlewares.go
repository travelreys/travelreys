package moodboard

import (
	"context"
	"errors"

	"github.com/travelreys/travelreys/pkg/reqctx"
	"go.uber.org/zap"
)

var (
	ErrRBAC = errors.New("auth.rbac.error")
)

type rbacMiddleware struct {
	next   Service
	logger *zap.Logger
}

func ServiceWithRBACMiddleware(svc Service, logger *zap.Logger) Service {
	return &rbacMiddleware{svc, logger}
}

func (mw rbacMiddleware) ReadAndCreateIfNotExists(ctx context.Context, id string) (Moodboard, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return Moodboard{}, ErrRBAC
	}
	return mw.next.ReadAndCreateIfNotExists(ctx, ci.UserID)
}

func (mw rbacMiddleware) Update(ctx context.Context, id, title string) error {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return ErrRBAC
	}
	return mw.next.Update(ctx, ci.UserID, title)
}

func (mw rbacMiddleware) AddPin(ctx context.Context, id string, url string) (string, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return "", ErrRBAC
	}
	return mw.next.AddPin(ctx, ci.UserID, url)
}

func (mw rbacMiddleware) UpdatePin(ctx context.Context, id, pinID string, notes string) error {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return ErrRBAC
	}
	return mw.next.UpdatePin(ctx, id, pinID, notes)
}

func (mw rbacMiddleware) DeletePin(ctx context.Context, id, pinID string) error {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return ErrRBAC
	}
	return mw.next.DeletePin(ctx, id, pinID)
}
