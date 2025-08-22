package econet

import (
	"fmt"
	"log/slog"

	"github.com/johnnewcombe/econet-simple-server/src/fs"
	"github.com/johnnewcombe/econet-simple-server/src/lib"
)

var fileXfer fs.FileTransfer // used to persist data about the current file transfer

func fc1Save(srcStationId byte, srcNetworkId byte, port byte, data []byte) (*FSReply, error) {

	var (
		reply       *FSReply
		session     *Session
		dataAckPort byte
	)

	// port represents the port that the request was sent on, this allows us to determine if we
	// are in the data phase or not.
	slog.Info(fmt.Sprintf("econet-f1-save: src=%02X/%02X, port=%02X, data=[% 02X]",
		srcStationId, srcNetworkId, port, data))

	// get the logged on status, we're not using .IsLoggedOn() here as we need the session later anyway
	session = ActiveSessions.GetSession(srcStationId, srcNetworkId)
	if session == nil {

		// user is not logged on so return 'who are you'
		reply = NewFSReply(CCComplete, RCWhoAreYou, ReplyCodeMap[RCWhoAreYou])
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

		// create a file transfer object to keep track of stuff
		// TODO get the name size start address, exec address from the dat
		fileXfer = fs.FileTransfer{
			Filename:       fmt.Sprintf("%-12s", "NOS"),
			StartAddress:   lib.StringToUint32(string(data[3:7])),
			ExecuteAddress: lib.StringToUint32(string(data[7:11])),
			Size:           lib.StringToUint32(string(data[11:14])),
			FileData:       []byte{},
		}

		// the data will give us the reply port. this needs to be stored perhaps in session?
		// as it will be checked for by the listener on each RX_TRANSMIT event
		dataAckPort = data[0]

		replyData := []byte{
			DataPort,
			byte(MaxBlockSize % 256), // this needs to be calculated from q constant (little endian i.e.maxBlockSize=0x0500-1280
			byte(MaxBlockSize / 256),
		}

		// pad the filename to 12 chars and add to the reply
		replyData = append(replyData, []byte(fileXfer.Filename)...)
		reply = NewFSReply(CCComplete, RCOk, replyData)

		//reply.data = append(reply.data, []byte("NOS\r")...)
		//slog.Info(fmt.Sprintf("econet-reply: src-stn=%02X, src-net=%02X, reply-code=%s, data=[% 02X]",
		//	srcStationId, srcNetworkId, string(ReplyCodeMap[reply.ReturnCode]), reply.ToBytes()))

	} else if port == DataPort {

		fileXfer.BytesTransferred += len(data)

		if fileXfer.BytesTransferred < int(fileXfer.Size) {

			// return data block reply
			fileXfer.FileData = append(fileXfer.FileData, data...)

			reply = NewFsReplyData([]byte{0x0})
			print(dataAckPort)

		} else if fileXfer.BytesTransferred == int(fileXfer.Size) {

			// return final reply
			// TODO determine access byte and file creation date
			accessByte := byte(0x00)
			fileCreationDate := []byte{0x00, 0x00, 0x00}
			reply = NewFSReply(CCComplete, RCOk, []byte{0x00, accessByte, fileCreationDate[0], fileCreationDate[1], fileCreationDate[2]})

		} else {
			//TODO reply with error
			return nil, fmt.Errorf("econet-f0-save: too much data received")
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
