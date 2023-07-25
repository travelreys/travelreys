package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/travelreys/travelreys/pkg/common"
	"github.com/travelreys/travelreys/pkg/reqctx"
	"go.uber.org/zap"
)

var (
	ErrValidation = errors.New("auth.ErrValidation")
	ErrRBAC       = errors.New("auth.ErrRBAC")
)

// validationMiddleware validates that the input for the service calls are acceptable
type validationMiddleware struct {
	next   Service
	logger *zap.Logger
}

func SvcWithValidationMw(svc Service, logger *zap.Logger) Service {
	return &validationMiddleware{svc, logger.Named("auth.validationMiddleware")}
}

func (mw validationMiddleware) Login(
	ctx context.Context,
	authCode,
	signature,
	provider string,
) (User, *http.Cookie, error) {
	if provider == "" || authCode == "" {
		return User{}, nil, ErrValidation
	}
	if provider == OIDCProviderOTP && signature == "" {
		return User{}, nil, ErrValidation
	}
	return mw.next.Login(ctx, authCode, signature, provider)
}

func (mw validationMiddleware) MagicLink(ctx context.Context, email string) error {
	if !common.IsEmailValid(email) {
		return ErrProviderOTPInvalidEmail
	}
	return mw.next.MagicLink(ctx, email)
}

func (mw validationMiddleware) Read(ctx context.Context, ID string) (User, error) {
	if ID == "" {
		return User{}, ErrValidation
	}
	return mw.next.Read(ctx, ID)
}

func (mw validationMiddleware) List(ctx context.Context, ff ListFilter) (UsersList, error) {
	if err := ff.Validate(); err != nil {
		return nil, ErrValidation
	}
	return mw.next.List(ctx, ff)
}

func (mw validationMiddleware) Update(ctx context.Context, ID string, ff UpdateFilter) error {
	if ID == "" {
		return ErrValidation
	}
	if err := ff.Validate(); err != nil {
		return ErrValidation
	}
	return mw.next.Update(ctx, ID, ff)
}

func (mw validationMiddleware) Delete(ctx context.Context, ID string) error {
	if ID == "" {
		return ErrValidation
	}
	return mw.next.Delete(ctx, ID)
}

func (mw validationMiddleware) UploadAvatarPresignedURL(ctx context.Context, ID, mimeType string) (string, string, error) {
	if ID == "" || mimeType == "" {
		return "", "", ErrValidation
	}
	return mw.next.UploadAvatarPresignedURL(ctx, ID, mimeType)
}

type rbacMiddleware struct {
	next   Service
	logger *zap.Logger
}

func SvcWithRBACMw(svc Service, logger *zap.Logger) Service {
	return &rbacMiddleware{svc, logger.Named("auth.rbacMiddleware")}
}

func (mw rbacMiddleware) Login(ctx context.Context, authCode, signature, provider string) (User, *http.Cookie, error) {
	return mw.next.Login(ctx, authCode, signature, provider)
}

func (mw rbacMiddleware) MagicLink(ctx context.Context, email string) error {
	return mw.next.MagicLink(ctx, email)
}

func (mw rbacMiddleware) Read(ctx context.Context, ID string) (User, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return User{}, ErrRBAC
	}
	return mw.next.Read(ctx, ID)
}

func (mw rbacMiddleware) List(ctx context.Context, ff ListFilter) (UsersList, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil || ci.HasEmptyID() {
		return UsersList{}, ErrRBAC
	}
	if ff.Email != nil || len(ff.IDs) > 0 {
		return UsersList{}, ErrRBAC
	}
	return mw.next.List(ctx, ff)
}

func (mw rbacMiddleware) Update(ctx context.Context, ID string, ff UpdateFilter) error {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil {
		return ErrRBAC
	}
	if ci.UserID != ID {
		return ErrRBAC
	}
	return mw.next.Update(ctx, ID, ff)
}

func (mw rbacMiddleware) Delete(ctx context.Context, ID string) error {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil {
		return ErrRBAC
	}
	if ci.UserID != ID {
		return ErrRBAC
	}
	return mw.next.Delete(ctx, ID)
}

func (mw rbacMiddleware) UploadAvatarPresignedURL(ctx context.Context, ID, mimeType string) (string, string, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil {
		return "", "", ErrRBAC
	}
	if ci.UserID != ID {
		return "", "", ErrRBAC
	}
	return mw.next.UploadAvatarPresignedURL(ctx, ID, mimeType)
}
