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
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/aurc/loggo/internal/config"
)

type LogicalOperator int

const (
	And LogicalOperator = iota
	Or
)

var (
	sqlLexer = lexer.MustSimple([]lexer.SimpleRule{
		{`Keyword`, `(?i)\b(MATCH|CONTAINSIC|CONTAINS|BETWEEN|AND|OR)\b`},
		{`Ident`, `[a-zA-Z_][a-zA-Z0-9_./]*`},
		{`Number`, `[-+]?\d*\.?\d+([eE][-+]?\d+)?`},
		{`String`, `'[^']*'|"[^"]*"`},
		{`Operators`, `<>|!=|<=|>=|==|[()=<>]`},
		{"whitespace", `\s+`},
	})

	cachedDef = make(map[string]Filter)

	parser = participle.MustBuild[Expression](
		participle.Lexer(sqlLexer),
		participle.Unquote("String"),
		participle.CaseInsensitive("Keyword"),
	)
)

func ParseFilterExpression(exp string) (*Expression, error) {
	return parser.ParseString("", exp)
}

func cachedOperation(op Operation, key string, v ...string) Filter {
	ck := fmt.Sprintf(`[%s:%s]:%+v`, op, key, v)
	if v, ok := cachedDef[ck]; ok {
		return v
	}
	var f Filter
	switch op {
	case OpNotEqual:
		f = NotEquals(key, v[0])
	case OpEquals:
		f = Equals(key, v[0])
	case OpEqualsIgnoreCase:
		f = EqualIgnoreCase(key, v[0])
	case OpLowerThan:
		f = LowerThan(key, v[0])
	case OpLowerOrEqualThan:
		f = LowerOrEqualThan(key, v[0])
	case OpGreaterThan:
		f = GreaterThan(key, v[0])
	case OpGreaterOrEqualThan:
		f = GreaterOrEqualThan(key, v[0])
	case OpContains:
		f = Contains(key, v[0])
	case OpContainsIgnoreCase:
		f = ContainsIgnoreCase(key, v[0])
	case OpMatchesRegex:
		f = MatchesRegex(key, v[0])
	case OpBetween:
		f = BetweenInclusive(key, v[0], v[1])
	}
	cachedDef[ck] = f
	return f
}

var operatorMap = map[string]LogicalOperator{"AND": And, "OR": Or}

func (o *LogicalOperator) Capture(s []string) error {
	*o = operatorMap[strings.ToUpper(s[0])]
	return nil
}

type Expression struct {
	Left  *Term     `@@`
	Right []*OpTerm `@@*`
}

type ConditionElement struct {
	Condition     *Condition   ` @@`
	GlobalToken   *GlobalToken `| @@ `
	Subexpression *Expression  `| "(" @@ ")"`
}

type GlobalToken struct {
	String *string `@String`
}

type Condition struct {
	Operand  string `@Ident`
	Operator string `@( "<>" | "<=" | ">=" | "=" | "==" | "<" | ">" | "!=" | "BETWEEN" | "CONTAINS" | "CONTAINSIC" | "MATCH" )`
	Value    *Value `@@`
	Value2   *Value `( "AND" @@ )*`
}

func (v *Value) ToString() string {
	if v.Number == nil {
		return *v.String
	} else {
		return fmt.Sprintf(`%f`, *v.Number)
	}
}

type Value struct {
	Number *float64 `( @Number`
	String *string  ` | @String )`
}

type OpValue struct {
	Operator         LogicalOperator   `@("AND")`
	ConditionElement *ConditionElement `@@`
}

type Term struct {
	Left  *ConditionElement `@@`
	Right []*OpValue        `@@*`
}

type OpTerm struct {
	Operator LogicalOperator `@("OR")`
	Term     *Term           `@@`
}

func (c LogicalOperator) Apply(l, r bool) bool {
	switch c {
	case And:
		return l && r
	case Or:
		return l || r
	}
	return false
}

func (c *ConditionElement) Apply(row map[string]interface{}, key map[string]*config.Key) (bool, error) {
	switch {
	case c.Condition != nil:
		return c.Condition.Apply(row, key)
	case c.GlobalToken != nil:
		return c.GlobalToken.Apply(row)
	default:
		return c.Subexpression.Apply(row, key)
	}
}

func (g *GlobalToken) Apply(row map[string]interface{}) (bool, error) {
	b, err := json.Marshal(row)
	if err != nil {
		return false, err
	}
	str := strings.ToLower(string(b))
	return strings.Contains(str, strings.ToLower(*g.String)), nil
}

func (c *Condition) Apply(row map[string]interface{}, key map[string]*config.Key) (bool, error) {
	var op Operation
	switch strings.ToUpper(c.Operator) {
	case "<>", "!=":
		op = OpNotEqual
	case "=":
		op = OpEqualsIgnoreCase
	case "==":
		op = OpEquals
	case "<":
		op = OpLowerThan
	case "<=":
		op = OpLowerOrEqualThan
	case ">":
		op = OpGreaterThan
	case ">=":
		op = OpGreaterOrEqualThan
	case "CONTAINS":
		op = OpContains
	case "CONTAINSIC":
		op = OpContainsIgnoreCase
	case "MATCH":
		op = OpMatchesRegex
	case "BETWEEN":
		op = OpBetween
	default:
		return false, fmt.Errorf("unrecognised operator %s", c.Operator)
	}
	v2 := ""
	if c.Value2 != nil {
		v2 = c.Value2.ToString()
	}
	fi := cachedOperation(op, c.Operand, c.Value.ToString(), v2)
	var k *config.Key
	if v, ok := key[fi.Name()]; ok {
		k = v
	} else {
		k = &config.Key{
			Name: fi.Name(),
			Type: config.TypeString,
		}
	}
	return fi.Apply(k.ExtractValue(row), key)
}

func (c *Term) Apply(row map[string]interface{}, key map[string]*config.Key) (bool, error) {
	lv, le := c.Left.Apply(row, key)
	if le != nil {
		return false, le
	}
	for _, r := range c.Right {
		rv, re := r.ConditionElement.Apply(row, key)
		if re != nil {
			return false, re
		}
		lv = r.Operator.Apply(lv, rv)
	}
	return lv, nil
}

func (c *Expression) Apply(row map[string]interface{}, key map[string]*config.Key) (bool, error) {
	lv, le := c.Left.Apply(row, key)
	if le != nil {
		return false, le
	}
	for _, r := range c.Right {
		rv, re := r.Term.Apply(row, key)
		if re != nil {
			return false, re
		}
		lv = r.Operator.Apply(lv, rv)
	}
	return lv, nil
}
