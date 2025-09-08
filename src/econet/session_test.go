package econet

import (
	"fmt"
	"os"
	"testing"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.

func Test_NewSession(t *testing.T) {
	// Defining the columns of the table

	//  create a new session, this should set up the default file handles
	session := *NewSession("JOHN", 100, 12)

	var tests = []struct {
		name        string
		handle      byte
		wantName    string
		wantStation byte
		wantNetwork byte
		wantCsd     string
		wantUrd     string
		wantCsl     string
	}{
		// the table itself
		{"Filename should be 'JOHN'", 1, "JOHN", 100, 12, "DISK0.$.JOHN", "DISK0.$.JOHN", "DISK0.$.LIBRARY"},
	}

	// The execution loop
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if session.Username != tt.wantName ||
				session.StationId != tt.wantStation ||
				session.NetworkId != tt.wantNetwork ||
				session.GetUrd() != tt.wantUrd ||
				session.GetCsd() != tt.wantCsd ||
				session.GetCsl() != tt.wantCsl {
				t.Errorf("got %s, want %s, got %d, want %d, got %d, want %d, got %s, want %s, got %s, want %s, got %s, want %s",
					session.Username,
					tt.wantName,
					session.StationId,
					tt.wantStation,
					session.NetworkId,
					tt.wantNetwork,
					session.GetUrd(),
					tt.wantUrd,
					session.GetCsd(),
					tt.wantCsd,
					session.GetCsl(),
					tt.wantCsl,
				)
			}
		})
	}
}

func Test_AddSession(t *testing.T) {

	sessions := Sessions{}

	var tests = []struct {
		name      string
		inputName string
		inputStn  byte
		wantName  string
		wantStn   byte
	}{
		{"JOHN at Station 100", "JOHN", 100, "JOHN", 100},
		{"JOHN at Station 64", "JOHN", 64, "JOHN", 64},
		{"SYST at Station 12", "SYST", 12, "SYST", 12},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ans := NewSession(tt.inputName, tt.inputStn, 0)
			sessions.AddSession(ans)

			// TODO use ans.SessionId to get the session back from the collection?

			if ans.Username != tt.wantName || ans.StationId != tt.wantStn || len(sessions.items) != i+1 {
				t.Errorf("got %s, want %s, got %d, want %d, got %d, want %d", ans.Username, tt.wantName, ans.StationId, tt.inputStn, len(sessions.items), i+1)
			}
		})
	}

}
func Test_RemoveSession(t *testing.T) {

	// two sessions from same user different machines
	session1 := *NewSession("JOHN", 100, 0)
	session2 := *NewSession("JOHN", 64, 0)
	session3 := *NewSession("SYST", 12, 0)

	sessions := Sessions{}
	sessions.items = append(sessions.items, session1)
	sessions.items = append(sessions.items, session2)
	sessions.items = append(sessions.items, session3)

	var tests = []struct {
		name  string
		input *Session
		want  int
	}{
		// one less session
		{"Remove JOHN at stn 100 from session", &session1, 2},
		{"Remove JOHN at stn 64 from session", &session2, 1},
		{"Not remove JOHN at stn 64 from session", &session2, 1},
		{"Remove SYST at stn 12 from session", &session3, 0},
		{"Not remove JOHN at stn 64 from session", &session1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			sessions.RemoveSession(tt.input)

			if len(sessions.items) != tt.want {
				t.Errorf("got %d, want %d", len(sessions.items), tt.want)
			}
		})
	}
}

func Test_getFreeHandle(t *testing.T) {

	//  create a new session, this should set up the default file handles
	session := *NewSession("JOHN", 100, 0)

	var tests = []struct {
		name string

		want byte
	}{
		// handles 1, 2 and 4 are already allocated to URD, CSD and CSL
		{"Handle should be 1", 4},
		{"Handle should be 2", 5},
		{"Handle should be 3", 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ans := session.getFreeHandle()
			session.AddHandle("$.MYFILE", File, false)

			if ans != tt.want {
				t.Errorf("got %d, want %d", ans, tt.want)
			}
		})
	}

	//session.RemoveHandle(1)
	//session.RemoveHandle(2)
	//session.RemoveHandle(3)
	//
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//
	//		ans := session.getFreeHandle()
	//		session.AddHandle("$.MYFILE", File, false)
	//
	//		if ans != tt.want {
	//			t.Errorf("got %d, want %d", ans, tt.want)
	//		}
	//	})
	//}

}

func Test_GetSession(t *testing.T) {

	// two sessions from same user different machines
	session1 := *NewSession("JOHN", 100, 0)
	session2 := *NewSession("JOHN", 64, 0)
	session3 := *NewSession("SYST", 12, 0)

	sessions := Sessions{}
	sessions.items = append(sessions.items, session1)
	sessions.items = append(sessions.items, session2)
	sessions.items = append(sessions.items, session3)

	var tests = []struct {
		name      string
		inputName string
		inputStn  byte
		wantName  string
		wantStn   byte
	}{
		{"JOHN at Station 100", "JOHN", 100, "JOHN", 100},
		{"JOHN at Station 64", "JOHN", 64, "JOHN", 64},
		{"SYST at Station 12", "SYST", 12, "SYST", 12},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ans := sessions.GetSession(tt.inputStn, 0)
			if ans.Username != tt.wantName {
				t.Errorf("got %s, want %s", ans.Username, tt.wantName)
			}
			if ans.StationId != tt.wantStn {
				t.Errorf("got %d, want %d", ans.StationId, tt.wantStn)
			}
		})
	}

}

func Test_AddHandle(t *testing.T) {
	// Test Data
	const (
		testUsername = "JOHN"
		testStation  = byte(100)
		testNetwork  = byte(0)
	)

	type testCase struct {
		name     string
		fileName string
		want     byte
		desc     string
	}

	tests := []testCase{
		{
			name:     "First file handle",
			fileName: "MYFILE1.txt",
			want:     4,
			desc:     "First handle should be 1",
		},
		{
			name:     "Second file handle",
			fileName: "MYFILE2.txt",
			want:     5,
			desc:     "Second handle should be 2",
		},
		{
			name:     "Third file handle",
			fileName: "MYFILE3.txt",
			want:     6,
			desc:     "Third handle should be 3",
		},
	}

	// Create a fresh session for all tests
	session := NewSession(testUsername, testStation, testNetwork)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Construct full file path according to the system convention
			fullPath := fmt.Sprintf("$.%s.%s", testUsername, tt.fileName)

			// When: adding a new file handle
			got := session.AddHandle(fullPath, File, false)

			// Then: verify the handle number
			if got != tt.want {
				t.Errorf("%s: got handle %d, want %d - %s",
					tt.name, got, tt.want, tt.desc)
			}
		})
	}
}
func Test_DeleteHandle(t *testing.T) {

	session := *NewSession("JOHN", 100, 0)
	//session.handles[0] = Handle{EconetPath: "$"}
	//session.handles[1] = Handle{EconetPath: "$"}
	//session.handles[2] = Handle{EconetPath: "$"}

	var tests = []struct {
		name  string
		input byte
		want  int
	}{
		{"Handle count should be 3", 0, 3},
		{"Handle count should be 2", 1, 2},
		{"Handle count should be 2", 1, 2},
		{"Handle count should be 1", 2, 1},
		{"Handle count should be 1", 2, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			session.RemoveHandle(tt.input)

			if len(session.handles) != tt.want {
				t.Errorf("got %d, want %d", len(session.handles), tt.want)
			}
		})
	}
}

func Test_EconetPathToLocalPath_Table(t *testing.T) {

	// Make LocalRootDiectory deterministic for tests
	LocalRootDiectory = "filestore"

	type tc struct {
		name       string
		input      string
		csd        string // CSD Econet path, e.g., "$.JOHN" (used only when input has no '$')
		expectDisc string
		expectPath string
		wantErr    bool
	}

	cwd, _ := os.Getwd()

	tests := []tc{
		{
			name:       "Valid disk prefix with $",
			input:      ":DISK1.$.MYDIR.FILE",
			csd:        "$.JOHN",
			expectPath: "filestore/DISK1/MYDIR/FILE",
			wantErr:    false,
		},
		{
			name:       "Valid disk prefix but no $",
			input:      ":DISK1.MYDIR.FILE",
			csd:        "$.JOHN",
			expectPath: "filestore/DISK1/MYDIR/FILE",
			wantErr:    false,
		},
		{
			name:       "No disk prefix, relative path",
			input:      "MYDIR.SUBDIR.FILE",
			csd:        "$.JOHN",
			expectPath: "filestore/DISK0/JOHN/MYDIR/SUBDIR/FILE",
			wantErr:    false,
		},
		{
			name:       "user root file, relative path",
			input:      "$.JOHN.FILE",
			csd:        "$.JOHN",
			expectPath: "filestore/DISK0/JOHN/FILE",
			wantErr:    false,
		},
		{
			name:       "filename only, relative path",
			input:      "FILE",
			csd:        "$.JOHN",
			expectPath: "filestore/DISK0/JOHN/FILE",
			wantErr:    false,
		},
		{
			name:       "Disk prefix, with double colons",
			input:      "::DISK1.MYDIR.SUBDIR.FILE",
			csd:        "$.JOHN",
			expectPath: "filestore/DISK1/MYDIR/SUBDIR/FILE",
			wantErr:    true,
		},
		{
			name:       "Disk prefix, with double dots",
			input:      ":DISK1..MYDIR.SUBDIR.FILE",
			csd:        "$.JOHN",
			expectPath: "filestore/DISK1/MYDIR/SUBDIR/FILE",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			s := &Session{handles: make(map[byte]Handle)}
			s.CurrentDisk = Disk0

			// Set CSD if provided, to emulate environment for relative paths
			if tt.csd != "" {
				s.handles[DefaultCurrentDirectoryHandle] = Handle{EconetPath: tt.csd, Type: CurrentSelectedDirectory}
			}

			got, err := s.EconetPathToLocalPath(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got none (got=%q)", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != cwd+"/"+tt.expectPath {
				t.Errorf("got %q, want %q", got, cwd+"/"+tt.expectPath)
			}
		})
	}
}
