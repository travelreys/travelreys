package flights

import (
	"net/http"

	"github.com/tiinyplanet/tiinyplanet/pkg/utils"
)

var (
	notFoundErrors     = []error{}
	appErrors          = []error{ErrInvalidSearchRequest}
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
