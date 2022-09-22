package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

func New(tokenType TokenType, ch byte) Token {
	return Token{tokenType, string(ch)}
}

func MultiCharToken(literal string) TokenType {
	switch literal {
	case "fn":
		return FUNCTION
	case "let":
		return LET
	default:
		return IDENT
	}
}
