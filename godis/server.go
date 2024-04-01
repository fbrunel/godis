package godis

import (
	"log"
	"net/http"
)

type Server struct {
	addr string
	path string
}

func NewServer(addr string) *Server {
	return &Server{
		addr: addr,
		path: "/cmd",
	}
}

func (s *Server) Serve() error {
	store := NewStandardStore()
	service := NewCommandService(store)
	handle := NewCommandHandler(service)

	http.Handle(s.path, handle)
	log.Printf("-- serv: %s", s.addr)
	return http.ListenAndServe(s.addr, nil)
}
