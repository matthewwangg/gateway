package parser

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
	// Implement this
	return []Token{}
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
