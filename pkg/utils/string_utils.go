package utils

import "strings"

func IsWhitespaceOrEmpty(value *string) bool {
	if value == nil {
		return true
	}

	return strings.TrimSpace(*value) == ""
}
