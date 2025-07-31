package econet

import (
	"reflect"
	"testing"
)

func Test_parseCmd(t *testing.T) {
	// Defining the columns of the table

	var tests = []struct {
		name        string
		input       string
		wantCmd     string
		wantCmdText string
		wantArgs    []string
	}{
		// the table itself
		{"Test 1", "I AM JOHN PASS", "I AM", "I AM JOHN PASS", []string{"JOHN", "PASS"}},
		{"Test 2", "I AM JOHN", "I AM", "I AM JOHN", []string{"JOHN"}},
	}

	// The execution loop
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := parseCommand(tt.input)
			if ans.Cmd != tt.wantCmd || ans.CmdText != tt.wantCmdText || !reflect.DeepEqual(ans.Args, tt.wantArgs) {
				t.Errorf("got %s, want %s, got %s, want %s, got %v, want %v", ans.Cmd, tt.wantCmd, ans.CmdText, tt.wantCmdText, ans.Args, tt.wantArgs)
			}
		})
	}
}
