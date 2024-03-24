package main

import (
	"bufio"
	"flag"
	"fmt"
	godis "godis/internal"
	"io"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

func dial(addr string, path string) (*websocket.Conn, error) {
	u := url.URL{Scheme: "ws", Host: addr, Path: path}
	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}
	return ws, nil
}

func hangup(ws *websocket.Conn) error {
	return ws.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}

func replyReader(ws *websocket.Conn, chreply chan<- godis.Reply, cherr chan<- error) {
	for {
		var r godis.Reply
		err := ws.ReadJSON(&r)
		if err != nil {
			cherr <- err
			break
		}
		log.Printf("<- recv: %v", r)
		chreply <- r
	}
}

func commandWriter(ws *websocket.Conn, chcmd <-chan godis.Command, cherr chan<- error) {
	for {
		cmd := <-chcmd
		err := ws.WriteJSON(cmd)
		if err != nil {
			cherr <- err
			break
		}
		log.Printf("-> sent: %v", cmd)
	}
}

func readPrompt(prefix string) string {
	var str string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stdout, prefix)
		str, _ = r.ReadString('\n')
		str = strings.TrimSpace(str)
		if str != "" {
			break
		}
	}
	return str
}

func fmtReply(r *godis.Reply) string {
	switch r.Type {
	case godis.TypeAck, godis.TypeStr:
		return fmt.Sprintf("%v", r.Data)
	}
	return fmt.Sprintf("%s %v", r.Type, r.Data)
}

func main() {
	addr := flag.String("addr", ":8080", "server address:port")
	verb := flag.Bool("v", false, "verbose")
	flag.Parse()
	log.SetFlags(0)

	if !*verb {
		log.SetOutput(io.Discard)
	}

	//

	ws, err := dial(*addr, "/cmd")
	if err != nil {
		log.Fatal(err)
	}

	cherr := make(chan error, 1)
	chreply := make(chan godis.Reply)
	go replyReader(ws, chreply, cherr)

	chcmd := make(chan godis.Command)
	go commandWriter(ws, chcmd, cherr)

	for {
		prompt := readPrompt("> ")
		if prompt == "exit" {
			break
		}
		tokens := strings.Split(prompt, " ")
		chcmd <- godis.MakeCommand(tokens[0], tokens[1:]...)

		select {
		case r := <-chreply:
			fmt.Println(fmtReply(&r))
		case err := <-cherr:
			fmt.Printf("EE %v", err)
			return
		}
	}

	hangup(ws)
}
