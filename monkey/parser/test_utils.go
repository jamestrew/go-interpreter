package parser

import (
	"fmt"
	"testing"

	"github.com/jamestrew/go-interpreter/monkey/ast"
	"github.com/jamestrew/go-interpreter/monkey/lexer"
)

func programSetup(t *testing.T, input string, stmtCnt int) (*ast.Program, *Parser) {
	lexer := lexer.New(input)
	parser := New(lexer)
	program := parser.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if stmtCnt != -1 && len(program.Statements) != stmtCnt {
		t.Errorf(
			"program.Statements does not contain %d statements. got=%d",
			stmtCnt,
			len(program.Statements),
		)
		t.Log(input)
	}

	return program, parser
}

func checkLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("s.TokenLiteral() not '%s'. got=%s", name, letStmt.Name)
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser, expectedErrCnt int) {
	errors := p.Errors()

	if len(errors) == expectedErrCnt {
		return
	}

	t.Errorf("parser.errors expected %d error(s). got=%d", expectedErrCnt, len(errors))
	for _, err := range errors {
		t.Errorf("parser error: %s", err)
	}
	t.FailNow()
}

func checkIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("il.Value not %d. got=%d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value, integ.TokenLiteral())
		return false
	}

	return true
}

func checkStringLiteral(t *testing.T, sl ast.Expression, value string) bool {
	literal, ok := sl.(*ast.StringLiteral)
	if !ok {
		t.Errorf("exp not *ast.StringLiteral. got=%T", sl)
		return false
	}

	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q. got=%q", "hello world", literal.Value)
		return false
	}
	return true
}

func checkIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp is not ast.Identifier. got=%T, value=%s", exp, value)
		return false
	}
	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}
	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value, ident.TokenLiteral())
		return false
	}
	return true
}

func checkBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s", value, bo.TokenLiteral())
		return false
	}

	return true
}

func checkLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return checkIntegerLiteral(t, exp, int64(v))
	case int64:
		return checkIntegerLiteral(t, exp, v)
	case string:
		return checkIdentifier(t, exp, v)
	case bool:
		return checkBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func checkPrefixExpression(
	t *testing.T,
	exp ast.Expression,
	value interface{},
	operator string,
) bool {
	opExp, ok := exp.(*ast.PrefixExpression)
	if !ok {
		t.Errorf("exp is not ast.OperatorExpression. got=%T(%s)", exp, exp)
		return false
	}
	if !checkLiteralExpression(t, opExp.Right, value) {
		return false
	}
	if opExp.Operator != operator {
		t.Errorf("opExp.Operator is not %s. got=%s", operator, opExp.Operator)
	}
	return true
}

func checkInfixExpression(
	t *testing.T,
	exp ast.Expression,
	left, right interface{},
	operator string,
) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.OperatorExpression. got=%T(%s)", exp, exp)
		return false
	}
	if !checkLiteralExpression(t, opExp.Left, left) {
		return false
	}
	if !checkLiteralExpression(t, opExp.Right, right) {
		return false
	}
	if opExp.Operator != operator {
		t.Errorf("opExp.Operator is not %s. got=%s", operator, opExp.Operator)
	}
	return true
}

func checkExpressionStatement(t *testing.T, program *ast.Program) *ast.ExpressionStatement {
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0],
		)
	}
	return stmt
}
