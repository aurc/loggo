package lexer

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Expression struct {
	Or []*OrCondition `@@ ( "OR" @@ )*`
}

type OrCondition struct {
	And []*Condition `@@ ( "AND" @@ )*`
}

type Condition struct {
	Operand  string `@Ident`
	Operator string `@( "<>" | "<=" | ">=" | "=" | "<" | ">" | "!=" )`
	Value    *Value `@@`
}

type Value struct {
	//Wildcard bool     `(  @"*"`
	Number *float64 `( @Number`
	String *string  ` | @String )`
	//Boolean  *Boolean ` | @("TRUE" | "FALSE")`
	//Null     bool     ` | @"NULL"`
	//Array    *Array   ` | @@ )`
}

type Boolean bool

func (b *Boolean) Capture(values []string) error {
	*b = values[0] == "TRUE"
	return nil
}

type Array struct {
	Expressions []*Expression `"(" @@ ( "," @@ )* ")"`
}

var (
	cli struct {
		SQL string `arg:"" required:"" help:"SQL to parse."`
	}

	sqlLexer = lexer.MustSimple([]lexer.SimpleRule{
		{`Keyword`, `(?i)\b(TRUE|FALSE|NOT|BETWEEN|AND|OR|IN)\b`},
		{`Ident`, `[a-zA-Z_][a-zA-Z0-9_]*`},
		{`Number`, `[-+]?\d*\.?\d+([eE][-+]?\d+)?`},
		{`String`, `'[^']*'|"[^"]*"`},
		{`Operators`, `<>|!=|<=|>=|[-+*/%,.()=<>]`},
		{"whitespace", `\s+`},
	})
	parser = participle.MustBuild[Expression](
		participle.Lexer(sqlLexer),
		participle.Unquote("String"),
		participle.CaseInsensitive("Keyword"),
		// participle.Elide("Comment"),
		// Need to solve left recursion detection first, if possible.
		// participle.UseLookahead(),
	)
)
