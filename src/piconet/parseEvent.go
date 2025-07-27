package piconet

import (
	"strings"
)

func ParseEvent(commandText string) Cmd {

	var (
		cmdArgs []string
		cmd     string
		events  []string
		ok      bool
		argText string
	)

	// list of piconet events commands
	events = []string{"STATUS", "ERROR", "MONITOR", "RX_BROADCAST", "RX_IMMEDIATE", "RX_TRANSMIT", "TX_RESULT"}

	for _, cmd = range events {
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
