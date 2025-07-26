package econet

import (
	"fmt"
	"github.com/johnnewcombe/econet-simple-server/src/utils"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

const (
	PasswordFile = "PASSWORD"
)

var (
	RootFolder string
	Userdata   Users
)

type Users struct {
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

// GetUser Returns the password for the specified user or nil if user does not exist
func (u *Users) GetUser(username string) *User {
	for _, password := range u.Users {
		if password.Username == username {
			return &password
		}
	}
	return nil

}

func (u *Users) ToString() string {
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

func NewUsers(pwFilePath string) (Users, error) {
	var (
		err      error
		userData string
		users    Users
	)

	slog.Info("Loading password file.", "password-file", pwFilePath)
	if userData, err = utils.ReadString(pwFilePath); err != nil {
		return Users{}, err
	}

	// load the users
	if users, err = parseUsers(userData); err != nil {
		return Users{}, err
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

		users.Users = append(users.Users, user)
	}
	return users, nil
}
