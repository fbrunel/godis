package internal

import "strconv"

type cmdFunc func(args []string, st *Store) (*Reply, error)

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
		"SET":    cmdSet,
		"GET":    cmdGet,
		"DEL":    cmdDelete,
		"EXISTS": cmdExists,
		"INCR":   cmdIncr,
		"DECR":   cmdDecr,
	}
}

func cmdSet(args []string, st *Store) (*Reply, error) {
	if len(args) < 2 {
		return NewReplyErr("not enough arguments for SET"), nil
	}

	k, v := args[0], args[1]
	st.Set(k, v)

	return NewReplyOK(), nil
}

func cmdGet(args []string, st *Store) (*Reply, error) {
	if len(args) < 1 {
		return NewReplyErr("not enough arguments for GET"), nil
	}

	k := args[0]
	if !st.Exists(k) {
		return NewReplyNil(), nil
	}

	v := st.Get(args[0])
	return NewReply(v), nil
}

func cmdDelete(args []string, st *Store) (*Reply, error) {
	if len(args) < 1 {
		return NewReplyErr("not enough arguments for DEL"), nil
	}

	var count int64 = 0
	for _, k := range args {
		if st.Exists(k) {
			st.Delete(k)
			count++
		}
	}

	return NewReplyInteger(count), nil
}

func cmdExists(args []string, st *Store) (*Reply, error) {
	if len(args) < 1 {
		return NewReplyErr("not enough arguments for EXISTS"), nil
	}

	exists := st.Exists(args[0])
	if exists {
		return NewReplyInteger(1), nil
	}
	return NewReplyInteger(0), nil
}

func cmdIncr(args []string, st *Store) (*Reply, error) {
	if len(args) < 1 {
		return NewReplyErr("not enough arguments for INCR"), nil
	}

	k := args[0]
	var val int64 = 0
	if st.Exists(k) {
		v, err := strconv.ParseInt(st.Get(k), 10, 64)
		if err != nil {
			return NewReplyErr("WRONGTYPE operation against a key holding the wrong kind of value"), nil
		}
		val = v
	}
	val = val + 1
	st.Set(k, strconv.FormatInt(val, 10))
	return NewReplyInteger(val), nil
}

func cmdDecr(args []string, st *Store) (*Reply, error) {
	if len(args) < 1 {
		return NewReplyErr("not enough arguments for DECR"), nil
	}

	k := args[0]
	var val int64 = 0
	if st.Exists(k) {
		v, err := strconv.ParseInt(st.Get(k), 10, 64)
		if err != nil {
			return NewReplyErr("WRONGTYPE operation against a key holding the wrong kind of value"), nil
		}
		val = v
	}
	val = val - 1
	st.Set(k, strconv.FormatInt(val, 10))
	return NewReplyInteger(val), nil
}

//

func (srv *CommandService) ExecCommand(c Command) (*Reply, error) {
	cmd, exists := srv.cmdFns[c.Op]
	if !exists {
		return NewReplyErr("unknown command"), nil
	}

	rep, err := cmd(c.Args, srv.store)
	if err != nil {
		return nil, err
	}

	return rep, nil
}
