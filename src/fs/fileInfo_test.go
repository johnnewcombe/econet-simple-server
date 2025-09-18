package fs

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
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
		want FileInfo
	}{
		{
			name: "StartPlusLength",
			// e.g. *SAVE MYDATA 3000+500
			args: []string{"MYDATA", "3000+4ff"},
			want: FileInfo{
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
			want: FileInfo{
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
			want: FileInfo{
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
			want: FileInfo{
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

			got, err := NewFileInfoFromCliCmdArgs(tc.args)
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
		fd   FileInfo
		want []byte
	}{
		{
			name: "all-zero-values-empty-name",
			fd: FileInfo{
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
			fd: FileInfo{
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
			fd: FileInfo{
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
			fd: FileInfo{
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
	fd := FileInfo{
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

func TestNewFileInfoFromLocalPath_Success(t *testing.T) {
	// Build a local path with expected underscore-separated parts (>=5 parts)
	// Format expected by code: name_start_end_... (at least 5 parts)
	dir := t.TempDir()

	// NAME__STAR_EXEC_SIZE_ACCESS_MISC
	base := "MYPROG__C000_C2B2_FF_0b"
	lp := filepath.Join(dir, base)

	wantName := "MYPROG"
	wantStartAddress := uint32(0xC000)
	wantExecAddress := uint32(0xC2B2)
	wantSize := uint32(0xFF)
	wantWriteByOthers := false
	wantReadByOthers := false
	wantLocked := true
	wantReadByOwner := true
	wantWriteByOwner := true

	// e.g. access = 00001011
	if err := os.WriteFile(lp, []byte("x"), 0o644); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	got, err := NewFileInfoFromLocalPath(lp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil {
		t.Fatalf("expected non-nil FileInfo")
	}

	if got.Name != wantName ||
		got.StartAddress != wantStartAddress ||
		got.ExecuteAddress != wantExecAddress ||
		got.Size != wantSize ||
		got.WriteByOthers != wantWriteByOthers ||
		got.ReadByOthers != wantReadByOthers ||
		got.Locked != wantLocked ||
		got.ReadByOwner != wantReadByOwner ||
		got.WriteByOwner != wantWriteByOwner {

		t.Fatalf("Name: got %v, want %v, SAddr: got %v, want %v, EAddr: got %v, want %v, Size: got %v, want %v, WPub: got %v, want %v, "+
			"RPub: got %v, want %v, Lkd: got %v, want %v, WOwn: got %v, want %v, ROwn: got %v, want %v, ",
			got.Name, wantName,
			got.StartAddress, wantStartAddress,
			got.ExecuteAddress, wantExecAddress,
			got.Size, wantSize,
			got.WriteByOthers, wantWriteByOthers,
			got.ReadByOthers, wantReadByOthers,
			got.Locked, wantLocked,
			got.ReadByOwner, wantReadByOwner,
			got.WriteByOwner, wantWriteByOwner)

	}
}

func TestNewFileInfoFromLocalPath_InvalidFormat(t *testing.T) {

	dir := t.TempDir()
	// Too few parts (only 3 when split by '_') should yield an error per code
	lp := filepath.Join(dir, "NAME_ONLY_BAD")

	fi, err := NewFileInfoFromLocalPath(lp)
	got := fi != nil
	if err == nil {
		t.Fatalf("expected error for invalid filename format, got nil (result: %#v)", got)
	}
}

// EconetFileExists checks for entries that start with filename+"__" in the same dir.
func TestEconetFileExists(t *testing.T) {
	// Each test should use its own temp dir; otherwise earlier files could leak into later cases
	tests := []struct {
		name  string
		setup func(dir string) string // returns localPath to query
		want  bool
	}{
		{
			name: "exact name with double underscore exists",
			setup: func(dir string) string {
				// create file that matches pattern: <name>__<attrs>
				fname := "prog__C000_C000_010"
				if err := os.WriteFile(filepath.Join(dir, fname), []byte("x"), 0o644); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				return filepath.Join(dir, "prog")
			},
			want: true,
		},
		{
			name: "prefix that is not followed by double underscore should not match",
			setup: func(dir string) string {
				// file without the double underscore after the prefix
				if err := os.WriteFile(filepath.Join(dir, "progC000.bin"), []byte("x"), 0o644); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				return filepath.Join(dir, "prog")
			},
			want: false,
		},
		{
			name: "directories do not count",
			setup: func(dir string) string {
				if err := os.Mkdir(filepath.Join(dir, "doc__F000"), 0o755); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				return filepath.Join(dir, "doc")
			},
			want: true, // Note: the function does not filter out dirs; it checks names only
		},
		{
			name: "non-existent directory returns false",
			setup: func(dir string) string {
				return filepath.Join(dir, "nope", "file")
			},
			want: false,
		},
		{
			name: "no matches returns false",
			setup: func(dir string) string {
				if err := os.WriteFile(filepath.Join(dir, "alpha__1000_1000_010"), []byte("x"), 0o644); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				return filepath.Join(dir, "beta")
			},
			want: false,
		},
		{
			name: "multiple entries one matches",
			setup: func(dir string) string {
				_ = os.WriteFile(filepath.Join(dir, "note1"), []byte("x"), 0o644)
				_ = os.WriteFile(filepath.Join(dir, "note__AAAA_BBBB_005"), []byte("x"), 0o644)
				_ = os.WriteFile(filepath.Join(dir, "other__C000_C000_001"), []byte("x"), 0o644)
				return filepath.Join(dir, "note")
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			localPath := tt.setup(dir)
			got := EconetFileExists(localPath)
			if got != tt.want {
				t.Fatalf("EconetFileExists(%q) = %v, want %v", localPath, got, tt.want)
			}
		})
	}
}
