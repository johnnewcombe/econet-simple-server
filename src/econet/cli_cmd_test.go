package econet

import (
	"reflect"
	"testing"
)

func TestCliCmdToBytes(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectBts []byte
	}{
		{
			name:      "EmptyInput",
			input:     "",
			expectBts: []byte(""),
		},
		{
			name:      "UnknownCommand",
			input:     "FOO BAR",
			expectBts: []byte(""), // NewCliCmd returns zero-value CliCmd (CmdText is empty)
		},
		{
			name:      "SimpleCommand",
			input:     "I AM JOHN PASS\rabcdef",
			expectBts: []byte("I AM JOHN PASS"), // tidyText uppercases, trims CR
		},
		{
			name:      "ExtraSpacesAndMixedCase",
			input:     "  sAvE   myfile   3000+500  \n",
			expectBts: []byte("SAVE MYFILE 3000+500"),
		},
		{
			name:      "UTF8AndSpacing",
			input:     "dir   Café/π\r\n",
			expectBts: []byte("DIR CAFÉ/Π"), // tidyText uppercases and collapses spaces, keeps UTF-8
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewCliCmd(tc.input)
			got := cmd.ToBytes()
			if !reflect.DeepEqual(got, tc.expectBts) {
				t.Fatalf("ToBytes mismatch. expect: %q got: %q", tc.expectBts, got)
			}
		})
	}
}

func TestNewCliCmd_Table(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantCmd     string
		wantCmdText string
		wantArgs    []string
	}{
		{
			name:        "IAMTwoArgs",
			input:       "I am  John   Pass",
			wantCmd:     "I AM",
			wantCmdText: "I AM JOHN PASS",
			wantArgs:    []string{"JOHN", "PASS"},
		},
		{
			name:        "SAVEWithStartPlusLen",
			input:       "  SAVE   MyFile   3000+500 ",
			wantCmd:     "SAVE",
			wantCmdText: "SAVE MYFILE 3000+500",
			wantArgs:    []string{"MYFILE", "3000+500"},
		},
		{
			name:        "LOADWithPath",
			input:       "load $.lib.basic",
			wantCmd:     "LOAD",
			wantCmdText: "LOAD $.LIB.BASIC",
			wantArgs:    []string{"$.LIB.BASIC"},
		},
		{
			name:        "DIRNoArgs",
			input:       "DIR",
			wantCmd:     "DIR",
			wantCmdText: "DIR",
			// Current implementation splits empty argText -> []string{""}
			wantArgs: []string{""},
		},
		{
			name:        "LibWithPath",
			input:       "  lib   $.foo ",
			wantCmd:     "LIB",
			wantCmdText: "LIB $.FOO",
			wantArgs:    []string{"$.FOO"},
		},
		{
			name:        "UnknownCommandGivesZeroValue",
			input:       "ECHO HELLO",
			wantCmd:     "",
			wantCmdText: "",
			wantArgs:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := NewCliCmd(tt.input)
			if ans.Cmd != tt.wantCmd || ans.CmdText != tt.wantCmdText || !reflect.DeepEqual(ans.Args, tt.wantArgs) {
				t.Errorf("got Cmd=%q CmdText=%q Args=%v, want Cmd=%q CmdText=%q Args=%v", ans.Cmd, ans.CmdText, ans.Args, tt.wantCmd, tt.wantCmdText, tt.wantArgs)
			}
		})
	}
}
