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
		if err := commsClient.Write([]byte(fmt.Sprintf("SET_STATION %d\r", stationID))); err != nil {
			logger.LogError.Println(err)
		}
	}
}

func SetMode(commsClient comms.CommunicationClient, mode string) {

	if commsClient != nil {
		if piconetMode[mode] {
			if err := commsClient.Write([]byte(fmt.Sprintf("SET_MODE %s\r", mode))); err != nil {
				logger.LogError.Println(err)
			}
		} else {
			logger.LogError.Printf(fmt.Errorf("invalid mode: %s\r", mode).Error())
		}
	}
}

func Status(commsClient comms.CommunicationClient) {

	if commsClient != nil {
		if err := commsClient.Write([]byte("STATUS\r")); err != nil {
			logger.LogError.Println(err)
		}
	}
}
