package admin

import (
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
func (p *Users) GetUser(username string) *User {
	for _, password := range p.Users {
		if password.Username == username {
			return &password
		}
	}
	return nil
}

func (p *Users) Load(passwordData string) error {

	for _, line := range strings.Split(passwordData, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

	}
	return nil
}

func (p *Users) ToString() string {
	result := strings.Builder{}
	for _, user := range p.Users {
		result.WriteString(user.Username)
		result.WriteString(",")
		result.WriteString(user.Password)
		result.WriteString(",")
		result.WriteString(strconv.Itoa(user.FreeSpace))
		result.WriteString(",")
		result.WriteByte(user.BootOption)
		result.WriteString(",")
		result.WriteByte(user.Privilege)
	}
	return result.String()
}
