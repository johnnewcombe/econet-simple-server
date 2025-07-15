package server

import (
	"encoding/base64"
	"fmt"
	"github.com/johnnewcombe/econet-simple-server/comms"
	"github.com/johnnewcombe/econet-simple-server/econet"
	"github.com/johnnewcombe/econet-simple-server/logger"
	"strings"
)

//import "github.com/johnnewcombe/econet-simple-server/logger"

func Listener(comms comms.CommunicationClient, ch chan byte) {

	var (
		ec   econet.Cmd
		s    strings.Builder
		err  error
		data []byte
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

			switch ec.Cmd {
			case "RX_TRANSMIT":
				//See https://www.npmjs.com/package/@jprayner/piconet-nodejs for protocol details for each response etc.

				// RX_TRANSMIT scout data
				d := []byte(ec.Args[0])
				print(d)

				/* Scout Frame
				   +------+------+-----+-----+---------+------+
				   | Dest | Dest | Src | Src | Control | Port |
				   | Stn  | Net  | Stn | Net |  Byte   |      |
				   +------+------+-----+-----+---------+------+
				    <-------- - - Packet Header - - ---------> <--- - - Packet Data - - --->
				*/

				if data, err = base64.StdEncoding.DecodeString(ec.Args[0]); err != nil {
					print(err)
				}
				print(data)
				// result from I AM SYST SYST
				// 252 96 200 0 128 153
				if data, err = base64.StdEncoding.DecodeString(ec.Args[1]); err != nil {
					print(err)
				}
				// result from I AM SYST SYST
				// 252 96 200 0 144 0 0 0 0 73 32 65 77 32 89 17 115 116 32 115 121 11 116 13
				print(data)
			}

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
