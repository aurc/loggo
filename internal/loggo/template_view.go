/*
Copyright © 2022 Aurelio Calegari

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
	"os"
	"strings"
	"time"

	"github.com/aurc/loggo/internal/color"
	"github.com/aurc/loggo/internal/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TemplateView struct {
	tview.Flex
	app                      Loggo
	config                   *config.Config
	table                    *tview.Table
	data                     *TemplateData
	contextMenu              *tview.List
	showQuit                 bool
	toggleFullScreenCallback func()
	closeCallback            func()
}

func NewTemplateView(app Loggo, showQuit bool, toggleFullScreenCallback, closeCallback func()) *TemplateView {
	tv := &TemplateView{
		Flex:                     *tview.NewFlex(),
		app:                      app,
		config:                   app.Config(),
		showQuit:                 showQuit,
		toggleFullScreenCallback: toggleFullScreenCallback,
		closeCallback:            closeCallback,
	}
	tv.makeUIComponents()
	tv.makeLayouts()
	go func() {
		time.Sleep(100 * time.Millisecond)
		if len(tv.config.Keys) > 0 {
			tv.table.Select(1, 0)
			tv.app.Draw()
		}
	}()
	return tv
}

func (t *TemplateView) makeUIComponents() {
	t.data = &TemplateData{
		templateView: t,
	}
	t.table = tview.NewTable().
		SetSelectable(true, false).
		SetSeparator(tview.Borders.Vertical).
		SetContent(t.data)
	t.table.
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			selected := false
			if r, _ := t.table.GetSelection(); r > 0 {
				selected = true
			}
			switch event.Rune() {
			case 'f', 'F':
				if t.toggleFullScreenCallback != nil {
					t.toggleFullScreenCallback()
					return nil
				}
			case 'a', 'A':
				t.addEntry()
				return nil
			case 's', 'S':
				t.saveForm()
				return nil
			case 'e', 'E':
				if selected {
					t.editEntry()
					return nil
				}
			case 'q':
				if t.showQuit {
					t.app.Stop()
					return nil
				}
			case 'u', 'U':
				if selected {
					t.moveUp()
					return nil
				}
			case 'd', 'D':
				if selected {
					t.moveDown()
					return nil
				}
			case 'x', 'X':
				if t.closeCallback != nil {
					t.closeCallback()
					return nil
				}
			case 'r', 'R':
				if selected {
					t.confirmDelete()
					return nil
				}
			}
			return event
		})
	t.table.SetSelectionChangedFunc(func(row, column int) {
		t.makeContextMenu()
	})

	t.contextMenu = tview.NewList()
	t.contextMenu.
		SetBorder(true).
		SetTitle("Context Menu").
		SetBackgroundColor(color.ColorBackgroundField)
}

func (t *TemplateView) makeLayouts() {
	t.makeContextMenu()
	t.Flex.Clear().SetDirection(tview.FlexColumn).
		AddItem(t.contextMenu, 30, 1, false).
		AddItem(t.table, 0, 2, true).
		SetBackgroundColor(color.ColorBackgroundField)
	t.app.SetFocus(t.table)
}

func (t *TemplateView) makeSaveLayouts() {
	bar, input := t.makeSaveUI()
	t.Flex.Clear().SetDirection(tview.FlexRow).
		AddItem(bar, 3, 1, false).
		AddItem(t.table, 0, 1, false)
	t.app.SetFocus(input)
}

func (t *TemplateView) makeSaveUI() (*tview.Flex, *tview.InputField) {
	dirName, _ := os.UserHomeDir()
	if len(t.config.LastSavedName) > 0 {
		dirName = t.config.LastSavedName
	} else {
		dirName = fmt.Sprintf(`%s%c`, dirName, os.PathSeparator)
	}
	saveBar := tview.NewFlex().SetDirection(tview.FlexColumn)
	saveBar.SetBackgroundColor(color.ColorBackgroundField).SetBorder(true).SetTitle("Save Template As...")

	saveInput := tview.NewInputField().SetText(dirName)
	saveInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			t.save(saveInput.GetText())
			t.makeLayouts()
		case tcell.KeyEscape:
			t.makeLayouts()
		}
		return event
	})
	buttonSave := tview.NewButton("Save").SetSelectedFunc(func() {
		// Save
		t.save(saveInput.GetText())
		t.makeLayouts()
	})
	buttonCancel := tview.NewButton("Cancel").SetSelectedFunc(func() {
		// Cancel
		t.makeLayouts()
	})
	saveBar.AddItem(tview.NewBox(), 1, 1, false).
		AddItem(saveInput, 0, 1, true).
		AddItem(tview.NewBox(), 1, 1, false).
		AddItem(buttonSave, 6, 1, false).
		AddItem(tview.NewBox(), 1, 1, false).
		AddItem(buttonCancel, 6, 1, false).
		AddItem(tview.NewBox(), 1, 1, false)
	return saveBar, saveInput
}

func (t *TemplateView) makeContextMenu() {
	t.contextMenu.Clear().ShowSecondaryText(false).SetBorderPadding(0, 0, 1, 1)
	t.contextMenu.
		ShowSecondaryText(false)
	if t.toggleFullScreenCallback != nil {
		t.contextMenu.AddItem("Toggle Full Screen", "", 'f', func() {
			t.toggleFullScreenCallback()
		})
	}
	t.contextMenu.AddItem("Add New", "", 'a', func() {
		t.addEntry()
	})
	if r, _ := t.table.GetSelection(); r > 0 {
		t.contextMenu.AddItem("Edit", "", 'e', func() {
			t.editEntry()
		})
		t.contextMenu.AddItem("Move Up", "", 'u', func() {
			t.moveUp()
		})
		t.contextMenu.AddItem("Move Down", "", 'd', func() {
			t.moveDown()
		})
		t.contextMenu.AddItem("Remove", "", 'r', func() {
			t.confirmDelete()
		})
	}
	t.contextMenu.AddItem("Save", "", 's', func() {
		t.saveForm()
	})
	if t.closeCallback != nil {
		t.contextMenu.AddItem("Done", "", 'x', func() {
			t.closeCallback()
		})
	}
	if t.showQuit {
		t.contextMenu.AddItem("Quit", "", 'q', func() {
			t.app.Stop()
		})
	}
}

func (t *TemplateView) save(fileName string) {
	if err := t.config.Save(fileName); err != nil {
		t.app.ShowPrefabModal(
			fmt.Sprintf(`Failed to save! Error: %v`, err), 40, 10,
			tview.NewButton("Ok").SetSelectedFunc(func() {
				t.app.DismissModal()
			}))
	} else {
		t.app.ShowPrefabModal(
			fmt.Sprintf(`File %v saved successfully!`, fileName), 40, 10,
			tview.NewButton("Ok").SetSelectedFunc(func() {
				t.app.DismissModal()
			}))
	}
}

func (t *TemplateView) Close() {
	if t.closeCallback != nil {
		t.closeCallback()
	}
}

func (t *TemplateView) saveForm() {
	t.makeSaveLayouts()
}

func (t *TemplateView) addEntry() {
	v := &config.Key{
		Type: config.TypeString,
		Color: config.Color{
			Foreground: "white",
			Background: "black",
		},
	}
	t.app.StackView(NewTemplateItemView(t.app, v, nil, func() {
		t.app.PopView()
		kn := strings.TrimSpace(v.Name)
		if len(kn) > 0 {
			v.Name = kn
			t.config.Keys = append(t.config.Keys, *v)
			t.table.Select(len(t.config.Keys), 0)
		}
	}))
}

func (t *TemplateView) editEntry() {
	r, _ := t.table.GetSelection()
	t.app.StackView(NewTemplateItemView(t.app, &t.config.Keys[r-1], nil, func() {
		t.app.PopView()
	}))
}

func (t *TemplateView) moveUp() {
	r, _ := t.table.GetSelection()
	finalRow := r
	r = r - 1
	keys := t.config.Keys
	if r > 0 {
		curr := keys[r]
		keys[r] = keys[r-1]
		keys[r-1] = curr
		finalRow = r
	} else if len(keys) > 1 {
		t.config.Keys = append(keys[1:], keys[r])
		finalRow = len(t.config.Keys)
	}
	t.makeLayouts()
	t.table.Select(finalRow, 0)
}

func (t *TemplateView) moveDown() {
	r, _ := t.table.GetSelection()
	finalRow := r
	r = r - 1
	keys := t.config.Keys
	if r < len(keys)-1 {
		curr := keys[r]
		keys[r] = keys[r+1]
		keys[r+1] = curr
		finalRow = r + 2
	} else if len(keys) > 1 {
		t.config.Keys = append([]config.Key{keys[r]}, keys[:len(keys)-1]...)
		finalRow = 1
	}
	t.makeLayouts()
	t.table.Select(finalRow, 0)
}

func (t *TemplateView) confirmDelete() {
	t.app.ShowPrefabModal("Are you sure you want to remove this entry", 40, 10,
		tview.NewButton("Delete").
			SetSelectedFunc(func() {
				var newKeys []config.Key
				r, _ := t.table.GetSelection()
				for i, k := range t.config.Keys {
					if i != r-1 {
						newKeys = append(newKeys, k)
					}
				}
				t.config.Keys = newKeys
				t.makeLayouts()
				t.app.DismissModal()
			}),
		tview.NewButton("Cancel").
			SetSelectedFunc(func() {
				t.app.DismissModal()
			}),
	)
}

type TemplateData struct {
	tview.TableContentReadOnly
	templateView *TemplateView
}

var columnNames = []string{" Key ", " Type ", " Layout ", " Text Color ", " Background ", " Color Match "}

func (d *TemplateData) GetCell(row, column int) *tview.TableCell {
	if row == -1 || len(d.templateView.config.Keys) < row-1 || column == -1 {
		return nil
	}
	c := d.templateView.config
	// Set Headers
	if row == 0 {
		return tview.NewTableCell(columnNames[column]).
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignCenter).
			SetBackgroundColor(tcell.ColorBlack).SetSelectable(false)
	}
	// Set Body Cells
	k := c.Keys[row-1]
	var cell *tview.TableCell
	switch column {
	case 0:
		cell = tview.NewTableCell(" " + k.Name + " ")
	case 1:
		cell = tview.NewTableCell(" " + string(k.Type) + " ").
			SetTextColor(k.Type.GetColor()).
			SetAlign(tview.AlignCenter)
	case 2:
		cell = tview.NewTableCell(" " + k.Layout + " ").
			SetAlign(tview.AlignCenter)
	case 3:
		fg := k.Color.Foreground
		cell = tview.NewTableCell(
			fmt.Sprintf(` [%s] ■ [-] │ %s `, fg, fg)).
			SetAlign(tview.AlignLeft)
	case 4:
		bg := k.Color.Background
		cell = tview.NewTableCell(
			fmt.Sprintf(` [%s] ■ [-] │ %s `, bg, bg)).
			SetAlign(tview.AlignLeft)
	case 5:
		caseWhen := strings.Builder{}
		caseWhen.WriteString(" ")
		for _, cw := range k.ColorWhen {
			caseWhen.WriteString(cw.Color.SetTextTagColor(cw.MatchValue))
			caseWhen.WriteString(" ")
		}
		cell = tview.NewTableCell(caseWhen.String())
	}
	return cell
}

func (d *TemplateData) GetRowCount() int {
	return len(d.templateView.config.Keys) + 1
}

func (d *TemplateData) GetColumnCount() int {
	return len(columnNames)
}
