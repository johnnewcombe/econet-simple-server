package fs

import (
	"errors"
	"fmt"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/johnnewcombe/econet-simple-server/src/lib"
)

var (
	localFilenameRegx = regexp.MustCompile(`^[A-Za-z0-9]+__[a-fA-F0-9]+_[a-fA-F0-9]+_[a-fA-F0-9]+_[a-fA-F0-9]+$`)
)

type FileInfo struct {
	StartAddress   uint32
	ExecuteAddress uint32
	Size           uint32
	Name           string
	WriteByOthers  bool
	ReadByOthers   bool
	Locked         bool
	ReadByOwner    bool
	WriteByOwner   bool
	IsDirectory    bool
	Exists         bool
	LocalPath      string
}

// NewFileInfoFromCliCmdArgs Accepts CLI (Function 0) command args and returns a FileInfo struct with default access permissions
func NewFileInfoFromCliCmdArgs(args []string) (*FileInfo, error) {

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

	// default permissions for new files
	fInfo.WriteByOthers = false
	fInfo.ReadByOthers = false
	fInfo.Locked = true
	fInfo.ReadByOwner = true
	fInfo.WriteByOwner = true

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
		filename   string
		dirName    string
		dirList    []os.DirEntry
		err        error
		accessByte uint64
	)

	dirName, filename = path.Split(localPath)
	parts := strings.Split(filename, "_")

	// we need to check that the local filename fits with the server's filename format
	// i.e. includes all the attributes etc.
	if !localFilenameRegx.MatchString(filename) {
		return nil, errors.New("invalid filename")
	}

	if accessByte, err = strconv.ParseUint(parts[5], 16, 8); err != nil {
		return nil, err
	}
	/* Definition of Access Byte

			   Bits NFS State   Meaning
			   --------------------------------------------
			    7               Undefined
			    6               Undefined
			    5    W    0     Not writable by other users
			              1     Writable by other users
			    4    R    0     Not readable by other users
			              1     Readable by other users
			    3    L    0     Not locked
			              1     Locked
			    2               Undefined
			    1    R    0     Not writable by owner
			              1     Writable by owner
			    0    W    0     Not readable by owner
			              1     Readable by owner

	//TODO descripion above has an error bit one shows R for read but the description states write
		00001011 = 13h

	*/

	fInfo := FileInfo{
		Name:           parts[0],
		StartAddress:   lib.StringToUint32(parts[2]),
		ExecuteAddress: lib.StringToUint32(parts[3]),
		Size:           lib.StringToUint32(parts[4]),
		WriteByOthers:  accessByte&0b00100000 > 0,
		ReadByOthers:   accessByte&0b00010000 > 0,
		Locked:         accessByte&0b00001000 > 0,
		WriteByOwner:   accessByte&0b00000010 > 0,
		ReadByOwner:    accessByte&0b00000001 > 0,
		IsDirectory:    false,
		Exists:         false, //EconetFileExists(localPath),
		LocalPath:      localPath,
	}

	// get a list for files in the directory
	dirList, err = lib.GetDirectoryList(dirName)
	if err != nil {
		return nil, err
	}

	// loop through dir entries looking for the file
	for _, entry := range dirList {

		// if the file exists, then set the flags
		if strings.HasPrefix(entry.Name(), filename+"__") {

			fInfo.Exists = true

			// file exists so update the local path in case the attributes are different
			fInfo.LocalPath = entry.Name()

		}
		if entry.Name() == filename && entry.IsDir() {
			fInfo.IsDirectory = true
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
