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

type ColorPickerButton struct {
	tview.Flex
	app            Loggo
	input          *tview.InputField
	button         *tview.Button
	label          string
	fieldWidth     int
	labelWidth     int
	labelColor     tcell.Color
	bgColor        tcell.Color
	fieldTextColor tcell.Color
	fieldBgColor   tcell.Color
	labelText      *tview.TextView
	colorLabel     *tview.TextView
	changedFunc    func(text string)
}

func NewColorPickerButton(app Loggo, label, value string, fieldWidth int, changedFunc func(string)) *ColorPickerButton {
	c := &ColorPickerButton{
		Flex:        *tview.NewFlex().SetDirection(tview.FlexColumn),
		app:         app,
		input:       tview.NewInputField().SetText(value),
		button:      tview.NewButton("Choose"),
		label:       label,
		labelText:   tview.NewTextView().SetDynamicColors(true).SetText(label),
		colorLabel:  tview.NewTextView().SetText(" ■ "),
		changedFunc: changedFunc,
		fieldWidth:  fieldWidth,
	}
	c.makeLayout()
	c.button.
		SetSelectedFunc(func() {
			cp := NewColorPickerView(app, "Choose a Color",
				func(s string) {
					c.input.SetText(s)
					if c.changedFunc != nil {
						c.changedFunc(s)
					}
					app.PopView()
				}, nil, func() {
					app.PopView()
				})
			app.StackView(cp)
			go cp.SelectColor(c.input.GetText())
		})
	return c
}

func (c *ColorPickerButton) makeLayout() {
	c.Flex.Clear().
		AddItem(c.labelText, c.labelWidth, 1, false).
		AddItem(c.input, c.fieldWidth-13, 1, false).
		AddItem(c.colorLabel, 3, 1, false).
		AddItem(c.button, 10, 1, false)
	c.Flex.SetBackgroundColor(c.bgColor)
	c.labelText.SetBackgroundColor(c.bgColor)
	c.labelText.SetTextColor(c.labelColor)
	c.input.SetBackgroundColor(c.bgColor)
	c.input.SetFieldBackgroundColor(c.fieldBgColor)
	c.input.SetFieldTextColor(c.fieldTextColor)
	c.input.SetChangedFunc(c.changedFunc)
	c.colorLabel.SetTextColor(tcell.GetColor(c.input.GetText()))
}

func (c *ColorPickerButton) SetText(text string) *ColorPickerButton {
	c.input.SetText(text)
	c.makeLayout()
	return c
}

func (c *ColorPickerButton) Focus(delegate func(p tview.Primitive)) {
	c.input.Focus(delegate)
}

func (c *ColorPickerButton) Blur() {
	c.input.Blur()
}

func (c *ColorPickerButton) HasFocus() bool {
	return c.input.HasFocus()
}

func (c *ColorPickerButton) SetFocusFunc(callback func()) *tview.Box {
	return c.input.SetFocusFunc(callback)
}

func (c *ColorPickerButton) SetBlurFunc(callback func()) *tview.Box {
	return c.input.SetBlurFunc(callback)
}

func (c *ColorPickerButton) SetLabel(label string) *ColorPickerButton {
	c.label = label
	c.labelText.SetText(label)
	return c
}

func (c *ColorPickerButton) SetFieldWidth(width int) *ColorPickerButton {
	c.fieldWidth = width
	c.makeLayout()
	return c
}

func (c *ColorPickerButton) GetLabel() string {
	return c.label
}

func (c *ColorPickerButton) SetFormAttributes(
	labelWidth int, labelColor, bgColor, fieldTextColor, fieldBgColor tcell.Color) tview.FormItem {
	c.labelWidth = labelWidth
	c.labelColor = labelColor
	c.bgColor = bgColor
	c.fieldTextColor = fieldTextColor
	c.fieldBgColor = fieldBgColor
	c.makeLayout()
	return c
}

func (c *ColorPickerButton) GetFieldWidth() int {
	return c.fieldWidth
}

func (c *ColorPickerButton) GetFieldHeight() int {
	return c.input.GetFieldHeight()
}

func (c *ColorPickerButton) SetFinishedFunc(handler func(key tcell.Key)) tview.FormItem {
	c.input.SetFinishedFunc(handler)
	return c
}

func (c *ColorPickerButton) SetDisabled(disabled bool) tview.FormItem {
	c.button.SetDisabled(disabled)
	return c
}
