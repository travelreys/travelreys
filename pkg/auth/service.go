package auth

import (
	"context"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/travelreys/travelreys/pkg/common"
	"github.com/travelreys/travelreys/pkg/storage"
	"go.uber.org/zap"
)

const (
	SvcLoggerName      = "auth.service"
	authCookieDuration = 365 * 24 * time.Hour
)

var (
	avatarFilePrefix = "avatar"
	avatarBucket     = os.Getenv("TRAVELREYS_MEDIA_BUCKET")
	mediaBucket      = os.Getenv("TRAVELREYS_MEDIA_BUCKET")
	mediaCDNDomain   = os.Getenv("TRAVELREYS_MEDIA_DOMAIN") // cdn.travelreys.com

	ErrProviderGoogleError   = errors.New("auth.service.google.error")
	ErrProviderFacebookError = errors.New("auth.service.facebook.error")
	ErrProviderNotSupported  = errors.New("auth.service.provider.notsupported")
)

type Service interface {
	Signup(ctx context.Context, authCode string) (User, *http.Cookie, error)
	Login(context.Context, string, string) (User, *http.Cookie, error)
	Read(context.Context, string) (User, error)
	List(ctx context.Context, ff ListFilter) (UsersList, error)
	Update(context.Context, string, UpdateFilter) error
	Delete(context.Context, string) error

	UploadAvatarPresignedURL(context.Context, string) (string, error)
	GenerateMediaPresignedCookie(ctx context.Context) (*http.Cookie, error)
}

type service struct {
	google       GoogleProvider
	fb           FacebookProvider
	emailpw      *EmailPasswordProvider
	store        Store
	secureCookie bool

	storageSvc storage.Service
	logger     *zap.Logger
}

func NewService(
	gp GoogleProvider,
	fb FacebookProvider,
	emailpw *EmailPasswordProvider,
	store Store,
	secureCookie bool,
	storageSvc storage.Service,
	logger *zap.Logger,
) Service {
	return &service{
		google:       gp,
		fb:           fb,
		emailpw:      emailpw,
		store:        store,
		secureCookie: secureCookie,
		storageSvc:   storageSvc,
		logger:       logger.Named(SvcLoggerName),
	}
}

func (svc service) Signup(ctx context.Context, authCode string) (User, *http.Cookie, error) {
	usr, err := svc.emailpw.Signup(ctx, authCode)
	if err != nil {
		return usr, nil, err
	}

	jwt, err := svc.issueJwtToken(usr)
	if err != nil {
		svc.logger.Error("Signup", zap.Error(err))
		return User{}, nil, err
	}

	if err := svc.store.Save(ctx, usr); err != nil {
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

func (svc service) Login(ctx context.Context, authCode, provider string) (User, *http.Cookie, error) {
	var (
		usr     User
		googUsr GoogleUser
		fbUsr   FacebookUser
		err     error
	)

	if provider == OIDCProviderGoogle {
		googUsr, err := svc.google.TokenToUserInfo(ctx, authCode)
		if err != nil {
			svc.logger.Error("Login", zap.String("provider", provider), zap.Error(err))
			return User{}, nil, err
		}
		usr = UserFromGoogleUser(googUsr)
	} else if provider == OIDCProviderFacebook {
		fbUsr, err := svc.fb.TokenToUserInfo(ctx, authCode)
		if err != nil {
			svc.logger.Error("Login", zap.String("provider", provider), zap.Error(err))
			return User{}, nil, err
		}
		usr = UserFromFBUser(fbUsr)
	} else if provider == OIDCProviderEmailPassword {
		emailUsr, err := svc.emailpw.TokenToUserInfo(ctx, authCode)
		if err != nil {
			svc.logger.Error("Login", zap.String("provider", provider), zap.Error(err))
			return User{}, nil, err
		}
		usr = emailUsr
	} else {
		return User{}, nil, ErrProviderNotSupported
	}

	if provider == OIDCProviderGoogle || provider == OIDCProviderFacebook {
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

		if provider == OIDCProviderGoogle {
			googUsr.AddLabelsToUser(&existUsr)
		} else if provider == OIDCProviderFacebook {
			fbUsr.AddLabelsToUser(&existUsr)
		}

		if err := svc.store.Save(ctx, existUsr); err != nil {
			return User{}, nil, err
		}
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

func (svc service) Delete(ctx context.Context, ID string) error {
	ff := UpdateFilter{
		Email:       common.StringPtr(""),
		Name:        common.StringPtr(""),
		PhoneNumber: &PhoneNumber{},
		Labels:      &common.Labels{},
	}
	return svc.store.Update(ctx, ID, ff)
}

func (svc service) UploadAvatarPresignedURL(ctx context.Context, ID string) (string, error) {
	return svc.storageSvc.PutPresignedURL(
		ctx,
		avatarBucket,
		filepath.Join(avatarFilePrefix, ID),
		ID,
	)
}

func (svc service) GenerateMediaPresignedCookie(ctx context.Context) (*http.Cookie, error) {
	return svc.storageSvc.GeneratePresignedCookie(ctx, mediaCDNDomain, mediaBucket)
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
