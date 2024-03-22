package internal

import (
	"errors"
)

type cmdFunc func(args []string, st *Store) (string, error)

type CommandService struct {
	store  *Store
	cmdFns map[string]cmdFunc
}

func NewCommandService() *CommandService {
	return &CommandService{
		store:  NewStore(),
		cmdFns: defaultCmdFns(),
	}
}

//

func defaultCmdFns() map[string]cmdFunc {
	return map[string]cmdFunc{
		"SET": cmdSet,
		"GET": cmdGet,
		"DEL": cmdDel,
	}
}

func cmdSet(args []string, st *Store) (string, error) {
	if len(args) < 2 {
		return "", errors.New("not enough arguments for SET")
	}
	k, v := args[0], args[1]
	st.Set(k, v)
	return v, nil
}

func cmdGet(args []string, st *Store) (string, error) {
	if len(args) < 1 {
		return "", errors.New("not enough arguments for GET")
	}
	return st.Get(args[0]), nil
}

func cmdDel(args []string, st *Store) (string, error) {
	if len(args) < 1 {
		return "", errors.New("not enough arguments for DEL")
	}
	for _, a := range args {
		st.Del(a)
	}
	return "", nil
}

//

func (srv *CommandService) ExecCommand(c Command) Reply {
	fn, exists := srv.cmdFns[c.Op]
	if !exists {
		return MakeReply("ERR", "unknown command")
	}

	val, err := fn(c.Args, srv.store)
	if err != nil {
		return MakeReply("ERR", err.Error())
	}
	return MakeReply("OK!", val)
}
