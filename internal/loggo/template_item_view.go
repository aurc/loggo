package loggo

import (
	"strings"

	"github.com/aurc/loggo/internal/color"
	"github.com/aurc/loggo/internal/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TemplateItemView struct {
	tview.Flex
	app                      Loggo
	contextMenu              *tview.List
	form                     *tview.Form
	toggleFullScreenCallback func()
	closeCallback            func()
	key                      *config.Key
}

func NewTemplateItemView(app Loggo, key *config.Key, toggleFullScreenCallback, closeCallback func()) *TemplateItemView {
	tv := &TemplateItemView{
		Flex:                     *tview.NewFlex(),
		app:                      app,
		key:                      key,
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
		SetBackgroundColor(color.ColorBackgroundField)
	//selectType
	typeDD := tview.NewDropDown().
		SetLabel("Type").
		SetListStyles(color.FieldStyle, color.SelectStyle).
		AddOption(config.TypeString+"  ", nil).
		AddOption(config.TypeDateTime+"  ", nil).
		AddOption(config.TypeBool+"  ", nil).
		AddOption(config.TypeNumber+"  ", nil).SetSelectedFunc(func(text string, index int) {
		t.key.Type = config.Type(strings.TrimSpace(text))
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
		AddInputField("Key", t.key.Name, 40, nil, func(text string) {
			t.key.Name = strings.TrimSpace(text)
		}).
		AddFormItem(typeDD).
		AddInputField("Layout", t.key.Layout, 40, nil, func(text string) {
			t.key.Layout = strings.TrimSpace(text)
		}).
		AddFormItem(NewColorPickerButton(t.app, "Text Color", t.key.Color.Foreground, 40, func(text string) {
			t.key.Color.Foreground = strings.TrimSpace(text)
		})).
		AddFormItem(NewColorPickerButton(t.app, "Background Color", t.key.Color.Background, 40, func(text string) {
			t.key.Color.Background = strings.TrimSpace(text)
		}))
}

func (t *TemplateItemView) makeLayouts() {
	t.makeContextMenu()
	t.Flex.Clear().SetDirection(tview.FlexColumn).
		AddItem(t.contextMenu, 30, 1, false).
		AddItem(t.form, 0, 2, true).
		SetBackgroundColor(color.ColorBackgroundField)
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
