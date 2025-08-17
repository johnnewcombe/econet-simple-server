package piconet

import (
	"fmt"
	"log/slog"
)

const (
	MODE_LISTEN  = "LISTEN"
	MODE_MONITOR = "MONITOR"
	MODE_STOP    = "STOP"
)

// TODO Need to be able to tell if client is connected.
// TODO All commands should return error
var piconetMode = map[string]bool{
	"LISTEN":  true,
	"MONITOR": true,
	"STOP":    true,
}

func SetStationID(commsClient CommunicationClient, stationID int) {

	if commsClient != nil {
		slog.Info(fmt.Sprintf("piconet-cmd=SET_STATION, stn=%d", stationID))
		if err := commsClient.Write([]byte(fmt.Sprintf("SET_STATION %d\r", stationID))); err != nil {
			slog.Error(err.Error())
		}

	}
}

func SetMode(commsClient CommunicationClient, mode string) {

	if commsClient != nil {
		if piconetMode[mode] {
			slog.Info(fmt.Sprintf("piconet-cmd=SET_MODE, mode=%s", mode))
			if err := commsClient.Write([]byte(fmt.Sprintf("SET_MODE %s\r", mode))); err != nil {
				slog.Error(err.Error())
			}
		} else {
			slog.Error("invalid mode", "mode", mode)
		}
	}
}

func Status(commsClient CommunicationClient) {

	if commsClient != nil {
		slog.Info(fmt.Sprintf("piconet-cmd=STATUS"))
		if err := commsClient.Write([]byte("STATUS\r")); err != nil {
			slog.Error(err.Error())
		}
	}
}

func Restart(commsClient CommunicationClient) {
	if commsClient != nil {
		slog.Info(fmt.Sprintf("piconet-cmd=RESTART"))
		if err := commsClient.Write([]byte("RESTART\r")); err != nil {
			slog.Error(err.Error())
		}
	}
}

func Transmit(commsClient CommunicationClient, stationId byte, network byte, controlByte byte, port byte, data []byte, extraScoutData []byte) {

	var (
		sReply string
		err    error
	)

	if commsClient != nil {

		encData := Base64Encode(data)

		// TODO I think that the first two bytes of data are always Func code (cli command codes are in the following
		//  byte) and return code adding these as parameters to the function we can report them clearer in the slog message

		slog.Info(fmt.Sprintf("piconet-command=TX, dst-stn=%d, dst-net=%d, ctrl-byte=%d, port=%02Xh, data=[% 02X]",
			stationId,
			network,
			controlByte,
			port,
			data))

		// The Piconet firmware adds the Source Station and Net bytes to the reply.
		sReply = fmt.Sprintf("TX %d %d %d %d %s\r",
			stationId,
			network,
			controlByte,
			port, encData)

		if len(extraScoutData) > 0 {
			encScoutExtraData := Base64Encode(extraScoutData)
			sReply += " " + encScoutExtraData
		}

		if err = commsClient.Write([]byte(sReply + "\r")); err != nil {
			slog.Error(err.Error())
		}

	}
}
func Broadcast(commsClient CommunicationClient, data []byte) {

	var err error

	if commsClient != nil {

		encData := Base64Encode(data)

		slog.Info(fmt.Sprintf("piconet-command=BCAST, data=[% 02X]", data))

		sReply := fmt.Sprintf("BCAST %s\r", encData)
		if err = commsClient.Write([]byte(sReply)); err != nil {
			slog.Error(err.Error())
		}
	}
}
