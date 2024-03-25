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

	"github.com/google/shlex"
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

func replyReader(ws *websocket.Conn, replych chan<- godis.Reply, errch chan<- error) {
	for {
		var r godis.Reply
		err := ws.ReadJSON(&r)
		if err != nil {
			errch <- err
			break
		}
		log.Printf("<- recv: %v", r)
		replych <- r
	}
}

func commandWriter(ws *websocket.Conn, cmdch <-chan godis.Command, errch chan<- error) {
	for {
		cmd := <-cmdch
		err := ws.WriteJSON(cmd)
		if err != nil {
			errch <- err
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
	case godis.TypeAck:
		return fmt.Sprintf("%v", r.Value)
	case godis.TypeNil:
		return "(nil)"
	case godis.TypeStr:
		return fmt.Sprintf("\"%s\"", r.Value)
	case godis.TypeArr:
		var lines []string
		for i, v := range r.Values() {
			lines = append(lines, fmt.Sprintf("%d \"%s\"", i+1, v.(string)))
		}
		if len(lines) == 0 {
			return "(empty)"
		}
		return strings.Join(lines, "\n")
	}
	return fmt.Sprintf("%s %v", r.Type, r.Value)
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

	var (
		errch   = make(chan error)
		replych = make(chan godis.Reply)
		cmdch   = make(chan godis.Command)
		donech  = make(chan struct{})
	)

	go replyReader(ws, replych, errch)
	go commandWriter(ws, cmdch, errch)
	go func() {
		for {
			prompt := readPrompt("> ")
			if prompt == "exit" {
				close(donech)
				return
			}
			tokens, _ := shlex.Split(prompt)
			cmdch <- godis.MakeCommand(tokens[0], tokens[1:]...)
			reply := <-replych
			fmt.Println(fmtReply(&reply))
		}
	}()

	select {
	case <-donech:
		break
	case err := <-errch:
		fmt.Printf("EE %v", err)
		os.Exit(1)
	}

	hangup(ws)
}
