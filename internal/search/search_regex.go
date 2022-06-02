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
