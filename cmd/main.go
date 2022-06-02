package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/aurc/loggo/internal/reader"
	"github.com/aurc/loggo/pkg/loggo"
	"github.com/google/uuid"
)

//func main() {
//	inputChan := make(chan string, 1)
//	go func() {
//
//		for {
//			line := <-inputChan
//			if len(line) > 0 {
//				fmt.Println("----------------")
//				fmt.Println(line)
//			}
//		}
//	}()
//	go func() {
//		time.Sleep(time.Second * 10)
//		close(inputChan)
//	}()
//	reader.ReadPipeStream(inputChan)
//}

//func main() {
//	app := tview.NewApplication()
//
//	jsonViewer := loggo.NewJsonView(app, false)
//
//	b, err := os.ReadFile("testdata/test3.txt")
//	//b, err := os.ReadFile("testdata/test1.json")
//	if err != nil {
//		panic(err)
//	}
//	jsonViewer.SetJson(b)
//	if err := app.
//		SetRoot(jsonViewer, true).
//		EnableMouse(true).
//		Run(); err != nil {
//		panic(err)
//	}
//}

//func main() {
//	app := tview.NewApplication()
//	c, err := config.MakeConfig("")
//	if err != nil {
//		panic(err)
//	}
//	jsonViewer := loggo.NewTemplateView(app, c)
//
//	if err := app.
//		SetRoot(jsonViewer, true).
//		EnableMouse(true).
//		Run(); err != nil {
//		panic(err)
//	}
//}

func main() {
	//app := tview.NewApplication()
	inputChan := make(chan string, 1)
	reader := reader.MakeReader("")
	oldStdIn := os.Stdin
	defer func() {
		os.Stdin = oldStdIn
	}()
	r, w, _ := os.Pipe()
	os.Stdin = r
	go func() {
		JsonGenerator(w)
	}()

	reader.StreamInto(inputChan)
	app := loggo.NewLoggoApp(inputChan, "")
	app.Run()
	//c, err := config.MakeConfig("")
	//if err != nil {
	//	panic(err)
	//}
	//jsonViewer := loggo.NewLogReader(app, inputChan, c)
	//
	//pages := tview.NewPages().
	//	AddPage("background", jsonViewer, true, true)
	//
	//if err := app.
	//	SetRoot(pages, true).
	//	EnableMouse(true).
	//	Run(); err != nil {
	//	panic(err)
	//}
}

func JsonGenerator(writer io.Writer) {
	b, err := ioutil.ReadFile("testdata/test1.json")
	if err != nil {
		panic(err)
	}
	jm := make(map[string]interface{})
	_ = json.Unmarshal(b, &jm)

	for {
		uid1 := uuid.New().String()
		uid2 := uuid.New().String()
		id3 := rand.Intn(30)
		jm["severity"] = []string{"INFO", "ERROR", "DEBUG", "WARN"}[rand.Intn(4)]
		jm["insertId"] = uid1
		jm["trace"] = uid2
		jm["spanId"] = fmt.Sprintf("%d", id3)
		jm["timestamp"] = time.Now().Format("2006-01-02T15:04:05-0700")
		b, _ = json.Marshal(jm)
		fmt.Fprintln(writer, string(b))
		time.Sleep(time.Second)
	}
}
