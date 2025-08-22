package fs

import (
	"fmt"
	"strings"

	"github.com/johnnewcombe/econet-simple-server/src/lib"
)

type FileDescriptor struct {
	StartAddress   uint32
	ExecuteAddress uint32
	Size           uint32
	Name           string
}

func CreateFileDescriptor(args []string) (*FileDescriptor, error) {

	argCount := len(args)
	if argCount < 2 {
		return nil, fmt.Errorf("econet-f0-save: invalid number of arguments")
	}

	fd := FileDescriptor{Name: args[0]}

	var (
		start uint32
		size  uint32
		exec  uint32
		load  uint32
	)

	if strings.Contains(args[1], "+") {

		parts := strings.SplitN(args[1], "+", 2)
		start = lib.StringToUint32(parts[0])
		size = lib.StringToUint32(parts[1])

		if argCount > 2 {
			exec = lib.StringToUint32(args[2])
		} else {
			exec = start
		}

		if argCount > 3 {
			load = lib.StringToUint32(args[3])
		} else {
			load = start
		}

	} else {

		if argCount < 3 {
			return nil, fmt.Errorf("econet-f0-save: invalid number cmd arguments")
		}

		start = lib.StringToUint32(args[1])
		end := lib.StringToUint32(args[2])
		size = end - start

		if argCount > 3 {
			exec = lib.StringToUint32(args[3])
		} else {
			exec = start
		}

		if argCount > 4 {
			load = lib.StringToUint32(args[4])
		} else {
			load = start
		}
	}

	// Load address updates the start address (preserve exec as per original logic)
	fd.StartAddress = load
	fd.Size = size
	fd.ExecuteAddress = exec

	return &fd, nil
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
