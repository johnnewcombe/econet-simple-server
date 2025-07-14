package server

import (
	"fmt"
	"github.com/johnnewcombe/econet-simple-server/comms"
	"github.com/johnnewcombe/econet-simple-server/econet"
	"github.com/johnnewcombe/econet-simple-server/logger"
	"strings"
)

//import "github.com/johnnewcombe/econet-simple-server/logger"

func Listener(comms comms.CommunicationClient, ch chan byte) {

	var (
		ec econet.Cmd
		s  strings.Builder
	)

	s = strings.Builder{}

	for {
		// TODO: this will block the next byte if not collected quickly
		// this retrieves data from the serial port e.g. /dev/econet but is blocking
		b := <-ch

		// add byte to string builder
		s.WriteByte(b)

		// end of the command?
		if b == 0x0d && s.Len() > 1 {

			// TODO do we do this asynchronously? maybe not as the user will be
			//  waiting for the results of the command
			//  do we handle the action to be performed  in 'ParseCommand'?
			//  or do we return a cmd object and do it here?
			ec = econet.ParseCommand(s.String())

			// unknown command so must be a message to be returned to the client via the serial port
			// TODO: This will need to be packaged up
			comms.Write([]byte(fmt.Sprintf("%s\n", ec)))

			logger.LogInfo.Printf("RX: %s, %s", ec.Cmd, ec.Args)

			// empty the string builder
			s.Reset()
		}

		//logger.LogDebug.Printf(" %s", string(b))
	}
}
