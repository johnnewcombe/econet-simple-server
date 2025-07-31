package piconet

import (
	"strings"
)

func ParseEvent(commandText string) Event {

	var (
		cmdArgs []string
		event   string
		events  []string
		ok      bool
		argText string
	)

	// list of piconet events commands
	events = []string{"STATUS", "ERROR", "MONITOR", "RX_BROADCAST", "RX_IMMEDIATE", "RX_TRANSMIT", "TX_RESULT"}

	for _, event = range events {
		if _, argText, ok = strings.Cut(commandText, event); ok { // i.e. if ok
			cmdArgs = strings.Split(strings.Trim(argText, " "), " ")
			break
		}
	}

	if !ok {
		// not understood so return an empty event object
		return Event{
			CmdText: commandText,
		}
	}

	// all good so populate the command and return
	return Event{
		Cmd:     event,
		Args:    cmdArgs,
		CmdText: commandText,
	}
}
