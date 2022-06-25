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
	"strconv"
	"strings"

	"github.com/aurc/loggo/internal/color"
	"github.com/aurc/loggo/internal/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	maxFieldWidth = 26
)

type TemplateItemView struct {
	tview.Flex
	app                      Loggo
	contextMenu              *tview.List
	form                     *tview.Form
	toggleFullScreenCallback func()
	closeCallback            func()
	key                      *config.Key
	caseWhenTable            *tview.Table
	caseWhenLayout           *tview.Flex
	caseWhenForm             *tview.Form
	caseWhenCurrent          *config.ColorWhen
}

func NewTemplateItemView(app Loggo, key *config.Key, toggleFullScreenCallback, closeCallback func()) *TemplateItemView {
	tv := &TemplateItemView{
		Flex:                     *tview.NewFlex(),
		app:                      app,
		key:                      key,
		toggleFullScreenCallback: toggleFullScreenCallback,
		closeCallback:            closeCallback,
		caseWhenCurrent:          &config.ColorWhen{},
	}
	tv.makeUIComponents()
	tv.makeLayouts()
	return tv
}

func (t *TemplateItemView) makeUIComponents() {
	t.contextMenu = tview.NewList()
	t.contextMenu.
		SetBorder(true).
		SetTitle("Context Menu").
		SetBackgroundColor(color.ColorBackgroundField)

	// Main Form
	// text color
	colorable := func() *config.Color {
		if t.key == nil {
			return nil
		}
		return &t.key.Color
	}
	textColor := NewColorPickerButton(t.app, "Text Color",
		config.GetForegroundColorName(colorable, "white"), maxFieldWidth,
		func(text string) {
			t.key.Color.Foreground = strings.TrimSpace(text)
		})
	//text bg color
	textBgColor := NewColorPickerButton(t.app, "Background Color",
		config.GetBackgroundColorName(colorable, "black"), maxFieldWidth,
		func(text string) {
			t.key.Color.Background = strings.TrimSpace(text)
		})
	//selectType
	typeDD := tview.NewDropDown().
		SetLabel("Type [red]*").
		SetListStyles(color.FieldStyle, color.SelectStyle).
		AddOption(config.TypeString+"  ", nil).
		AddOption(config.TypeDateTime+"  ", nil).
		AddOption(config.TypeBool+"  ", nil).
		AddOption(config.TypeNumber+"  ", nil).
		SetSelectedFunc(func(text string, index int) {
			t.key.Type = config.Type(strings.TrimSpace(text))
			t.key.Color.Foreground = t.key.Type.GetColorName()
			t.key.Color.Background = "black"
			textColor.SetText(t.key.Color.Foreground)
			textBgColor.SetText(t.key.Color.Background)
		})
	currOpt := 0
	switch t.key.Type {
	case config.TypeString:
		currOpt = 0
	case config.TypeDateTime:
		currOpt = 1
	case config.TypeBool:
		currOpt = 2
	case config.TypeNumber:
		currOpt = 3
	}
	typeDD.SetCurrentOption(currOpt)

	t.form = tview.NewForm().
		SetFieldBackgroundColor(tcell.ColorDarkGray).
		SetFieldTextColor(tcell.ColorBlack).
		AddInputField("Key [red]*", t.key.Name, maxFieldWidth, nil, func(text string) {
			t.key.Name = strings.TrimSpace(text)
		}).
		AddFormItem(typeDD).
		AddInputField("Layout", t.key.Layout, maxFieldWidth, nil, func(text string) {
			t.key.Layout = strings.TrimSpace(text)
		}).
		AddFormItem(textColor).
		AddFormItem(textBgColor).
		AddInputField("Max Width", fmt.Sprintf("%d", t.key.MaxWidth), maxFieldWidth,
			func(textToCheck string, lastChar rune) bool {
				switch lastChar {
				case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
					return true
				default:
					return false
				}
			},
			func(text string) {
				w, _ := strconv.ParseInt(text, 10, 64)
				t.key.MaxWidth = int(w)
			})

	t.makeCaseWhenForm()
	t.caseWhenLayout = tview.NewFlex().SetDirection(tview.FlexRow)
	t.caseWhenLayout.SetBackgroundColor(tcell.ColorBlack)
	t.caseWhenTable = tview.NewTable()

	t.app.SetFocus(t.form)
	t.form.SetFocus(0)
}

func (t *TemplateItemView) makeCaseWhenForm() {
	// Case When Form
	caseWhenColorable := func() *config.Color {
		if t.caseWhenCurrent == nil {
			return nil
		}
		return &t.caseWhenCurrent.Color
	}
	caseWhenTextColor := NewColorPickerButton(t.app, "[::iu]then[::-], Text Color",
		config.GetForegroundColorName(caseWhenColorable, "white"), maxFieldWidth,
		func(text string) {
			t.caseWhenCurrent.Color.Foreground = strings.TrimSpace(text)
		})
	//text bg color
	caseWhenTextBgColor := NewColorPickerButton(t.app, "[::iu]and[::-], Background Color",
		config.GetBackgroundColorName(caseWhenColorable, "black"), maxFieldWidth,
		func(text string) {
			t.caseWhenCurrent.Color.Background = strings.TrimSpace(text)
		})

	t.caseWhenForm = tview.NewForm().
		SetFieldBackgroundColor(tcell.ColorDarkGray).
		SetFieldTextColor(tcell.ColorBlack).
		AddInputField("[::iu]when[::-] Value Matches", t.caseWhenCurrent.MatchValue, maxFieldWidth, nil, func(text string) {
			t.caseWhenCurrent.MatchValue = strings.TrimSpace(text)
		}).
		AddFormItem(caseWhenTextColor).
		AddFormItem(caseWhenTextBgColor).
		AddButton("Add", func() {
			t.key.ColorWhen = append(t.key.ColorWhen, *t.caseWhenCurrent)
			t.makeCaseWhenData()
		}).
		AddButton("Update Selected", func() {
			r, _ := t.caseWhenTable.GetSelection()
			if r > 0 && r-1 < len(t.key.ColorWhen) {
				t.key.ColorWhen[r-1] = *t.caseWhenCurrent
			}
			t.makeCaseWhenData()
		}).
		AddButton("Delete Selected", func() {
			r, _ := t.caseWhenTable.GetSelection()
			if r > 0 {
				idx := r - 1
				cw := make([]config.ColorWhen, 0)
				for i := range t.key.ColorWhen {
					if i != idx {
						cw = append(cw, t.key.ColorWhen[i])
					}
				}
				t.key.ColorWhen = cw
				t.makeCaseWhenData()
			}
		})
}
func (t *TemplateItemView) makeLayouts() {
	t.makeContextMenu()

	t.caseWhenLayout.Clear().
		AddItem(t.caseWhenForm, 9, 1, false).
		AddItem(t.caseWhenTable, 0, 1, false)

	mainForm := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(t.contextMenu, 3, 1, false).
		AddItem(t.form, 0, 1, false)
	formLayout := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(mainForm, maxFieldWidth+20, 1, false).
		AddItem(t.caseWhenLayout, 0, 1, false)

	t.Flex.Clear().SetDirection(tview.FlexRow).
		//AddItem(t.contextMenu, 3, 1, false).
		AddItem(formLayout, 0, 2, true).
		SetBackgroundColor(color.ColorBackgroundField)

	t.makeCaseWhenData()
}

func (t *TemplateItemView) makeCaseWhenData() {
	t.caseWhenTable.Clear().
		SetSelectable(true, true).
		SetFixed(1, 1).
		SetSeparator(tview.Borders.Vertical)
	t.caseWhenTable.SetCell(0, 0,
		tview.NewTableCell(" Match Value ").
			SetTextColor(tcell.ColorLightGray).
			SetSelectable(false).
			SetAlign(tview.AlignCenter))
	t.caseWhenTable.SetCell(0, 1,
		tview.NewTableCell(" Text Color ").
			SetTextColor(tcell.ColorLightGray).
			SetSelectable(false).
			SetAlign(tview.AlignCenter))
	t.caseWhenTable.SetCell(0, 2,
		tview.NewTableCell(" Background ").
			SetTextColor(tcell.ColorLightGray).
			SetSelectable(false).
			SetAlign(tview.AlignCenter))
	t.caseWhenTable.SetCell(0, 3,
		tview.NewTableCell(" ✎ ").
			SetTextColor(tcell.ColorLightGray).
			SetSelectable(false).
			SetAlign(tview.AlignCenter))
	for i, k := range t.key.ColorWhen {
		t.caseWhenTable.SetCell(i+1, 0,
			tview.NewTableCell(k.Color.SetTextTagColor(k.MatchValue)).
				SetTextColor(tcell.ColorYellow).
				SetSelectable(false).
				SetAlign(tview.AlignCenter))
		t.caseWhenTable.SetCell(i+1, 1,
			tview.NewTableCell(fmt.Sprintf(` [%s] ■ [-]│ %s `, k.Color.Foreground, k.Color.Foreground)).
				SetSelectable(false).
				SetAlign(tview.AlignLeft))
		t.caseWhenTable.SetCell(i+1, 2,
			tview.NewTableCell(fmt.Sprintf(` [%s] ■ [-]│ %s `, k.Color.Background, k.Color.Background)).
				SetSelectable(false).
				SetAlign(tview.AlignLeft))
		t.caseWhenTable.SetCell(i+1, 3,
			tview.NewTableCell(" ✎ ").
				SetTextColor(tcell.ColorLightGray).
				SetSelectable(true).
				SetTextColor(tcell.ColorYellow).
				SetAlign(tview.AlignCenter))

		t.caseWhenTable.SetSelectionChangedFunc(func(row, column int) {
			if row > 0 {
				v := t.key.ColorWhen[row-1]
				t.caseWhenCurrent = &v
				t.makeCaseWhenForm()
				t.makeLayouts()
			}
		})
	}
}

func (t *TemplateItemView) makeContextMenu() {
	t.contextMenu.Clear().ShowSecondaryText(false).SetBorderPadding(0, 0, 1, 1)
	t.contextMenu.
		ShowSecondaryText(false)
	if t.toggleFullScreenCallback != nil {
		t.contextMenu.AddItem("Toggle Full Screen", "", 'f', func() {
			t.toggleFullScreenCallback()
		})
	}

	if t.closeCallback != nil {
		t.contextMenu.AddItem("Done", "", 'x', func() {
			t.closeCallback()
		})
	}
}
