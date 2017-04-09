package utils

import (
	"github.com/alioygur/godash"
	"strings"
)

func ToUpperFirstCamelCase(s string) string {
	return strings.ToUpper(string(s[0])) + godash.ToCamelCase(s)[1:]
}
func ToUpperFirst(s string) string {
	return strings.ToUpper(string(s[0])) + s[1:]
}
func ToLowerSnakeCase(s string) string {
	return strings.ToLower(godash.ToSnakeCase(s))
}

func ToCamelCase(s string) string {
	return godash.ToCamelCase(s)
}
