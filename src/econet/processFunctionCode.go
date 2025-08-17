package econet

import (
	"errors"
)

func ProcessFunctionCode(functionCode byte, data []byte, srcStationId byte, srcNetworkId byte) (*FSReply, error) {

	var (
		reply *FSReply
		err   error
	)

	switch functionCode {
	case 0:
		// tidy the command string
		reply, err = fc0CliDecode(srcStationId, srcNetworkId, data)
		break
	case 1:
		reply, err = fc1Save(srcStationId, srcNetworkId, data)
		break
	case 2:
		break
	case 3:
		break
	case 4:
		break
	case 5:
		break
	case 6:
		break
	case 7:
		break
	case 8:
		break
	case 9:
		break
	case 10:
		break
	case 11:
		break
	case 12:
		break
	case 13:
		break
	case 14:
		break
	case 15:
		break
	case 16:
		break
	case 17:
		break
	case 18:
		break
	case 19:
		break
	case 20:
		break
	case 21:
		break
	case 22:
		break
	case 23:
		break
	case 24:
		break
	case 25:
		break
	case 26:
		break
	case 27:
		break
	case 28:
		break
	case 29:
		break
	case 30:
		break
	case 31:
		break
	case 32:
		break
	case 33:
		break
	case 34:
		break
	case 35:
		break
	case 36:
		break
	case 37:
		break
	case 38:
		break
	case 39:
		break
	case 40:
		break
	case 41:
		break
	case 42:
		break
	case 43:
		break
	case 44:
		break
	case 45:
		break
	case 46:
		break
	default:
		reply = NewFSReply(CCIam, RCBadCommmand, ReplyCodeMap[RCBadCommmand])
		err = errors.New("bad command")
		break
	}

	if reply == nil {
		return nil, err
	}

	return reply, err
}
