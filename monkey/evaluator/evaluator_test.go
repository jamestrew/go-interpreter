package evaluator

import (
	"testing"

	"github.com/jamestrew/go-interpreter/monkey/object"
)

func TestEvalIntegerObject(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.input, tt.expected)
	}
}

func TestBooleanObject(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.input, tt.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!!true", true},
		{"!false", true},
		{"!!false", false},
		{"!5", false},
		{"!!5", true},
		{"!0", true},
		{"!!0", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.input, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, tt.input, int64(integer))
		} else {
			testNullObject(t, evaluated, tt.input)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{`
		if (10 > 1) {
				if (10 > 1) {
					return 10;
				}
			return 1;
		}`,
			10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.input, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input       string
		expectedMsg string
	}{
		{"5 + true;", "type mismatch: INTEGER + BOOLEAN"},
		{"5 + true; 5;", "type mismatch: INTEGER + BOOLEAN"},
		{"-true", "unknown operator: -BOOLEAN"},
		{"true + false;", "unknown infix operation: BOOLEAN + BOOLEAN"},
		{"5; true + false; 5", "unknown infix operation: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) { true + false; }", "unknown infix operation: BOOLEAN + BOOLEAN"},
		{
			`
if (10 > 1) {
  if (10 > 1) {
    return true + false;
  }

  return 1;
}
`,
			"unknown infix operation: BOOLEAN + BOOLEAN",
		},
		{"foobar", "identifier not found: foobar"},
		{`"hello" - "world"`, "unknown infix operation: STRING - STRING"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned for `%s`. got=%T", tt.input, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMsg {
			t.Errorf(
				"wrong error message for `%s`. expected=%q, got=%q",
				tt.input,
				tt.expectedMsg,
				errObj.Message,
			)
		}

	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; a;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.input, tt.expected)
	}

}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; }"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	if fn.Body.String() != "(x + 2)" {
		t.Fatalf("body is not \"(x + 2)\". got=%q", fn.Body.String())
	}

}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.input, tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
let newAdder = fn(x) {
	fn(y) { x + y };
};

let addTwo = newAdder(2);
addTwo(2);
	`
	testIntegerObject(t, testEval(input), input, 4)
}

func TestEvalStringObject(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello world"`, "hello world"},
		{`"this is" + " cool"`, "this is cool"},
	}

	for _, tt := range tests {
		testStringObject(t, testEval(tt.input), tt.input, tt.expected)
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, tt.input, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T", evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", errObj.Message, expected)
				continue
			}
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{`[1, 2]`, []interface{}{1, 2}},
		{`["foo", "bar"]`, []interface{}{"foo", "bar"}},
		{`[]`, []interface{}{}},
		{`[1, "two", false]`, []interface{}{1, "two", false}},
		{`[1, 1 * 2, 1 + 2]`, []interface{}{1, 2, 3}},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		arr, ok := evaluated.(*object.Array)
		if !ok {
			t.Errorf("object is not Array. got=%T", evaluated)
			continue
		}

		for idx, expectedElem := range tt.expected {
			testElem := arr.Elements[idx]
			switch expectedElem := expectedElem.(type) {
			case int:
				testIntegerObject(t, testElem, tt.input, int64(expectedElem))
			case string:
				testStringObject(t, testElem, tt.input, expectedElem)
			case bool:
				testBooleanObject(t, testElem, tt.input, expectedElem)
			default:
				t.Errorf("test doesn't support element of type %T", expectedElem)
			}
		}
	}
}
