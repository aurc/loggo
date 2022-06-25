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

type operator string

const (
	and = operator("and")
	or  = operator("or")
)

type filterGroup struct {
	filters  []Filter
	groups   []FilterGroup
	operator operator
}

type FilterGroup interface {
	Resolve(row map[string]interface{}) (bool, error)
}

func (f *filterGroup) Resolve(row map[string]interface{}) (bool, error) {
	initVal := f.operator == and
	if len(f.groups) > 0 {
		for _, fg := range f.groups {
			val, err := fg.Resolve(row)
			if err != nil {
				return false, err
			}
			switch f.operator {
			case and:
				initVal = initVal && val
			case or:
				initVal = initVal || val
			}
		}
	} else if len(f.filters) > 0 {
		for _, fi := range f.filters {
			k := fi.Key()
			val, err := fi.Apply(k.ExtractValue(row))
			if err != nil {
				return false, err
			}
			switch f.operator {
			case and:
				initVal = initVal && val
			case or:
				initVal = initVal || val
			}
		}
	}
	return initVal, nil
}

func And(group ...FilterGroup) *filterGroup {
	return &filterGroup{
		groups:   group,
		operator: and,
	}
}

func Or(group ...FilterGroup) *filterGroup {
	return &filterGroup{
		groups:   group,
		operator: or,
	}
}

func AndFilters(filter ...Filter) *filterGroup {
	return &filterGroup{
		filters:  filter,
		operator: and,
	}
}

func OrFilters(filter ...Filter) *filterGroup {
	return &filterGroup{
		filters:  filter,
		operator: or,
	}
}
