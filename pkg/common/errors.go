package common

import (
	"errors"
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
