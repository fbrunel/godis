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
	http     http.Server
	urlpath  string
	dumpfile string
	store    *StandardStore
	service  *CommandService
	handler  *CommandHandler
}

func NewServer(addr string) *Server {
	return &Server{
		http:     http.Server{Addr: addr},
		urlpath:  "/cmd",
		dumpfile: "/tmp/godis.dump",
	}
}

func (s *Server) Start() error {
	s.init()

	errch := make(chan error)
	go func() {
		log.Printf("-- serv: %s", s.http.Addr)
		err := s.http.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			errch <- err
		}
	}()

	select {
	case err := <-errch:
		return err
	case <-signals():
		return s.shutdown()
	}
}

//

func (s *Server) init() {
	log.Printf("-- load: %s", s.dumpfile)
	store, err := LoadStoreFromFile(s.dumpfile)
	if err != nil {
		log.Printf("EE %v", err)
		store = NewStandardStore()
	}

	s.store = store
	s.service = NewCommandService(s.store)

	ctx, cancel := context.WithCancel(context.Background())
	s.handler = NewCommandHandler(ctx, s.service)

	router := http.NewServeMux()
	router.Handle(s.urlpath, s.handler)
	s.http.Handler = router
	s.http.RegisterOnShutdown(func() { cancel() })
}

func (s *Server) shutdown() error {
	log.Printf("-- shutting down")
	err := s.http.Shutdown(context.Background())
	if err != nil {
		return err
	}

	log.Printf("-- save: %s", s.dumpfile)
	err = SaveStoreToFile(s.store, s.dumpfile)
	if err != nil {
		return err
	}

	time.Sleep(1 * time.Second)
	return nil
}

func signals() chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	return ch
}
