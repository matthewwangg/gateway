package parser

import (
	"bufio"
	"os"
	"strings"

	logger "github.com/matthewwangg/gateway/internal/logger"
)

type Lexer struct {
	Tokens   []Token
	Position int
}

func NewLexer(filepath string) *Lexer {
	tokens := Tokenize(filepath)
	return &Lexer{
		Tokens:   tokens,
		Position: 0,
	}
}

func Tokenize(filepath string) []Token {
	file, err := os.Open(filepath)
	if err != nil {
		logger.Log.Error("error opening service definition file at " + filepath + "\n")
		return nil
	}
	defer func() {
		if err := file.Close(); err != nil {
			logger.Log.Error("error closing service definition file at " + filepath + "\n")
			return
		}
	}()

	tokens := make([]Token, 0)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "- ") {
			tokens = append(tokens, Token{
				Type:   TOKEN_DASH,
				Lexeme: "DASH",
			})
			tokens = append(tokens, Token{
				Type:   TOKEN_VALUE,
				Lexeme: strings.TrimPrefix(line, "- "),
			})
			tokens = append(tokens, Token{
				Type:   TOKEN_NEWLINE,
				Lexeme: "NEWLINE",
			})
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			tokens = append(tokens, Token{
				Type:   TOKEN_KEY,
				Lexeme: key,
			})
			tokens = append(tokens, Token{
				Type:   TOKEN_COLON,
				Lexeme: "COLON",
			})
			if len(parts[1]) > 0 {
				tokens = append(tokens, Token{
					Type:   TOKEN_VALUE,
					Lexeme: value,
				})
			}
			tokens = append(tokens, Token{
				Type:   TOKEN_NEWLINE,
				Lexeme: "NEWLINE",
			})
		}
	}

	tokens = append(tokens, Token{
		Type:   TOKEN_EOF,
		Lexeme: "EOF",
	})

	return tokens
}

func (l *Lexer) GetToken() Token {
	if l.Position >= len(l.Tokens) {
		return Token{
			Type:   TOKEN_EOF,
			Lexeme: "EOF",
		}
	}
	token := l.Tokens[l.Position]
	l.Position++
	return token
}

func (l *Lexer) Peek() Token {
	if l.Position >= len(l.Tokens) {
		return Token{
			Type:   TOKEN_EOF,
			Lexeme: "EOF",
		}
	}
	return l.Tokens[l.Position]
}
