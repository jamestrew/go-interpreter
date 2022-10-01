package parser

import (
	"testing"

	"github.com/jamestrew/go-interpreter/monkey/ast"
	"github.com/jamestrew/go-interpreter/monkey/lexer"
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

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0],
		)
	}

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

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0],
		)
	}

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

func programSetup(t *testing.T, input string, stmtCnt int) (*ast.Program, *Parser) {
	lexer := lexer.New(input)
	parser := New(lexer)
	program := parser.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if stmtCnt != -1 && len(program.Statements) != stmtCnt {
		t.Fatalf(
			"program.Statements does not contain %d statements. got=%d",
			stmtCnt,
			len(program.Statements),
		)
	}

	return program, parser
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
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

