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

const (
	TypeAck = "ACK"
	TypeNil = "NIL"
	TypeStr = "STR"
	TypeInt = "INT"
	TypeArr = "ARR"
	TypeErr = "ERR"
)

type Reply struct {
	Type string `json:"status"`
	Data any    `json:"data"`
}

func NewReplyOK() *Reply {
	return &Reply{
		Type: TypeAck,
		Data: "OK",
	}
}

func NewReplyNil() *Reply {
	return &Reply{
		Type: TypeNil,
		Data: nil,
	}
}

func NewReply(data string) *Reply {
	return &Reply{
		Type: TypeStr,
		Data: data,
	}
}

func NewReplyInteger(i int64) *Reply {
	return &Reply{
		Type: TypeInt,
		Data: i,
	}
}

func NewReplyArray(data ...string) *Reply {
	return &Reply{
		Type: TypeArr,
		Data: data,
	}
}

func NewReplyErr(msg string) *Reply {
	return &Reply{
		Type: TypeErr,
		Data: msg,
	}
}
