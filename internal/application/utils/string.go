package utils

import (
	"strings"
)

func CapitalFirstLetter(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}
