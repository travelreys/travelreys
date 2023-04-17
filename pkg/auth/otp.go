package auth

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"strings"

	"github.com/travelreys/travelreys/pkg/common"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrProviderEmailError         = errors.New("auth.service.emailpassword.error")
	ErrProviderEmailEmailNotFound = errors.New("auth.service.emailpassword.notfound")
	ErrProviderEmailEmailExists   = errors.New("auth.service.emailpassword.exists")
	ErrProviderEmailNotSet        = errors.New("auth.service.emailpassword.notset")
	ErrProviderEmailInvalidEmail  = errors.New("auth.service.emailpassword.invalidemail")
	ErrProviderEmailInvalidPw     = errors.New("auth.service.emailpassword.invalidpw")

	defaultOTPPeriod = 60
	b32NoPadding     = base32.StdEncoding.WithPadding(base32.NoPadding)
)

type OTPProvider struct {
	store      Store
	randReader io.Reader
}

func NewOTPProvider(store Store, randReader io.Reader) *OTPProvider {
	return &OTPProvider{store, randReader}
}

func (prv OTPProvider) parseAuthCode(code string) (string, []byte, error) {
	tkns := strings.Split(code, "|")
	if len(tkns) < 2 {
		return "", []byte{}, ErrProviderEmailError
	}
	email := tkns[0]
	if !common.IsEmailValid(email) {
		return "", []byte{}, ErrProviderEmailInvalidEmail
	}
	hashedOTP := tkns[1]
	pw, err := base64.StdEncoding.DecodeString(hashedOTP)
	if err != nil {
		return "", []byte{}, ErrProviderEmailError
	}
	return email, pw, nil
}

func (prv OTPProvider) TokenToUserInfo(ctx context.Context, code string) (User, error) {
	email, pw, err := prv.parseAuthCode(code)
	if err != nil {
		return User{}, err
	}
	usr, err := prv.store.Read(ctx, ReadFilter{Email: email})
	if err != nil {
		return User{}, ErrProviderEmailEmailNotFound
	}

	hashedPw, err := prv.store.GetOTP(ctx, usr.ID)
	if err := prv.ValidateOTP(pw, []byte(hashedPw)); err != nil {
		return User{}, err
	}
	return usr, nil
}

func (prv OTPProvider) GenerateOTP(maxDigits uint32) (string, string, error) {
	bi, err := rand.Int(
		prv.randReader,
		big.NewInt(int64(math.Pow(10, float64(maxDigits)))),
	)
	if err != nil {
		return "", "", err
	}

	pw := fmt.Sprintf("%0*d", maxDigits, bi)
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return "", "", ErrProviderEmailError
	}
	return pw, string(hashedPw), nil
}

func (prv OTPProvider) ValidateOTP(otp, hashedOTP []byte) error {
	if err := bcrypt.CompareHashAndPassword(hashedOTP, otp); err != nil {
		return ErrProviderEmailInvalidPw
	}
	return nil
}

func (prv OTPProvider) GenerateMagicLinkEmail(usr User, otp string) (string, error) {
	authCode := fmt.Sprintf("%s|%s", usr.Email, otp)
	sEnc := base64.StdEncoding.EncodeToString([]byte(authCode))

	magicLink := fmt.Sprintf("https://www.travelreys.com/magic-link?c=%s", sEnc)
	bodyTmpl := `
	<div>
	<p>Welcome to travelreys. Click on the following magic link to login.</p>
	<br />
	<a href="%s" target="_blank" rel="noopener noreferrer">%s</a>
	</div>
	`
	body := fmt.Sprintf(bodyTmpl, magicLink, magicLink)
	fmt.Println(body)
	return body, nil
}
