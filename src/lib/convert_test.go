package lib

import (
	"reflect"
	"testing"
)

func Test_IntToLittleEndianBytes32(t *testing.T) {
	// Defining the columns of the table

	var tests = []struct {
		name  string
		input uint32
		want  []byte
	}{
		// the table itself
		{"Test 1", 49152, []byte{0x0, 0xC0, 0x0, 0x0}},         // 00 C0 00 00,
		{"Test 2", 0xc000, []byte{0x0, 0xC0, 0x0, 0x0}},        // 00 C0 00 00,
		{"Test 3", 0x10, []byte{0x10, 0x0, 0x0, 0x00}},         // 00 C0 00 00,
		{"Test 4", 0xffffe003, []byte{0x03, 0xe0, 0xff, 0xff}}, // 00 C0 00 00,
	}

	// The execution loop
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := IntToLittleEndianBytes32(tt.input)

			if !reflect.DeepEqual(ans, tt.want) {
				t.Errorf("got % 02X, want % 02X", ans, tt.want)
			}
		})
	}
}

func Test_IntToLittleEndianBytes24(t *testing.T) {
	// Defining the columns of the table

	var tests = []struct {
		name  string
		input uint32
		want  []byte
	}{
		// the table itself
		{"Test 1", 49152, []byte{0x0, 0xC0, 0x0}},  // 00 C0 00 00,
		{"Test 2", 0xc000, []byte{0x0, 0xC0, 0x0}}, // 00 C0 00 00,
		{"Test 3", 0x10, []byte{0x10, 0x0, 0x0}},   // 00 C0 00 00,
	}

	// The execution loop
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := IntToLittleEndianBytes24(tt.input)

			if !reflect.DeepEqual(ans, tt.want) {
				t.Errorf("got % 02X, want % 02X", ans, tt.want)
			}
		})
	}
}

func Test_LittleEndianBytesToInt(t *testing.T) {

	var tests = []struct {
		name  string
		input []byte
		want  uint32
	}{
		// the table itself
		{"Test 1", []byte{0x0, 0xC0, 0x0, 0x0}, 49152},     // 00 C0 00 00,
		{"Test 2", []byte{0x0, 0xC0, 0x0, 0x0}, 0xc000},    // 00 C0 00 00,
		{"Test 3", []byte{0x10, 0x0, 0x0, 0x0}, 0x10},      // 00 C0 00 00,
		{"Test 4", []byte{0x10, 0x0, 0x0}, 0x10},           // 00 C0 00 00,
		{"Test 5", []byte{0x10, 0x0}, 0x10},                // 00 C0 00 00,
		{"Test 6", []byte{0x10, 0x0, 0x0, 0x0, 0x0}, 0x10}, // 00 C0 00 00,
	}

	// The execution loop
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := LittleEndianBytesToInt(tt.input)

			if !reflect.DeepEqual(ans, tt.want) {
				t.Errorf("got % 02X, want % 02X", ans, tt.want)
			}
		})
	}
}
