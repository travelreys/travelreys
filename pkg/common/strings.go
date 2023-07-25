package common

import (
	"fmt"
	"regexp"
)

func FmtString(val interface{}) string {
	return fmt.Sprintf("%+v", val)
}

func StringContains(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}

var EmailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// IsEmailValid checks if the email provided passes the required structure and length.
func IsEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return EmailRegexp.MatchString(e)
}
