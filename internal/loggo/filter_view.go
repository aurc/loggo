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
	"strings"
)

type FilterView struct {
	tview.Flex
	app             Loggo
	expressionField *tview.InputField
	buttonSearch    *tview.Button
	buttonClear     *tview.Button
	keyFinderField  *tview.InputField
	showQuit        bool
}

func NewFilterView(app Loggo, showQuit bool) *FilterView {
	tv := &FilterView{
		Flex:     *tview.NewFlex(),
		app:      app,
		showQuit: showQuit,
	}
	tv.makeUIComponents()
	tv.makeLayouts()
	return tv
}

func (t *FilterView) makeUIComponents() {
	t.expressionField = tview.NewInputField().
		SetPlaceholder("Filter Expression...").
		SetFieldStyle(color.FieldStyle).SetPlaceholderStyle(color.PlaceholderStyle)
	t.expressionField.
		SetBackgroundColor(color.ColorBackgroundField)

	t.buttonSearch = tview.NewButton("Search").SetSelectedFunc(func() {

	})
	t.buttonClear = tview.NewButton("Clear").SetSelectedFunc(func() {
		t.expressionField.SetText("")
		t.app.SetFocus(t.expressionField)
	})

	t.keyFinderField = tview.NewInputField().SetPlaceholder("Start typing to find a key...")
	t.keyFinderField.SetAutocompleteFunc(func(currentText string) (entries []string) {
		matches := make([]string, 0)
		for _, v := range t.app.Config().Keys {
			vt := strings.ToLower(strings.TrimSpace(v.Name))
			ct := strings.ToLower(strings.TrimSpace(currentText))
			if strings.Contains(vt, ct) && len(ct) > 0 || ct == "*" {
				matches = append(matches, v.Name)
			}
		}
		return matches
	})

	t.keyFinderField.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter, tcell.KeyTAB:
			t.addKey()
		case tcell.KeyEsc:
			t.keyFinderField.SetText("")
		}
	})
}

func (t *FilterView) addKey() {
	tex := t.expressionField.GetText()
	t.expressionField.SetText(tex + " " + t.keyFinderField.GetText())
	t.keyFinderField.SetText("")
	t.app.SetFocus(t.expressionField)
}

func (t *FilterView) makeLayouts() {
	t.Flex.Clear()
	filterRow := tview.NewFlex().SetDirection(tview.FlexColumn)
	filterField := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(tview.NewTextView().SetText("🔎").SetTextAlign(tview.AlignCenter), 4, 1, true).
		AddItem(t.expressionField, 0, 1, true)
	filterField.SetBorder(true)
	filterRow.
		AddItem(filterField, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewBox(), 1, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
				AddItem(tview.NewBox(), 1, 1, false).
				AddItem(t.buttonSearch, 10, 1, false).
				AddItem(tview.NewBox(), 1, 1, false).
				AddItem(t.buttonClear, 10, 1, false), 1, 1, false).
			AddItem(tview.NewBox(), 1, 1, false),
			23, 1, true)

	okButton := tview.NewButton("OK").SetSelectedFunc(t.addKey)
	okButton.SetBackgroundColor(tcell.ColorGreen)
	actionBar := tview.NewFlex().SetDirection(tview.FlexColumn)
	actionBar.AddItem(tview.NewTextView().SetText(" 🔑 Finder:"), 12, 0, false)
	actionBar.AddItem(t.keyFinderField, 0, 1, false).
		AddItem(tview.NewBox(), 1, 1, false).
		AddItem(okButton, 4, 1, false).
		AddItem(tview.NewTextView().SetText(" |"), 2, 0, false)
	t.addButton(actionBar, "=")
	t.addButton(actionBar, "==")
	t.addButton(actionBar, "!=")
	t.addButton(actionBar, ">")
	t.addButton(actionBar, "<")
	t.addButton(actionBar, ">=")
	t.addButton(actionBar, "<=")
	actionBar.AddItem(tview.NewTextView().SetText(" |"), 2, 0, false)
	t.addButton(actionBar, "CONTAINS")
	t.addButton(actionBar, "BETWEEN")
	t.addButton(actionBar, "MATCH")
	actionBar.AddItem(tview.NewTextView().SetText(" |"), 2, 0, false)
	t.addButton(actionBar, "AND")
	t.addButton(actionBar, "OR")
	actionBar.AddItem(tview.NewBox(), 24, 1, false)

	t.Flex.Clear().SetDirection(tview.FlexRow).
		AddItem(filterRow, 3, 1, false).
		AddItem(actionBar, 1, 1, false)

}

func (t *FilterView) addButton(ab *tview.Flex, title string) {
	b := tview.NewButton(title).SetSelectedFunc(func() {
		t.expressionField.SetText(fmt.Sprintf(`%s %s `, t.expressionField.GetText(), title))
		t.app.SetFocus(t.expressionField)
	})
	b.SetBackgroundColor(tcell.ColorGray).SetTitleColor(tcell.ColorWhite)
	ab.
		AddItem(tview.NewBox(), 1, 1, false).
		AddItem(b, len(title)+2, 1, false)
}
