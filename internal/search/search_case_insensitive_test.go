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
