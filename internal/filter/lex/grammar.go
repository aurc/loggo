package lex

/*
predicate:		identifier operator value | identifier opFunction
operator: 		[=|!|<|<=|...]
opFunction:     contains(value) | containsIgnoreCase(value) | between(value | number, value | number) | match(value)
identifier: 	key name
value:   		string | number
group: 			(predicate [chainer group | predicate])
chainer: 		&& | ||
*/

type GrammarItem interface {
	IsTerminal() bool
	NextGrammarItem() []GrammarItem
	StartChar()
}

type GrPredicateWithOp struct {
	Identifier GrIdentifier
	Operator   GrOperator
	Value      GrValue
}

type GrPredicateWithFunc struct {
	Identifier   GrIdentifier
	OperatorFunc GrOpFunction
}

type GrOperator struct {
}

type GrOpFunction struct {
}

type GrIdentifier struct {
}

type GrValue struct {
}

type GrString struct {
}

type GrNumber struct {
}

type GrGroup struct {
}

type GrChainer struct {
}
