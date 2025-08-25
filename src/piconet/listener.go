package piconet

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/johnnewcombe/econet-simple-server/src/econet"
	"github.com/johnnewcombe/econet-simple-server/src/lib"
)

//import "github.com/johnnewcombe/econet-simple-server/logger"

func Listener(comms CommunicationClient, ch chan byte) {

	var (
		ec             Event
		s              strings.Builder
		err            error
		rxTransmit     *RxTransmit
		statusResponse *StatusResponse
		monitor        Monitor
		reply          *econet.FSReply
		functionCode   byte
	)

	s = strings.Builder{}

	// TODO Understand Broadcasts
	//piconet.Broadcast(comms, []byte("Piconet Simple File Server Active"))

	for {
		// this retrieves data from the serial port e.g. /dev/econet but is blocking
		b := <-ch

		// add byte to string builder
		s.WriteByte(b)

		// end of the command?
		if b == 0x0d && s.Len() > 1 {

			// this could be a response from a piconet command e.g. STATUS or a Piconet Event due to data received over Econet
			ec = ParseEvent(tidyText(s.String()))

			switch ec.Cmd {

			case "MONITOR":

				if monitor, err = NewMonitor(ec.Args); err != nil {
					slog.Error(err.Error())
				}
				slog.Info("piconet-event=MONITOR")
				lib.LogData(monitor.Frame)
				break

			case "STATUS":

				if statusResponse, err = NewStatusResponse(ec.Args); err != nil {
					slog.Error(err.Error())
				}
				slog.Info("piconet-event=STATUS",
					"major-ver", statusResponse.MajorVersion,
					"minor-ver", statusResponse.MinorVersion,
					"patch", statusResponse.Patch,
					"station", statusResponse.Station,
					"status-reg", statusResponse.StatusReg,
					"mode", statusResponse.Mode,
					"mode-name", ModeMap[statusResponse.Mode])

				// check CTS status as this will be high if the clock is missing
				if statusResponse.StatusReg&0b00010000 > 0 {
					slog.Error("piconet-event=STATUS, msg=Missing clock?")
				}
				break

			case "RX_TRANSMIT":

				//See https://www.npmjs.com/package/@jprayner/piconet-nodejs for protocol details for each response etc.

				if rxTransmit, err = NewRxTransmit(ec.Args); err != nil {
					slog.Error(err.Error())
				}

				// TODO: consider that the reply port and function code are invalid if it is a frame received when in
				//   data transfer mode

				// get logged in status of the machine could this user or a previous one
				//session := econet.ActiveSessions.GetSession(rxTransmit.ScoutFrame.SrcStn, rxTransmit.ScoutFrame.SrcNet)

				// TODO use ToBytes() to log the data and this should include all properties
				slog.Info("piconet-event=RX_TRANSMIT",
					"dst-stn", rxTransmit.ScoutFrame.DstStn,
					"dst-net", rxTransmit.ScoutFrame.DstNet,
					"src-stn", rxTransmit.ScoutFrame.SrcStn,
					"src-net", rxTransmit.ScoutFrame.SrcNet,
					"control-byte", rxTransmit.ScoutFrame.ControlByte,
					"port", rxTransmit.ScoutFrame.Port,
					"port-desc", econet.PortMap[rxTransmit.ScoutFrame.Port])
				// log-level=debug only
				lib.LogData(rxTransmit.ScoutFrame.ToBytes())

				slog.Info("piconet-event=RX_TRANSMIT",
					"dst-stn", rxTransmit.DataFrame.DstStn,
					"dst-net", rxTransmit.DataFrame.DstNet,
					"src-stn", rxTransmit.DataFrame.SrcStn,
					"src-net", rxTransmit.DataFrame.SrcNet)
				// log-level=debug only
				lib.LogData(rxTransmit.DataFrame.ToBytes())

				if rxTransmit.ScoutFrame.ControlByte != econet.CtrlByte {
					slog.Error("piconet-event=RX_TRANSMIT, msg=ignoring request due to unexpected control byte")
					break
				}

				// when in data transfer mode data would come in using a data port determined by the initial
				// request from the client request e.g. fc1 (Save) so we need to handle this by checking the
				// data port in the users session. A data port > 0 means the user is in data transfer mode.

				if len(rxTransmit.DataFrame.Data) < 5 {
					slog.Error("piconet-event=RX_TRANSMIT, msg=data frame too short")
					break
				} // else if session == nil && rxTransmit.ScoutFrame.Port != kPort {
				// TODO fix this
				//slog.Error("piconet-event=RX_TRANSMIT, msg=ignoring request due to unexpected port")
				//break
				//} else if session != nil && rxTransmit.ScoutFrame.Port != session.DataPort {
				//	slog.Error("piconet-event=RX_TRANSMIT, msg=ignoring request due to unexpected port")
				//	break
				//}

				// TODO fix this as we do not need to send received port and function code if we are also sending the
				// data frames data

				// need to set function code explicitly as a parameter to processFunction code as it is not
				// always present in the data e.g. when processing data blocks
				if rxTransmit.ScoutFrame.Port == econet.DataPort && econet.FileXfer != nil {
					functionCode = econet.FileXfer.FunctionCode
				} else {
					functionCode = rxTransmit.DataFrame.Data[1]
				}

				if reply, err = econet.ProcessFunctionCode(
					rxTransmit.DataFrame.SrcStn,
					rxTransmit.DataFrame.SrcNet,
					functionCode,               // TODO function code IS NEEDED as it doesn't always appear in data
					rxTransmit.ScoutFrame.Port, // port is the port that the request was sent on
					rxTransmit.DataFrame.Data); err != nil {
					slog.Error(err.Error())
				}

				if reply != nil {

					//					slog.Info(fmt.Sprintf("piconet-eventREPLY: dst-stn=%02X, dst-net=%02X, return-code=%s",
					//						rxTransmit.ScoutFrame.SrcStn, rxTransmit.ScoutFrame.SrcNet, string(econet.ReplyCodeMap[reply.ReturnCode])))

					//
					// The Piconet firmware adds the Source Station and Net bytes to the reply.
					Transmit(comms,
						rxTransmit.ScoutFrame.SrcStn, // this is the client's station id and now becomes the destination
						rxTransmit.ScoutFrame.SrcNet,
						econet.CtrlByte,
						reply.ReplyPort,
						reply.Data,
						[]byte{})

				} else {
					slog.Error("piconet-event=RX_TRANSMIT: msg=server error, server's reply is nil")
				}

				break
			case "TX_RESULT":
				slog.Info(fmt.Sprintf("piconet-event=TX_RESULT, msg=%s", ec.Args[0]))
				break
			case "ERROR":
				slog.Info(fmt.Sprintf("piconet-event=ERROR, error=%s", ec.CmdText))
				break
			case "RX_BROADCAST":
				slog.Info(fmt.Sprintf("piconet-event=RX_BROADCAST,msg= %s", ec.Args[0]))
				break
			case "RX_IMMEDIATE":
				slog.Info(fmt.Sprintf("piconet-event=RX_IMMEDIATE, msg=%s", ec.Args[0]))
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

// tidyText Removes whitespace e.g. 'I AM' and ' I   AM ' are both valid.
func tidyText(text string) string {

	text = strings.Trim(text, "\x00\n\r ")

	s := strings.Builder{}
	items := lib.Split(text, " ")
	for _, item := range items {
		s.WriteString(item)
		s.WriteString(" ")
	}

	return strings.TrimRight(s.String(), " ")
}
