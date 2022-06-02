package search

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
)

type search struct {
	startIndexes   [][]int
	selectionCount int
	searchWordIdx  int
	statusBar      *tview.TextView
	searchStrategy Searchable
}

type Searchable interface {
	Search(word, text string) ([][]int, error)
	TagWord(withTag, val string) string
	SetCurrentStatus()
	Clear()
	GetSearchCount() int
	GetSearchPosition() int
	Next() int
	Prev() int
}

func (s *search) Clear() {
	s.selectionCount = 0
	s.searchWordIdx = 0
	if s.statusBar != nil {
		s.statusBar.Clear()
	}
}

func (s *search) Search(word, text string) ([][]int, error) {
	if s.statusBar != nil {
		s.statusBar.Clear()
	}
	return nil, nil
}

func (s *search) Next() int {
	if s.searchWordIdx >= s.selectionCount-1 {
		s.searchWordIdx = 0
	} else {
		s.searchWordIdx++
	}
	return s.searchWordIdx
}
func (s *search) Prev() int {
	if s.searchWordIdx == 0 {
		s.searchWordIdx = s.selectionCount - 1
	} else {
		s.searchWordIdx--
	}
	return s.searchWordIdx
}

func (s *search) GetSearchPosition() int {
	return s.searchWordIdx + 1
}

func (s *search) GetSearchCount() int {
	return s.selectionCount
}

func (s *search) TagWord(withTag, val string) string {
	idxs, _ := s.searchStrategy.Search(withTag, val)
	if len(idxs) > 0 {
		taggedWord := strings.Builder{}
		preIdx := 0
		for i, idx := range idxs {
			if i == 0 {
				taggedWord.WriteString(val[0:idx[0]])
			} else {
				taggedWord.WriteString(val[preIdx:idx[0]])
			}
			tagID := fmt.Sprintf("%d", s.selectionCount)
			s.selectionCount++
			taggedWord.WriteString(
				fmt.Sprintf(
					`[:brown:]["%s"]%v[""][:-:]`,
					tagID,
					val[idx[0]:idx[1]]))
			preIdx = idx[1]
		}
		taggedWord.WriteString(val[preIdx:])
		return taggedWord.String()
	}
	return ""
}

func (s *search) SetCurrentStatus() {
	if s.selectionCount == 0 {
		s.resetStatus(`[yellow]No results returned`)
	} else if s.selectionCount > 0 {
		s.resetStatus(
			fmt.Sprintf(`[white]Showing result [green::b]%d[white:-:-] out of [green::b]%d[white:-:-]`,
				s.searchWordIdx+1, s.selectionCount))
	}
}

func (s *search) setErrorStatus(err error) {
	s.resetStatus(fmt.Sprintf(`[red::b]%s`, err.Error()))
}

func (s *search) resetStatus(text string) {
	if s.statusBar != nil {
		s.statusBar.Clear().SetTextAlign(tview.AlignCenter).
			SetDynamicColors(true).
			SetText(text)
	}
}
