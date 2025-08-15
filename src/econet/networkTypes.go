package econet

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
