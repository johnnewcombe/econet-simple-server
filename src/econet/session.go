package econet

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/johnnewcombe/econet-simple-server/src/utils"
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
	DefaultLibraryDirectory            = "LIBRARY"
	PasswordFile                       = "PASSWORD"
)

// public variables
var (
	// RootFolder this is the local folder name NOT econet name
	RootFolder     string
	Userdata       Passwords
	ActiveSessions Sessions
)

/*
type HandleType byte

const (

	UserRootDirectory HandleType = iota
	CurrentSelectedDirectory
	CurrrentSelectedLibrary
	StandardDirectore
	StandardFile

)

	type FileHandle struct {
		EconetPath string
		HandleType HandleType
	}
*/
type Sessions struct {
	items []Session
}

type Session struct {
	SessionId  uuid.UUID
	Username   string
	StationId  byte
	NetworkId  byte
	handles    map[byte]string
	BootOption byte
}

func NewSession(username string, stationId byte, networkId byte) *Session {

	// create a new map of file handles for the session and set the three defaults
	handles := make(map[byte]string)
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

// GetSession AuthenticateUser Returns the password for the specified user or nil if user does not exist
func (s *Sessions) GetSession(username string, stationId byte, networkId byte) *Session {
	for _, session := range s.items {
		if session.StationId == stationId && session.NetworkId == networkId && session.Username == username {
			return &session
		}
	}
	return nil
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
func (s *Session) AddHandle(econetFile string) byte {

	handle := s.getFreeHandle()
	s.handles[handle] = econetFile
	return handle

}

// RemoveHandle called when a user closes a file or directory in this session
func (s *Session) RemoveHandle(handle byte) {

	delete(s.handles, handle)
}

type Passwords struct {
	Items []User
}
type User struct {
	Username  string // 20 bytes max
	Password  string // 6 bytes max
	FreeSpace int    // max users space

	// these two are combined into a single byte in Acorn fileservers
	BootOption byte // uses the lower four bits
	Privilege  byte // uses the upper four bits
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

	if userData, err = utils.ReadString(pwFilePath); err != nil {
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
