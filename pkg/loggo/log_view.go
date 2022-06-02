package loggo

import (
	"encoding/json"
	"fmt"

	"github.com/aurc/loggo/internal/colour"
	"github.com/aurc/loggo/pkg/config"
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
	logFullScreen      bool
	templateFullScreen bool
	inSlice            []map[string]interface{}
}

func NewLogReader(app *LoggoApp, input <-chan string, config *config.Config) *LogView {
	lv := &LogView{
		Flex:   *tview.NewFlex(),
		app:    app,
		config: config,
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
		for {
			t := <-l.input
			if len(t) > 0 {
				m := make(map[string]interface{})
				err := json.Unmarshal([]byte(t), &m)
				if err != nil {
					m[parseErr] = err.Error()
					m[textPayload] = t
				}
				l.inSlice = append(l.inSlice, m)
				l.app.Draw()
			}
		}
	}()
}

func (l *LogView) makeUIComponents() {
	l.templateView = NewTemplateView(l.app, l.config, func() {
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
		SetBackgroundColor(colour.ColourBackgroundField)
	l.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyF1:
			l.makeLayoutsWithTemplateView()
			return nil
		}
		return event
	})

	l.mainMenu = tview.NewFlex().SetDirection(tview.FlexColumn)
	l.mainMenu.
		SetBackgroundColor(colour.ColourBackgroundField).SetTitleAlign(tview.AlignCenter)
	l.mainMenu.
		AddItem(tview.NewTextView().
			SetDynamicColors(true).
			SetText("[yellow::b](ENTER)[-::-] Display selected"), 0, 1, false).
		AddItem(tview.NewTextView().
			SetDynamicColors(true).
			SetText("[yellow::b](↓ ↑ ← →)[-::-] Navigate"), 0, 1, false).
		AddItem(tview.NewTextView().
			SetDynamicColors(true).
			SetText("[yellow::b](F1)[-::-] Edit Template"), 0, 1, false).
		AddItem(tview.NewTextView().
			SetDynamicColors(true).
			SetText("[yellow::b](CTRL-C)[-::-] Quit"), 0, 1, false)
}

func (l *LogView) makeLayouts() {
	l.Flex.Clear().SetDirection(tview.FlexRow).
		AddItem(l.table, 0, 2, true).
		AddItem(l.mainMenu, 1, 1, false).
		SetBackgroundColor(colour.ColourBackgroundField)
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
	tc := tview.NewTableCell(k.Name)
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
	var bgColour, fgColour tcell.Color
	if len(k.Color.Foreground) == 0 {
		fgColour = k.Type.GetColor()
	} else {
		fgColour = k.Color.GetForegroundColor()
	}
	bgColour = k.Color.GetBackgroundColor()
	if len(k.ColorWhen) > 0 {
	OUT:
		for _, kv := range k.ColorWhen {
			if cellValue == kv.MatchValue {
				bgColour = kv.Color.GetBackgroundColor()
				fgColour = kv.Color.GetForegroundColor()
				break OUT
			}
		}
	}
	return tc.
		SetBackgroundColor(bgColour).
		SetTextColor(fgColour).
		SetText(fmt.Sprintf("%s", cellValue))
}

func (d *LogData) GetRowCount() int {
	return len(d.logView.inSlice) + 1
}

func (d *LogData) GetColumnCount() int {
	c := d.logView.config
	return len(c.Keys)
}
