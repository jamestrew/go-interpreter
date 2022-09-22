package lexer

import (
	"fmt"
	"testing"

	"github.com/jamestrew/go-interpreter/monkey/token"
)

func TestNextToken(t *testing.T) {
	input := `let five = 5;
	let ten = 10;
	let add = fn(x, y) {
		x + y;
	}
	let result = add(five, ten);`

	test := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
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
