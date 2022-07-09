/*
Copyright © 2022 Aurelio Calegari, et al.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package loggo

import (
	"fmt"
	"github.com/aurc/loggo/internal/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"time"
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
	cfg, err := config.MakeConfig(configFile)
	if err != nil {
		panic(err)
	}
	return NewAppWithConfig(cfg)
}

func NewAppWithConfig(cfg *config.Config) *appScaffold {
	scaffold := &appScaffold{}
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

func (a *appScaffold) ShowPopMessage(text string, waitSecs int64) {
	modal := tview.NewFlex().SetDirection(tview.FlexRow)
	modal.SetBackgroundColor(tcell.ColorDarkBlue)
	countdownText := tview.NewTextView().SetTextAlign(tview.AlignRight)
	mainContent := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetWordWrap(true).
		SetDynamicColors(true).
		SetText(text)
	mainContent.SetBackgroundColor(tcell.ColorDarkBlue).SetBorderPadding(0, 0, 2, 2)
	modal.AddItem(mainContent, 0, 1, false)
	modal.AddItem(countdownText, 1, 1, false)
	a.ShowModal(modal, len(text)/2, 5, tcell.ColorDarkBlue)
	countdownText.SetTextColor(tcell.ColorLightGrey).SetBackgroundColor(tcell.ColorDarkBlue)
	go func() {
		for i := waitSecs; i >= 0; i-- {
			countdownText.SetText(fmt.Sprintf(`(%ds)`, i))
			a.Draw()
			time.Sleep(time.Second)
		}
		a.DismissModal()
		a.Draw()
	}()
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
	a.ShowModal(modal, width, height, tcell.ColorDarkBlue)
}

func (a *appScaffold) ShowModal(p tview.Primitive, width, height int, bgColor tcell.Color) {
	modContainer := tview.NewFlex().AddItem(p, 0, 1, false)
	modContainer.SetBorder(true).SetBackgroundColor(bgColor)
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
