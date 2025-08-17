package econet

import (
	"fmt"
	"log/slog"
)

func fc1Save(srcStationId byte, srcNetworkId byte, data []byte) (*FSReply, error) {
	var (
		reply   *FSReply
		session *Session
	)

	slog.Info(fmt.Sprintf("econet-f1-save: src-stn=%02X, src-net=%02X, data=[% 02X]",
		srcStationId, srcNetworkId, data))

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

	// needs to be at least 15 chars
	if len(data) < 15 {

		// TODO should this be a Reply? and an error or just a reply
		// error
		return nil, fmt.Errorf("econet-f0-save: not enough data")
	}

	/* Receives

	    Byte 6 -   Data Acknowledge Port
	    Byte 7 -   Current Selected Directory (CSD)
	    Byte 8 -   Current Selected Library (CSL)
	    Byte 9-12  32-bit Load Address                        <----- is this little endian?
	    Byte 13-16 32-bit Execute Address                      <----- is this little endian?
	    Byte 17-19 24-bit file size
	    Byte 20-n  File Name in ASCII followed by CR

	The reply from the server to the client is as follows.

	    Byte 4   - Command code
	    Byte 5   - Return Code (0 for success)
	    Byte 6   - Data Port
	    Byte 7-8 - Maximum Block Size
	    Byte 9   - File Leaf Name (this is not sent by ArduinoFS but is sent by L3FS padded with spaces to 12 bytes)

	Then its the data transfer phase

	*/

	// the data will give us the reply port. this needs to be stored perhaps in session?
	// as it will be checked for by the listener on each RX_TRANSMIT event
	//dataAckPort := data[0]

	const dataPort byte = 0x9a

	reply = NewFSReply(CCComplete, RCOk, []byte{
		dataPort,
		byte(MaxBlockSize % 256), // this needs to be calculated from q constant (little endian i.e.maxBlockSize=0x0500-1280
		byte(MaxBlockSize / 256),
	})

	// pad the filename to 12 chars and add to the reply
	filename := fmt.Sprintf("%-12s", "NOS")
	reply.Data = append(reply.Data, []byte(filename)...)

	slog.Info(fmt.Sprintf("econet-reply: src-stn=%02X, src-net=%02X, reply-code=%s, data=[% 02X]",
		srcStationId, srcNetworkId, string(ReplyCodeMap[reply.ReturnCode]), reply.ToBytes()))

	return reply, nil

}
