package server

import (
	"github.com/johnnewcombe/econet-simple-server/comms"
	"github.com/johnnewcombe/econet-simple-server/logger"
	"github.com/johnnewcombe/econet-simple-server/piconet"
	"strings"
)

//import "github.com/johnnewcombe/econet-simple-server/logger"

func Listener(comms comms.CommunicationClient, ch chan byte) {

	var (
		ec           piconet.Cmd
		s            strings.Builder
		err          error
		scoutFrame   piconet.EconetFrame
		messageFrame piconet.EconetFrame
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
			//  do we handle the action to be performed  in 'ParseEvent'?
			//  or do we return a cmd object and do it here?
			ec = piconet.ParseEvent(s.String())

			switch ec.Cmd {
			case "RX_TRANSMIT":
				//See https://www.npmjs.com/package/@jprayner/piconet-nodejs for protocol details for each response etc.

				/* Scout Frame
				   +------+------+-----+-----+---------+------+
				   | Dest | Dest | Src | Src | Control | Port |
				   | Stn  | Net  | Stn | Net |  Byte   |      |
				   +------+------+-----+-----+---------+------+
				    <-------- - - Packet Header - - --------->
				*/

				if scoutFrame, err = piconet.CreateFrame(ec.Args[0]); err != nil {
					logger.LogError.Println(err)
				}
				logger.LogInfo.Printf("Received Scout:'%s'", scoutFrame.ToString())

				if messageFrame, err = piconet.CreateFrame(ec.Args[1]); err != nil {
					logger.LogError.Println(err)
				}
				logger.LogInfo.Printf("Received Data:'%s'", messageFrame.ToString())

			}

			// empty the string builder
			s.Reset()
		}

		//logger.LogDebug.Printf(" %s", string(b))
	}
}
