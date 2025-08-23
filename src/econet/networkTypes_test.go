package econet

import (
	"reflect"
	"testing"
)

func TestNewFsReply(t *testing.T) {
	// Test cases
	testCases := []struct {
		name        string
		replyPort   byte
		commandCode CommandCode
		returnCode  ReturnCode
		expectedMsg string
	}{
		{
			name:        "Bad Username Error",
			replyPort:   0x9a,
			commandCode: CCIam,
			returnCode:  RCBadUserName,
			expectedMsg: "BAD USERNAME\r",
		},
		{
			name:        "Not Logged In Error",
			replyPort:   0x9a,
			commandCode: CCLoad,
			returnCode:  RCNotLoggedIn,
			expectedMsg: "NOT LOGGED ON\r",
		},
		{
			name:        "Insufficient Access Error",
			replyPort:   0x9a,
			commandCode: CCSave,
			returnCode:  RCInsufficientAccess,
			expectedMsg: "INSUFFICIENT ACCESS\r",
		},
		{
			name:        "Invalid",
			replyPort:   0x9a,
			commandCode: CCIam,
			returnCode:  0xA8,
			expectedMsg: "",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the function being tested
			result := NewFSReply(tc.replyPort, tc.commandCode, tc.returnCode, ReplyCodeMap[tc.returnCode])

			if result.ReplyPort != byte(tc.replyPort) {
				t.Errorf("Expected ReplyPort %v, got %v", tc.commandCode, result.ReplyPort)
			}

			// Check the command code
			if result.data[0] != byte(tc.commandCode) {
				t.Errorf("Expected CommandCode %v, got %v", tc.commandCode, result.data[0])
			}

			// Check the return code
			if result.data[1] != byte(tc.returnCode) {
				t.Errorf("Expected ReturnCode %v, got %v", tc.returnCode, result.data[1])
			}

			// Check the error message
			actualMsg := string(result.data[2:])
			if actualMsg != tc.expectedMsg {
				t.Errorf("Expected error message '%s', got '%s'", tc.expectedMsg, actualMsg)
			}
		})
	}
}

func TestScoutFrameToBytes(t *testing.T) {
	cases := []struct {
		name   string
		frame  ScoutFrame
		expect []byte
	}{
		{
			name: "NoData",
			frame: ScoutFrame{
				NetHeader:   NetHeader{DstStn: 0x11, DstNet: 0x22, SrcStn: 0x33, SrcNet: 0x44},
				ControlByte: 0x55,
				Port:        0x66,
			},
			expect: []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66},
		},
		{
			name: "WithDataIgnored",
			frame: ScoutFrame{
				NetHeader:   NetHeader{DstStn: 0xAA, DstNet: 0xBB, SrcStn: 0xCC, SrcNet: 0xDD},
				ControlByte: 0xEE,
				Port:        0x01,
				Data:        []byte{0x99, 0x98, 0x97},
			},
			expect: []byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0x01, 0x99, 0x98, 0x97},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.frame.ToBytes()
			if !reflect.DeepEqual(got, tc.expect) {
				t.Fatalf("ToBytes mismatch.\nexpect: % X\n     got: % X", tc.expect, got)
			}
		})
	}
}

func TestDataFrameToBytes(t *testing.T) {
	cases := []struct {
		name   string
		frame  DataFrame
		expect []byte
	}{
		{
			name: "NoData",
			frame: DataFrame{
				NetHeader: NetHeader{DstStn: 0x10, DstNet: 0x20, SrcStn: 0x30, SrcNet: 0x40},
			},
			expect: []byte{0x10, 0x20, 0x30, 0x40},
		},
		{
			name: "WithData",
			frame: DataFrame{
				NetHeader: NetHeader{DstStn: 0xA1, DstNet: 0xB2, SrcStn: 0xC3, SrcNet: 0xD4},
				Data:      []byte{0x01, 0x02, 0x03, 0x04},
			},
			expect: []byte{0xA1, 0xB2, 0xC3, 0xD4, 0x01, 0x02, 0x03, 0x04},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.frame.ToBytes()
			if !reflect.DeepEqual(got, tc.expect) {
				t.Fatalf("ToBytes mismatch.\nexpect: % X\n     got: % X", tc.expect, got)
			}
		})
	}
}

func TestScoutFrameString(t *testing.T) {
	cases := []struct {
		name   string
		frame  ScoutFrame
		expect string
	}{
		{
			name: "WithDataKnownPort",
			frame: ScoutFrame{
				NetHeader:   NetHeader{DstStn: 0x11, DstNet: 0x22, SrcStn: 0x33, SrcNet: 0x44},
				ControlByte: 0x55,
				Port:        0x99, // FileServer Command
				Data:        []byte{0xDE, 0xAD, 0xBE, 0xEF},
			},
			expect: "scout-dst=11/22, scout-src=33/44, scout-ctrl-byte=55, scout-port=99, scout-port-desc=FileServer Command, data=[DE AD BE EF]",
		},
		{
			name: "NoDataKnownPort",
			frame: ScoutFrame{
				NetHeader:   NetHeader{DstStn: 0x01, DstNet: 0x02, SrcStn: 0x03, SrcNet: 0x04},
				ControlByte: 0x05,
				Port:        0x90, // FileServer Reply
				Data:        []byte{},
			},
			expect: "scout-dst=01/02, scout-src=03/04, scout-ctrl-byte=05, scout-port=90, scout-port-desc=FileServer Reply",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.frame.String()
			if got != tc.expect {
				t.Fatalf("String() mismatch.\nexpect: %s\n     got: %s", tc.expect, got)
			}
		})
	}
}

func TestDataFrameString(t *testing.T) {
	cases := []struct {
		name   string
		frame  DataFrame
		expect string
	}{
		{
			name: "WithData",
			frame: DataFrame{
				NetHeader: NetHeader{DstStn: 0x10, DstNet: 0x20, SrcStn: 0x30, SrcNet: 0x40},
				Data:      []byte{0x01, 0x02, 0xA0},
			},
			expect: "data-dst=10/20, data-src=30/40, data=[01 02 A0]",
		},
		{
			name: "NoData",
			frame: DataFrame{
				NetHeader: NetHeader{DstStn: 0xAA, DstNet: 0xBB, SrcStn: 0xCC, SrcNet: 0xDD},
				Data:      nil,
			},
			expect: "data-dst=AA/BB, data-src=CC/DD",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.frame.String()
			if got != tc.expect {
				t.Fatalf("String() mismatch.\nexpect: %s\n     got: %s", tc.expect, got)
			}
		})
	}
}
