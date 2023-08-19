package auth

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/lucasepe/codename"
	"github.com/travelreys/travelreys/pkg/common"
	"github.com/travelreys/travelreys/pkg/email"
	"github.com/travelreys/travelreys/pkg/storage"
	"go.uber.org/zap"
)

const (
	authCookieDuration = 365 * 24 * time.Hour
	otpDuration        = 120 * time.Second

	loginMailSender         = "login@travelreys.com"
	welcomEmailTmplFilePath = "assets/welcomeEmail.tmpl.html"
	welcomEmailTmplFileName = "welcomeEmail.tmpl.html"
)

var (
	avatarFilePrefix = "avatar"
	avatarBucket     = os.Getenv("TRAVELREYS_PUBLIC_BUCKET")

	ErrProviderGoogleError   = errors.New("auth.ErrProviderGoogleError")
	ErrProviderFacebookError = errors.New("auth.ErrProviderFacebookError")
	ErrProviderNotSupported  = errors.New("auth.ErrProviderNotSupported")
)

type Service interface {
	Login(ctx context.Context, authCode, signature, provider string) (User, *http.Cookie, error)
	MagicLink(ctx context.Context, email string) error

	Read(ctx context.Context, ID string) (User, error)
	List(ctx context.Context, ff ListFilter) (UsersList, error)
	Update(ctx context.Context, ID string, ff UpdateFilter) error
	Delete(ctx context.Context, ID string) error

	UploadAvatarPresignedURL(context.Context, string, string) (string, string, error)
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
		logger:       logger.Named("auth.service"),
	}
}

func (svc service) Login(
	ctx context.Context,
	authCode,
	signature,
	provider string,
) (User, *http.Cookie, error) {
	var (
		usr     User
		googUsr GoogleUser
		fbUsr   FacebookUser
		err     error
	)
	switch provider {
	case OIDCProviderGoogle:
		googUsr, err = svc.google.TokenToUserInfo(ctx, authCode)
		if err != nil {
			return User{}, nil, err
		}
		usr = UserFromGoogleUser(googUsr)
	case OIDCProviderFacebook:
		fbUsr, err = svc.fb.TokenToUserInfo(ctx, authCode)
		if err != nil {
			return User{}, nil, err
		}
		usr = UserFromFBUser(fbUsr)
	case OIDCProviderOTP:
		usr, err = svc.otp.TokenToUserInfo(ctx, authCode, signature)
		if err != nil {
			return User{}, nil, err
		}
	default:
		return User{}, nil, ErrProviderNotSupported
	}

	isNewUser := false
	existUsr, err := svc.store.Read(ctx, ReadFilter{Email: usr.Email})
	if err == ErrUserNotFound {
		isNewUser = true
	} else if err != nil {
		return User{}, nil, err
	} else {
		usr = existUsr
	}

	if provider == OIDCProviderGoogle {
		googUsr.AddLabelsToUser(&usr)
	} else if provider == OIDCProviderFacebook {
		fbUsr.AddLabelsToUser(&usr)
	}
	if isNewUser {
		usr.CreatedAt = time.Now()
	}

	if err := svc.store.Save(ctx, usr); err != nil {
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

	if isNewUser {
		go svc.sendWelcomeEmail(ctx, usr.Name, usr.Email)
	}

	return usr, cookie, nil
}

func (svc service) MagicLink(ctx context.Context, email string) error {
	otp, hashedOTP, err := svc.otp.GenerateOTP(6)
	if err != nil {
		return err
	}
	if err = svc.store.SaveOTP(ctx, email, hashedOTP, otpDuration); err != nil {
		return err
	}
	go func() {
		svc.sendMagicLinkEmail(ctx, otp, email)
	}()
	return nil
}

func (svc service) Read(ctx context.Context, ID string) (User, error) {
	return svc.store.Read(ctx, ReadFilter{ID: ID})
}

func (svc service) Update(ctx context.Context, ID string, ff UpdateFilter) error {
	origUser, err := svc.store.Read(ctx, ReadFilter{ID: ID})
	if err != nil {
		return err
	}
	if err := svc.store.Update(ctx, ID, ff); err != nil {
		return err
	}

	if ff.Labels != nil &&
		(*ff.Labels)[LabelAvatarImage] != origUser.GetAvatarImgURL() &&
		origUser.GetAvatarImgURL() != "" {
		svc.storageSvc.Remove(ctx, origUser.MakeUserAvatarObject())
	}
	return nil
}

func (svc service) List(ctx context.Context, ff ListFilter) (UsersList, error) {
	return svc.store.List(ctx, ff)
}

func (svc service) Delete(ctx context.Context, ID string) error {
	return svc.store.Delete(ctx, ID)
}

func (svc service) UploadAvatarPresignedURL(ctx context.Context, ID, mimeType string) (string, string, error) {
	rng, _ := codename.DefaultRNG()
	suffixToken := common.RandomToken(rng, 8)
	presignedURL, err := svc.storageSvc.PutPresignedURL(
		ctx,
		avatarBucket,
		filepath.Join(avatarFilePrefix, fmt.Sprintf("%s-%s", ID, suffixToken)),
		ID,
		mimeType,
	)
	return suffixToken, presignedURL, err
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

func (svc service) sendWelcomeEmail(ctx context.Context, name, to string) {
	svc.logger.Info("sending welcome email", zap.String("to", to))

	t, err := template.
		New(welcomEmailTmplFileName).
		ParseFiles(welcomEmailTmplFilePath)
	if err != nil {
		svc.logger.Error("sendWelcomeEmail", zap.Error(err))
		return
	}

	var doc bytes.Buffer
	data := struct {
		Name string
	}{name}
	if err := t.Execute(&doc, data); err != nil {
		svc.logger.Error("sendWelcomeEmail", zap.Error(err))
		return
	}

	mailBody, err := svc.mailSvc.InsertContentOnTemplate(doc.String())
	if err != nil {
		svc.logger.Error("sendWelcomeEmail", zap.Error(err))
		return
	}

	subj := "Welcome to Travelreys!"
	if err := svc.mailSvc.SendMail(ctx, to, loginMailSender, subj, mailBody); err != nil {
		svc.logger.Error("sendWelcomeEmail", zap.Error(err))
	}
}

func (svc service) sendMagicLinkEmail(
	ctx context.Context,
	otp,
	email string,
) {
	mailContentBody, err := svc.otp.GenerateMagicLinkEmail(email, otp)
	if err != nil {
		svc.logger.Error("sendMagicLinkEmail", zap.Error(err))
		return
	}
	mailBody, err := svc.mailSvc.InsertContentOnTemplate(mailContentBody)
	if err != nil {
		svc.logger.Error("sendMagicLinkEmail", zap.Error(err))
		return
	}

	loginSubj := "Login to Travelreys!"
	if err := svc.mailSvc.SendMail(
		ctx,
		email,
		loginMailSender,
		loginSubj,
		mailBody,
	); err != nil {
		svc.logger.Error("sendMagicLinkEmail", zap.Error(err))
	}

}
