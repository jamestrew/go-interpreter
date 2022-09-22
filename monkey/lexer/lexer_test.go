package lexer

import (
	"github.com/jamestrew/go-interpreter/monkey/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `=+(){},;`

	test := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	lexer := New(input)
	for i, test_token := range test {
		token := lexer.NextToken()

		var errMsg string
		if token.Type != test_token.expectedType {
			errMsg = fmt.Sprintf("tests[%d] - tokentype wrong. expected=%q, got=%q\n", i, test_token.expectedType, token.Type)
		}
		if token.Literal != test_token.expectedLiteral {
			errMsg = fmt.Sprintf("%stests[%d] - tokenLiteral wrong. expected=%q, got=%q", errMsg, i, test_token.expectedLiteral, token.Literal)
		}

		if errMsg != "" {
			t.Fatalf(errMsg)
		}
	}
}
