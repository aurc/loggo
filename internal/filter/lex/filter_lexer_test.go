package lex

import (
	"fmt"
	"strings"
	"testing"
)

func TestNewLexer(t *testing.T) {
	tests := []struct {
		name       string
		expression string
	}{
		{
			name:       "Name",
			expression: `keyA = "abc" OR keyB = "x"`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reader := strings.NewReader(test.expression)
			l := NewLexer(reader)
			tokens := getTokens(l)
			fmt.Printf(`%v`, tokens)
		})
	}
}

func getTokens(l *Lexer) []Token {
	tokens := make([]Token, 0)
	for {
		_, tok, _ := l.Lex()
		if tok == EOF {
			break
		}

		tokens = append(tokens, tok)
	}

	return tokens
}
