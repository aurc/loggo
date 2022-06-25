/*
Copyright © 2022 Aurelio Calegari

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
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aurc/loggo/internal/config"
)

type Filter interface {
	Apply(value string) (bool, error)
	Key() *config.Key
}

func Equals(key *config.Key, expression string) *equals {
	return &equals{
		predicate: predicate{
			key:        key,
			expression: expression,
		},
	}
}

func EqualIgnoreCase(key *config.Key, expression string) *equalsIgnoreCase {
	return &equalsIgnoreCase{
		predicate: predicate{
			key:        key,
			expression: expression,
		},
	}
}

func Contains(key *config.Key, expression string) *contains {
	return &contains{
		predicate: predicate{
			key:        key,
			expression: expression,
		},
	}
}

func ContainsIgnoreCase(key *config.Key, expression string) *containsIgnoreCase {
	return &containsIgnoreCase{
		predicate: predicate{
			key:        key,
			expression: expression,
		},
	}
}

func LowerThan(key *config.Key, expression string) *lowerThan {
	return &lowerThan{
		predicate: predicate{
			key:        key,
			expression: expression,
		},
	}
}

func GreaterThan(key *config.Key, expression string) *greaterThan {
	return &greaterThan{
		predicate: predicate{
			key:        key,
			expression: expression,
		},
	}
}

func LowerOrEqualThan(key *config.Key, expression string) *lowerOrEqualThan {
	return &lowerOrEqualThan{
		predicate: predicate{
			key:        key,
			expression: expression,
		},
	}
}

func GreaterOrEqualThan(key *config.Key, expression string) *greaterOrEqualThan {
	return &greaterOrEqualThan{
		predicate: predicate{
			key:        key,
			expression: expression,
		},
	}
}

func MatchesRegex(key *config.Key, expression string) *matchRegex {
	return &matchRegex{
		predicate: predicate{
			key:        key,
			expression: expression,
		},
	}
}

func Between(key *config.Key, expression, expression2 string) *between {
	return &between{
		predicate: predicate{
			key:        key,
			expression: expression,
		},
		expression2: expression2,
	}
}

func BetweenInclusive(key *config.Key, expression, expression2 string) *betweenInclusive {
	return &betweenInclusive{
		between: between{
			predicate: predicate{
				key:        key,
				expression: expression,
			},
			expression2: expression2,
		},
	}
}

type predicate struct {
	key        *config.Key
	expression string
}

func (p *predicate) Key() *config.Key {
	return p.key
}

type equals struct {
	predicate
}

func (f *equals) Apply(value string) (bool, error) {
	switch f.key.Type {
	case config.TypeString:
		return f.expression == value, nil
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
	predicate
}

func (f *contains) Apply(value string) (bool, error) {
	return strings.Index(value, f.expression) != -1, nil
}

type equalsIgnoreCase struct {
	predicate
}

func (f *equalsIgnoreCase) Apply(value string) (bool, error) {
	return strings.ToLower(f.expression) == strings.ToLower(value), nil
}

type containsIgnoreCase struct {
	predicate
}

func (f *containsIgnoreCase) Apply(value string) (bool, error) {
	return strings.Index(strings.ToLower(value), strings.ToLower(f.expression)) != -1, nil
}

type matchRegex struct {
	predicate
}

func (f *matchRegex) Apply(value string) (bool, error) {
	reg, err := regexp.Compile(f.expression)
	if err != nil {
		return false, err
	}
	return reg.Match([]byte(value)), nil
}

type lowerThan struct {
	predicate
}

func (f *lowerThan) Apply(value string) (bool, error) {
	switch f.key.Type {
	case config.TypeString:
		return strings.Compare(value, f.expression) < 0, nil
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
	predicate
}

func (f *greaterThan) Apply(value string) (bool, error) {
	switch f.key.Type {
	case config.TypeString:
		return strings.Compare(value, f.expression) > 0, nil
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
	predicate
}

func (f *lowerOrEqualThan) Apply(value string) (bool, error) {
	switch f.key.Type {
	case config.TypeString:
		return strings.Compare(value, f.expression) <= 0, nil
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
	predicate
}

func (f *greaterOrEqualThan) Apply(value string) (bool, error) {
	switch f.key.Type {
	case config.TypeString:
		return strings.Compare(value, f.expression) >= 0, nil
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
	predicate
	expression2 string
}

func (f *between) Apply(value string) (bool, error) {
	switch f.key.Type {
	case config.TypeString:
		return strings.Compare(value, f.expression) > 0 && strings.Compare(value, f.expression2) < 0, nil
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
		return strings.Compare(value, f.expression) >= 0 && strings.Compare(value, f.expression2) <= 0, nil
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

func (f *predicate) parseNumberAndCheck(value string, check func(number, expression float64) (bool, error)) (bool, error) {
	var n, e float64
	var err error
	tv := strings.TrimSpace(value)
	if len(tv) == 0 {
		value = "0"
	}
	n, err = strconv.ParseFloat(value, 64)
	if err == nil {
		e, err = strconv.ParseFloat(f.expression, 64)
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
		e, err = strconv.ParseFloat(f.expression, 64)
		if err == nil {
			e2, err = strconv.ParseFloat(f.expression2, 64)
			if err == nil {
				return check(v, e, e2)
			}
		}
	}
	return false, err
}

func (f *predicate) parseBoolAndCheck(value string, check func(value, expression bool) (bool, error)) (bool, error) {
	var v, e bool
	var err error
	v, err = strconv.ParseBool(value)
	if err == nil {
		e, err = strconv.ParseBool(f.expression)
		if err == nil {
			return check(v, e)
		}
	}
	return false, err
}

func (f *predicate) parseDateTimeAndCheck(value string, check func(value, expression time.Time) (bool, error)) (bool, error) {
	var v, e time.Time
	var err error
	v, err = time.Parse(f.key.Layout, value)
	if err == nil {
		e, err = time.Parse(f.key.Layout, f.expression)
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
		e, err = time.Parse(f.key.Layout, f.expression)
		if err == nil {
			e2, err = time.Parse(f.key.Layout, f.expression2)
			if err == nil {
				return check(v, e, e2)
			}
		}
	}
	return false, err
}
