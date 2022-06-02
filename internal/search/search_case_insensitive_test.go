package search

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCaseInsensitiveSearch_Search(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		word  string
		count int
	}{
		{
			name:  "simple text",
			text:  "insert",
			word:  "s",
			count: 1,
		},
		{
			name:  "double text",
			text:  "message",
			word:  "s",
			count: 2,
		},
		{
			name:  "start with word",
			text:  "sam",
			word:  "s",
			count: 1,
		},
		{
			name:  "end with word",
			text:  "seas",
			word:  "s",
			count: 2,
		},
		{
			name:  "url",
			text:  "POST_/api/internal/notification-events",
			word:  "s",
			count: 2,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := MakeCaseInsensitiveSearch(nil)
			idx, err := s.Search(test.word, test.text)
			assert.NoError(t, err)
			assert.Len(t, idx, test.count)
			for _, i := range idx {
				assert.Equal(t, strings.ToLower(test.word),
					strings.ToLower(test.text[i[0]:i[1]]))
			}
			fmt.Println(idx)
		})
	}
}
