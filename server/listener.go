package server

import (
	"strings"
)

//import "github.com/johnnewcombe/econet-simple-server/logger"

func Listener(ch chan byte) {
	var (
		b byte
	)
	for {
		// TODO: need to listen to channel from commClient and process
		b = <-ch

		s := strings.Builder{}
		s.WriteByte(b)

		//logger.LogDebug.Printf("RX: %s", string(b))
	}
}
