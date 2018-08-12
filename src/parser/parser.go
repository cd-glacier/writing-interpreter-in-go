package parser

import (
	"fmt"
	"strconv"

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

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODCU      // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
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
	p.logger.Info("[parser] New")

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) ParseProgram() *ast.Program {
	p.logger.WithFields(logrus.Fields{
		"current_token": p.curToken.Literal,
	}).Info("[parser] ParseProgram")

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

func (p *Parser) parseExpression(precedence int) ast.Expression {
	p.logger.WithFields(logrus.Fields{
		"current_token": p.curToken.Literal,
	}).Info("[parser] parseExpression")

	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		return nil
	}
	leftExp := prefix()
	return leftExp
}

func (p *Parser) parseStatement() ast.Statement {
	p.logger.WithFields(logrus.Fields{
		"current_token": p.curToken.Literal,
	}).Info("[parser] parseStatement")

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
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	p.logger.WithFields(logrus.Fields{
		"current_token": p.curToken.Literal,
	}).Info("[parser] parseExpressionStatement")

	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	p.logger.WithFields(logrus.Fields{
		"current_token_literal": p.curToken.Literal,
		"peek_token_literal":    p.peekToken.Literal,
	}).Info("[parser] parseLetStatement")

	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.logger.WithFields(logrus.Fields{
		"current_token": p.curToken,
		"peek_token":    p.peekToken,
	}).Info("[parser] parseLetStatement")

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
	p.logger.WithFields(logrus.Fields{
		"current_token": p.curToken.Literal,
		"peek_token":    p.peekToken.Literal,
	}).Info("[parser] parseReturnStatement")

	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	// here is where to implement to parse value

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseIdentifier() ast.Expression {
	p.logger.WithFields(logrus.Fields{
		"current_token_literal": p.curToken.Literal,
		"current_token":         p.curToken,
	}).Info("[parser] parseIdentifier")
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	p.logger.WithFields(logrus.Fields{
		"current_token": p.curToken.Literal,
	}).Info("[parser] parseIntegerLiteral")

	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("cloud not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		p.logger.Errorf("[parser] %s", msg)
		return nil
	}
	lit.Value = value
	return lit
}
