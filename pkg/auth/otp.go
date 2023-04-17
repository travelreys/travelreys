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
	ErrProviderOTPError         = errors.New("auth.service.otp.error")
	ErrProviderOTPEmailNotFound = errors.New("auth.service.otp.notfound")
	ErrProviderOTPEmailExists   = errors.New("auth.service.otp.exists")
	ErrProviderOTPNotSet        = errors.New("auth.service.otp.notset")
	ErrProviderOTPInvalidEmail  = errors.New("auth.service.otp.invalidemail")
	ErrProviderOTPInvalidPw     = errors.New("auth.service.otp.invalidpw")

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

func (prv OTPProvider) parseAuthCode(code string) (string, string, error) {
	sDec, err := base64.StdEncoding.DecodeString(code)
	if err != nil {
		return "", "", err
	}
	tkns := strings.Split(string(sDec), "|")
	if len(tkns) < 2 {
		return "", "", ErrProviderOTPError
	}
	email := tkns[0]
	if !common.IsEmailValid(email) {
		return "", "", ErrProviderOTPInvalidEmail
	}
	otp := tkns[1]
	return email, otp, nil
}

func (prv OTPProvider) TokenToUserInfo(ctx context.Context, code string) (User, error) {
	email, pw, err := prv.parseAuthCode(code)
	if err != nil {
		return User{}, err
	}
	usr, err := prv.store.Read(ctx, ReadFilter{Email: email})
	if err != nil {
		return User{}, ErrProviderOTPEmailNotFound
	}

	hashedPw, err := prv.store.GetOTP(ctx, usr.ID)
	if err := prv.ValidateOTP([]byte(pw), []byte(hashedPw)); err != nil {
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
		return "", "", ErrProviderOTPError
	}
	return pw, string(hashedPw), nil
}

func (prv OTPProvider) ValidateOTP(otp, hashedOTP []byte) error {
	if err := bcrypt.CompareHashAndPassword(hashedOTP, otp); err != nil {
		return ErrProviderOTPInvalidPw
	}
	return nil
}

func (prv OTPProvider) GenerateMagicLinkEmail(usr User, otp string) (string, error) {
	authCode := fmt.Sprintf("%s|%s", usr.Email, otp)
	sEnc := base64.StdEncoding.EncodeToString([]byte(authCode))

	magicLink := fmt.Sprintf("https://www.travelreys.com/magic-link?c=%s", sEnc)
	bodyTmpl := `
	<div>
	<p>Welcome to travelreys! Click the following magic link to login.</p>
	<br />
	<a href="%s" target="_blank" rel="noopener noreferrer">%s</a>
	</div>
	`
	body := fmt.Sprintf(bodyTmpl, magicLink, magicLink)
	return body, nil
}
