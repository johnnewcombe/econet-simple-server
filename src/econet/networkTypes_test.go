package econet

import (
	"testing"
)

func TestNewFsReply(t *testing.T) {
	// Test cases
	testCases := []struct {
		name        string
		commandCode CommandCode
		returnCode  ReturnCode
		expectedMsg string
	}{
		{
			name:        "Bad Username Error",
			commandCode: CCIam,
			returnCode:  RCBadUserName,
			expectedMsg: "BAD USERNAME\r",
		},
		{
			name:        "Not Logged In Error",
			commandCode: CCLoad,
			returnCode:  RCNotLoggedIn,
			expectedMsg: "NOT LOGGED ON\r",
		},
		{
			name:        "Insufficient Access Error",
			commandCode: CCSave,
			returnCode:  RCInsufficientAccess,
			expectedMsg: "INSUFFICIENT ACCESS\r",
		},
		{
			name:        "Invalid",
			commandCode: CCIam,
			returnCode:  0xA8,
			expectedMsg: "",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the function being tested
			result := NewFSReply(tc.commandCode, tc.returnCode, ReplyCodeMap[tc.returnCode])

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
