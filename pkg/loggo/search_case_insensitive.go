package loggo

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
