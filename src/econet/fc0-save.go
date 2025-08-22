package econet

import (
	"fmt"
	"log/slog"

	"github.com/johnnewcombe/econet-simple-server/src/fs"
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

	slog.Info(fmt.Sprintf("econet-f0-save: src=%02X/%02X, cmd-text=%s, data=[% 02X]", srcStationId, srcNetworkId, cmd.ToString(), cmd.ToBytes()))

	if fd, err = fs.NewFileDescriptor(cmd.Args); err != nil {
		return nil, err
	}

	reply = NewFSReply(CCSave, RCOk, fd.ToBytes())
	return reply, nil
}
