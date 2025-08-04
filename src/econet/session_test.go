package econet

import (
	"testing"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.

func Test_NewSession(t *testing.T) {
	// Defining the columns of the table

	//  create a new session, this should set up the default file handles
	session := *NewSession("JOHN", 100, 0)

	var tests = []struct {
		name        string
		handle      byte
		wantName    string
		wantStation byte
	}{
		// the table itself
		{"Name should be 'JOHN'", 1, "JOHN", 100},
	}

	// The execution loop
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans1 := session.Username
			ans2 := session.StationId
			if ans1 != tt.wantName || ans2 != tt.wantStation {
				t.Errorf("got %s, want %s, got %d, want %d", ans1, tt.wantName, ans2, tt.wantStation)
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

			ans := sessions.AddSession(tt.inputName, tt.inputStn, 0)

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
		{"Handle should be 1", 1},
		{"Handle should be 2", 2},
		{"Handle should be 3", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ans := session.getFreeHandle()
			session.handles[ans] = ""
			if ans != tt.want {
				t.Errorf("got %d, want %d", ans, tt.want)
			}
		})
	}
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
	//  create a new session, this should set up the default file handles
	session := *NewSession("JOHN", 100, 0)

	var tests = []struct {
		name  string
		input Session
		want  byte
	}{
		{"Handle should be 1", session, 1},
		{"Handle should be 2", session, 2},
		{"Handle should be 3", session, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ans := session.AddHandle("$.JOHN.MYFILE")
			session.handles[ans] = ""
			if ans != tt.want {
				t.Errorf("got %d, want %d", ans, tt.want)
			}
		})
	}

}
func Test_DeleteHandle(t *testing.T) {

	session := *NewSession("JOHN", 100, 0)
	session.handles[0] = "$"
	session.handles[1] = "$"
	session.handles[2] = "$"

	var tests = []struct {
		name  string
		input byte
		want  int
	}{
		{"Handle count should be 3", 0, 2},
		{"Handle count should be 2", 1, 1},
		{"Handle count should be 2", 1, 1},
		{"Handle count should be 1", 2, 0},
		{"Handle count should be 1", 2, 0},
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
