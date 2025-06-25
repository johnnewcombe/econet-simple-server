package piconet

import (
	"fmt"
	"github.com/johnnewcombe/econet-simple-server/comms"
	"github.com/johnnewcombe/econet-simple-server/logger"
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
		commsClient.Write([]byte(fmt.Sprintf("SET_STATION %d\r", stationID)))
	}
}

func SetMode(commsClient comms.CommunicationClient, mode string) {

	if commsClient != nil {
		if piconetMode[mode] {
			commsClient.Write([]byte(fmt.Sprintf("SET_STATION %s\r", mode)))
		} else {
			logger.LogError.Println(fmt.Errorf("invalid mode: %s", mode))
		}
	}
}

func GetStatus(commsClient comms.CommunicationClient) {

	if commsClient != nil {
		commsClient.Write([]byte("STATUS\r"))
	}
}
