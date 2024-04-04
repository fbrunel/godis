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

type Options struct {
	Addr     string
	URLPath  string
	Dumpfile string
}

func DefaultOptions() Options {
	return Options{
		Addr:     ":8080",
		URLPath:  "/cmd",
		Dumpfile: "/tmp/godis.dump",
	}
}

//

type Server struct {
	opt     Options
	http    http.Server
	store   *StandardStore
	service *CommandService
	handler *CommandHandler
}

func NewServer(opt Options) *Server {
	return &Server{
		opt:  opt,
		http: http.Server{Addr: opt.Addr},
	}
}

func (s *Server) Start() error {
	s.init()

	errch := make(chan error)
	go func() {
		log.Printf("-- serv: %s", s.opt.Addr)
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
	log.Printf("-- load: %s", s.opt.Dumpfile)
	store, err := LoadStoreFromFile(s.opt.Dumpfile)
	if err != nil {
		log.Printf("EE %v", err)
		store = NewStandardStore()
	}

	s.store = store
	s.service = NewCommandService(s.store)

	ctx, cancel := context.WithCancel(context.Background())
	s.handler = NewCommandHandler(ctx, s.service)

	router := http.NewServeMux()
	router.Handle(s.opt.URLPath, s.handler)
	s.http.Handler = router
	s.http.RegisterOnShutdown(func() { cancel() })
}

func (s *Server) shutdown() error {
	log.Printf("-- shutting down")
	err := s.http.Shutdown(context.Background())
	if err != nil {
		return err
	}

	log.Printf("-- save: %s", s.opt.Dumpfile)
	err = SaveStoreToFile(s.store, s.opt.Dumpfile)
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
