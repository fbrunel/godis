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
	case godis.TypeAck, godis.TypeNil, godis.TypeStr:
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
		fmt.Printf("EE %v", err)
		os.Exit(1)
	}

	cherr := make(chan error, 1)
	chreply := make(chan godis.Reply)
	go replyReader(ws, chreply, cherr)

	chcmd := make(chan godis.Command)
	go commandWriter(ws, chcmd, cherr)

	chdone := make(chan struct{})
	go func() {
		for {
			prompt := readPrompt("> ")
			if prompt == "exit" {
				close(chdone)
				return
			}
			tokens := strings.Split(prompt, " ")
			chcmd <- godis.MakeCommand(tokens[0], tokens[1:]...)

			r := <-chreply
			fmt.Println(fmtReply(&r))
		}
	}()

	select {
	case <-chdone:
		break
	case err := <-cherr:
		fmt.Printf("EE %v", err)
		os.Exit(1)
	}

	hangup(ws)
}
