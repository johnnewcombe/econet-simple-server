package admin

import (
	"fmt"
	"strconv"
	"strings"
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

func (u *Users) Load(passwordData string) error {

	var (
		err  error
		i    int
		user User
	)

	for _, line := range strings.Split(passwordData, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		lines := strings.Split(line, ":")
		if len(lines) != 4 {
			return fmt.Errorf("bad password file")
		}

		// create the user from the line
		user = User{
			Username: lines[0],
			Password: lines[1],
		}

		// add the free space
		i, err = strconv.Atoi(lines[2])
		if err != nil {
			return err
		}
		user.FreeSpace = i

		//add the Option
		i, err = strconv.Atoi(lines[3])
		if err != nil {
			return err
		}
		user.BootOption = byte(i) & 0b00001111

		//add the Privilege
		i, err = strconv.Atoi(lines[3])
		if err != nil {
			return err
		}
		user.Privilege = byte(i) & 0b11110000

		u.Users = append(u.Users, user)
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
