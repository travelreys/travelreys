package trips

import (
	"net/http"

	"github.com/awhdesmond/tiinyplanet/pkg/utils"
)

var (
	notFoundErrors     = []error{}
	appErrors          = []error{}
	unauthorisedErrors = []error{}
)

func ErrorToHTTPCode(err error) int {
	if utils.ErrorContains(notFoundErrors, err) {
		return http.StatusNotFound
	}
	if utils.ErrorContains(appErrors, err) {
		return http.StatusUnprocessableEntity
	}
	if utils.ErrorContains(unauthorisedErrors, err) {
		return http.StatusUnauthorized
	}
	return http.StatusInternalServerError
}
