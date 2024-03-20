package main

import (
	"flag"
	"godis/internal"
	"log"
	"net/url"

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

	ws := Dial(addr)

	done := make(chan struct{})
	go func() {
		for {
			_, bytes, err := ws.ReadMessage()
			if err != nil {
				close(done)
				return
			}
			r, _ := internal.DecodeReply(bytes)
			log.Printf("<- recv: %s", r)
		}
	}()

	c := internal.MakeCommand(args...)
	bytes, _ := internal.EncodeCommand(c)
	err := ws.WriteMessage(websocket.TextMessage, bytes)
	if err != nil {
		log.Println("EE werr:", err)
		return
	}
	log.Printf("-> send: %s", c)

	Hangup(ws)
	<-done
}
