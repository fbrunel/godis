package main

import (
	"flag"
	godis "godis/internal"
	"log"
)

func main() {
	addr := flag.String("addr", ":8080", "")
	flag.Parse()
	log.SetFlags(0)

	srv := godis.NewCommandService()
	api := godis.NewAPIServer(srv)
	log.Fatal(api.Serve(*addr))
}
