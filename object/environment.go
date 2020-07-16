package object

//Environment Holds environments in the interpreter
type Environment struct {
	store map[string]Object
	outer *Environment
}

//Get gets the value associated with a key in the environment store
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

//Set sets an entry to the environment's store
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

//NewEnvironment Create and returns a new environment
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

//NewEnclosedEnvironment creates a new environment that contains a parent environment in it
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}
