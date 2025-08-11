package econet

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/johnnewcombe/econet-simple-server/src/lib"
)

func f0_Save(cmd CliCmd, srcStationId byte, srcNetworkId byte) (*FSReply, error) {

	var (
		reply        *FSReply
		filename     string
		startAddress uint32
		endAddress   uint32
		execAddress  uint32
		loadAddress  uint32
		length       uint32
		data         []byte
	)

	slog.Info(fmt.Sprintf("econet-f0-save: src-stn%02X, src-net:%02X, data=[% 02X]", srcStationId, srcNetworkId, cmd.ToBytes()))

	/* Example of command
	2025/08/02 10:25:53.017511 INFO piconet-event=RX_TRANSMIT scout-dst-stn=FE, scout-dst-net=00, scout-src-stn=BF, scout-scr-net=00, scout-ctrl-byte=80, scout-port=99, scout-port-desc=FileServer Command, data-dst-stn=FE, data-dst-net=00, data-src-stn=BF, data-scr-net=00, reply-port=88, function-code=00, usd=01,csd=02,cslV=03, data-bytes=[53 41 56 45 20 43 30 30 30 20 43 30 31 30 0D E0 45 4C 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20]
	2025/08/02 10:25:53.017551 INFO econet-f0-cli:, data=[53 41 56 45 20 43 30 30 30 20 43 30 31 30 0D EF BF BD 45 4C]
	*/

	// TODO need to handle following syntax. This can be done by further splitting Args[1] by a "+"
	//
	//     *SAVE MYDATA 3000+500
	//     *SAVE MYDATA 3000 3500
	//     *SAVE BASIC C000+1000 C2B2      // adds execution address OF C2B2
	//     *SAVE PROG 3000 3500 5050 5000  // adds execution address and load address

	argCount := len(cmd.Args)

	// sort the first arg out
	if argCount > 0 {

		if strings.Contains(cmd.Args[0], "+") {

			// we have the length specified
			arg := strings.Split(cmd.Args[0], "+")

			// get the start address and length
			startAddress = lib.LittleEndianBytesToInt([]byte(arg[0]))
			length = lib.LittleEndianBytesToInt([]byte(arg[1]))

		} else {
			// just the start address
			startAddress = lib.LittleEndianBytesToInt([]byte(cmd.Args[0]))
		}
	}

	//		if len(argCount)>1

	// these are string values
	//filename := cmd.Args[0]
	//startAddress := cmd.Args[1] // could include a '+' followed by length
	//endAddress := cmd.Args[2]
	//execAddress := cmd.Args[3]
	//loadAddress := cmd.Args[4] // same as start address if not specified

	//only if length has not already been determined from Arg[1] See above.
	//length := endAddress-startAddress

	//startAddress := lib.LittleEndianBytesToInt(data[:4])
	// TODO determine if a load address has been specified in which case the start address is the load address

	print(startAddress)
	print(execAddress)
	print(endAddress)
	print(loadAddress)
	print(filename)

	data = append(data, lib.IntToLittleEndianBytes32(loadAddress)...)
	data = append(data, lib.IntToLittleEndianBytes32(execAddress)...)
	data = append(data, lib.IntToLittleEndianBytes24(length)...)
	data = append(data, []byte(filename)...)
	data = append(data, 0x0d)

	reply = NewFSReply(CCSave, RCOk, data)
	return reply, nil

	//reply := NewFSReply(CCSave, WrongPassword, []byte("NOT IMPLEMENTED\r"))
	//return reply, errors.New("not implemented")
}
