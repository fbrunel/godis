package core

import (
	"encoding/json"
)

type Command struct {
	Cmd  string   `json:"cmd"`
	Args []string `json:"args"`
}

func MakeCommand(args ...string) Command {
	return Command{args[0], args[1:]}
}

func EncodeCommand(c Command) ([]byte, error) {
	bytes, err := json.Marshal(c)
	return bytes, err
}

func DecodeCommand(b []byte) (Command, error) {
	var c Command
	err := json.Unmarshal(b, &c)
	return c, err
}

//

type Reply struct {
	Code string   `json:"code"`
	Data []string `json:"data"`
}

func MakeReply(c string, data ...string) Reply {
	return Reply{c, data}
}

func EncodeReply(r Reply) ([]byte, error) {
	bytes, err := json.Marshal(r)
	return bytes, err
}

func DecodeReply(b []byte) (Reply, error) {
	var r Reply
	err := json.Unmarshal(b, &r)
	return r, err
}
