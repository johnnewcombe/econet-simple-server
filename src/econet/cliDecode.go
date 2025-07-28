package econet

import (
	"fmt"
	"log/slog"
	"strings"
)

// CLIDecode Function code 0 CLI Decode
func CLIDecode(command string, srcStationId byte, srcNetworkId byte) []byte {

	// PROCESS RX_TRANSMIT
	var (
		reply []byte
	)

	if strings.HasPrefix(command, "I AM") {
		// parse the command
		reply = iAm(command, srcStationId, srcNetworkId)

	} else if strings.HasPrefix(command, "CAT") {
	}

	return reply
}

func iAm(command string, srcStationId byte, srcNetworkId byte) []byte {
	var (
		//cmds     []string
		password      string
		username      string
		returnCode    string
		reply         *FSReply
		authenticated bool
		session       *Session
	)

	// get logged in status
	// TODO: This needs to be a user session as well as machine
	session = ActiveSessions.GetSession(username, srcStationId, srcNetworkId)

	//parse the command itself
	args := strings.Split(command, " ")
	if len(args) == 3 {
		username = args[2]
		password = ""
	} else if len(args) == 4 {
		username = args[2]
		password = args[3]
	}

	// check user against users
	if session != nil {
		// already logged in
		//TODO is this correct behaviour i.e. if we are already logged on from this station then
		// just say OK and keep current session? Or do we remove old session and create a new one
		slog.Info(fmt.Sprintf("FC0 CLI Decoding, econet-command=I AM %s, authenticated=%v, return-code=OK", username, authenticated))
		reply = NewFSReply(CCIam, RCOk, []byte{
			userRootDir,
			currentSelectedDirectory,
			currentSelectedLibrary,
			bootOption,
		})

		returnCode = "OK"
		authenticated = true

	} else if !Userdata.UserExists(username) {

		reply = NewFSReply(CCIam, UserNotKnown, []byte{
			userRootDir,
			currentSelectedDirectory,
			currentSelectedLibrary,
			bootOption,
		})

		returnCode = "USER NOT KNOWN"

	} else {
		if user := Userdata.AuthenticateUser(username, password); user != nil {
			// user good
			reply = NewFSReply(CCIam, RCOk, []byte{
				userRootDir,
				currentSelectedDirectory,
				currentSelectedLibrary,
				bootOption,
			})

			//reply = []byte{byte(CCIam),
			//	0x00,
			//	0x01,
			//	0x02,
			//	0x04,
			//	0x00}

			returnCode = "OK"
			authenticated = true

			// add the new session
			ActiveSessions.Sessions = append(ActiveSessions.Sessions, *NewSession(username, srcStationId, srcNetworkId, userRootDir, currentSelectedDirectory, currentSelectedLibrary))

		} else {
			// TODO Sort out correct responses for failed login
			reply = NewFSReply(CCIam, WrongPassword, []byte{
				userRootDir,
				currentSelectedDirectory,
				currentSelectedLibrary,
				bootOption,
			})
			returnCode = "WRONG PASSWORD"
		}
	}

	slog.Info(fmt.Sprintf("FC0 CLI Decoding, econet-command=I AM %s, authenticated=%v, retturn-code=%s", username, authenticated, returnCode))

	return reply.ToBytes()
}
