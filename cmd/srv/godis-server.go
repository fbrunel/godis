package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"

	"github.com/fbrunel/godis/godis"
)

func run(addr string, storefn string) error {
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM)
	defer stop()

	store, err := godis.LoadStoreFromFile(storefn)
	if err != nil {
		store = godis.NewStandardStore()
	}

	server := godis.NewServer(store)
	go func() {
		err := server.Start(ctx, addr)
		if err != nil {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	err = server.Shutdown(context.Background())
	if err != nil {
		return err
	}

	err = godis.SaveStoreToFile(store, storefn)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	addr := flag.String("addr", ":8080", "listening address in the form of host:port")
	storefn := flag.String("store", "/tmp/godis.dump", "pathname of the data store file")

	flag.Parse()
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.LUTC)

	err := run(*addr, *storefn)
	if err != nil {
		log.Printf("EE %v", err)
	}
}
