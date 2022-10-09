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

func testIntegerObject(t *testing.T, obj object.Object, input string, expected int64) bool {
	myInt, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf(
			"object (%s) is not an Integer with expected value %d. got=%T (%+v)",
			input,
			expected,
			obj,
			obj,
		)
		return false
	}

	if myInt.Value != expected {
		t.Errorf("object (%s) has wrong value. got=%d, expected=%d", input, myInt.Value, expected)
		return false
	}
	return true
}

func testBooleanObject(t *testing.T, obj object.Object, input string, expected bool) bool {
	myBool, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf(
			"object (%s) is not an Boolean with expected value %t. got=%T (%+v)",
			input,
			expected,
			obj,
			obj,
		)
		return false
	}

	if myBool.Value != expected {
		t.Errorf("object (%s) has wrong value. got=%t, expected=%t", input, myBool.Value, expected)
		return false
	}
	return true
}

func testNullObject(t *testing.T, obj object.Object, input string) bool {
	_, ok := obj.(*object.Null)
	if !ok {
		t.Errorf("object (%s) is not null. got=%T (%+v)", input, obj, obj)
		return false
	}
	return true
}
