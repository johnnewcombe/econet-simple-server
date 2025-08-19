package econet

import (
	"fmt"
)

// NetHeader
type NetHeader struct {
	DstStn byte
	DstNet byte
	SrcStn byte
	SrcNet byte
}

type DataFrame struct {
	NetHeader
	ReplyPort    byte
	FunctionCode byte
	Data         []byte
}

type ScoutFrame struct {
	NetHeader
	ControlByte byte
	Port        byte
	Data        []byte
}

func (s *ScoutFrame) ToBytes() []byte {

	return []byte{
		s.DstStn,
		s.DstNet,
		s.SrcStn,
		s.SrcNet,
		s.ControlByte,
		s.Port,
	}
}

func (s *ScoutFrame) String() string {

	var data []byte
	if len(s.Data) > 0 {
		data = s.Data
	}
	return fmt.Sprintf("scout-dst-stn=%02X, scout-dst-net=%02X, scout-src-stn=%02X, scout-scr-net=%02X, scout-ctrl-byte=%02X, scout-port=%02X, scout-port-desc=%s, data=[% 02X]",
		s.DstStn, s.DstNet, s.SrcStn, s.SrcNet, s.ControlByte, s.Port, PortMap[s.Port], data)
}

func (s *DataFrame) ToBytes() []byte {

	result := []byte{
		s.DstStn,
		s.DstNet,
		s.SrcStn,
		s.SrcNet,
		s.ReplyPort,
		s.FunctionCode,
	}
	if len(s.Data) > 0 {
		result = append(result, s.Data...)
	}
	return result
}

func (d *DataFrame) String() string {

	var data []byte
	if len(d.Data) > 0 {
		data = d.Data
	}

	return fmt.Sprintf("data-dst-stn=%02X, data-dst-net=%02X, data-src-stn=%02X, data-scr-net=%02X, "+
		"reply-port=%02X, function-code=%02x, data=[% 02X]",
		d.DstStn, d.DstNet, d.SrcStn, d.SrcNet, d.ReplyPort, d.FunctionCode, data)

}

type CliCmd struct {
	Cmd     string
	Args    []string
	CmdText string
}

type FSReply struct {
	data []byte
}

//func (f *FSReply) Append(data []byte) {
//	f.data = append(f.data, data...)
//}

func (f *FSReply) ToBytes() []byte {
	return f.data
}
func NewFsReplyData(data []byte) *FSReply {
	reply := FSReply{}
	reply.data = data
	return &reply
}

func NewFSReply(commandCode CommandCode, returnCode ReturnCode, data []byte) *FSReply {

	reply := FSReply{}
	reply.data = []byte{
		byte(commandCode),
		byte(returnCode),
	}

	if data != nil {
		reply.data = append(reply.data, data...)
	}

	return &reply

}

func (c *CliCmd) ToBytes() []byte {
	return []byte(c.CmdText)
}

//func NewFsReplyError(commandCode CommandCode, returnCode ReturnCode) *FSReply {
//
//	data := ReplyCodeMap[returnCode]
//
//	return &FSReply{
//		commandCode: commandCode,
//		returnCode:  returnCode,
//		data:        data,
//	}
//}
