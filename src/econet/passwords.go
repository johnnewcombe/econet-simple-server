package econet

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/johnnewcombe/econet-simple-server/src/lib"
)

const (
	DefaultPrivilege byte = 0b11000000
	MaxFreeSpace     int  = 1024e3
)

// Users is a structure that encompasses the econet PASSWORDS file
type Users struct {
	passwordFilePath  string // just so the o
	HomeDirectoryPath string
	Items             []User
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

func (p *Users) saveToDisk() error {

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

func (p *Users) UserExists(username string) bool {
	for _, pwd := range p.Items {
		if pwd.Username == username {
			return true
		}
	}
	return false
}

// AuthenticateUser Returns the password for the specified user or nil if user does not exist
func (p *Users) AuthenticateUser(username string, password string) *User {
	for _, pwd := range p.Items {
		if pwd.Username == username && pwd.Password == password {
			return &pwd
		}
	}
	return nil
}

func (p *Users) AddUser(username string, password string) error {

	var (
		err error
	)

	// create a new user obj
	user := User{
		Username:   strings.ToUpper(username),
		Password:   strings.ToUpper(password),
		FreeSpace:  MaxFreeSpace, // TODO is this an OK value to return to a client
		BootOption: DefaultBootOption,
		Privilege:  DefaultPrivilege,
	}

	p.Items = append(p.Items, user)

	if err = p.saveToDisk(); err != nil {
		return err
	}

	return nil
}

func (p *Users) ToString() string {
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

func NewUsers(pwFilePath string) (Users, error) {
	var (
		err      error
		userData string
		users    Users
	)

	// new file
	users = Users{}
	// set the path that was used so that we cal save it later
	users.passwordFilePath = pwFilePath

	if lib.Exists(pwFilePath) {
		if userData, err = lib.ReadString(pwFilePath); err != nil {
			return Users{}, err
		}

		// load the users
		if users, err = parseUsers(userData); err != nil {
			return Users{}, err
		}

	} else {

		// add user, this will
		if err = users.AddUser(DefaultSystemUserName, DefaultSystemPassword); err != nil {
			return Users{}, err
		}

		if err = users.saveToDisk(); err != nil {
			return Users{}, err
		}
	}

	return users, nil
}

func parseUsers(passwordData string) (Users, error) {

	var (
		err   error
		i     int
		user  User
		users Users
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
			return Users{}, fmt.Errorf("bad password file")
		}

		// create the user from the line
		user = User{
			Username: lines[0],
			Password: lines[1],
		}

		// add the free space
		i, err = strconv.Atoi(lines[2])
		if err != nil {
			return Users{}, err
		}
		user.FreeSpace = i

		//add the Option
		if i, err = strconv.Atoi(lines[3]); err != nil {
			return Users{}, err
		}
		user.BootOption = byte(i) & 0b00001111

		//add the Privilege
		if i, err = strconv.Atoi(lines[3]); err != nil {
			return Users{}, err
		}
		user.Privilege = byte(i) & 0b11110000

		users.Items = append(users.Items, user)
	}
	return users, nil
}

/*

	if !lib.Exists(pwFile) {

		slog.Info("Creating new password file.", "password-file", pwFile)
		// create a new password file
		user := econet.User{
			Username:   "SYST",
			Password:   "SYST",
			FreeSpace:  1024e3,
			BootOption: 0b00000000,
			Privilege:  0b11000000,
		}

		// add the user to the userData
		userData := econet.Users{
			Items: []econet.User{user},
		}

		// write the userData to disk
		s := userData.ToString()
		if err = lib.WriteString(pwFile, s); err != nil {
			return err
		}

		// create a home directory for the new user
		if err = lib.CreateDirectoryIfNotExists(econet.LocalDisk0 + user.Username); err != nil {
			return err
		}
	}

*/
