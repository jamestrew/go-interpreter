package parser

import (
	"testing"

	"github.com/jamestrew/go-interpreter/monkey/ast"
	"github.com/jamestrew/go-interpreter/monkey/token"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		program, parser := programSetup(t, tt.input, 1)
		checkParserErrors(t, parser, 0)
		stmt := program.Statements[0]
		if !checkLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.LetStatement).Value
		if !checkLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
	}

	for _, tt := range tests {
		program, parser := programSetup(t, tt.input, 1)
		checkParserErrors(t, parser, 0)
		stmt := program.Statements[0]
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("ReturnStatement.TokenLiteral not 'return', got %q", returnStmt.TokenLiteral())
		}

		val := stmt.(*ast.ReturnStatement).Value
		if !checkLiteralExpression(t, val, tt.expectedValue) {
			return
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
		integerValue interface{}
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
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
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		program, parser := programSetup(t, tt.input, 1)
		checkParserErrors(t, parser, 0)

		stmt := checkExpressionStatement(t, program)
		checkInfixExpression(t, stmt.Expression, tt.leftValue, tt.rightValue, tt.operator)
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
		{"true", "true", 1},
		{"false", "false", 1},
		{"3 > 5 == false", "((3 > 5) == false)", 1},
		{"3 < 5 == true", "((3 < 5) == true)", 1},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)", 1},
		{"(5 + 5) * 2", "((5 + 5) * 2)", 1},
		{"2 / (5 + 5)", "(2 / (5 + 5))", 1},
		{"-(5 + 5)", "(-(5 + 5))", 1},
		{"!(true == true)", "(!(true == true))", 1},
		{"a + add(b * c) + d", "((a + add((b * c))) + d)", 1},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
			1,
		},
		{"add(a + b + c * d / f + g)", "add((((a + b) + ((c * d) / f)) + g))", 1},
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

func TestIfExpression(t *testing.T) {
	input := "if (x < y) { x }"
	program, parser := programSetup(t, input, 1)
	checkParserErrors(t, parser, 0)

	stmt := checkExpressionStatement(t, program)

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T",
			stmt.Expression)
	}

	if !checkInfixExpression(t, exp.Condition, "x", "y", "<") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	if !checkIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil. got=%+v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	program, parser := programSetup(t, input, 1)
	checkParserErrors(t, parser, 0)

	stmt := checkExpressionStatement(t, program)

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if !checkInfixExpression(t, exp.Condition, "x", "y", "<") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	if !checkIdentifier(t, consequence.Expression, "x") {
		return
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("exp.Alternative.Statements does not contain 1 statements. got=%d\n",
			len(exp.Alternative.Statements))
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Alternative.Statements[0])
	}

	if !checkIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestFunctionLiteral(t *testing.T) {
	input := "fn(x, y) { x + y; }"

	program, parser := programSetup(t, input, 1)
	checkParserErrors(t, parser, 0)

	stmt := checkExpressionStatement(t, program)
	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not a ast.FunctionLiteral. got=%T", stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function expected 2 params. got=%d", len(function.Parameters))
	}

	checkLiteralExpression(t, function.Parameters[0], "x")
	checkLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf(
			"function.Body.Statements expects 1 statement. got=%d",
			len(function.Body.Statements),
		)
	}

	body, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"function.Body.Statements[0] is not ast.ExpressionStatement. got=%T",
			function.Body.Statements[0],
		)
	}
	checkInfixExpression(t, body.Expression, "x", "y", "+")
}

func TestFunctionParamter(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{"fn() {};", []string{}},
		{"fn(x) {};", []string{"x"}},
		{"fn(x, y, z) {};", []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		program, parser := programSetup(t, tt.input, 1)
		checkParserErrors(t, parser, 0)

		stmt := checkExpressionStatement(t, program)
		function := stmt.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(tt.expectedParams) {
			t.Fatalf(
				"%s expected %d params. got=%d",
				tt.input,
				len(tt.expectedParams),
				len(function.Parameters),
			)
		}

		for i, param := range tt.expectedParams {
			checkLiteralExpression(t, function.Parameters[i], param)
		}
	}
}

func TestCallExpression(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"
	program, parser := programSetup(t, input, 1)
	checkParserErrors(t, parser, 0)

	stmt := checkExpressionStatement(t, program)
	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt is not ast.CallExpression. got=%T", stmt.Expression)
	}

	if !checkIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("%s has 3 arguments. got=%d", input, len(exp.Arguments))
	}

	checkLiteralExpression(t, exp.Arguments[0], 1)
	checkInfixExpression(t, exp.Arguments[1], 2, 3, "*")
	checkInfixExpression(t, exp.Arguments[2], 4, 5, "+")
}

func TestStringLiteral(t *testing.T) {
	input := `"hello world";`
	program, parser := programSetup(t, input, 1)
	checkParserErrors(t, parser, 0)

	stmt := checkExpressionStatement(t, program)
	checkStringLiteral(t, stmt.Expression, "hello world")
}

	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q. got=%q", "hello world", literal.Value)
	}

}
