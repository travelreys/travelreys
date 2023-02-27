package common

import "fmt"

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