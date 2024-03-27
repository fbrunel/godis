package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/fbrunel/godis/godis"
)

func main() {
	addr := flag.String("addr", ":8080", "server address:port")
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
