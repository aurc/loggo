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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aurc/loggo/internal/config"
)

func TestFilterGroup_Resolve(t *testing.T) {
	tests := []struct {
		name        string
		whenJsonRow string
		givenFilter Group
		wantsResult bool
		wantsError  bool
	}{
		{
			name: `Given ((a/b = 'x' OR a/b = 'y') AND (c between 1 AND 3 OR c > 5)) with a/b = 'y', c = 2, wants true`,
			whenJsonRow: `
					{
						"a": {
							"b": "y"
						},
						"c": "2"
					}`,
			givenFilter: And(
				OrFilters(
					Equals(&config.Key{
						Name: "a/b",
						Type: config.TypeString,
					}, "x"),
					Equals(&config.Key{
						Name: "a/b",
						Type: config.TypeString,
					}, "y")),
				OrFilters(
					Between(&config.Key{
						Name: "c",
						Type: config.TypeNumber,
					}, "1", "3"),
					GreaterThan(&config.Key{
						Name: "c",
						Type: config.TypeNumber,
					}, "5"))),
			wantsResult: true,
		},
		{
			name: `Given ((a/b = 'x' OR a/b = 'y') AND (c between 1 AND 3 OR c > 5)) with a/b = 'x', c = 7, wants true`,
			whenJsonRow: `
					{
						"a": {
							"b": "x"
						},
						"c": "7"
					}`,
			givenFilter: And(
				OrFilters(
					Equals(&config.Key{
						Name: "a/b",
						Type: config.TypeString,
					}, "x"),
					Equals(&config.Key{
						Name: "a/b",
						Type: config.TypeString,
					}, "y")),
				OrFilters(
					Between(&config.Key{
						Name: "c",
						Type: config.TypeNumber,
					}, "1", "3"),
					GreaterThan(&config.Key{
						Name: "c",
						Type: config.TypeNumber,
					}, "5"))),
			wantsResult: true,
		},
		{
			name: `Given ((a/b = 'x' OR a/b = 'y') AND (c between 1 AND 3 OR c > 5)) with a/b = 'n', c = 7, wants false`,
			whenJsonRow: `
					{
						"a": {
							"b": "n"
						},
						"c": "7"
					}`,
			givenFilter: And(
				OrFilters(
					Equals(&config.Key{
						Name: "a/b",
						Type: config.TypeString,
					}, "x"),
					Equals(&config.Key{
						Name: "a/b",
						Type: config.TypeString,
					}, "y")),
				OrFilters(
					Between(&config.Key{
						Name: "c",
						Type: config.TypeNumber,
					}, "1", "3"),
					GreaterThan(&config.Key{
						Name: "c",
						Type: config.TypeNumber,
					}, "5"))),
			wantsResult: false,
		},
		{
			name: `Given ((a/b = 'x' OR a/b = 'y') AND (c between 1 AND 3 OR c > 5)) with a/b = 'x', c = 3, wants false`,
			whenJsonRow: `
					{
						"a": {
							"b": "x"
						},
						"c": "3"
					}`,
			givenFilter: And(
				OrFilters(
					Equals(&config.Key{
						Name: "a/b",
						Type: config.TypeString,
					}, "x"),
					Equals(&config.Key{
						Name: "a/b",
						Type: config.TypeString,
					}, "y")),
				OrFilters(
					Between(&config.Key{
						Name: "c",
						Type: config.TypeNumber,
					}, "1", "3"),
					GreaterThan(&config.Key{
						Name: "c",
						Type: config.TypeNumber,
					}, "5"))),
			wantsResult: false,
		},
		{
			name: `Given ((a/b = 'x' OR a/b = 'y') OR (c between 1 AND 3 OR c > 5)) with a/b = 'x', c = 3, wants true`,
			whenJsonRow: `
					{
						"a": {
							"b": "x"
						},
						"c": "3"
					}`,
			givenFilter: Or(
				OrFilters(
					Equals(&config.Key{
						Name: "a/b",
						Type: config.TypeString,
					}, "x"),
					Equals(&config.Key{
						Name: "a/b",
						Type: config.TypeString,
					}, "y")),
				OrFilters(
					Between(&config.Key{
						Name: "c",
						Type: config.TypeNumber,
					}, "1", "3"),
					GreaterThan(&config.Key{
						Name: "c",
						Type: config.TypeNumber,
					}, "5"))),
			wantsResult: true,
		},
		{
			name: `Given ((a/b = 'x' OR a/b = 'y') OR (c between 1 AND 3 AND r < 5)) with a/b = 'x', c = 3, r = 5, wants true`,
			whenJsonRow: `
					{
						"a": {
							"b": "x"
						},
						"c": "2",
						"r": "4"
					}`,
			givenFilter: Or(
				OrFilters(
					Equals(&config.Key{
						Name: "a/b",
						Type: config.TypeString,
					}, "x"),
					Equals(&config.Key{
						Name: "a/b",
						Type: config.TypeString,
					}, "y")),
				AndFilters(
					Between(&config.Key{
						Name: "c",
						Type: config.TypeNumber,
					}, "1", "3"),
					LowerThan(&config.Key{
						Name: "r",
						Type: config.TypeNumber,
					}, "5"))),
			wantsResult: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var row map[string]interface{}
			err := json.Unmarshal([]byte(test.whenJsonRow), &row)
			assert.NoError(t, err)
			result, err := test.givenFilter.Resolve(row)
			if test.wantsError {
				assert.NotNil(t, err)
				assert.Error(t, err)
			} else {
				assert.Equal(t, test.wantsResult, result)
			}
		})
	}
}
