package auth

import (
	"context"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/travelreys/travelreys/pkg/common"
	"github.com/travelreys/travelreys/pkg/email"
	"github.com/travelreys/travelreys/pkg/storage"
	"go.uber.org/zap"
)

const (
	authCookieDuration = 365 * 24 * time.Hour
	otpDuration        = 60 * time.Second
	defaultLoginSender = "login@travelreys.com"

	svcLoggerName = "auth.service"
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
	Login(context.Context, string, string) (User, *http.Cookie, error)
	MagicLink(ctx context.Context, email string) error
	Read(context.Context, string) (User, error)
	List(ctx context.Context, ff ListFilter) (UsersList, error)
	Update(context.Context, string, UpdateFilter) error
	Delete(context.Context, string) error

	UploadAvatarPresignedURL(context.Context, string) (string, error)
	GenerateMediaPresignedCookie(ctx context.Context) (*http.Cookie, error)
}

type service struct {
	google GoogleProvider
	fb     FacebookProvider
	otp    *OTPProvider

	store        Store
	secureCookie bool

	mailSvc    email.Service
	storageSvc storage.Service
	logger     *zap.Logger
}

func NewService(
	gp GoogleProvider,
	fb FacebookProvider,
	otp *OTPProvider,
	store Store,
	secureCookie bool,
	mailSvc email.Service,
	storageSvc storage.Service,
	logger *zap.Logger,
) Service {
	return &service{
		google:       gp,
		fb:           fb,
		otp:          otp,
		store:        store,
		secureCookie: secureCookie,
		mailSvc:      mailSvc,
		storageSvc:   storageSvc,
		logger:       logger.Named(svcLoggerName),
	}
}

func (svc service) Login(ctx context.Context, authCode, provider string) (User, *http.Cookie, error) {
	var (
		usr User
		err error
	)
	switch provider {
	case OIDCProviderFacebook, OIDCProviderGoogle:
		usr, err = svc.socialLogin(ctx, authCode, provider)
	case OIDCProviderOTP:
		usr, err = svc.otp.TokenToUserInfo(ctx, authCode)
	default:
		return User{}, nil, ErrProviderNotSupported
	}

	if err != nil {
		svc.logger.Error("Login", zap.String("provider", provider), zap.Error(err))
		return User{}, nil, err
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

func (svc service) socialLogin(ctx context.Context, authCode, provider string) (User, error) {
	var (
		usr     User
		googUsr GoogleUser
		fbUsr   FacebookUser
		err     error
	)

	if provider == OIDCProviderGoogle {
		googUsr, err = svc.google.TokenToUserInfo(ctx, authCode)
		if err != nil {
			return User{}, err
		}
		usr = UserFromGoogleUser(googUsr)
	} else if provider == OIDCProviderFacebook {
		fbUsr, err = svc.fb.TokenToUserInfo(ctx, authCode)
		if err != nil {
			return User{}, err
		}
		usr = UserFromFBUser(fbUsr)
	}

	existUsr, err := svc.store.Read(ctx, ReadFilter{Email: usr.Email})
	if err == ErrUserNotFound {
		usr.CreatedAt = time.Now()
		if err := svc.store.Save(ctx, usr); err != nil {
			return User{}, err
		}
	} else if err != nil {
		return User{}, err
	} else {
		usr = existUsr
	}

	if provider == OIDCProviderGoogle {
		googUsr.AddLabelsToUser(&existUsr)
	} else if provider == OIDCProviderFacebook {
		fbUsr.AddLabelsToUser(&existUsr)
	}

	if err := svc.store.Save(ctx, existUsr); err != nil {
		return User{}, err
	}
	return existUsr, nil
}

func (svc service) MagicLink(ctx context.Context, email string) error {
	var (
		usr User
		err error
	)
	usr, err = svc.store.Read(ctx, ReadFilter{Email: email})
	if err == ErrUserNotFound {
		newUsr, err := svc.createUser(ctx, email)
		if err != nil {
			return err
		}
		usr = newUsr
	}

	otp, hashedOTP, err := svc.otp.GenerateOTP(6)
	if err != nil {
		return err
	}
	if err = svc.store.SaveOTP(ctx, usr.ID, hashedOTP, otpDuration); err != nil {
		return err
	}
	go func() {
		mailBody, _ := svc.otp.GenerateMagicLinkEmail(usr, otp)
		loginSubj := "Login to Travelreys!"
		if err := svc.mailSvc.SendMail(ctx, email, defaultLoginSender, loginSubj, mailBody); err != nil {
			svc.logger.Error("send mail", zap.Error(err))
		}
	}()
	return nil
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

func (svc service) issueJwtToken(usr User) (string, error) {
	token := jwt.NewWithClaims(common.JWTDefaultSigningMethod, jwt.MapClaims{
		common.JwtClaimIss:   common.JwtIssuer,
		common.JwtClaimSub:   usr.ID,
		common.JwtClaimEmail: usr.Email,
		common.JwtClaimIat:   time.Now().Unix(),
	})
	return token.SignedString(common.GetJwtSecretBytes())
}

func (svc service) createUser(ctx context.Context, email string) (User, error) {
	if _, err := svc.store.Read(ctx, ReadFilter{Email: email}); err == nil {
		return User{}, ErrProviderOTPEmailExists
	}
	newusr := User{
		ID:          uuid.NewString(),
		Email:       email,
		Name:        email,
		CreatedAt:   time.Now(),
		PhoneNumber: PhoneNumber{},
		Labels: common.Labels{
			LabelAvatarImage:   "https://cdn.travelreys.com/travelreys-media-demo/avatar/account.png",
			LabelDefaultLocale: "en",
		},
	}
	if err := svc.store.Save(ctx, newusr); err != nil {
		return User{}, err
	}
	return newusr, nil
}
