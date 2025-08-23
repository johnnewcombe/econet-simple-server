package lib

import (
	"strings"
)

// HexDump formats a byte slice into lines of 16 hexadecimal byte pairs
// followed by their ASCII representation. Non-printable bytes are shown
// as '.'. Lines are separated by a newline character. If b is empty,
// an empty string is returned.
func HexDump(b []byte) string {
	if len(b) == 0 {
		return ""
	}

	var sb strings.Builder

	for i := 0; i < len(b); i += 16 {
		end := i + 16
		if end > len(b) {
			end = len(b)
		}

		// Hex portion (padded to 16 bytes width)
		for j := 0; j < 16; j++ {
			idx := i + j
			if idx < end {
				sb.WriteString(hexByte(b[idx]))
				sb.WriteByte(' ')
			} else {
				// three spaces to account for two hex chars + space
				sb.WriteString("   ")
			}
		}

		// Space between hex and ASCII
		sb.WriteByte(' ')
		sb.WriteByte(' ')

		// ASCII portion
		for j := i; j < end; j++ {
			c := b[j]
			if c >= 32 && c <= 126 { // printable ASCII range
				sb.WriteByte(c)
			} else {
				sb.WriteByte('.')
			}
		}

		if end < i+16 {
			// No need to pad ASCII; it naturally shortens.
		}

		if end < len(b) {
			sb.WriteByte('\n')
		}
	}

	return sb.String()
}

// hexByte returns the uppercase hexadecimal representation of a single byte as two characters.
func hexByte(b byte) string {
	const hexdigits = "0123456789ABCDEF"
	var buf [2]byte
	buf[0] = hexdigits[b>>4]
	buf[1] = hexdigits[b&0x0F]
	return string(buf[:])
}
