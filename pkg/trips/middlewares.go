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
	ErrRBAC = errors.New("trips.ErrRBAC")
)

// validationMiddleware validates that the input for the service calls are acceptable
type validationMiddleware struct {
	next   Service
	logger *zap.Logger
}

func SvcWithValidationMw(svc Service, logger *zap.Logger) Service {
	return &validationMiddleware{svc, logger.Named("trips.validationMiddleware")}
}

func (mw validationMiddleware) Create(
	ctx context.Context,
	creatorID,
	name string,
	start,
	end time.Time,
) (*Trip, error) {
	if creatorID == "" {
		return nil, common.ErrValidation
	}
	return mw.next.Create(ctx, creatorID, name, start, end)
}

func (mw validationMiddleware) Save(ctx context.Context, trip *Trip) error {
	return common.ErrValidation
}

func (mw validationMiddleware) Read(ctx context.Context, ID string) (*Trip, error) {
	if ID == "" {
		return nil, common.ErrValidation
	}
	return mw.next.Read(ctx, ID)
}

func (mw validationMiddleware) ReadOGP(ctx context.Context, ID string) (TripOGP, error) {
	if ID == "" {
		return TripOGP{}, common.ErrValidation
	}
	return mw.next.ReadOGP(ctx, ID)
}

func (mw validationMiddleware) ReadMembers(ctx context.Context, ID string) (MembersMap, error) {
	if ID == "" {
		return nil, common.ErrValidation
	}
	return mw.next.ReadMembers(ctx, ID)
}

func (mw validationMiddleware) List(ctx context.Context, ff ListFilter) (TripsList, error) {
	if err := ff.Validate(); err != nil {
		return nil, common.ErrValidation
	}
	return mw.next.List(ctx, ff)
}

func (mw validationMiddleware) ListWithMembers(
	ctx context.Context,
	ff ListFilter,
) (TripsList, auth.UsersMap, error) {
	if err := ff.Validate(); err != nil {
		return nil, nil, common.ErrValidation
	}
	return mw.next.ListWithMembers(ctx, ff)
}

func (mw validationMiddleware) Delete(ctx context.Context, ID string) error {
	if ID == "" {
		return common.ErrValidation
	}
	return mw.next.Delete(ctx, ID)
}

func (mw validationMiddleware) DeleteAttachment(ctx context.Context, ID string, obj storage.Object) error {
	if ID == "" || obj.Path == "" || obj.Name == "" {
		return common.ErrValidation
	}
	return mw.next.DeleteAttachment(ctx, ID, obj)
}

func (mw validationMiddleware) DownloadAttachmentPresignedURL(
	ctx context.Context,
	ID,
	path,
	filename string,
) (string, error) {
	if ID == "" || path == "" || filename == "" {
		return "", common.ErrValidation
	}
	return mw.next.DownloadAttachmentPresignedURL(ctx, ID, path, filename)
}

func (mw validationMiddleware) UploadAttachmentPresignedURL(
	ctx context.Context,
	ID,
	filename,
	fileType string,
) (string, error) {
	if ID == "" || filename == "" {
		return "", common.ErrValidation
	}
	return mw.next.UploadAttachmentPresignedURL(ctx, ID, filename, fileType)
}

func (mw validationMiddleware) GenerateMediaItems(
	ctx context.Context,
	ID,
	userID string,
	params []media.NewMediaItemParams,
) (media.MediaItemList, media.MediaPresignedUrlList, error) {
	if ID == "" || userID == "" {
		return nil, nil, common.ErrValidation
	}
	for _, p := range params {
		if p.Name == "" {
			return nil, nil, common.ErrValidation
		}
	}
	return mw.next.GenerateMediaItems(ctx, ID, userID, params)
}

func (mw validationMiddleware) SaveMediaItems(ctx context.Context, ID string, items media.MediaItemList) error {
	if ID == "" {
		return common.ErrValidation
	}
	for _, item := range items {
		if item.Name == "" || item.Path == "" {
			return common.ErrValidation
		}
	}
	return mw.next.SaveMediaItems(ctx, ID, items)
}

func (mw validationMiddleware) DeleteMediaItems(
	ctx context.Context,
	ID string,
	items media.MediaItemList,
) error {
	if ID == "" {
		return common.ErrValidation
	}
	for _, item := range items {
		if item.Name == "" || item.Path == "" {
			return common.ErrValidation
		}
	}
	return mw.next.DeleteMediaItems(ctx, ID, items)
}

func (mw validationMiddleware) GenerateGetSignedURLs(
	ctx context.Context,
	ID string,
	items media.MediaItemList,
) (media.MediaPresignedUrlList, error) {
	if ID == "" {
		return nil, common.ErrValidation
	}
	for _, item := range items {
		if item.Name == "" || item.Path == "" {
			return nil, common.ErrValidation
		}
	}
	return mw.next.GenerateGetSignedURLs(ctx, ID, items)

}

type rbacMiddleware struct {
	next   Service
	logger *zap.Logger
}

func SvcWithRBACMw(
	svc Service,
	logger *zap.Logger,
) Service {
	return &rbacMiddleware{svc, logger}
}

func (mw rbacMiddleware) Create(
	ctx context.Context,
	creatorID,
	name string,
	start,
	end time.Time,
) (*Trip, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return nil, ErrRBAC
	}
	return mw.next.Create(ctx, creatorID, name, start, end)
}

func (mw rbacMiddleware) Save(ctx context.Context, trip *Trip) error {
	return ErrRBAC
}

func (mw rbacMiddleware) Read(ctx context.Context, ID string) (*Trip, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return nil, ErrRBAC
	}

	trip, err := mw.next.Read(ctx, ID)
	if err != nil {
		return nil, err
	}
	if !common.StringContains(trip.GetMemberIDs(), ci.UserID) {
		return nil, ErrRBAC
	}
	return trip, nil
}

func (mw rbacMiddleware) ReadOGP(ctx context.Context, ID string) (TripOGP, error) {
	return mw.next.ReadOGP(ctx, ID)
}

func (mw rbacMiddleware) ReadMembers(ctx context.Context, ID string) (MembersMap, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return nil, ErrRBAC
	}
	trip, err := mw.next.Read(ctx, ID)
	if err != nil {
		return nil, err
	}
	if !common.StringContains(trip.GetMemberIDs(), ci.UserID) {
		return nil, ErrRBAC
	}
	return mw.next.ReadMembers(ContextWithTripInfo(ctx, trip), ID)
}

func (mw rbacMiddleware) List(ctx context.Context, ff ListFilter) (TripsList, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return nil, ErrRBAC
	}
	return mw.next.List(ctx, ff)
}

func (mw rbacMiddleware) ListWithMembers(
	ctx context.Context,
	ff ListFilter,
) (TripsList, auth.UsersMap, error) {
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

	trip, err := mw.next.Read(ctx, ID)
	if err != nil && err != ErrTripNotFound {
		return err
	}
	if trip.Creator.ID != ci.UserID {
		return ErrRBAC
	}
	return mw.next.Delete(ContextWithTripInfo(ctx, trip), ID)
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

func (mw rbacMiddleware) UploadAttachmentPresignedURL(ctx context.Context, ID, filename, fileType string) (string, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return "", ErrRBAC
	}
	return mw.next.UploadAttachmentPresignedURL(ctx, ID, filename, fileType)
}

func (mw rbacMiddleware) GenerateMediaItems(ctx context.Context, ID, userID string, params []media.NewMediaItemParams) (media.MediaItemList, media.MediaPresignedUrlList, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return nil, nil, ErrRBAC
	}
	if ci.UserID != userID {
		return nil, nil, ErrRBAC
	}
	return mw.next.GenerateMediaItems(ctx, ID, userID, params)
}

func (mw rbacMiddleware) SaveMediaItems(ctx context.Context, ID string, items media.MediaItemList) error {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return ErrRBAC
	}
	return mw.next.SaveMediaItems(ctx, ID, items)
}

func (mw rbacMiddleware) DeleteMediaItems(ctx context.Context, ID string, items media.MediaItemList) error {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return ErrRBAC
	}
	return mw.next.DeleteMediaItems(ctx, ID, items)
}

func (mw rbacMiddleware) GenerateGetSignedURLs(ctx context.Context, ID string, items media.MediaItemList) (media.MediaPresignedUrlList, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return nil, ErrRBAC
	}
	return mw.next.GenerateGetSignedURLs(ctx, ID, items)

}
