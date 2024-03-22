package internal

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
	}
}

func cmdSet(args []string, st *Store) (*Reply, error) {
	if len(args) < 2 {
		return NewReplyErr("not enough arguments for SET"), nil
	}

	k, v := args[0], args[1]
	st.Set(k, v)

	return NewReplyOnce(v), nil
}

func cmdGet(args []string, st *Store) (*Reply, error) {
	if len(args) < 1 {
		return NewReplyErr("not enough arguments for GET"), nil
	}

	k := args[0]
	if !st.Exists(k) {
		return NewReplyOnce(nil), nil
	}

	v := st.Get(args[0])
	return NewReplyOnce(v), nil
}

func cmdDelete(args []string, st *Store) (*Reply, error) {
	if len(args) < 1 {
		return NewReplyErr("not enough arguments for DEL"), nil
	}

	count := 0
	for _, k := range args {
		if st.Exists(k) {
			st.Delete(k)
			count++
		}
	}

	return NewReplyOnce(count), nil
}

func cmdExists(args []string, st *Store) (*Reply, error) {
	if len(args) < 1 {
		return NewReplyErr("not enough arguments for EXISTS"), nil
	}

	return NewReplyOnce(st.Exists(args[0])), nil
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
