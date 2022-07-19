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

	"github.com/aurc/loggo/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestParseFilterExpression(t *testing.T) {
	tests := []struct {
		name            string
		whenJsonRow     string
		givenExpression string
		keySet          map[string]*config.Key
		wantsResult     bool
		wantsError      bool
	}{
		{
			name: `wants true - 2 between 1 and 3`,
			whenJsonRow: `
					{
						"a": {
							"b": "y"
						},
						"c": "2"
					}`,
			keySet: map[string]*config.Key{
				"a/b": {
					Name: "a/b",
					Type: config.TypeString,
				},
				"c": {
					Name: "c",
					Type: config.TypeNumber,
				},
			},
			givenExpression: `((a/b = "x" OR a/b = "y") AND (c between 1 AND 3 OR c > 5))`,
			wantsResult:     true,
		},
		{
			name: `wants true - 7 is greater than 5`,
			whenJsonRow: `
					{
						"a": {
							"b": "x"
						},
						"c": "7"
					}`,
			keySet: map[string]*config.Key{
				"a/b": {
					Name: "a/b",
					Type: config.TypeString,
				},
				"c": {
					Name: "c",
					Type: config.TypeNumber,
				},
			},
			givenExpression: `((a/b = "x" OR a/b = "y") AND (c between 1 AND 3 OR c > 5))`,
			wantsResult:     true,
		},
		{
			name: `wants false - b is not in range`,
			whenJsonRow: `
					{
						"a": {
							"b": "n"
						},
						"c": "7"
					}`,
			givenExpression: `((a/b = "x" OR a/b = "y") AND (c between 1 AND 3 OR c > 5))`,
			keySet: map[string]*config.Key{
				"a/b": {
					Name: "a/b",
					Type: config.TypeString,
				},
				"c": {
					Name: "c",
					Type: config.TypeNumber,
				},
			},
			wantsResult: false,
		},
		{
			name: `wants true - between inclusive of 3`,
			whenJsonRow: `
					{
						"a": {
							"b": "x"
						},
						"c": "3"
					}`,
			givenExpression: `((a/b == "x" OR a/b = "y") AND (c between 1 AND 3 OR c > 5))`,
			keySet: map[string]*config.Key{
				"a/b": {
					Name: "a/b",
					Type: config.TypeString,
				},
				"c": {
					Name: "c",
					Type: config.TypeNumber,
				},
			},
			wantsResult: true,
		},
		{
			name: `wants false - group items resolve to false`,
			whenJsonRow: `
					{
						"a": {
							"b": "x"
						},
						"c": "4"
					}`,
			givenExpression: `a/b = "y" OR (a/b = "x" AND (c between 1 AND 3 OR c > 5))`,
			keySet: map[string]*config.Key{
				"a/b": {
					Name: "a/b",
					Type: config.TypeString,
				},
				"c": {
					Name: "c",
					Type: config.TypeNumber,
				},
			},
			wantsResult: false,
		},
		{
			name: `wants true - all groups resolve to true`,
			whenJsonRow: `
					{
						"a": {
							"b": "X"
						},
						"c": "2",
						"r": "4"
					}`,
			givenExpression: `((a/b = "x" OR a/b = "y") AND (c between 1 AND 3 AND r < 5))`,
			keySet: map[string]*config.Key{
				"a/b": {
					Name: "a/b",
					Type: config.TypeString,
				},
				"c": {
					Name: "c",
					Type: config.TypeNumber,
				},
				"r": {
					Name: "r",
					Type: config.TypeNumber,
				},
			},
			wantsResult: true,
		},
		{
			name: `wants true - when bool and contains`,
			whenJsonRow: `
					{
						"b": "true",
						"s": "banana"
					}`,
			givenExpression: `b = 'true' and s contains "ana"`,
			keySet: map[string]*config.Key{
				"b": {
					Name: "b",
					Type: config.TypeBool,
				},
				"s": {
					Name: "s",
					Type: config.TypeString,
				},
			},
			wantsResult: true,
		},
		{
			name: `wants false - contains does not match`,
			whenJsonRow: `
					{
						"b": "true",
						"s": "banana"
					}`,
			givenExpression: `b = 'true' and s contains "aNa"`,
			keySet: map[string]*config.Key{
				"b": {
					Name: "b",
					Type: config.TypeBool,
				},
				"s": {
					Name: "s",
					Type: config.TypeString,
				},
			},
			wantsResult: false,
		},
		{
			name: `wants false - contains ignore case`,
			whenJsonRow: `
					{
						"b": "true",
						"s": "banana"
					}`,
			givenExpression: `b = 'true' and s containsIC "aNa"`,
			keySet: map[string]*config.Key{
				"b": {
					Name: "b",
					Type: config.TypeBool,
				},
				"s": {
					Name: "s",
					Type: config.TypeString,
				},
			},
			wantsResult: true,
		},
		{
			name: `wants true - numb and regex`,
			whenJsonRow: `
					{
						"b": 1,
						"s": "abb333"
					}`,
			givenExpression: `b = 1 and s match "[a-z]+[0-9]+"`,
			keySet: map[string]*config.Key{
				"b": {
					Name: "b",
					Type: config.TypeNumber,
				},
				"s": {
					Name: "s",
					Type: config.TypeString,
				},
			},
			wantsResult: true,
		},
		{
			name: `wants true - not equals and missing key`,
			whenJsonRow: `
					{
						"b": "b val",
						"s": "some"
					}`,
			givenExpression: `b != "c val" and s == "some"`,
			keySet: map[string]*config.Key{
				"b": {
					Name: "b",
					Type: config.TypeString,
				},
			},
			wantsResult: true,
		},
		{
			name: `wants true - between if lower-greater than`,
			whenJsonRow: `
					{
						"a": 2
					}`,
			givenExpression: `a >= 1 and a <=3`,
			keySet: map[string]*config.Key{
				"a": {
					Name: "a",
					Type: config.TypeNumber,
				},
			},
			wantsResult: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var row map[string]interface{}
			err := json.Unmarshal([]byte(test.whenJsonRow), &row)
			assert.NoError(t, err)
			exp, err := ParseFilterExpression(test.givenExpression)
			assert.NoError(t, err)
			result, err := exp.Apply(row, test.keySet)

			if test.wantsError {
				assert.NotNil(t, err)
				assert.Error(t, err)
			} else {
				assert.Equal(t, test.wantsResult, result)
			}
		})
	}
}
