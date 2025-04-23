package utils

import (
	"fmt"
	"unicode/utf8"
)

func IsValidMaxLengthUnicode(text string, maxLength int) bool {
	return utf8.RuneCountInString(text) > maxLength
}
func GenerateUserLastname() string {
	return fmt.Sprintf("%s-%s", "gumusic", GenerateRandomString(8))
}
