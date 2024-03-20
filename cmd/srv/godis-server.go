package main

import (
	"flag"
	"godis/internal"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type CommandHandler struct {
	Upgrader websocket.Upgrader
}

func NewCommandHandler() *CommandHandler {
	return &CommandHandler{
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
		_, bytes, err := ws.ReadMessage()
		if err != nil {
			break
		}

		c, _ := internal.DecodeCommand(bytes)
		log.Printf("<- recv: %s", c)

		r := internal.MakeReply("OK")
		bytes, _ = internal.EncodeReply(r)
		err = ws.WriteMessage(websocket.TextMessage, bytes)
		if err != nil {
			break
		}
	}
}

func main() {
	addr := flag.String("addr", "localhost:8080", "")
	flag.Parse()
	log.SetFlags(0)

	http.Handle("/cmd", NewCommandHandler())
	log.Printf("-- serv: %s", *addr)
	http.ListenAndServe(*addr, nil)
}
