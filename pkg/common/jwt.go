package common

import (
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

const (
	JwtIssuer     = "tiinyplanet"
	JwtClaimIss   = "iss"
	JwtClaimSub   = "sub"
	JwtClaimEmail = "email"
	JwtClaimIat   = "iat"

	EnvJwtSecret = "TIINYPLANET_JWT_SECRET"
)

var (
	JWTDefaultSigningMethod = jwt.SigningMethodHS512
)

func GetJwtSecret() string {
	return os.Getenv(EnvJwtSecret)
}

func GetJwtSecretBytes() []byte {
	return []byte(os.Getenv(EnvJwtSecret))
}

func ParseBearerAndToken(header string) (string, error) {
	bearerAndToken := strings.Split(header, " ")
	if len(bearerAndToken) < 2 {
		return "", ErrInvalidAuthToken
	}
	if bearerAndToken[0] != "Bearer" {
		return "", ErrInvalidAuthToken
	}
	return bearerAndToken[1], nil
}

func ParseJWT(jwtToken, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrMissingJWTClaims
	}
	return claims, nil
}
