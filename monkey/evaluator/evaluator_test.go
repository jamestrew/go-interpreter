package evaluator

import (
	"testing"

	"github.com/jamestrew/go-interpreter/monkey/object"
	"github.com/jamestrew/go-interpreter/monkey/parser"
)

func TestEvalIntegerObject(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestBooleanObject(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	program := parser.ParseInput(input)
	return Eval(program)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	myInt, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf(
			"object is not an Integer with expected value %d. got=%T (%+v)",
			expected,
			obj,
			obj,
		)
		return false
	}

	if myInt.Value != expected {
		t.Errorf("object has wrong value. got=%d, expected=%d", myInt.Value, expected)
		return false
	}
	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	myBool, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf(
			"object is not an Boolean with expected value %t. got=%T (%+v)",
			expected,
			obj,
			obj,
		)
		return false
	}

	if myBool.Value != expected {
		t.Errorf("object has wrong value. got=%t, expected=%t", myBool.Value, expected)
		return false
	}
	return true
}
