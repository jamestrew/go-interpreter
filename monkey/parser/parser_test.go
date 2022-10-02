package parser

import (
	"testing"

	"github.com/jamestrew/go-interpreter/monkey/ast"
	"github.com/jamestrew/go-interpreter/monkey/token"
)

func TestLetStatements(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foobar = 838383;
	`
	program, parser := programSetup(t, input, 3)
	checkParserErrors(t, parser, 0)

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func TestReturnStatement(t *testing.T) {
	input := `
	return 5;
	return 10;
	return add(15);
	`
	program, parser := programSetup(t, input, 3)
	checkParserErrors(t, parser, 0)

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("ReturnStatement.TokenLiteral not 'return', got %q", returnStmt.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	program, parser := programSetup(t, input, 1)
	checkParserErrors(t, parser, 0)

	stmt := checkExpressionStatement(t, program)

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("ident is not ast.Identifier. got=%T", stmt.Expression)
	}
	if ident.Token.Type != token.IDENT {
		t.Fatalf("ident.Token.Type is not IDENT. got=%q", ident.Token.Type)
	}
	if ident.Value != "foobar" {
		t.Fatalf("ident.Value not %s. got=%s", input, ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Fatalf("ident.TokenLiteral() not %s. got=%s", input, ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	program, parser := programSetup(t, input, 1)
	checkParserErrors(t, parser, 0)

	stmt := checkExpressionStatement(t, program)

	intLiteral, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("intLiteral is not ast.IntegerLiteral. got=%T", stmt.Expression)
	}
	if intLiteral.Token.Type != token.INT {
		t.Fatalf("intLiteral.Token.Type is not INT. got=%q", intLiteral.Token.Type)
	}
	if intLiteral.Value != 5 {
		t.Fatalf("intLiteral.Value not %s. got=%d", input, intLiteral.Value)
	}
	if intLiteral.TokenLiteral() != "5" {
		t.Fatalf("intLiteral.TokenLiteral() not %s. got=%s", input, intLiteral.TokenLiteral())
	}
}

func TestPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	for _, tt := range prefixTests {
		program, parser := programSetup(t, tt.input, 1)
		checkParserErrors(t, parser, 0)

		stmt := checkExpressionStatement(t, program)
		checkPrefixExpression(t, stmt.Expression, tt.integerValue, tt.operator)
	}
}

func TestInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
	}

	for _, tt := range infixTests {
		program, parser := programSetup(t, tt.input, 1)
		checkParserErrors(t, parser, 0)

		stmt := checkExpressionStatement(t, program)
		checkInflixExpression(t, stmt.Expression, tt.leftValue, tt.rightValue, tt.operator)
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		stmtCnt  int
	}{
		{"-a * b", "((-a) * b)", 1},
		{"!-a", "(!(-a))", 1},
		{"a + b + c", "((a + b) + c)", 1},
		{"a + b - c", "((a + b) - c)", 1},
		{"a * b * c", "((a * b) * c)", 1},
		{"a * b / c", "((a * b) / c)", 1},
		{"a + b / c", "(a + (b / c))", 1},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)", 1},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)", 2},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))", 1},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))", 1},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))", 1},
	}

	for _, tt := range tests {
		program, parser := programSetup(t, tt.input, tt.stmtCnt)
		checkParserErrors(t, parser, 0)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input           string
		expectedBoolean bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		program, parser := programSetup(t, tt.input, 1)
		checkParserErrors(t, parser, 0)

		stmt := checkExpressionStatement(t, program)

		boolean, ok := stmt.Expression.(*ast.Boolean)
		if !ok {
			t.Fatalf("exp not *ast.Boolean. got=%T", stmt.Expression)
		}
		if boolean.Value != tt.expectedBoolean {
			t.Errorf("boolean.Value not %t. got=%t", tt.expectedBoolean, boolean.Value)
		}
	}
}
