package econet

import (
	"slices"
	"strings"
)

func tidyText(text string) string {

	text = strings.Trim(text, "\x00")
	text = strings.Trim(text, "\n")
	text = strings.Trim(text, "\r")
	text = strings.ToUpper(text)

	s := strings.Builder{}
	items := split(text, " ")
	for _, item := range items {
		s.WriteString(item)
		s.WriteString(" ")
	}

	return strings.TrimRight(s.String(), " ")
}
func split(commandText string, separator string) []string {

	items := slices.DeleteFunc(strings.Split(commandText, separator), func(e string) bool {
		return e == ""
	})
	return items
}
