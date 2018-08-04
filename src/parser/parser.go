package parser

import (
	"github.com/g-hyoga/writing-interpreter-in-go/src/ast"
	"github.com/g-hyoga/writing-interpreter-in-go/src/lexer"
	"github.com/g-hyoga/writing-interpreter-in-go/src/logger"
	"github.com/g-hyoga/writing-interpreter-in-go/src/token"
	"github.com/sirupsen/logrus"
)

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token
	logger    *logrus.Logger
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.logger = logger.New()

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()

	p.logger.WithFields(logrus.Fields{
		"currentToken": p.curToken.Literal,
		"peekToken":    p.peekToken.Literal,
	}).Debug("[parser] nextToken")
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

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		return false
	}
}
