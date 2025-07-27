package econet

import (
	"fmt"
	"log/slog"
	"strings"
)

func ParseCommand(command string, srcStationId byte, srcNetworkId byte) []byte {

	// PROCESS RX_TRANSMIT
	var (
		//cmds     []string
		password      string
		username      string
		data          []byte
		authenticated bool
	)

	// get logged in status
	session := ActiveSessions.GetSession(srcStationId, srcNetworkId)

	/*
		cmds = []string{"I AM", "BYE"}

		for _, cmd := range cmds {
			if strings.HasPrefix(command, cmd) {
				break
			}
		}
	*/

	if strings.HasPrefix(command, "I AM") {
		// parse the command
		args := strings.Split(command, " ")
		if len(args) == 3 {
			username = args[2]
			password = ""
		} else if len(args) == 4 {
			username = args[2]
			password = args[3]
		}

		// check user against users
		if session != nil {
			// already logged in
		}

		if user := Userdata.AuthenticateUser(username, password); user != nil {
			// user good
			data = []byte{0x05, 0x00, 0x01, 0x02, 0x04, 0x00}
			authenticated = true

			// TODO set session

		} else {
			// TODO Sort out correct responses for failed login
			// user not good
			data = []byte{0x05, 0x00, 0x01, 0x02, 0x04, 0x00}
		}

		slog.Info(fmt.Sprintf("econet-command=I AM %s, authenticated=%v", username, authenticated))

	} else if strings.HasPrefix(command, "NOTIFY") {

	}

	// TODO Remove dummy reply for a real one
	// TODO Better understand the control port
	// issue a dummy successful reply

	/*
		0x05, // indicates a successful login
		0x00, // return code of zero indicates success
		0x01, // user root dir handle
		0x02, // currently selected dir handle
		0x04, // library dir handle
		0x00, // boot option (0 = none)
	*/
	return data
}
func IAM() {

}
