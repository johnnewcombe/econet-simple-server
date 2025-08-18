package econet

import (
	"fmt"
	"log/slog"
)

func fc1Save(srcStationId byte, srcNetworkId byte, port byte, data []byte) (*FSReply, error) {

	var (
		reply       *FSReply
		session     *Session
		dataAckPort byte
	)

	// port represents the port that the request was sent on, this allows us to determine if we
	// are in the data phase or not.
	slog.Info(fmt.Sprintf("econet-f1-save: src-stn=%02X, src-net=%02X, port=%02X, data=[% 02X]",
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
			return nil, fmt.Errorf("econet-f0-save: not enough data")
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
		filename := fmt.Sprintf("%-12s", "NOS")
		replyData = append(replyData, []byte(filename)...)
		reply = NewFSReply(CCComplete, RCOk, replyData)

		//reply.data = append(reply.data, []byte("NOS\r")...)
		//slog.Info(fmt.Sprintf("econet-reply: src-stn=%02X, src-net=%02X, reply-code=%s, data=[% 02X]",
		//	srcStationId, srcNetworkId, string(ReplyCodeMap[reply.ReturnCode]), reply.ToBytes()))

	} else if port == DataPort {

		// data transfer mode
		// store data in memory (extend file descriptor???)
		// check bytes received against size ( this needs to be stored outside of this function)
		// send short reply after each block
		// send long reply at end of file

		// a one-byte reply to acknowledge a block this will be sent on the data acknowledge port
		reply = NewFsReplyData([]byte{0x0})
		print(dataAckPort)

		//reply = NewFSReply(CCComplete, RCOk, []byte{0x00})

	}

	return reply, nil

}
