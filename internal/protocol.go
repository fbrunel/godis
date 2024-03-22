package internal

import (
	"strings"
)

type Command struct {
	Op   string   `json:"op"`
	Args []string `json:"args"`
}

func MakeCommand(op string, args ...string) Command {
	return Command{
		Op:   strings.ToUpper(op),
		Args: args,
	}
}

//

type Status string

const (
	StatusOK  = "OK!"
	StatusErr = "ERR"
)

type Reply struct {
	Status string `json:"status"`
	Data   any    `json:"data"`
}

func NewReplyMany(data ...any) *Reply {
	return &Reply{
		Status: StatusOK,
		Data:   data,
	}
}

func NewReplyOnce(data any) *Reply {
	return &Reply{
		Status: StatusOK,
		Data:   data,
	}
}

func NewReplyErr(msg string) *Reply {
	return &Reply{
		Status: StatusErr,
		Data:   msg,
	}
}
