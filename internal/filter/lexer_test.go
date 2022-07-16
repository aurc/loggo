/*
Copyright © 2022 Aurelio Calegari, et al.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software AND associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, AND/OR sell
copies of the Software, AND to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice AND this permission notice shall be included in
all copies OR substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package filter

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFilterExpression(t *testing.T) {
	tests := []struct {
		name       string
		expression string
	}{
		{
			name:       "Test",
			expression: `a = 3 or b = 4`,
		},
		{
			name:       "Test",
			expression: `(a = 3 or b between 1 and 2) and c = "r"`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fg, err := ParseFilterExpression(test.expression)
			assert.NoError(t, err)
			b, err := json.MarshalIndent(fg, "", "  ")
			assert.NoError(t, err)
			fmt.Println(string(b))
		})
	}
}
