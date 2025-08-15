package econet

import (
	"fmt"
	"log/slog"
)

func fc1Save(srcStationId byte, srcNetworkId byte, data []byte) (*FSReply, error) {
	var (
		reply   *FSReply
		session *Session
		//returnCode string
	)

	slog.Info(fmt.Sprintf("econet-f1-save: src-stn=%02X, src-net=%02X, data=[% 02X]", srcStationId, srcNetworkId, data))

	// get the logged on status not using .IsLoggedOn() as we need the session later
	session = ActiveSessions.GetSession(srcStationId, srcNetworkId)
	if session == nil {
		// TODO Reply with Who Are You? instead of an error
		return nil, fmt.Errorf("econet-f0-save: user not authenticated")
	}

	// user logged on

	// 00 C0 00 00 00 C0 00 00 10 00 00 44 41 54 41 0D

	//SEE aun-filestore, FileServer.php Line 1466

	// needs to be at least 13 chars to include a one letter filename followed by CR
	if len(data) < 13 {
		// error
		return nil, fmt.Errorf("econet-f0-save: not enough data")
	}

	// the data will give us the reply port. this needs to be stored perhaps in session?
	// as it will be checked for by the listener on each RX_TRANSMIT event
	session.DataPort = data[0]

	// already logged in
	//TODO is this correct behaviour i.e. if we are already logged on from this station then
	// just say OK and keep current session? Or do we remove old session and create a new one
	//slog.Info(fmt.Sprintf("FC0 CLI Decoding, econet-command=I AM %s, authenticated=%v, return-code=OK", username, authenticated))

	// TODO Change FSReply as the Command code is only used for FunctionCode 0 calls
	reply = NewFSReply(CCSave, RCOk, []byte{
		session.DataPort,
		// max block size
		// file eaf name terminated by CR
	})

	// these are confusing
	//returnCode = "OK"
	//slog.Info(fmt.Sprintf("fc1-save: return-code=%s", returnCode))

	return reply, nil

}
