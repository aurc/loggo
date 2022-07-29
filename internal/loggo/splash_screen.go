/*
Copyright © 2022 Aurelio Calegari, et al.

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
	"strings"
	"time"

	"github.com/aurc/loggo/internal/char"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type SplashScreen struct {
	tview.Flex
	app          Loggo
	titleView    *tview.TextView
	subtitleView *tview.TextView
	canvas       [][]rune
}

func NewSplashScreen(app Loggo) *SplashScreen {
	tv := &SplashScreen{
		Flex: *tview.NewFlex(),
		app:  app,
	}
	tv.makeUIComponents()
	tv.renderLogo()
	tv.makeLayouts()
	return tv
}

func (t *SplashScreen) makeUIComponents() {
	t.Flex.SetBackgroundColor(tcell.ColorBlack)
	c := char.NewCanvas().WithWord(char.LoggoLogo...).WithDimensions(69, 11)
	t.canvas = c.PrintCanvas()
	t.titleView = tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignCenter)
	t.subtitleView = tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignCenter)
	t.subtitleView.SetText(fmt.Sprintf(`
[white:black:b]l'oGGo %s[::-]: [yellow::u]Rich Terminal User Interface for following JSON logs
[gray::-]Copyright © 2022 Aurelio Calegari, et al.
[lightgray::u]https://github.com/aurc/loggo
`, BuildVersion)).SetBackgroundColor(tcell.ColorBlack)
}

func (t *SplashScreen) renderLogo() {
	steps := []struct {
		bg string
		fg string
		sh string
		tx string
	}{
		{"#000000", "#000000", "#000000", "#000000"},
		{"#000000", "#000203", "#000203", "#000406"},
		{"#000000", "#000405", "#000405", "#00090b"},
		{"#000000", "#000708", "#000708", "#000d11"},
		{"#000000", "#00090b", "#00090b", "#001216"},
		{"#000000", "#000b0e", "#000b0e", "#00161c"},
		{"#000000", "#000d10", "#000d10", "#001a21"},
		{"#000000", "#000f13", "#000f13", "#001f27"},
		{"#000000", "#001116", "#001116", "#00232c"},
		{"#000000", "#001419", "#001419", "#002832"},
		{"#000000", "#00161b", "#00161b", "#002c37"},
		{"#000000", "#00181e", "#00181e", "#00313d"},
		{"#000000", "#001a21", "#001a21", "#003542"},
		{"#000000", "#001c24", "#001c24", "#003948"},
		{"#000000", "#001f26", "#001f26", "#003e4d"},
		{"#000000", "#002129", "#002129", "#004253"},
		{"#000000", "#00232c", "#00232c", "#004758"},
		{"#000000", "#00252f", "#00252f", "#004b5e"},
		{"#000000", "#002731", "#002731", "#004f63"},
		{"#000000", "#002934", "#002934", "#005469"},
		{"#000000", "#002c37", "#002c37", "#00586e"},
		{"#000000", "#002e3a", "#002e3a", "#005d74"},
		{"#000000", "#00303c", "#00303c", "#006179"},
		{"#000000", "#00323f", "#00323f", "#00657f"},
		{"#000000", "#003442", "#003442", "#006a84"},
		{"#000000", "#003645", "#003645", "#006e8a"},
		{"#000000", "#003947", "#003947", "#00738f"},
		{"#000000", "#003b4a", "#003b4a", "#007795"},
		{"#000000", "#003d4d", "#003d4d", "#007b9a"},
		{"#000000", "#003f50", "#003f50", "#0080a0"},
		{"#000000", "#004152", "#004152", "#0084a5"},
		{"#000000", "#004455", "#004455", "#0089ab"},
		{"#000000", "#004658", "#004658", "#008db0"},
		{"#000000", "#00485b", "#00485b", "#0092b6"},
		{"#000000", "#004a5d", "#004a5d", "#0096bb"},
		{"#000000", "#004c60", "#004c60", "#009ac1"},
		{"#000000", "#004e63", "#004e63", "#009fc6"},
		{"#000000", "#005166", "#005166", "#00a3cc"},
		{"#000000", "#005368", "#005368", "#00a8d1"},
		{"#000000", "#00556b", "#00556b", "#00acd7"},
	}
	go func() {
		shColor := ""
		txColor := ""
		for _, s := range steps {
			bgColor := fmt.Sprintf(`[%s:%s]`, s.fg, s.bg)
			txColor = fmt.Sprintf(`[%s:%s]`, s.tx, s.bg)
			shColor = fmt.Sprintf(`[%s:%s]`, s.sh, s.bg)
			text := t.PrintCanvasAsColorString('▓', '░', txColor, shColor, bgColor)
			t.titleView.SetText(text)
			time.Sleep(15 * time.Millisecond)
			t.app.Draw()
		}
		for i := len(steps) - 1; i >= 0; i-- {
			s := steps[i]
			bgColor := fmt.Sprintf(`[%s:%s]`, s.fg, s.bg)
			//txColor := fmt.Sprintf(`[%s:%s]`, s.tx, s.bg)
			//shColor := fmt.Sprintf(`[%s:%s]`, s.sh, s.bg)
			text := t.PrintCanvasAsColorString('▓', '░', txColor, shColor, bgColor)
			t.titleView.SetText(text)
			time.Sleep(25 * time.Millisecond)
			t.app.Draw()
		}
	}()
}

func (t *SplashScreen) makeLayouts() {
	t.Flex.Clear().SetDirection(tview.FlexRow).
		AddItem(t.titleView, 10, 1, false).
		AddItem(t.subtitleView, 0, 1, false)
}

func (t *SplashScreen) PrintCanvasAsColorString(foreground, shade rune, foregroundColor, shadeColor, backgroundColor string) string {
	sb := strings.Builder{}
	const (
		fg = "fg"
		bg = "bg"
		sh = "sh"
	)
	prev := ""
	for row := 0; row < len(t.canvas); row++ {
		for col := 0; col < len(t.canvas[row]); col++ {
			switch t.canvas[row][col] {
			case foreground:
				if prev != fg {
					sb.WriteString(fmt.Sprintf(foregroundColor))
					prev = fg
				}
			case shade:
				if prev != sh {
					sb.WriteString(fmt.Sprintf(shadeColor))
					prev = sh
				}
			default:
				if prev != bg {
					sb.WriteString(fmt.Sprintf(backgroundColor))
					prev = bg
				}
			}
			sb.WriteString(string(t.canvas[row][col]))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
