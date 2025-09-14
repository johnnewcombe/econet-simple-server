package econet

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/johnnewcombe/econet-simple-server/src/fs"
	"github.com/johnnewcombe/econet-simple-server/src/lib"
)

const (
	defaultAccessByte byte = 0b00010011
)

// FileXfer is used to persist Data about the current file transfer
var FileXfer *fs.FileTransfer

// fc1-save is the function code for saving a file and is called by Acorn System and Atoms via f0-save
// and bt BBC computers and later, directly. It always returns a reply even in the case of an error
// so that the error can be propagated to the client.
func fc1Save(srcStationId byte, srcNetworkId byte, port byte, data []byte) (*FSReply, error) {

	var (
		reply      *FSReply
		session    *Session
		replyPort  byte
		localPath  string
		filename   string
		diskName   string
		returnCode ReturnCode
		err        error
	)

	// port represents the port that the request was sent on, this allows us to determine if we
	// are in the Data phase or not.
	slog.Info("econet-f0-save:",
		"src-stn", srcStationId,
		"src-net", srcNetworkId,
		"port", port)

	// normal reply port not the Data acknowledge reply port
	// get the replyPort, this is not used for Data block frames
	replyPort = data[0]

	// get the logged on status, we're not using .IsLoggedOn() here as we need the session later anyway
	session = ActiveSessions.GetSession(srcStationId, srcNetworkId)

	// return WHO ARE YOU if the user is not logged on
	if session == nil {

		// user is not logged on so return 'who are you'
		reply = NewFSReply(replyPort, CCComplete, RCWhoAreYou, ReplyCodeMap[RCWhoAreYou])
		return reply, nil
	}

	if port == ServerPort {

		// needs to be at least 15 chars
		if len(data) < 15 {

			reply = NewFSReply(replyPort, CCIam, RCBadCommmand, ReplyCodeMap[RCBadCommmand])
			return reply, fmt.Errorf("not enough data received")
		}

		// create a file transfer object to keep track of stuff, the Data received is passes in as a parameter and this
		// is parsed and used to populate the object
		//FileXfer = fs.NewFileTransferOld(byte(FCSave), replyPort, Data[5:])

		// get the filename element from data
		filename = strings.Split(string(data[16:]), "\r")[0]

		// expand filename to full name as specified from $ (root), the diskName is returned
		// separately
		if filename, diskName, err = session.ExpandEconetPath(filename); err != nil {
			reply = NewFSReply(replyPort, CCIam, RCBadFileName, ReplyCodeMap[RCBadFileName])
			return reply, err
		}

		// not hat the current disk makes no difference as each user has space on each disk.
		// TODO: Check what is the difference between RCInsufficientAccess ans RCInsufficientPrivilege
		//  is in this case, check with BBC Level 3 server
		if !fs.IsOwner(filename, session.User.Username) && !session.User.IsPrivileged {
			returnCode = RCInsufficientAccess
			reply = NewFSReply(replyPort, CCSave, returnCode, ReplyCodeMap[returnCode])
		}

		// the file transfer object is created here and parses the data allowing simple access
		// to the file transfer parameters
		FileXfer = fs.NewFileTransfer(byte(FCSave), replyPort,
			lib.LittleEndianBytesToInt(data[5:9]),
			lib.LittleEndianBytesToInt(data[9:13]),
			lib.LittleEndianBytesToInt(data[13:16]),
			filename, diskName,
		)

		if FileXfer == nil {
			reply = NewFSReply(replyPort, CCIam, RCBadCommmand, ReplyCodeMap[RCBadCommmand])
			return reply, fmt.Errorf("could not create file transfer object, bad command")
		}
		// capture from the Data port that needs to be used to acknowledge future received Data blocks
		FileXfer.DataAckPort = data[2]

		// send a reply to the client with the max block size
		replyData := []byte{
			DataPort,
			byte(MaxBlockSize % 256), // this needs to be calculated from q constant (little endian i.e.maxBlockSize=0x0500-1280
			byte(MaxBlockSize / 256),
		}

		// get the filename leaf name and pad to 12 chars
		// this is only needed for function code 3 (*CAT)
		//replyData = append(replyData, []byte(FileXfer.GetLeafName())...)
		reply = NewFSReply(replyPort, CCComplete, RCOk, replyData)

	} else if port == DataPort {

		// record the number of bytes received
		FileXfer.BytesTransferred += len(data)

		// check if we have received all the Data
		if FileXfer.BytesTransferred < int(FileXfer.Size) {

			// return Data block reply
			FileXfer.FileData = append(FileXfer.FileData, data...)

			reply = NewFsReplyData(FileXfer.DataAckPort)
			//slog.Warn("econet-f1-save: Data-ack-port", dataAckPort)

		} else if FileXfer.BytesTransferred == int(FileXfer.Size) {

			// return final reply
			// TODO: set the file creation date from current time/date
			accessByte := defaultAccessByte                    // unlocked, r/w for the owner and ro for others
			fileCreationDate := []byte{0b00001100, 0b10000011} //  12th March 1989

			// all good so save the file
			// save the file
			if localPath, err = session.EconetPathToLocalPath(FileXfer.Filename); err != nil {
				reply = NewFSReply(replyPort, CCIam, RCBadFileName, ReplyCodeMap[RCBadFileName])
				return reply, err
			}

			// add the attributes
			localPath = fmt.Sprintf("%s_%4X_%4X_%2X",
				localPath,
				FileXfer.StartAddress,
				FileXfer.ExecuteAddress,
				accessByte)

			// check for an open handle (the file may exist)
			if !session.HandleExists(FileXfer.DiskName, FileXfer.Filename) {
				// add the handle
				if _, err = session.AddHandle(FileXfer.DiskName, FileXfer.Filename, File, false); err != nil {
					reply = NewFSReply(replyPort, CCIam, RCTooManyOpenFiles, ReplyCodeMap[RCTooManyOpenFiles])
					return reply, fmt.Errorf("cannot save, no file hanles available")
				}
			} else {
				// file open for read or write so cannot save
				reply = NewFSReply(replyPort, CCIam, RCObjectInUse, ReplyCodeMap[RCObjectInUse])
				return reply, fmt.Errorf("cannot save, file exists and is open")
			}

			// TODO handle these all
			//  No free network ports
			//  If File Exists:
			//    is object locked "Access Violation"
			//    is object ia directory?
			// 	  PWEntry does not have write access

			// create/overwrite the file
			if err = lib.WriteBytes(localPath, FileXfer.FileData); err != nil {
				reply = NewFSReply(replyPort, CCIam, RCDiscFault, ReplyCodeMap[RCDiscFault])
				return reply, err
			}

			reply = NewFSReply(FileXfer.ReplyPort, CCComplete, RCOk, []byte{accessByte, fileCreationDate[0], fileCreationDate[1]})

		} else {
			reply = NewFSReply(replyPort, CCIam, RCTooMuchDataSentFromClient, ReplyCodeMap[RCTooMuchDataSentFromClient])
			return reply, fmt.Errorf("too much data received")
		}
	}

	return reply, nil
}
