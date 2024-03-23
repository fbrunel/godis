package internal

import "strconv"

type operation func(args []string, st *Store) (*Reply, error)

type CommandService struct {
	store      *Store
	operations map[string]operation
}

func NewCommandService() *CommandService {
	return &CommandService{
		store:      NewStore(),
		operations: defaultOperations(),
	}
}

//

func defaultOperations() map[string]operation {
	return map[string]operation{
		"SET":    operationSet,
		"GET":    operationGet,
		"DEL":    operationDelete,
		"EXISTS": operationExists,
		"INCR":   operationIncr,
		"DECR":   operationDecr,
	}
}

func operationSet(args []string, st *Store) (*Reply, error) {
	if len(args) < 2 {
		return NewReplyErr("not enough arguments for SET"), nil
	}

	k, v := args[0], args[1]
	st.Set(k, v)

	return NewReplyOK(), nil
}

func operationGet(args []string, st *Store) (*Reply, error) {
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

func operationDelete(args []string, st *Store) (*Reply, error) {
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

func operationExists(args []string, st *Store) (*Reply, error) {
	if len(args) < 1 {
		return NewReplyErr("not enough arguments for EXISTS"), nil
	}

	exists := st.Exists(args[0])
	if exists {
		return NewReplyInteger(1), nil
	}
	return NewReplyInteger(0), nil
}

func operationIncr(args []string, st *Store) (*Reply, error) {
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

func operationDecr(args []string, st *Store) (*Reply, error) {
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
	op, exists := srv.operations[c.Op]
	if !exists {
		return NewReplyErr("unknown command"), nil
	}
	return op(c.Args, srv.store)
}
