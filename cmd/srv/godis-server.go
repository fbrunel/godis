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
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.LUTC)

	store := godis.NewStandardStore()
	service := godis.NewCommandService(store)
	handle := godis.NewCommandHandler(service)

	http.Handle("/cmd", handle)
	log.Printf("-- serv: %s", *addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
