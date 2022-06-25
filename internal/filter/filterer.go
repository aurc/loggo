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
