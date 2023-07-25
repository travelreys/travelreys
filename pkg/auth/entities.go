package auth

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lucasepe/codename"
	"github.com/travelreys/travelreys/pkg/storage"
)

const (
	OIDCProviderGoogle   = "google"
	OIDCProviderFacebook = "facebook"
	OIDCProviderOTP      = "otp"
)

const (
	LabelAvatarImage     = "avatarImage"
	LabelBaseCountry     = "basecountry"
	LabelDefaultCurrency = "currency"
	LabelDefaultLocale   = "locale"
	LabelOnboarded       = "onboarded"

	LabelGoogleID            = "google|id"
	LabelGoogleUserPicture   = "google|picture"
	LabelGoogleVerifiedEmail = "google|verifiedEmail"

	LabelFacebookID          = "facebook|id"
	LabelFacebookUserPicture = "facebook|picure"
)

type PhoneNumber struct {
	CountryCode string
	Number      string
}

type User struct {
	ID       string `json:"id" bson:"id"`
	Email    string `json:"email" bson:"email"`
	Name     string `json:"name" bson:"name"`
	Username string `json:"username" bson:"username"`

	CreatedAt   time.Time   `json:"createdAt" bson:"createdAt"`
	PhoneNumber PhoneNumber `json:"phoneNumber" bson:"phonenumber"`

	Labels map[string]string `json:"labels" bson:"labels"`
}

func (user User) MakeUserAvatarObject() storage.Object {
	avatarImgURL := user.GetAvatarImgURL()
	tkns := strings.Split(avatarImgURL, "/")

	return storage.Object{
		Bucket: avatarBucket,
		Path:   filepath.Join(avatarFilePrefix, tkns[len(tkns)-1]),
	}
}

func (user User) GetAvatarImgURL() string {
	return user.Labels[LabelAvatarImage]
}

func (user User) GetProfileImgURL() string {
	if user.Labels[LabelAvatarImage] != "" {
		return user.Labels[LabelAvatarImage]
	}
	if user.Labels[LabelGoogleUserPicture] != "" {
		return user.Labels[LabelGoogleUserPicture]
	}
	return ""
}

type UsersList []User
type UsersMap map[string]User

func (m UsersMap) Scrub() {
	for k := range m {
		copy := m[k]
		copy.Email = ""
		m[k] = copy
	}
}

func RandomUsernameGenerator() string {
	rng, _ := codename.DefaultRNG()
	return codename.Generate(rng, 8)
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
		Username:    RandomUsernameGenerator(),
		Labels: map[string]string{
			LabelGoogleID:            gusr.ID,
			LabelGoogleUserPicture:   gusr.Picture,
			LabelGoogleVerifiedEmail: fmt.Sprintf("%t", gusr.VerifiedEmail),
			LabelDefaultLocale:       gusr.Locale,
		},
	}
}

func (gUsr GoogleUser) AddLabelsToUser(usr *User) {
	usr.Labels[LabelGoogleID] = gUsr.ID
	usr.Labels[LabelGoogleUserPicture] = gUsr.Picture
	usr.Labels[LabelGoogleVerifiedEmail] = fmt.Sprintf("%t", gUsr.VerifiedEmail)
}

type FacebookUser struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture struct {
		Data struct {
			URL string `json:"url"`
		}
	} `json:"picture"`
}

func (fbUsr FacebookUser) AddLabelsToUser(usr *User) {
	usr.Labels[LabelFacebookID] = fbUsr.ID
	usr.Labels[LabelFacebookUserPicture] = fbUsr.Picture.Data.URL
}

func UserFromFBUser(fbUsr FacebookUser) User {
	return User{
		ID:          uuid.New().String(),
		Email:       fbUsr.Email,
		Name:        fbUsr.Name,
		PhoneNumber: PhoneNumber{},
		Username:    RandomUsernameGenerator(),
		Labels: map[string]string{
			LabelFacebookID:          fbUsr.ID,
			LabelFacebookUserPicture: fbUsr.Picture.Data.URL,
		},
	}
}

var UsernameRegexp = regexp.MustCompile("^[a-zA-Z0-9_]*$")

func IsValidUsername(userName string) bool {
	if len(userName) < 5 && len(userName) > 255 {
		return false
	}
	return UsernameRegexp.MatchString(userName)
}
