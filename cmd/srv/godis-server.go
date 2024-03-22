package main

import (
	"flag"
	godis "godis/internal"
	"log"
	"net/http"
)

func main() {
	addr := flag.String("addr", ":8080", "")
	flag.Parse()
	log.SetFlags(0)

	srv := godis.NewCommandService()
	hdl := godis.NewCommandHandler(srv)

	http.Handle("/cmd", hdl)
	log.Printf("-- serv: %s", *addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
