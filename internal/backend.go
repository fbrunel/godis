package internal

import (
	"errors"
)

type evalFunc func(args []string, st *Store) (string, error)

type Backend struct {
	Store *Store
	FMap  map[string]evalFunc
}

func NewDefaultBackend() *Backend {
	b := &Backend{
		NewStore(),
		make(map[string]evalFunc),
	}
	prepareEvalFunctions(b)
	return b
}

//

func prepareEvalFunctions(b *Backend) {
	b.FMap["SET"] = evalSET
	b.FMap["GET"] = evalGET
	b.FMap["DEL"] = evalDEL
}

func evalSET(args []string, st *Store) (string, error) {
	if len(args) < 2 {
		return "", errors.New("")
	}
	k, v := args[0], args[1]
	st.Set(k, v)
	return v, nil
}

func evalGET(args []string, st *Store) (string, error) {
	if len(args) < 1 {
		return "", errors.New("")
	}
	val := st.Get(args[0])
	return val, nil
}

func evalDEL(args []string, st *Store) (string, error) {
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
	eval, ok := b.FMap[c.Op]
	if ok {
		val, err := eval(c.Args, b.Store)
		if err != nil {
			return MakeReply("ERR", err.Error())
		}
		return MakeReply("OK!", val)
	}
	return MakeReply("ERR")
}
