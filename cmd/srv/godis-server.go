package main

import (
	"flag"
	"log"

	"github.com/fbrunel/godis/godis"
)

func main() {
	addr := flag.String("addr", ":8080", "server address:port")
	dump := flag.String("dump", "/tmp/godis.dump", "pathname of the dump file")

	flag.Parse()
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.LUTC)

	options := godis.DefaultOptions()
	options.Addr = *addr
	options.Dumpfile = *dump

	server := godis.NewServer(options)
	err := server.Start()
	if err != nil {
		log.Fatalf("EE %v", err)
	}
}
