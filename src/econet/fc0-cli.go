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
		cmd   CliCmd
	)

	cmd = parseCommand(command)

	// these are all * commands
	switch cmd.Cmd {

	case "SAVE":
		break
	case "LOAD":
		break
	case "CAT":
		break
	case "INFO":
		break
	case "I AM":
		reply = iAm(cmd, srcStationId, srcNetworkId)
		break
	case "SDISK":
		break
	case "DIR":
		break
	case "LIB":
		break
	default:

	}

	return reply
}

func iAm(cmd CliCmd, srcStationId byte, srcNetworkId byte) []byte {
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
		password = cmd.Args[1]
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

			//handles[DefaultUserRootDirHandle] = DefaultRootDirectory
			//handles[DefaultCurrentDirectoryHandle] = DefaultRootDirectory + "." + username
			//handles[DefaultCurrentLibraryHandle] = DefaultRootDirectory + "." + DefaultLibraryDirectory

			reply = NewFSReply(CCIam, RCOk, []byte{
				session.AddHandle(DefaultRootDirectory),
				session.AddHandle(DefaultRootDirectory + "." + username),
				session.AddHandle(DefaultRootDirectory + "." + DefaultLibraryDirectory),
				session.BootOption,
			})

		} else {

			returnCode = "WRONG PASSWORD"
			reply = NewFSReply(CCIam, WrongPassword, []byte(returnCode+"\r"))

		}
	}

	slog.Info(fmt.Sprintf("FC0 CLI Decoding, econet-command=I AM %s, authenticated=%v, return-code=%s", username, authenticated, returnCode))

	return reply.ToBytes()
}

func parseCommand(commandText string) CliCmd {

	var (
		commands []string
		cmdArgs  []string
		cmd      string
		ok       bool
		argText  string
	)

	// TODO Reinstate this ?
	//commandText = strings.ToUpper(commandText)

	// list of piconet * commands
	commands = []string{"SAVE", "LOAD", "CAT", "INFO", "I AM", "SDISK", "DIR", "LIB"}

	for _, cmd = range commands {
		if _, argText, ok = strings.Cut(commandText, cmd); ok { // i.e. if ok
			cmdArgs = strings.Split(strings.Trim(argText, " "), " ")
			break
		}
	}

	return CliCmd{
		Cmd:     cmd,
		CmdText: commandText,
		Args:    cmdArgs,
	}
}
