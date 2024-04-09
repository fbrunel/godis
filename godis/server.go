package godis

import (
	"context"
	"log"
	"net"
	"net/http"
)

type Server struct {
	http    http.Server
	store   Store
	service *CommandService
	handler *CommandHandler
}

func NewServer(store Store) *Server {
	srv := Server{
		http:  http.Server{},
		store: store,
	}

	srv.service = NewCommandService(srv.store)
	srv.handler = NewCommandHandler(srv.service)

	router := http.NewServeMux()
	router.Handle("/cmd", srv.handler)
	srv.http.Handler = router

	return &srv
}

func (srv *Server) Start(ctx context.Context, addr string) error {
	srv.http.Addr = addr
	srv.http.BaseContext = func(_ net.Listener) context.Context { return ctx }

	log.Printf("-- serv: %s", srv.http.Addr)
	err := srv.http.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (srv *Server) Shutdown(ctx context.Context) error {
	defer srv.handler.WaitClose()

	log.Printf("-- stop")
	err := srv.http.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}
