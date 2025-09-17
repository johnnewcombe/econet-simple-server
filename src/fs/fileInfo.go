package fs

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/johnnewcombe/econet-simple-server/src/lib"
)

type FileInfo struct {
	StartAddress   uint32
	ExecuteAddress uint32
	Size           uint32
	Name           string
	Locked         bool
	ReadAccess     bool
	WriteAccess    bool
	IsDirectory    bool
	Exists         bool
}

func NewFileInfo(args []string) (*FileInfo, error) {

	argCount := len(args)
	if argCount < 2 {
		return nil, fmt.Errorf("econet-f0-save: invalid number of arguments")
	}

	fInfo := FileInfo{Name: args[0]}

	var (
		start uint32
		size  uint32
		exec  uint32
		load  uint32
	)

	if strings.Contains(args[1], "+") {

		parts := strings.SplitN(args[1], "+", 2)
		start = lib.StringToUint32(parts[0])
		size = lib.StringToUint32(parts[1]) + 1

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
		size = end - start + 1

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
	fInfo.StartAddress = load
	fInfo.Size = size
	fInfo.ExecuteAddress = exec
	fInfo.IsDirectory = false

	// defaults for new files
	// TODO get defaults e.g. If directory then ... etc
	fInfo.Locked = true
	fInfo.ReadAccess = true
	fInfo.WriteAccess = false

	//TDOD parse the access byte
	//fInfo.LocalPath = fmt.Sprintf("%s__%4X_%4X_%3X_%2X",
	//	fInfo.Name,
	//	fInfo.StartAddress,
	//	fInfo.ExecuteAddress,
	//	fInfo.Size,
	//	0x00)

	return &fInfo, nil
}

func (f *FileInfo) ToBytes() []byte {

	var (
		data []byte
	)

	// returns bytes in order required by econet protocol and name is terminated with a CR
	data = append(data, lib.IntToLittleEndianBytes32(f.StartAddress)...)
	data = append(data, lib.IntToLittleEndianBytes32(f.ExecuteAddress)...)
	data = append(data, lib.IntToLittleEndianBytes24(f.Size)...)
	data = append(data, []byte(f.Name)...)
	data = append(data, 0x0d)

	return data

}

func NewFileInfoFromLocalPath(localPath string) (*FileInfo, error) {

	var (
		filename string
		dirName  string
		dirList  []os.DirEntry
		err      error
	)

	dirName, filename = path.Split(localPath)
	parts := strings.Split(filename, "_")

	// note that part [1] should be a blank line due to the double underscore in the filename
	if len(parts) < 5 || len(parts[1]) != 0 {
		return nil, errors.New("invalid filename")
	}

	fInfo := FileInfo{
		Name:           parts[0],
		StartAddress:   lib.StringToUint32(parts[2]),
		ExecuteAddress: lib.StringToUint32(parts[3]),
		Size:           lib.StringToUint32(parts[4]),
		Locked:         false, // TODO this needs to be implemented
		ReadAccess:     false, // TODO this needs to be implemented
		WriteAccess:    false, // TODO this needs to be implemented
		IsDirectory:    false,
		Exists:         false,
	}

	// get a list for files in the directory
	dirList, err = lib.GetDirectoryList(dirName)
	if err != nil {
		return nil, err
	}

	// loop through dir entries looking for the file
	for _, entry := range dirList {

		// if the file exists, then set the flags
		if entry.Name() == filename {
			fInfo.Exists = true
			fInfo.IsDirectory = entry.IsDir()
		}
	}

	return &fInfo, nil
}

// EconetFileExists Accepts a local path and returns true if the file exists with
// the second return value indicating if the file is a directory.
func EconetFileExists(localPath string) bool {

	var (
		dirList []os.DirEntry
		err     error
	)

	// get the directory and filenam from the path
	dirName, filename := path.Split(localPath)

	// get a list for files in the directory
	dirList, err = lib.GetDirectoryList(dirName)
	if err != nil {
		return false
	}

	// loop through dir entries looking for the file
	for _, entry := range dirList {

		// the filename on disk is followed by attributes such as start and execute addresses
		// these are separated from the filename by a double underscore
		if strings.HasPrefix(entry.Name(), filename+"__") {
			return true
		}
	}

	return false
}
