package common

import "errors"

var (
	ErrInvalidEndpointRequestType    = errors.New("endpoint-invalid-req-type")
	ErrInvalidEndpointRequestContext = errors.New("endpoint-invalid-req-context")
)

type Errorer interface {
	Error() error
}
