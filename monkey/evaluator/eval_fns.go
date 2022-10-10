package evaluator

import (
	"github.com/jamestrew/go-interpreter/monkey/ast"
	"github.com/jamestrew/go-interpreter/monkey/object"
)

func evalProgram(statements []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range statements {
		result = Eval(stmt)
		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

func evalBangOperator(right object.Object) object.Object {
	if isObjTruthy(right) {
		return FALSE
	}
	return TRUE
}

func evalMinusPrefixOperator(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}
	intObj := right.(*object.Integer)
	value := intObj.Value
	intObj.Value = -value
	return intObj
}

func evalPrefixExpression(pe *ast.PrefixExpression) object.Object {
	right := Eval(pe.Right)

	if isError(right) {
		return right
	}

	switch pe.Operator {
	case "!":
		return evalBangOperator(right)
	case "-":
		return evalMinusPrefixOperator(right)
	default:
		return newError("unknown operator: %s%s", pe.Operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.Integer{Value: leftValue + rightValue}
	case "-":
		return &object.Integer{Value: leftValue - rightValue}
	case "*":
		return &object.Integer{Value: leftValue * rightValue}
	case "/":
		return &object.Integer{Value: leftValue / rightValue}
	case ">":
		return nativeBoolToBooleanObject(leftValue > rightValue)
	case "<":
		return nativeBoolToBooleanObject(leftValue < rightValue)
	case "==":
		return nativeBoolToBooleanObject(leftValue == rightValue)
	case "!=":
		return nativeBoolToBooleanObject(leftValue != rightValue)
	default:
		return newError("unknown operator: %s %s %s", left.Inspect(), operator, right.Type())
	}
}

func evalInfixExpression(ie *ast.InfixExpression) object.Object {
	left := Eval(ie.Left)
	if isError(left) {
		return left
	}

	right := Eval(ie.Right)
	if isError(right) {
		return right
	}

	leftType := left.Type()
	rightType := right.Type()
	switch {
	case leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(ie.Operator, left, right)
	case ie.Operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case ie.Operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case leftType != rightType:
		return newError("type mismatch: %s %s %s", leftType, ie.Operator, rightType)
	default:
		return newError("unknown operator: %s %s %s", leftType, ie.Operator, rightType)
	}
}

func evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)
	if isError(condition) {
		return condition
	}
	if isObjTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	} else {
		return NULL
	}
}

func evalBlockStatement(statements []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range statements {
		result = Eval(stmt)
		switch result := result.(type) {
		case *object.ReturnValue:
			return result
		case *object.Error:
			return result
		}
	}

	return result
}

func evalReturnStatement(rs *ast.ReturnStatement) object.Object {
	value := Eval(rs.Value)
	if isError(value) {
		return value
	}
	ret := &object.ReturnValue{Value: value}
	return ret
}
