package econetCommands

import (
	"github.com/johnnewcombe/econet-simple-server/logger"
	"slices"
	"strings"
)

func ParseCommand(commandText string) error {

	var (
		cmdArgs []string
	)
	commandText = tidyText(commandText)

	// need to remove whitespace e.g. I AM and I  AM should both work.
	if strings.Contains(commandText, "I AM") {
		cmdArgs = strings.Split(commandText, " ")
		logger.LogInfo.Printf("RX: %s, %s", commandText, cmdArgs)
	}

	return nil

}
func split(commandText string, separator string) []string {

	items := slices.DeleteFunc(strings.Split(commandText, separator), func(e string) bool {
		return e == ""
	})
	return items
}

func tidyText(text string) string {

	text = strings.TrimRight(text, "\r")

	s := strings.Builder{}
	items := split(text, " ")
	for _, item := range items {
		s.WriteString(item)
		s.WriteString(" ")
	}

	return strings.TrimRight(strings.ToUpper(s.String()), " ")
}
