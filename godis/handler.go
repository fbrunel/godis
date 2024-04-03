package godis

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type CommandHandler struct {
	ctx     context.Context
	service *CommandService
}

func NewCommandHandler(ctx context.Context, srv *CommandService) *CommandHandler {
	return &CommandHandler{
		ctx:     ctx,
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

	errch := make(chan error)
	go func() {
		for {
			var c Command
			err := conn.ReadJSON(&c)
			if err != nil {
				log.Printf("EE (%s) %v", conn.RemoteAddr(), err)
				errch <- err
				break
			}
			log.Printf("<- recv: %v", c)

			start := time.Now()
			rep, _ := h.service.ExecCommand(c)
			delta := time.Since(start)

			err = conn.WriteJSON(*rep)
			if err != nil {
				log.Printf("EE (%s) %v", conn.RemoteAddr(), err)
				errch <- err
				break
			}
			log.Printf("-> sent: %v (%v)", *rep, delta)
		}
	}()

	select {
	case <-h.ctx.Done():
		hangup(conn)
		<-time.After(500 * time.Millisecond)
	case <-errch:
		log.Printf("DONE")
	}
}

func hangup(conn *websocket.Conn) error {
	return conn.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}
