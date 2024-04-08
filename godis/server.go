package godis

import (
	"context"
	"log"
	"net"
	"net/http"
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
	http    http.Server
	opt     Options
	store   *StandardStore
	service *CommandService
	handler *CommandHandler
}

func NewServer(opt Options) *Server {
	return &Server{
		http: http.Server{Addr: opt.Addr},
		opt:  opt,
	}
}

func (srv *Server) Start(ctx context.Context) error {
	srv.setup(ctx)

	errch := make(chan error)
	go func() {
		log.Printf("-- serv: %s", srv.opt.Addr)
		err := srv.http.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			errch <- err
		}
	}()

	select {
	case err := <-errch:
		return err
	case <-ctx.Done():
		return srv.shutdown()
	}
}

//

func (srv *Server) setup(ctx context.Context) {
	log.Printf("-- load: %s", srv.opt.Dumpfile)
	store, err := LoadStoreFromFile(srv.opt.Dumpfile)
	if err != nil {
		log.Printf("EE %v", err)
		store = NewStandardStore()
	}

	srv.store = store
	srv.service = NewCommandService(srv.store)
	srv.handler = NewCommandHandler(srv.service)

	router := http.NewServeMux()
	router.Handle(srv.opt.URLPath, srv.handler)
	srv.http.Handler = router
	srv.http.BaseContext = func(_ net.Listener) context.Context { return ctx }
}

func (srv *Server) shutdown() error {
	log.Printf("-- shutting down")
	err := srv.http.Shutdown(context.Background())
	if err != nil {
		return err
	}

	log.Printf("-- save: %s", srv.opt.Dumpfile)
	err = SaveStoreToFile(srv.store, srv.opt.Dumpfile)
	if err != nil {
		return err
	}

	time.Sleep(1 * time.Second)
	return nil
}
