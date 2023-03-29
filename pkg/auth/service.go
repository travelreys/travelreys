package auth

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/travelreys/travelreys/pkg/common"
	"go.uber.org/zap"
)

const (
	SvcLoggerName      = "auth.service"
	authCookieDuration = 365 * 24 * time.Hour
)

var (
	ErrProviderGoogleError  = errors.New("auth.service.google.error")
	ErrProviderNotSupported = errors.New("auth.service.provider.notsupported")
)

type Service interface {
	Login(context.Context, string, string) (User, *http.Cookie, error)
	Read(context.Context, string) (User, error)
	List(ctx context.Context, ff ListFilter) (UsersList, error)
	Update(context.Context, string, UpdateFilter) error
}

type service struct {
	google       GoogleProvider
	store        Store
	secureCookie bool

	logger *zap.Logger
}

func NewService(gp GoogleProvider, store Store, secureCookie bool, logger *zap.Logger) Service {
	return &service{
		google:       gp,
		store:        store,
		secureCookie: secureCookie,
		logger:       logger.Named(SvcLoggerName),
	}
}

func (svc service) googleLogin(ctx context.Context, authCode string) (User, error) {
	gusr, err := svc.google.TokenToUserInfo(ctx, authCode)
	if err != nil {
		svc.logger.Error("google login failed", zap.Error(err))
		return User{}, ErrProviderGoogleError
	}
	usr := UserFromGoogleUser(gusr)
	return usr, nil
}

func (svc service) Login(ctx context.Context, authCode, provider string) (User, *http.Cookie, error) {
	var (
		usr User
		err error
	)
	if provider == OIDCProviderGoogle {
		usr, err = svc.googleLogin(ctx, authCode)
		if err != nil {
			svc.logger.Error("Login", zap.String("provider", provider), zap.Error(err))
			return User{}, nil, err
		}
	} else {
		return User{}, nil, ErrProviderNotSupported
	}

	existUsr, err := svc.store.Read(ctx, ReadFilter{Email: usr.Email})
	if err == ErrUserNotFound {
		usr.CreatedAt = time.Now()
		if err := svc.createUser(ctx, usr); err != nil {
			return User{}, nil, err
		}
	} else if err != nil {
		return User{}, nil, err
	} else {
		usr = existUsr
	}

	jwt, err := svc.issueJwtToken(usr)
	if err != nil {
		svc.logger.Error("Login", zap.String("provider", provider), zap.Error(err))
		return User{}, nil, err
	}
	cookie := &http.Cookie{
		Name:     AccessCookieName,
		Value:    jwt,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   int(authCookieDuration.Seconds()),
	}
	if svc.secureCookie {
		cookie.SameSite = http.SameSiteNoneMode
		cookie.Secure = true
	}

	return usr, cookie, nil
}

func (svc service) Read(ctx context.Context, ID string) (User, error) {
	return svc.store.Read(ctx, ReadFilter{ID: ID})
}

func (svc service) Update(ctx context.Context, ID string, ff UpdateFilter) error {
	return svc.store.Update(ctx, ID, ff)
}

func (svc service) List(ctx context.Context, ff ListFilter) (UsersList, error) {
	return svc.store.List(ctx, ff)
}

func (svc service) createUser(ctx context.Context, usr User) error {
	return svc.store.Save(ctx, usr)
}

func (svc service) issueJwtToken(usr User) (string, error) {
	token := jwt.NewWithClaims(common.JWTDefaultSigningMethod, jwt.MapClaims{
		common.JwtClaimIss:   common.JwtIssuer,
		common.JwtClaimSub:   usr.ID,
		common.JwtClaimEmail: usr.Email,
		common.JwtClaimIat:   time.Now().Unix(),
	})
	return token.SignedString(common.GetJwtSecretBytes())
}
