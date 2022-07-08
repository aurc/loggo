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
	"github.com/aurc/loggo/internal/filter"
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
	filterGroup  filter.Group
	showQuit     bool
}

func NewFilterView(app Loggo, filter filter.Group, showQuit bool) *FilterView {
	tv := &FilterView{
		Flex:        *tview.NewFlex(),
		app:         app,
		showQuit:    showQuit,
		filterGroup: filter,
	}
	tv.makeUIComponents()
	tv.makeLayouts()
	return tv
}

func (t *FilterView) makeUIComponents() {

}

func (t *FilterView) makeLayouts() {

}
