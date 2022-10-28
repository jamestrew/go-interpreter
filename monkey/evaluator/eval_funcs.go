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
		return infixOperatorError(left, right, operator)
	}
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	if operator != "+" {
		return infixOperatorError(left, right, operator)
	}
	return &object.String{Value: left.(*object.String).Value + right.(*object.String).Value}
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
	case leftType == object.STRING_OBJ && rightType == object.STRING_OBJ:
		return evalStringInfixExpression(ie.Operator, left, right)
	case ie.Operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case ie.Operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case leftType != rightType:
		return newError("type mismatch: %s %s %s", leftType, ie.Operator, rightType)
	default:
		return infixOperatorError(left, right, ie.Operator)
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
	if val, ok := e.env.Get(i.Value); ok {
		return val
	}

	if builtin, ok := builtins[i.Value]; ok {
		return builtin
	}
	return newError("identifier not found: %s", i.Value)
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

	args := e.evalExpressions(ce.Arguments)
	if len(args) == 1 && isError(args[0]) {
		return args[0]
	}

	return execFunction(function, args...)
}

func (e *Evaluator) evalArrayLiteral(al *ast.ArrayLiteral) object.Object {
	elements := e.evalExpressions(al.Elements)
	if len(elements) == 1 && isError(elements[0]) {
		return elements[0]
	}
	return &object.Array{Elements: elements}
}

func evalArrayIndex(array, index object.Object) object.Object {
	arr := array.(*object.Array)
	idx := index.(*object.Integer)

	maxIdx := len(arr.Elements)

	if idx.Value >= 0 && idx.Value < int64(maxIdx) {
		return arr.Elements[idx.Value]
	} else if idx.Value < 0 && -idx.Value <= int64(maxIdx) {
		return arr.Elements[maxIdx+int(idx.Value)]
	}
	return NULL
}

func (e *Evaluator) evalArrayIndexExpression(ie *ast.IndexExpression) object.Object {
	left := e.Eval(ie.Left)
	if isError(left) {
		return left
	}
	index := e.Eval(ie.Index)
	if isError(index) {
		return index
	}

	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndex(left, index)
	default:
		return newError("index operator not supported: %s", ie.String())
	}
}

func (e *Evaluator) evalHashLiteral(hl *ast.HashLiteral) object.Object {
	pairs := map[object.HashKey]object.HashPair{}

	for keyNode, valueNode := range hl.Pairs {
		key := e.Eval(keyNode)
		if isError(key) {
			return key
		}
		hashKey, ok := key.(object.Hashable)
		if !ok {
			return hashKeyError(key)
		}

		value := e.Eval(valueNode)
		if isError(value) {
			return value
		}
		pair := object.HashPair{Key: key, Value: value}
		pairs[hashKey.HashKey()] = pair
	}

	return &object.Hash{Pairs: pairs}
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

func extendFunctionEnv(fn *object.Function, args ...object.Object) *object.Environment {
	newEnv := object.NewEnclosedEnvironment(fn.Env)
	for paramIdx, param := range fn.Parameters {
		newEnv.Set(param.Value, args[paramIdx])
	}
	return newEnv
}

func execFunction(obj object.Object, args ...object.Object) object.Object {
	switch fn := obj.(type) {
	case *object.Function:
		newEnv := extendFunctionEnv(fn, args...)
		return New(newEnv).Eval(fn.Body)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", obj.Type())
	}
}
