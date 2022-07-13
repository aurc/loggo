package lexer

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	sel, err := parser.ParseString("", `aab > 3.0 or x = 'xxx'`)
	assert.NoError(t, err)
	fmt.Println(sel)
}
