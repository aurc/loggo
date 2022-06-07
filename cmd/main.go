package main

import (
	"github.com/aurc/loggo/internal/loggo"
	"github.com/aurc/loggo/internal/reader"
)

func main() {
	inputChan := make(chan string, 1)
	reader := reader.MakeReader("testdata/log3.txt")

	reader.StreamInto(inputChan)
	app := loggo.NewLoggoApp(inputChan, "")
	app.Run()
}
