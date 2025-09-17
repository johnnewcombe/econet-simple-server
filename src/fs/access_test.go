package fs

import "testing"

func TestIsOwner(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		user     string
		want     bool
	}{
		{name: "owner_simple", filename: "$.alice.FILE", user: "alice", want: true},
		{name: "owner_nested", filename: "$.bob.dir.sub.FILE", user: "bob", want: true},
		{name: "different_user", filename: "$.charlie.FILE", user: "dave", want: false},
		{name: "username_prefix_not_exact", filename: "$.usernamex.FILE", user: "username", want: true},
		{name: "missing_dollar_prefix", filename: "alice.FILE", user: "alice", want: false},
		{name: "username_later_not_owner", filename: "$.x.alice.FILE", user: "alice", want: false},
		{name: "empty_username_exact_prefix", filename: "$.", user: "", want: false},
		{name: "username_exact_prefix", filename: "$.", user: "Alice", want: false},
		{name: "empty_username_non_matching", filename: "$.x", user: "", want: false},
		{name: "case_sensitive_differs", filename: "$.Alice.FILE", user: "alice", want: false},
	}

	for _, tc := range tests {
		// capture tc
		t.Run(tc.name, func(t *testing.T) {
			got := IsOwner(tc.filename, tc.user)
			if got != tc.want {
				t.Errorf("IsOwner(%q, %q) = %v, want %v", tc.filename, tc.user, got, tc.want)
			}
		})
	}
}
