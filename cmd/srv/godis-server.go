package main

import (
	"flag"
	godis "godis/internal"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type CommandHandler struct {
	Backend  *godis.Backend
	Upgrader websocket.Upgrader
}

func NewCommandHandler(backend *godis.Backend) *CommandHandler {
	return &CommandHandler{
		Backend:  backend,
		Upgrader: websocket.Upgrader{},
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
	addr := flag.String("addr", ":8080", "")
	flag.Parse()
	log.SetFlags(0)

	be := godis.NewBackend()
	http.Handle("/cmd", NewCommandHandler(be))
	log.Printf("-- serv: %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
