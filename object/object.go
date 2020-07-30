package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"monkey/ast"
	"strings"
)

type ObjectType string

const (
	INTEGER_OBJ       = "INTEGER"
	FLOAT_OBJ         = "FLOAT"
	BOOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ          = "NULL"
	RETURN_VALUE_OBJ  = "RETURN_VALUE"
	FUNCTION_OBJ      = "FUNCTION"
	STRING_OBJ        = "STRING"
	ERROR_OBJ         = "ERROR"
	BUILTIN_OBJ       = "BUILTIN"
	ARRAY_OBJ         = "ARRAY"
	HASH_OBJ          = "HASH"
	CLASS_OBJ         = "CLASS"
	CLASSINSTANCE_OBJ = "CLASS_INSTANCE"
	MODULE_OBJ        = "MODULE"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

//BuiltinFunction a function representation of builtin functions
type BuiltinFunction func(args ...Object) Object

//Hashable interface for object that implement a HashKey function
type Hashable interface {
	HashKey() HashKey
}

//Integer object to hold integers
type Integer struct {
	Value int64
}

//Inspect returns a string representation of the object
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

//Type returns the type of the object
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

//HashKey function to generate a HashKey object from a Integer
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: float64(i.Value)}
}

//Float object to hold integers
type Float struct {
	Value float64
}

//Inspect returns a string representation of the object
func (f *Float) Inspect() string { return fmt.Sprintf("%f", f.Value) }

//Type returns the type of the object
func (f *Float) Type() ObjectType { return FLOAT_OBJ }

//HashKey function to generate a HashKey object from a Integer
func (f *Float) HashKey() HashKey {
	return HashKey{Type: f.Type(), Value: float64(f.Value)}
}

//Boolean object to hold boolean values true and false
type Boolean struct {
	Value bool
}

//Inspect returns a string representation of the object
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

//Type returns the type of the object
func (b *Boolean) Type() ObjectType { return BOOOLEAN_OBJ }

//HashKey function to generate a HashKey object from a boolean
func (b *Boolean) HashKey() HashKey {
	var value float64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}

//Null struct representation of null values
type Null struct{}

//Inspect returns a string representation of the object
func (n *Null) Inspect() string { return "null" }

//Type returns the type of the object
func (n *Null) Type() ObjectType { return NULL_OBJ }

//ReturnValue object for holding return values
type ReturnValue struct {
	Value Object
}

//Inspect returns a string representation of the object
func (rv *ReturnValue) Inspect() string { return rv.Value.Inspect() }

//Type returns the type of the object
func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }

//Error base error object
type Error struct {
	Message string
}

//Inspect returns a string representation of the object
func (e *Error) Inspect() string { return "Error: " + e.Message }

//Type returns the type of the object
func (e *Error) Type() ObjectType { return ERROR_OBJ }

//Function type for functions
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

//Inspect returns a string representation of the object
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n")

	return out.String()
}

//Type returns the type of the object
func (f *Function) Type() ObjectType { return FUNCTION_OBJ }

//String object representation of a string
type String struct {
	Value string
}

//Type returns the type of the object
func (s *String) Type() ObjectType { return STRING_OBJ }

//Inspect returns the value of the object
func (s *String) Inspect() string { return s.Value }

//HashKey function to generate a HashKey object from a boolean
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: float64(h.Sum64())}
}

//Builtin node representation of all builtin functions and objects
type Builtin struct {
	Fn BuiltinFunction
}

//Type returns the type of the object
func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }

//Inspect returns a string representation of the node
func (b *Builtin) Inspect() string { return "builtin function" }

//Array node representation of array objects
type Array struct {
	Elements []Object
}

//Type returns the type of the object
func (ao *Array) Type() ObjectType { return ARRAY_OBJ }

//Inspect returns a string representation of the node
func (ao *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

//HashKey node representation of hash objects keys
type HashKey struct {
	Type  ObjectType
	Value float64
}

//HashPair a pair of entries in a hashmap
type HashPair struct {
	Key   Object
	Value Object
}

//Hash a hashmap/dictionary/map implementation
type Hash struct {
	Pairs map[HashKey]HashPair
}

//Type returns the type of the object
func (h *Hash) Type() ObjectType { return HASH_OBJ }

//Inspect returns a string representation of the Hash object
func (h *Hash) Inspect() string {
	var out bytes.Buffer
	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

//Class Base handler for class
type Class struct {
	Name string
	Env  *Environment
}

//Type returns the type of the object
func (C *Class) Type() ObjectType { return CLASS_OBJ }

//Inspect returns a string representation of the node
func (C *Class) Inspect() string { return "class " + C.Name }

//ClassInstance an instance of a class
type ClassInstance struct {
	Name string
	Env  *Environment
}

//Type returns the type of the object
func (Ci *ClassInstance) Type() ObjectType { return CLASSINSTANCE_OBJ }

//Inspect returns a string representation of the node
func (Ci *ClassInstance) Inspect() string { return "<Instance of Class " + Ci.Name + ">" }

//Module Base handler for class
type Module struct {
	Name string
	Env  *Environment
}

//Type returns the type of the object
func (M *Module) Type() ObjectType { return MODULE_OBJ }

//Inspect returns a string representation of the node
func (M *Module) Inspect() string { return "module " + M.Name }
