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

	"github.com/aurc/loggo/internal/color"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (l *LogView) populateMenu() {
	l.navMenu = tview.NewFlex().SetDirection(tview.FlexRow)
	l.navMenu.
		SetBackgroundColor(color.ColorBackgroundField).SetBorderPadding(0, 0, 1, 0)
	sepForeground := tview.Styles.ContrastBackgroundColor
	sepStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(sepForeground)
	l.navMenu.
		//////////////////////////////////////////////////////////////////
		// Stream Menu
		//////////////////////////////////////////////////////////////////
		AddItem(NewHorizontalSeparator(sepStyle, LineHThick, "Stream", sepForeground), 1, 2, false).
		AddItem(l.followingView, 1, 2, false).
		AddItem(l.textViewMenuControl(tview.NewTextView().
			SetDynamicColors(true).SetRegions(true).
			SetText(`[yellow::b] ^t      [-::u]["1"]Template[""]`), func() {
			if l.isTemplateViewShown() {
				// TODO: Find a reliable way to respond to external closure
			} else {
				l.makeLayoutsWithTemplateView()
			}
		}), 1, 2, false).
		AddItem(l.textViewMenuControl(tview.NewTextView().
			SetDynamicColors(true).SetRegions(true).
			SetText(`[yellow::b] :       [-::u]["1"]Local Filter[""]`), func() {
			if l.isTemplateViewShown() {
				// TODO: Find a reliable way to respond to external closure
			} else {
				l.makeLayoutsWithTemplateView()
			}
		}), 1, 2, false).
		//////////////////////////////////////////////////////////////////
		// Navigation Menu
		//////////////////////////////////////////////////////////////////
		AddItem(
			NewHorizontalSeparator(sepStyle, LineHThick, "Navigation", sepForeground), 1, 2, false).
		AddItem(tview.NewTextView().
			SetDynamicColors(true).
			SetText("[yellow::b] ↓ ↑ ← →[-::-] Navigate"), 1, 3, false).
		AddItem(l.textViewMenuControl(tview.NewTextView().
			SetDynamicColors(true).SetRegions(true).
			SetText(`[yellow::b] g       [-::u]["1"]Top[""]`), func() {
			l.isFollowing = false
			l.table.ScrollToBeginning()
			if len(l.inSlice) > 1 {
				go l.table.Select(1, 0)
			}
		}), 1, 1, false).
		AddItem(l.textViewMenuControl(tview.NewTextView().
			SetDynamicColors(true).SetRegions(true).
			SetText(`[yellow::b] G       [-::u]["1"]Bottom[""]`), func() {
			l.isFollowing = false
			l.table.ScrollToEnd()
			go l.table.Select(len(l.inSlice), 0)
		}), 1, 2, false).
		AddItem(l.textViewMenuControl(tview.NewTextView().
			SetDynamicColors(true).SetRegions(true).
			SetText(`[yellow::b] ^b      [-::u]["1"]Pg Up[""]`), func() {
			l.isFollowing = false
			l.table.InputHandler()(tcell.NewEventKey(tcell.KeyPgUp, '0', 0), func(p tview.Primitive) {})
		}), 1, 2, false).
		AddItem(l.textViewMenuControl(tview.NewTextView().
			SetDynamicColors(true).SetRegions(true).
			SetText(`[yellow::b] ^f      [-::u]["1"]Pg Down[""]`), func() {
			l.isFollowing = false
			l.table.InputHandler()(tcell.NewEventKey(tcell.KeyPgDn, '0', 0), func(p tview.Primitive) {})
		}), 1, 2, false).
		//////////////////////////////////////////////////////////////////
		// Selection Menu
		//////////////////////////////////////////////////////////////////
		AddItem(NewHorizontalSeparator(sepStyle, LineHThick, "Selection", sepForeground), 1, 2, false).
		AddItem(tview.NewTextView().
			SetDynamicColors(true).
			SetText("[yellow::b] ⌥ 🖱    [-::-] Horizontal"), 1, 3, false).
		AddItem(tview.NewTextView().
			SetDynamicColors(true).
			SetText("[yellow::b] ⌥ ⌘ 🖱  [-::-] Vertical"), 1, 3, false).
		//////////////////////////////////////////////////////////////////
		// Application Menu
		//////////////////////////////////////////////////////////////////
		AddItem(NewHorizontalSeparator(sepStyle, LineHThick, "Application", sepForeground), 1, 2, false).
		AddItem(l.textViewMenuControl(tview.NewTextView().
			SetDynamicColors(true).SetRegions(true).
			SetText(`[yellow::b] ^a      [-::u]["1"]About[""]`), func() {
			go func() {
				l.showAbout()
			}()
		}), 1, 2, false).
		AddItem(l.textViewMenuControl(tview.NewTextView().SetRegions(true).
			SetDynamicColors(true).
			SetText(`[yellow::b] ^c      [-::u]["1"]Quit[""]`), func() {
			l.app.Stop()
		}), 0, 1, false).
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(l.linesView, 1, 1, false)

	l.mainMenu = tview.NewFlex().SetDirection(tview.FlexColumn)
	l.mainMenu.
		SetBackgroundColor(color.ColorBackgroundField).SetTitleAlign(tview.AlignCenter)
	l.mainMenu.
		AddItem(l.textViewMenuControl(tview.NewTextView().
			SetDynamicColors(true).SetRegions(true).
			SetText(`[yellow::b](^t) [-::u]["1"]Template[""]`), func() {
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
		if len(removed) == 0 {
			onFocus()
		}
	})
	//onFocus()
	return tv
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
		l.followingView.SetText(`[yellow::b] ^Space  [-::u]["1"]Stream[::-] [green::bi]ON[-::-][""]`)
	} else {
		l.followingView.SetText(`[yellow::b] ^Space  [-::u]["1"]Stream[::-] [red::bi]OFF[-::-][""]`)
	}
}
