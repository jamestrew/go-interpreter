package evaluator

import (
	"fmt"

	"github.com/jamestrew/go-interpreter/monkey/object"
)

var builtins = map[string]*object.Builtin{
	"len":   {Fn: __len},
	"print": {Fn: __print},
	"first": {Fn: __first},
}

func __len(args ...object.Object) object.Object {
	if len(args) != 1 {
		return wrongArgCountError(1, len(args))
	}

	switch argObj := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(len(argObj.Value))}
	case *object.Array:
		return &object.Integer{Value: int64(len(argObj.Elements))}
	default:
		return newError("argument to `len` not supported, got %s", args[0].Type())
	}
}

func __print(args ...object.Object) object.Object {
	// TODO: string formatting
	for _, arg := range args {
		fmt.Println(arg.Inspect())
	}
	return NULL
}

func __first(args ...object.Object) object.Object {
	if len(args) != 1 {
		return wrongArgCountError(1, len(args))
	}

	switch argObj := args[0].(type) {
	case *object.String:
		return &object.String{Value: string(argObj.Value[0])}
	case *object.Array:
		return argObj.Elements[0]
	default:
		return newError("argument to `len` not supported, got %s", args[0].Type())
	}
}
