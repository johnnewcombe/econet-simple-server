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
	Usd          byte
	Csd          byte
	Csl          byte
	Data         []byte
}

type ScoutFrame struct {
	NetHeader
	ControlByte byte
	Port        byte
	Data        []byte
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

	return &FSReply{
		CommandCode: commandCode,
		ReturnCode:  returnCode,
		Data:        data,
	}

}

type CliCmd struct {
	Cmd     string
	Args    []string
	CmdText string
}
