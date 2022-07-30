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
	"regexp"
	"strings"

	"github.com/aurc/loggo/internal/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type LogData struct {
	tview.TableContentReadOnly
	logView *LogView
}

func (d *LogData) GetCell(row, column int) *tview.TableCell {
	d.logView.filterLock.RLock()
	defer d.logView.filterLock.RUnlock()
	if row == -1 || len(d.logView.finSlice) < row-1 || column == -1 {
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
			if _, ok := d.logView.finSlice[row-1][config.ParseErr]; ok {
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
	cellValue := k.ExtractValue(d.logView.finSlice[row-1])
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

	if k.Name == config.TextPayload {
		if _, ok := d.logView.finSlice[row-1][config.ParseErr]; ok {
			fgColor = tcell.ColorBlue
		}
	}

	return tc.
		SetBackgroundColor(bgColor).
		SetTextColor(fgColor).
		SetText(fmt.Sprintf("%s", cellValue))
}

func (d *LogData) GetRowCount() int {
	d.logView.filterLock.RLock()
	defer d.logView.filterLock.RUnlock()
	return len(d.logView.finSlice) + 1
}

func (d *LogData) GetColumnCount() int {
	d.logView.filterLock.RLock()
	defer d.logView.filterLock.RUnlock()
	c := d.logView.config
	return len(c.Keys) + 1
}
