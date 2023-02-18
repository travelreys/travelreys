package auth

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	OIDCProviderGoogle = "google"
)

const (
	LabelBaseCountry     = "basecountry"
	LabelDefaultCurrency = "currency"
	LabelDefaultLocale   = "locale"

	LabelGoogleID            = "google|id"
	LabelGoogleUserPicture   = "google|picture"
	LabelGoogleVerifiedEmail = "google|verifiedEmail"
)

type PhoneNumber struct {
	CountryCode string
	Number      string
}

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`

	PhoneNumber PhoneNumber `json:"phoneNumber"`

	Labels map[string]string `json:"labels"`
}

type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func UserFromGoogleUser(gusr GoogleUser) User {
	return User{
		ID:          uuid.New().String(),
		Email:       gusr.Email,
		Name:        gusr.Name,
		PhoneNumber: PhoneNumber{},
		Labels: map[string]string{
			LabelGoogleID:            gusr.ID,
			LabelGoogleUserPicture:   gusr.Picture,
			LabelGoogleVerifiedEmail: fmt.Sprintf("%t", gusr.VerifiedEmail),
			LabelDefaultLocale:       gusr.Locale,
		},
	}
}
