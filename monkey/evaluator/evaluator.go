package evaluator

import (
	"github.com/jamestrew/go-interpreter/monkey/ast"
	"github.com/jamestrew/go-interpreter/monkey/object"
)

func Eval(program *ast.Program) object.Object {
	stmt := program.Statements[0]
	exp := stmt.(*ast.ExpressionStatement)
	myInt := exp.Expression.(*ast.IntegerLiteral)
	obj := &object.Integer{Value: myInt.Value}
	return obj
}
