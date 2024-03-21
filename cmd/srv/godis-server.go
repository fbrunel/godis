package main

import (
	"flag"
	"godis/internal"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type CommandHandler struct {
	Upgrader *websocket.Upgrader
	Backend  *internal.Backend
}

func NewCommandHandler(backend *internal.Backend) *CommandHandler {
	return &CommandHandler{
		&websocket.Upgrader{},
		backend,
	}
}

func (h *CommandHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := h.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	for {
		var c internal.Command
		err := ws.ReadJSON(&c)
		if err != nil {
			break
		}
		log.Printf("<- recv: %s", c)

		r := h.Backend.EvalCommand(c)

		err = ws.WriteJSON(r)
		if err != nil {
			break
		}
		log.Printf("-> sent: %s", r)
	}

	ws.Close()
}

func main() {
	addr := flag.String("addr", "127.0.0.1:8080", "")
	flag.Parse()
	log.SetFlags(0)

	b := internal.NewDefaultBackend()
	http.Handle("/cmd", NewCommandHandler(b))
	log.Printf("-- serv: %s", *addr)
	http.ListenAndServe(*addr, nil)
}
