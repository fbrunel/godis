package godis

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type CommandHandler struct {
	service *CommandService
}

func NewCommandHandler(srv *CommandService) *CommandHandler {
	return &CommandHandler{
		service: srv,
	}
}

func (h *CommandHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	log.Printf("<< conn: %s", conn.RemoteAddr())

	errch := make(chan error, 1)
	go func() {
		for {
			var c Command
			err := conn.ReadJSON(&c)
			if err != nil {
				errch <- err
				break
			}
			log.Printf("<- recv: %v", c)

			start := time.Now()
			rep, _ := h.service.ExecCommand(c)
			delta := time.Since(start)

			err = conn.WriteJSON(*rep)
			if err != nil {
				errch <- err
				break
			}
			log.Printf("-> sent: %v (%v)", *rep, delta)
		}
		log.Printf("Go routine ended")
	}()

	select {
	case <-r.Context().Done():
		hangup(conn)
		time.Sleep(250 * time.Millisecond)
	case err = <-errch:
		log.Printf("EE (%s) %v", conn.RemoteAddr(), err)
	}
}

func hangup(conn *websocket.Conn) error {
	return conn.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}
