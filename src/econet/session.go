package econet

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/johnnewcombe/econet-simple-server/src/lib"

	//"github.com/johnnewcombe/econet-simple-server/src/cobra"
	"strconv"
	"strings"
	"time"
)

// private constants
const (
	DefaultBootOption             byte = 0
	DefaultUserRootDirHandle      byte = 1
	DefaultCurrentDirectoryHandle byte = 2
	DefaultCurrentLibraryHandle   byte = 4
	DefaultRootDirectory               = "$"
	DefaultLibraryDirectory            = "DISK0/LIBRARY"
	PasswordFile                       = "PASSWORD"
	Disk0                              = "DISK0"
	Disk1                              = "DISK1"
	Disk2                              = "DISK2"
	Disk3                              = "DISK3"
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

type Passwords struct {
	Items []User
}
type User struct {
	Username   string // 20 bytes max
	Password   string // 6 bytes max
	FreeSpace  int    // max users space
	Urd        byte
	Csd        byte
	Csl        byte
	BootOption byte // combined with Privilege uses the lower four bits
	Privilege  byte // combined with BootOption uses the upper four bits
	LoggedIn   bool
	LoggedInAt time.Time
}

func (p *Passwords) UserExists(username string) bool {
	for _, pwd := range p.Items {
		if pwd.Username == username {
			return true
		}
	}
	return false
}

// AuthenticateUser Returns the password for the specified user or nil if user does not exist
func (p *Passwords) AuthenticateUser(username string, password string) *User {
	for _, pwd := range p.Items {
		if pwd.Username == username && pwd.Password == password {
			return &pwd
		}
	}
	return nil
}

func (p *Passwords) ToString() string {
	result := strings.Builder{}
	for _, user := range p.Items {
		result.WriteString(user.Username)
		result.WriteString(":")
		result.WriteString(user.Password)
		result.WriteString(":")
		result.WriteString(strconv.Itoa(user.FreeSpace))
		result.WriteString(":")
		bootPriv := user.BootOption | user.Privilege
		result.WriteString(strconv.Itoa(int(bootPriv)))
	}
	return result.String()
}

func NewUsers(pwFilePath string) (Passwords, error) {
	var (
		err      error
		userData string
		users    Passwords
	)

	if userData, err = lib.ReadString(pwFilePath); err != nil {
		return Passwords{}, err
	}

	// load the users
	if users, err = parseUsers(userData); err != nil {
		return Passwords{}, err
	}

	return users, nil
}

func parseUsers(passwordData string) (Passwords, error) {

	var (
		err   error
		i     int
		user  User
		users Passwords
	)

	// TODO: Check specification in the comments within the password file
	//  and implement fully if appropriate
	for _, line := range strings.Split(passwordData, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		lines := strings.Split(line, ":")
		if len(lines) != 4 {
			return Passwords{}, fmt.Errorf("bad password file")
		}

		// create the user from the line
		user = User{
			Username: lines[0],
			Password: lines[1],
		}

		// add the free space
		i, err = strconv.Atoi(lines[2])
		if err != nil {
			return Passwords{}, err
		}
		user.FreeSpace = i

		//add the Option
		if i, err = strconv.Atoi(lines[3]); err != nil {
			return Passwords{}, err
		}
		user.BootOption = byte(i) & 0b00001111

		//add the Privilege
		if i, err = strconv.Atoi(lines[3]); err != nil {
			return Passwords{}, err
		}
		user.Privilege = byte(i) & 0b11110000

		users.Items = append(users.Items, user)
	}
	return users, nil
}
