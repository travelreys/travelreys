package auth

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/travelreys/travelreys/pkg/common"
	"github.com/travelreys/travelreys/pkg/reqctx"
	"go.uber.org/zap"
)

var (
	ErrRBAC = errors.New("auth.ErrRBAC")
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
		mw.logger.Warn("Login")
		return User{}, nil, common.ErrValidation
	}
	if provider == OIDCProviderOTP && signature == "" {
		return User{}, nil, common.ErrValidation
	}
	return mw.next.Login(ctx, authCode, signature, provider)
}

func (mw validationMiddleware) MagicLink(ctx context.Context, email string) error {
	if !common.IsEmailValid(email) {
		mw.logger.Warn("MagicLink")
		return ErrProviderOTPInvalidEmail
	}
	return mw.next.MagicLink(ctx, email)
}

func (mw validationMiddleware) GenerateOTPAuthCodeAndSig(
	ctx context.Context,
	email string,
	duration time.Duration,
) (string, string, error) {
	if email == "" {
		mw.logger.Warn("GenerateOTPAuthCodeAndSig")
		return "", "", common.ErrValidation
	}
	return mw.next.GenerateOTPAuthCodeAndSig(ctx, email, duration)
}

func (mw validationMiddleware) Read(ctx context.Context, ID string) (User, error) {
	if ID == "" {
		mw.logger.Warn("Read")
		return User{}, common.ErrValidation
	}
	return mw.next.Read(ctx, ID)
}

func (mw validationMiddleware) List(ctx context.Context, ff ListFilter) (UsersList, error) {
	if err := ff.Validate(); err != nil {
		mw.logger.Warn("List")
		return nil, common.ErrValidation
	}
	return mw.next.List(ctx, ff)
}

func (mw validationMiddleware) Update(ctx context.Context, ID string, ff UpdateFilter) error {
	if ID == "" {
		mw.logger.Warn("Update")
		return common.ErrValidation
	}
	if err := ff.Validate(); err != nil {
		return common.ErrValidation
	}
	return mw.next.Update(ctx, ID, ff)
}

func (mw validationMiddleware) Delete(ctx context.Context, ID string) error {
	if ID == "" {
		mw.logger.Warn("Delete")
		return common.ErrValidation
	}
	return mw.next.Delete(ctx, ID)
}

func (mw validationMiddleware) UploadAvatarPresignedURL(ctx context.Context, ID, mimeType string) (string, string, error) {
	if ID == "" || mimeType == "" {
		mw.logger.Warn("UploadAvatarPresignedURL")
		return "", "", common.ErrValidation
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

func (mw rbacMiddleware) GenerateOTPAuthCodeAndSig(
	ctx context.Context,
	email string,
	duration time.Duration,
) (string, string, error) {
	return "", "", nil
}

func (mw rbacMiddleware) EmailLogin(
	ctx context.Context,
	authCode,
	signature string,
	isLoggedIn bool,
) (User, *http.Cookie, error) {
	return User{}, nil, ErrRBAC
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
