package econet

import (
	"log/slog"
	"strings"

	"github.com/johnnewcombe/econet-simple-server/src/lib"
)

// fc0cli Function code 0 CLI Decode

func f0Iam(cmd CliCmd, srcStationId byte, srcNetworkId byte, replyPort byte) (*FSReply, error) {
	var (
		password      string
		username      string
		returnCode    ReturnCode
		reply         *FSReply
		authenticated bool
		session       *Session
		err           error
	)

	slog.Info("econet-f0-iam:",
		"src-stn", srcStationId,
		"src-net", srcNetworkId,
		"reply-port", replyPort,
		"cmd", cmd.ToString())

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

	// if the station is logged on, log it off
	if session != nil {
		ActiveSessions.RemoveSession(session)
		//slog.Info("econet-f0-iam:", "previous session", "removed", "user", username)
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
		reply = NewFSReply(replyPort, CCIam, returnCode, ReplyCodeMap[returnCode])

	} else {

		// user exists so authenticate them
		if user := Userdata.AuthenticateUser(username, password); user != nil {

			// add the new session
			session = NewSession(*user, srcStationId, srcNetworkId)
			ActiveSessions.AddSession(session)

			// create user directories on the local file system
			if err = lib.CreateDirectoryIfNotExists(LocalRootDiectory + "." + user.Username); err != nil {
				return nil, err
			}
			if err = lib.CreateDirectoryIfNotExists(LocalRootDiectory + "." + user.Username); err != nil {
				return nil, err
			}
			if err = lib.CreateDirectoryIfNotExists(LocalRootDiectory + "." + user.Username); err != nil {
				return nil, err
			}
			if err = lib.CreateDirectoryIfNotExists(LocalRootDiectory + "." + user.Username); err != nil {
				return nil, err
			}

			returnCode = RCOk

			// TODO Fix this as  NewSession() also creates these handles
			reply = NewFSReply(replyPort, CCIam, returnCode, []byte{
				session.GetUrd(),
				session.GetCsd(),
				session.GetCsl(),
				session.BootOption,
			})

			authenticated = true

		} else {

			rc := RCWrongPassword
			reply = NewFSReply(replyPort, CCIam, rc, ReplyCodeMap[rc])

		}
	}

	slog.Info("econet-f0-iam:", "authenticated", authenticated, "user", username, "reply", ReplyCodeMap[returnCode])

	return reply, nil
}
