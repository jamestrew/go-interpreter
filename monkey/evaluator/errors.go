package evaluator

import (
	"fmt"

	"github.com/jamestrew/go-interpreter/monkey/object"
)

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func infixOperatorError(left, right object.Object, operator string) *object.Error {
	return newError("unknown infix operation: %s %s %s", left.Type(), operator, right.Type())
}

func wrongArgCountError(want, got int) *object.Error {
	return newError("wrong number of arguments. got=%d, want=%d", got, want)
}
