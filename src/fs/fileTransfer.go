package fs

import (
	"fmt"
	"strings"
)

type FileTransfer struct {
	Filename         string
	DiskName         string
	StartAddress     uint32
	ExecuteAddress   uint32
	Size             uint32
	BytesTransferred int
	CurrentDirectory byte
	CurrentLibrary   byte
	FileData         []byte
	DataAckPort      byte
	ReplyPort        byte
	FunctionCode     byte
}

func NewFileTransfer(functionCode byte, replyPort byte, startAddress uint32, execAddress uint32, fileSize uint32, filename string, diskname string) *FileTransfer {

	result := FileTransfer{
		FunctionCode:   functionCode,
		ReplyPort:      replyPort,
		StartAddress:   startAddress,
		ExecuteAddress: execAddress,
		Size:           fileSize,
		Filename:       strings.Split(filename, "\r")[0], // belts and braces
		DiskName:       diskname,
		FileData:       []byte{},
	}
	return &result
}

// GetLeafName returns the leaf name of the filename padded with spaces to 12 characters.
func (f *FileTransfer) GetLeafName() string {
	parts := strings.Split(f.Filename, ".")
	leaf := parts[len(parts)-1]
	return fmt.Sprintf("%-12s", leaf)
}

// TODO Consider refactoring this to not parse the filename from the data block etc.
//func NewFileTransferOld(functionCode byte, replyPort byte, dataBlock []byte) *FileTransfer {
//
//	if len(dataBlock) < 11 {
//		return nil
//	}
//
//	filename := strings.Split(string(dataBlock[11:]), "\r")[0]
//	result := FileTransfer{
//		FunctionCode:   functionCode,
//		ReplyPort:      replyPort,
//		StartAddress:   lib.LittleEndianBytesToInt(dataBlock[0:4]),
//		ExecuteAddress: lib.LittleEndianBytesToInt(dataBlock[4:8]),
//		Size:           lib.LittleEndianBytesToInt(dataBlock[8:11]),
//		Filename:       filename,
//		FileData:       []byte{},
//	}
//	return &result
//}
