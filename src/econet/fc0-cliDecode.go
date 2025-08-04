package econet

import (
	"fmt"
	"github.com/johnnewcombe/econet-simple-server/src/lib"
	"log/slog"
	"strings"
)

func fc0cliDecode(srcStationId byte, srcNetworkId byte, data []byte) (*FSReply, error) {

	var (
		reply *FSReply
		cmd   CliCmd
		err   error
	)

	var command string = ""
	if len(data) > 0 {
		command = strings.TrimRight(string(data), "\r")
	}

	cmd = parseCommand(tidyText(command))

	slog.Info(fmt.Sprintf("econet-f0-cli:, data=[% 02X]", cmd.ToBytes()))

	// these are all * commands
	switch cmd.Cmd {

	case "SAVE":
		reply, err = f0_Save(cmd, srcStationId, srcNetworkId)
		break
	case "LOAD":
		break
	case "CAT":
		break
	case "INFO":
		break
	case "I AM":
		reply, err = f0_Iam(cmd, srcStationId, srcNetworkId)
		break
	case "SDISK":
		break
	case "DIR":
		break
	case "LIB":
		break
	default:

	}

	return reply, err
}

func parseCommand(commandText string) CliCmd {

	var (
		commands []string
		cmdArgs  []string
		cmd      string
		ok       bool
		argText  string
	)

	commandText = tidyText(commandText)

	// list of piconet * commands
	commands = []string{"SAVE", "LOAD", "CAT", "INFO", "I AM", "SDISK", "DIR", "LIB"}

	for _, cmd = range commands {
		if _, argText, ok = strings.Cut(commandText, cmd); ok { // i.e. if ok
			cmdArgs = strings.Split(strings.Trim(argText, " "), " ")
			return CliCmd{
				Cmd:     cmd,
				CmdText: commandText,
				Args:    cmdArgs,
			}

		}
	}
	return CliCmd{}
}

func tidyText(text string) string {

	text = strings.Trim(text, "\x00")
	text = strings.Trim(text, "\n")
	text = strings.Trim(text, "\r")
	text = strings.ToUpper(text)

	s := strings.Builder{}
	items := lib.Split(text, " ")
	for _, item := range items {
		s.WriteString(item)
		s.WriteString(" ")
	}

	return strings.TrimRight(s.String(), " ")
}
