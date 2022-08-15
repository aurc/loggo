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

package char

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

type Canvas struct {
	Width        int
	Height       int
	PaintChar    rune
	Word         []Char
	CanvasBorder CanvasBorder
}

type CanvasBorder struct {
	TopBorderChar         rune
	BottomBorderChar      rune
	LeftBorderChar        rune
	RightBorderChar       rune
	TopLeftCornerChar     rune
	TopRightCornerChar    rune
	BottomLeftCornerChar  rune
	BottomRightCornerChar rune
}

func NewCanvas() *Canvas {
	c := &Canvas{}
	return c.WithPaintChar('╬', CanvasBorder{
		TopBorderChar:         '╦',
		BottomBorderChar:      '╩',
		LeftBorderChar:        '╠',
		RightBorderChar:       '╣',
		TopLeftCornerChar:     '╔',
		TopRightCornerChar:    '╗',
		BottomLeftCornerChar:  '╚',
		BottomRightCornerChar: '╝',
	}).WithDimensions(40, 10)
}

func (c *Canvas) WithPaintChar(r rune, border CanvasBorder) *Canvas {
	c.PaintChar = r
	c.CanvasBorder = border
	return c
}

func (c *Canvas) WithDimensions(width, height int) *Canvas {
	c.Width = width
	c.Height = height
	return c
}

func (c *Canvas) WithWord(word ...Char) *Canvas {
	c.Word = word
	// Recalculate canvas size
	if c.Height < 11 {
		c.Height = 11
	}
	width := 0
	for _, ch := range word {
		width += ch.GetWidth()
	}
	c.Width = width
	return c
}

func (c *Canvas) BlankCanvas() [][]rune {
	topRow := make([]rune, 1)
	topRow[0] = c.CanvasBorder.TopLeftCornerChar
	topRow = append(topRow, []rune(strings.Repeat(string(c.CanvasBorder.TopBorderChar), c.Width-2))...)
	topRow = append(topRow, c.CanvasBorder.TopRightCornerChar)

	canvas := [][]rune{topRow}
	for i := 1; i < c.Height-1; i++ {
		middleRow := make([]rune, 1)
		middleRow[0] = c.CanvasBorder.LeftBorderChar
		middleRow = append(middleRow, []rune(strings.Repeat(string(c.PaintChar), c.Width-2))...)
		middleRow = append(middleRow, c.CanvasBorder.RightBorderChar)
		canvas = append(canvas, middleRow)
	}

	bottomRow := make([]rune, 1)
	bottomRow[0] = c.CanvasBorder.BottomLeftCornerChar
	bottomRow = append(bottomRow, []rune(strings.Repeat(string(c.CanvasBorder.BottomBorderChar), c.Width-2))...)
	bottomRow = append(bottomRow, c.CanvasBorder.BottomRightCornerChar)

	canvas = append(canvas, bottomRow)

	return canvas
}

func (c *Canvas) BlankCanvasAsString() string {
	return c.toString(c.BlankCanvas())
}

func (c *Canvas) PrintCanvas() [][]rune {
	bc := c.BlankCanvas()
	currWidth := 1
	for _, w := range c.Word {
		for _, coord := range w.Coordinates {
			for i := 0; i < coord.L; i++ {
				x := coord.X + i + currWidth
				y := coord.Y
				bc[y][x] = w.PaintChar
				bc[y][x+1] = w.Shade
			}
		}
		currWidth = currWidth + w.Next
	}
	return bc
}

func (c *Canvas) PrintCanvasAsHtml() string {
	str := c.PrintCanvasAsString()
	buf := bytes.NewBufferString(str)
	reader := bufio.NewReader(buf)
	builder := strings.Builder{}
	convMap := map[rune]string{
		'▓': "&blk34;",
		'░': "&blk14;",
		'╬': "&boxVH;",
		'╦': "&boxHD;",
		'╩': "&boxHU;",
		'╠': "&boxVR;",
		'╣': "&boxVL;",
		'╔': "&boxDR;",
		'╗': "&boxDL;",
		'╚': "&boxUR;",
		'╝': "&boxUL;",
	}
	paintChar := '▓'
	shade := '░'
	for {
		str, err := reader.ReadString('\n')
		if err == nil {
			for _, char := range str {
				switch char {
				case paintChar, shade:
					builder.WriteString(fmt.Sprintf(`<span class="fgCol">%s</span>`, convMap[char]))
				default:
					builder.WriteString(fmt.Sprintf(`<span class="bgCol">%s</span>`, convMap[char]))
				}
			}
			builder.WriteString("<br>\n")
		} else {
			break
		}
	}
	return builder.String()
}

func (c *Canvas) PrintCanvasAsString() string {
	return c.toString(c.PrintCanvas())
}

func (c *Canvas) toString(rc [][]rune) string {
	sb := strings.Builder{}
	for i, row := range rc {
		sb.WriteString(string(row))
		if i < len(rc)-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}
