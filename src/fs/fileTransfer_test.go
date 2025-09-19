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

func TestNewFileTransfer(t *testing.T) {
	tests := []struct {
		name         string
		functionCode byte
		replyPort    byte
		start        uint32
		exec         uint32
		size         uint32
		access       byte
		filenameIn   string
		diskName     string
		want         FileTransfer
	}{
		{
			name:         "Tidy_Data",
			functionCode: 0x01,
			replyPort:    0x10,
			start:        0x00003000,
			exec:         0x00004000,
			size:         0x000500,
			access:       0x13,
			filenameIn:   "NOS\r",
			diskName:     "DISK0",
			want: FileTransfer{
				FunctionCode:   0x01,
				ReplyPort:      0x10,
				StartAddress:   0x00003000,
				ExecuteAddress: 0x00004000,
				Size:           0x000500,
				Filename:       "NOS",
				DiskName:       "DISK0",
				FileData:       []byte{},
			},
		},
		{
			name:         "Post_Data",
			functionCode: 0x02,
			replyPort:    0x11,
			start:        0x0000C000,
			exec:         0x0000C2B2,
			size:         0x001000,
			access:       0x13,
			filenameIn:   "BASIC\rEXTRA",
			diskName:     "DISK0",
			want: FileTransfer{
				FunctionCode:   0x02,
				ReplyPort:      0x11,
				StartAddress:   0x0000C000,
				ExecuteAddress: 0x0000C2B2,
				Size:           0x001000,
				Filename:       "BASIC",
				DiskName:       "DISK0",
				FileData:       []byte{},
			},
		},
		{
			name:         "No_CR_In_Name",
			functionCode: 0x03,
			replyPort:    0x12,
			start:        0x12345678,
			exec:         0x9ABCDEF0,
			size:         0x000123,
			access:       0x13,
			filenameIn:   "DATA",
			diskName:     "DISK0",
			want: FileTransfer{
				FunctionCode:   0x03,
				ReplyPort:      0x12,
				StartAddress:   0x12345678,
				ExecuteAddress: 0x9ABCDEF0,
				Size:           0x000123,
				Filename:       "DATA",
				DiskName:       "DISK0",
				FileData:       []byte{},
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := NewFileTransfer(tc.functionCode, tc.replyPort, tc.start, tc.exec, tc.size, 0x13, tc.filenameIn, tc.diskName)
			if got == nil {
				t.Fatalf("NewFileTransfer returned nil")
			}
			if got.FunctionCode != tc.want.FunctionCode {
				t.Errorf("FunctionCode: got %02x, want %02x", got.FunctionCode, tc.want.FunctionCode)
			}
			if got.ReplyPort != tc.want.ReplyPort {
				t.Errorf("ReplyPort: got %02x, want %02x", got.ReplyPort, tc.want.ReplyPort)
			}
			if got.Filename != tc.want.Filename {
				t.Errorf("Filename: got %q, want %q", got.Filename, tc.want.Filename)
			}
			if got.DiskName != tc.want.DiskName {
				t.Errorf("Diskname: got %q, want %q", got.DiskName, tc.want.DiskName)
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
			if got.BytesTransferred != 0 {
				t.Errorf("BytesTransferred: got %d, want %d", got.BytesTransferred, 0)
			}
		})
	}
}

//func TestGetLeafName(t *testing.T) {
//	tests := []struct {
//		name     string
//		filename string
//		want     string
//	}{
//		{name: "no_dot", filename: "BASIC", want: "BASIC       "},
//		{name: "single_dot", filename: "LIB.FILE", want: "FILE        "},
//		{name: "multiple_dots", filename: "LIB.SUB.FILE", want: "FILE        "},
//		{name: "leading_dot", filename: ".HIDDEN", want: "HIDDEN      "},
//		{name: "trailing_dot", filename: "NAME.", want: "            "},
//		{name: "empty", filename: "", want: "            "},
//	}
//
//	for _, tc := range tests {
//		tc := tc
//		t.Run(tc.name, func(t *testing.T) {
//			ft := &FileTransfer{Filename: tc.filename}
//			got := ft.GetLeafName()
//			if got != tc.want {
//				t.Errorf("GetLeafName(%q) = %q, want %q", tc.filename, got, tc.want)
//			}
//		})
//	}
//}
//
