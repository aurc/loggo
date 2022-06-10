package loggo

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aurc/loggo/internal/color"
	"github.com/aurc/loggo/internal/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type LogView struct {
	tview.Flex
	app                *LoggoApp
	input              <-chan string
	table              *tview.Table
	jsonView           *JsonView
	data               *LogData
	templateView       *TemplateView
	layout             *tview.Flex
	config             *config.Config
	mainMenu           *tview.Flex
	linesView          *tview.TextView
	logFullScreen      bool
	templateFullScreen bool
	inSlice            []map[string]interface{}
	globalCount        int64
}

func NewLogReader(app *LoggoApp, input <-chan string) *LogView {
	lv := &LogView{
		Flex:   *tview.NewFlex(),
		app:    app,
		config: app.Config(),
		input:  input,
	}
	lv.makeUIComponents()
	lv.makeLayouts()
	lv.read()
	return lv
}

const (
	parseErr    = "$_parseErr"
	textPayload = "$_textPayload"
)

func (l *LogView) read() {
	go func() {
		var sampling []map[string]interface{}
		samplingCount := 0
		if len(l.config.LastSavedName) == 0 {
			samplingCount = 50
		}
		for {
			t := <-l.input
			if len(t) > 0 {
				l.globalCount++
				m := make(map[string]interface{})
				err := json.Unmarshal([]byte(t), &m)
				if err != nil {
					m[parseErr] = err.Error()
					m[textPayload] = t
				}
				if l.globalCount <= int64(samplingCount) {
					sampling = append(sampling, m)
				} else if len(sampling) <= samplingCount {
					l.processSampleForConfig(sampling)
				}
				l.inSlice = append(l.inSlice, m)
				l.updateLineView()
				l.app.Draw()
			}
		}
	}()
}

func (l *LogView) processSampleForConfig(sampling []map[string]interface{}) {
	l.config = config.MakeConfigFromSample(sampling)
	l.app.config = l.config
}

func (l *LogView) makeUIComponents() {
	l.templateView = NewTemplateView(l.app, func() {
		// Toggle full screen func
		l.templateFullScreen = !l.templateFullScreen
		l.makeLayoutsWithTemplateView()
	}, l.makeLayouts)
	l.templateView.SetBorder(true).SetTitle("Template Editor")
	l.data = &LogData{
		logView: l,
	}
	selection := func(row, column int) {
		if row > 0 {
			l.jsonView = NewJsonView(l.app, false,
				func() {
					// Toggle full screen func
					l.logFullScreen = !l.logFullScreen
					l.makeLayoutsWithJsonView()
				}, l.makeLayouts)
			l.jsonView.SetBorder(true).SetTitle("Log Entry")
			b, _ := json.Marshal(l.inSlice[row-1])
			l.jsonView.SetJson(b)
			l.makeLayoutsWithJsonView()
		} else {
			l.makeLayouts()
		}
	}
	l.table = tview.NewTable().
		SetSelectable(true, false).
		SetFixed(1, 1).
		SetSeparator(tview.Borders.Vertical).
		SetContent(l.data)
	l.table.SetSelectedFunc(selection).
		SetBackgroundColor(color.ColorBackgroundField)
	l.table.SetSelectionChangedFunc(func(row, column int) {
		// stop scrolling!
		r, c := l.table.GetOffset()
		l.updateLineView()
		l.table.SetOffset(r, c)
	})

	l.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyF1:
			l.makeLayoutsWithTemplateView()
			return nil
		}
		return event
	})

	l.linesView = tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignRight)
	l.mainMenu = tview.NewFlex().SetDirection(tview.FlexColumn)
	l.mainMenu.
		SetBackgroundColor(color.ColorBackgroundField).SetTitleAlign(tview.AlignCenter)
	l.mainMenu.
		AddItem(tview.NewTextView().
			SetDynamicColors(true).
			SetText("[yellow::b](↲)[-::-] View"), 0, 1, false).
		AddItem(tview.NewTextView().
			SetDynamicColors(true).
			SetText("[yellow::b](↓ ↑ ← →)[-::-] Navigate"), 0, 1, false).
		AddItem(tview.NewTextView().
			SetDynamicColors(true).
			SetText("[yellow::b](g/G)[-::-] Top/Bottom"), 0, 1, false).
		AddItem(tview.NewTextView().
			SetDynamicColors(true).
			SetText("[yellow::b](^f/^b)[-::-] Page Up/Down"), 0, 1, false).
		AddItem(tview.NewTextView().
			SetDynamicColors(true).
			SetText("[yellow::b](F1)[-::-] Template"), 0, 1, false).
		AddItem(tview.NewTextView().
			SetDynamicColors(true).
			SetText("[yellow::b](^C)[-::-] Quit"), 0, 1, false).
		AddItem(l.linesView, 0, 1, false)
	l.updateLineView()
}

func (l *LogView) updateLineView() {
	r, _ := l.table.GetSelection()
	if r > 0 {
		l.linesView.SetText(
			fmt.
				Sprintf(`[yellow::]Line [green::b]%d[yellow::-] ([green::b]%d[yellow::-] lines)`,
					r,
					l.globalCount))
	} else {
		l.linesView.SetText(
			fmt.
				Sprintf(`[green::b]%d[yellow::-] lines`,
					l.globalCount))
	}

}

func (l *LogView) makeLayouts() {
	l.Flex.Clear().SetDirection(tview.FlexRow).
		AddItem(l.table, 0, 2, true).
		AddItem(l.mainMenu, 1, 1, false).
		SetBackgroundColor(color.ColorBackgroundField)
	l.app.SetFocus(l.table)
}

func (l *LogView) makeLayoutsWithJsonView() {
	l.Flex.Clear().SetDirection(tview.FlexRow)
	if !l.logFullScreen {
		l.Flex.AddItem(l.table, 0, 1, false)
	}
	l.Flex.
		AddItem(l.jsonView, 0, 2, false).
		AddItem(l.mainMenu, 1, 1, false)

	l.app.SetFocus(l.jsonView.textView)
}

func (l *LogView) makeLayoutsWithTemplateView() {
	l.Flex.Clear().SetDirection(tview.FlexRow)
	if !l.templateFullScreen {
		l.Flex.AddItem(l.table, 0, 1, false)
	}
	l.templateView.config = l.config
	l.Flex.
		AddItem(l.templateView, 0, 2, false).
		AddItem(l.mainMenu, 1, 1, false)

	l.app.SetFocus(l.templateView.table)
}

type LogData struct {
	tview.TableContentReadOnly
	logView *LogView
}

func (d *LogData) GetCell(row, column int) *tview.TableCell {
	if row == -1 || len(d.logView.inSlice) < row-1 || column == -1 {
		return nil
	}
	c := d.logView.config
	k := c.Keys[column]
	tc := tview.NewTableCell(" " + k.Name + " ")
	if k.MaxWidth > 0 && k.MaxWidth-len(k.Name) >= len(k.Name) {
		spaces := strings.Repeat(" ", k.MaxWidth-len(k.Name))
		tc.SetText(" " + k.Name + spaces)
	}
	// Set Headers
	if row == 0 {
		tc.SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignCenter).
			SetBackgroundColor(tcell.ColorBlack).
			SetSelectable(false)
		return tc
	}
	// Set Body Cells
	cellValue := k.ExtractValue(d.logView.inSlice[row-1])
	var bgColor, fgColor tcell.Color
	if len(k.Color.Foreground) == 0 {
		fgColor = k.Type.GetColor()
	} else {
		fgColor = k.Color.GetForegroundColor()
	}
	bgColor = k.Color.GetBackgroundColor()
	if len(k.ColorWhen) > 0 {
	OUT:
		for _, kv := range k.ColorWhen {
			if strings.ToLower(cellValue) == strings.ToLower(kv.MatchValue) {
				bgColor = kv.Color.GetBackgroundColor()
				fgColor = kv.Color.GetForegroundColor()
				break OUT
			}
		}
	}
	switch k.Type {
	case config.TypeNumber, config.TypeBool:
		tc.SetAlign(tview.AlignRight)
	}
	if k.MaxWidth > 0 {
		tc.MaxWidth = k.MaxWidth
	}

	return tc.
		SetBackgroundColor(bgColor).
		SetTextColor(fgColor).
		SetText(fmt.Sprintf("%s", cellValue))
}

func (d *LogData) GetRowCount() int {
	return len(d.logView.inSlice) + 1
}

func (d *LogData) GetColumnCount() int {
	c := d.logView.config
	return len(c.Keys)
}
