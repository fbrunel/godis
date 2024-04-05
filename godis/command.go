package godis

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
)

type operation func(args []string, st Store) (*Reply, error)

type CommandService struct {
	store      Store
	operations map[string]operation
}

func NewCommandService(store Store) *CommandService {
	return &CommandService{
		store:      store,
		operations: defaultOperations(),
	}
}

func defaultOperations() map[string]operation {
	return map[string]operation{
		"SET":     operationSet,
		"GET":     operationGet,
		"MGET":    operationMGet,
		"DEL":     operationDelete,
		"EXISTS":  operationExists,
		"INCR":    operationIncr,
		"DECR":    operationDecr,
		"KEYS":    operationKeys,
		"DBSIZE":  operationDbSize,
		"FLUSHDB": operationFlushDb,
		"INFO":    operationInfo,
	}
}

func operationSet(args []string, st Store) (*Reply, error) {
	if len(args) != 2 {
		return NewReplyErr(ErrWrongArgs), nil
	}

	st.Set(args[0], args[1])

	return NewReplyOK(), nil
}

func operationGet(args []string, st Store) (*Reply, error) {
	if len(args) != 1 {
		return NewReplyErr(ErrWrongArgs), nil
	}

	if !st.Exists(args[0]) {
		return NewReplyNil(), nil
	}

	return NewReply(st.Get(args[0])), nil
}

func operationMGet(args []string, st Store) (*Reply, error) {
	if len(args) < 1 {
		return NewReplyErr(ErrWrongArgs), nil
	}

	keys := make([]any, 0, len(args))
	for _, k := range args {
		if st.Exists(k) {
			keys = append(keys, st.Get(k))
		} else {
			keys = append(keys, nil)
		}
	}

	return NewReplyArray(keys), nil
}

func operationDelete(args []string, st Store) (*Reply, error) {
	if len(args) < 1 {
		return NewReplyErr(ErrWrongArgs), nil
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

func operationExists(args []string, st Store) (*Reply, error) {
	if len(args) < 1 {
		return NewReplyErr(ErrWrongArgs), nil
	}

	var count int64 = 0
	for _, k := range args {
		if st.Exists(k) {
			count++
		}
	}

	return NewReplyInteger(count), nil
}

func operationIncr(args []string, st Store) (*Reply, error) {
	if len(args) != 1 {
		return NewReplyErr(ErrWrongArgs), nil
	}

	k := args[0]

	var val int64 = 0
	if st.Exists(k) {
		v, err := strconv.ParseInt(st.Get(k), 10, 64)
		if err != nil {
			return NewReplyErr(ErrWrongType), nil
		}
		val = v
	}

	val = val + 1
	st.Set(k, strconv.FormatInt(val, 10))

	return NewReplyInteger(val), nil
}

func operationDecr(args []string, st Store) (*Reply, error) {
	if len(args) != 1 {
		return NewReplyErr(ErrWrongArgs), nil
	}

	k := args[0]

	var val int64 = 0
	if st.Exists(k) {
		v, err := strconv.ParseInt(st.Get(k), 10, 64)
		if err != nil {
			return NewReplyErr(ErrWrongType), nil
		}
		val = v
	}

	val = val - 1
	st.Set(k, strconv.FormatInt(val, 10))

	return NewReplyInteger(val), nil
}

// pattern:
//
//	{ term }
//
// term:
//
//	'*'         matches any sequence of non-Separator characters
//	'?'         matches any single non-Separator character
//	'[' [ '^' ] { character-range } ']'
//	            character class (must be non-empty)
//	c           matches character c (c != '*', '?', '\\', '[')
//	'\\' c      matches character c
//
// character-range:
//
//	c           matches character c (c != '\\', '-', ']')
//	'\\' c      matches character c
//	lo '-' hi   matches character c for lo <= c <= hi
func operationKeys(args []string, st Store) (*Reply, error) {
	if len(args) != 1 {
		return NewReplyErr(ErrWrongArgs), nil
	}

	keys := st.Keys()
	pattern := args[0]
	matches := make([]any, 0, len(keys))
	for _, k := range keys {
		match, err := filepath.Match(pattern, k)
		if err != nil {
			return NewReplyErr(err.Error()), nil
		}
		if match {
			matches = append(matches, k)
		}
	}

	return NewReplyArray(matches), nil
}

func operationDbSize(args []string, st Store) (*Reply, error) {
	if len(args) != 0 {
		return NewReplyErr(ErrWrongArgs), nil
	}
	return NewReplyInteger(st.Size()), nil
}

func operationFlushDb(args []string, st Store) (*Reply, error) {
	if len(args) != 0 {
		return NewReplyErr(ErrWrongArgs), nil
	}

	st.Flush()

	return NewReplyOK(), nil
}

func operationInfo(args []string, st Store) (*Reply, error) {
	if len(args) != 0 {
		return NewReplyErr(ErrWrongArgs), nil
	}

	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	fmt := fmt.Sprintf("HeapAlloc: %d bytes", mem.HeapAlloc)

	return NewReply(fmt), nil
}

//

func (srv *CommandService) ExecCommand(c Command) (*Reply, error) {
	op, exists := srv.operations[c.Op]
	if !exists {
		return NewReplyErr(ErrUnknownCmd), nil
	}
	return op(c.Args, srv.store)
}
