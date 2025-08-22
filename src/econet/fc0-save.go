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

	slog.Info(fmt.Sprintf("econet-f0-save: src=%02X/%02X, cmd=%s, data=[% 02X]", srcStationId, srcNetworkId, cmd.ToString(), cmd.ToBytes()))

	if fd, err = createFileDescriptor(cmd); err != nil {
		return nil, err
	}

	reply = NewFSReply(CCSave, RCOk, fd.ToBytes())
	return reply, nil
}

func createFileDescriptor(cmd CliCmd) (*fs.FileDescriptor, error) {
	argCount := len(cmd.Args)
	if argCount < 2 {
		return nil, fmt.Errorf("econet-f0-save: invalid number of arguments")
	}

	fd := fs.FileDescriptor{Name: cmd.Args[0]}

	var (
		start uint32
		size  uint32
		exec  uint32
		load  uint32
	)

	if strings.Contains(cmd.Args[1], "+") {

		parts := strings.SplitN(cmd.Args[1], "+", 2)
		start = lib.StringToUint32(parts[0])
		size = lib.StringToUint32(parts[1])

		if argCount > 2 {
			exec = lib.StringToUint32(cmd.Args[2])
		} else {
			exec = start
		}

		if argCount > 3 {
			load = lib.StringToUint32(cmd.Args[3])
		} else {
			load = start
		}

	} else {

		if argCount < 3 {
			return nil, fmt.Errorf("econet-f0-save: invalid number cmd arguments")
		}

		start = lib.StringToUint32(cmd.Args[1])
		end := lib.StringToUint32(cmd.Args[2])
		size = end - start

		if argCount > 3 {
			exec = lib.StringToUint32(cmd.Args[3])
		} else {
			exec = start
		}

		if argCount > 4 {
			load = lib.StringToUint32(cmd.Args[4])
		} else {
			load = start
		}
	}

	// Load address updates the start address (preserve exec as per original logic)
	fd.StartAddress = load
	fd.Size = size
	fd.ExecuteAddress = exec

	return &fd, nil
}
