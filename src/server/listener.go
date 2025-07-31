package server

import (
	"fmt"
	"github.com/johnnewcombe/econet-simple-server/src/comms"
	"github.com/johnnewcombe/econet-simple-server/src/econet"
	"github.com/johnnewcombe/econet-simple-server/src/piconet"
	"log/slog"
	"slices"
	"strings"
)

//import "github.com/johnnewcombe/econet-simple-server/logger"

func Listener(comms comms.CommunicationClient, ch chan byte) {

	var (
		ec             piconet.Event
		s              strings.Builder
		err            error
		rxTransmit     *piconet.RxTransmit
		statusResponse *piconet.StatusResponse
		monitor        piconet.Monitor
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
			ec = piconet.ParseEvent(tidyText(s.String()))

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
				// process RX_TRANSMIT
				reply := econet.CLIDecode(tidyText(rxTransmit.Command()), rxTransmit.DataFrame.SrcStn, rxTransmit.DataFrame.SrcNet)

				replyPort := rxTransmit.DataFrame.Data[0]
				piconet.Transmit(comms, rxTransmit.ScoutFrame.SrcStn, rxTransmit.ScoutFrame.SrcNet, kCtrlByte, replyPort, reply, []byte{})

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
	items := split(text, " ")
	for _, item := range items {
		s.WriteString(item)
		s.WriteString(" ")
	}

	return strings.TrimRight(s.String(), " ")
}
func split(commandText string, separator string) []string {

	items := slices.DeleteFunc(strings.Split(commandText, separator), func(e string) bool {
		return e == ""
	})
	return items
}
