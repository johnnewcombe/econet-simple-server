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
		err error
		ec  econet.Cmd
		s   strings.Builder
	)

	s = strings.Builder{}

	for {
		// TODO: this will block the next byte if not collected quickly
		b := <-ch

		s.WriteByte(b)
		if b == 0x0d || b == 0x10 {

			// TODO do we do this asynchronously? maybe not as the user will be
			//  waiting for the results of the command
			//  do we handle the action to be performed  in 'ParseCommand'?
			//  or do we return a cmd object and do it here?
			if ec, err = econet.ParseCommand(s.String()); err != nil {
				logger.LogError.Printf(err.Error())
				comms.Write([]byte(fmt.Sprintf("%s\r\n", strings.ToUpper(err.Error()))))
			}
			logger.LogInfo.Printf("RX: %s, %s", ec.Cmd, ec.Args)

			// empty the string builder
			s.Reset()
		}

		//logger.LogDebug.Printf(" %s", string(b))
	}
}
