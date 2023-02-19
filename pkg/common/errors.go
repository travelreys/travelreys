package common

import (
	"errors"
)

var (
	ErrorInvalidEndpointRequestType = errors.New("endpoint.invalidrequest")
)
var (
	ErrBadPath         = errors.New("http.badpath")
	ErrInvalidRequest  = errors.New("http.invalidrequest")
	ErrInvalidJSONBody = errors.New("http.invalidjson")
)

type Errorer interface {
	Error() error
}

func ErrorContains(slice []error, target error) bool {
	for _, err := range slice {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}
