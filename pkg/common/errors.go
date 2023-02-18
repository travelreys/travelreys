package common

import (
	"errors"
)

var (
	ErrorInvalidEndpointRequestType  = errors.New("endpoint-invalid-req-type")
	ErrInvalidEndpointRequestContext = errors.New("endpoint-invalid-req-context")
)
var (
	ErrBadPath         = errors.New("http-bad-path")
	ErrInvalidRequest  = errors.New("http-invalid-request")
	ErrInvalidJSONBody = errors.New("http-invalid-json-body")
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
