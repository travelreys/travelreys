package trips

import (
	context "context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/travelreys/travelreys/pkg/auth"
	"github.com/travelreys/travelreys/pkg/reqctx"
	"github.com/travelreys/travelreys/pkg/storage"
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

func (mw rbacMiddleware) Create(ctx context.Context, creator Member, name string, start, end time.Time) (Trip, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return Trip{}, ErrRBAC
	}
	return mw.next.Create(ctx, creator, name, start, end)
}

func (mw rbacMiddleware) Read(ctx context.Context, ID string) (Trip, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return Trip{}, ErrRBAC
	}
	return mw.next.Read(ctx, ID)
}

func (mw rbacMiddleware) ReadWithMembers(ctx context.Context, ID string) (Trip, auth.UsersMap, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return Trip{}, nil, ErrRBAC
	}
	return mw.next.ReadWithMembers(ctx, ID)
}

func (mw rbacMiddleware) ReadMembers(ctx context.Context, ID string) (auth.UsersMap, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return nil, ErrRBAC
	}
	return mw.next.ReadMembers(ctx, ID)
}

func (mw rbacMiddleware) List(ctx context.Context, ff ListFilter) (TripsList, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return nil, ErrRBAC
	}
	return mw.next.List(ctx, ff)
}

func (mw rbacMiddleware) Delete(ctx context.Context, ID string) error {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return ErrRBAC
	}
	return mw.next.Delete(ctx, ID)
}

func (mw rbacMiddleware) Upload(ctx context.Context, ID string, filename string, filesize int64, mimeType, attachmentType string, file io.Reader) error {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return ErrRBAC
	}
	return mw.next.Upload(ctx, ID, filename, filesize, mimeType, attachmentType, file)
}

func (mw rbacMiddleware) Download(ctx context.Context, ID string, obj storage.Object) (storage.Object, io.ReadCloser, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return storage.Object{}, nil, ErrRBAC
	}
	return mw.next.Download(ctx, ID, obj)
}

func (mw rbacMiddleware) DeleteFile(ctx context.Context, ID string, obj storage.Object) error {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return ErrRBAC
	}
	fmt.Println(obj)
	return mw.next.DeleteFile(ctx, ID, obj)
}
