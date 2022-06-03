package main

import (
	"os"

	"github.com/aurc/loggo/internal/reader"
	"github.com/aurc/loggo/internal/uitest/helper"
	"github.com/aurc/loggo/pkg/loggo"
)

func main() {
	inputChan := make(chan string, 1)
	rd := reader.MakeReader("")
	oldStdIn := os.Stdin
	defer func() {
		os.Stdin = oldStdIn
	}()
	r, w, _ := os.Pipe()
	os.Stdin = r
	go func() {
		helper.JsonGenerator(w)
	}()

	_ = rd.StreamInto(inputChan)
	app := loggo.NewLoggoApp(inputChan, "")
	app.Run()
}
