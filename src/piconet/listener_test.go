package piconet

import (
	"testing"
)

func Test_tidyText(t *testing.T) {
	// Define test cases
	var tests = []struct {
		name  string
		input string
		want  string
	}{
		{"Empty string", "", ""},
		{"String with null characters", "\x00test\x00", "test"},
		{"String with multiple spaces", "test  with    spaces", "test with spaces"},
		{"String with leading and trailing spaces", "  test with spaces  ", "test with spaces"},
		{"Normal string", "test string", "test string"},
	}

	// Execute tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tidyText(tt.input)
			if got != tt.want {
				t.Errorf("tidyText(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
