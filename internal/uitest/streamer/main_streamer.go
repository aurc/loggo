package main

import (
	"fmt"

	"github.com/aurc/loggo/internal/reader"
)

func main() {
	streamReader := reader.MakeReader("")
	streamReceiver := make(chan string, 1)
	go streamReader.StreamInto(streamReceiver)
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
