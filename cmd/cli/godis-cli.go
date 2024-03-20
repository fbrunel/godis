package main

import (
	"flag"
	"godis/internal"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

func Dial(addr *string) *websocket.Conn {
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/cmd"}
	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}
	return ws
}

func Hangup(ws *websocket.Conn) {
	ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}

func main() {
	addr := flag.String("addr", "localhost:8080", "")
	flag.Parse()
	log.SetFlags(0)

	args := flag.Args()
	if len(args) == 0 {
		return
	}

	//

	ostart := time.Now()
	ws := Dial(addr)

	istart := time.Now()
	done := make(chan struct{})
	go func() {
		for {
			var r internal.Reply
			err := ws.ReadJSON(&r)
			if err != nil {
				close(done)
				return
			}
			log.Printf("<- recv: %s", r)
			log.Printf("-- time: %s", time.Since(istart))
		}
	}()

	c := internal.MakeCommand(args...)
	err := ws.WriteJSON(c)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("-> send: %s", c)

	Hangup(ws)
	<-done

	log.Printf("-- time: %s", time.Since(ostart))
}
