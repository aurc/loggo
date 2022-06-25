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

package search

import (
	"strings"

	"github.com/rivo/tview"
)

type caseInsensitiveSearch struct {
	search
}

func MakeCaseInsensitiveSearch(statusBar *tview.TextView) Searchable {
	s := &caseInsensitiveSearch{}
	s.searchStrategy = s
	s.search.statusBar = statusBar
	s.search.Clear()
	return s
}

func (c *caseInsensitiveSearch) Search(word, text string) ([][]int, error) {
	_, _ = c.search.Search(word, text)
	word = strings.ToLower(word)
	c.startIndexes = [][]int{}
	text = strings.ToLower(text)
	c.searchAll(word, text, 0)
	return c.startIndexes, nil
}

func (c *caseInsensitiveSearch) searchAll(word, text string, currPointerIdx int) {
	idx := strings.Index(text, word)
	if idx > -1 {
		newPointerIdx := currPointerIdx + idx
		c.startIndexes = append(c.startIndexes, []int{newPointerIdx, newPointerIdx + len(word)})
		c.searchAll(word, text[idx+len(word):], newPointerIdx+len(word))
	}
}
