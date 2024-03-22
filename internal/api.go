package internal

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type CommandHandler struct {
	service  *CommandService
	upgrader websocket.Upgrader
}

func NewCommandHandler(srv *CommandService) *CommandHandler {
	return &CommandHandler{
		service:  srv,
		upgrader: websocket.Upgrader{},
	}
}

func (h *CommandHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	for {
		var c Command
		err := ws.ReadJSON(&c)
		if err != nil {
			break
		}
		log.Printf("<- recv: %v", c)

		r := h.service.ExecCommand(c)

		err = ws.WriteJSON(r)
		if err != nil {
			break
		}
		log.Printf("-> sent: %v", r)
	}
}
