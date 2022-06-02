package loggo

import (
	"github.com/aurc/loggo/pkg/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type LoggoApp struct {
	app     *tview.Application
	input   <-chan string
	logView *LogView
	config  *config.Config
	pages   *tview.Pages
	modal   *tview.Flex
}

func NewLoggoApp(input <-chan string, configFile string) *LoggoApp {
	cfg, err := config.MakeConfig(configFile)
	if err != nil {
		panic(err)
	}
	app := tview.NewApplication()
	lapp := &LoggoApp{
		app:    app,
		input:  input,
		config: cfg,
	}

	lapp.logView = NewLogReader(lapp, input, cfg)

	lapp.pages = tview.NewPages().
		AddPage("background", lapp.logView, true, true)

	return lapp
}

func (a *LoggoApp) Draw() {
	a.app.Draw()
}

func (a *LoggoApp) SetInputCapture(cap func(event *tcell.EventKey) *tcell.EventKey) {
	a.app.SetInputCapture(cap)
}

func (a *LoggoApp) Stop() {
	a.app.Stop()
}

func (a *LoggoApp) SetFocus(primitive tview.Primitive) {
	a.app.SetFocus(primitive)
}

func (a *LoggoApp) ShowPrefabModal(text string, width, height int, buttons ...*tview.Button) {
	modal := tview.NewFlex().SetDirection(tview.FlexRow)
	modal.SetBackgroundColor(tcell.ColorDarkBlue)
	mainContent := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetWordWrap(true).
		SetText(text)
	mainContent.SetBackgroundColor(tcell.ColorDarkBlue).SetBorderPadding(1, 0, 2, 2)

	buts := tview.NewFlex().SetDirection(tview.FlexColumn)
	for _, b := range buttons {
		buts.AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDarkBlue), 2, 1, false)
		buts.AddItem(b, 0, 1, false)
	}
	buts.AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDarkBlue), 2, 1, false)

	modal.AddItem(mainContent, 0, 1, false)
	modal.AddItem(buts, 1, 1, false)
	a.ShowModal(modal, width, height)
}

func (a *LoggoApp) ShowModal(p tview.Primitive, width, height int) {
	modContainer := tview.NewFlex().AddItem(p, 0, 1, false)
	modContainer.SetBorder(true).SetBackgroundColor(tcell.ColorDarkBlue)
	a.modal = tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(modContainer, height, 1, false).
			AddItem(nil, 0, 1, false), width, 1, false).
		AddItem(nil, 0, 1, false)
	a.pages.AddPage("modal", a.modal, true, true)
}

func (a *LoggoApp) DismissModal() {
	a.pages.RemovePage("modal")
}

func (a *LoggoApp) Run() {
	if err := a.app.
		SetRoot(a.pages, true).
		EnableMouse(true).
		Run(); err != nil {
		panic(err)
	}
}
