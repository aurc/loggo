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
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const LineHThick = '\u2501'
const LineHThin = '\u2500'

func NewHorizontalSeparator(lineStyle tcell.Style, lineRune rune, text string, textColor tcell.Color) *tview.Box {
	return tview.NewBox().
		SetDrawFunc(func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
			// Draw a horizontal line across the middle of the box.
			centerY := y + height/2
			for cx := x; cx < x+width; cx++ {

				screen.SetContent(cx, centerY, lineRune, nil, lineStyle)
			}

			// Write some text along the horizontal line.
			if len(text) > 0 {
				tview.Print(screen, text, x+1, centerY, width-2, tview.AlignCenter, textColor)
			}

			// Space for other content.
			return x + 1, centerY + 1, width - 2, height - (centerY + 1 - y)
		})
}
