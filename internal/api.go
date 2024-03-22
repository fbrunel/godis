package internal

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type APIServer struct {
	service  *CommandService
	upgrader websocket.Upgrader
}

func NewAPIServer(srv *CommandService) *APIServer {
	return &APIServer{
		service:  NewCommandService(),
		upgrader: websocket.Upgrader{},
	}
}

func (api *APIServer) handleCommand(w http.ResponseWriter, r *http.Request) {
	ws, err := api.upgrader.Upgrade(w, r, nil)
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

		r := api.service.ExecCommand(c)

		err = ws.WriteJSON(r)
		if err != nil {
			break
		}
		log.Printf("-> sent: %v", r)
	}
}

func (api *APIServer) Serve(addr string) error {
	http.HandleFunc("/cmd", api.handleCommand)
	log.Printf("-- serv: %s", addr)
	return http.ListenAndServe(addr, nil)
}
