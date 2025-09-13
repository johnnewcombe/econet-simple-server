package econet

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/johnnewcombe/econet-simple-server/src/lib"
)

const (
	DefaultPrivilege byte = 0b11000000
	MaxFreeSpace     int  = 1024e3
)

// Passwords is a structure that encompasses the econet PASSWORDS file
type Passwords struct {
	passwordFilePath  string // just so the o
	HomeDirectoryPath string
	Items             []PWEntry
}
type PWEntry struct {
	Username     string // 20 bytes max
	Password     string // 6 bytes max
	FreeSpace    int    // max users space
	Option       byte   // combined with IsPrivileged uses the lower four bits
	IsPrivileged bool   // combined with BootOption uses the upper four bits
	//LoggedInAt time.Time
}

func (p *Passwords) saveToDisk() error {

	var (
		err error
	)

	if len(p.passwordFilePath) > 0 {
		if err = lib.WriteString(p.passwordFilePath, p.ToString()); err != nil {
			return err
		}
	}
	// write the userData to disk
	return nil
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
func (p *Passwords) AuthenticateUser(username string, password string) *PWEntry {
	for _, pwd := range p.Items {
		if pwd.Username == username && pwd.Password == password {
			return &pwd
		}
	}
	return nil
}

func (p *Passwords) AddUser(username string, password string) error {

	var (
		err error
	)

	// create a new user obj
	user := PWEntry{
		Username:     strings.ToUpper(username),
		Password:     strings.ToUpper(password),
		FreeSpace:    MaxFreeSpace, // TODO is this an OK value to return to a client
		Option:       0,
		IsPrivileged: false,
	}

	p.Items = append(p.Items, user)

	if err = p.saveToDisk(); err != nil {
		return err
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
		result.WriteString(strconv.Itoa(int(user.Option)))
		result.WriteString(strconv.FormatBool(user.IsPrivileged))

	}
	return result.String()
}

func NewUsers(pwFilePath string) (Passwords, error) {
	var (
		err      error
		userData string
		users    Passwords
	)

	// new file
	users = Passwords{}

	// set the path that was used so that we cal save it later
	users.passwordFilePath = pwFilePath

	if lib.Exists(pwFilePath) {
		if userData, err = lib.ReadString(pwFilePath); err != nil {
			return Passwords{}, err
		}

		// load the users
		if users, err = parseUsers(userData); err != nil {
			return Passwords{}, err
		}

	} else {

		// add user, this will
		if err = users.AddUser(DefaultSystemUserName, DefaultSystemPassword); err != nil {
			return Passwords{}, err
		}

		if err = users.saveToDisk(); err != nil {
			return Passwords{}, err
		}
	}

	return users, nil
}

func parseUsers(passwordData string) (Passwords, error) {

	var (
		err   error
		i     int
		user  PWEntry
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
		user = PWEntry{
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
		user.Option = byte(i)

		//add the IsPrivileged
		if i, err = strconv.Atoi(lines[3]); err != nil {
			return Passwords{}, err
		}
		user.IsPrivileged = byte(i)&0b01000000 > 0
		users.Items = append(users.Items, user)
	}
	return users, nil
}
