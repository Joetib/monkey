package evaluator

import (
	"monkey/object"
)

var OBJECT = &object.ClassInstance{
	Name: "BuiltinObject",
	Env: object.NewEnvironment().SetMultiple(
		map[string]object.Object{
			"__str__": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("__str__ takes only one argument, but more than one was given")
					}
					return &object.String{Value: args[0].Inspect()}
				},
			},
			"__repr__": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) != 1 {
						return newError("__repr__ takes only one argument, but more than one was given")
					}
					return &object.String{Value: args[0].Inspect()}
				},
			},
		},
	),
} /*
var REFLECTOR = &object.ClassInstance{
	Name: "BuiltinReflector",
	Env: object.NewEnvironment().SetMultiple(
		map[string]object.Object{
			"get_method": &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					if len(args) < 2 {
						return newError("__str__ takes only one argument, but more than one was given")
					}
					httpServerObject := http.ReadRequest
					argObject, ok := args[0].(*object.ClassInstance)
					e := reflect.ValueOf(http)

					if !ok {
						return newError("The first argument to reflector must be a class Instance")
					}
					argObject.Env.Set("__value__", httpServerObject)
					return &object.String{Value: args[0].Inspect()}
				},
			},
		},
	),
}

type Reflector struct {
}
*/
