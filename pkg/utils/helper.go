package utils

import (
	"runtime/debug"
	"strings"
)

func IsWhitespaceOrEmpty(value *string) bool {
	if value == nil {
		return true
	}

	return strings.TrimSpace(*value) == ""
}

func GetVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		return info.Main.Version
	}

	return "local-dev"
}
