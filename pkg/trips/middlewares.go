package trips

import (
	context "context"
	"errors"
	"time"

	"github.com/travelreys/travelreys/pkg/auth"
	"github.com/travelreys/travelreys/pkg/common"
	"github.com/travelreys/travelreys/pkg/media"
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

func ServiceWithRBACMiddleware(
	svc Service,
	logger *zap.Logger,
) Service {
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

	trip, err := mw.next.Read(ctx, ID)
	if err != nil {
		return Trip{}, err
	}
	if !common.StringContains(trip.GetAllMembersID(), ci.UserID) {
		return Trip{}, ErrRBAC
	}
	return trip, nil
}

func (mw rbacMiddleware) ReadOGP(ctx context.Context, ID string) (TripOGP, error) {
	return mw.next.ReadOGP(ctx, ID)
}

func (mw rbacMiddleware) ReadWithMembers(ctx context.Context, ID string) (Trip, auth.UsersMap, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return Trip{}, nil, ErrRBAC
	}
	trip, err := mw.next.Read(ctx, ID)
	if err != nil {
		return Trip{}, nil, err
	}
	if !common.StringContains(trip.GetAllMembersID(), ci.UserID) {
		return Trip{}, nil, ErrRBAC
	}
	ctxWithTripInfo := ContextWithTripInfo(ctx, trip)
	return mw.next.ReadWithMembers(ctxWithTripInfo, ID)
}

func (mw rbacMiddleware) ReadMembers(ctx context.Context, ID string) (auth.UsersMap, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return nil, ErrRBAC
	}
	trip, err := mw.next.Read(ctx, ID)
	if err != nil {
		return nil, err
	}
	if !common.StringContains(trip.GetAllMembersID(), ci.UserID) {
		return nil, ErrRBAC
	}
	ctxWithTripInfo := ContextWithTripInfo(ctx, trip)
	return mw.next.ReadMembers(ctxWithTripInfo, ID)
}

func (mw rbacMiddleware) List(ctx context.Context, ff ListFilter) (TripsList, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return nil, ErrRBAC
	}
	return mw.next.List(ctx, ff)
}

func (mw rbacMiddleware) ListWithMembers(ctx context.Context, ff ListFilter) (TripsList, auth.UsersMap, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return nil, nil, ErrRBAC
	}
	return mw.next.ListWithMembers(ctx, ff)
}

func (mw rbacMiddleware) Delete(ctx context.Context, ID string) error {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return ErrRBAC
	}
	return mw.next.Delete(ctx, ID)
}

func (mw rbacMiddleware) DeleteAttachment(ctx context.Context, ID string, obj storage.Object) error {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return ErrRBAC
	}
	return mw.next.DeleteAttachment(ctx, ID, obj)
}

func (mw rbacMiddleware) DownloadAttachmentPresignedURL(ctx context.Context, ID, path, filename string) (string, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return "", ErrRBAC
	}
	return mw.next.DownloadAttachmentPresignedURL(ctx, ID, path, filename)
}

func (mw rbacMiddleware) UploadAttachmentPresignedURL(ctx context.Context, ID, filename string) (string, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return "", ErrRBAC
	}
	return mw.next.UploadAttachmentPresignedURL(ctx, ID, filename)
}

func (mw rbacMiddleware) GenerateMediaItems(ctx context.Context, id, userID string, params []media.NewMediaItemParams) (media.MediaItemList, media.MediaPresignedUrlList, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return nil, nil, ErrRBAC
	}
	if ci.UserID != userID {
		return nil, nil, ErrRBAC
	}
	return mw.next.GenerateMediaItems(ctx, id, userID, params)
}

func (mw rbacMiddleware) SaveMediaItems(ctx context.Context, id string, items media.MediaItemList) error {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return ErrRBAC
	}
	return mw.next.SaveMediaItems(ctx, id, items)
}

func (mw rbacMiddleware) DeleteMediaItems(ctx context.Context, id string, items media.MediaItemList) error {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return ErrRBAC
	}
	return mw.next.DeleteMediaItems(ctx, id, items)
}

func (mw rbacMiddleware) GenerateGetSignedURLs(ctx context.Context, ID string, items media.MediaItemList) (media.MediaPresignedUrlList, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return nil, ErrRBAC
	}
	return mw.next.GenerateGetSignedURLs(ctx, ID, items)

}
