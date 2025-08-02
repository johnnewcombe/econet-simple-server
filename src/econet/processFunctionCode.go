package econet

import (
	"fmt"
	"github.com/johnnewcombe/econet-simple-server/src/lib"
	"log/slog"
	"strings"
)

func ProcessFunctionCode(functionCode byte, command string, srcStationId byte, srcNetworkId byte) []byte {

	var (
		reply []byte
	)

	// tidy the command string
	command = tidyText(command)

	switch functionCode {
	case 0:

		reply = fc0ProcessCommand(command, srcStationId, srcNetworkId)
		break
	case 1:
		reply = fc1save(command, srcStationId, srcNetworkId)
		break
	case 2:
		break
	case 3:
		break
	case 4:
		break
	case 5:
		break
	case 6:
		break
	case 7:
		break
	case 8:
		break
	case 9:
		break
	case 10:
		break
	case 11:
		break
	case 12:
		break
	case 13:
		break
	case 14:
		break
	case 15:
		break
	case 16:
		break
	case 17:
		break
	case 18:
		break
	case 19:
		break
	case 20:
		break
	case 21:
		break
	case 22:
		break
	case 23:
		break
	case 24:
		break
	case 25:
		break
	case 26:
		break
	case 27:
		break
	case 28:
		break
	case 29:
		break
	case 30:
		break
	case 31:
		break
	case 32:
		break
	case 33:
		break
	case 34:
		break
	case 35:
		break
	case 36:
		break
	case 37:
		break
	case 38:
		break
	case 39:
		break
	case 40:
		break
	case 41:
		break
	case 42:
		break
	case 43:
		break
	case 44:
		break
	case 45:
		break
	case 46:
		break
	}

	return reply
}

func fc0ProcessCommand(command string, srcStationId byte, srcNetworkId byte) []byte {

	var (
		reply []byte
		cmd   CliCmd
	)

	cmd = parseCommand(command)

	slog.Info(fmt.Sprintf("econet-f0-cli:, data=[% 02X]", cmd.ToBytes()))

	// these are all * commands
	switch cmd.Cmd {

	case "SAVE":
		break
	case "LOAD":
		break
	case "CAT":
		break
	case "INFO":
		break
	case "I AM":
		reply = f0_Iam(cmd, srcStationId, srcNetworkId)
		break
	case "SDISK":
		break
	case "DIR":
		break
	case "LIB":
		break
	default:

	}

	return reply
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
