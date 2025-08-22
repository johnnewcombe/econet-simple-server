package econet

import "strings"

type CliCmd struct {
	Cmd     string
	Args    []string
	CmdText string
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
