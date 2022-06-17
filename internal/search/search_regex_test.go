/*
Copyright © 2022 Aurelio Calegari

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
