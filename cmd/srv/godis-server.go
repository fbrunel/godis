package main

import (
	"flag"
	godis "godis/internal"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type CommandHandler struct {
	Runner   *godis.CommandRunner
	Upgrader websocket.Upgrader
}

func NewCommandHandler(runner *godis.CommandRunner) *CommandHandler {
	return &CommandHandler{
		Runner:   runner,
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

		r := h.Runner.RunCommand(c)

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

	r := godis.NewCommandRunner()
	http.Handle("/cmd", NewCommandHandler(r))
	log.Printf("-- serv: %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
