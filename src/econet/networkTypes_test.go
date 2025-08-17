package econet

import (
	"testing"
)

func TestNewFsReplyWithError(t *testing.T) {
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
			if result.CommandCode != tc.commandCode {
				t.Errorf("Expected CommandCode %v, got %v", tc.commandCode, result.CommandCode)
			}

			// Check the return code
			if result.ReturnCode != tc.returnCode {
				t.Errorf("Expected ReturnCode %v, got %v", tc.returnCode, result.ReturnCode)
			}

			// Check the error message
			actualMsg := string(result.Data)
			if actualMsg != tc.expectedMsg {
				t.Errorf("Expected error message '%s', got '%s'", tc.expectedMsg, actualMsg)
			}
		})
	}
}
