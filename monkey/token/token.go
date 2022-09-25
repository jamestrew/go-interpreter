package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

func New(tokenType TokenType, ch byte) Token {
	return Token{tokenType, string(ch)}
}

func KeywordOrIdent(literal string) TokenType {
	switch literal {
	case "fn":
		return FUNCTION
	case "let":
		return LET
	case "true":
		return TRUE
	case "false":
		return FALSE
	case "return":
		return RETURN
	case "if":
		return IF
	case "else":
		return ELSE

	default:
		return IDENT
	}
}
