package godis

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type CommandHandler struct {
	wg      sync.WaitGroup
	service *CommandService
}

func NewCommandHandler(srv *CommandService) *CommandHandler {
	return &CommandHandler{
		service: srv,
	}
}

func (h *CommandHandler) WaitClose() {
	h.wg.Wait()
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
	h.wg.Add(1)
	go func() {
		defer h.wg.Done()
		var err error
		for {
			var c Command
			// ReadJSON blocks undefinitely waiting for data to be read
			// but will exit when conn.Close() is called, when ServeHTTP()
			// terminates.
			err = conn.ReadJSON(&c)
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
		log.Printf(">> conn: %s (%v)", conn.RemoteAddr(), err)
	}()

	select {
	case <-r.Context().Done():
		hangup(conn)
		time.Sleep(250 * time.Millisecond)
	case <-errch:
		// Do nothing
	}
}

func hangup(conn *websocket.Conn) error {
	return conn.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}
