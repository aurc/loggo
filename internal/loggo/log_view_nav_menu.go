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
	"runtime"

	"github.com/aurc/loggo/internal/color"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	selectionMouseEnabledMenu  = `[yellow::b] ^n      [-::u]["1"]Enable Selection[""]`
	selectionMouseDisabledMenu = `[yellow::b] ^n      [-::u]["1"]Enable Mouse[""]`
	templateMenu               = `[yellow::b] ^t      [-::u]["1"]Template[""]`
	localFilterMenu            = `[yellow::b] :       [-::u]["1"]Local Filter[""]`
	viewEntryMenu              = `[yellow::b] Enter[-::-]   View Entry`
	navigateMenu               = `[yellow::b] ↓ ← ↑ →[-::-] Navigate`
	goTopMenu                  = `[yellow::b] g       [-::u]["1"]Top[""]`
	goBottomMenu               = `[yellow::b] G       [-::u]["1"]Bottom[""]`
	pageUpMenu                 = `[yellow::b] ^b      [-::u]["1"]Pg Up[""]`
	pageDownMenu               = `[yellow::b] ^f      [-::u]["1"]Pg Down[""]`
	mouseHoMenu                = `[yellow::b] ⌥ 🖱    [-::-]Horizontal`
	mouseVeMenu                = `[yellow::b] ⌥ ⌘ 🖱  [-::-]Vertical`
	aboutMenu                  = `[yellow::b] ^a      [-::u]["1"]About[""]`
	quitMenu                   = `[yellow::b] ^c      [-::u]["1"]Quit[""]`
	autoScrollOnMenu           = `[yellow::b] ^Space  [-::u]["1"]Auto-Scroll[::-] [green::bi]ON[-::-][""]`
	autoScrollOffMenu          = `[yellow::b] ^Space  [-::u]["1"]Auto-Scroll[::-] [red::bi]OFF[-::-][""]`
)

func (l *LogView) populateMenu() {
	l.mouseSel = tview.NewTextView().
		SetDynamicColors(true).SetRegions(true).
		SetText(selectionMouseEnabledMenu)

	l.navMenu = tview.NewFlex().SetDirection(tview.FlexRow)
	l.navMenu.
		SetBackgroundColor(color.ColorBackgroundField).SetBorderPadding(0, 0, 0, 0)
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
			SetText(templateMenu), func() {
			if l.isTemplateViewShown() {
				// TODO: Find a reliable way to respond to external closure
			} else {
				l.makeLayoutsWithTemplateView()
			}
		}), 1, 2, false).
		AddItem(l.textViewMenuControl(tview.NewTextView().
			SetDynamicColors(true).SetRegions(true).
			SetText(localFilterMenu), func() {
			l.toggleFilter()
		}), 1, 2, false).
		//////////////////////////////////////////////////////////////////
		// Navigation Menu
		//////////////////////////////////////////////////////////////////
		AddItem(
			NewHorizontalSeparator(sepStyle, LineHThick, "Navigation", sepForeground), 1, 2, false).
		AddItem(tview.NewTextView().
			SetDynamicColors(true).
			SetText(viewEntryMenu), 1, 3, false).
		AddItem(tview.NewTextView().
			SetDynamicColors(true).
			SetText(navigateMenu), 1, 3, false).
		AddItem(l.textViewMenuControl(tview.NewTextView().
			SetDynamicColors(true).SetRegions(true).
			SetText(goTopMenu), func() {
			l.isFollowing = false
			l.table.ScrollToBeginning()
			if len(l.inSlice) > 1 {
				go l.table.Select(1, 0)
			}
		}), 1, 1, false).
		AddItem(l.textViewMenuControl(tview.NewTextView().
			SetDynamicColors(true).SetRegions(true).
			SetText(goBottomMenu), func() {
			l.isFollowing = false
			l.table.ScrollToEnd()
			go l.table.Select(len(l.inSlice), 0)
		}), 1, 2, false).
		AddItem(l.textViewMenuControl(tview.NewTextView().
			SetDynamicColors(true).SetRegions(true).
			SetText(pageUpMenu), func() {
			l.isFollowing = false
			l.table.InputHandler()(tcell.NewEventKey(tcell.KeyPgUp, '0', 0), func(p tview.Primitive) {})
		}), 1, 2, false).
		AddItem(l.textViewMenuControl(tview.NewTextView().
			SetDynamicColors(true).SetRegions(true).
			SetText(pageDownMenu), func() {
			l.isFollowing = false
			l.table.InputHandler()(tcell.NewEventKey(tcell.KeyPgDn, '0', 0), func(p tview.Primitive) {})
		}), 1, 2, false)
	//////////////////////////////////////////////////////////////////
	// Selection Menu
	//////////////////////////////////////////////////////////////////
	l.navMenu.
		AddItem(NewHorizontalSeparator(sepStyle, LineHThick, "Selection", sepForeground), 1, 2, false).
		AddItem(l.textViewMenuControl(l.mouseSel, l.toggleSelectionMouse), 1, 2, false)
	if runtime.GOOS != "windows" {
		l.navMenu.
			AddItem(tview.NewTextView().
				SetDynamicColors(true).
				SetText(mouseHoMenu), 1, 3, false).
			AddItem(tview.NewTextView().
				SetDynamicColors(true).
				SetText(mouseVeMenu), 1, 3, false)
	}
	//////////////////////////////////////////////////////////////////
	// Application Menu
	//////////////////////////////////////////////////////////////////
	l.navMenu.
		AddItem(NewHorizontalSeparator(sepStyle, LineHThick, "Application", sepForeground), 1, 2, false).
		AddItem(l.textViewMenuControl(tview.NewTextView().
			SetDynamicColors(true).SetRegions(true).
			SetText(aboutMenu), func() {
			go func() {
				l.showAbout()
			}()
		}), 1, 2, false).
		AddItem(l.textViewMenuControl(tview.NewTextView().SetRegions(true).
			SetDynamicColors(true).
			SetText(quitMenu), func() {
			l.app.Stop()
		}), 1, 1, false).
		AddItem(NewHorizontalSeparator(sepStyle, LineHThick, "", sepForeground), 1, 2, false).
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(l.linesView, 1, 1, false)

	l.mainMenu = tview.NewFlex().SetDirection(tview.FlexColumn)
	l.updateBottomBarMenu()
}

func (l *LogView) updateBottomBarMenu() {
	l.mainMenu.Clear().
		SetBackgroundColor(color.ColorBackgroundField).SetTitleAlign(tview.AlignCenter)
	l.mainMenu.
		AddItem(l.textViewMenuControl(tview.NewTextView().
			SetDynamicColors(true).SetRegions(true).
			SetText(`[yellow::b](^t) [-::u]["1"]Template[""]`), func() {
			if l.isTemplateViewShown() {
				// TODO: Find a reliable way to respond to external closure
			} else {
				l.makeLayoutsWithTemplateView()
				l.updateBottomBarMenu()
			}
		}), 0, 3, false).
		AddItem(l.followingView, 0, 5, false)
	if l.isJsonViewShown() && !l.jsonView.HasFocus() {
		l.mainMenu.
			AddItem(l.textViewMenuControl(tview.NewTextView().SetRegions(true).
				SetDynamicColors(true).
				SetText(`[yellow::b](TAB) [-::u]["1"]Focus Log Entry[""]`), func() {
				go l.app.SetFocus(l.jsonView.textView)
			}), 0, 3, false)
	} else if l.isJsonViewShown() && l.jsonView.HasFocus() {
		l.mainMenu.
			AddItem(l.textViewMenuControl(tview.NewTextView().SetRegions(true).
				SetDynamicColors(true).
				SetText(`[yellow::b](TAB) [-::u]["1"]Focus Stream Table[""]`), func() {
				go l.app.SetFocus(l.table)
			}), 0, 3, false)
	}
	l.mainMenu.
		AddItem(l.textViewMenuControl(tview.NewTextView().SetRegions(true).
			SetDynamicColors(true).
			SetText(`[yellow::b](^c) [-::u]["1"]Quit[""]`), func() {
			l.app.Stop()
		}), 0, 2, false).
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
		l.followingView.SetText(autoScrollOnMenu)
	} else {
		l.followingView.SetText(autoScrollOffMenu)
	}
}

func (l *LogView) toggleSelectionMouse() {
	l.selectionEnabled = !l.selectionEnabled
	l.app.app.EnableMouse(!l.selectionEnabled)
	go func() {
		if l.selectionEnabled {
			l.app.ShowPopMessage("Mouse disabled! Click and drag to select...", 2, l.table)
			l.mouseSel.SetText(selectionMouseDisabledMenu)
		} else {
			l.app.ShowPopMessage("Selection disabled! Mouse input active...", 2, l.table)
			l.mouseSel.SetText(selectionMouseEnabledMenu)
		}
		l.app.Draw()
	}()
}
