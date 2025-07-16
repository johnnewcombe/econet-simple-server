package piconet

import (
	"github.com/johnnewcombe/econet-simple-server/logger"
	"slices"
	"strings"
)

func ParseEvent(commandText string) Cmd {

	var (
		cmdArgs      []string
		cmd          string
		econetEvents []string
		ok           bool
		argText      string
	)

	// list of piconet events commands
	econetEvents = []string{"RX_TRANSMIT", "TX_TRANSMIT"}

	logger.LogInfo.Printf("RX: %s", commandText)

	for _, cmd = range econetEvents {
		if _, argText, ok = strings.Cut(commandText, cmd); ok { // i.e. if ok
			cmdArgs = strings.Split(strings.Trim(argText, " "), " ")
			break
		}
	}

	if !ok {
		// not understood so return an empty cmd object
		return Cmd{
			CmdText: commandText,
		}
	}

	// all good so populate the command and return
	return Cmd{
		Cmd:     cmd,
		Args:    cmdArgs,
		CmdText: commandText,
	}
}

func split(commandText string, separator string) []string {

	items := slices.DeleteFunc(strings.Split(commandText, separator), func(e string) bool {
		return e == ""
	})
	return items
}

// tidyText Removes whitespace e.g. 'I AM' and ' I   AM ' are both valid.
func tidyText(text string) string {

	text = strings.Trim(text, "\r")
	text = strings.Trim(text, "\n")

	s := strings.Builder{}
	items := split(text, " ")
	for _, item := range items {
		s.WriteString(item)
		s.WriteString(" ")
	}

	return strings.TrimRight(strings.ToUpper(s.String()), " ")
}
