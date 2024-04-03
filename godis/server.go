package godis

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	http http.Server
	path string
}

func NewServer(addr string) *Server {
	return &Server{
		http: http.Server{Addr: addr},
		path: "/cmd",
	}
}

func (s *Server) Serve() error {
	store := NewStandardStore()
	service := NewCommandService(store)

	ctx, cancel := context.WithCancel(context.Background())
	handle := NewCommandHandler(ctx, service)

	router := http.NewServeMux()
	router.Handle(s.path, handle)
	s.http.Handler = router
	s.http.RegisterOnShutdown(func() { cancel() })

	go func() {
		log.Printf("-- serv: %s", s.http.Addr)
		err := s.http.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("EE %v", err)
		}
	}()

	signalch := make(chan os.Signal, 1)
	signal.Notify(signalch, syscall.SIGINT, syscall.SIGTERM)
	<-signalch

	log.Printf("-- shutting down")
	err := s.http.Shutdown(context.Background())
	time.Sleep(1 * time.Second)
	return err
}
