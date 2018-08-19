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
	p.logger.Debug("[parser] New")

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)

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
	}).Debug("[parser] ParseProgram")

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
		"current_token_literal": p.curToken.Literal,
		"current_token_type":    p.curToken.Type,
		"precedence":            precedence,
	}).Debug("[parser] enter parseExpression")

	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParserFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	p.logger.WithFields(logrus.Fields{
		"current_token": p.curToken,
	}).Debug("[parser] exit parseExpression")
	return leftExp
}

func (p *Parser) parseStatement() ast.Statement {
	p.logger.WithFields(logrus.Fields{
		"current_token": p.curToken.Literal,
	}).Debug("[parser] parseStatement")

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
	}).Debug("[parser] parseExpressionStatement")

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
	}).Debug("[parser] parseLetStatement")

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
	}).Debug("[parser] parseLetStatement")

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
	}).Debug("[parser] parseReturnStatement")

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
	}).Debug("[parser] parseIdentifier")
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	p.logger.WithFields(logrus.Fields{
		"current_token": p.curToken.Literal,
	}).Debug("[parser] parseIntegerLiteral")

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

func (p *Parser) parsePrefixExpression() ast.Expression {
	p.logger.WithFields(logrus.Fields{
		"current_token":         p.curToken,
		"current_token_literal": p.curToken.Literal,
	}).Debug("[parser] parsePrefixExpression")

	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}
