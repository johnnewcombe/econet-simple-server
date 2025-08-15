package econet

import (
	"github.com/google/uuid"
)

// private constants
const (
	DefaultSystemUserName         string = "SYST"
	DefaultSystemPassword         string = "SYST"
	DefaultBootOption             byte   = 0
	DefaultUserRootDirHandle      byte   = 1
	DefaultCurrentDirectoryHandle byte   = 2
	DefaultCurrentLibraryHandle   byte   = 4
	DefaultRootDirectory                 = "$"
	DefaultLibraryDirectory              = "DISK0/LIBRARY"
	PasswordFile                         = "PASSWORD"
	Disk0                                = "DISK0"
	Disk1                                = "DISK1"
	Disk2                                = "DISK2"
	Disk3                                = "DISK3"
)

// public variables
var (
	// LocalRootDiectory this is the local folder name NOT econet name
	LocalRootDiectory string
	LocalDisk0        string
	LocalDisk1        string
	LocalDisk2        string
	LocalDisk3        string

	Userdata       Users
	ActiveSessions Sessions
)

type HandleType byte

const (
	File HandleType = iota
	UserRootDirectory
	CurrentSelectedDirectory
	CurrentSelectedLibrary
)

type Sessions struct {
	items []Session
}

type Session struct {
	SessionId  uuid.UUID
	Username   string
	StationId  byte
	NetworkId  byte
	handles    map[byte]Handle
	BootOption byte
	DataPort   byte // this will be non-zero during a data transfer
}

type Handle struct {
	EconetPath string
	Type       HandleType
}

func NewSession(username string, stationId byte, networkId byte) *Session {

	// create a new map of file handles for the session and set the three defaults
	handles := make(map[byte]Handle)
	//handles[DefaultUserRootDirHandle] = DefaultRootDirectory
	//handles[DefaultCurrentDirectoryHandle] = DefaultRootDirectory + "." + username
	//handles[DefaultCurrentLibraryHandle] = DefaultRootDirectory + "." + DefaultLibraryDirectory

	return &Session{
		SessionId:  uuid.New(),
		Username:   username,
		StationId:  stationId,
		NetworkId:  networkId,
		handles:    handles,
		BootOption: DefaultBootOption,
	}
}

// GetSession  Returns the session for the specified user or nil if the session doesn't exist
func (s *Sessions) GetSession(stationId byte, networkId byte) *Session {
	for _, session := range s.items {
		if session.StationId == stationId && session.NetworkId == networkId {
			return &session
		}
	}
	return nil
}

// IsLoggedOn Returns true id the station has been logged on. It does not check for any specific user
func (s *Sessions) IsLoggedOn(stationId byte, networkId byte) bool {
	for _, session := range s.items {
		if session.StationId == stationId && session.NetworkId == networkId {
			return true
		}
	}
	return false
}

func (s *Sessions) AddSession(username string, stationId byte, networkId byte) *Session {

	session := *NewSession(username, stationId, networkId)
	s.items = append(s.items, session)
	return &session

}

func (s *Sessions) RemoveSession(session *Session) {
	for i, ses := range s.items {
		if ses.SessionId == session.SessionId {
			// to remove an item simply append everything before the specified index with everything after it
			s.items = append(s.items[:i], s.items[i+1:]...)
		}
	}
}

func (s *Session) getFreeHandle() byte {

	var (
		f byte
	)
	// looking for a free handle. Note that true is used in the for
	// statement simply because f<=255 is always true
	for f = 1; true; f++ {
		var _, ok = s.handles[f]

		// If the key exists ok will be true, we are looking for the
		// non-existence of a key i.e. ok=false
		if !ok {
			return f
		}
	}
	// no free handle so return zero, i.e. invalid handle
	return 0
}

// AddHandle Creates and adds a new file handle to the session for the specified
// file. Returns the file handle.
func (s *Session) AddHandle(econetFilePath string, handleType HandleType) byte {

	handle := Handle{
		EconetPath: econetFilePath,
		Type:       handleType,
	}

	key := s.getFreeHandle()
	s.handles[key] = handle

	return key
}

// RemoveHandle called when a user closes a file or directory in this session
func (s *Session) RemoveHandle(handle byte) {

	delete(s.handles, handle)
}

func (s *Session) GetUrd() string {

	for key, value := range s.handles {
		if s.handles[key].Type == UserRootDirectory {
			return value.EconetPath
		}
	}
	return ""
}

func (s *Session) GetCsd() string {

	for key, value := range s.handles {
		if s.handles[key].Type == CurrentSelectedDirectory {
			return value.EconetPath
		}
	}
	return ""
}
func (s *Session) GetCsl() string {

	for key, value := range s.handles {
		if s.handles[key].Type == CurrentSelectedLibrary {
			return value.EconetPath
		}
	}
	return ""
}
