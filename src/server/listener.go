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
		rxTransmit     piconet.RxTransmit
		statusResponse piconet.StatusResponse
		monitor        piconet.Monitor
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

			case "MONITOR":

				if monitor, err = piconet.NewMonitor(ec.Args); err != nil {
					slog.Error(err.Error())
				}
				slog.Info(fmt.Sprintf("piconet-event=MONITOR, frame=[% 02X]", monitor.Frame))
				break

			case "STATUS":

				if statusResponse, err = piconet.NewStatusResponse(ec.Args); err != nil {
					slog.Error(err.Error())
				}
				slog.Info(fmt.Sprintf("piconet-event=STATUS, %s", statusResponse.String()))

				// check CTS status as this will be high if the clock is missing
				if statusResponse.StatusReg&0b00010000 > 0 {
					slog.Error("piconet-event=STATUS, msg=Missing clock?")
				}
				break

			case "RX_TRANSMIT":

				//See https://www.npmjs.com/package/@jprayner/piconet-nodejs for protocol details for each response etc.

				if rxTransmit, err = piconet.NewRxTransmit(ec.Args); err != nil {
					slog.Error(err.Error())
				}

				slog.Info(fmt.Sprintf("piconet-event=RX_TRANSMIT %s", rxTransmit.String()))

				// PROCESS RX_TRANSMIT
				// TODO rework this into some form of command parser and check with user data
				//  *I AM
				const kCtrlByte = 0x80
				const kPort = 0x99

				if rxTransmit.ScoutFrame.ControlByte != kCtrlByte {
					slog.Error("piconet-event=RX_TRANSMIT, msg=ignoring request due to unexpected control byte")
				}
				if rxTransmit.ScoutFrame.Port != kPort {
					slog.Error("piconet-event=RX_TRANSMIT, msg=ignoring request due to unexpected port")
				}
				if len(rxTransmit.DataFrame.Data) < 5 {
					slog.Error("piconet-event=RX_TRANSMIT, msg=data frame too short")
				}

				//slog.Info(fmt.Sprintf("piconet-event=RX_TRANSMIT, %s", rxTransmit.DataFrame.Data))
				// TODO Investigate the dataFrame (piconetPacket?) as the byte that is laced in the ControlByte
				//   property may be the replyPort (try monitoring ecoclient with the BBC FS3? Plug BBC into other
				//   clock output).
				replyPort := rxTransmit.DataFrame.Data[0]

				//[64 00 FB 00 80 90]
				// TODO Remove dummy reply for a real one
				// issue a dummy successful reply
				//encodedData := []byte(base64.StdEncoding.EncodeToString([]byte{0x05, 0x00, 0x01, 0x02, 0x04, 0x00}))
				data := []byte{0x05, 0x00, 0x01, 0x02, 0x04, 0x00}

				// send the reply, this will generate a TX_RESULT event
				piconet.Transmit(comms, rxTransmit.ScoutFrame.SrcStn, rxTransmit.ScoutFrame.SrcNet, kCtrlByte, replyPort, data)

				/*
					Exmple Response for successful login to BBC L3 FS
					64 00 FB 00 05 00 01 02 04 00
				*/

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
			case "TX_RESULT":
				slog.Info(fmt.Sprintf("piconet-event=TX_RESULT, msg=%s", ec.CmdText))
				break
			case "ERROR":
				slog.Info(fmt.Sprintf("piconet-event=ERROR, error=%s", ec.CmdText))
				break
			case "RX_BROADCAST":
				slog.Info(fmt.Sprintf("piconet-event=RX_BROADCAST,msg= %s", ec.CmdText))
				break
			case "RX_IMMEDIATE":
				slog.Info(fmt.Sprintf("piconet-event=RX_IMMEDIATE, msg=%s", ec.CmdText))
				break
			default:
				slog.Info(fmt.Sprintf("piconet-event=UNKNOWN, msg=%s", ec.CmdText))
			}

			// empty the string builder
			s.Reset()

		}

		//logger.LogDebug.Printf(" %s", string(b))
	}
}
