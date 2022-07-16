package lexer

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	sel, err := parser.ParseString("", `aab > 3.0 or x = 'xxx' and (a.b >= 2 or a.c between 2 and 3)`)
	assert.NoError(t, err)
	fmt.Println(sel)
}

func Test_GrammarLexer(t *testing.T) {
	tests := []struct {
		name       string
		expression string
	}{
		{
			name:       "simple operator and operand",
			expression: `aab > 3.0`,
		},
		{
			name:       "simple operators and operands",
			expression: `aab > 3.0 or x = 'xxx' and a.b >= 2 or a.c != 'n'`,
		},
		{
			name:       "complex operators and operands",
			expression: `aab > 3.0 or x = 'xxx' and (a.b >= 2 or a.c != 'n')`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sel, err := parser.ParseString("", test.expression)
			assert.NoError(t, err)
			b, _ := json.MarshalIndent(sel, "", "  ")
			fmt.Println(string(b))
		})
	}
}
