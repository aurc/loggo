package loggo

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"

	"github.com/rivo/tview"
)

type JsonView struct {
	tview.Flex
	app            *tview.Application
	textView       *tview.TextView
	searchInput    *tview.InputField
	searchType     *tview.DropDown
	statusBar      *tview.TextView
	contextMenu    *tview.List
	jMap           map[string]interface{}
	jText          []byte
	tagValToKey    map[string][]string
	tagValues      []string
	searchWord     string
	isSearching    bool
	indent         string
	searchStrategy Searchable
	withSearchTag  string
}

func NewJsonView(app *tview.Application) *JsonView {
	v := &JsonView{
		Flex:   *tview.NewFlex(),
		app:    app,
		indent: "  ",
	}
	v.makeUIComponents()
	v.makeLayouts(false)
	v.makeContextMenu()
	return v
}

func (j *JsonView) makeUIComponents() {
	j.textView = tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetRegions(true).
		SetDynamicColors(true).SetWrap(false)
	j.textView.
		SetBackgroundColor(ColourBackgroundField).
		SetBorderPadding(0, 0, 1, 1).
		SetBorder(true).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if !j.isSearching {
			switch event.Rune() {
			case 's', 'S':
				j.prepareCaseInsensitiveSearch()
				return nil
			case 'r', 'R':
				j.prepareRegexSearch()
				return nil
			case 'q':
				j.app.Stop()
				return nil
			}
			switch event.Key() {
			case tcell.KeyEsc:
				j.clearSearch()
				return nil
			}
		} else {
			switch event.Rune() {
			case 'n', 'N':
				j.Next()
				return nil
			case 'p', 'P':
				j.Prev()
				return nil
			case 'c', 'C':
				j.clearSearch()
				return nil
			}
			switch event.Key() {
			case tcell.KeyEsc:
				j.clearSearch()
				return nil
			case tcell.KeyTAB, tcell.KeyEnter:
				j.Next()
				return nil
			}
		}
		return event
	})

	j.contextMenu = tview.NewList()
	j.contextMenu.
		SetBorder(true).
		//SetBorderPadding(0, 0, 1, 1).
		SetTitle("Context Menu").
		SetBackgroundColor(ColourBackgroundField)

	j.searchInput = tview.NewInputField()
	j.searchInput.SetFieldStyle(FieldStyle).
		SetBorder(true).
		SetBackgroundColor(ColourBackgroundField)
	j.searchInput.SetChangedFunc(func(text string) {
		j.Search(text)
	})
	j.searchInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			j.clearSearch()
			return nil
		case tcell.KeyEnter:
			if len(j.searchInput.GetText()) > 0 {
				j.app.SetFocus(j.contextMenu)
				j.Next()
				return nil
			} else {
				j.clearSearch()
				return nil
			}
		}
		return event
	})

	j.statusBar = tview.NewTextView()
	j.statusBar.SetBackgroundColor(ColourSecondaryBorder)
}

func (j *JsonView) makeLayouts(search bool) {
	var rightMenuLayout *tview.Flex
	if !search {
		rightMenuLayout = tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(j.contextMenu, 0, 2, false)
	} else {
		rightMenuLayout = tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(j.searchInput, 3, 1, false).
			AddItem(j.contextMenu, 0, 2, false)
	}
	mainContent := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(j.textView, 0, 2, false).
		AddItem(rightMenuLayout, 40, 1, false)
	j.Flex.Clear().SetDirection(tview.FlexRow).
		AddItem(mainContent, 0, 2, false).
		AddItem(j.statusBar.Clear(), 1, 1, false)
}

func (j *JsonView) makeContextMenu() {
	if j.isSearching {
		j.contextMenu.Clear().
			ShowSecondaryText(false).
			AddItem("Next Result", "", 'n', func() {
				j.Next()
			}).
			AddItem("Previous Result", "", 'p', func() {
				j.Prev()
			}).
			AddItem("Clear Search", "", 'c', func() {
				j.clearSearch()
			}).
			SetBorderPadding(1, 1, 1, 1)
	} else {
		j.contextMenu.Clear().
			ShowSecondaryText(false).
			AddItem("Search Word", "", 's', func() {
				j.prepareCaseInsensitiveSearch()
			}).
			AddItem("Search Regex", "", 'r', func() {
				j.prepareRegexSearch()
			}).
			AddItem("Quit", "", 'q', func() {
				j.app.Stop()
			}).
			SetBorderPadding(1, 1, 1, 1)
	}
}

func (j *JsonView) prepareCaseInsensitiveSearch() {
	j.searchStrategy = MakeCaseInsensitiveSearch(j.statusBar)
	j.makeLayouts(true)
	j.searchInput.SetTitle("Search Word")
	j.app.SetFocus(j.searchInput)
}

func (j *JsonView) prepareRegexSearch() {
	j.searchStrategy = MakeRegexSearch(j.statusBar)
	j.makeLayouts(true)
	j.searchInput.SetTitle("Search Regex")
	j.app.SetFocus(j.searchInput)
}

func (j *JsonView) clearSearch() {
	j.app.SetFocus(j.textView)
	j.searchInput.SetText("")
	j.isSearching = false
	j.withSearchTag = ""
	j.setJson()
	j.makeLayouts(false)
	j.makeContextMenu()
}

func (j *JsonView) Search(word string) []int {
	j.isSearching = true
	j.makeContextMenu()
	j.searchStrategy.Clear()
	j.withSearchTag = word
	j.setJson()
	j.textView.
		Highlight(fmt.Sprintf(`%d`, j.searchStrategy.GetSearchPosition()-1)).
		ScrollToHighlight()

	j.searchStrategy.SetCurrentStatus()
	return nil
}

func (j *JsonView) Next() {
	j.searchStrategy.Next()
	j.textView.
		Highlight(fmt.Sprintf(`%d`, j.searchStrategy.GetSearchPosition()-1)).
		ScrollToHighlight()

	j.searchStrategy.SetCurrentStatus()
}
func (j *JsonView) Prev() {
	j.searchStrategy.Prev()
	j.textView.
		Highlight(fmt.Sprintf(`%d`, j.searchStrategy.GetSearchPosition()-1)).
		ScrollToHighlight()

	j.searchStrategy.SetCurrentStatus()
}

// SetJson sets a JSON and colourise accordingly, replacing any existing content.
func (j *JsonView) SetJson(jText []byte) *JsonView {
	j.jText = jText
	return j.setJson()
}

func (j *JsonView) setJson() *JsonView {
	j.jMap = make(map[string]interface{})
	if err := json.Unmarshal(j.jText, &j.jMap); err != nil {
		//j.SetText("[yellow]")
		panic(err)
	} else {
		text := &strings.Builder{}
		text.WriteString("{" + j.newLine())
		kc := len(j.jMap)
		i := 0
		keys := j.extractKeys(j.jMap)
		for _, k := range keys {
			v := j.jMap[k]
			j.processNode(k, v, j.indent, text, i+1 == kc)
			text.WriteString(j.newLine())
			i++
		}
		text.WriteString("}" + j.newLine())
		markedText := text.String()
		j.textView.SetText(markedText)
	}

	j.tagValues = []string{}
	for k := range j.tagValToKey {
		j.tagValues = append(j.tagValues, k)
	}

	return j
}

func (j *JsonView) processNode(k, v interface{}, indent string, text *strings.Builder, last bool) {
	word := j.captureWordSection(k, j.withSearchTag)
	if word != "" {
		k = word
	}
	key := fmt.Sprintf(`%s%s"%v"%s: `, indent, clField, k, clWhite)
	text.WriteString(key)
	switch tp := v.(type) {
	case int,
		float64,
		bool:
		j.processNumeric(text, v, "")
	case string:
		j.processString(text, v, "")
	case map[string]interface{}:
		j.processObject(text, v, j.indent+indent)
	case []interface{}:
		j.processArray(text, tp, j.indent+indent)
	}
	if !last {
		text.WriteString(",")
	}
}

func (j *JsonView) processArray(text *strings.Builder, tp []interface{}, indent string) {
	text.WriteString("[" + j.newLine())
	kc := len(tp)
	i := 0
	for _, n := range tp {
		j.processArrayItem(n, indent+j.indent, text, i+1 == kc)
		text.WriteString(j.newLine())
		i++
	}
	text.WriteString(j.computeIndent(indent[len(j.indent):]) + "]")
}

func (j *JsonView) processObject(text *strings.Builder, val interface{}, indent string) {
	text.WriteString(clString)
	text.WriteString(fmt.Sprintf(`[white::]{%s`, j.newLine()))

	vmap := val.(map[string]interface{})
	kc := len(vmap)
	i := 0

	keys := j.extractKeys(vmap)
	for _, k := range keys {
		v := vmap[k]
		j.processNode(k, v, indent+j.indent, text, i+1 == kc)
		text.WriteString(j.newLine())
		i++
	}
	text.WriteString(indent[len(j.indent):] + `}`)
}

func (j *JsonView) processString(text *strings.Builder, v interface{}, indent string) {
	val := fmt.Sprintf(`%v`, v)
	val = strings.ReplaceAll(val, "\"", "\\\"")
	val = strings.ReplaceAll(val, "\n", "\\n")
	if word := j.captureWordSection(v, j.withSearchTag); len(word) > 0 {
		val = word
	}
	text.WriteString(clString)
	text.WriteString(fmt.Sprintf(`%s"%v"`, j.computeIndent(indent), val))
	text.WriteString(clWhite)
}

func (j *JsonView) processNumeric(text *strings.Builder, v interface{}, indent string) {
	if word := j.captureWordSection(v, j.withSearchTag); len(word) > 0 {
		v = word
	}
	text.WriteString(clNumeric)
	text.WriteString(fmt.Sprintf("%s%v", j.computeIndent(indent), v))
	text.WriteString(clWhite)
}

func (j *JsonView) processArrayItem(v interface{}, indent string, text *strings.Builder, last bool) {
	switch tp := v.(type) {
	case int,
		float64,
		bool:
		j.processNumeric(text, v, indent)
	case string:
		j.processString(text, v, indent)
	case map[string]interface{}:
		j.processObject(text, v, indent)
	case []interface{}:
		j.processArray(text, tp, indent)
	}
	if !last {
		text.WriteString(",")
	}
}

func (j *JsonView) extractKeys(m map[string]interface{}) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys
}

func (j *JsonView) computeIndent(indent string) string {
	if len(j.indent) > 0 {
		return indent
	}
	return ""
}

func (j *JsonView) newLine() string {
	if len(j.indent) > 0 {
		return "\n"
	}
	return ""
}

func (j *JsonView) captureWordSection(text interface{}, withTag string) string {
	val := fmt.Sprintf("%v", text)
	tagged := len(withTag) > 0
	sel := ""
	if tagged {
		sel = j.searchStrategy.TagWord(withTag, val)
	}
	return sel
}
