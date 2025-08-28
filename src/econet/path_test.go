package econet

import "testing"

func Test_HasDiskPrefix_Table(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want bool
	}{
		{name: "empty", in: "", want: false},
		{name: "just colon", in: ":", want: false},
		{name: "single digit then colon", in: "0:$.FILE", want: false},
		{name: "multi-digit not allowed", in: ":DISK0.MYFILE", want: true},
		{name: "non-digit before colon", in: "a:$", want: false},
		{name: "space before digit", in: " 0:$", want: false},
		{name: "digit without colon", in: ":0", want: false},
		{name: "digit-colon-semicolon", in: ":0$", want: false},
		{name: "colon then digit", in: ":0", want: false},
		{name: "dollar first", in: "$.X", want: false},
		{name: "double colon after digit", in: "0::REST", want: true},
		{name: "single digit then colon 2", in: ":XYZ.$", want: true},
	}

	for _, tc := range cases {
		// use tc.name as subtest label for clarity
		t.Run(tc.name, func(t *testing.T) {
			got := HasDiskPrefix(tc.in)
			if got != tc.want {
				t.Errorf("HasDiskPrefix(%q) = %v, want %v", tc.in, got, tc.want)
			}
		})
	}
}
