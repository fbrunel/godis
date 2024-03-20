package internal

type Backend struct {
	Store *Store
}

func NewDefaultBackend() *Backend {
	return &Backend{
		NewStore(),
	}
}

func (b *Backend) EvalCommand(c Command) Reply {
	switch c.Op {
	case "SET":
		b.Store.Set(c.Args[0], c.Args[1])
		return MakeReply("OK")
	case "GET":
		val := b.Store.Get(c.Args[0])
		return MakeReply("OK", val)
	case "DEL":
		for _, a := range c.Args {
			b.Store.Del(a)
		}
		return MakeReply("OK")
	default:
		return MakeReply("ERR")
	}
}
