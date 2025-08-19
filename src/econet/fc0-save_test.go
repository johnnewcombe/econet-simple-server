package econet

import (
	"testing"

	"github.com/johnnewcombe/econet-simple-server/src/fs"
	"github.com/johnnewcombe/econet-simple-server/src/lib"
)

func le(s string) uint32 {
	// Helper to mirror how createFileDescriptor() parses numeric strings
	return lib.StringToUint32(s)
}

func TestCreateFileDescriptor(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want fs.FileDescriptor
	}{
		{
			name: "StartPlusLength",
			// e.g. *SAVE MYDATA 3000+500
			args: []string{"MYDATA", "3000+500"},
			want: fs.FileDescriptor{
				Name:           "MYDATA",
				StartAddress:   le("3000"),
				ExecuteAddress: le("3000"),
				Size:           le("500"),
			},
		},
		{
			name: "StartAndSize",
			// e.g. *SAVE MYDATA 3000 3500
			args: []string{"MYDATA", "3000", "3500"},
			want: fs.FileDescriptor{
				Name:           "MYDATA",
				StartAddress:   le("3000"),
				ExecuteAddress: le("3000"),
				Size:           le("500"),
			},
		},
		{
			name: "StartPlusLengthAndExec",
			// e.g. *SAVE BASIC C000+1000 C2B2
			args: []string{"BASIC", "C000+1000", "C2B2"},
			want: fs.FileDescriptor{
				Name:           "BASIC",
				StartAddress:   le("C000"),
				ExecuteAddress: le("C2B2"),
				Size:           le("1000"),
			},
		},
		{
			name: "StartSizeExecAndLoad",
			// e.g. *SAVE PROG 3000 3500 5050 5000
			args: []string{"PROG", "3000", "3500", "5050", "5000"},
			want: fs.FileDescriptor{
				Name:           "PROG",
				StartAddress:   le("5000"), // load address overrides start
				ExecuteAddress: le("5050"),
				Size:           le("500"),
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cmd := CliCmd{Args: tc.args}
			got, err := createFileDescriptor(cmd)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.Name != tc.want.Name {
				t.Errorf("Name: got %q, want %q", got.Name, tc.want.Name)
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
