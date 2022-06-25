/*
Copyright © 2022 Aurelio Calegari, et al.

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

func TestEqual_Apply(t *testing.T) {
	tests := []testFilter{
		{
			name: "Wants exact STRING match",
			filter: Equals(&config.Key{
				Name: "abc",
				Type: config.TypeString,
			}, "minion"),
			whenValue:   "minion",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name: "No STRING match",
			filter: Equals(&config.Key{
				Name: "abc",
				Type: config.TypeString,
			}, "min"),
			whenValue:   "minion",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name: "Wants exact BOOL match",
			filter: Equals(&config.Key{
				Name: "abc",
				Type: config.TypeString,
			}, "true"),
			whenValue:   "true",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name: "No BOOL match",
			filter: Equals(&config.Key{
				Name: "abc",
				Type: config.TypeBool,
			}, "true"),
			whenValue:   "bubbles",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name: "Wants BAD BOOL on value",
			filter: Equals(&config.Key{
				Name: "abc",
				Type: config.TypeBool,
			}, "false"),
			whenValue:   "bananas",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name: "Wants BAD BOOL on expression",
			filter: Equals(&config.Key{
				Name: "abc",
				Type: config.TypeBool,
			}, "bananas"),
			whenValue:   "false",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name: "Wants exact NUMBER match",
			filter: Equals(&config.Key{
				Name: "abc",
				Type: config.TypeNumber,
			}, "0.01"),
			whenValue:   "0.01",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name: "No NUMBER match",
			filter: Equals(&config.Key{
				Name: "abc",
				Type: config.TypeNumber,
			}, "0.0109"),
			whenValue:   "0.01",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name: "Wants BAD number on value",
			filter: Equals(&config.Key{
				Name: "abc",
				Type: config.TypeNumber,
			}, "0.0109"),
			whenValue:   "bananas",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name: "Wants BAD number on expression",
			filter: Equals(&config.Key{
				Name: "abc",
				Type: config.TypeNumber,
			}, "bananas"),
			whenValue:   "10",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name: "Wants exact DATE match",
			filter: Equals(&config.Key{
				Name:   "abc",
				Type:   config.TypeDateTime,
				Layout: "2006-01-02T15:04:05-0700",
			}, "2006-01-02T15:04:05-0700"),
			whenValue:   "2006-01-02T15:04:05-0700",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name: "No DATE match",
			filter: Equals(&config.Key{
				Name:   "abc",
				Type:   config.TypeDateTime,
				Layout: "2006-01-02T15:04:05-0700",
			}, "2006-01-02T15:04:05-0710"),
			whenValue:   "2006-01-02T15:04:05-0700",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name: "Wants BAD DATE value",
			filter: Equals(&config.Key{
				Name:   "abc",
				Type:   config.TypeDateTime,
				Layout: "2006-01-02T15:04:05-0700",
			}, "2006-01-02T15:04:05-0700"),
			whenValue:   "bananas",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name: "Wants BAD DATE expression",
			filter: Equals(&config.Key{
				Name:   "abc",
				Type:   config.TypeDateTime,
				Layout: "2006-01-02T15:04:05-0700",
			}, "bananas"),
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
			name: "Wants exact STRING match",
			filter: MatchesRegex(&config.Key{
				Name: "abc",
				Type: config.TypeString,
			}, `\d+[a-zA-Z]+`),
			whenValue:   "123LoGGo",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name: "No STRING match",
			filter: MatchesRegex(&config.Key{
				Name: "abc",
				Type: config.TypeString,
			}, "onion"),
			whenValue:   "minion",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name: "BAD Regex",
			filter: MatchesRegex(&config.Key{
				Name: "abc",
				Type: config.TypeString,
			}, `\`),
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
			name: "Wants exact STRING match",
			filter: Contains(&config.Key{
				Name: "abc",
				Type: config.TypeString,
			}, "io"),
			whenValue:   "minion",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name: "No STRING match",
			filter: Contains(&config.Key{
				Name: "abc",
				Type: config.TypeString,
			}, "onion"),
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
			name: "Wants exact STRING match",
			filter: EqualIgnoreCase(&config.Key{
				Name: "abc",
				Type: config.TypeString,
			}, "minion"),
			whenValue:   "mInioN",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name: "No STRING match",
			filter: EqualIgnoreCase(&config.Key{
				Name: "abc",
				Type: config.TypeString,
			}, "mInio"),
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
			name: "Wants exact STRING match",
			filter: ContainsIgnoreCase(&config.Key{
				Name: "abc",
				Type: config.TypeString,
			}, "minion"),
			whenValue:   "miNion",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name: "Wants contains STRING match",
			filter: ContainsIgnoreCase(&config.Key{
				Name: "abc",
				Type: config.TypeString,
			}, "mIN"),
			whenValue:   "minion",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name: "No STRING match",
			filter: Equals(&config.Key{
				Name: "abc",
				Type: config.TypeString,
			}, "m1N"),
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
			name: "Wants exact STRING match",
			filter: Between(&config.Key{
				Name: "abc",
				Type: config.TypeString,
			}, "minion", "zorg"),
			whenValue:   "onion",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name: "No STRING match",
			filter: Between(&config.Key{
				Name: "abc",
				Type: config.TypeString,
			}, "minion", "zorg"),
			whenValue:   "alf",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name: "Wants exact NUMBER match",
			filter: Between(&config.Key{
				Name: "abc",
				Type: config.TypeNumber,
			}, "1", "2"),
			whenValue:   "1.5",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name: "No NUMBER match",
			filter: Between(&config.Key{
				Name: "abc",
				Type: config.TypeNumber,
			}, "1", "2"),
			whenValue:   "2.5",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name: "No NUMBER match - not inclusive",
			filter: Between(&config.Key{
				Name: "abc",
				Type: config.TypeNumber,
			}, "1", "2"),
			whenValue:   "2",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name: "Wants BAD number on value",
			filter: Between(&config.Key{
				Name: "abc",
				Type: config.TypeNumber,
			}, "1", "2"),
			whenValue:   "bananas",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name: "Wants BAD number on expression",
			filter: Between(&config.Key{
				Name: "abc",
				Type: config.TypeNumber,
			}, "1ogg0", "3"),
			whenValue:   "10",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name: "Wants BAD number on expression2",
			filter: Between(&config.Key{
				Name: "abc",
				Type: config.TypeNumber,
			}, "1", "ba"),
			whenValue:   "10",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name: "Wants exact DATE match",
			filter: Between(&config.Key{
				Name:   "abc",
				Type:   config.TypeDateTime,
				Layout: "2006-01-02T15:04:05-0700",
			}, "2006-01-02T15:04:05-0700", "2006-03-02T15:04:05-0700"),
			whenValue:   "2006-02-02T15:04:05-0700",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name: "No DATE match",
			filter: Between(&config.Key{
				Name:   "abc",
				Type:   config.TypeDateTime,
				Layout: "2006-01-02T15:04:05-0700",
			}, "2006-01-02T15:04:05-0700", "2006-03-02T15:04:05-0700"),
			whenValue:   "2020-01-02T15:04:05-0700",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name: "Wants BAD DATE value",
			filter: Between(&config.Key{
				Name:   "abc",
				Type:   config.TypeDateTime,
				Layout: "2006-01-02T15:04:05-0700",
			}, "2006-01-02T15:04:05-0700", "2006-03-02T15:04:05-0700"),
			whenValue:   "asz",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name: "Wants BAD DATE expression",
			filter: Between(&config.Key{
				Name:   "abc",
				Type:   config.TypeDateTime,
				Layout: "2006-01-02T15:04:05-0700",
			}, "xyz", "2006-03-02T15:04:05-0700"),
			whenValue:   "2020-01-02T15:04:05-0700",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name: "Wants BAD DATE expression2",
			filter: Between(&config.Key{
				Name:   "abc",
				Type:   config.TypeDateTime,
				Layout: "2006-01-02T15:04:05-0700",
			}, "2006-03-02T15:04:05-0700", "abc"),
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
			name: "Wants exact STRING match",
			filter: BetweenInclusive(&config.Key{
				Name: "abc",
				Type: config.TypeString,
			}, "minion", "zorg"),
			whenValue:   "minion",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name: "No STRING match",
			filter: BetweenInclusive(&config.Key{
				Name: "abc",
				Type: config.TypeString,
			}, "minion", "zorg"),
			whenValue:   "alf",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name: "Wants exact NUMBER match",
			filter: BetweenInclusive(&config.Key{
				Name: "abc",
				Type: config.TypeNumber,
			}, "1", "2"),
			whenValue:   "1",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name: "No NUMBER match",
			filter: BetweenInclusive(&config.Key{
				Name: "abc",
				Type: config.TypeNumber,
			}, "1", "2"),
			whenValue:   "2.5",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name: "Wants BAD number on value",
			filter: BetweenInclusive(&config.Key{
				Name: "abc",
				Type: config.TypeNumber,
			}, "1", "2"),
			whenValue:   "bananas",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name: "Wants BAD number on expression",
			filter: BetweenInclusive(&config.Key{
				Name: "abc",
				Type: config.TypeNumber,
			}, "1ogg0", "3"),
			whenValue:   "10",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name: "Wants BAD number on expression2",
			filter: BetweenInclusive(&config.Key{
				Name: "abc",
				Type: config.TypeNumber,
			}, "1", "ba"),
			whenValue:   "10",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name: "Wants exact DATE match",
			filter: BetweenInclusive(&config.Key{
				Name:   "abc",
				Type:   config.TypeDateTime,
				Layout: "2006-01-02T15:04:05-0700",
			}, "2006-01-02T15:04:05-0700", "2006-03-02T15:04:05-0700"),
			whenValue:   "2006-01-02T15:04:05-0700",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name: "No DATE match",
			filter: BetweenInclusive(&config.Key{
				Name:   "abc",
				Type:   config.TypeDateTime,
				Layout: "2006-01-02T15:04:05-0700",
			}, "2006-01-02T15:04:05-0700", "2006-03-02T15:04:05-0700"),
			whenValue:   "2020-01-02T15:04:05-0700",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name: "Wants BAD DATE value",
			filter: BetweenInclusive(&config.Key{
				Name:   "abc",
				Type:   config.TypeDateTime,
				Layout: "2006-01-02T15:04:05-0700",
			}, "2006-01-02T15:04:05-0700", "2006-03-02T15:04:05-0700"),
			whenValue:   "asz",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name: "Wants BAD DATE expression",
			filter: BetweenInclusive(&config.Key{
				Name:   "abc",
				Type:   config.TypeDateTime,
				Layout: "2006-01-02T15:04:05-0700",
			}, "xyz", "2006-03-02T15:04:05-0700"),
			whenValue:   "2020-01-02T15:04:05-0700",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name: "Wants BAD DATE expression2",
			filter: BetweenInclusive(&config.Key{
				Name:   "abc",
				Type:   config.TypeDateTime,
				Layout: "2006-01-02T15:04:05-0700",
			}, "2006-03-02T15:04:05-0700", "abc"),
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
			name: "Wants exact STRING match",
			filter: LowerThan(&config.Key{
				Name: "abc",
				Type: config.TypeString,
			}, "z"),
			whenValue:   "a",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name: "No STRING match",
			filter: LowerThan(&config.Key{
				Name: "abc",
				Type: config.TypeString,
			}, "a"),
			whenValue:   "z",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name: "Wants exact NUMBER match",
			filter: LowerThan(&config.Key{
				Name: "abc",
				Type: config.TypeNumber,
			}, "0.02"),
			whenValue:   "0.01",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name: "No NUMBER match",
			filter: LowerThan(&config.Key{
				Name: "abc",
				Type: config.TypeNumber,
			}, "0.01"),
			whenValue:   "0.02",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name: "Wants BAD number on value",
			filter: LowerThan(&config.Key{
				Name: "abc",
				Type: config.TypeNumber,
			}, "0.0109"),
			whenValue:   "bananas",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name: "Wants BAD number on expression",
			filter: LowerThan(&config.Key{
				Name: "abc",
				Type: config.TypeNumber,
			}, "bananas"),
			whenValue:   "10",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name: "Wants exact DATE match",
			filter: LowerThan(&config.Key{
				Name:   "abc",
				Type:   config.TypeDateTime,
				Layout: "2006-01-02T15:04:05-0700",
			}, "2006-01-02T15:04:05-0700"),
			whenValue:   "2006-01-02T14:04:05-0700",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name: "No DATE match",
			filter: LowerThan(&config.Key{
				Name:   "abc",
				Type:   config.TypeDateTime,
				Layout: "2006-01-02T15:04:05-0700",
			}, "2006-01-02T14:04:05-0700"),
			whenValue:   "2006-01-02T15:04:05-0700",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name: "Wants BAD DATE value",
			filter: LowerThan(&config.Key{
				Name:   "abc",
				Type:   config.TypeDateTime,
				Layout: "2006-01-02T15:04:05-0700",
			}, "2006-01-02T15:04:05-0700"),
			whenValue:   "bananas",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name: "Wants BAD DATE expression",
			filter: LowerThan(&config.Key{
				Name:   "abc",
				Type:   config.TypeDateTime,
				Layout: "2006-01-02T15:04:05-0700",
			}, "bananas"),
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
			name: "Wants exact STRING match",
			filter: LowerOrEqualThan(&config.Key{
				Name: "abc",
				Type: config.TypeString,
			}, "z"),
			whenValue:   "a",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name: "No STRING match",
			filter: LowerOrEqualThan(&config.Key{
				Name: "abc",
				Type: config.TypeString,
			}, "a"),
			whenValue:   "z",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name: "Wants exact NUMBER match",
			filter: LowerOrEqualThan(&config.Key{
				Name: "abc",
				Type: config.TypeNumber,
			}, "0.02"),
			whenValue:   "0.01",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name: "No NUMBER match",
			filter: LowerOrEqualThan(&config.Key{
				Name: "abc",
				Type: config.TypeNumber,
			}, "0.01"),
			whenValue:   "0.02",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name: "Wants BAD number on value",
			filter: LowerOrEqualThan(&config.Key{
				Name: "abc",
				Type: config.TypeNumber,
			}, "0.0109"),
			whenValue:   "bananas",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name: "Wants BAD number on expression",
			filter: LowerOrEqualThan(&config.Key{
				Name: "abc",
				Type: config.TypeNumber,
			}, "bananas"),
			whenValue:   "10",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name: "Wants exact DATE match",
			filter: LowerOrEqualThan(&config.Key{
				Name:   "abc",
				Type:   config.TypeDateTime,
				Layout: "2006-01-02T15:04:05-0700",
			}, "2006-01-02T15:04:05-0700"),
			whenValue:   "2006-01-02T14:04:05-0700",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name: "No DATE match",
			filter: LowerOrEqualThan(&config.Key{
				Name:   "abc",
				Type:   config.TypeDateTime,
				Layout: "2006-01-02T15:04:05-0700",
			}, "2006-01-02T14:04:05-0700"),
			whenValue:   "2006-01-02T15:04:05-0700",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name: "Wants BAD DATE value",
			filter: LowerOrEqualThan(&config.Key{
				Name:   "abc",
				Type:   config.TypeDateTime,
				Layout: "2006-01-02T15:04:05-0700",
			}, "2006-01-02T15:04:05-0700"),
			whenValue:   "bananas",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name: "Wants BAD DATE expression",
			filter: LowerOrEqualThan(&config.Key{
				Name:   "abc",
				Type:   config.TypeDateTime,
				Layout: "2006-01-02T15:04:05-0700",
			}, "bananas"),
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
			name: "Wants exact STRING match",
			filter: GreaterThan(&config.Key{
				Name: "abc",
				Type: config.TypeString,
			}, "a"),
			whenValue:   "z",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name: "No STRING match",
			filter: GreaterThan(&config.Key{
				Name: "abc",
				Type: config.TypeString,
			}, "z"),
			whenValue:   "a",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name: "Wants exact NUMBER match",
			filter: GreaterThan(&config.Key{
				Name: "abc",
				Type: config.TypeNumber,
			}, "0.01"),
			whenValue:   "0.02",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name: "No NUMBER match",
			filter: GreaterThan(&config.Key{
				Name: "abc",
				Type: config.TypeNumber,
			}, "0.02"),
			whenValue:   "0.01",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name: "Wants BAD number on value",
			filter: GreaterThan(&config.Key{
				Name: "abc",
				Type: config.TypeNumber,
			}, "0.0109"),
			whenValue:   "bananas",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name: "Wants BAD number on expression",
			filter: GreaterThan(&config.Key{
				Name: "abc",
				Type: config.TypeNumber,
			}, "bananas"),
			whenValue:   "10",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name: "Wants exact DATE match",
			filter: GreaterThan(&config.Key{
				Name:   "abc",
				Type:   config.TypeDateTime,
				Layout: "2006-01-02T15:04:05-0700",
			}, "2006-01-02T14:04:05-0700"),
			whenValue:   "2006-01-02T15:04:05-0700",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name: "No DATE match",
			filter: GreaterThan(&config.Key{
				Name:   "abc",
				Type:   config.TypeDateTime,
				Layout: "2006-01-02T15:04:05-0700",
			}, "2006-01-02T15:04:05-0700"),
			whenValue:   "2006-01-02T14:04:05-0700",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name: "Wants BAD DATE value",
			filter: GreaterThan(&config.Key{
				Name:   "abc",
				Type:   config.TypeDateTime,
				Layout: "2006-01-02T15:04:05-0700",
			}, "2006-01-02T15:04:05-0700"),
			whenValue:   "bananas",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name: "Wants BAD DATE expression",
			filter: GreaterThan(&config.Key{
				Name:   "abc",
				Type:   config.TypeDateTime,
				Layout: "2006-01-02T15:04:05-0700",
			}, "bananas"),
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
			name: "Wants exact STRING match",
			filter: GreaterOrEqualThan(&config.Key{
				Name: "abc",
				Type: config.TypeString,
			}, "a"),
			whenValue:   "z",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name: "No STRING match",
			filter: GreaterOrEqualThan(&config.Key{
				Name: "abc",
				Type: config.TypeString,
			}, "z"),
			whenValue:   "a",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name: "Wants exact NUMBER match",
			filter: GreaterOrEqualThan(&config.Key{
				Name: "abc",
				Type: config.TypeNumber,
			}, "0.01"),
			whenValue:   "0.02",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name: "No NUMBER match",
			filter: GreaterOrEqualThan(&config.Key{
				Name: "abc",
				Type: config.TypeNumber,
			}, "0.02"),
			whenValue:   "0.01",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name: "Wants BAD number on value",
			filter: GreaterOrEqualThan(&config.Key{
				Name: "abc",
				Type: config.TypeNumber,
			}, "0.0109"),
			whenValue:   "bananas",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name: "Wants BAD number on expression",
			filter: GreaterOrEqualThan(&config.Key{
				Name: "abc",
				Type: config.TypeNumber,
			}, "bananas"),
			whenValue:   "10",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name: "Wants exact DATE match",
			filter: GreaterOrEqualThan(&config.Key{
				Name:   "abc",
				Type:   config.TypeDateTime,
				Layout: "2006-01-02T15:04:05-0700",
			}, "2006-01-02T14:04:05-0700"),
			whenValue:   "2006-01-02T15:04:05-0700",
			shouldMatch: true,
			wantError:   false,
		},
		{
			name: "No DATE match",
			filter: GreaterOrEqualThan(&config.Key{
				Name:   "abc",
				Type:   config.TypeDateTime,
				Layout: "2006-01-02T15:04:05-0700",
			}, "2006-01-02T15:04:05-0700"),
			whenValue:   "2006-01-02T14:04:05-0700",
			shouldMatch: false,
			wantError:   false,
		},
		{
			name: "Wants BAD DATE value",
			filter: GreaterOrEqualThan(&config.Key{
				Name:   "abc",
				Type:   config.TypeDateTime,
				Layout: "2006-01-02T15:04:05-0700",
			}, "2006-01-02T15:04:05-0700"),
			whenValue:   "bananas",
			shouldMatch: false,
			wantError:   true,
		},
		{
			name: "Wants BAD DATE expression",
			filter: GreaterOrEqualThan(&config.Key{
				Name:   "abc",
				Type:   config.TypeDateTime,
				Layout: "2006-01-02T15:04:05-0700",
			}, "bananas"),
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
	got, err := test.filter.Apply(test.whenValue)
	if test.wantError {
		assert.NotNil(t, err)
		assert.Error(t, err)
	} else {
		assert.Equal(t, test.shouldMatch, got)
	}
}
