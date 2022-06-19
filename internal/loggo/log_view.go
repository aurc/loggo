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
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/aurc/loggo/internal/color"
	"github.com/aurc/loggo/internal/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type LogView struct {
	tview.Flex
	app                *LoggoApp
	input              <-chan string
	table              *tview.Table
	jsonView           *JsonView
	data               *LogData
	templateView       *TemplateView
	layout             *tview.Flex
	config             *config.Config
	navMenu            *tview.Flex
	mainMenu           *tview.Flex
	linesView          *tview.TextView
	followingView      *tview.TextView
	logFullScreen      bool
	templateFullScreen bool
	inSlice            []map[string]interface{}
	globalCount        int64
	isFollowing        bool
}

func NewLogReader(app *LoggoApp, input <-chan string) *LogView {
	lv := &LogView{
		Flex:        *tview.NewFlex(),
		app:         app,
		config:      app.Config(),
		input:       input,
		isFollowing: true,
	}
	lv.makeUIComponents()
	lv.makeLayouts()
	lv.read()
	go func() {
		lv.app.ShowModal(NewSplashScreen(lv.app), 71, 16, tcell.ColorBlack)
		lv.app.Draw()
		time.Sleep(4 * time.Second)
		lv.app.DismissModal()
		lv.app.Draw()
	}()
	go func() {
		time.Sleep(10 * time.Millisecond)
		lv.isFollowing = true
	}()
	return lv
}

func (l *LogView) read() {
	go func() {
		var sampling []map[string]interface{}
		samplingCount := 0
		if len(l.config.LastSavedName) == 0 {
			samplingCount = 50
		}
		lastUpdate := time.Now().Add(-time.Minute)
		for {
			t := <-l.input
			if len(t) > 0 {
				l.globalCount++
				m := make(map[string]interface{})
				err := json.Unmarshal([]byte(t), &m)
				if err != nil {
					m[config.ParseErr] = err.Error()
					m[config.TextPayload] = t
				}
				if l.globalCount <= int64(samplingCount) {
					sampling = append(sampling, m)
					l.processSampleForConfig(sampling)
				}
				// TODO: Review at some stage if sampling gets to accumulate
				//} else if len(sampling) <= samplingCount {
				// l.processSampleForConfig(sampling)
				//}
				l.inSlice = append(l.inSlice, m)
				l.updateLineView()
				now := time.Now()
				if now.Sub(lastUpdate)*time.Millisecond > 500 && l.isFollowing {
					lastUpdate = now
					l.app.Draw()
					l.table.ScrollToEnd()
				}
			}
		}
	}()
}

func (l *LogView) processSampleForConfig(sampling []map[string]interface{}) {
	if len(l.config.LastSavedName) > 0 || l.isTemplateViewShown() {
		return
	}
	l.config = config.MakeConfigFromSample(sampling)
	l.app.config = l.config
}

func (l *LogView) textViewMenuControl(tv *tview.TextView, onFocus func()) *tview.TextView {
	tv.SetBlurFunc(func() {
		tv.Highlight("")
	}).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			onFocus()
			return nil
		case tcell.KeyESC:
			l.app.SetFocus(l.table)
			return nil
		}
		return event
	})
	tv.SetHighlightedFunc(func(added, removed, remaining []string) {
		onFocus()
	})
	onFocus()
	return tv
}

func (l *LogView) makeUIComponents() {
	l.templateView = NewTemplateView(l.app, false, func() {
		// Toggle full screen func
		l.templateFullScreen = !l.templateFullScreen
		l.makeLayoutsWithTemplateView()
	}, l.makeLayouts)
	l.templateView.SetBorder(true).SetTitle("Template Editor")
	l.data = &LogData{
		logView: l,
	}
	selection := func(row, column int) {
		if row > 0 {
			l.jsonView = NewJsonView(l.app, false,
				func() {
					// Toggle full screen func
					l.logFullScreen = !l.logFullScreen
					l.makeLayoutsWithJsonView()
				}, l.makeLayouts)
			l.jsonView.SetBorder(true).SetTitle("Log Entry")
			b, _ := json.Marshal(l.inSlice[row-1])
			l.jsonView.SetJson(b)
			l.makeLayoutsWithJsonView()
		} else {
			l.makeLayouts()
		}
	}
	l.table = tview.NewTable().
		SetSelectable(true, false).
		SetFixed(1, 1).
		SetSeparator(tview.Borders.Vertical).
		SetContent(l.data)
	l.table.SetSelectedFunc(selection).
		SetBackgroundColor(color.ColorBackgroundField)
	l.table.SetSelectionChangedFunc(func(row, column int) {
		// stop scrolling!
		if l.isFollowing {
			l.isFollowing = false

			go func() {
				r, c := l.table.GetOffset()
				l.updateLineView()
				l.table.SetOffset(r, c)
				l.table.Select(r, c)
				go l.app.Draw()
			}()
		} else {
			r, c := l.table.GetOffset()
			l.updateLineView()
			l.table.SetOffset(r, c)
		}
	})

	l.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyF1:
			l.makeLayoutsWithTemplateView()
			return nil
		case tcell.KeyF2:
			l.toggledFollowing()
			return nil
		}
		return event
	})

	l.linesView = tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignRight)
	l.followingView = tview.NewTextView().
		SetRegions(true).
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	l.followingView.SetFocusFunc(func() {
		go l.toggledFollowing()
	})
	l.followingView.SetBlurFunc(func() {
		l.followingView.Highlight("")
	})
	l.navMenu = tview.NewFlex().SetDirection(tview.FlexColumn)
	l.navMenu.
		SetBackgroundColor(color.ColorBackgroundField).SetTitleAlign(tview.AlignCenter)

	l.navMenu.
		AddItem(tview.NewTextView().
			SetDynamicColors(true).
			SetText("[yellow::b](↲)[-::-] View"), 0, 2, false).
		AddItem(tview.NewTextView().
			SetDynamicColors(true).
			SetText("[yellow::b](↓ ↑ ← →)[-::-] Navigate"), 0, 3, false).
		AddItem(l.textViewMenuControl(tview.NewTextView().
			SetDynamicColors(true).SetRegions(true).
			SetText(`[yellow::b](g) [-::u]["1"]Top[""]`), func() {
			l.isFollowing = false
			l.table.ScrollToBeginning()
			if len(l.inSlice) > 1 {
				go l.table.Select(1, 0)
			}
		}), 0, 1, false).
		AddItem(l.textViewMenuControl(tview.NewTextView().
			SetDynamicColors(true).SetRegions(true).
			SetText(`[yellow::b](G) [-::u]["1"]Bottom[""]`), func() {
			l.isFollowing = false
			l.table.ScrollToEnd()
			go l.table.Select(len(l.inSlice), 0)
		}), 0, 2, false).
		AddItem(l.textViewMenuControl(tview.NewTextView().
			SetDynamicColors(true).SetRegions(true).
			SetText(`[yellow::b](^f) [-::u]["1"]Pg Up[""]`), func() {
			l.isFollowing = false
			l.table.InputHandler()(tcell.NewEventKey(tcell.KeyPgUp, '0', 0), func(p tview.Primitive) {})
		}), 0, 2, false).
		AddItem(l.textViewMenuControl(tview.NewTextView().
			SetDynamicColors(true).SetRegions(true).
			SetText(`[yellow::b](^b) [-::-]["1"]Pg Down[""]`), func() {
			l.isFollowing = false
			l.table.InputHandler()(tcell.NewEventKey(tcell.KeyPgDn, '0', 0), func(p tview.Primitive) {})
		}), 0, 2, false)
	l.mainMenu = tview.NewFlex().SetDirection(tview.FlexColumn)
	l.mainMenu.
		SetBackgroundColor(color.ColorBackgroundField).SetTitleAlign(tview.AlignCenter)
	l.mainMenu.
		AddItem(l.textViewMenuControl(tview.NewTextView().
			SetDynamicColors(true).SetRegions(true).
			SetText(`[yellow::b](F1) [-::u]["1"]Template[""]`), func() {
			if l.isTemplateViewShown() {
				// TODO: Find a reliable way to respond to external closure
			} else {
				l.makeLayoutsWithTemplateView()
			}
		}), 0, 2, false).
		AddItem(l.followingView, 0, 5, false).
		AddItem(l.textViewMenuControl(tview.NewTextView().SetRegions(true).
			SetDynamicColors(true).
			SetText(`[yellow::b](^C) [-::u]["1"]Quit[""]`), func() {
			l.app.Stop()
		}), 0, 1, false).
		AddItem(l.linesView, 0, 3, false)
	l.updateLineView()
}

func (l *LogView) isTemplateViewShown() bool {
	return l.Flex.GetItemCount() > 0 && l.Flex.GetItem(0) == l.templateView ||
		l.Flex.GetItemCount() > 1 && l.Flex.GetItem(1) == l.templateView
}

func (l *LogView) toggledFollowing() {
	l.isFollowing = !l.isFollowing
	l.updateLineView()
	go l.app.Draw()
}
func (l *LogView) updateLineView() {
	r, _ := l.table.GetSelection()
	if r > 0 {
		l.linesView.SetText(
			fmt.
				Sprintf(`[yellow::]Line [green::b]%d[yellow::-] ([green::b]%d[yellow::-] lines)`,
					r,
					l.globalCount))
	} else {
		l.linesView.SetText(
			fmt.
				Sprintf(`[green::b]%d[yellow::-] lines`,
					l.globalCount))
	}
	if l.isFollowing {
		l.followingView.SetText(`[yellow::b](F2) [-::u]["1"]Toggle Auto-Scroll[""][::-] ([green::bi]ON[-::-])`)
	} else {
		l.followingView.SetText(`[yellow::b](F2) [-::u]["1"]Toggle Auto-Scroll[""][::-] ([red::bi]OFF[-::-])`)
	}
}

func (l *LogView) makeLayouts() {
	l.Flex.Clear().SetDirection(tview.FlexRow).
		AddItem(l.navMenu, 1, 1, false).
		AddItem(l.table, 0, 2, true).
		AddItem(l.mainMenu, 1, 1, false).
		SetBackgroundColor(color.ColorBackgroundField)
	l.app.SetFocus(l.table)
}

func (l *LogView) makeLayoutsWithJsonView() {
	l.Flex.Clear().SetDirection(tview.FlexRow)
	if !l.logFullScreen {
		l.Flex.AddItem(l.table, 0, 1, false)
	}
	l.Flex.
		AddItem(l.jsonView, 0, 2, false).
		AddItem(l.mainMenu, 1, 1, false)

	l.app.SetFocus(l.jsonView.textView)
}

func (l *LogView) makeLayoutsWithTemplateView() {
	l.isFollowing = false
	l.Flex.Clear().SetDirection(tview.FlexRow)
	if !l.templateFullScreen {
		l.Flex.AddItem(l.table, 0, 1, false)
	}
	l.templateView.config = l.config
	l.Flex.
		AddItem(l.templateView, 0, 2, false).
		AddItem(l.mainMenu, 1, 1, false)

	l.app.SetFocus(l.templateView.table)
}

type LogData struct {
	tview.TableContentReadOnly
	logView *LogView
}

func (d *LogData) GetCell(row, column int) *tview.TableCell {
	if row == -1 || len(d.logView.inSlice) < row-1 || column == -1 {
		return nil
	}
	if column == 0 {
		if row == 0 {
			tc := tview.NewTableCell("[yellow] ☀[white] / [blue]☂ ").
				SetAlign(tview.AlignCenter).
				SetBackgroundColor(tcell.ColorBlack).
				SetSelectable(false)
			return tc
		} else {
			if _, ok := d.logView.inSlice[row-1][config.ParseErr]; ok {
				tc := tview.NewTableCell(" ︎  ☂   ").
					SetTextColor(tcell.ColorBlue).
					SetAlign(tview.AlignCenter).
					SetBackgroundColor(tcell.ColorBlack)
				return tc
			} else {
				tc := tview.NewTableCell("   ︎☀  ︎ ").
					SetTextColor(tcell.ColorYellow).
					SetAlign(tview.AlignCenter).
					SetBackgroundColor(tcell.ColorBlack)
				return tc
			}
		}
	}
	c := d.logView.config
	if len(c.Keys) == 0 {
		return nil
	}
	k := c.Keys[column-1]
	tc := tview.NewTableCell(" " + k.Name + " ")
	if k.MaxWidth > 0 && k.MaxWidth-len(k.Name) >= len(k.Name) {
		spaces := strings.Repeat(" ", k.MaxWidth-len(k.Name))
		tc.SetText(" " + k.Name + spaces)
	}
	// Set Headers
	if row == 0 {
		tc.SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignCenter).
			SetBackgroundColor(tcell.ColorBlack).
			SetSelectable(false)
		return tc
	}
	// Set Body Cells
	cellValue := k.ExtractValue(d.logView.inSlice[row-1])
	var bgColor, fgColor tcell.Color
	if len(k.Color.Foreground) == 0 {
		fgColor = k.Type.GetColor()
	} else {
		fgColor = k.Color.GetForegroundColor()
	}
	bgColor = k.Color.GetBackgroundColor()
	if len(k.ColorWhen) > 0 {
	OUT:
		for _, kv := range k.ColorWhen {
			reg, err := regexp.Compile(kv.MatchValue)
			if err == nil && reg.FindIndex([]byte(cellValue)) != nil {
				bgColor = kv.Color.GetBackgroundColor()
				fgColor = kv.Color.GetForegroundColor()
				break OUT
			}
		}
	}
	switch k.Type {
	case config.TypeNumber, config.TypeBool:
		tc.SetAlign(tview.AlignRight)
	}
	if k.MaxWidth > 0 {
		tc.MaxWidth = k.MaxWidth
	}

	return tc.
		SetBackgroundColor(bgColor).
		SetTextColor(fgColor).
		SetText(fmt.Sprintf("%s", cellValue))
}

func (d *LogData) GetRowCount() int {
	return len(d.logView.inSlice) + 1
}

func (d *LogData) GetColumnCount() int {
	c := d.logView.config
	return len(c.Keys) + 1
}
