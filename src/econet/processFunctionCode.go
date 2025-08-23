package econet

import (
	"errors"
)

// func ProcessFunctionCode(functionCode byte, port byte, data []byte, srcStationId byte, srcNetworkId byte) (*FSReply, error) {
func ProcessFunctionCode(srcStationId byte, srcNetworkId byte, functionCode byte, receivePort byte, data []byte) (*FSReply, error) {

	var (
		replyPort byte
		reply     *FSReply
		err       error
	)

	replyPort = data[0] // this is ignored for data blocks that don't send a reply port

	switch functionCode {
	case 0:
		// tidy the command string
		reply, err = fc0CliDecode(srcStationId, srcNetworkId, data)
		break
	case 1:
		reply, err = fc1Save(srcStationId, srcNetworkId, receivePort, data)
		break
	case 2:
	case 3:
	case 4:
	case 5:
	case 6:
	case 7:
	case 8:
	case 9:
	case 10:
	case 11:
	case 12:
	case 13:
	case 14:
	case 15:
	case 16:
	case 17:
	case 18:
	case 19:
	case 20:
	case 21:
	case 22:
	case 23:
	case 24:
	case 25:
	case 26:
	case 27:
	case 28:
	case 29:
	case 30:
	case 31:
	case 32:
	case 33:
	case 34:
	case 35:
	case 36:
	case 37:
	case 38:
	case 39:
	case 40:
	case 41:
	case 42:
	case 43:
	case 44:
	case 45:
	case 46:
	default:
		reply = NewFSReply(replyPort, CCIam, RCBadCommmand, ReplyCodeMap[RCBadCommmand])
		err = errors.New("bad command")
		break
	}

	if reply == nil {
		return nil, err
	}

	return reply, err
}
