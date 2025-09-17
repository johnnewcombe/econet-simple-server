package lib

import (
	"slices"
	"strings"
)

// Split Splits a string based on the specified separator and removes empty elements.
func Split(commandText string, separator string) []string {

	items := slices.DeleteFunc(strings.Split(commandText, separator), func(e string) bool {
		return e == ""
	})
	return items
}

func LeftString(s string, n int) string {

	if len(s) > 10 {
		return s[:10]
	}
	return s
}
