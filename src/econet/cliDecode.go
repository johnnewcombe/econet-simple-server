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
			DefaultUserRootDirHandle,
			DefaultCurrentDirectoryHandle,
			DefaultCurrentLibraryHandle,
			DefaultBootOption,
		})

		returnCode = "OK"
		authenticated = true

	}

	if !Userdata.UserExists(username) {

		/*
			//TODO Fixme.... FSErrorReply needed!
			The return code is an indication to the client of any error status which has
			arisen, as a result of attempting to execute the command. A return code of
			zero indicates that the command step completed successfully; otherwise the
			return code is the error number indicating what error has occurred. If the
			return code is non-zero, then the remainder of the message contains an ASCII
			string terminated by a carriage return, which describes the error.
		*/

		returnCode = "USER NOT KNOWN"
		reply = NewFSReply(CCIam, UserNotKnown, []byte(returnCode+"\r"))

	} else {

		// if logged on at this machine already then logg them off
		session = ActiveSessions.GetSession(username, srcStationId, srcNetworkId)
		if session != nil {
			ActiveSessions.RemoveSession(session)
		}

		// authenticate user
		if user := Userdata.AuthenticateUser(username, password); user != nil {

			returnCode = "OK"
			authenticated = true

			// add the new session
			session = ActiveSessions.AddSession(username, srcStationId, srcNetworkId)

			// note that these default handles are already set in the newly created
			// session object, as is the default boot option
			reply = NewFSReply(CCIam, RCOk, []byte{
				DefaultUserRootDirHandle,
				DefaultCurrentDirectoryHandle,
				DefaultCurrentLibraryHandle,
				DefaultBootOption,
			})

		} else {

			returnCode = "WRONG PASSWORD"
			// TODO Sort out correct responses for failed login
			reply = NewFSReply(CCIam, WrongPassword, []byte(returnCode+"\r"))

		}
	}

	slog.Info(fmt.Sprintf("FC0 CLI Decoding, econet-command=I AM %s, authenticated=%v, retturn-code=%s", username, authenticated, returnCode))

	return reply.ToBytes()
}
