package evaluator

import (
	"fmt"
	"reflect"

	"github.com/jamestrew/go-interpreter/monkey/ast"
	"github.com/jamestrew/go-interpreter/monkey/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	default:
		fmt.Printf("Eval: node type not handled: %s\n", reflect.TypeOf(node))
	}
	return nil
}

func evalStatements(statements []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range statements {
		result = Eval(stmt)
	}
	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func objectTruthy(obj object.Object) *object.Boolean {
	switch obj := obj.(type) {
	case *object.Integer:
		if obj.Value == 0 {
			return FALSE
		} else {
			return TRUE
		}
	case *object.Boolean:
		return obj
	case *object.Null:
		return FALSE
	default:
		return TRUE
	}
}

func evalBangOperator(right object.Object) object.Object {
	if objectTruthy(right) == TRUE {
		return FALSE
	}
	return TRUE
}

func evalMinusPrefixOperator(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return NULL
	}
	intObj := right.(*object.Integer)
	value := intObj.Value
	intObj.Value = -value
	return intObj
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperator(right)
	case "-":
		return evalMinusPrefixOperator(right)
	default:
		return NULL
	}
}

func evalInfixExpression(
	operator string,
	left object.Object,
	right object.Object,
) object.Object {
	switch operator {
	case "+":
		value := left.(*object.Integer).Value + right.(*object.Integer).Value
		return &object.Integer{Value: value}
	case "-":
		value := left.(*object.Integer).Value - right.(*object.Integer).Value
		return &object.Integer{Value: value}
	case "*":
		value := left.(*object.Integer).Value * right.(*object.Integer).Value
		return &object.Integer{Value: value}
	case "/":
		value := left.(*object.Integer).Value / right.(*object.Integer).Value
		return &object.Integer{Value: value}
	default:
		return NULL // TODO
	}
}
