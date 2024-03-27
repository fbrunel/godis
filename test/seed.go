package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/fbrunel/godis/godis"
)

func ReadFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	log.Println("read lines", len(lines))
	return lines, nil
}

func main() {
	addr := flag.String("addr", ":8080", "server address:port")
	verb := flag.Bool("v", false, "verbose")
	filename := flag.String("words", "/usr/share/dict/words", "words file")
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

	words, _ := ReadFile(*filename)
	for i, w := range words {
		_, err := client.SendCommand("SET", w, fmt.Sprint(len(w)))
		fmt.Printf("%d\t%s\033[0K\r", i, w)
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	client.Hangup()
}
