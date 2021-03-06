package object

//Environment Holds environments in the interpreter
type Environment struct {
	store map[string]Object
	outer *Environment
}

//SetOuter sets a value for the outer field of an environment
func (e *Environment) SetOuter(env *Environment) {
	e.outer = env
}

//GetOuter returns the value of the outer field of an environment
func (e *Environment) GetOuter() *Environment {
	return e.outer
}

//shallowCopy  copies the values in an environment to another
func (e *Environment) ShallowCopy(env *Environment) *Environment {
	for key, value := range env.store {
		if _, ok := e.Get(key); ok {
			continue
		}
		e.Set(key, value)
	}
	return e
}

//Get gets the value associated with a key in the environment store
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

//Closed makes a copy of the environment excluding the outer
func (e *Environment) Closed() *Environment {
	return &Environment{store: e.store, outer: nil}
}

//Set sets an entry to the environment's store
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

//SetMultiple sets multiple key values to evironment's store
func (e *Environment) SetMultiple(values map[string]Object) *Environment {
	for key, value := range values {
		e.store[key] = value
	}
	return e
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
