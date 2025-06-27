package server

import (
	"github.com/johnnewcombe/econet-simple-server/econetCommands"
	"strings"
)

//import "github.com/johnnewcombe/econet-simple-server/logger"

func Listener(ch chan byte) {

	s := strings.Builder{}

	for {
		// TODO: this will block the next byte if not collected quickly
		b := <-ch

		s.WriteByte(b)
		if b == 0x0d || b == 0x10 {
			econetCommands.ParseCommand(s.String())
			s = strings.Builder{}
		}

		//logger.LogDebug.Printf(" %s", string(b))
	}
}
