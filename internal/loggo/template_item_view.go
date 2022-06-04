package loggo

import (
	"fmt"
	"sort"
	"strings"

	"github.com/aurc/loggo/internal/colour"
	"github.com/aurc/loggo/internal/config"
	"github.com/gdamore/tcell/v2"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/rivo/tview"
)

type TemplateItemView struct {
	tview.Flex
	app                      Loggo
	contextMenu              *tview.List
	form                     *tview.Form
	toggleFullScreenCallback func()
	closeCallback            func()
}

func NewTemplateItemView(app Loggo, toggleFullScreenCallback, closeCallback func()) *TemplateItemView {
	tv := &TemplateItemView{
		Flex:                     *tview.NewFlex(),
		app:                      app,
		toggleFullScreenCallback: toggleFullScreenCallback,
		closeCallback:            closeCallback,
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
		SetBackgroundColor(colour.ColourBackgroundField)
	typeDD := tview.NewDropDown().
		SetLabel("Type").
		SetListStyles(colour.FieldStyle, colour.SelectStyle).
		AddOption(config.TypeString+"  ", nil).
		AddOption(config.TypeDateTime+"  ", nil).
		AddOption(config.TypeBool+"  ", nil).
		AddOption(config.TypeNumber+"  ", nil)
	//colorDD := tview.NewDropDown().
	//	SetLabel("Text Colour").
	//	SetListStyles(colour.FieldStyle, colour.SelectStyle)
	var cols []string
	for col := range tcell.ColorNames {
		label := fmt.Sprintf(` [%s] ■ [-] %s `, col, col)
		cols = append(cols, label)
	}
	colorDD := tview.NewInputField().SetFieldStyle(colour.FieldStyle).
		SetLabel("Text Colour").
		SetAutocompleteStyles(tcell.Color236, colour.FieldStyle, colour.SelectStyle)
	colorDD.
		SetAutocompleteFunc(func(currentText string) (entries []string) {
			if len(currentText) == 0 {
				return
			}
			ranks := fuzzy.RankFind(currentText, cols)
			sort.Sort(ranks)
			var results []string
			for _, r := range ranks {
				results = append(results, r.Target)
			}
			return results
		}).SetChangedFunc(func(text string) {
		nt := strings.TrimSpace(text)
		idx := strings.LastIndex(nt, " ")
		if idx != -1 {
			nt = nt[idx+1:]
			colorDD.SetText(nt)
		}
	})
	t.form = tview.NewForm().
		SetFieldBackgroundColor(tcell.ColorDarkGray).
		SetFieldTextColor(tcell.ColorBlack).
		AddInputField("Key", "", 40, nil, nil).
		AddFormItem(typeDD).
		AddInputField("Layout", "", 40, nil, nil).
		AddFormItem(colorDD).
		AddInputField("Background Colour", "", 40, nil, nil)
}

func (t *TemplateItemView) makeLayouts() {
	t.makeContextMenu()
	t.Flex.Clear().SetDirection(tview.FlexColumn).
		AddItem(t.contextMenu, 30, 1, false).
		AddItem(t.form, 0, 2, true).
		SetBackgroundColor(colour.ColourBackgroundField)
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
		t.contextMenu.AddItem("Close", "", 'x', func() {
			t.closeCallback()
		})
	}
}
