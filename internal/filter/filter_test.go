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
	"testing"

	"github.com/aurc/loggo/internal/config"
	"github.com/stretchr/testify/assert"
)

type testFilter struct {
	name        string
	filter      Filter
	whenValue   string
	shouldMatch bool
	wantError   bool
}

var keySet = map[string]*config.Key{
	"strName": {
		Name: "strName",
		Type: config.TypeString,
	},
	"boolKey": {
		Name: "boolKey",
		Type: config.TypeBool,
	},
	"numbKey": {
		Name: "numbKey",
		Type: config.TypeNumber,
	},
	"dateTimeKey": {
		Name:   "abc",
		Type:   config.TypeDateTime,
		Layout: "2006-01-02T15:04:05-0700",
	},
}

func TestEqual_Apply(t *testing.T) {
	tests := []testFilter{
		{
			name:        "Wants exact STRING match",
			filter:      Equals("strName", "minion"),
			whenValue:   "minion",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name:        "No STRING match",
			filter:      Equals("strName", "min"),
			whenValue:   "minion",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name:        "Wants exact BOOL match",
			filter:      Equals("strName", "true"),
			whenValue:   "true",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name:        "No BOOL match",
			filter:      Equals("boolKey", "true"),
			whenValue:   "bubbles",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name:        "Wants BAD BOOL on value",
			filter:      Equals("boolKey", "false"),
			whenValue:   "bananas",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name:        "Wants BAD BOOL on expression",
			filter:      Equals("boolKey", "bananas"),
			whenValue:   "false",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name:        "Wants exact NUMBER match",
			filter:      Equals("numbKey", "0.01"),
			whenValue:   "0.01",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name:        "No NUMBER match",
			filter:      Equals("numbKey", "0.0109"),
			whenValue:   "0.01",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name:        "Wants BAD number on value",
			filter:      Equals("numbKey", "0.0109"),
			whenValue:   "bananas",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name:        "Wants BAD number on expression",
			filter:      Equals("numbKey", "bananas"),
			whenValue:   "10",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name:        "Wants exact DATE match",
			filter:      Equals("dateTimeKey", "2006-01-02T15:04:05-0700"),
			whenValue:   "2006-01-02T15:04:05-0700",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name:        "No DATE match",
			filter:      Equals("dateTimeKey", "2006-01-02T15:04:05-0710"),
			whenValue:   "2006-01-02T15:04:05-0700",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name:        "Wants BAD DATE value",
			filter:      Equals("dateTimeKey", "2006-01-02T15:04:05-0700"),
			whenValue:   "bananas",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name:        "Wants BAD DATE expression",
			filter:      Equals("dateTimeKey", "bananas"),
			whenValue:   "2006-01-02T15:04:05-0700",
			shouldMatch: false,
			wantError:   true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testFilterFunc(t, test)
		})
	}
}

func TestNotEqual_Apply(t *testing.T) {
	tests := []testFilter{
		{
			name:        "Wants exact STRING match",
			filter:      NotEquals("strName", "minion"),
			whenValue:   "minion",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name:        "No STRING match",
			filter:      NotEquals("strName", "min"),
			whenValue:   "minion",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name:        "Wants exact BOOL match",
			filter:      NotEquals("strName", "true"),
			whenValue:   "true",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name:        "No BOOL match",
			filter:      NotEquals("boolKey", "true"),
			whenValue:   "bubbles",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name:        "Wants BAD BOOL on value",
			filter:      Equals("boolKey", "false"),
			whenValue:   "bananas",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name:        "Wants BAD BOOL on expression",
			filter:      Equals("boolKey", "bananas"),
			whenValue:   "false",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name:        "Wants exact NUMBER match",
			filter:      NotEquals("numbKey", "0.01"),
			whenValue:   "0.01",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name:        "No NUMBER match",
			filter:      NotEquals("numbKey", "0.0109"),
			whenValue:   "0.01",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name:        "Wants BAD number on value",
			filter:      Equals("numbKey", "0.0109"),
			whenValue:   "bananas",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name:        "Wants BAD number on expression",
			filter:      Equals("numbKey", "bananas"),
			whenValue:   "10",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name:        "Wants exact DATE match",
			filter:      NotEquals("dateTimeKey", "2006-01-02T15:04:05-0700"),
			whenValue:   "2006-01-02T15:04:05-0700",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name:        "No DATE match",
			filter:      NotEquals("dateTimeKey", "2006-01-02T15:04:05-0710"),
			whenValue:   "2006-01-02T15:04:05-0700",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name:        "Wants BAD DATE value",
			filter:      Equals("dateTimeKey", "2006-01-02T15:04:05-0700"),
			whenValue:   "bananas",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name:        "Wants BAD DATE expression",
			filter:      Equals("dateTimeKey", "bananas"),
			whenValue:   "2006-01-02T15:04:05-0700",
			shouldMatch: false,
			wantError:   true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testFilterFunc(t, test)
		})
	}
}

func TestMatchRegex_Apply(t *testing.T) {
	tests := []testFilter{
		{
			name:        "Wants exact STRING match",
			filter:      MatchesRegex("strName", `\d+[a-zA-Z]+`),
			whenValue:   "123LoGGo",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name:        "No STRING match",
			filter:      MatchesRegex("strName", "onion"),
			whenValue:   "minion",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name:        "BAD Regex",
			filter:      MatchesRegex("strName", `\`),
			whenValue:   "minion",
			shouldMatch: false,
			wantError:   true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testFilterFunc(t, test)
		})
	}
}

func TestContains_Apply(t *testing.T) {
	tests := []testFilter{
		{
			name:        "Wants exact STRING match",
			filter:      Contains("strName", "io"),
			whenValue:   "minion",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name:        "No STRING match",
			filter:      Contains("strName", "onion"),
			whenValue:   "minion",
			shouldMatch: false,
			wantError:   false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testFilterFunc(t, test)
		})
	}
}

func TestEqualsIgnoreCase_Apply(t *testing.T) {
	tests := []testFilter{
		{
			name:        "Wants exact STRING match",
			filter:      EqualIgnoreCase("strName", "minion"),
			whenValue:   "mInioN",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name:        "No STRING match",
			filter:      EqualIgnoreCase("strName", "mInio"),
			whenValue:   "minion",
			shouldMatch: false,
			wantError:   false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testFilterFunc(t, test)
		})
	}
}

func TestContainsIgnoreCase_Apply(t *testing.T) {
	tests := []testFilter{
		{
			name:        "Wants exact STRING match",
			filter:      ContainsIgnoreCase("strName", "minion"),
			whenValue:   "miNion",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name:        "Wants contains STRING match",
			filter:      ContainsIgnoreCase("strName", "mIN"),
			whenValue:   "minion",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name:        "No STRING match",
			filter:      Equals("strName", "m1N"),
			whenValue:   "minion",
			shouldMatch: false,
			wantError:   false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testFilterFunc(t, test)
		})
	}
}

func TestBetween(t *testing.T) {
	tests := []testFilter{
		{
			name:        "Wants exact STRING match",
			filter:      Between("strName", "minion", "zorg"),
			whenValue:   "onion",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name:        "No STRING match",
			filter:      Between("strName", "minion", "zorg"),
			whenValue:   "alf",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name:        "Wants exact NUMBER match",
			filter:      Between("numbKey", "1", "2"),
			whenValue:   "1.5",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name:        "No NUMBER match",
			filter:      Between("numbKey", "1", "2"),
			whenValue:   "2.5",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name:        "No NUMBER match - not inclusive",
			filter:      Between("numbKey", "1", "2"),
			whenValue:   "2",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name:        "Wants BAD number on value",
			filter:      Between("numbKey", "1", "2"),
			whenValue:   "bananas",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name:        "Wants BAD number on expression",
			filter:      Between("numbKey", "1ogg0", "3"),
			whenValue:   "10",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name:        "Wants BAD number on expression2",
			filter:      Between("numbKey", "1", "ba"),
			whenValue:   "10",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name:        "Wants exact DATE match",
			filter:      Between("dateTimeKey", "2006-01-02T15:04:05-0700", "2006-03-02T15:04:05-0700"),
			whenValue:   "2006-02-02T15:04:05-0700",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name:        "No DATE match",
			filter:      Between("dateTimeKey", "2006-01-02T15:04:05-0700", "2006-03-02T15:04:05-0700"),
			whenValue:   "2020-01-02T15:04:05-0700",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name:        "Wants BAD DATE value",
			filter:      Between("dateTimeKey", "2006-01-02T15:04:05-0700", "2006-03-02T15:04:05-0700"),
			whenValue:   "asz",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name:        "Wants BAD DATE expression",
			filter:      Between("dateTimeKey", "xyz", "2006-03-02T15:04:05-0700"),
			whenValue:   "2020-01-02T15:04:05-0700",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name:        "Wants BAD DATE expression2",
			filter:      Between("dateTimeKey", "2006-03-02T15:04:05-0700", "abc"),
			whenValue:   "2020-01-02T15:04:05-0700",
			shouldMatch: false,
			wantError:   true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testFilterFunc(t, test)
		})
	}
}

func TestBetweenInclusive(t *testing.T) {
	tests := []testFilter{
		{
			name:        "Wants exact STRING match",
			filter:      BetweenInclusive("strName", "minion", "zorg"),
			whenValue:   "minion",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name:        "No STRING match",
			filter:      BetweenInclusive("strName", "minion", "zorg"),
			whenValue:   "alf",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name:        "Wants exact NUMBER match",
			filter:      BetweenInclusive("numbKey", "1", "2"),
			whenValue:   "1",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name:        "No NUMBER match",
			filter:      BetweenInclusive("numbKey", "1", "2"),
			whenValue:   "2.5",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name:        "Wants BAD number on value",
			filter:      BetweenInclusive("numbKey", "1", "2"),
			whenValue:   "bananas",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name:        "Wants BAD number on expression",
			filter:      BetweenInclusive("numbKey", "1ogg0", "3"),
			whenValue:   "10",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name:        "Wants BAD number on expression2",
			filter:      BetweenInclusive("numbKey", "1", "ba"),
			whenValue:   "10",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name:        "Wants exact DATE match",
			filter:      BetweenInclusive("dateTimeKey", "2006-01-02T15:04:05-0700", "2006-03-02T15:04:05-0700"),
			whenValue:   "2006-01-02T15:04:05-0700",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name:        "No DATE match",
			filter:      BetweenInclusive("dateTimeKey", "2006-01-02T15:04:05-0700", "2006-03-02T15:04:05-0700"),
			whenValue:   "2020-01-02T15:04:05-0700",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name:        "Wants BAD DATE value",
			filter:      BetweenInclusive("dateTimeKey", "2006-01-02T15:04:05-0700", "2006-03-02T15:04:05-0700"),
			whenValue:   "asz",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name:        "Wants BAD DATE expression",
			filter:      BetweenInclusive("dateTimeKey", "xyz", "2006-03-02T15:04:05-0700"),
			whenValue:   "2020-01-02T15:04:05-0700",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name:        "Wants BAD DATE expression2",
			filter:      BetweenInclusive("dateTimeKey", "2006-03-02T15:04:05-0700", "abc"),
			whenValue:   "2020-01-02T15:04:05-0700",
			shouldMatch: false,
			wantError:   true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testFilterFunc(t, test)
		})
	}
}

func TestLowerThan_Apply(t *testing.T) {
	tests := []testFilter{
		{
			name:        "Wants exact STRING match",
			filter:      LowerThan("strName", "z"),
			whenValue:   "a",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name:        "No STRING match",
			filter:      LowerThan("strName", "a"),
			whenValue:   "z",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name:        "Wants exact NUMBER match",
			filter:      LowerThan("numbKey", "0.02"),
			whenValue:   "0.01",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name:        "No NUMBER match",
			filter:      LowerThan("numbKey", "0.01"),
			whenValue:   "0.02",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name:        "Wants BAD number on value",
			filter:      LowerThan("numbKey", "0.0109"),
			whenValue:   "bananas",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name:        "Wants BAD number on expression",
			filter:      LowerThan("numbKey", "bananas"),
			whenValue:   "10",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name:        "Wants exact DATE match",
			filter:      LowerThan("dateTimeKey", "2006-01-02T15:04:05-0700"),
			whenValue:   "2006-01-02T14:04:05-0700",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name:        "No DATE match",
			filter:      LowerThan("dateTimeKey", "2006-01-02T14:04:05-0700"),
			whenValue:   "2006-01-02T15:04:05-0700",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name:        "Wants BAD DATE value",
			filter:      LowerThan("dateTimeKey", "2006-01-02T15:04:05-0700"),
			whenValue:   "bananas",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name:        "Wants BAD DATE expression",
			filter:      LowerThan("dateTimeKey", "bananas"),
			whenValue:   "2006-01-02T15:04:05-0700",
			shouldMatch: false,
			wantError:   true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testFilterFunc(t, test)
		})
	}
}

func TestLowerOrEqualThan_Apply(t *testing.T) {
	tests := []testFilter{
		{
			name:        "Wants exact STRING match",
			filter:      LowerOrEqualThan("strName", "z"),
			whenValue:   "a",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name:        "No STRING match",
			filter:      LowerOrEqualThan("strName", "a"),
			whenValue:   "z",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name:        "Wants exact NUMBER match",
			filter:      LowerOrEqualThan("numbKey", "0.02"),
			whenValue:   "0.01",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name:        "No NUMBER match",
			filter:      LowerOrEqualThan("numbKey", "0.01"),
			whenValue:   "0.02",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name:        "Wants BAD number on value",
			filter:      LowerOrEqualThan("numbKey", "0.0109"),
			whenValue:   "bananas",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name:        "Wants BAD number on expression",
			filter:      LowerOrEqualThan("numbKey", "bananas"),
			whenValue:   "10",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name:        "Wants exact DATE match",
			filter:      LowerOrEqualThan("dateTimeKey", "2006-01-02T15:04:05-0700"),
			whenValue:   "2006-01-02T14:04:05-0700",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name:        "No DATE match",
			filter:      LowerOrEqualThan("dateTimeKey", "2006-01-02T14:04:05-0700"),
			whenValue:   "2006-01-02T15:04:05-0700",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name:        "Wants BAD DATE value",
			filter:      LowerOrEqualThan("dateTimeKey", "2006-01-02T15:04:05-0700"),
			whenValue:   "bananas",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name:        "Wants BAD DATE expression",
			filter:      LowerOrEqualThan("dateTimeKey", "bananas"),
			whenValue:   "2006-01-02T15:04:05-0700",
			shouldMatch: false,
			wantError:   true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testFilterFunc(t, test)
		})
	}
}

func TestGreaterThan_Apply(t *testing.T) {
	tests := []testFilter{
		{
			name:        "Wants exact STRING match",
			filter:      GreaterThan("strName", "a"),
			whenValue:   "z",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name:        "No STRING match",
			filter:      GreaterThan("strName", "z"),
			whenValue:   "a",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name:        "Wants exact NUMBER match",
			filter:      GreaterThan("numbKey", "0.01"),
			whenValue:   "0.02",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name:        "No NUMBER match",
			filter:      GreaterThan("numbKey", "0.02"),
			whenValue:   "0.01",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name:        "Wants BAD number on value",
			filter:      GreaterThan("numbKey", "0.0109"),
			whenValue:   "bananas",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name:        "Wants BAD number on expression",
			filter:      GreaterThan("numbKey", "bananas"),
			whenValue:   "10",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name:        "Wants exact DATE match",
			filter:      GreaterThan("dateTimeKey", "2006-01-02T14:04:05-0700"),
			whenValue:   "2006-01-02T15:04:05-0700",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name:        "No DATE match",
			filter:      GreaterThan("dateTimeKey", "2006-01-02T15:04:05-0700"),
			whenValue:   "2006-01-02T14:04:05-0700",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name:        "Wants BAD DATE value",
			filter:      GreaterThan("dateTimeKey", "2006-01-02T15:04:05-0700"),
			whenValue:   "bananas",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name:        "Wants BAD DATE expression",
			filter:      GreaterThan("dateTimeKey", "bananas"),
			whenValue:   "2006-01-02T15:04:05-0700",
			shouldMatch: false,
			wantError:   true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testFilterFunc(t, test)
		})
	}
}

func TestGreaterOrEqualThan_Apply(t *testing.T) {
	tests := []testFilter{
		{
			name:        "Wants exact STRING match",
			filter:      GreaterOrEqualThan("strName", "a"),
			whenValue:   "z",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name:        "No STRING match",
			filter:      GreaterOrEqualThan("strName", "z"),
			whenValue:   "a",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name:        "Wants exact NUMBER match",
			filter:      GreaterOrEqualThan("numbKey", "0.01"),
			whenValue:   "0.02",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name:        "No NUMBER match",
			filter:      GreaterOrEqualThan("numbKey", "0.02"),
			whenValue:   "0.01",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name:        "Wants BAD number on value",
			filter:      GreaterOrEqualThan("numbKey", "0.0109"),
			whenValue:   "bananas",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name:        "Wants BAD number on expression",
			filter:      GreaterOrEqualThan("numbKey", "bananas"),
			whenValue:   "10",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name:        "Wants exact DATE match",
			filter:      GreaterOrEqualThan("dateTimeKey", "2006-01-02T14:04:05-0700"),
			whenValue:   "2006-01-02T15:04:05-0700",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name:        "No DATE match",
			filter:      GreaterOrEqualThan("dateTimeKey", "2006-01-02T15:04:05-0700"),
			whenValue:   "2006-01-02T14:04:05-0700",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name:        "Wants BAD DATE value",
			filter:      GreaterOrEqualThan("dateTimeKey", "2006-01-02T15:04:05-0700"),
			whenValue:   "bananas",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name:        "Wants BAD DATE expression",
			filter:      GreaterOrEqualThan("dateTimeKey", "bananas"),
			whenValue:   "2006-01-02T15:04:05-0700",
			shouldMatch: false,
			wantError:   true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testFilterFunc(t, test)
		})
	}
}

func testFilterFunc(t *testing.T, test testFilter) {
	got, err := test.filter.Apply(test.whenValue, keySet)
	if test.wantError {
		assert.NotNil(t, err)
		assert.Error(t, err)
	} else {
		assert.Equal(t, test.shouldMatch, got)
	}
}
