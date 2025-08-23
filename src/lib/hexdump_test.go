package lib

import (
	"strings"
	"testing"
)

// helper to build expected string mirroring the documented format
func buildExpected(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	var sb strings.Builder
	for i := 0; i < len(b); i += 16 {
		end := i + 16
		if end > len(b) {
			end = len(b)
		}
		for j := 0; j < 16; j++ {
			idx := i + j
			if idx < end {
				sb.WriteString(hexByte(b[idx]))
				sb.WriteByte(' ')
			} else {
				sb.WriteString("   ")
			}
		}
		sb.WriteByte(' ')
		sb.WriteByte(' ')
		for j := i; j < end; j++ {
			c := b[j]
			if c >= 32 && c <= 126 {
				sb.WriteByte(c)
			} else {
				sb.WriteByte('.')
			}
		}
		if end < len(b) {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func TestHexDump_Empty(t *testing.T) {
	got := HexDump(nil)
	if got != "" {
		t.Fatalf("expected empty string for nil input, got %q", got)
	}
	got = HexDump([]byte{})
	println(got)
	if got != "" {
		t.Fatalf("expected empty string for empty slice, got %q", got)
	}
}

func TestHexDump_LessThan16(t *testing.T) {
	in := []byte{0x41, 0x42, 0x00, 0x7F, 0x20}
	expect := buildExpected(in)
	got := HexDump(in)
	println(got)
	if got != expect {
		t.Fatalf("hexdump mismatch (len<16)\nexpect: %q\n     got: %q", expect, got)
	}
}

func TestHexDump_Exactly16(t *testing.T) {
	in := make([]byte, 16)
	for i := 0; i < 16; i++ {
		in[i] = byte(i)
	}
	expect := buildExpected(in)
	got := HexDump(in)
	println(got)
	if got != expect {
		t.Fatalf("hexdump mismatch (exactly 16)\nexpect: %q\n     got: %q", expect, got)
	}
}

func TestHexDump_MoreThan16_MultiLine(t *testing.T) {

	in := make([]byte, 20)
	for i := 0; i < 20; i++ {
		in[i] = 0x30 + byte(i) // '0', '1', ... printable ASCII
	}
	expect := buildExpected(in)
	got := HexDump(in)
	println(got)
	if got != expect {
		t.Fatalf("hexdump mismatch (multi-line)\nexpect: %q\n     got: %q", expect, got)
	}

}
