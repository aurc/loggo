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
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (l *LogView) keyEvents() {
	l.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if l.app.inputCapture != nil {
			return l.app.inputCapture(event)
		}
		switch event.Key() {
		case tcell.KeyCtrlN:
			l.toggleSelectionMouse()
			return nil
		case tcell.KeyCtrlA:
			go func() {
				l.showAbout()
			}()
			return nil
		case tcell.KeyCtrlT:
			l.makeLayoutsWithTemplateView()
			return nil
		case tcell.KeyCtrlSpace:
			l.toggledFollowing()
			return nil
		case tcell.KeyTAB:
			if l.isJsonViewShown() {
				if l.jsonView.textView.HasFocus() {
					l.app.SetFocus(l.table)
					go func() {
						time.Sleep(time.Millisecond)
						l.updateBottomBarMenu()
					}()
				} else {
					l.app.SetFocus(l.jsonView.textView)
					go func() {
						time.Sleep(time.Millisecond)
						l.updateBottomBarMenu()
					}()
				}
				return nil
			}
			return event
		}
		prim := l.app.app.GetFocus()
		if _, ok := prim.(*tview.InputField); ok {
			return event
		}
		switch event.Rune() {
		case ':':
			l.toggleFilter()
			return nil
		}
		if prim == l.table && l.isJsonViewShown() {
			switch event.Rune() {
			case 'f', '`', 's', 'r', 'g', 'G', 'w', 'x':
				return l.jsonView.textView.GetInputCapture()(event)
			}
		}

		return event
	})
}
