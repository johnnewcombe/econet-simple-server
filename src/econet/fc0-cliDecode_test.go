package econet

import (
	"reflect"
	"testing"
)

func Test_tidyText(t *testing.T) {
	// Table-driven test cases for tidyText
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"TrimNullCharacter", "\x00HELLO\x00", "HELLO"},
		{"TrimNewlines", "\nHELLO\n", "HELLO"},
		{"TrimCarriageReturns", "\rHELLO\r", "HELLO"},
		{"TrimSpaces", "  HELLO  ", "HELLO"},
		{"ConvertToUppercase", "hello world", "HELLO WORLD"},
		{"MultipleSpaces", "HELLO    WORLD", "HELLO WORLD"},
		{"EmptyString", "", ""},
		{"SingleSpace", " ", ""},
		{"SingleCharacter", "a", "A"},
		{"MixedFormatting", "\x00 HeLlO  wOrLd \n\r", "HELLO WORLD"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tidyText(tt.input)
			if got != tt.want {
				t.Errorf("tidyText(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

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
		{"Test 3", "I  aM  JOhN pass", "I AM", "I AM JOHN PASS", []string{"JOHN", "PASS"}},
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
