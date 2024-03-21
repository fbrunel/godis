package main

import (
	"flag"
	godis "godis/internal"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

func dial(addr *string) (*websocket.Conn, error) {
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/cmd"}
	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}
	return ws, nil
}

func hangup(ws *websocket.Conn) {
	_ = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}

func main() {
	addr := flag.String("addr", ":8080", "")
	flag.Parse()
	log.SetFlags(0)

	args := flag.Args()
	if len(args) == 0 {
		return
	}

	//

	ostart := time.Now()
	ws, err := dial(addr)
	if err != nil {
		log.Fatal(err)
	}
	defer hangup(ws)

	istart := time.Now()
	done := make(chan struct{})
	go func() {
		for {
			var r godis.Reply
			err := ws.ReadJSON(&r)
			if err != nil {
				close(done)
				return
			}
			log.Printf("<- recv: %s", r)
			log.Printf("-- time: %s", time.Since(istart))
		}
	}()

	c := godis.MakeCommand(args[0], args[1:]...)
	err = ws.WriteJSON(c)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("-> sent: %s", c)

	hangup(ws)
	<-done

	log.Printf("-- time: %s", time.Since(ostart))
}
