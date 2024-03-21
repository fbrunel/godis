package internal

import (
	"errors"
)

type evalFunc func(args []string, st *Store) (string, error)

type Backend struct {
	Store   *Store
	EvalFns map[string]evalFunc
}

func NewDefaultBackend() *Backend {
	return &Backend{
		Store:   NewStore(),
		EvalFns: defaultEvalFns(),
	}
}

//

func defaultEvalFns() map[string]evalFunc {
	return map[string]evalFunc{
		"SET": evalSet,
		"GET": evalGet,
		"DEL": evalDel,
	}
}

func evalSet(args []string, st *Store) (string, error) {
	if len(args) < 2 {
		return "", errors.New("not enough arguments for SET")
	}
	k, v := args[0], args[1]
	st.Set(k, v)
	return v, nil
}

func evalGet(args []string, st *Store) (string, error) {
	if len(args) < 1 {
		return "", errors.New("not enough arguments for GET")
	}
	return st.Get(args[0]), nil
}

func evalDel(args []string, st *Store) (string, error) {
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
	fn, exists := b.EvalFns[c.Op]
	if !exists {
		return MakeReply("ERR", "unknown command")
	}

	val, err := fn(c.Args, b.Store)
	if err != nil {
		return MakeReply("ERR", err.Error())
	}
	return MakeReply("OK!", val)
}
