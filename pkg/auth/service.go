package auth

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/tiinyplanet/tiinyplanet/pkg/common"
)

type Service interface {
	Login(context.Context, string, string) (string, error)
	ReadUser(context.Context, string) (User, error)
	UpdateUser(context.Context, string, UpdateUserFilter) error
}

type service struct {
	google GoogleProvider
	store  Store
}

func NewService(gp GoogleProvider, store Store) Service {
	return &service{
		google: gp,
		store:  store,
	}
}

func (svc service) Login(ctx context.Context, authCode, provider string) (string, error) {
	var usr User
	if provider == OIDCProviderGoogle {
		tkn, err := svc.google.AuthCodeToToken(ctx, authCode)
		if err != nil {
			return "", err
		}
		gusr, err := svc.google.TokenToUserInfo(ctx, tkn)
		if err != nil {
			return "", err
		}
		usr = UserFromGoogleUser(gusr)
	} else {
		return "", errors.New("not-supported")
	}

	existUsr, err := svc.store.ReadUserByEmail(ctx, usr.Email)
	if err == ErrUserNotFound {
		if err := svc.createUser(ctx, usr); err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	} else {
		usr = existUsr
	}

	jwtTkn, err := svc.issueJWTToken(usr)
	if err != nil {
		return "", err
	}
	return jwtTkn, nil
}

func (svc service) ReadUser(ctx context.Context, ID string) (User, error) {
	return svc.store.ReadUserByID(ctx, ID)
}

func (svc service) UpdateUser(ctx context.Context, ID string, ff UpdateUserFilter) error {
	return svc.store.UpdateUser(ctx, ID, ff)
}

func (svc service) createUser(ctx context.Context, usr User) error {
	return svc.store.SaveUser(ctx, usr)
}

func (svc service) issueJWTToken(usr User) (string, error) {
	token := jwt.NewWithClaims(common.JWTDefaultSigningMethod, jwt.MapClaims{
		common.JWTClaimIss:   common.JWTIssuer,
		common.JWTClaimSub:   usr.ID,
		common.JWTClaimEmail: usr.Email,
		common.JWTClaimIat:   time.Now().Unix(),
	})

	jwtSecret := os.Getenv("TIINYPLANET_JWT_SECRET")
	return token.SignedString([]byte(jwtSecret))
}
