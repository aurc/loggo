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

type Filter interface {
	Apply(value string) (bool, error)
	Key() *config.Key
	Expression() []string
	Name() string
}

type Predicate struct {
	key        *config.Key `json:"-" yaml:"-"`
	name       string      `json:"key" yaml:"key"`
	expression []string    `json:"expression" yaml:"expression"`
}

func (p *Predicate) Apply(value string) (bool, error) {
	return true, nil
}

func (p *Predicate) Key() *config.Key {
	return p.key
}

func (p *Predicate) Expression() []string {
	return p.expression
}

func (p *Predicate) Name() string {
	return p.name
}

func Equals(key *config.Key, expression string) *equals {
	return &equals{
		Predicate: Predicate{
			key:        key,
			expression: []string{expression},
		},
	}
}

func EqualIgnoreCase(key *config.Key, expression string) *equalsIgnoreCase {
	return &equalsIgnoreCase{
		Predicate: Predicate{
			key:        key,
			expression: []string{expression},
		},
	}
}

func Contains(key *config.Key, expression string) *contains {
	return &contains{
		Predicate: Predicate{
			key:        key,
			expression: []string{expression},
		},
	}
}

func ContainsIgnoreCase(key *config.Key, expression string) *containsIgnoreCase {
	return &containsIgnoreCase{
		Predicate: Predicate{
			key:        key,
			expression: []string{expression},
		},
	}
}

func LowerThan(key *config.Key, expression string) *lowerThan {
	return &lowerThan{
		Predicate: Predicate{
			key:        key,
			expression: []string{expression},
		},
	}
}

func GreaterThan(key *config.Key, expression string) *greaterThan {
	return &greaterThan{
		Predicate: Predicate{
			key:        key,
			expression: []string{expression},
		},
	}
}

func LowerOrEqualThan(key *config.Key, expression string) *lowerOrEqualThan {
	return &lowerOrEqualThan{
		Predicate: Predicate{
			key:        key,
			expression: []string{expression},
		},
	}
}

func GreaterOrEqualThan(key *config.Key, expression string) *greaterOrEqualThan {
	return &greaterOrEqualThan{
		Predicate: Predicate{
			key:        key,
			expression: []string{expression},
		},
	}
}

func MatchesRegex(key *config.Key, expression string) *matchRegex {
	return &matchRegex{
		Predicate: Predicate{
			key:        key,
			expression: []string{expression},
		},
	}
}

func Between(key *config.Key, expression, expression2 string) *between {
	return &between{
		Predicate: Predicate{
			key:        key,
			expression: []string{expression, expression2},
		},
	}
}

func BetweenInclusive(key *config.Key, expression, expression2 string) *betweenInclusive {
	return &betweenInclusive{
		between: between{
			Predicate: Predicate{
				key:        key,
				expression: []string{expression, expression2},
			},
		},
	}
}

type equals struct {
	Predicate
}

func (f *equals) Apply(value string) (bool, error) {
	switch f.key.Type {
	case config.TypeString:
		return f.expression[0] == value, nil
	case config.TypeNumber:
		return f.parseNumberAndCheck(value, func(number, expression float64) (bool, error) {
			return number == expression, nil
		})
	case config.TypeBool:
		return f.parseBoolAndCheck(value, func(value, expression bool) (bool, error) {
			return value == expression, nil
		})
	case config.TypeDateTime:
		return f.parseDateTimeAndCheck(value, func(value, expression time.Time) (bool, error) {
			return value.Equal(expression), nil
		})
	}
	return false, nil
}

type contains struct {
	Predicate
}

func (f *contains) Apply(value string) (bool, error) {
	return strings.Index(value, f.expression[0]) != -1, nil
}

type equalsIgnoreCase struct {
	Predicate
}

func (f *equalsIgnoreCase) Apply(value string) (bool, error) {
	return strings.ToLower(f.expression[0]) == strings.ToLower(value), nil
}

type containsIgnoreCase struct {
	Predicate
}

func (f *containsIgnoreCase) Apply(value string) (bool, error) {
	return strings.Index(strings.ToLower(value), strings.ToLower(f.expression[0])) != -1, nil
}

type matchRegex struct {
	Predicate
}

func (f *matchRegex) Apply(value string) (bool, error) {
	reg, err := regexp.Compile(f.expression[0])
	if err != nil {
		return false, err
	}
	return reg.Match([]byte(value)), nil
}

type lowerThan struct {
	Predicate
}

func (f *lowerThan) Apply(value string) (bool, error) {
	switch f.key.Type {
	case config.TypeString:
		return strings.Compare(value, f.expression[0]) < 0, nil
	case config.TypeNumber:
		return f.parseNumberAndCheck(value, func(number, expression float64) (bool, error) {
			return number < expression, nil
		})
	case config.TypeDateTime:
		return f.parseDateTimeAndCheck(value, func(value, expression time.Time) (bool, error) {
			return value.Before(expression), nil
		})
	}
	return false, nil
}

type greaterThan struct {
	Predicate
}

func (f *greaterThan) Apply(value string) (bool, error) {
	switch f.key.Type {
	case config.TypeString:
		return strings.Compare(value, f.expression[0]) > 0, nil
	case config.TypeNumber:
		return f.parseNumberAndCheck(value, func(number, expression float64) (bool, error) {
			return number > expression, nil
		})
	case config.TypeDateTime:
		return f.parseDateTimeAndCheck(value, func(value, expression time.Time) (bool, error) {
			return value.After(expression), nil
		})
	}
	return false, nil
}

type lowerOrEqualThan struct {
	Predicate
}

func (f *lowerOrEqualThan) Apply(value string) (bool, error) {
	switch f.key.Type {
	case config.TypeString:
		return strings.Compare(value, f.expression[0]) <= 0, nil
	case config.TypeNumber:
		return f.parseNumberAndCheck(value, func(number, expression float64) (bool, error) {
			return number <= expression, nil
		})
	case config.TypeDateTime:
		return f.parseDateTimeAndCheck(value, func(value, expression time.Time) (bool, error) {
			return value.Before(expression) || value.Equal(expression), nil
		})
	}
	return false, nil
}

type greaterOrEqualThan struct {
	Predicate
}

func (f *greaterOrEqualThan) Apply(value string) (bool, error) {
	switch f.key.Type {
	case config.TypeString:
		return strings.Compare(value, f.expression[0]) >= 0, nil
	case config.TypeNumber:
		return f.parseNumberAndCheck(value, func(number, expression float64) (bool, error) {
			return number >= expression, nil
		})
	case config.TypeDateTime:
		return f.parseDateTimeAndCheck(value, func(value, expression time.Time) (bool, error) {
			return value.After(expression) || value.Equal(expression), nil
		})
	}
	return false, nil
}

type between struct {
	Predicate
}

func (f *between) Apply(value string) (bool, error) {
	switch f.key.Type {
	case config.TypeString:
		return strings.Compare(value, f.expression[0]) > 0 && strings.Compare(value, f.expression[1]) < 0, nil
	case config.TypeNumber:
		return f.parseNumberAndCheck(value, func(number, expression, expression2 float64) (bool, error) {
			return number > expression && number < expression2, nil
		})
	case config.TypeDateTime:
		return f.parseDateTimeAndCheck(value, func(value, expression, expression2 time.Time) (bool, error) {
			return value.After(expression) && value.Before(expression2), nil
		})
	}
	return false, nil
}

type betweenInclusive struct {
	between
}

func (f *betweenInclusive) Apply(value string) (bool, error) {
	switch f.key.Type {
	case config.TypeString:
		return strings.Compare(value, f.expression[0]) >= 0 && strings.Compare(value, f.expression[1]) <= 0, nil
	case config.TypeNumber:
		return f.parseNumberAndCheck(value, func(number, expression, expression2 float64) (bool, error) {
			return number >= expression && number <= expression2, nil
		})
	case config.TypeDateTime:
		return f.parseDateTimeAndCheck(value, func(value, expression, expression2 time.Time) (bool, error) {
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
		e, err = strconv.ParseFloat(p.expression[0], 64)
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
		e, err = strconv.ParseFloat(f.expression[0], 64)
		if err == nil {
			e2, err = strconv.ParseFloat(f.expression[1], 64)
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
		e, err = strconv.ParseBool(p.expression[0])
		if err == nil {
			return check(v, e)
		}
	}
	return false, err
}

func (p *Predicate) parseDateTimeAndCheck(value string, check func(value, expression time.Time) (bool, error)) (bool, error) {
	var v, e time.Time
	var err error
	v, err = time.Parse(p.key.Layout, value)
	if err == nil {
		e, err = time.Parse(p.key.Layout, p.expression[0])
		if err == nil {
			return check(v, e)
		}
	}
	return false, err
}

func (f *between) parseDateTimeAndCheck(value string, check func(value, expression, expression2 time.Time) (bool, error)) (bool, error) {
	var v, e, e2 time.Time
	var err error
	v, err = time.Parse(f.key.Layout, value)
	if err == nil {
		e, err = time.Parse(f.key.Layout, f.expression[0])
		if err == nil {
			e2, err = time.Parse(f.key.Layout, f.expression[1])
			if err == nil {
				return check(v, e, e2)
			}
		}
	}
	return false, err
}
