package lex

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

type Token int

const (
	EOF = iota
	Illegal
	Identifier
	Predicate
	Operator
	Value2
	OpFunction
	Group
	Chainer
	INT
	SEMI            // ;
	RoundBraketL    // (
	RoundBraketR    // )
	Quote           // "
	OpEqual         // =
	OpNotEqual      // !=
	OpLowerThan     // <
	OpGreaterThan   // >
	OpLowerEqThan   // <=
	OpGreaterEqThan // >=
	OpBetween       // BETWEEN
	OpRegex         // MATCH
	OpAnd           // &&
	OpOr            // ||
)

var tokens = []string{
	EOF:             "EOF",
	Illegal:         "Illegal",
	Identifier:      "Identifier",
	INT:             "INT",
	SEMI:            ";",
	RoundBraketL:    `(`,
	RoundBraketR:    `)`,
	Quote:           `"`,
	Operator:        "Operator",
	OpEqual:         `=`,
	OpNotEqual:      `!=`,
	OpLowerThan:     `<`,
	OpGreaterThan:   `>`,
	OpLowerEqThan:   `<=`,
	OpGreaterEqThan: `>=`,
	OpBetween:       `BETWEEN`,
	OpRegex:         `MATCH`,
	OpAnd:           `&&`,
	OpOr:            `||`,
}

var operators = map[string]bool{
	tokens[OpEqual]:         true,
	tokens[OpNotEqual]:      true,
	tokens[OpLowerThan]:     true,
	tokens[OpGreaterThan]:   true,
	tokens[OpLowerEqThan]:   true,
	tokens[OpGreaterEqThan]: true,
	tokens[OpBetween]:       true,
	tokens[OpRegex]:         true,
	tokens[OpAnd]:           true,
	tokens[OpOr]:            true,
}

func (t Token) String() string {
	return tokens[t]
}

type Position struct {
	line   int
	column int
}

type Lexer struct {
	pos    Position
	reader *bufio.Reader
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		pos:    Position{line: 1, column: 0},
		reader: bufio.NewReader(reader),
	}
}

// Lex scans the input for the next token. It returns the position of the token,
// the token's type, and the literal value.
func (l *Lexer) Lex() (Position, Token, string) {
	// keep looping until we return a token
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return l.pos, EOF, ""
			}
			// at this point there isn't much we can do, and the compiler
			// should just return the raw error to the user
			panic(err)
		}
		// update the column to the position of the newly read in rune
		l.pos.column++

		switch r {
		case '\n':
			l.resetPosition()
		case ';':
			return l.pos, SEMI, ";"
		default:
			if unicode.IsSpace(r) {
				continue // nothing to do here, just move on
			} else if unicode.IsDigit(r) {
				// backup and let lexInt rescan the beginning of the int
				startPos := l.pos
				l.backup()
				lit := l.lexInt()
				return startPos, INT, lit
			} else if unicode.IsLetter(r) {
				// backup and let lexIdent rescan the beginning of the ident
				startPos := l.pos
				l.backup()
				lit := l.lexIdent()
				return startPos, Identifier, lit
			} else if isOperator(fmt.Sprintf(`%c`, r)) {
				// backup and let lexIdent rescan the beginning of the ident
				startPos := l.pos
				l.backup()
				lit := l.lexOperator()
				return startPos, Operator, lit
			} else {
				return l.pos, Illegal, string(r)
			}
		}
	}
}

func isOperator(s string) bool {
	for k := range operators {
		if strings.Index(k, s) == 0 {
			return true
		}
	}
	return false
}

func (l *Lexer) resetPosition() {
	l.pos.line++
	l.pos.column = 0
}

func (l *Lexer) backup() {
	if err := l.reader.UnreadRune(); err != nil {
		panic(err)
	}

	l.pos.column--
}

// lexInt scans the input until the end of an integer and then returns the
// literal.
func (l *Lexer) lexInt() string {
	var lit string
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// at the end of the int
				return lit
			}
		}

		l.pos.column++
		if unicode.IsDigit(r) {
			lit = lit + string(r)
		} else {
			// scanned something not in the integer
			l.backup()
			return lit
		}
	}
}

// lexOperator scans the input until the end of an operator and then returns the
// operator.
func (l *Lexer) lexOperator() string {
	var op string
	opStarted := false
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// at the end of the identifier
				return op
			}
		}

		l.pos.column++
		if !opStarted && isOperator(fmt.Sprintf(`%c`, r)) {
			opStarted = true
			op = op + string(r)
		} else if opStarted && isOperator(fmt.Sprintf(`%s%c`, op, r)) {
			op = op + string(r)
		} else {
			// scanned something not in the identifier
			l.backup()
			return op
		}
	}
}

// lexIdent scans the input until the end of an identifier and then returns the
// literal.
func (l *Lexer) lexIdent() string {
	var lit string
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// at the end of the identifier
				return lit
			}
		}

		l.pos.column++
		if unicode.IsLetter(r) {
			lit = lit + string(r)
		} else {
			// scanned something not in the identifier
			l.backup()
			return lit
		}
	}
}

func main() {
	file, err := os.Open("input.test")
	if err != nil {
		panic(err)
	}

	lexer := NewLexer(file)
	for {
		pos, tok, lit := lexer.Lex()
		if tok == EOF {
			break
		}

		fmt.Printf("%d:%d\t%s\t%s\n", pos.line, pos.column, tok, lit)
	}
}
