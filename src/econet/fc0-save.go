package econet

import (
	"log/slog"

	"github.com/johnnewcombe/econet-simple-server/src/fs"
)

func f0Save(cmd CliCmd, srcStationId byte, srcNetworkId byte, replyPort byte) (*FSReply, error) {

	var (
		reply *FSReply
		err   error
		fd    *fs.FileInfo
	)

	// don't need to ensure we are logged on as that will take place in f1-save anyway

	slog.Info("econet-f0-save:",
		"src-stn", srcStationId,
		"src-net", srcNetworkId,
		"cmd-text", cmd.ToString())

	// the fileDescriptor will parse the args and gives us an easy way to return them as a byte slice
	if fd, err = fs.NewFileInfo(cmd.Args); err != nil {
		return nil, err
	}

	reply = NewFSReply(replyPort, CCSave, RCOk, fd.ToBytes())
	return reply, nil
}
