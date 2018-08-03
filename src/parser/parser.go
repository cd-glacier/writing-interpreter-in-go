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
	return nil
}
