package parser

import (
	"testing"

	"github.com/g-hyoga/writing-interpreter-in-go/src/ast"
	"github.com/g-hyoga/writing-interpreter-in-go/src/lexer"
)

func TestLetStatement(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrros(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() return nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
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
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
			continue
		}
		expectedLiteral := "return"
		if returnStmt.TokenLiteral() != expectedLiteral {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q", returnStmt.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrros(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}

	expectedValue := "foobar"

	if ident.Value != expectedValue {
		t.Errorf("ident.Value not %s. got=%s", expectedValue, ident.Value)
	}
	if ident.TokenLiteral() != expectedValue {
		t.Errorf("ident.TokenLiteral not %s. got=%s", expectedValue, ident.TokenLiteral())
	}

}
