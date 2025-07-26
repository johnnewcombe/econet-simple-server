package econet

import (
	"github.com/johnnewcombe/econet-simple-server/src/comms"
)

func ParseCommand(comms comms.CommunicationClient, command string) []byte {

	// PROCESS RX_TRANSMIT

	print(command)

	/*
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
	*/

	//replyPort := rxTransmit.DataFrame.Data[0]

	// TODO Remove dummy reply for a real one
	// TODO Better understand the control port
	// issue a dummy successful reply
	data := []byte{0x05, 0x00, 0x01, 0x02, 0x04, 0x00}
	/*
		0x05, // indicates a successful login
		0x00, // return code of zero indicates success
		0x01, // user root dir handle
		0x02, // currently selected dir handle
		0x04, // library dir handle
		0x00, // boot option (0 = none)
	*/

	return data
	// send the reply, this will generate a TX_RESULT event
	//piconet.Transmit(comms, rxTransmit.ScoutFrame.SrcStn, rxTransmit.ScoutFrame.SrcNet, kCtrlByte, replyPort, data, []byte{})
}
func IAM() {

}
