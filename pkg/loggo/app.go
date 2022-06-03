package loggo

import (
	"github.com/aurc/loggo/pkg/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type LoggoApp struct {
	appScaffold
	input   <-chan string
	logView *LogView
}

type Loggo interface {
	Draw()
	SetInputCapture(cap func(event *tcell.EventKey) *tcell.EventKey)
	Stop()
	SetFocus(primitive tview.Primitive)
	ShowPrefabModal(text string, width, height int, buttons ...*tview.Button)
	ShowModal(p tview.Primitive, width, height int)
	DismissModal()
	Config() *config.Config
}

func NewLoggoApp(input <-chan string, configFile string) *LoggoApp {
	app := NewApp(configFile)
	lapp := &LoggoApp{
		appScaffold: *app,
		input:       input,
	}

	lapp.logView = NewLogReader(lapp, input)

	lapp.pages = tview.NewPages().
		AddPage("background", lapp.logView, true, true)

	return lapp
}

func (a *LoggoApp) Run() {
	if err := a.app.
		SetRoot(a.pages, true).
		EnableMouse(true).
		Run(); err != nil {
		panic(err)
	}
}
