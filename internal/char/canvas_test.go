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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanvas_BlankCanvas(t *testing.T) {
	t.Run("Test Blank Canvas", func(t *testing.T) {
		c := NewCanvas()
		canvas := c.BlankCanvasAsString()
		want := `╔╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╗
╠╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╣
╠╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╣
╠╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╣
╠╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╣
╠╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╣
╠╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╣
╠╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╣
╠╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╣
╚╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╝`
		assert.Equal(t, want, canvas)
	})
}

func TestCanvas_BlankCanvasAsString(t *testing.T) {
	tests := []struct {
		name  string
		words []Char
		wants string
	}{
		{
			name:  "Test l",
			words: []Char{CharacterL},
			wants: `╔╦╦╦╦╦╦╗
╠▓▓▓░╬╬╣
╠╬╬▓▓░╬╣
╠╬╬▓▓░╬╣
╠╬╬▓▓░╬╣
╠╬╬▓▓░╬╣
╠╬╬▓▓░╬╣
╠╬╬▓▓░╬╣
╠╬╬▓▓░╬╣
╠╬╬╬▓▓▓░
╚╩╩╩╩╩╩╝`,
		},
		{
			name:  "Test `",
			words: []Char{CharacterApostrophe},
			wants: `╔▓▓░
╠▓░╣
╠╬╬╣
╠╬╬╣
╠╬╬╣
╠╬╬╣
╠╬╬╣
╠╬╬╣
╠╬╬╣
╠╬╬╣
╚╩╩╝`,
		},
		{
			name:  "Test o",
			words: []Char{CharacterO},
			wants: `╔╦╦╦╦╦╦╦╦╦╦╦╗
╠╬╬╬╬╬╬╬╬╬╬╬╣
╠╬╬╬╬╬╬╬╬╬╬╬╣
╠╬╬╬╬╬╬╬╬╬╬╬╣
╠╬╬╬▓▓▓▓▓░╬╬╣
╠╬▓▓░╬╬╬╬▓▓░╣
╠▓▓░╬╬╬╬╬╬▓▓░
╠▓▓░╬╬╬╬╬╬▓▓░
╠╬▓▓░╬╬╬╬▓▓░╣
╠╬╬╬▓▓▓▓▓░╬╬╣
╚╩╩╩╩╩╩╩╩╩╩╩╝`,
		},
		{
			name:  "Test G",
			words: []Char{CharacterG},
			wants: `╔╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╗
╠╬╬╬╬╬▓▓▓▓▓▓░╬╬╬╣
╠╬╬╬▓▓░╬╬╬╬▓▓▓░╬╣
╠╬╬▓▓░╬╬╬╬╬╬╬╬╬╬╣
╠╬▓▓░╬╬╬╬╬╬╬╬╬╬╬╣
╠▓▓░╬╬╬╬╬╬╬╬╬╬╬╬╣
╠▓▓░╬╬╬╬╬╬╬▓▓▓▓▓░
╠▓▓░╬╬╬╬╬╬╬╬╬▓▓░╣
╠╬╬▓▓░╬╬╬╬╬╬▓▓▓░╣
╠╬╬╬╬▓▓▓▓▓▓▓░╬╬╬╣
╚╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╝`,
		},
		{
			name:  "Test rev G",
			words: []Char{CharacterRevG},
			wants: `╔╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╗
╠╬╬╬╬▓▓▓▓▓▓░╬╬╬╬╣
╠╬╬▓▓▓░╬╬╬╬▓▓░╬╬╣
╠╬╬╬╬╬╬╬╬╬╬╬▓▓░╬╣
╠╬╬╬╬╬╬╬╬╬╬╬╬▓▓░╣
╠╬╬╬╬╬╬╬╬╬╬╬╬╬▓▓░
╠▓▓▓▓▓░╬╬╬╬╬╬╬▓▓░
╠╬▓▓░╬╬╬╬╬╬╬╬╬▓▓░
╠╬▓▓▓░╬╬╬╬╬╬▓▓░╬╣
╠╬╬╬╬▓▓▓▓▓▓▓░╬╬╬╣
╚╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╝`,
		},
		{
			name:  "Test l`oGGo",
			words: LoggoLogo,
			wants: `╔╦╦╦╦╦╦╦╦▓▓░╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╗
╠▓▓▓░╬╬╬╬▓░╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬▓▓▓▓▓▓░╬╬╬╬╬╬╬╬╬╬╬▓▓▓▓▓▓░╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╣
╠╬╬▓▓░╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬▓▓▓░╬╬╬╬▓▓░╬╬╬╬╬╬╬▓▓░╬╬╬╬▓▓▓░╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╣
╠╬╬▓▓░╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬▓▓░╬╬╬╬╬▓▓░╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╣
╠╬╬▓▓░╬╬╬╬╬╬▓▓▓▓▓░╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬▓▓░╬╬╬▓▓░╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬▓▓▓▓▓░╬╬╬╬╬╬╣
╠╬╬▓▓░╬╬╬╬▓▓░╬╬╬╬▓▓░╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬▓▓░╬▓▓░╬╬╬╬╬╬╬╬╬╬╬╬╬╬╬▓▓░╬╬╬╬▓▓░╬╬╬╬╣
╠╬╬▓▓░╬╬╬▓▓░╬╬╬╬╬╬▓▓░╬▓▓▓▓▓░╬╬╬╬╬╬╬▓▓░╬▓▓░╬╬╬╬╬╬╬▓▓▓▓▓░╬▓▓░╬╬╬╬╬╬▓▓░╬╬╬╣
╠╬╬▓▓░╬╬╬▓▓░╬╬╬╬╬╬▓▓░╬╬▓▓░╬╬╬╬╬╬╬╬╬▓▓░╬▓▓░╬╬╬╬╬╬╬╬╬▓▓░╬╬▓▓░╬╬╬╬╬╬▓▓░╬╬╬╣
╠╬╬▓▓░╬╬╬╬▓▓░╬╬╬╬▓▓░╬╬╬▓▓▓░╬╬╬╬╬╬▓▓░╬╬╬╬╬▓▓░╬╬╬╬╬╬▓▓▓░╬╬╬▓▓░╬╬╬╬▓▓░╬╬╬╬╣
╠╬╬╬▓▓▓░╬╬╬╬▓▓▓▓▓░╬╬╬╬╬╬╬╬▓▓▓▓▓▓▓░╬╬╬╬╬╬╬╬╬▓▓▓▓▓▓▓░╬╬╬╬╬╬╬╬▓▓▓▓▓░╬╬╬╬╬╬╣
╚╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╝`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := NewCanvas().WithWord(test.words...)
			str := c.PrintCanvasAsString()
			fmt.Println(str)
			assert.Equal(t, test.wants, str)
		})
	}
}

func TestCanvas_PrintCanvasAsHtml(t *testing.T) {
	c := NewCanvas().WithWord(LoggoLogo...)
	str := c.PrintCanvasAsHtml()
	fmt.Println(str)
}
