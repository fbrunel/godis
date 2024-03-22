package internal

import (
	"errors"
)

type cmdFunc func(args []string, st *Store) (string, error)

type Backend struct {
	store  *Store
	cmdFns map[string]cmdFunc
}

func NewBackend() *Backend {
	return &Backend{
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

func (b *Backend) EvalCommand(c Command) Reply {
	fn, exists := b.cmdFns[c.Op]
	if !exists {
		return MakeReply("ERR", "unknown command")
	}

	val, err := fn(c.Args, b.store)
	if err != nil {
		return MakeReply("ERR", err.Error())
	}
	return MakeReply("OK!", val)
}
