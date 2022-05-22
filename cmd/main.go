package main

import (
	"os"

	"github.com/gdamore/tcell/v2"

	"github.com/aurc/loggo/pkg/ui"
	"github.com/rivo/tview"
)

var (
	app        *tview.Application
	layout     *tview.Flex
	jsonViewer *ui.JsonViewer
)

func main() {
	app = tview.NewApplication()

	//newPrimitive := func(text string) tview.Primitive {
	//	return tview.NewTextView().
	//		SetTextAlign(tview.AlignCenter).
	//		SetText(text)
	//}
	//searchBar := newPrimitive("JSON Viewer")
	//textView := ui.NewJsonRenderer().
	//	SetJsonConfigIndent(ui.OrderSorted, "  ")
	//grid := tview.NewGrid().
	//	SetRows(0, 3, 3).
	//	SetColumns(0).
	//	SetBorders(true).
	//	AddItem(textView, 0, 0, 1, 1, 0, 0, true).
	//	AddItem(searchBar, 1, 0, 1, 1, 3, 0, false).
	//	AddItem(newPrimitive("Footer"), 2, 0, 1, 1, 0, 0, false)

	jsonViewer = ui.MakeJsonViewer(app)
	layout = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(jsonViewer, 0, 1, true)

	b, err := os.ReadFile("testdata/test1.json")
	if err != nil {
		panic(err)
	}

	setKeyboardShortcuts()

	jsonViewer.SetJson(b)
	//textView.SetJson(b).SetChangedFunc(func() {
	//	app.Draw()
	//})
	if err := app.
		SetRoot(layout, true).
		SetFocus(jsonViewer).
		EnableMouse(true).
		Run(); err != nil {
		panic(err)
	}
}

func setKeyboardShortcuts() *tview.Application {
	return app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		//sc := unicode.ToLower(event.Rune())
		//fmt.Printf("%v, %s\n", sc, string(sc))
		// Global shortcuts
		//switch unicode.ToLower(event.Rune()) {
		//case 's':
		//	app.SetFocus(jsonViewer)
		//	return nil
		//case 'q':
		//case 't':
		//
		//	os.Exit(0)
		//	return nil
		//}
		if e := jsonViewer.HandleShortcuts(event); e == nil {
			return nil
		}

		//// Handle based on current focus. Handlers may modify event
		//switch {
		//case projectPane.HasFocus():
		//	event = projectPane.handleShortcuts(event)
		//case taskPane.HasFocus():
		//	event = taskPane.handleShortcuts(event)
		//	if event != nil && projectDetailPane.isShowing() {
		//		event = projectDetailPane.handleShortcuts(event)
		//	}
		//case taskDetailPane.HasFocus():
		//	event = taskDetailPane.handleShortcuts(event)
		//}

		return event
	})
}
