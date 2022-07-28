package main

import (
	"fmt"

	"github.com/aurc/loggo/internal/reader"
)

func main() {
	streamReceiver := make(chan string, 1)
	streamReader := reader.MakeReader("", streamReceiver)
	go streamReader.StreamInto()
	for {
		line, ok := <-streamReceiver
		if !ok {
			break
		}
		if len(line) > 0 {
			fmt.Printf("READER: %s", line)
		}
	}
}
