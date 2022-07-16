package lexer

import (
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

var operatorMap = map[string]LogicalOperator{"AND": OpAND, "OR": OpOR}

func (o *LogicalOperator) Capture(s []string) error {
	*o = operatorMap[strings.ToUpper(s[0])]
	return nil
}

type ConditionElement struct {
	Condition     *Condition  ` @@`
	Subexpression *Expression `| "(" @@ ")"`
}

type Condition struct {
	Operand  string `@Ident`
	Operator string `@( "<>" | "<=" | ">=" | "=" | "<" | ">" | "!=" | "BETWEEN" | "CONTAINS" | "CONTAINS_I" | "MATCH" )`
	Value    *Value `@@`
	Value2   *Value `( "AND" @@ )*`
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

type Expression struct {
	Left  *Term     `@@`
	Right []*OpTerm `@@*`
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
