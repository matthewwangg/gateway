package parser

type TokenType string

const (
	TOKEN_KEY     TokenType = "key"
	TOKEN_VALUE   TokenType = "value"
	TOKEN_COLON   TokenType = "colon"
	TOKEN_DASH    TokenType = "dash"
	TOKEN_NEWLINE TokenType = "newline"
	TOKEN_EOF     TokenType = "eof"
)

type Token struct {
	Type   TokenType
	Lexeme string
}
