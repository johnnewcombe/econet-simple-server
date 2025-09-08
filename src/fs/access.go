package fs

import "strings"

func IsOwner(econetFilename string, username string) bool {

	if len(username) == 0 || len(econetFilename) < 3 {
		return false
	}

	if strings.HasPrefix(econetFilename, "$."+username) {
		return true
	}
	return false
}
