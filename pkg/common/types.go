package common

import "errors"

const (
	DefaultChSize = 512
)

type Labels map[string]string
type Tags map[string]string

type GenericJSON map[string]interface{}

func UInt64Ptr(i uint64) *uint64 { return &i }
func Int64Ptr(i int64) *int64    { return &i }
func StringPtr(i string) *string { return &i }
func BoolPtr(i bool) *bool       { return &i }

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
