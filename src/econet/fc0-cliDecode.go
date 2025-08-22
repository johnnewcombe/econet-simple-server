package econet

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/johnnewcombe/econet-simple-server/src/lib"
)

func fc0CliDecode(srcStationId byte, srcNetworkId byte, data []byte) (*FSReply, error) {

	var (
		reply *FSReply
		cmd   CliCmd
		err   error
	)

	var command string = ""
	if len(data) > 0 {
		command = strings.TrimRight(string(data), "\r")
	}

	// TODO Do we need to support abbreviated commands e.g. *. or *S. etc
	cmd = *NewCliCmd(tidyText(command))

	slog.Info(fmt.Sprintf("econet-f0-cli: src=%02X/%02X, data=[% 02X], cmd=%s", srcStationId, srcNetworkId, cmd.ToString(), cmd.ToBytes()))

	// these are all * commands
	switch cmd.Cmd {

	case "SAVE":
		reply, err = f0Save(cmd, srcStationId, srcNetworkId)
		break
	case "LOAD":
		break
	case "CAT":
		break
	case "INFO":
		break
	case "I AM":
		reply, err = f0Iam(cmd, srcStationId, srcNetworkId)
		break
	case "SDISK":
		break
	case "DIR":
		break
	case "LIB":
		break
	default:
		reply = NewFSReply(CCIam, RCBadCommmand, ReplyCodeMap[RCBadCommmand])
		err = errors.New("not implemented")
	}

	return reply, err
}

func NewCliCmd(commandText string) *CliCmd {

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
			return &CliCmd{
				Cmd:     cmd,
				CmdText: commandText,
				Args:    cmdArgs,
			}

		}
	}
	return &CliCmd{}
}

func tidyText(text string) string {
	// Trim null characters, newlines, carriage returns, and spaces from both ends
	text = strings.Trim(text, "\x00\n\r ")

	// Convert to uppercase
	text = strings.ToUpper(text)

	// Replace multiple spaces with a single space
	text = strings.Join(lib.Split(text, " "), " ")

	return text
}
