package evaluator

import (
	"github.com/jamestrew/go-interpreter/monkey/ast"
	"github.com/jamestrew/go-interpreter/monkey/object"
)

func (e *Evaluator) evalProgram(statements []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range statements {
		result = e.Eval(stmt)
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

func (e *Evaluator) evalPrefixExpression(pe *ast.PrefixExpression) object.Object {
	right := e.Eval(pe.Right)

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

func (e *Evaluator) evalInfixExpression(ie *ast.InfixExpression) object.Object {
	left := e.Eval(ie.Left)
	if isError(left) {
		return left
	}

	right := e.Eval(ie.Right)
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

func (e *Evaluator) evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := e.Eval(ie.Condition)
	if isError(condition) {
		return condition
	}
	if isObjTruthy(condition) {
		return e.Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return e.Eval(ie.Alternative)
	} else {
		return NULL
	}
}

func (e *Evaluator) evalBlockStatement(statements []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range statements {
		result = e.Eval(stmt)
		switch result := result.(type) {
		case *object.ReturnValue:
			return result
		case *object.Error:
			return result
		}
	}

	return result
}

func (e *Evaluator) evalReturnStatement(rs *ast.ReturnStatement) object.Object {
	value := e.Eval(rs.Value)
	if isError(value) {
		return value
	}
	ret := &object.ReturnValue{Value: value}
	return ret
}

func (e *Evaluator) evalLetStatement(ls *ast.LetStatement) object.Object {
	val := e.Eval(ls.Value)
	if isError(val) {
		return val
	}

	e.env.Set(ls.Name.Value, val)
	return val
}

func (e *Evaluator) evalIdentifier(i *ast.Identifier) object.Object {
	val, ok := e.env.Get(i.Value)
	if !ok {
		return newError("identifier not found: %s", i.Value)
	}
	return val
}

func (e *Evaluator) evalFunctionLiteral(fl *ast.FunctionLiteral) object.Object {
	params := fl.Parameters
	body := fl.Body
	return &object.Function{Parameters: params, Body: body, Env: e.env}
}

func (e *Evaluator) evalCallExpression(ce *ast.CallExpression) object.Object {
	function := e.Eval(ce.Function)
	if isError(function) {
		return function
	}

	fn, ok := function.(*object.Function)
	if !ok {
		return newError("not a function: %s", function.Type())
	}

	args := e.evalExpressions(ce.Arguments)
	if len(args) == 1 && isError(args[0]) {
		return args[0]
	}

	newEnv := object.NewEnclosedEnvironment(fn.Env)
	for paramIdx, param := range fn.Parameters {
		newEnv.Set(param.Value, args[paramIdx])
	}

	return New(newEnv).Eval(fn.Body)
}

func (e *Evaluator) evalExpressions(expressions []ast.Expression) []object.Object {
	var result []object.Object
	for _, expression := range expressions {
		expObj := e.Eval(expression)
		if isError(expObj) {
			return []object.Object{expObj}
		}
		result = append(result, expObj)
	}

	return result
}
