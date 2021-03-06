package evaluator

import (
	"fmt"
	"monkey/object"
)

var builtins = map[string]*object.Builtin{
	"str": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. to str got=%d, want=1", len(args))
			}

			return &object.String{Value: args[0].Inspect()}

		},
	},
	"env": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. to str got=%d, want=1", len(args))
			}
			obj := args[0]
			switch obj.(type) {
			case *object.ClassInstance:
				newObj, _ := obj.(*object.ClassInstance)
				fmt.Println(newObj.Env)
			case *object.Class:
				newObj, _ := obj.(*object.Class)
				fmt.Println(newObj.Env)
			case *object.Module:
				newObj, _ := obj.(*object.Module)
				fmt.Println(newObj.Env)
			default:
				fmt.Println("Error : >>> Object has no Env")
			}
			return NULL
		},
	},
	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"first": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}
			return NULL
		},
	},
	"last": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `last` must be ARRAY, got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)

			if length > 0 {
				return arr.Elements[length-1]
			}
			return NULL
		},
	},
	"rest": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `rest` must be ARRAY, got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				newElements := make([]object.Object, length-1, length-1)
				copy(newElements, arr.Elements[1:length])
				return &object.Array{Elements: newElements}
			}
			return NULL
		},
	},
	"push": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("first argument to `push` must be ARRAY, got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)

			// create and return a new array so that modification to the new array
			// does not affect the old array.

			newElements := make([]object.Object, length+1, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]
			return &object.Array{Elements: newElements}
		},
	},
	"puts": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NULL
		},
	},
	"hasattr": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			left, ok := args[0].(*object.ClassInstance)
			if !ok {
				return newError("Argument 1 must be a `CLASSINSTANCE`")
			}
			right, ok := args[0].(*object.String)
			if !ok {
				return newError("Argument 2 must be a `STRING`")
			}
			_, ok = left.Env.Closed().Get(right.Inspect())
			if ok {
				return TRUE
			}
			return FALSE

		},
	},
	"setattr": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 3 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			left, ok := args[0].(*object.ClassInstance)
			if !ok {
				return newError("Argument 1 must be a `CLASSINSTANCE`")
			}
			right, ok := args[1].(*object.String)
			if !ok {
				return newError("Argument 2 must be a `STRING`")
			}
			left.Env.Set(right.Inspect(), args[2])
			fmt.Println(left.Env.Get(right.Inspect()))
			return args[2]

		},
	},
}
