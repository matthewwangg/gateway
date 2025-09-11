package parser

import (
	"strconv"

	logger "github.com/matthewwangg/gateway/internal/logger"
	models "github.com/matthewwangg/gateway/internal/models"
)

type Parser struct {
	Lexer *Lexer
}

func NewParser(filepath string) *Parser {
	return &Parser{Lexer: NewLexer(filepath)}
}

func (p *Parser) Expect(tokenType TokenType) Token {
	token := p.Lexer.GetToken()
	if tokenType != token.Type {
		logger.Log.Error("parser failed to correctly parse service definition: expected " + string(tokenType) + " got " + string(token.Type))
		return Token{
			Type:   TOKEN_EOF,
			Lexeme: "EOF",
		}
	}
	return token
}

func (p *Parser) Parse() *models.ServiceDefinition {
	return p.ParseServiceDefinition()
}

func (p *Parser) ParseServiceDefinition() *models.ServiceDefinition {
	serviceDefinition := &models.ServiceDefinition{}
	p.ParseElements(serviceDefinition)
	p.Expect(TOKEN_EOF)
	return serviceDefinition
}

func (p *Parser) ParseElements(serviceDefinition *models.ServiceDefinition) {
	token := p.Lexer.Peek()
	if token.Type == TOKEN_KEY {
		p.ParseElement(serviceDefinition)
		p.ParseElements(serviceDefinition)
	}
}

func (p *Parser) ParseElement(serviceDefinition *models.ServiceDefinition) {
	keyToken := p.Expect(TOKEN_KEY)
	p.Expect(TOKEN_COLON)

	key := keyToken.Lexeme
	values := make([]string, 0)

	token := p.Lexer.Peek()
	if token.Type == TOKEN_NEWLINE {
		p.Expect(TOKEN_NEWLINE)
		values = p.ParseListValues()
	} else {
		valueToken := p.Expect(TOKEN_VALUE)
		p.Expect(TOKEN_NEWLINE)
		values = append(values, valueToken.Lexeme)
	}
	p.AddField(serviceDefinition, key, values)
}

func (p *Parser) ParseListValues() []string {
	p.Expect(TOKEN_DASH)
	valueToken := p.Expect(TOKEN_VALUE)
	p.Expect(TOKEN_NEWLINE)

	values := []string{valueToken.Lexeme}

	token := p.Lexer.Peek()
	if token.Type == TOKEN_DASH {
		values = append(values, p.ParseListValues()...)
	}
	return values
}

func (p *Parser) AddField(serviceDefinition *models.ServiceDefinition, key string, values []string) {
	switch key {
	case "name":
		serviceDefinition.Name = values[0]
		break
	case "replicas":
		replicas, err := strconv.Atoi(values[0])
		if err != nil {
			break
		}
		serviceDefinition.Replicas = replicas
		break
	case "addresses":
		serviceDefinition.Addresses = values
		break
	case "api_type":
		switch values[0] {
		case "GRPC":
			serviceDefinition.APIType = models.GRPC
			break
		case "REST":
			serviceDefinition.APIType = models.REST
			break
		default:
			serviceDefinition.APIType = models.GRPC
			break
		}
	case "endpoints":
		serviceDefinition.Endpoints = values
		break
	default:
		break
	}
}
