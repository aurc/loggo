package lex

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexer_Lex(t *testing.T) {
	sel, err := parser.ParseString("", `a = 3`)
	assert.NoError(t, err)
	fmt.Println(sel)
}
