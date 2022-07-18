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
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aurc/loggo/internal/config"
)

type Operation string

const (
	OpEquals             = Operation("OpEquals")
	OpEqualsIgnoreCase   = Operation("OpEqualsIgnoreCase")
	OpNotEquals          = Operation("OpNotEquals")
	OpEqualIgnoreCase    = Operation("OpEqualIgnoreCase")
	OpContains           = Operation("OpContains")
	OpContainsIgnoreCase = Operation("OpContainsIgnoreCase")
	OpLowerThan          = Operation("OpLowerThan")
	OpGreaterThan        = Operation("OpGreaterThan")
	OpLowerOrEqualThan   = Operation("OpLowerOrEqualThan")
	OpGreaterOrEqualThan = Operation("OpGreaterOrEqualThan")
	OpMatchesRegex       = Operation("OpMatchesRegex")
	OpBetween            = Operation("OpBetween")
	OpBetweenInclusive   = Operation("OpBetweenInclusive")
)

type Filter interface {
	Apply(value string, key map[string]*config.Key) (bool, error)
	Expression() []string
	Name() string
}

type Predicate struct {
	KeyName       string    `json:"key" yaml:"key"`
	KeyExpression []string  `json:"expression" yaml:"expression"`
	Operation     Operation `json:"operation" yaml:"operation"`
	Right         []Filter  `json:"right,omitempty" yaml:"right"`
}

func (p *Predicate) Apply(value string, key map[string]*config.Key) (bool, error) {
	return true, nil
}

func (p *Predicate) Expression() []string {
	return p.KeyExpression
}

func (p *Predicate) Name() string {
	return p.KeyName
}

func Equals(key string, expression string) *equals {
	return &equals{
		Predicate: Predicate{
			KeyName:       key,
			KeyExpression: []string{expression},
			Operation:     OpEquals,
		},
	}
}

func NotEquals(key string, expression string) *notEquals {
	return &notEquals{
		equals: equals{
			Predicate: Predicate{
				KeyName:       key,
				KeyExpression: []string{expression},
				Operation:     OpNotEquals,
			},
		},
	}
}

func EqualIgnoreCase(key string, expression string) *equalsIgnoreCase {
	return &equalsIgnoreCase{
		Predicate: Predicate{
			KeyName:       key,
			KeyExpression: []string{expression},
			Operation:     OpEqualIgnoreCase,
		},
	}
}

func Contains(key string, expression string) *contains {
	return &contains{
		Predicate: Predicate{
			KeyName:       key,
			KeyExpression: []string{expression},
			Operation:     OpContains,
		},
	}
}

func ContainsIgnoreCase(key string, expression string) *containsIgnoreCase {
	return &containsIgnoreCase{
		Predicate: Predicate{
			KeyName:       key,
			KeyExpression: []string{expression},
			Operation:     OpContainsIgnoreCase,
		},
	}
}

func LowerThan(key string, expression string) *lowerThan {
	return &lowerThan{
		Predicate: Predicate{
			KeyName:       key,
			KeyExpression: []string{expression},
			Operation:     OpLowerThan,
		},
	}
}

func GreaterThan(key string, expression string) *greaterThan {
	return &greaterThan{
		Predicate: Predicate{
			KeyName:       key,
			KeyExpression: []string{expression},
			Operation:     OpGreaterThan,
		},
	}
}

func LowerOrEqualThan(key string, expression string) *lowerOrEqualThan {
	return &lowerOrEqualThan{
		Predicate: Predicate{
			KeyName:       key,
			KeyExpression: []string{expression},
			Operation:     OpLowerOrEqualThan,
		},
	}
}

func GreaterOrEqualThan(key string, expression string) *greaterOrEqualThan {
	return &greaterOrEqualThan{
		Predicate: Predicate{
			KeyName:       key,
			KeyExpression: []string{expression},
			Operation:     OpGreaterOrEqualThan,
		},
	}
}

func MatchesRegex(key string, expression string) *matchRegex {
	return &matchRegex{
		Predicate: Predicate{
			KeyName:       key,
			KeyExpression: []string{expression},
			Operation:     OpMatchesRegex,
		},
	}
}

func Between(key string, expression, expression2 string) *between {
	return &between{
		Predicate: Predicate{
			KeyName:       key,
			KeyExpression: []string{expression, expression2},
			Operation:     OpBetween,
		},
	}
}

func BetweenInclusive(key string, expression, expression2 string) *betweenInclusive {
	return &betweenInclusive{
		between: between{
			Predicate: Predicate{
				KeyName:       key,
				KeyExpression: []string{expression, expression2},
				Operation:     OpBetweenInclusive,
			},
		},
	}
}

type equals struct {
	Predicate
}

func (f *equals) Apply(value string, key map[string]*config.Key) (bool, error) {
	var tp config.Type = config.TypeString
	var k *config.Key
	if v, ok := key[f.KeyName]; ok {
		tp = v.Type
		k = v
	}
	switch tp {
	case config.TypeString:
		return f.KeyExpression[0] == value, nil
	case config.TypeNumber:
		return f.parseNumberAndCheck(value, func(number, expression float64) (bool, error) {
			return number == expression, nil
		})
	case config.TypeBool:
		return f.parseBoolAndCheck(value, func(value, expression bool) (bool, error) {
			return value == expression, nil
		})
	case config.TypeDateTime:
		return f.parseDateTimeAndCheck(value, k, func(value, expression time.Time) (bool, error) {
			return value.Equal(expression), nil
		})
	}
	return false, nil
}

type notEquals struct {
	equals
}

func (f *notEquals) Apply(value string, key map[string]*config.Key) (bool, error) {
	v, err := f.equals.Apply(value, key)
	return !v, err
}

type contains struct {
	Predicate
}

func (f *contains) Apply(value string, key map[string]*config.Key) (bool, error) {
	return strings.Index(value, f.KeyExpression[0]) != -1, nil
}

type equalsIgnoreCase struct {
	Predicate
}

func (f *equalsIgnoreCase) Apply(value string, key map[string]*config.Key) (bool, error) {
	var tp config.Type = config.TypeString
	var k *config.Key
	if v, ok := key[f.KeyName]; ok {
		tp = v.Type
		k = v
	}
	switch tp {
	case config.TypeString:
		return strings.ToLower(f.KeyExpression[0]) == strings.ToLower(value), nil
	case config.TypeNumber:
		return f.parseNumberAndCheck(value, func(number, expression float64) (bool, error) {
			return number == expression, nil
		})
	case config.TypeBool:
		return f.parseBoolAndCheck(value, func(value, expression bool) (bool, error) {
			return value == expression, nil
		})
	case config.TypeDateTime:
		return f.parseDateTimeAndCheck(value, k, func(value, expression time.Time) (bool, error) {
			return value.Equal(expression), nil
		})
	}
	return false, nil
}

type containsIgnoreCase struct {
	Predicate
}

func (f *containsIgnoreCase) Apply(value string, key map[string]*config.Key) (bool, error) {
	return strings.Index(strings.ToLower(value), strings.ToLower(f.KeyExpression[0])) != -1, nil
}

type matchRegex struct {
	Predicate
}

func (f *matchRegex) Apply(value string, key map[string]*config.Key) (bool, error) {
	reg, err := regexp.Compile(f.KeyExpression[0])
	if err != nil {
		return false, err
	}
	return reg.Match([]byte(value)), nil
}

type lowerThan struct {
	Predicate
}

func (f *lowerThan) Apply(value string, key map[string]*config.Key) (bool, error) {
	var tp config.Type = config.TypeString
	var k *config.Key
	if v, ok := key[f.KeyName]; ok {
		tp = v.Type
		k = v
	}
	switch tp {
	case config.TypeString:
		return strings.Compare(value, f.KeyExpression[0]) < 0, nil
	case config.TypeNumber:
		return f.parseNumberAndCheck(value, func(number, expression float64) (bool, error) {
			return number < expression, nil
		})
	case config.TypeDateTime:
		return f.parseDateTimeAndCheck(value, k, func(value, expression time.Time) (bool, error) {
			return value.Before(expression), nil
		})
	}
	return false, nil
}

type greaterThan struct {
	Predicate
}

func (f *greaterThan) Apply(value string, key map[string]*config.Key) (bool, error) {
	var tp config.Type = config.TypeString
	var k *config.Key
	if v, ok := key[f.KeyName]; ok {
		tp = v.Type
		k = v
	}
	switch tp {
	case config.TypeString:
		return strings.Compare(value, f.KeyExpression[0]) > 0, nil
	case config.TypeNumber:
		return f.parseNumberAndCheck(value, func(number, expression float64) (bool, error) {
			return number > expression, nil
		})
	case config.TypeDateTime:
		return f.parseDateTimeAndCheck(value, k, func(value, expression time.Time) (bool, error) {
			return value.After(expression), nil
		})
	}
	return false, nil
}

type lowerOrEqualThan struct {
	Predicate
}

func (f *lowerOrEqualThan) Apply(value string, key map[string]*config.Key) (bool, error) {
	var tp config.Type = config.TypeString
	var k *config.Key
	if v, ok := key[f.KeyName]; ok {
		tp = v.Type
		k = v
	}
	switch tp {
	case config.TypeString:
		return strings.Compare(value, f.KeyExpression[0]) <= 0, nil
	case config.TypeNumber:
		return f.parseNumberAndCheck(value, func(number, expression float64) (bool, error) {
			return number <= expression, nil
		})
	case config.TypeDateTime:
		return f.parseDateTimeAndCheck(value, k, func(value, expression time.Time) (bool, error) {
			return value.Before(expression) || value.Equal(expression), nil
		})
	}
	return false, nil
}

type greaterOrEqualThan struct {
	Predicate
}

func (f *greaterOrEqualThan) Apply(value string, key map[string]*config.Key) (bool, error) {
	var tp config.Type = config.TypeString
	var k *config.Key
	if v, ok := key[f.KeyName]; ok {
		tp = v.Type
		k = v
	}
	switch tp {
	case config.TypeString:
		return strings.Compare(value, f.KeyExpression[0]) >= 0, nil
	case config.TypeNumber:
		return f.parseNumberAndCheck(value, func(number, expression float64) (bool, error) {
			return number >= expression, nil
		})
	case config.TypeDateTime:
		return f.parseDateTimeAndCheck(value, k, func(value, expression time.Time) (bool, error) {
			return value.After(expression) || value.Equal(expression), nil
		})
	}
	return false, nil
}

type between struct {
	Predicate
}

func (f *between) Apply(value string, key map[string]*config.Key) (bool, error) {
	var tp config.Type = config.TypeString
	var k *config.Key
	if v, ok := key[f.KeyName]; ok {
		tp = v.Type
		k = v
	}
	switch tp {
	case config.TypeString:
		return strings.Compare(value, f.KeyExpression[0]) > 0 && strings.Compare(value, f.KeyExpression[1]) < 0, nil
	case config.TypeNumber:
		return f.parseNumberAndCheck(value, func(number, expression, expression2 float64) (bool, error) {
			return number > expression && number < expression2, nil
		})
	case config.TypeDateTime:
		return f.parseDateTimeAndCheck(value, k, func(value, expression, expression2 time.Time) (bool, error) {
			return value.After(expression) && value.Before(expression2), nil
		})
	}
	return false, nil
}

type betweenInclusive struct {
	between
}

func (f *betweenInclusive) Apply(value string, key map[string]*config.Key) (bool, error) {
	var tp config.Type = config.TypeString
	var k *config.Key
	if v, ok := key[f.KeyName]; ok {
		tp = v.Type
		k = v
	}
	switch tp {
	case config.TypeString:
		return strings.Compare(value, f.KeyExpression[0]) >= 0 && strings.Compare(value, f.KeyExpression[1]) <= 0, nil
	case config.TypeNumber:
		return f.parseNumberAndCheck(value, func(number, expression, expression2 float64) (bool, error) {
			return number >= expression && number <= expression2, nil
		})
	case config.TypeDateTime:
		return f.parseDateTimeAndCheck(value, k, func(value, expression, expression2 time.Time) (bool, error) {
			return (value.After(expression) || value.Equal(expression)) &&
				(value.Before(expression2) || value.Equal(expression2)), nil
		})
	}
	return false, nil
}

func (p *Predicate) parseNumberAndCheck(value string, check func(number, expression float64) (bool, error)) (bool, error) {
	var n, e float64
	var err error
	tv := strings.TrimSpace(value)
	if len(tv) == 0 {
		value = "0"
	}
	n, err = strconv.ParseFloat(value, 64)
	if err == nil {
		e, err = strconv.ParseFloat(p.KeyExpression[0], 64)
		if err == nil {
			return check(n, e)
		}
	}
	return false, err
}

func (f *between) parseNumberAndCheck(value string, check func(number, expression, expression2 float64) (bool, error)) (bool, error) {
	var v, e, e2 float64
	var err error
	tv := strings.TrimSpace(value)
	if len(tv) == 0 {
		value = "0"
	}
	v, err = strconv.ParseFloat(value, 64)
	if err == nil {
		e, err = strconv.ParseFloat(f.KeyExpression[0], 64)
		if err == nil {
			e2, err = strconv.ParseFloat(f.KeyExpression[1], 64)
			if err == nil {
				return check(v, e, e2)
			}
		}
	}
	return false, err
}

func (p *Predicate) parseBoolAndCheck(value string, check func(value, expression bool) (bool, error)) (bool, error) {
	var v, e bool
	var err error
	v, err = strconv.ParseBool(value)
	if err == nil {
		e, err = strconv.ParseBool(p.KeyExpression[0])
		if err == nil {
			return check(v, e)
		}
	}
	return false, err
}

func (p *Predicate) parseDateTimeAndCheck(value string, key *config.Key, check func(value, expression time.Time) (bool, error)) (bool, error) {
	var v, e time.Time
	var err error
	v, err = time.Parse(key.Layout, value)
	if err == nil {
		e, err = time.Parse(key.Layout, p.KeyExpression[0])
		if err == nil {
			return check(v, e)
		}
	}
	return false, err
}

func (f *between) parseDateTimeAndCheck(value string, key *config.Key, check func(value, expression, expression2 time.Time) (bool, error)) (bool, error) {
	var v, e, e2 time.Time
	var err error
	v, err = time.Parse(key.Layout, value)
	if err == nil {
		e, err = time.Parse(key.Layout, f.KeyExpression[0])
		if err == nil {
			e2, err = time.Parse(key.Layout, f.KeyExpression[1])
			if err == nil {
				return check(v, e, e2)
			}
		}
	}
	return false, err
}
