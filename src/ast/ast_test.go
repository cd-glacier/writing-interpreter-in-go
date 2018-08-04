package ast

import (
	"testing"

	"github.com/g-hyoga/writing-interpreter-in-go/src/lexer"
	"github.com/g-hyoga/writing-interpreter-in-go/src/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != "let myVar = anotherVar;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
return 5;	
return 10;	
return 993322;	
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrros(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain statements. got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
	}
}
