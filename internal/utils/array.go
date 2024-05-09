package utils

import "strings"

// InArray checks whether the value is in an array.
func InArray(arr []string, value string) bool {
	for _, v := range arr {
		if strings.EqualFold(v, value) {
			return true
		}
	}

	return false
}
