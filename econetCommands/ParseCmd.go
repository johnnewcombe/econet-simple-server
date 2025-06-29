package econetCommands

import (
	"errors"
	"github.com/johnnewcombe/econet-simple-server/logger"
	"slices"
	"strings"
)

func ParseCommand(commandText string) (Cmd, error) {

	var (
		cmdArgs        []string
		cmd            string
		econetCommands []string
		ok             bool
		argText        string
		ec             Cmd
	)

	// list of available commands
	econetCommands = []string{"I AM", "CAT", "NOTIFY"}

	commandText = tidyText(commandText)
	logger.LogInfo.Printf("RX: %s", commandText)

	for _, cmd = range econetCommands {
		if _, argText, ok = strings.Cut(commandText, cmd); ok {
			cmdArgs = strings.Split(strings.Trim(argText, " "), " ")
			break
		}
	}

	if !ok {
		// return "bad command"
		return Cmd{}, errors.New("bad command")
	}

	ec = Cmd{
		Cmd:     cmd,
		Args:    cmdArgs,
		CmdText: commandText,
	}

	return ec, nil

}

func split(commandText string, separator string) []string {

	items := slices.DeleteFunc(strings.Split(commandText, separator), func(e string) bool {
		return e == ""
	})
	return items
}

// tidyText Removes whitespace e.g. 'I AM' and ' I   AM ' are both valid.
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
