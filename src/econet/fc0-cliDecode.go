package econet

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/johnnewcombe/econet-simple-server/src/lib"
)

func fc0CliDecode(srcStationId byte, srcNetworkId byte, data []byte) (*FSReply, error) {

	var (
		reply     *FSReply
		cmd       CliCmd
		err       error
		replyPort byte
	)

	replyPort = data[0]

	var command string = ""
	if len(data) > 0 {
		command = strings.Split(string(data[5:]), "\r")[0]
	}

	// TODO Do we need to support abbreviated commands e.g. *. or *S. etc
	cmd = *NewCliCmd(tidyText(command))

	slog.Info("econet-f0-cli-decode:",
		"src-stn", srcStationId,
		"src-net", srcNetworkId,
		"reply-port", replyPort,
		"cmd", cmd.ToString())

	// these are all * commands
	switch cmd.Cmd {

	case "SAVE":
		reply, err = f0Save(cmd, srcStationId, srcNetworkId, replyPort)
		break
	case "LOAD":
		break
	case "CAT":
		break
	case "INFO":
		break
	case "I AM":
		reply, err = f0Iam(cmd, srcStationId, srcNetworkId, replyPort)
		break
	case "SDISK":
		break
	case "DIR":
		break
	case "LIB":
		break
	default:
		//reply = NewFSReply(replyPort, CCIam, RCBadCommand, ReplyCodeMap[RCBadCommand])
		err = errors.New("not implemented")
	}

	return reply, err
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
