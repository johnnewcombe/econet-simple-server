package fs

import (
	"github.com/johnnewcombe/econet-simple-server/src/lib"
)

type FileDescriptor struct {
	StartAddress   uint32
	ExecuteAddress uint32
	Size           uint32
	Name           string
}

func (fd *FileDescriptor) ToBytes() []byte {

	var (
		data []byte
	)

	// returns bytes in order required by econet protocol and name is terminated with a CR
	data = append(data, lib.IntToLittleEndianBytes32(fd.StartAddress)...)
	data = append(data, lib.IntToLittleEndianBytes32(fd.ExecuteAddress)...)
	data = append(data, lib.IntToLittleEndianBytes24(fd.Size)...)
	data = append(data, []byte(fd.Name)...)
	data = append(data, 0x0d)

	return data

}
