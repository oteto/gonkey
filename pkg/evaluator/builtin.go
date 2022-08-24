package evaluator

import "github.com/oteto/gonkey/pkg/object"

const (
	BUILTIN_NUMBER_OF_ARGUMENT_ERROR = "wrong number of arguments. got=%d, want=%d"
	BUILTIN_ARGUMENT_TYPE_ERRROR     = "argument to `%s` not supported, got %s"
)

var builtins = map[string]*object.Builtin{
	"len": {Fn: builtinLen},
}

func builtinLen(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError(BUILTIN_NUMBER_OF_ARGUMENT_ERROR, len(args), 1)
	}
	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	default:
		return newError(BUILTIN_ARGUMENT_TYPE_ERRROR, "len", arg.Type())
	}
}
