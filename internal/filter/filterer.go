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

import "github.com/aurc/loggo/internal/config"

type Operator string

const (
	AND = Operator("AND")
	OR  = Operator("OR")
)

type filterGroup struct {
	Filters  []Filter `json:"filters,omitempty" yaml:"filters"`
	Groups   []Group  `json:"groups,omitempty" yaml:"groups"`
	Operator Operator `json:"Operator,omitempty" yaml:"Operator"`
}

type Group interface {
	//Filters() []Filter
	//Groups() []Group
	//Operator() Operator
	Resolve(row map[string]interface{}, keys map[string]*config.Key) (bool, error)
}

//func (f *filterGroup) Filters() []Filter {
//	return f.filters
//}
//
//func (f *filterGroup) Groups() []Group {
//	return f.groups
//}
//
//func (f *filterGroup) Operator() Operator {
//	return f.operator
//}

func (f *filterGroup) Resolve(row map[string]interface{}, keys map[string]*config.Key) (bool, error) {
	initVal := f.Operator == AND
	if len(f.Groups) > 0 {
		for _, fg := range f.Groups {
			val, err := fg.Resolve(row, keys)
			if err != nil {
				return false, err
			}
			switch f.Operator {
			case AND:
				initVal = initVal && val
			case OR:
				initVal = initVal || val
			}
		}
	} else if len(f.Filters) > 0 {
		for _, fi := range f.Filters {
			var k *config.Key
			if v, ok := keys[fi.Name()]; ok {
				k = v
			} else {
				k = &config.Key{
					Name: fi.Name(),
					Type: config.TypeString,
				}
			}
			val, err := fi.Apply(k.ExtractValue(row), keys)
			if err != nil {
				return false, err
			}
			switch f.Operator {
			case AND:
				initVal = initVal && val
			case OR:
				initVal = initVal || val
			}
		}
	}
	return initVal, nil
}

func And(group ...Group) *filterGroup {
	return &filterGroup{
		Groups:   group,
		Operator: AND,
	}
}

func Or(group ...Group) *filterGroup {
	return &filterGroup{
		Groups:   group,
		Operator: OR,
	}
}

func AndFilters(filter ...Filter) *filterGroup {
	return &filterGroup{
		Filters:  filter,
		Operator: AND,
	}
}

func OrFilters(filter ...Filter) *filterGroup {
	return &filterGroup{
		Filters:  filter,
		Operator: OR,
	}
}
