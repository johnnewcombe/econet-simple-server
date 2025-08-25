package econet

import (
	"fmt"
)

type NetHeader struct {
	DstStn byte
	DstNet byte
	SrcStn byte
	SrcNet byte
}

type DataFrame struct {
	NetHeader
	//ReplyPort    byte
	//FunctionCode byte
	Data []byte
}

type ScoutFrame struct {
	NetHeader
	ControlByte byte
	Port        byte
	Data        []byte
}

func (s *ScoutFrame) ToBytes() []byte {

	result := []byte{
		s.DstStn,
		s.DstNet,
		s.SrcStn,
		s.SrcNet,
		s.ControlByte,
		s.Port,
	}

	if len(s.Data) > 0 {
		result = append(result, s.Data...)
	}
	return result
}

func (s *ScoutFrame) String() string {

	msg := fmt.Sprintf("scout-dst=%02X/%02X, scout-src=%02X/%02X, scout-ctrl-byte=%02X, scout-port=%02X, scout-port-desc=%s",
		s.DstStn, s.DstNet, s.SrcStn, s.SrcNet, s.ControlByte, s.Port, PortMap[s.Port])

	if len(s.Data) > 0 {
		msg += fmt.Sprintf(", Data=[% 02X]", s.Data)
	}
	return msg
}

func (d *DataFrame) ToBytes() []byte {

	result := []byte{
		d.DstStn,
		d.DstNet,
		d.SrcStn,
		d.SrcNet,
		//d.ReplyPort,
		//d.FunctionCode,
	}
	if len(d.Data) > 0 {
		result = append(result, d.Data...)
	}
	return result
}

func (d *DataFrame) String() string {

	msg := fmt.Sprintf("Data-dst=%02X/%02X, Data-src=%02X/%02X",
		d.DstStn, d.DstNet, d.SrcStn, d.SrcNet)

	if len(d.Data) > 0 {
		msg += fmt.Sprintf(", Data=[% 02X]", d.Data)
	}
	return msg
}

type FSReply struct {
	ReplyPort byte
	Data      []byte
}

func NewFsReplyData(replyPort byte) *FSReply {
	reply := FSReply{}
	reply.ReplyPort = replyPort
	reply.Data = []byte{0}
	return &reply
}

func NewFSReply(replyPort byte, commandCode CommandCode, returnCode ReturnCode, data []byte) *FSReply {

	reply := FSReply{}
	reply.ReplyPort = replyPort
	reply.Data = []byte{
		byte(commandCode),
		byte(returnCode),
	}

	if data != nil {
		reply.Data = append(reply.Data, data...)
	}

	return &reply

}

func (f *FSReply) ToBytes() []byte {

	result := []byte{f.ReplyPort}
	result = append(result, f.Data...)
	return result
}

//func NewFsReplyError(commandCode CommandCode, returnCode ReturnCode) *FSReply {
//
//	Data := ReplyCodeMap[returnCode]
//
//	return &FSReply{
//		commandCode: commandCode,
//		returnCode:  returnCode,
//		Data:        Data,
//	}
//}
