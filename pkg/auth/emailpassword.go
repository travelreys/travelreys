package auth

import (
	"context"
	"encoding/base64"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/travelreys/travelreys/pkg/common"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrProviderEmailPasswordError         = errors.New("auth.service.emailpassword.error")
	ErrProviderEmailPasswordEmailNotFound = errors.New("auth.service.emailpassword.notfound")
	ErrProviderEmailPasswordEmailExists   = errors.New("auth.service.emailpassword.exists")
	ErrProviderEmailPasswordNotSet        = errors.New("auth.service.emailpassword.notset")
	ErrProviderEmailPasswordInvalidEmail  = errors.New("auth.service.emailpassword.invalidemail")
	ErrProviderEmailPasswordInvalidPw     = errors.New("auth.service.emailpassword.invalidpw")
)

type EmailPasswordProvider struct {
	store Store
}

func NewEmailPasswordProvider(store Store) *EmailPasswordProvider {
	return &EmailPasswordProvider{store}
}

func (prv EmailPasswordProvider) parseAuthCode(code string) (string, []byte, error) {
	tkns := strings.Split(code, "|")
	if len(tkns) < 2 {
		return "", []byte{}, ErrProviderEmailPasswordError
	}
	email := tkns[0]
	if !common.IsEmailValid(email) {
		return "", []byte{}, ErrProviderEmailPasswordInvalidEmail
	}
	base64Hash := tkns[1]
	pw, err := base64.StdEncoding.DecodeString(base64Hash)
	if err != nil {
		return "", []byte{}, ErrProviderEmailPasswordError
	}
	return email, pw, nil
}

func (prv EmailPasswordProvider) TokenToUserInfo(ctx context.Context, code string) (User, error) {
	email, pw, err := prv.parseAuthCode(code)
	if err != nil {
		return User{}, err
	}

	usr, err := prv.store.Read(ctx, ReadFilter{Email: email})
	if err != nil {
		return User{}, ErrProviderEmailPasswordEmailNotFound
	}

	bcryptHash, ok := usr.Labels[LabelPasswordHash]
	if !ok {
		return User{}, ErrProviderEmailPasswordNotSet
	}
	if err := bcrypt.CompareHashAndPassword([]byte(bcryptHash), []byte(pw)); err != nil {
		return User{}, ErrProviderEmailPasswordInvalidPw
	}

	return usr, nil
}

func (prv EmailPasswordProvider) Signup(ctx context.Context, code string) (User, error) {
	email, pw, err := prv.parseAuthCode(code)
	if err != nil {
		return User{}, err
	}

	if _, err = prv.store.Read(ctx, ReadFilter{Email: email}); err == nil {
		return User{}, ErrProviderEmailPasswordEmailExists
	}

	hashedPw, err := bcrypt.GenerateFromPassword(pw, bcrypt.DefaultCost)
	if err != nil {
		return User{}, ErrProviderEmailPasswordError
	}

	newusr := User{
		ID:          uuid.NewString(),
		Email:       email,
		Name:        email,
		CreatedAt:   time.Now(),
		PhoneNumber: PhoneNumber{},
		Labels: common.Labels{
			LabelVerifiedEmail: "false",
			LabelAvatarImage:   "https://cdn.travelreys.com/travelreys-media-demo/avatar/account.png",
			LabelPasswordHash:  string(hashedPw),
		},
	}
	return newusr, nil

}
