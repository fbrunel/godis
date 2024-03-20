package internal

import (
	"strings"
)

type Command struct {
	Cmd  string   `json:"cmd"`
	Args []string `json:"args"`
}

func MakeCommand(args ...string) Command {
	return Command{strings.ToUpper(args[0]), args[1:]}
}

//

type Reply struct {
	Code string `json:"code"`
	Data []any  `json:"data"`
}

func MakeReply(c string, data ...any) Reply {
	return Reply{c, data}
}
