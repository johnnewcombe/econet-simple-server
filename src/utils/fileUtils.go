package utils

import (
	"errors"
	"os"
)

// Returns true if folder exists
func Exists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func CreateDirectoryIfNotExists(path string) error {

	if !Exists(path) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	}
	return nil
}

func WriteString(path string, content string) error {

	return WriteBytes(path, []byte(content))
}

func WriteBytes(path string, content []byte) error {

	var file *os.File

	// if it doesn't exist then create it
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does not exist
		if file, err = os.Create(path); err != nil {
			return err
		}
		if _, err = file.Write(content); err != nil {
			return err
		}
	}
	return nil
}

func ReadString(path string) (string, error) {

	b, err := ReadBytes(path)
	return string(b), err
}

func ReadBytes(path string) ([]byte, error) {

	b, err := os.ReadFile(path)
	if err != nil {
		return []byte{}, err
	}

	return b, nil
}

func EconetToLocalPath(rootFolder string, econetPath string) string {

	// TODO convert the path from ':$.MYPATH.' etc...
	return ""
}
