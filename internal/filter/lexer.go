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
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type LogicalOperator int

const (
	OpAND LogicalOperator = iota
	OpOR
)

var (
	sqlLexer = lexer.MustSimple([]lexer.SimpleRule{
		{`Keyword`, `(?i)\b(MATCH|CONTAINS_I|CONTAINS|BETWEEN|AND|OR|IN)\b`},
		{`Ident`, `[a-zA-Z_][a-zA-Z0-9_./]*`},
		{`Number`, `[-+]?\d*\.?\d+([eE][-+]?\d+)?`},
		{`String`, `'[^']*'|"[^"]*"`},
		{`Operators`, `<>|!=|<=|>=|[()=<>]`},
		{"whitespace", `\s+`},
	})

	parser = participle.MustBuild[Expression](
		participle.Lexer(sqlLexer),
		participle.Unquote("String"),
		participle.CaseInsensitive("Keyword"),
	)
)

func ParseFilterExpression(exp string) (*filterGroup, error) {
	sel, err := parser.ParseString("", exp)
	if err != nil {
		return nil, err
	}
	fg := sel.Eval()
	return fg, nil
}

var operatorMap = map[string]LogicalOperator{"AND": OpAND, "OR": OpOR}

func (o *LogicalOperator) Capture(s []string) error {
	*o = operatorMap[strings.ToUpper(s[0])]
	return nil
}

func (c *ConditionElement) Eval() *filterGroup {
	if c.Condition != nil {
		return c.Condition.Eval()
	} else {
		return c.Subexpression.Eval()
	}
}

type ConditionElement struct {
	Condition     *Condition  ` @@`
	Subexpression *Expression `| "(" @@ ")"`
}

func (c *Condition) Eval() *filterGroup {
	switch strings.ToUpper(c.Operator) {
	case "<>", "!=":
		return AndFilters(NotEquals(c.Operand, c.Value.ToString()))
	case "=":
		return AndFilters(Equals(c.Operand, c.Value.ToString()))
	case "<":
		return AndFilters(LowerThan(c.Operand, c.Value.ToString()))
	case "<=":
		return AndFilters(LowerOrEqualThan(c.Operand, c.Value.ToString()))
	case ">":
		return AndFilters(GreaterThan(c.Operand, c.Value.ToString()))
	case ">=":
		return AndFilters(GreaterOrEqualThan(c.Operand, c.Value.ToString()))
	case "CONTAINS":
		return AndFilters(Contains(c.Operand, c.Value.ToString()))
	case "CONTAINS_I":
		return AndFilters(ContainsIgnoreCase(c.Operand, c.Value.ToString()))
	case "MATCH":
		return AndFilters(MatchesRegex(c.Operand, c.Value.ToString()))
	case "BETWEEN":
		return AndFilters(Between(c.Operand, c.Value.ToString(), c.Value2.ToString()))
	}
	return AndFilters()
}

type Condition struct {
	Operand  string `@Ident`
	Operator string `@( "<>" | "<=" | ">=" | "=" | "<" | ">" | "!=" | "BETWEEN" | "CONTAINS" | "CONTAINS_I" | "MATCH" )`
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

func (t *Term) Eval() *filterGroup {
	fg := t.Left.Eval()
	if len(t.Right) > 0 {
		for _, v := range t.Right {
			switch v.Operator {
			case OpAND:
				fg.Groups = append(fg.Groups, And(v.ConditionElement.Eval()))
			case OpOR:
				fg.Groups = append(fg.Groups, Or(v.ConditionElement.Eval()))
			}
		}
	}
	return fg
}

type OpTerm struct {
	Operator LogicalOperator `@("OR")`
	Term     *Term           `@@`
}

func (t *OpTerm) Eval() *filterGroup {
	fg := &filterGroup{}
	switch t.Operator {
	case OpAND:
		fg.Groups = append(fg.Groups, And(t.Term.Eval()))
	case OpOR:
		fg.Groups = append(fg.Groups, Or(t.Term.Eval()))
	}
	return fg
}

type Expression struct {
	Left  *Term     `@@`
	Right []*OpTerm `@@*`
}

func (t *Expression) Eval() *filterGroup {
	fg := &filterGroup{}
	if t.Left != nil {
		fg.Groups = append(fg.Groups, t.Left.Eval())
	}
	if len(t.Right) > 0 {
		for _, v := range t.Right {
			f := v.Eval()
			fg.Groups = append(fg.Groups, f.Groups...)
		}
	}
	return fg
}

func (o LogicalOperator) String() string {
	switch o {
	case OpAND:
		return "AND"
	case OpOR:
		return "OR"
	}
	panic("unsupported operator")
}
