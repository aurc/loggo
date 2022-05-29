package main

import (
	"github.com/aurc/loggo/pkg/loggo"
	"github.com/rivo/tview"
	"os"
)

//var (
//	app        *tview.Application
//	layout     *tview.Flex
//	jsonViewer *ui.JsonViewer
//)

func main() {
	app := tview.NewApplication()

	jsonViewer := loggo.NewJsonView(app)

	b, err := os.ReadFile("testdata/test1.json")
	if err != nil {
		panic(err)
	}
	jsonViewer.SetJson(b)
	if err := app.
		SetRoot(jsonViewer, true).
		EnableMouse(true).
		Run(); err != nil {
		panic(err)
	}
}

//func main() {
//	app = tview.NewApplication()
//
//	jsonViewer = ui.MakeJsonViewer(app)
//	layout = tview.NewFlex().
//		SetDirection(tview.FlexRow).
//		AddItem(jsonViewer, 0, 1, true)
//
//	b, err := os.ReadFile("testdata/test1.json")
//	if err != nil {
//		panic(err)
//	}
//
//	setKeyboardShortcuts()
//
//	jsonViewer.SetJson(b)
//
//	if err := app.
//		SetRoot(layout, true).
//		SetFocus(jsonViewer).
//		EnableMouse(true).
//		Run(); err != nil {
//		panic(err)
//	}
//}
//
//func setKeyboardShortcuts() *tview.Application {
//	return app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
//		if e := jsonViewer.HandleShortcuts(event); e == nil {
//			return nil
//		}
//		return event
//	})
//}
