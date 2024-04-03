package main

import (
	"flag"
	"log"

	"github.com/fbrunel/godis/godis"
)

func main() {
	addr := flag.String("addr", ":8080", "server address:port")
	flag.Parse()
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.LUTC)

	server := godis.NewServer(*addr)
	err := server.Start()
	if err != nil {
		log.Fatalf("EE %v", err)
	}
}
