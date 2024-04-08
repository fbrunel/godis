package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"

	"github.com/fbrunel/godis/godis"
)

func run(options godis.Options) error {
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM)
	defer stop()

	return godis.NewServer(options).Start(ctx)
}

func main() {
	addr := flag.String("addr", ":8080", "server address:port")
	dump := flag.String("dump", "/tmp/godis.dump", "pathname of the dump file")

	flag.Parse()
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.LUTC)

	options := godis.DefaultOptions()
	options.Addr = *addr
	options.Dumpfile = *dump

	err := run(options)
	if err != nil {
		log.Printf("EE %v", err)
	}
}
