package econet

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/johnnewcombe/econet-simple-server/src/lib"
)

// fc0cli Function code 0 CLI Decode

func f0Iam(cmd CliCmd, srcStationId byte, srcNetworkId byte) (*FSReply, error) {
	var (
		//cmds     []string
		password      string
		username      string
		returnCode    ReturnCode
		reply         *FSReply
		authenticated bool
		session       *Session
		err           error
	)

	slog.Info(fmt.Sprintf("econet-f0-iam: src%02X/%02X, data=[% 02X], cmd=%s", srcStationId, srcNetworkId, cmd.ToString(), cmd.ToBytes()))

	// TODO need to sort out commands such as the following with or without passwords
	// the clients NFS probably handles all of this
	// I AM JOHN
	// I AM 247 JOHN
	// I AM 3.247 JOHN

	argCount := len(cmd.Args)

	if argCount > 0 {
		username = cmd.Args[0]

	}
	if argCount > 1 {
		// the password is everything upto the CR if there is one
		// note that the System 3 send additional bytes after the
		// password.
		password = strings.Split(cmd.Args[1], "\r")[0]
	}

	// get logged in status of the machine could this user or a previous one
	session = ActiveSessions.GetSession(srcStationId, srcNetworkId)

	// if station is logged on, log off
	if session != nil {
		ActiveSessions.RemoveSession(session)
		slog.Info(fmt.Sprintf("econet-f0-iam: econet-command=I AM %s, msg=previous session removed", username))
	}

	// check user
	if !Userdata.UserExists(username) {
		/*
			The return code is an indication to the client of any error status which has
			arisen, as a result of attempting to execute the command. A return code of
			zero indicates that the command step completed successfully; otherwise the
			return code is the error number indicating what error has occurred. If the
			return code is non-zero, then the remainder of the message contains an ASCII
			string terminated by a carriage return, which describes the error.
		*/

		returnCode = RCUserNotKnown
		reply = NewFSReply(CCIam, returnCode, ReplyCodeMap[returnCode])

	} else {

		// user exists so all good

		// authenticate user
		if user := Userdata.AuthenticateUser(username, password); user != nil {

			// add the new session
			session = ActiveSessions.AddSession(username, srcStationId, srcNetworkId)

			// TODO is it correct that the current selected dir will be the same as
			//  user root dir but have a separate handle?
			urd := DefaultRootDirectory + "." + Disk0 + user.Username
			csd := DefaultRootDirectory + "." + Disk0 + user.Username
			csl := DefaultRootDirectory + "." + DefaultLibraryDirectory

			if err = lib.CreateDirectoryIfNotExists(urd); err != nil {
				return nil, err
			}

			returnCode = RCOk
			reply = NewFSReply(CCIam, returnCode, []byte{
				session.AddHandle(urd, UserRootDirectory),
				session.AddHandle(csd, CurrentSelectedDirectory),
				session.AddHandle(csl, CurrentSelectedDirectory),
				session.BootOption,
			})

			authenticated = true

		} else {

			rc := RCWrongPassword
			reply = NewFSReply(CCIam, rc, ReplyCodeMap[rc])

		}
	}

	slog.Info(fmt.Sprintf("econet-f0-iam: econet-command=I AM %s, msg=authenticated=%v, return-code=%s", username, authenticated, ReplyCodeMap[returnCode]))

	return reply, nil
}
