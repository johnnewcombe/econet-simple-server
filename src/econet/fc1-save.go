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

var FileXfer *fs.FileTransfer // used to persist Data about the current file transfer

func fc1Save(srcStationId byte, srcNetworkId byte, port byte, data []byte) (*FSReply, error) {

	var (
		reply     *FSReply
		session   *Session
		replyPort byte
		localPath string
		filename  string
		diskname  string
		err       error
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
		//		slog.Info(fmt.Sprintf("econet-f1-save: src-stn=%02X, src-net=%02X, return-code=%s",
		//			srcStationId, srcNetworkId, string(ReplyCodeMap[RCWhoAreYou])))

		return reply, nil
	}

	if port == ServerPort {

		// needs to be at least 15 chars
		if len(data) < 15 {

			// TODO should this be a Reply? and an error or just a reply
			// error
			return nil, fmt.Errorf("econet-f0-save: not enough Data received")
		}

		// create a file transfer object to keep track of stuff, the Data received is passes in as a parameter and this
		// is parsed and used to populate the object
		//FileXfer = fs.NewFileTransferOld(byte(FCSave), replyPort, Data[5:])

		// get the filename element from data
		filename = strings.Split(string(data[16:]), "\r")[0]

		// expand filename to full name as specified from $ (root), the diskname is returned
		// separately
		if filename, diskname, err = session.ExpandEconetPath(filename); err != nil {
			return nil, err
		}

		// not hat the current disk makes no difference as each user has space on each disk.
		if !fs.IsOwner(filename, session.Username) {

		}

		FileXfer = fs.NewFileTransfer(byte(FCSave), replyPort,
			lib.LittleEndianBytesToInt(data[5:9]),
			lib.LittleEndianBytesToInt(data[9:13]),
			lib.LittleEndianBytesToInt(data[13:16]),
			filename, diskname,
		)

		if FileXfer == nil {
			return nil, fmt.Errorf("econet-f0-save: could not create file transfer object")
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
			// TODO inorder to check the handles the filename will need expanding to include the disk and directory etc
			if localPath, err = session.EconetPathToLocalPath(FileXfer.Filename); err != nil {
				//TODO reply with CORRECT error
				reply = NewFSReply(replyPort, CCIam, RCBadCommmand, ReplyCodeMap[RCBadCommmand])
				return nil, err
			}

			// add the attributes
			// TODO should this be handled inside the FileXferObject and applied to the EconetPath
			localPath = fmt.Sprintf("%s_%4X_%4X_%2X",
				localPath,
				FileXfer.StartAddress,
				FileXfer.ExecuteAddress,
				accessByte)

			// TODO handle these
			//  Insufficient Access to directory (are there directory permissions?)
			//	  does user own parent directory
			//  Too many open files
			//  Max handles
			//  Max files on serer
			//  No free network ports
			//  Server error unable to open file for writing
			//  If File Exists:
			//    is object locked "Access Violation"
			//    is object ia directory?
			// 	  User does not have write access

			// check for an open handle (file may exist)
			if !session.HandleExists(FileXfer.Filename) {
				// all good so add the handle
				session.AddHandle(FileXfer.Filename, File, false)
			} else {
				//TODO reply with CORRECT error
				reply = NewFSReply(replyPort, CCIam, RCInsufficientAccess, ReplyCodeMap[RCInsufficientAccess])
				return nil, fmt.Errorf("econet-f1-save: cannot save, file exists and is open")
			}

			// all good so create/overwrite the file
			if err = lib.WriteBytes(localPath, FileXfer.FileData); err != nil {
				//TODO reply with CORRECT error
				reply = NewFSReply(replyPort, CCIam, RCInsufficientAccess, ReplyCodeMap[RCInsufficientAccess])
				return nil, err
			}

			reply = NewFSReply(FileXfer.ReplyPort, CCComplete, RCOk, []byte{accessByte, fileCreationDate[0], fileCreationDate[1]})

		} else {
			reply = NewFSReply(replyPort, CCIam, RCTooMuchDataSentFromClient, ReplyCodeMap[RCTooMuchDataSentFromClient])
			return nil, fmt.Errorf("econet-f1-save: too much data received")
		}
	}

	return reply, nil
}
