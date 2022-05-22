package ui

import (
	"fmt"
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	SearchNoResults  = "[#ff5f00]No results"
	SearchFuzzyFound = `[yellow::b]%d[::-] potential results`
	SearchMatch      = `Match [yellow::b]%d[white::-] out of [green::b]%d[white::-] results`
)

const (
	ColourBackgroundField    = tcell.Color236
	ColourForegroundField    = tcell.ColorWhite
	ColourSelectedBackground = tcell.Color69
	ColourSelectedForeground = tcell.ColorWhite
	ColourSecondaryBorder    = tcell.Color240
)

var (
	fieldStyle = tcell.StyleDefault.
			Background(ColourBackgroundField).
			Foreground(ColourForegroundField)
	selectStyle = tcell.StyleDefault.
			Background(ColourSelectedBackground).
			Foreground(ColourSelectedForeground)
)

type JsonViewer struct {
	tview.Flex
	jsonRenderer   *JsonRenderer
	searchPane     *tview.Flex
	searchField    *tview.InputField
	searchInfo     *tview.TextView
	focusDelegator FocusDelegator
}

func MakeJsonViewer(focusDelegator FocusDelegator) *JsonViewer {
	jv := &JsonViewer{
		Flex:         *tview.NewFlex(),
		jsonRenderer: NewJsonRenderer().SetJsonConfigIndent(OrderSorted, "  "),
		searchPane:   tview.NewFlex().SetDirection(tview.FlexColumn),
		searchField:  tview.NewInputField(),
		searchInfo: tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(SearchNoResults).
			SetDynamicColors(true),
		focusDelegator: focusDelegator,
	}
	jv.SetDirection(tview.FlexRow)
	jv.ResetUI()
	jv.setupSearch()
	return jv
}

func (j *JsonViewer) setupSearch() {
	j.searchInfo.SetBackgroundColor(ColourBackgroundField).
		SetBorderColor(ColourSecondaryBorder).
		SetBorder(true)
	j.searchField.SetFieldStyle(fieldStyle).
		SetPlaceholder("Start typing the search...").
		SetAutocompleteStyles(tcell.Color236, fieldStyle, selectStyle).
		SetBorder(true).
		SetBackgroundColor(ColourBackgroundField)
	j.searchField.SetAutocompleteFunc(func(currentText string) (entries []string) {
		res := j.jsonRenderer.Search(currentText)
		if len(res) == 0 {
			j.searchInfo.SetText(SearchNoResults)
		} else {
			j.searchInfo.SetText(fmt.Sprintf(SearchFuzzyFound, len(res)))
		}
		return res
	})
	j.searchField.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			res := j.jsonRenderer.SearchTraversalSetup(j.searchField.GetText())
			j.searchInfo.SetText(fmt.Sprintf(SearchMatch,
				res.CurrentPosition,
				res.TotalPositions))
			j.focusDelegator.SetFocus(j.jsonRenderer)
			j.jsonRenderer.SetDoneFunc(func(key tcell.Key) {
				if j.jsonRenderer.isSearching {
					switch key {
					case tcell.KeyEnter, tcell.KeyTab:
						res := j.jsonRenderer.SearchTraverseNext()
						j.searchInfo.SetText(fmt.Sprintf(SearchMatch,
							res.CurrentPosition,
							res.TotalPositions))
					case tcell.KeyBacktab:
						res := j.jsonRenderer.SearchTraversePrev()
						j.searchInfo.SetText(fmt.Sprintf(SearchMatch,
							res.CurrentPosition,
							res.TotalPositions))
					case tcell.KeyEsc:
						j.searchField.SetText("")
						j.ResetUI()
						j.jsonRenderer.SearchTraversalReset()
					}
				}
			})
		case tcell.KeyEsc:
			j.searchField.SetText("")
			j.ResetUI()
			j.jsonRenderer.SearchTraversalReset()
		}

	})
}

func (j *JsonViewer) ResetUI() {
	j.Flex.Clear().
		AddItem(j.jsonRenderer, 0, 2, true)
	go func() {
		j.focusDelegator.SetFocus(j.jsonRenderer)
	}()
}

func (j *JsonViewer) SearchUI() {
	j.Flex.Clear().
		AddItem(j.searchPane.
			AddItem(j.searchField, 0, 5, true).
			AddItem(j.searchInfo, 0, 1, false),
			3, 0, false).
		AddItem(j.jsonRenderer, 0, 2, false)
	go func() {
		j.focusDelegator.SetFocus(j.searchField)
	}()
}

func (j *JsonViewer) SetJson(jText []byte) *JsonViewer {
	j.jsonRenderer.SetJson(jText)
	return j
}

func (j *JsonViewer) SetChangedFunc(handler func()) *JsonViewer {
	j.jsonRenderer.SetChangedFunc(handler)
	return j
}

func (j *JsonViewer) HandleShortcuts(event *tcell.EventKey) *tcell.EventKey {
	k := unicode.ToLower(event.Rune())
	if !j.searchField.HasFocus() && (k == 's' || k == 'S') {
		j.SearchUI()
		return nil
	} else if j.jsonRenderer.isSearching {
		//return nil
	}

	return event
}
