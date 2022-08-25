package evaluator

import (
	"fmt"

	"github.com/oteto/gonkey/pkg/object"
)

const (
	BUILTIN_NUMBER_OF_ARGUMENT_ERROR = "wrong number of arguments. got=%d, want=%d"
	BUILTIN_ARGUMENT_TYPE_ERRROR     = "argument to `%s` not supported, got %s"
)

var builtins = map[string]*object.Builtin{
	"len":   {Fn: builtinLen},
	"first": {Fn: builtinFirst},
	"last":  {Fn: builtinLast},
	"rest":  {Fn: builtinRest},
	"push":  {Fn: builtinPush},
	"puts":  {Fn: builtinPuts},
}

func builtinPuts(args ...object.Object) object.Object {
	for _, arg := range args {
		fmt.Println(arg.Inspect())
	}
	return NULL
}

func builtinPush(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError(BUILTIN_NUMBER_OF_ARGUMENT_ERROR, len(args), 2)
	}
	arr, ok := args[0].(*object.Array)
	if !ok {
		return newError("first "+BUILTIN_ARGUMENT_TYPE_ERRROR, "push", args[0].Type())
	}
	newElements := make([]object.Object, len(arr.Elements)+1)
	copy(newElements, arr.Elements)
	newElements[len(arr.Elements)] = args[1]
	return &object.Array{Elements: newElements}
}

func builtinRest(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError(BUILTIN_NUMBER_OF_ARGUMENT_ERROR, len(args), 1)
	}
	arr, ok := args[0].(*object.Array)
	if !ok {
		return newError(BUILTIN_ARGUMENT_TYPE_ERRROR, "last", args[0].Type())
	}
	if len(arr.Elements) < 1 {
		return NULL
	}
	newElements := make([]object.Object, len(arr.Elements)-1)
	copy(newElements, arr.Elements[1:])
	return &object.Array{Elements: newElements}
}

func builtinLast(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError(BUILTIN_NUMBER_OF_ARGUMENT_ERROR, len(args), 1)
	}
	arr, ok := args[0].(*object.Array)
	if !ok {
		return newError(BUILTIN_ARGUMENT_TYPE_ERRROR, "last", args[0].Type())
	}
	if len(arr.Elements) < 1 {
		return NULL
	}
	return arr.Elements[len(arr.Elements)-1]
}

func builtinFirst(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError(BUILTIN_NUMBER_OF_ARGUMENT_ERROR, len(args), 1)
	}
	arr, ok := args[0].(*object.Array)
	if !ok {
		return newError(BUILTIN_ARGUMENT_TYPE_ERRROR, "first", args[0].Type())
	}
	if len(arr.Elements) < 1 {
		return NULL
	}
	return arr.Elements[0]
}

func builtinLen(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError(BUILTIN_NUMBER_OF_ARGUMENT_ERROR, len(args), 1)
	}
	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	case *object.Array:
		return &object.Integer{Value: int64(len(arg.Elements))}
	default:
		return newError(BUILTIN_ARGUMENT_TYPE_ERRROR, "len", arg.Type())
	}
}
