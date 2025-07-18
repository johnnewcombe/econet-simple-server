package server

import (
	"fmt"
	"github.com/johnnewcombe/econet-simple-server/comms"
	"github.com/johnnewcombe/econet-simple-server/piconet"
	"log/slog"
	"strings"
)

//import "github.com/johnnewcombe/econet-simple-server/logger"

func Listener(comms comms.CommunicationClient, ch chan byte) {

	var (
		ec         piconet.Cmd
		s          strings.Builder
		err        error
		scoutFrame piconet.Frame
		dataFrame  piconet.Frame
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

				if scoutFrame, err = piconet.NewDataFrame(ec.Args[0]); err != nil {

					slog.Error(err.Error())
				}

				slog.Info(fmt.Sprintf("RX_TRANSMIT: frame=scout, %s", scoutFrame.String()))

				if dataFrame, err = piconet.NewDataFrame(ec.Args[1]); err != nil {
					slog.Error(err.Error())
				}
				slog.Info(fmt.Sprintf("RX_TRANSMIT: frame=data, %s", dataFrame.String()))

				/*
					example of a response to *I AM

					  // issue a dummy successful reply
					  const txResult = await driver.transmit(
					    scout.fromStation,
					    scout.fromNetwork,
					    controlByte,
					    replyPort,
					    Buffer.from([
					      0x05, // indicates a successful login
					      0x00, // return code of zero indicates success
					      0x01, // user root dir handle
					      0x02, // currently selected dir handle
					      0x04, // library dir handle
					      0x00, // boot option (0 = none)
					    ]),
					  );
				*/
			}

			// empty the string builder
			s.Reset()
		}

		//logger.LogDebug.Printf(" %s", string(b))
	}
}
