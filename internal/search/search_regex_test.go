package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegexSearch_Search(t *testing.T) {
	tests := []struct {
		name       string
		text       string
		word       string
		wants      string
		wantsError bool
	}{
		{
			name:  "simple text",
			text:  "insert",
			word:  `.+s`,
			wants: "ins",
		},
		{
			name:  "double text",
			text:  "message",
			word:  `s+`,
			wants: "ss",
		},
		{
			name:  "url",
			text:  "POST_/api/internal/notification-events",
			word:  `/[a-z]+/`,
			wants: "/api/",
		},
		{
			name:       "bad pattern",
			text:       "POST_/api/internal/notification-events",
			word:       `\`,
			wantsError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := MakeRegexSearch(nil)
			idx, err := s.Search(test.word, test.text)
			if test.wantsError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.wants, test.text[idx[0][0]:idx[0][1]])
			}
		})
	}
}
