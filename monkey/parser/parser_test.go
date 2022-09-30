package parser

import (
	"testing"

	"github.com/jamestrew/go-interpreter/monkey/ast"
	"github.com/jamestrew/go-interpreter/monkey/lexer"
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

func TestLetStatementsErrors(t *testing.T) {
	input := `
	let x 5;
	let = 10;
	let 838383;
	`
	program, parser := programSetup(t, input, 3)
	checkParserErrors(t, parser, 3)
	if len(program.Statements) != 3 {
		t.Log(program.Statements)
		t.Fatalf("program.Statements should contain 3 statements. got=%d", len(program.Statements))
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

func programSetup(t *testing.T, input string, stmtCnt int) (*ast.Program, *Parser) {
	lexer := lexer.New(input)
	parser := New(lexer)
	program := parser.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != stmtCnt {
		t.Fatalf(
			"program.Statements does not contain 3 statements. got=%d",
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
