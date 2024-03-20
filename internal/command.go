package internal

import (
	"strings"
)

type Command struct {
	Op   string   `json:"op"`
	Args []string `json:"args"`
}

func MakeCommand(op string, args ...string) Command {
	return Command{strings.ToUpper(op), args}
}

//

type Reply struct {
	Status string   `json:"status"`
	Data   []string `json:"data"`
}

func MakeReply(status string, data ...string) Reply {
	return Reply{status, data}
}
