package main

import (
	"flag"
	godis "godis/internal"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type CommandHandler struct {
	Upgrader *websocket.Upgrader
	Backend  *godis.Backend
}

func NewCommandHandler(backend *godis.Backend) *CommandHandler {
	return &CommandHandler{
		Upgrader: &websocket.Upgrader{},
		Backend:  backend,
	}
}

func (h *CommandHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := h.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	for {
		var c godis.Command
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
}

func main() {
	addr := flag.String("addr", "127.0.0.1:8080", "")
	flag.Parse()
	log.SetFlags(0)

	be := godis.NewDefaultBackend()
	http.Handle("/cmd", NewCommandHandler(be))
	log.Printf("-- serv: %s", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}
