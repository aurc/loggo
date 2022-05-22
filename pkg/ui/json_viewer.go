package ui

import (
	"fmt"
	"strings"
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
	HelpText = ` ([yellow::b]%s[white:-:-]) %s`
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
	jsonRenderer     *JsonRenderer
	searchPane       *tview.Flex
	searchFuzzyField *tview.InputField
	searchInfo       *tview.TextView
	hotKeys          *tview.TextView
	focusDelegator   FocusDelegator
}

func MakeJsonViewer(focusDelegator FocusDelegator) *JsonViewer {
	jv := &JsonViewer{
		Flex:             *tview.NewFlex(),
		jsonRenderer:     NewJsonRenderer().SetJsonConfigIndent(OrderSorted, "  "),
		searchPane:       tview.NewFlex().SetDirection(tview.FlexColumn),
		searchFuzzyField: tview.NewInputField(),
		searchInfo: tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(SearchNoResults).
			SetDynamicColors(true),
		hotKeys: tview.NewTextView().
			SetTextAlign(tview.AlignLeft).
			SetRegions(false).
			SetDynamicColors(true),
		focusDelegator: focusDelegator,
	}
	jv.ResetUI()
	jv.setupFuzzySearch()
	return jv
}

func (j *JsonViewer) setupFuzzySearch() {

	j.searchInfo.SetBackgroundColor(ColourBackgroundField).
		SetBorderColor(ColourSecondaryBorder).
		SetBorder(true)
	j.searchFuzzyField.SetFieldStyle(fieldStyle).
		SetPlaceholder("Start typing the search...").
		SetAutocompleteStyles(tcell.Color236, fieldStyle, selectStyle).
		SetBorder(true).
		SetTitle("Fuzzy Search").
		SetBackgroundColor(ColourBackgroundField)
	j.searchFuzzyField.SetAutocompleteFunc(func(currentText string) (entries []string) {
		if len(currentText) == 0 {
			return
		}
		res := j.jsonRenderer.SearchFuzzy(currentText)
		if len(res) == 0 {
			j.searchInfo.SetText(SearchNoResults)
		} else {
			j.searchInfo.SetText(fmt.Sprintf(SearchFuzzyFound, len(res)))
		}
		return res
	})
	j.searchFuzzyField.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			res := j.jsonRenderer.SearchTraversalSetup(j.searchFuzzyField.GetText())
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
						j.searchFuzzyField.SetText("")
						j.ResetUI()
						j.jsonRenderer.SearchTraversalReset()
					}
				}
			})
		case tcell.KeyEsc:
			j.searchFuzzyField.SetText("")
			j.ResetUI()
			j.jsonRenderer.SearchTraversalReset()
		}

	})
}

func (j *JsonViewer) ResetUI() {
	columnsMain := j.Clear().SetDirection(tview.FlexColumn)
	rowMain := tview.NewFlex().SetDirection(tview.FlexRow)
	rowMain.
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexColumn).
			AddItem(j.jsonRenderer, 0, 1, true),
			0, 1, false)
	columnsMain.
		AddItem(rowMain, 0, 2, false).
		AddItem(j.ResetHotKeysUI(), 35, 0, false)

	go func() {
		j.focusDelegator.SetFocus(j.jsonRenderer)
	}()
}

func (j *JsonViewer) ResetHotKeysUI() *tview.TextView {
	j.hotKeys.Clear()
	j.hotKeys.SetBorder(true).SetTitle(" Hot Keys ")
	w := strings.Builder{}
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf(HelpText, "f", "  Fuzzy Word Search\n"))
	w.WriteString(fmt.Sprintf(HelpText, "s", "  Word Search\n"))
	w.WriteString(fmt.Sprintf(HelpText, "r", "  Regulat Expression Search\n"))
	w.WriteString(fmt.Sprintf(HelpText, "c", "  Copy All Clipboard\n"))
	w.WriteString(fmt.Sprintf(HelpText, "g", "  Start Of File\n"))
	w.WriteString(fmt.Sprintf(HelpText, "G", "  End Of File\n"))
	w.WriteString(fmt.Sprintf(HelpText, "↓", "  Scroll Down\n"))
	w.WriteString(fmt.Sprintf(HelpText, "↑", "  Scroll Up\nG"))
	return j.hotKeys.SetText(w.String())
}

func (j *JsonViewer) SearchFuzzySearchUI() {
	columnsMain := j.Clear().SetDirection(tview.FlexColumn)
	rowMain := tview.NewFlex().SetDirection(tview.FlexRow)
	rowMain.
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(j.searchPane.Clear().SetDirection(tview.FlexColumn).
				AddItem(j.searchFuzzyField, 0, 5, true).
				AddItem(j.searchInfo, 0, 1, false),
				3, 0, false).
			AddItem(j.jsonRenderer, 0, 1, true),
			0, 1, false)
	columnsMain.
		AddItem(rowMain, 0, 2, false).
		AddItem(j.ResetHotKeysUI(), 35, 0, false)

	j.hotKeys.Clear()
	j.hotKeys.SetBorder(true).SetTitle(" Hot Keys ")
	w := strings.Builder{}
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf(HelpText, "↲", "  Select Match (Enter)\n"))
	w.WriteString(fmt.Sprintf(HelpText, "↲", "  Next Result (Enter)\n"))
	w.WriteString(fmt.Sprintf(HelpText, "⇥", "  Next Result (Tab)\n"))
	w.WriteString(fmt.Sprintf(HelpText, "⇪[white::]+[yellow::b]⇥", "Prev Result (Shit+Tab)\n"))
	w.WriteString(fmt.Sprintf(HelpText, "ESC", "Quit Search"))
	j.hotKeys.SetText(w.String())

	go func() {
		j.focusDelegator.SetFocus(j.searchFuzzyField)
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
	if !j.searchFuzzyField.HasFocus() && (k == 'f') {
		j.SearchFuzzySearchUI()
		return nil
	} else if j.jsonRenderer.isSearching {
		//return nil
	}

	return event
}
