package parser

import (
	"github.com/g-hyoga/writing-interpreter-in-go/src/ast"
	"github.com/g-hyoga/writing-interpreter-in-go/src/lexer"
	"github.com/g-hyoga/writing-interpreter-in-go/src/logger"
	"github.com/g-hyoga/writing-interpreter-in-go/src/token"
	"github.com/sirupsen/logrus"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn

	errors []string

	logger *logrus.Logger
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	p.logger = logger.New()

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			p.logger.WithFields(logrus.Fields{
				"stmt": stmt,
			}).Debug("[parser] statement")
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		stmt := p.parseLetStatement()
		if stmt == nil {
			p.logger.Error("[parser] failed to parse let statement")
		}
		return stmt
	case token.RETURN:
		stmt := p.parseReturnStatement()
		if stmt == nil {
			p.logger.Error("[parser] failed to parse return statement")
		}
		return stmt
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// here is where to implement to parse value

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	p.logger.WithFields(logrus.Fields{
		"stmt.Token": stmt.Token,
		"stmt.Name":  stmt.Name,
		"stmt.Value": stmt.Value,
	}).Debug("[parser] parseLetStatement")

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}
