package evaluator

import (
	"testing"

	"github.com/jamestrew/go-interpreter/monkey/object"
	"github.com/jamestrew/go-interpreter/monkey/parser"
)

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
