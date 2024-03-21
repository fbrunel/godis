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
	b := &Backend{
		NewStore(),
		defaultEvalFns(),
	}
	return b
}

//

func defaultEvalFns() map[string]evalFunc {
	fns := make(map[string]evalFunc)

	fns["SET"] = evalSet
	fns["GET"] = evalGet
	fns["DEL"] = evalDel

	return fns
}

func evalSet(args []string, st *Store) (string, error) {
	if len(args) < 2 {
		return "", errors.New("")
	}
	k, v := args[0], args[1]
	st.Set(k, v)
	return v, nil
}

func evalGet(args []string, st *Store) (string, error) {
	if len(args) < 1 {
		return "", errors.New("")
	}
	val := st.Get(args[0])
	return val, nil
}

func evalDel(args []string, st *Store) (string, error) {
	if len(args) < 1 {
		return "", errors.New("")
	}
	for _, a := range args {
		st.Del(a)
	}
	return "", nil
}

//

func (b *Backend) EvalCommand(c Command) Reply {
	fn, ok := b.EvalFns[c.Op]
	if ok {
		val, err := fn(c.Args, b.Store)
		if err != nil {
			return MakeReply("ERR", err.Error())
		}
		return MakeReply("OK!", val)
	}
	return MakeReply("ERR")
}
