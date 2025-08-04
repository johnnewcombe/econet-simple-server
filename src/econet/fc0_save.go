package econet

import "errors"

func f0_Save(cmd CliCmd, srcStationId byte, srcNetworkId byte) (*FSReply, error) {

	reply := NewFSReply(CCIam, WrongPassword, []byte("NOT IMPLEMENTED\r"))
	return reply, errors.New("not implemented")
}
