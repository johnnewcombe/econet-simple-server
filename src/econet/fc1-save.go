package econet

import (
	"fmt"
	"log/slog"

	"github.com/johnnewcombe/econet-simple-server/src/fs"
)

var FileXfer *fs.FileTransfer // used to persist data about the current file transfer

func fc1Save(srcStationId byte, srcNetworkId byte, port byte, data []byte) (*FSReply, error) {

	var (
		reply     *FSReply
		session   *Session
		replyPort byte
	)

	// port represents the port that the request was sent on, this allows us to determine if we
	// are in the data phase or not.
	slog.Info(fmt.Sprintf("econet-f1-save: src=%02X/%02X, port=%02X, data=[% 02X]",
		srcStationId, srcNetworkId, port, data))

	// normal reply port not the data acknowledge reply port
	// get the replyPort, this is not used for data block frames
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
			return nil, fmt.Errorf("econet-f0-save: not enough data received")
		}

		// create a file transfer object to keep track of stuff, the data received is passes in as a parameter and this
		// is parsed and used to populate the object
		FileXfer = fs.NewFileTransfer(byte(FCSave), replyPort, data[5:])

		if FileXfer == nil {
			return nil, fmt.Errorf("econet-f0-save: could not create file transfer object")
		}
		// capture from the data port that needs to be used to acknowledge future received data blocks
		FileXfer.DataAckPort = data[2]

		// send a reply to the client with the max block size
		replyData := []byte{
			DataPort,
			byte(MaxBlockSize % 256), // this needs to be calculated from q constant (little endian i.e.maxBlockSize=0x0500-1280
			byte(MaxBlockSize / 256),
		}

		// pad the filename to 12 chars and add to the reply

		replyData = append(replyData, []byte(FileXfer.Filename)...)
		reply = NewFSReply(replyPort, CCComplete, RCOk, replyData)

	} else if port == DataPort {

		// record the number of bytes received
		FileXfer.BytesTransferred += len(data)

		// check if we have received all the data
		if FileXfer.BytesTransferred < int(FileXfer.Size) {

			// return data block reply
			FileXfer.FileData = append(FileXfer.FileData, data...)

			reply = NewFsReplyData(FileXfer.DataAckPort)
			//slog.Warn("econet-f1-save: data-ack-port", dataAckPort)

		} else if FileXfer.BytesTransferred == int(FileXfer.Size) {

			// return final reply
			// TODO determine access byte and file creation date
			accessByte := byte(0b00010011)                     // unlocked, r/w for the owner and ro for others
			fileCreationDate := []byte{0b00001100, 0b10000011} //  12th March 1989

			reply = NewFSReply(FileXfer.ReplyPort, CCComplete, RCOk, []byte{0x00, accessByte, fileCreationDate[0], fileCreationDate[1]})

		} else {
			//TODO reply with error
			return nil, fmt.Errorf("econet-f1-save: too much data received")
		}

		// data transfer mode
		// store data in memory (extend file descriptor???)
		// check bytes received against size ( this needs to be stored outside of this function)
		// send short reply after each block
		// send long reply at end of file

		// a one-byte reply to acknowledge a block this will be sent on the data acknowledge port

		// this is a reply for all but the last block of data
		// TODO we need to keep track of the current filesize and compare to what is expected
		// i.e. we need to store the save activity status somewhere transient
		// we can only serve one client at a time so maybe this can just be stored at package level
		// we need to store size, filename, start, exec etc and the name and permissions that will be applied
		// all of this could be held in a file descriptor and saved with the file perhaps in the filename?

	}

	return reply, nil

}
