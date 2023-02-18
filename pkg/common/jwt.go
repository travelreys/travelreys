package common

import (
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

const (
	JWTIssuer     = "tiinyplanet"
	JWTClaimIss   = "iss"
	JWTClaimSub   = "sub"
	JWTClaimEmail = "email"
	JWTClaimIat   = "iat"
)

var (
	JWTDefaultSigningMethod = jwt.SigningMethodHS512
)

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
