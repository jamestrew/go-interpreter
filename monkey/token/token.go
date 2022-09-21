package token


type TokenType string

type Token struct {
	Type TokenType
	Literal string
}

func New(tokenType TokenType, ch byte) Token {
	return Token{tokenType, string(ch)}
}
