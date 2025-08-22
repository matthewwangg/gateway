package parser

import (
	models "github.com/matthewwangg/gateway/internal/models"
	"log"
)

type Parser struct {
	Lexer *Lexer
}

func NewParser(filepath string) *Parser {
	return &Parser{Lexer: NewLexer(filepath)}
}

func (p *Parser) Expect(tokenType TokenType) {
	token := p.Lexer.GetToken()
	if tokenType != token.Type {
		log.Fatalf("parser failed to correctly parse service definition: expected %s got %s", tokenType, token.Type)
	}
}

func (p *Parser) Parse() models.ServiceDefinition {
	return p.ParseServiceDefinition()
}

func (p *Parser) ParseServiceDefinition() models.ServiceDefinition {
	p.ParseElements()
	p.Expect(TOKEN_EOF)
	return models.ServiceDefinition{}
}

func (p *Parser) ParseElements() {
	token := p.Lexer.Peek()
	if token.Type == TOKEN_KEY {
		p.ParseElement()
		p.ParseElements()
	}
}

func (p *Parser) ParseElement() {
	p.Expect(TOKEN_KEY)
	p.Expect(TOKEN_COLON)

	token := p.Lexer.Peek()
	if token.Type == TOKEN_NEWLINE {
		p.Expect(TOKEN_NEWLINE)
		p.ParseListValues()
	} else {
		p.Expect(TOKEN_VALUE)
		p.Expect(TOKEN_NEWLINE)
	}
}

func (p *Parser) ParseListValues() {
	p.Expect(TOKEN_DASH)
	p.Expect(TOKEN_VALUE)
	p.Expect(TOKEN_NEWLINE)

	token := p.Lexer.Peek()
	if token.Type == TOKEN_DASH {
		p.ParseListValues()
	}
}
