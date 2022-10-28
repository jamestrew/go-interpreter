package evaluator

import (
	"fmt"

	"github.com/jamestrew/go-interpreter/monkey/ast"
	"github.com/jamestrew/go-interpreter/monkey/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

type Evaluator struct {
	env *object.Environment
}

func New(env *object.Environment) *Evaluator {
	return &Evaluator{env: env}
}

func (e *Evaluator) Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return e.evalProgram(node.Statements)
	case *ast.ExpressionStatement:
		return e.Eval(node.Expression)
	case *ast.PrefixExpression:
		return e.evalPrefixExpression(node)
	case *ast.InfixExpression:
		return e.evalInfixExpression(node)
	case *ast.IfExpression:
		return e.evalIfExpression(node)
	case *ast.BlockStatement:
		return e.evalBlockStatement(node.Statements)
	case *ast.ReturnStatement:
		return e.evalReturnStatement(node)
	case *ast.LetStatement:
		return e.evalLetStatement(node)
	case *ast.Identifier:
		return e.evalIdentifier(node)
	case *ast.FunctionLiteral:
		return e.evalFunctionLiteral(node)
	case *ast.CallExpression:
		return e.evalCallExpression(node)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.ArrayLiteral:
		return e.evalArrayLiteral(node)
	case *ast.IndexExpression:
		return e.evalIndexExpression(node)
	case *ast.HashLiteral:
		return e.evalHashLiteral(node)
	default:
		fmt.Printf("Eval: node type not handled: %T\n", node)
	}
	return nil
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func isObjTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Integer:
		return obj.Value != 0
	case *object.Boolean:
		return obj.Value
	case *object.Null:
		return false
	default:
		return true
	}
}
