package piconet

import (
	"fmt"
	"github.com/johnnewcombe/econet-simple-server/src/econet"
	"github.com/johnnewcombe/econet-simple-server/src/lib"
	"log/slog"
	"strings"
)

//import "github.com/johnnewcombe/econet-simple-server/logger"

func Listener(comms CommunicationClient, ch chan byte) {

	const (
		kCtrlByte = 0x80
		kPort     = 0x99
	)

	var (
		ec             Event
		s              strings.Builder
		err            error
		rxTransmit     *RxTransmit
		statusResponse *StatusResponse
		monitor        Monitor
		reply          []byte
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
				slog.Info(fmt.Sprintf("piconet-event=MONITOR, frame=[% 02X]", monitor.Frame))
				break

			case "STATUS":

				if statusResponse, err = NewStatusResponse(ec.Args); err != nil {
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

				if rxTransmit, err = NewRxTransmit(ec.Args); err != nil {
					slog.Error(err.Error())
				}

				slog.Info(fmt.Sprintf("piconet-event=RX_TRANSMIT %s", rxTransmit.String()))

				if rxTransmit.ScoutFrame.ControlByte != kCtrlByte {
					slog.Error("piconet-event=RX_TRANSMIT, msg=ignoring request due to unexpected control byte")
				}
				if rxTransmit.ScoutFrame.Port != kPort {
					slog.Error("piconet-event=RX_TRANSMIT, msg=ignoring request due to unexpected port")
				}
				if len(rxTransmit.DataFrame.Data) < 5 {
					slog.Error("piconet-event=RX_TRANSMIT, msg=data frame too short")
				}

				reply = econet.ProcessFunctionCode(
					rxTransmit.DataFrame.FunctionCode,
					rxTransmit.Command(),
					rxTransmit.DataFrame.SrcStn,
					rxTransmit.DataFrame.SrcNet)

				// Function Code 0 - CLI Decoding
				Transmit(comms,
					rxTransmit.ScoutFrame.SrcStn,
					rxTransmit.ScoutFrame.SrcNet,
					kCtrlByte,
					rxTransmit.DataFrame.ReplyPort,
					reply,
					[]byte{})

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

	text = strings.Trim(text, "\x00")
	text = strings.Trim(text, "\n")
	text = strings.Trim(text, "\r")

	s := strings.Builder{}
	items := lib.Split(text, " ")
	for _, item := range items {
		s.WriteString(item)
		s.WriteString(" ")
	}

	return strings.TrimRight(s.String(), " ")
}
