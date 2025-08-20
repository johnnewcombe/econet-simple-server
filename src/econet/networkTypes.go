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
		msg += fmt.Sprintf(", data=[% 02X]", s.Data)
	}
	return msg
}

func (d *DataFrame) ToBytes() []byte {

	result := []byte{
		d.DstStn,
		d.DstNet,
		d.SrcStn,
		d.SrcNet,
		d.ReplyPort,
		d.FunctionCode,
	}
	if len(d.Data) > 0 {
		result = append(result, d.Data...)
	}
	return result
}

func (d *DataFrame) String() string {

	msg := fmt.Sprintf("data-dst=%02X/%02X, data-src=%02X/%02X, reply-port=%02X, function-code=%02x",
		d.DstStn, d.DstNet, d.SrcStn, d.SrcNet, d.ReplyPort, d.FunctionCode)

	if len(d.Data) > 0 {
		msg += fmt.Sprintf(", data=[% 02X]", d.Data)
	}
	return msg
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
