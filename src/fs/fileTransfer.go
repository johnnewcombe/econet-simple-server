package fs

import (
	"strings"

	"github.com/johnnewcombe/econet-simple-server/src/lib"
)

type FileTransfer struct {
	Filename         string
	StartAddress     uint32
	ExecuteAddress   uint32
	Size             uint32
	BytesTransferred int
	CurrentDirectory byte
	CurrentLibrary   byte
	FileData         []byte
}

func NewFileTransfer(dataBlock []byte) *FileTransfer {

	if len(dataBlock) < 11 {
		return nil
	}

	filename := strings.Split(string(dataBlock[11:]), "\r")[0]
	result := FileTransfer{
		StartAddress:   lib.LittleEndianBytesToInt(dataBlock[0:4]),
		ExecuteAddress: lib.LittleEndianBytesToInt(dataBlock[4:8]),
		Size:           lib.LittleEndianBytesToInt(dataBlock[8:11]),
		Filename:       filename,
		FileData:       []byte{},
	}
	return &result
}
