/*
Copyright © 2022 Aurelio Calegari

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
	"github.com/aurc/loggo/internal/config"
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
	StackView(p tview.Primitive)
	PopView()
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
