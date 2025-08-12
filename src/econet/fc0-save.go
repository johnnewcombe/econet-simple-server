package econet

import (
	"fmt"
	"log/slog"
)

func fc1save(srcStationId byte, srcNetworkId byte, data []byte) (*FSReply, error) {
	var (
		reply      *FSReply
		session    *Session
		returnCode string
	)

	slog.Info(fmt.Sprintf("econet-f1-save: src-stn=%02X, src-net=%02X, data=[% 02X]", srcStationId, srcNetworkId, data))

	// get the logged on status
	session = ActiveSessions.GetSession(srcStationId, srcNetworkId)
	if session != nil {

		// user logged on

		// 00 C0 00 00 00 C0 00 00 10 00 00 44 41 54 41 0D

		//SEE aun-filestore, FileServer.php Line 1466

		// needs to be at least 13 chars to include a one letter filename followed by CR
		if len(data) < 13 {
			// error
			// TODO not enough data return the correct error
			returnCode = "SOME ERROR OR ANOTHER"
			//TODO change the reply code to something suitable
			reply = NewFSReply(CCIam, RCWrongPassword, ReplyCodeMap[RCWrongPassword])

		} else {

			// TODO Create a file and/or handle or something
			//loadAddress := lib.LittleEndianBytesToInt(data[:4])
			//execAddress := lib.LittleEndianBytesToInt(data[4:8])
			//length := lib.LittleEndianBytesToInt(data[8:11])
			//filename := lib.LittleEndianBytesToInt(data[11:16])

			//print(loadAddress)
			//print(execAddress)
			//print(length)
			//print(filename)

			// already logged in
			//TODO is this correct behaviour i.e. if we are already logged on from this station then
			// just say OK and keep current session? Or do we remove old session and create a new one
			//slog.Info(fmt.Sprintf("FC0 CLI Decoding, econet-command=I AM %s, authenticated=%v, return-code=OK", username, authenticated))

			// TODO Change FSReply as the Command code is only used for FunctionCode 0 calls
			reply = NewFSReply(CCSave, RCOk, []byte{
				DefaultUserRootDirHandle,
				DefaultCurrentDirectoryHandle,
				DefaultCurrentLibraryHandle,
				DefaultBootOption,
			})

			returnCode = "OK"
		}

	} else {
		// TODO Fix me
		//$oReply->setError(0xbf,"Who are you?");
	}

	slog.Info(fmt.Sprintf("fc1-save: return-code=%s", returnCode))

	return reply, nil

}
