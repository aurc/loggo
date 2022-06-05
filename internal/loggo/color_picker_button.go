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
		labelText:   tview.NewTextView().SetText(label),
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

func (c *ColorPickerButton) SetFinishedFunc(handler func(key tcell.Key)) tview.FormItem {
	c.input.SetFinishedFunc(handler)
	return c
}
