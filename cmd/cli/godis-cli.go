package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/fbrunel/godis/godis"

	"github.com/google/shlex"
)

func ReadPrompt(prefix string) string {
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

func FmtReply(r *godis.Reply) string {
	switch r.Type {
	case godis.TypeAck:
		return fmt.Sprintf("%s", r.Value)
	case godis.TypeNil:
		return "(nil)"
	case godis.TypeStr:
		return fmt.Sprintf("\"%s\"", r.Value)
	case godis.TypeArr:
		var lines []string
		for i, v := range r.Values() {
			var vfmt string
			switch v.(type) {
			case nil:
				vfmt = "(nil)"
			case string:
				vfmt = fmt.Sprintf("\"%s\"", v)
			default:
				vfmt = fmt.Sprintf("%v", v)
			}
			lines = append(lines, fmt.Sprintf("%d "+vfmt, i+1))
		}
		if len(lines) == 0 {
			return "(empty)"
		}
		return strings.Join(lines, "\n")
	}
	return fmt.Sprintf("%s %v", r.Type, r.Value)
}

//

func main() {
	addr := flag.String("addr", ":8080", "server address:port")
	verb := flag.Bool("v", false, "verbose")
	flag.Parse()
	log.SetFlags(0)

	if !*verb {
		log.SetOutput(io.Discard)
	}

	//

	client := godis.NewClient(*addr)

	err := client.Dial()
	if err != nil {
		fmt.Println("EE", err)
		os.Exit(1)
	}

	donech := make(chan struct{})
	go func() {
		for {
			prompt := ReadPrompt("> ")
			if prompt == "exit" {
				close(donech)
				return
			}
			tokens, _ := shlex.Split(prompt)
			reply, err := client.SendCommand(tokens[0], tokens[1:]...)
			if err != nil {
				fmt.Println("EE", err)
				close(donech)
				return
			}
			fmt.Println(FmtReply(reply))
		}
	}()

	<-donech
	client.Hangup()
}
