package econet

import (
	"fmt"
	"os"
	"regexp"
	"strings"

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
	DefaultLibraryDirectory              = "LIBRARY"
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

	Userdata       Passwords
	ActiveSessions Sessions
)

var (
	diskRootPathRegx = regexp.MustCompile(`^:[A-Za-z0-9]+\.\$\.[A-Za-z0-9]+`)
	diskNoRootRegx   = regexp.MustCompile(`^:[A-Za-z0-9]+\.[A-Za-z0-9]+`)
	rooRegx          = regexp.MustCompile(`^\$\.[A-Za-z0-9]+\.[A-Za-z0-9]+`)
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
	SessionId   uuid.UUID
	User        PWEntry
	StationId   byte
	NetworkId   byte
	handles     map[byte]Handle
	BootOption  byte
	CurrentDisk string
}

type Handle struct {
	EconetPath string
	Type       HandleType
	ReadOnly   bool
}

func NewSession(user PWEntry, stationId byte, networkId byte) *Session {

	// create a new map of file handles for the session and set the three defaults
	handles := make(map[byte]Handle)

	session := Session{
		SessionId:   uuid.New(),
		User:        user,
		StationId:   stationId,
		NetworkId:   networkId,
		handles:     handles,
		BootOption:  DefaultBootOption,
		CurrentDisk: Disk0,
	}

	// Note that the disk is part of the handles stored path
	urd := Disk0 + "." + DefaultRootDirectory + "." + user.Username
	csd := Disk0 + "." + DefaultRootDirectory + "." + user.Username
	csl := Disk0 + "." + DefaultRootDirectory + "." + DefaultLibraryDirectory

	session.AddHandle(urd, UserRootDirectory, false)
	session.AddHandle(csd, CurrentSelectedDirectory, false)
	session.AddHandle(csl, CurrentSelectedLibrary, false)

	return &session
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

func (s *Sessions) AddSession(session *Session) {
	s.items = append(s.items, *session)
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
func (s *Session) AddHandle(econetFilePath string, handleType HandleType, readOnly bool) byte {

	handle := Handle{
		EconetPath: econetFilePath,
		Type:       handleType,
		ReadOnly:   readOnly,
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

func (s *Session) HandleExists(econetPath string) bool {
	return s.GetHandle(econetPath) != nil
}

func (s *Session) GetHandle(econetPath string) *Handle {

	for _, value := range s.handles {
		if value.EconetPath == econetPath {
			return &value
		}
	}
	return nil
}

func (s *Session) EconetPathToLocalPath(econetPath string) (string, error) {

	var (
		diskName  string
		localRoot string
		localPath string
		cwd       string
		err       error
	)

	// Determine disk from optional leading "<digit>:" prefix.
	cwd, err = os.Getwd()
	if err != nil {
		return "", err
	}

	// get full path with disk name
	if diskName, econetPath, err = s.ExpandEconetPath(econetPath); err != nil {
		return "", err
	}

	// belts and braces check
	if !rooRegx.MatchString(econetPath) {
		return "", fmt.Errorf("invalid econet path")
	}

	localRoot = fmt.Sprintf("%s/%s", LocalRootDiectory, diskName)
	localPath = strings.Replace(econetPath, "$", localRoot, -1)
	localPath = cwd + "/" + strings.Replace(localPath, ".", "/", -1)

	// TODO This funcion only allows a-zA-Z0-9 as valid characters in the path.
	//  however the following chars are valid in econet paths
	//   ! % & = - ~ ^ | \ @ { [ £ _ + ; } ] < > ? / a-z A-Z 0-9
	//  need to consider linux,mac and windows characters to see what can safely be used

	return localPath, nil
}

// ExpandEconetPath Returns the  and root econet path for the specified file
func (s *Session) ExpandEconetPath(econetPath string) (string, string, error) {

	var (
		diskName string
	)

	// full path with disk name
	if diskRootPathRegx.MatchString(econetPath) {

		i := strings.Index(econetPath, ".")
		diskName = econetPath[1:i]
		econetPath = econetPath[i+1:]

	} else if diskNoRootRegx.MatchString(econetPath) {

		// when a disk is specified, the "$" is optional so if its missing, pop it back so that
		i := strings.Index(econetPath, ".")
		diskName = econetPath[1:i]
		econetPath = econetPath[i+1:]
		econetPath = "$." + econetPath

	} else if rooRegx.MatchString(econetPath) {
		diskName = s.CurrentDisk
	} else {

		// relative or invalid path so expand with csd and check
		// TODO the s.Csd() returns the current directory and includes the disk name
		diskName = s.CurrentDisk
		econetPath = s.GetCsd() + "." + econetPath

	}

	// belts and braces check
	if !rooRegx.MatchString(econetPath) {
		return "", "", fmt.Errorf("invalid econet path")
	}

	return diskName, econetPath, nil
}
