package server

import (
	"fmt"
	"github.com/johnnewcombe/econet-simple-server/src/comms"
	"github.com/johnnewcombe/econet-simple-server/src/piconet"
	"log/slog"
	"strings"
)

//import "github.com/johnnewcombe/econet-simple-server/logger"

func Listener(comms comms.CommunicationClient, ch chan byte) {

	var (
		ec             piconet.Cmd
		s              strings.Builder
		err            error
		scoutFrame     piconet.Frame
		dataFrame      piconet.Frame
		statusResponse piconet.StatusResponse
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
			// this could be a response from a piconet command e.g. STATUS or a Piconet Event due to data received over Econet
			ec = piconet.ParseEvent(s.String())

			switch ec.Cmd {

			case "STATUS":

				if statusResponse, err = piconet.NewStatusResponse(ec.Args); err != nil {
					slog.Error(err.Error())
				}
				slog.Info(fmt.Sprintf("STATUS: %s", statusResponse.String()))

				break

			case "RX_TRANSMIT":

				//See https://www.npmjs.com/package/@jprayner/piconet-nodejs for protocol details for each response etc.

				if scoutFrame, err = piconet.NewFrame(ec.Args[0]); err != nil {

					slog.Error(err.Error())
				}

				slog.Info(fmt.Sprintf("RX_TRANSMIT: frame=scout, %s", scoutFrame.String()))

				if dataFrame, err = piconet.NewFrame(ec.Args[1]); err != nil {
					slog.Error(err.Error())
				}
				slog.Info(fmt.Sprintf("RX_TRANSMIT: frame=data, %s", dataFrame.String()))

				// PROCESS RX_TRANSMIT
				// TODO rework this into some form of command parser
				//  *I AM
				const kCtrlByte = 0x80
				const kPort = 0x99

				if scoutFrame.ControlByte != kCtrlByte {
					slog.Error("ignoring request due to unexpected control byte")
				}
				if scoutFrame.Port != kPort {
					slog.Error("ignoring request due to unexpected port")
				}
				if len(dataFrame.Data) < 5 {
					slog.Error("data frame too short")
				}

				slog.Info(fmt.Sprintf("RX_TRANSMIT: %s", dataFrame.Data))
				// TODO Investigate the dataFrame (piconetPacket?) as the byte that is laced in the ControlByte
				//   property may be the replyPort (try monitoring ecoclient with the BBC FS3? Plug BBC into other
				//   clock output).
				replyPort := dataFrame.ControlByte
				reply := []byte{scoutFrame.SrcStn, scoutFrame.SrcNet, kCtrlByte, replyPort, 0x05, 0x00, 0x01, 0x02, 0x04, 0x00}
				comms.Write(reply)

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
