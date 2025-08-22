package fs

import (
	"reflect"
	"testing"

	"github.com/johnnewcombe/econet-simple-server/src/lib"
)

func le(s string) uint32 {
	// Helper to mirror how createFileDescriptor() parses numeric strings
	return lib.StringToUint32(s)
}

func TestCreateFileTransfer(t *testing.T) {

	tests := []struct {
		name string
		args []byte
		want FileTransfer
	}{
		{
			name: "Tidy Data",
			args: []byte{0x77, 0x02, 0x03, 0x00, 0xE0, 0x00, 0x00, 0x00, 0xE0, 0x00, 0x00, 0xFF, 0x0F, 0x00, 0x4E, 0x4F, 0x53, 0x0D},
			want: FileTransfer{
				Filename:       "NOS",
				StartAddress:   0x00E000,
				ExecuteAddress: 0x00E000,
				Size:           0xFFF,
				FileData:       []byte{},
			},
		},
		{name: "Post Data",
			args: []byte{0x77, 0x02, 0x03, 0x00, 0xE0, 0x00, 0x00, 0x00, 0xE0, 0x00, 0x00, 0xFF, 0x0F, 0x00, 0x4E, 0x4F, 0x53, 0x0D, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
			want: FileTransfer{
				Filename:       "NOS",
				StartAddress:   0x00E000,
				ExecuteAddress: 0x00E000,
				Size:           0xFFF,
				FileData:       []byte{},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := NewFileTransfer(tc.args[3:])
			//filename := strings.Split(string(tc.args[14:]), "\r")[0]
			//got := FileTransfer{
			//	StartAddress:   lib.LittleEndianBytesToInt(tc.args[3:7]),
			//	ExecuteAddress: lib.LittleEndianBytesToInt(tc.args[7:11]),
			//	Size:           lib.LittleEndianBytesToInt(tc.args[11:14]),
			//	Filename:       filename,
			//	FileData:       []byte{},
			//}
			if got.Filename != tc.want.Filename {
				t.Errorf("Filename: got %q, want %q", got.Filename, tc.want.Filename)
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
			if !reflect.DeepEqual(got.FileData, tc.want.FileData) {
				t.Errorf("FileData: got % X, want % X", got.FileData, tc.want.FileData)
			}
		})
	}
}
