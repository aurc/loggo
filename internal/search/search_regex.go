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
	"regexp"
	"strings"

	"github.com/rivo/tview"
)

type regexSearch struct {
	search
	regex *regexp.Regexp
}

func MakeRegexSearch(statusBar *tview.TextView) Searchable {
	s := &regexSearch{}
	s.searchStrategy = s
	s.search.statusBar = statusBar
	s.search.Clear()
	return s
}

func (c *regexSearch) Search(word, text string) ([][]int, error) {
	_, _ = c.search.Search(word, text)
	var err error
	c.regex, err = regexp.Compile(word)
	if err != nil {
		c.search.selectionCount = -1
		c.setErrorStatus(err)
		return nil, err
	}
	c.startIndexes = c.regex.FindAllIndex([]byte(text), -1)
	text = strings.ToLower(text)

	return c.startIndexes, nil
}
