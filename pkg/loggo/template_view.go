package loggo

import (
	"strings"

	"github.com/aurc/loggo/internal/colour"
	"github.com/aurc/loggo/pkg/config"
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
	toggleFullScreenCallback func()
	closeCallback            func()
}

func NewTemplateView(app Loggo, toggleFullScreenCallback, closeCallback func()) *TemplateView {
	tv := &TemplateView{
		Flex:                     *tview.NewFlex(),
		app:                      app,
		config:                   app.Config(),
		toggleFullScreenCallback: toggleFullScreenCallback,
		closeCallback:            closeCallback,
	}
	tv.makeUIComponents()
	tv.makeLayouts()
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
		SetBackgroundColor(colour.ColourBackgroundField)
}

func (t *TemplateView) makeLayouts() {
	t.makeContextMenu()
	t.Flex.Clear().SetDirection(tview.FlexColumn).
		AddItem(t.contextMenu, 30, 1, false).
		AddItem(t.table, 0, 2, true).
		SetBackgroundColor(colour.ColourBackgroundField)
	t.app.SetFocus(t.table)
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
	if r, _ := t.table.GetSelection(); r > 0 {
		t.contextMenu.AddItem("Move Up", "", 'u', func() {
			t.moveUp()
		})
		t.contextMenu.AddItem("Move Down", "", 'd', func() {
			t.moveDown()
		})
		t.contextMenu.AddItem("Remove Item", "", 'r', func() {
			t.confirmDelete()
		})
	}
	if t.closeCallback != nil {
		t.contextMenu.AddItem("Close", "", 'x', func() {
			t.closeCallback()
		})
	}
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
		cell = tview.NewTableCell(" " + k.Color.SetTextTagColor(k.Color.Foreground) + " ").
			SetAlign(tview.AlignCenter)
	case 4:
		cell = tview.NewTableCell(" " + k.Color.SetTextTagColor(k.Color.Background) + " ").
			SetAlign(tview.AlignCenter)
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
