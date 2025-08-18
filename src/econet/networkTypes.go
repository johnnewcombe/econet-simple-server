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

	return fmt.Sprintf("scout-dst-stn=%02X, scout-dst-net=%02X, scout-src-stn=%02X, scout-scr-net=%02X, scout-ctrl-byte=%02X, scout-port=%02X, scout-port-desc=%s",
		s.DstStn, s.DstNet, s.SrcStn, s.SrcNet, s.ControlByte, s.Port, PortMap[s.Port])
}

func (s *DataFrame) ToBytes() []byte {

	return []byte{
		s.DstStn,
		s.DstNet,
		s.SrcStn,
		s.SrcNet,
		s.ReplyPort,
		s.FunctionCode,
	}
}

func (d *DataFrame) String() string {
	return fmt.Sprintf("data-dst-stn=%02X, data-dst-net=%02X, data-src-stn=%02X, data-scr-net=%02X, "+
		"reply-port=%02X, function-code=%02x",
		d.DstStn, d.DstNet, d.SrcStn, d.SrcNet, d.ReplyPort, d.FunctionCode)
}

type CliCmd struct {
	Cmd     string
	Args    []string
	CmdText string
}

type FSReply struct {
	CommandCode CommandCode
	ReturnCode  ReturnCode
	Data        []byte
}

func (f *FSReply) ToBytes() []byte {
	result := []byte{
		byte(f.CommandCode),
		byte(f.ReturnCode),
	}
	return append(result, f.Data...)
}

func NewFSReply(commandCode CommandCode, returnCode ReturnCode, data []byte) *FSReply {

	if data == nil {
		data = []byte{}
	}

	return &FSReply{
		CommandCode: commandCode,
		ReturnCode:  returnCode,
		Data:        data,
	}

}

func (c *CliCmd) ToBytes() []byte {
	return []byte(c.CmdText)
}

/*
func NewFsReplyWithError(commandCode CommandCode, returnCode ReturnCode) *FSReply {

	data := ReplyCodeMap[returnCode]

	return &FSReply{
		CommandCode: commandCode,
		ReturnCode:  returnCode,
		Data:        data,
	}
}
*/
