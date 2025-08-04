package lib

import (
	"encoding/binary"
)

// IntToLittleEndianBytes32 returns a little endian 4 byte slice representing the specified 32bit integer.
func IntToLittleEndianBytes32(x uint32) []byte {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, x)
	return bs
}

// IntToLittleEndianBytes24 returns a little endian 3 byte slice representing the specified 32bit integer.
func IntToLittleEndianBytes24(x uint32) []byte {
	return IntToLittleEndianBytes32(x)[:3]
}
func LittleEndianBytesToInt(b []byte) uint32 {

	for ex := 4 - len(b); ex > 0; ex-- {
		b = append(b, 0)
	}

	return binary.LittleEndian.Uint32(b)
}
