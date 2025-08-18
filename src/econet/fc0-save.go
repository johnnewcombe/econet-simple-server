package econet

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/johnnewcombe/econet-simple-server/src/fs"
	"github.com/johnnewcombe/econet-simple-server/src/lib"
)

func f0Save(cmd CliCmd, srcStationId byte, srcNetworkId byte) (*FSReply, error) {

	var (
		reply *FSReply
		err   error
		fd    *fs.FileDescriptor
	)

	if !ActiveSessions.IsLoggedOn(srcStationId, srcNetworkId) {
		// TODO Reply with Who Are You? instead of an error
		return nil, fmt.Errorf("econet-f0-save: user not authenticated")
	}

	slog.Info(fmt.Sprintf("econet-f0-save: src-stn=%02X, src-net=%02X, data=[% 02X]", srcStationId, srcNetworkId, cmd.ToBytes()))

	if fd, err = createFileDescriptor(cmd); err != nil {
		return nil, err
	}

	reply = NewFSReply(CCSave, RCOk, fd.ToBytes())
	return reply, nil
}

func createFileDescriptor(cmd CliCmd) (*fs.FileDescriptor, error) {

	var (
		fd fs.FileDescriptor
	)
	argCount := len(cmd.Args)

	// sort the first arg out
	if argCount > 1 {

		// TODO need to handle following syntax. This can be done by further splitting Args[1] by a "+"

		//	Possible SAVE command syntax
		//
		//     *SAVE MYDATA 3000+500
		//     *SAVE MYDATA 3000 3500
		//     *SAVE BASIC C000+1000 C2B2      // adds execution address OF C2B2
		//     *SAVE PROG 3000 3500 5050 5000  // adds execution address and load address

		// TODO check that arg[0] is a valid filename
		fd = fs.FileDescriptor{
			Name: cmd.Args[0],
		}

		if strings.Contains(cmd.Args[1], "+") {

			// we have the length specified
			arg := strings.Split(cmd.Args[1], "+")

			// get the start address and length
			fd.StartAddress = lib.StringToUint32(arg[0])
			fd.Size = lib.StringToUint32(arg[1])

			if argCount > 2 {
				fd.ExecuteAddress = lib.StringToUint32(cmd.Args[2])
			}

			// load address updates the start address
			if argCount > 3 {
				fd.StartAddress = lib.StringToUint32(cmd.Args[3])
			}

		} else {
			// just the start address
			fd.StartAddress = lib.StringToUint32(cmd.Args[1])
			if argCount > 2 {
				fd.Size = lib.StringToUint32(cmd.Args[2])
			}
		}
		if argCount > 3 {
			fd.ExecuteAddress = lib.StringToUint32(cmd.Args[3])
		}

		// load address updates the start address
		if argCount > 4 {
			fd.StartAddress = lib.StringToUint32(cmd.Args[4])
		}

	} else {
		return nil, fmt.Errorf("econet-f0-save: invalid number of arguments")
	}

	return &fd, nil
}
