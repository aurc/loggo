package ui

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/rivo/tview"
)

const (
	clField   = "[#ffaf00::b]"
	clWhite   = "[#ffffff::-]"
	clNumeric = "[#00afff]"
	clString  = "[#6A9F59]"
)

type JsonRenderer struct {
	tview.TextView
	jMap           map[string]interface{}
	viewerConfig   *viewerConfig
	selectionCount int
	tagValToKey    map[string][]string
	fuzzyTags      map[string]bool
	tagValues      []string
	searchWord     string
	searchWordIdx  int
	isSearching    bool
}

// NewJsonRenderer returns a new json view.
func NewJsonRenderer() *JsonRenderer {
	jv := &JsonRenderer{
		TextView: *tview.NewTextView().
			SetDynamicColors(true).
			SetRegions(true),
		viewerConfig: &viewerConfig{
			renderType:   RenderJSON,
			shouldIndent: false,
			indent:       "",
			ordering:     OrderNone,
		},
		tagValToKey: make(map[string][]string),
	}
	jv.SetBackgroundColor(tcell.Color16)
	return jv
}

func (j *JsonRenderer) SearchFuzzy(text string) []string {
	ranks := fuzzy.RankFind(text, j.tagValues)
	sort.Sort(ranks)
	var results []string
	for _, r := range ranks {
		results = append(results, r.Target)
	}
	return results
}

type SearchTraversalState struct {
	CurrentPosition int
	TotalPositions  int
}

func (j *JsonRenderer) SearchTraversalSetup(text string) *SearchTraversalState {
	jText := j.GetText(true)
	j.Clear().SetDynamicColors(true).
		SetRegions(true)
	j.setJson([]byte(jText), text)
	j.searchWord = text
	regions := j.tagValToKey[j.searchWord]
	j.searchWord = text
	j.searchWordIdx = len(regions) - 1
	j.isSearching = true
	return j.SearchTraverseNext()
}

func (j *JsonRenderer) SearchTraversalReset() {
	j.searchWord = ""
	j.searchWordIdx = 0
	j.Highlight()
	j.isSearching = false
}

func (j *JsonRenderer) SearchTraverseNext() *SearchTraversalState {
	return j.searchTraversal(true)
}

func (j *JsonRenderer) SearchTraversePrev() *SearchTraversalState {
	return j.searchTraversal(false)
}

func (j *JsonRenderer) SetJsonConfig(ordering Ordering) *JsonRenderer {
	j.viewerConfig = &viewerConfig{
		renderType:   RenderJSON,
		shouldIndent: false,
		ordering:     ordering,
	}
	return j
}

func (j *JsonRenderer) SetJsonConfigIndent(ordering Ordering, indent string) *JsonRenderer {
	j.viewerConfig = &viewerConfig{
		renderType:   RenderJSON,
		shouldIndent: true,
		indent:       indent,
		ordering:     ordering,
	}
	return j
}

// SetJson sets a JSON and colourise accordingly, replacing any existing content.
func (j *JsonRenderer) SetJson(jText []byte) *JsonRenderer {
	j.selectionCount = 0
	return j.setJson(jText, "")
}

func (j *JsonRenderer) setJson(jText []byte, withTag string) *JsonRenderer {
	j.jMap = make(map[string]interface{})
	j.fuzzyTags = make(map[string]bool)
	if err := json.Unmarshal(jText, &j.jMap); err != nil {
		//j.SetText("[yellow]")
		panic(err)
	} else {
		cfg := j.viewerConfig
		text := &strings.Builder{}
		if cfg.renderType == RenderJSON {
			text.WriteString("{" + j.newLine())
		}
		kc := len(j.jMap)
		i := 0
		keys := j.extractKeys(j.jMap)
		for _, k := range keys {
			v := j.jMap[k]
			j.processNode(k, v, cfg.indent, text, i+1 == kc, strings.ToLower(withTag))
			text.WriteString(j.newLine())
			i++
		}
		if cfg.renderType == RenderJSON {
			text.WriteString("}" + j.newLine())
		}
		markedText := text.String()
		j.SetText(markedText)
	}

	j.tagValues = []string{}
	for k := range j.fuzzyTags {
		j.tagValues = append(j.tagValues, k)
	}
	for k := range j.tagValToKey {
		j.tagValues = append(j.tagValues, k)
	}

	return j
}

func (j *JsonRenderer) searchTraversal(isNext bool) *SearchTraversalState {
	regions := j.tagValToKey[j.searchWord]
	isPrev := !isNext
	if isNext {
		if j.searchWordIdx+1 == len(regions) {
			j.searchWordIdx = 0
		} else {
			j.searchWordIdx++
		}
	} else if isPrev {
		if j.searchWordIdx == 0 {
			j.searchWordIdx = len(regions) - 1
		} else {
			j.searchWordIdx--
		}
	}
	regionID := regions[j.searchWordIdx]
	j.Highlight(regionID).ScrollToHighlight()
	return &SearchTraversalState{
		CurrentPosition: j.searchWordIdx + 1,
		TotalPositions:  len(regions),
	}
}

func (j *JsonRenderer) handleShortcuts(event *tcell.EventKey) *tcell.EventKey {

	return event
}

func (j *JsonRenderer) processNode(k, v interface{}, indent string, text *strings.Builder, last bool,
	withTag string) {
	word := j.captureWordSection(k, withTag)
	if word != "" {
		k = word
	}
	key := fmt.Sprintf(`%s%s"%v[""]"%s: `, indent, clField, k, clWhite)
	text.WriteString(key)
	cfg := j.viewerConfig
	switch tp := v.(type) {
	case int,
		float64,
		bool:
		j.processNumeric(text, v, "", withTag)
	case string:
		j.processString(text, v, "", withTag)
	case map[string]interface{}:
		j.processObject(text, v, cfg.indent+indent, withTag)
	case []interface{}:
		j.processArray(text, tp, cfg.indent+indent, withTag)
	}
	if !last {
		text.WriteString(",")
	}
}

func (j *JsonRenderer) processArray(text *strings.Builder, tp []interface{}, indent, withTag string) {
	text.WriteString("[" + j.newLine())
	kc := len(tp)
	i := 0
	for _, n := range tp {
		j.processArrayItem(n, indent+j.viewerConfig.indent, text, i+1 == kc, withTag)
		text.WriteString(j.newLine())
		i++
	}
	text.WriteString(j.computeIndent(indent[len(j.viewerConfig.indent):]) + "]")
}

func (j *JsonRenderer) processObject(text *strings.Builder, val interface{}, indent, withTag string) {
	text.WriteString(clString)
	text.WriteString(fmt.Sprintf(`{%s`, j.newLine()))
	cfg := j.viewerConfig

	vmap := val.(map[string]interface{})
	kc := len(vmap)
	i := 0

	keys := j.extractKeys(vmap)
	for _, k := range keys {
		v := vmap[k]
		j.processNode(k, v, indent+j.viewerConfig.indent, text, i+1 == kc, withTag)
		text.WriteString(j.newLine())
		i++
	}
	text.WriteString(indent[len(cfg.indent):] + `}`)
}

func (j *JsonRenderer) processString(text *strings.Builder, v interface{}, indent, withTag string) {
	val := fmt.Sprintf(`%v`, v)
	val = strings.ReplaceAll(val, "\"", "\\\"")
	val = strings.ReplaceAll(val, "\n", "\\n")
	if word := j.captureWordSection(v, withTag); len(word) > 0 {
		val = word
	}
	text.WriteString(clString)
	text.WriteString(fmt.Sprintf(`%s"%v"`, j.computeIndent(indent), val))
	text.WriteString(clWhite)
}

func (j *JsonRenderer) processNumeric(text *strings.Builder, v interface{}, indent string, withTag string) {
	if word := j.captureWordSection(v, withTag); len(word) > 0 {
		v = word
	}
	text.WriteString(clNumeric)
	text.WriteString(fmt.Sprintf("%s%v", j.computeIndent(indent), v))
	text.WriteString(clWhite)
}

func (j *JsonRenderer) processArrayItem(v interface{}, indent string, text *strings.Builder, last bool, withTag string) {
	switch tp := v.(type) {
	case int,
		float64,
		bool:
		j.processNumeric(text, v, indent, withTag)
	case string:
		j.processString(text, v, indent, withTag)
	case map[string]interface{}:
		j.processObject(text, v, indent, withTag)
	case []interface{}:
		j.processArray(text, tp, indent, withTag)
	}
	if !last {
		text.WriteString(",")
	}
}

func (j *JsonRenderer) extractKeys(m map[string]interface{}) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	if j.viewerConfig.ordering == OrderSorted {
		sort.Strings(keys)
	}
	return keys
}

func (j *JsonRenderer) computeIndent(indent string) string {
	cfg := j.viewerConfig
	if cfg.shouldIndent {
		return indent
	}
	return ""
}

func (j *JsonRenderer) newLine() string {
	cfg := j.viewerConfig
	if cfg.shouldIndent {
		return "\n"
	}
	return ""
}

func (j *JsonRenderer) captureWordSection(text interface{}, withTag string) string {
	val := fmt.Sprintf("%v", text)
	tagged := len(withTag) > 0
	sel := ""
	if tagged && strings.ToLower(val) == withTag {
		tagID := fmt.Sprintf("%d", j.selectionCount)
		if _, ok := j.tagValToKey[val]; !ok {
			j.tagValToKey[val] = []string{tagID}
		} else {
			j.tagValToKey[val] = append(j.tagValToKey[val], tagID)
		}
		sel = fmt.Sprintf(`["%s"]%v[""]`, tagID, text)
		j.selectionCount++
	}
	if _, ok := j.fuzzyTags[val]; !ok {
		j.fuzzyTags[val] = true
	}
	return sel
}

//func (j *JsonRenderer) captureWordSection(text interface{}) string {
//	tagID := fmt.Sprintf("%d", j.selectionCount)
//	val := fmt.Sprintf("%v", text)
//	if _, ok := j.tagValToKey[val]; !ok {
//		j.tagValToKey[val] = []string{tagID}
//	} else {
//		j.tagValToKey[val] = append(j.tagValToKey[val], tagID)
//	}
//	sel := fmt.Sprintf(`["%s"]%v[""]`, tagID, text)
//	j.selectionCount++
//	return sel
//}

type Ordering string

const (
	OrderSorted = "OrderSorted"
	OrderNone   = "OrderNone"
)

type RenderType string

const (
	RenderJSON = "JSON"
)

type viewerConfig struct {
	renderType   RenderType
	shouldIndent bool
	indent       string
	ordering     Ordering
}
