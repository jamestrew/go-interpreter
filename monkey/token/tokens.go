package token

const (
	ILLEGAL = "ILLEGAL"
	EOF = "EOF"

	IDENT = "IDENT"
	INT = "INT"

	ASSIGN = "="
	PLUS = "+"
	MINUS = "-"
	BANG = "!"
	ASTERISK = "*"
	SLASH = "/"

	EQ = "=="
	NOT_EQ = "!="
	LT = "<"
	GT = ">"
	LT_EQ = "<="
	GT_EQ = ">="


	COMMA = ","
	SEMICOLON = ";"


	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// keywords
	FUNCTION = "FUNCTION"
	LET = "LET"
	TRUE = "TRUE"
	FALSE = "FALSE"
	RETURN = "RETURN"
	IF = "IF"
	ELSE = "ELSE"
)
