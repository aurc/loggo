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
	"github.com/rivo/tview"
)

type FilterView struct {
	tview.Flex
	app          Loggo
	nameField    *tview.InputField
	filterTree   *tview.TreeView
	keysTable    *tview.Table
	keyNameField *tview.InputField
	operation    *tview.DropDown
	value1       *tview.InputField
	value2       *tview.InputField
	setupPane    *tview.Flex
	showQuit     bool
}

func NewFilterView(app Loggo, showQuit bool) *FilterView {
	tv := &FilterView{
		Flex:     *tview.NewFlex(),
		app:      app,
		showQuit: showQuit,
	}
	tv.makeUIComponents()
	tv.makeLayouts()
	tv.makeKeysTableData()
	tv.app.SetFocus(tv.keysTable)
	return tv
}

func (t *FilterView) makeUIComponents() {
	t.nameField = tview.NewInputField()
	t.filterTree = tview.NewTreeView()
	t.keysTable = tview.NewTable().SetSelectable(true, false)
	t.keyNameField = tview.NewInputField()
	t.operation = tview.NewDropDown()
	t.value1 = tview.NewInputField()
	t.value2 = tview.NewInputField()
}

func (t *FilterView) makeLayouts() {
	t.Flex.Clear()

	body := tview.NewFlex().SetDirection(tview.FlexRow)
	body.SetBorder(true)
	t.Flex.Clear().SetDirection(tview.FlexColumn).
		AddItem(body, 0, 1, false).
		AddItem(t.keysTable, 40, 1, true)

	t.setupPane = tview.NewFlex().SetDirection(tview.FlexRow)
	body.Clear().
		AddItem(t.filterTree, 0, 1, false).
		AddItem(t.setupPane, 5, 1, false)
	t.makeSetupPane()
}

func (t *FilterView) makeSetupPane() {
	t.setupPane.Clear()
}

func (t *FilterView) makeKeysTableData() {
	t.keysTable.SetCell(0, 0, tview.NewTableCell("[yellow::b]Key Name").SetSelectable(false))
	for i, v := range t.app.Config().Keys {
		t.keysTable.SetCell(i+1, 0, tview.NewTableCell(v.Name).SetSelectable(true))
	}
	t.keysTable.SetSelectionChangedFunc(func(row, column int) {
		if row > 0 {

		}
	})
}
