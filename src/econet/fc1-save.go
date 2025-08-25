package econet

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/johnnewcombe/econet-simple-server/src/fs"
	"github.com/johnnewcombe/econet-simple-server/src/lib"
)

var FileXfer *fs.FileTransfer // used to persist Data about the current file transfer

func fc1Save(srcStationId byte, srcNetworkId byte, port byte, data []byte) (*FSReply, error) {

	var (
		reply     *FSReply
		session   *Session
		replyPort byte
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

	// TODO handle these
	//  is object locked "Access Violation"
	//  object is a directory
	//  Insufficient Access
	// 		User does not have write access to existing file
	//		File does not exist, and user does not own parent directory
	//  Too many open files
	//      Max handles
	//      Max files on serer
	//      No free network ports
	//  Server error unable to open file for writing

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

		FileXfer = fs.NewFileTransfer(byte(FCSave), replyPort,
			lib.LittleEndianBytesToInt(data[5:9]),
			lib.LittleEndianBytesToInt(data[9:13]),
			lib.LittleEndianBytesToInt(data[13:16]),
			strings.Split(string(data[16:]), "\r")[0],
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
		replyData = append(replyData, []byte(FileXfer.GetLeafName())...)
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
			// TODO determine access byte and file creation date
			accessByte := byte(0b00010011)                     // unlocked, r/w for the owner and ro for others
			fileCreationDate := []byte{0b00001100, 0b10000011} //  12th March 1989

			// all good so save the file
			// TODO save the file
			//lib.WriteBytes(FileXfer.Filename, FileXfer.FileData)

			reply = NewFSReply(FileXfer.ReplyPort, CCComplete, RCOk, []byte{0x00, accessByte, fileCreationDate[0], fileCreationDate[1]})

		} else {
			//TODO reply with error
			return nil, fmt.Errorf("econet-f1-save: too much Data received")
		}

		// Data transfer mode
		// store Data in memory (extend file descriptor???)
		// check bytes received against size ( this needs to be stored outside of this function)
		// send short reply after each block
		// send long reply at end of file

		// a one-byte reply to acknowledge a block this will be sent on the Data acknowledge port

		// this is a reply for all but the last block of Data
		// TODO we need to keep track of the current filesize and compare to what is expected
		// i.e. we need to store the save activity status somewhere transient
		// we can only serve one client at a time so maybe this can just be stored at package level
		// we need to store size, filename, start, exec etc and the name and permissions that will be applied
		// all of this could be held in a file descriptor and saved with the file perhaps in the filename?

	}

	return reply, nil

}
