package econet

import (
	"fmt"
	"github.com/johnnewcombe/econet-simple-server/src/utils"
	"strconv"
	"strings"
	"time"
)

var (
	RootFolder     string
	Userdata       Passwords
	ActiveSessions Sessions
)

type Sessions struct {
	Sessions []Session
}

type Session struct {
	Username  string
	StationId byte
	NetworkId byte
	// TODO work out how these are defined etc. clearly they change somehow when a user changes directory etc.
	//  but what about the UserRoot directory? Can this change?
	UserRootDirectory        byte
	CurrentSelectedDirectory byte
	CurrentSelectedLibrary   byte
}

func NewSession(username string, stationId byte, networkId byte, userRootDirectory byte, currentSelectedDirectory byte, currentSelectedLibrary byte) *Session {

	return &Session{
		Username:                 username,
		StationId:                stationId,
		NetworkId:                networkId,
		UserRootDirectory:        userRootDirectory,
		CurrentSelectedDirectory: currentSelectedDirectory,
		CurrentSelectedLibrary:   currentSelectedLibrary,
	}
}

// AuthenticateUser Returns the password for the specified user or nil if user does not exist
func (s *Sessions) GetSession(username string, stationId byte, networkId byte) *Session {
	for _, session := range s.Sessions {
		if session.StationId == stationId && session.NetworkId == networkId && session.Username == username {
			return &session
		}
	}
	return nil
}

type Passwords struct {
	Users []User
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
	for _, pwd := range p.Users {
		if pwd.Username == username {
			return true
		}
	}
	return false
}

// AuthenticateUser Returns the password for the specified user or nil if user does not exist
func (p *Passwords) AuthenticateUser(username string, password string) *User {
	for _, pwd := range p.Users {
		if pwd.Username == username && pwd.Password == password {
			return &pwd
		}
	}
	return nil
}

func (u *Passwords) ToString() string {
	result := strings.Builder{}
	for _, user := range u.Users {
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

		users.Users = append(users.Users, user)
	}
	return users, nil
}
