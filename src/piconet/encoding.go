package piconet

import (
	"encoding/base64"
)

func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}
func Base64Decode(data string) ([]byte, error) {

	var (
		result []byte
		err    error
	)

	if result, err = base64.StdEncoding.DecodeString(data); err != nil {
		return []byte{}, err
	}
	return result, nil
}
