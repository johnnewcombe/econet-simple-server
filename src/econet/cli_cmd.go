package econet

import "strings"

type CliCmd struct {
	Cmd     string
	Args    []string
	CmdText string
}

func NewCliCmd(commandText string) *CliCmd {

	var (
		commands []string
		cmdArgs  []string
		cmd      string
		ok       bool
		argText  string
	)

	commandText = tidyText(commandText)

	// list of piconet * commands
	commands = []string{"SAVE", "LOAD", "CAT", "INFO", "I AM", "SDISK", "DIR", "LIB"}

	for _, cmd = range commands {
		if _, argText, ok = strings.Cut(commandText, cmd); ok { // i.e. if ok
			cmdArgs = strings.Split(strings.Trim(argText, " "), " ")

			resultArgs := []string{}

			for _, arg := range cmdArgs {
				newArg := strings.Split(arg, "\r")[0]
				resultArgs = append(resultArgs, newArg)
			}

			return &CliCmd{
				Cmd:     cmd,
				CmdText: commandText,
				Args:    resultArgs,
			}

		}
	}
	return &CliCmd{}
}

func (c *CliCmd) ToBytes() []byte {
	return []byte(c.ToString())
}

func (c *CliCmd) ToString() string {

	str := strings.Builder{}
	str.WriteString(c.Cmd)

	if str.Len() > 0 {
		str.WriteString(" ")
	}

	for _, arg := range c.Args {
		str.WriteString(arg)
		str.WriteString(" ")
	}

	result := str.String()
	if len(result) > 0 {
		result = result[:len(result)-1]
	}

	return result
}
