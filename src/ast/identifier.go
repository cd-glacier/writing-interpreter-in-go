package ast

import "github.com/g-hyoga/writing-interpreter-in-go/src/token"

// Identifier implements ast.Expression interface.
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return i.Value
}
