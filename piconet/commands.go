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
		commsClient.Write([]byte(fmt.Sprintf("SET_STATION %d\r\n", stationID)))
	}
}

func SetMode(commsClient comms.CommunicationClient, mode string) {

	if commsClient != nil {
		if piconetMode[mode] {
			commsClient.Write([]byte(fmt.Sprintf("SET_STATION %s\r\n", mode)))
		} else {
			logger.LogError.Printf(fmt.Errorf("invalid mode: %s\r\n", mode).Error())
		}
	}
}

func GetStatus(commsClient comms.CommunicationClient) {

	if commsClient != nil {
		commsClient.Write([]byte("STATUS\r"))
	}
}
