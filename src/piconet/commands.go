package piconet

import (
	"fmt"
	"github.com/johnnewcombe/econet-simple-server/src/comms"
	"log/slog"
)

const (
	MODE_LISTEN  = "LISTEN"
	MODE_MONITOR = "MONITOR"
	MODE_STOP    = "STOP"
)

// TODO Need to be able to tell if client is connected.
var piconetMode = map[string]bool{
	"LISTEN":  true,
	"MONITOR": true,
	"STOP":    true,
}

func SetStationID(commsClient comms.CommunicationClient, stationID int) {

	if commsClient != nil {
		if err := commsClient.Write([]byte(fmt.Sprintf("SET_STATION %d\r", stationID))); err != nil {
			slog.Error(err.Error())
		}
	}
}

func SetMode(commsClient comms.CommunicationClient, mode string) {

	if commsClient != nil {
		if piconetMode[mode] {
			if err := commsClient.Write([]byte(fmt.Sprintf("SET_MODE %s\r", mode))); err != nil {
				slog.Error(err.Error())
			}
		} else {
			slog.Error("invalid mode", "mode", mode)
		}
	}
}

func Status(commsClient comms.CommunicationClient) {

	if commsClient != nil {
		if err := commsClient.Write([]byte("STATUS\r")); err != nil {
			slog.Error(err.Error())
		}
	}
}

func Restart(commsClient comms.CommunicationClient) {
	if commsClient != nil {
		if err := commsClient.Write([]byte("RESTART\r")); err != nil {
			slog.Error(err.Error())
		}
	}
}

func Transmit(commsClient comms.CommunicationClient, stationId byte, network byte, controlByte byte, port byte, data []byte) {
	if commsClient != nil {

	}
}
func Broadcast(commsClient comms.CommunicationClient, data []byte) {
	if commsClient != nil {

	}
}
