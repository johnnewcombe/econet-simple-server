package lib

import (
	"encoding/binary"
	"strconv"
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
func StringToLittleEndianBytes(s string) []byte {
	var (
		i   uint64
		err error
	)
	//s = Reverse(s)

	if i, err = strconv.ParseUint(s, 16, 32); err != nil {
		i = 0
	}
	return IntToLittleEndianBytes32(uint32(i))
}

func StringToUint32(s string) uint32 {
	var (
		i   uint64
		err error
	)
	if i, err = strconv.ParseUint(s, 16, 32); err != nil {
		i = 0
	}
	return uint32(i)
	//return LittleEndianBytesToInt(StringToLittleEndianBytes(s))
}

// Reverse reverses a string e.g. "C123" becomes "321C"
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
