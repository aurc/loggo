package loggo

import (
	"fmt"

	"github.com/aurc/loggo/internal/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type appScaffold struct {
	app        *tview.Application
	config     *config.Config
	pages      *tview.Pages
	modal      *tview.Flex
	stackPages []tview.Primitive
}

type App interface {
	Stop()
	Run(p tview.Primitive)
}

func NewApp(configFile string) *appScaffold {
	scaffold := &appScaffold{}
	cfg, err := config.MakeConfig(configFile)
	if err != nil {
		panic(err)
	}
	app := tview.NewApplication()

	scaffold.app = app
	scaffold.config = cfg
	scaffold.stackPages = []tview.Primitive{}
	scaffold.pages = tview.NewPages()

	return scaffold
}

func (a *appScaffold) Config() *config.Config {
	return a.config
}

func (a *appScaffold) Draw() {
	a.app.Draw()
}

func (a *appScaffold) SetInputCapture(cap func(event *tcell.EventKey) *tcell.EventKey) {
	a.app.SetInputCapture(cap)
}

func (a *appScaffold) Stop() {
	a.app.Stop()
}

func (a *appScaffold) SetFocus(primitive tview.Primitive) {
	a.app.SetFocus(primitive)
}

func (a *appScaffold) StackView(p tview.Primitive) {
	a.stackPages = append(a.stackPages, p)
	a.pages.AddPage(fmt.Sprintf(`_%d`, len(a.stackPages)), p, true, true)
}

func (a *appScaffold) PopView() {
	a.pages.RemovePage(fmt.Sprintf(`_%d`, len(a.stackPages)))
	a.stackPages = a.stackPages[:len(a.stackPages)-1]
}

func (a *appScaffold) ShowPrefabModal(text string, width, height int, buttons ...*tview.Button) {
	modal := tview.NewFlex().SetDirection(tview.FlexRow)
	modal.SetBackgroundColor(tcell.ColorDarkBlue)
	mainContent := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetWordWrap(true).
		SetText(text)
	mainContent.SetBackgroundColor(tcell.ColorDarkBlue).SetBorderPadding(1, 0, 2, 2)

	buts := tview.NewFlex().SetDirection(tview.FlexColumn)
	for _, b := range buttons {
		b.SetBackgroundColor(tcell.ColorWhite)
		b.SetLabelColor(tcell.ColorBlack)
		buts.AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDarkBlue), 2, 1, false)
		buts.AddItem(b, 0, 1, false)
	}
	buts.AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDarkBlue), 2, 1, false)

	modal.AddItem(mainContent, 0, 1, false)
	modal.AddItem(buts, 1, 1, false)
	a.ShowModal(modal, width, height)
}

func (a *appScaffold) ShowModal(p tview.Primitive, width, height int) {
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

func (a *appScaffold) DismissModal() {
	a.pages.RemovePage("modal")
}

func (a *appScaffold) Run(p tview.Primitive) {
	a.pages.AddPage("background", p, true, true)
	if err := a.app.
		SetRoot(a.pages, true).
		EnableMouse(true).
		Run(); err != nil {
		panic(err)
	}
}
