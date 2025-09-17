package fs

import (
	"os"
	"path"
	"strings"

	"github.com/johnnewcombe/econet-simple-server/src/lib"
)

// EconetFileExists Accepts
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
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), filename+"__") {
			return true
		}
	}

	return false
}
