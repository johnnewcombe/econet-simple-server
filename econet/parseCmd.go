package econet

import (
	"github.com/johnnewcombe/econet-simple-server/logger"
	"slices"
	"strings"
)

func ParseCommand(commandText string) Cmd {

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
		// TODO Parse messages? i.e. responses that are not commands
		// not a command could be an error message such as 'ERROR WHAT??' or a STATUS MESSAGE e.g. 'STATUS 2.0.20 254 04 1'
		return Cmd{}
	}

	ec = Cmd{
		Cmd:     cmd,
		Args:    cmdArgs,
		CmdText: commandText,
	}

	return ec
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
