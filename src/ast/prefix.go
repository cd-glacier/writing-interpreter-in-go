package ast

import (
	"bytes"

	"github.com/g-hyoga/writing-interpreter-in-go/src/token"
)

// PrefixExpression implements ast.Expression.
type PrefixExpression struct {
	Token    token.Token // the prefixx token, e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}

func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}
