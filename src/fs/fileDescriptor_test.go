package fs

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"
)

// le32 returns a 4-byte little-endian representation of v.
func le32(v uint32) []byte {
	return []byte{
		byte(v),
		byte(v >> 8),
		byte(v >> 16),
		byte(v >> 24),
	}
}

// le24 returns a 3-byte little-endian representation of the low 24 bits of v.
func le24(v uint32) []byte {
	return []byte{
		byte(v),
		byte(v >> 8),
		byte(v >> 16),
	}
}

func TestCreateFileDescriptor(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want FileDescriptor
	}{
		{
			name: "StartPlusLength",
			// e.g. *SAVE MYDATA 3000+500
			args: []string{"MYDATA", "3000+4ff"},
			want: FileDescriptor{
				Name:           "MYDATA",
				StartAddress:   le("3000"),
				ExecuteAddress: le("3000"),
				Size:           le("500"),
			},
		},
		{
			name: "StartAndSize",
			// e.g. *SAVE MYDATA 3000 3500
			args: []string{"MYDATA", "3000", "34ff"},
			want: FileDescriptor{
				Name:           "MYDATA",
				StartAddress:   le("3000"),
				ExecuteAddress: le("3000"),
				Size:           le("500"),
			},
		},
		{
			name: "StartPlusLengthAndExec",
			// e.g. *SAVE BASIC C000+1000 C2B2
			args: []string{"BASIC", "C000+FFF", "C2B2"},
			want: FileDescriptor{
				Name:           "BASIC",
				StartAddress:   le("C000"),
				ExecuteAddress: le("C2B2"),
				Size:           le("1000"),
			},
		},
		{
			name: "StartSizeExecAndLoad",
			// e.g. *SAVE PROG 3000 3500 5050 5000
			args: []string{"PROG", "3000", "34ff", "5050", "5000"},
			want: FileDescriptor{
				Name:           "PROG",
				StartAddress:   le("5000"), // load address overrides start
				ExecuteAddress: le("5050"),
				Size:           le("500"),
			},
		},
	}

	for _, tc := range tests {
		//tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewFileDescriptor(tc.args)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.Name != tc.want.Name {
				t.Errorf("Filename: got %q, want %q", got.Name, tc.want.Name)
			}
			if got.StartAddress != tc.want.StartAddress {
				t.Errorf("StartAddress: got %08x, want %08x", got.StartAddress, tc.want.StartAddress)
			}
			if got.ExecuteAddress != tc.want.ExecuteAddress {
				t.Errorf("ExecuteAddress: got %08x, want %08x", got.ExecuteAddress, tc.want.ExecuteAddress)
			}
			if got.Size != tc.want.Size {
				t.Errorf("Size: got %06x, want %06x", got.Size, tc.want.Size)
			}
		})
	}
}

func TestFileDescriptor_ToBytes(t *testing.T) {
	tests := []struct {
		name string
		fd   FileDescriptor
		want []byte
	}{
		{
			name: "all-zero-values-empty-name",
			fd: FileDescriptor{
				StartAddress:   0x00000000,
				ExecuteAddress: 0x00000000,
				Size:           0x000000,
				Name:           "",
			},
			want: append(append(append(append(
				le32(0x00000000),
				le32(0x00000000)...),
				le24(0x000000)...),
				0x0D),
			),
		},
		{
			name: "simple-values-with-name",
			fd: FileDescriptor{
				StartAddress:   0x00003000,
				ExecuteAddress: 0x00003500,
				Size:           0x000500,
				Name:           "MYDATA",
			},
			want: func() []byte {
				b := append([]byte{}, le32(0x00003000)...)
				b = append(b, le32(0x00003500)...)
				b = append(b, le24(0x000500)...)
				b = append(b, []byte("MYDATA")...)
				b = append(b, 0x0D)
				return b
			}(),
		},
		{
			name: "mixed-values-and-24bit-size",
			fd: FileDescriptor{
				StartAddress:   0x12345678,
				ExecuteAddress: 0x9ABCDEF0,
				Size:           0x00ABCDEF, // only low 24 bits used
				Name:           "DATA",
			},
			want: func() []byte {
				b := append([]byte{}, le32(0x12345678)...)
				b = append(b, le32(0x9ABCDEF0)...)
				b = append(b, le24(0x00ABCDEF)...)
				b = append(b, []byte("DATA")...)
				b = append(b, 0x0D)
				return b
			}(),
		},
		{
			name: "max-24bit-size",
			fd: FileDescriptor{
				StartAddress:   0xFFFFFFFF,
				ExecuteAddress: 0x00000001,
				Size:           0xFFFFFF,
				Name:           "X",
			},
			want: func() []byte {
				b := append([]byte{}, le32(0xFFFFFFFF)...)
				b = append(b, le32(0x00000001)...)
				b = append(b, le24(0xFFFFFF)...)
				b = append(b, []byte("X")...)
				b = append(b, 0x0D)
				return b
			}(),
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := tc.fd.ToBytes()
			if !bytes.Equal(got, tc.want) {
				t.Fatalf("ToBytes() mismatch.\n got: %s\nwant: %s\n got(raw): %v\nwant(raw): %v",
					hex.EncodeToString(got),
					hex.EncodeToString(tc.want),
					got,
					tc.want,
				)
			}

			// Ensure CR terminator is present at the end
			if len(got) == 0 || got[len(got)-1] != 0x0D {
				t.Fatalf("ToBytes() missing CR terminator at end: last=%#x, len=%d", got[len(got)-1], len(got))
			}

			// Quick structural sanity checks
			const headerLen = 4 + 4 + 3 // start + exec + size
			if len(got) < headerLen+1 { // at least header + CR
				t.Fatalf("ToBytes() too short, got len=%d", len(got))
			}
			nameBytes := got[headerLen : len(got)-1]
			if string(nameBytes) != tc.fd.Name {
				t.Fatalf("ToBytes() name mismatch: got=%q want=%q (bytes=%v)", string(nameBytes), tc.fd.Name, []byte(tc.fd.Name))
			}
		})
	}
}

// Optional: a focused test to help with debugging if needed.
func TestFileDescriptor_ToBytes_DebugExample(t *testing.T) {
	fd := FileDescriptor{
		StartAddress:   0x0000E000,
		ExecuteAddress: 0x0000EFFF,
		Size:           0x000FFF,
		Name:           "MYDATA",
	}
	got := fd.ToBytes()
	want := append(append(append(append(
		le32(0x0000E000),
		le32(0x0000EFFF)...),
		le24(0x000FFF)...),
		[]byte("MYDATA")...),
		0x0D)

	if !bytes.Equal(got, want) {
		t.Fatalf("debug example mismatch:\n got=%s\nwant=%s", hex.EncodeToString(got), hex.EncodeToString(want))
	}
	fmt.Println("debug example:", hex.EncodeToString(got))
}
