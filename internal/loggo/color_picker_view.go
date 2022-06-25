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
	"sort"

	"github.com/aurc/loggo/internal/color"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ColorPickerView struct {
	tview.Flex
	app                      Loggo
	contextMenu              *tview.List
	onSelect                 func(color string)
	toggleFullScreenCallback func()
	closeCallback            func()
	table                    *tview.Table
	data                     *ColorPickerData
	colors                   [][]string
	title                    string
	colorToCell              map[string][]int
}

func NewColorPickerView(app Loggo, title string, onSelect func(string),
	toggleFullScreenCallback, closeCallback func()) *ColorPickerView {
	tv := &ColorPickerView{
		Flex:                     *tview.NewFlex(),
		app:                      app,
		onSelect:                 onSelect,
		toggleFullScreenCallback: toggleFullScreenCallback,
		closeCallback:            closeCallback,
		title:                    title,
	}
	tv.makeColorTable()
	tv.makeUIComponents()
	tv.makeLayouts()
	return tv
}

func (t *ColorPickerView) SelectColor(color string) {
	if rc, ok := t.colorToCell[color]; ok {
		t.table.Select(rc[0], rc[1])
	}
}

func (t *ColorPickerView) makeColorTable() {
	const columns = 5
	t.colorToCell = make(map[string][]int)
	col := 0
	row := 0
	var currRow []string
	t.colors = [][]string{}
	var sortedCols []string
	for c := range tcell.ColorNames {
		sortedCols = append(sortedCols, c)
	}
	sort.Strings(sortedCols)
	for _, c := range sortedCols {
		if col < columns {
			currRow = append(currRow, c)
			t.colorToCell[c] = []int{row, col}
			col++
			if col == columns {
				t.colors = append(t.colors, currRow)
				currRow = []string{}
				col = 0
				row++
			}
		}
	}
	if col > 0 && col < columns {
		t.colors = append(t.colors, currRow)
	}
}

func (t *ColorPickerView) makeUIComponents() {
	t.data = &ColorPickerData{
		colourPickerView: t,
	}
	t.contextMenu = tview.NewList()
	t.contextMenu.
		SetBorder(true).
		SetTitle("Context Menu").
		SetBackgroundColor(color.ColorBackgroundField)

	t.table = tview.NewTable().
		SetSelectable(true, true).
		SetSeparator(tview.Borders.Vertical).
		SetContent(t.data)

	t.table.SetSelectionChangedFunc(func(row, column int) {
		t.makeContextMenu()
	})

	t.table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if (event.Key() == tcell.KeyEnter ||
			event.Rune() == 's' ||
			event.Rune() == 'S') && t.onSelect != nil {
			r, c := t.table.GetSelection()
			col := t.colors[r][c]
			t.onSelect(col)
			return nil
		}
		switch event.Rune() {
		case 'x', 'X':
			if t.closeCallback != nil {
				t.closeCallback()
			}
		case 'f', 'F':
			if t.toggleFullScreenCallback != nil {
				t.toggleFullScreenCallback()
			}
		}
		return event
	})
}

func (t *ColorPickerView) makeLayouts() {
	t.makeContextMenu()
	t.Flex.Clear().SetDirection(tview.FlexColumn).
		AddItem(t.contextMenu, 30, 1, false).
		AddItem(t.table, 0, 2, true).
		SetBackgroundColor(color.ColorBackgroundField).
		SetBorder(true).
		SetTitle(t.title)
}

func (t *ColorPickerView) makeContextMenu() {
	t.contextMenu.Clear().ShowSecondaryText(false).SetBorderPadding(0, 0, 1, 1)
	t.contextMenu.
		ShowSecondaryText(false)
	if t.toggleFullScreenCallback != nil {
		t.contextMenu.AddItem("Toggle Full Screen", "", 'f', func() {
			t.toggleFullScreenCallback()
		})
	}
	if t.onSelect != nil {
		t.contextMenu.AddItem("/ [yellow::](ENTER)[-::-] Select Color", "", 's', func() {
			r, c := t.table.GetSelection()
			col := t.colors[r][c]
			t.onSelect(col)
		})
	}
	if t.closeCallback != nil {
		t.contextMenu.AddItem("Close", "", 'x', func() {
			t.closeCallback()
		})
	}
}

type ColorPickerData struct {
	tview.TableContentReadOnly
	colourPickerView *ColorPickerView
}

func (d *ColorPickerData) GetCell(row, column int) *tview.TableCell {
	if column+1 <= len(d.colourPickerView.colors[row]) {
		c := d.colourPickerView.colors[row][column]
		label := fmt.Sprintf(` [%s] ■ [-] %s `, c, c)
		return tview.NewTableCell(label).
			SetAlign(tview.AlignLeft).
			SetBackgroundColor(tcell.ColorBlack)
	}
	return nil
}

func (d *ColorPickerData) GetRowCount() int {
	if d.colourPickerView.colors == nil {
		return 0
	}
	return len(d.colourPickerView.colors)
}

func (d *ColorPickerData) GetColumnCount() int {
	if d.colourPickerView.colors == nil {
		return 0
	}
	return len(d.colourPickerView.colors[0])
}
