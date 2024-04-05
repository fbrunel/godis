package godis

import (
	"strings"
)

type Command struct {
	Op   string   `json:"o"`
	Args []string `json:"a"`
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

const (
	ErrUnknownCmd = "unknown command"
	ErrWrongArgs  = "wrong number of arguments for command"
	ErrWrongType  = "WRONGTYPE operation against a key holding the wrong kind of value"
)

type Reply struct {
	Value any    `json:"v"`
	Type  string `json:"t"`
}

func (r *Reply) Values() []any {
	return r.Value.([]any)
}

func NewReplyOK() *Reply {
	return &Reply{
		Type:  TypeAck,
		Value: "OK",
	}
}

func NewReplyNil() *Reply {
	return &Reply{
		Type:  TypeNil,
		Value: nil,
	}
}

func NewReply(str string) *Reply {
	return &Reply{
		Type:  TypeStr,
		Value: str,
	}
}

func NewReplyInteger(i int64) *Reply {
	return &Reply{
		Type:  TypeInt,
		Value: i,
	}
}

func NewReplyArray(values []any) *Reply {
	return &Reply{
		Type:  TypeArr,
		Value: values,
	}
}

func NewReplyErr(str string) *Reply {
	return &Reply{
		Type:  TypeErr,
		Value: str,
	}
}
