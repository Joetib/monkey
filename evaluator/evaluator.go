package evaluator

import (
	"fmt"
	"io/ioutil"
	"monkey/ast"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
)

var (
	//TRUE single Boolean object for true
	//This is intended to reduce recreating
	//more such boolean objects hence reducing space
	TRUE = &object.Boolean{Value: true}
	//FALSE single boolean object for false
	FALSE = &object.Boolean{Value: false}

	//NULL single Null object to reduce creating more in memory
	NULL = &object.Null{}
)

//newError constructs a new error object
func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

//isError checks if an object is an error
func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

//Eval main evaluator function
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.ImportStatement:
		module := evalImportStatement(node.Value.String(), node.Alias, env)
		if isError(module) {
			return module
		}
		moduleObj, ok := module.(*object.Module)
		if ok {
			env.Set(moduleObj.Name, module)
		}

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		if node.Property != nil {
			property, ok := node.Property.(*ast.Identifier)
			if !ok {
				return newError("Error converting node.Property into an *ast.Identifier")
			}
			cls, ok := env.Get(node.Name.Value)
			if !ok {
				return newError("unknown identifier: %s", node.Name.Value)
			}
			classInstance, ok := cls.(*object.ClassInstance)
			if !ok {
				return newError("Dot assignment allowed only on class Instances.Cannot use dot assignment on %T: %s", cls, node.Name.Value)
			}
			_, ok = classInstance.Env.Get(property.Value)
			if !ok {
				return newError("%s is not an instance variable of class %s", property.Value, cls.Inspect())
			}
			classInstance.Env.Set(property.Value, val)
		} else {
			env.Set(node.Name.Value, val)
		}
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right, env)
	case *ast.InfixExpression:
		if node.Operator == "." {
			left := Eval(node.Left, env)
			if isError(left) {
				return left
			}
			return evalDotInfixExpression(left, node.Right, env)
		}
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.WhileExpression:
		return evalWhileExpression(node, env)
	case *ast.FunctionLiteral:
		return &object.Function{Parameters: node.Parameters, Env: env, Body: node.Body}
	case *ast.ClassStatement:
		newEnv := object.NewEnclosedEnvironment(env)
		for _, value := range node.Parents {
			pResult := Eval(value, env)
			cls, ok := pResult.(*object.Class)
			if !ok {
				return &object.Error{Message: fmt.Sprintf("parent to be inherited from must be a class. got %T", pResult)}
			}
			newEnv.ShallowCopy(cls.Env)
		}

		// Let every statement in the block get its environment from outside the class,
		// this hides instance variables and methods
		evalClassBlockStatement(node.Body, newEnv)

		class := &object.Class{Name: node.Name.String(), Env: newEnv}
		env.Set(class.Name, class)
		return class
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	}
	return nil
}

//evalImportStatement evaluates an import statement
func evalImportStatement(Name string, alias ast.Expression, env *object.Environment) object.Object {
	content, err := ioutil.ReadFile(Name + ".monkey")
	if err != nil {
		fmt.Printf("Could not open file : %s", Name)
		panic(err)
	}
	newEnv := object.NewEnvironment()
	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		errorMessage := ""
		for _, er := range p.Errors() {
			errorMessage = errorMessage + er
		}
		return &object.Error{Message: errorMessage}
	}
	Eval(program, newEnv)
	if alias != nil {
		aliasString := alias.(*ast.StringLiteral).String()
		return &object.Module{Env: newEnv, Name: aliasString}

	}
	return &object.Module{Env: newEnv, Name: Name}
}

//evalProgram evaluate a list of statements
func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)
		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	// this line allows return last values in a block of code
	// even though it does not explicitly use the return keyword
	// to avoid this behaviour, and return NULL instead, replace this
	// line with
	// return NULL
	return result
}

//evalBlockStatement evaluates a block of statements.
// this is implemented different from evaluate program so we can handle return statements
// properly. See page 130 of the book `writing an interpreter in go` for explanation
func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}
	// this line allows return last values in a block of code
	// even though it does not explicitly use the return keyword
	// to avoid this behaviour, and return NULL instead, replace this
	// line with
	// return NULL
	return result
}

//evalExpressions evaluates multiple expressions
func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object
	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

//applyFunction runs a function call
func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	//check whether the function is a builtin function first.
	// this is done first so that the end user cannot overide builtin objects
	// if this check is done second, then the user defined objects will have precedence
	// over builtin types.
	// this is a little deviation of my own from original monkey representation
	// where user defined values have more precedence over builtin types
	case *object.Builtin:
		return fn.Fn(args...)
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Class:
		cls := &object.ClassInstance{Name: fn.Name, Env: fn.Env}
		if value, ok := cls.Env.Get("__New__"); ok {
			if value.Type() == object.FUNCTION_OBJ {
				function, _ := value.(*object.Function)
				applyMethod(function, cls, args)
			}
		}
		return cls
	default:
		return newError("not a function: %s", fn.Type())
	}
}

//extendFunctionEnv creates a new running environment for a function that encloses
// the environment where the function was created
func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}
	return env
}

//unwrapReeturnValue unwraps the return value of a function
func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

//evalIdentifer evaluates an identifier by getting its value from the env
func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	//check whether the identifier is a builtin function first.
	// this is done first so that the end user cannot overide builtin objects
	// if this check is done second, then the user defined objects will have precedence
	// over builtin types.
	// this is a little deviation of my own from original monkey representation
	// where user defined values have more precedence over builtin types
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	return newError("identifier not found: " + node.Value)
}

//nativeBoolToBooleanOjbect creates a boolean object
func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

//evalPrefixExpression evaluates a prefix expression
func evalPrefixExpression(operator string, right object.Object, env *object.Environment) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

//evalBangOperatorExpression evaluates a bang operator
func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

//evalIndexExpression evaluates the result of indexing an object
func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

//evalHashLiteral evaluates and create a Hash object
func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)
	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}
		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as a hash key: %s", key.Type())
		}
		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}
		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}
	return &object.Hash{Pairs: pairs}
}

//evalHashIndexExpression evaluates the results of indexing the Hash
func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObect := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}
	pair, ok := hashObect.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}
	return pair.Value
}

//evalArrayIndexExpression evaluate the result of indexing an array
func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value

	max := int64(len(arrayObject.Elements) - 1)
	/*
		this conditional block returns null for indexes greater than the max and less than zero
		an alternate behavior where negative indexes represent counting objects from right to left
		can be implemented by uncommenting the code

		if idx < 0 {
			// to raise an error when the index is out of bounds rather than returning NULL,
			// we raise this error
			// you can uncomment the block if you want such behaviour
			// if -idx > max+1 {
			//	 return newError("index out of bounds : maximum index: %d, requested index : %d", -max-1, idx)
			// }
			idx = -idx -1 // gets the corresponding positive index
		}

	*/

	if idx < 0 || idx > max {
		// this returns NULL to the caller
		// alternatively, we may want to raise an indexoutofbounds exception
		// this is done by
		//return newError("index out of bounds : maximum index: %d, requested index : %d", max, idx)
		return NULL
	}
	return arrayObject.Elements[idx]
}

//evalMinusPrefixOperatorExpression evaluates negating a value
func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

//evalInfixExpression evaluates infix operations
func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(operator, left, right)

	case (left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ) || (right.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ):
		return evalFloatIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

//evalDotInfixOperation evaluates a method
func evalDotInfixExpression(left object.Object, right ast.Node, env *object.Environment) object.Object {
	switch left.(type) {
	case *object.ClassInstance:
		left, ok := left.(*object.ClassInstance)
		if !ok {
			return nil
		}
		return evalClassDotOperation(left, right, env)
	case *object.Module:
		left, ok := left.(*object.Module)
		if !ok {
			return nil
		}
		return evalModuleDotOperation(left, right, env)
	default:
		return &object.Error{Message: fmt.Sprintf("Dot operation not supported for %T", left.Type())}
	}
}

//evalClassBlockStatement evaluates a block of statements.
// this is implemented different from evaluate program so we can handle return statements
// properly. See page 130 of the book `writing an interpreter in go` for explanation
func evalClassBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
			if rt == object.FUNCTION_OBJ {
				fn, ok := result.(*object.Function)

				if ok {
					// Set the outer Environment to the outside of class
					// this way, class methods only see what is outside the class
					// not what is inside the class
					// This means that class methods can only access instance or class
					// properties from the self variable
					fn.Env.SetOuter(env.GetOuter())
				}
			}
		}
	}
	// this line allows return last values in a block of code
	// even though it does not explicitly use the return keyword
	// to avoid this behaviour, and return NULL instead, replace this
	// line with
	// return NULL
	return result
}

//evalModuleDotOperator evaluates dot operation between and object
func evalModuleDotOperation(left *object.Module, right ast.Node, env *object.Environment) object.Object {
	switch right.(type) {
	case *ast.CallExpression:
		right, ok := right.(*ast.CallExpression)
		if !ok {
			return nil
		}
		function := Eval(right.Function, left.Env.Closed())
		if isError(function) {
			return function
		}
		args := evalExpressions(right.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)

	case *ast.Identifier:
		return Eval(right, left.Env.Closed())
	default:
		return &object.Error{Message: "Cannot perform Dot operation"}
	}

}

//evalClassDotOperator evaluates dot operation between and object
func evalClassDotOperation(left *object.ClassInstance, right ast.Node, env *object.Environment) object.Object {
	switch right.(type) {
	case *ast.CallExpression:
		right, ok := right.(*ast.CallExpression)
		if !ok {
			return nil
		}
		function := Eval(right.Function, left.Env.Closed())
		if isError(function) {
			return function
		}
		args := evalExpressions(right.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyMethod(function, left, args)
	case *ast.Identifier:
		return Eval(right, left.Env.Closed())
	default:
		return &object.Error{Message: "Cannot perform Dot operation"}
	}
}

//applyMethod runs a function call
func applyMethod(fn object.Object, left object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	//check whether the function is a builtin function first.
	// this is done first so that the end user cannot overide builtin objects
	// if this check is done second, then the user defined objects will have precedence
	// over builtin types.
	// this is a little deviation of my own from original monkey representation
	// where user defined values have more precedence over builtin types

	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		class, ok := left.(*object.ClassInstance)
		if ok {
			extendedEnv.Set("self", class)
		}
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	default:
		return newError("not a function: %s", fn.Type())
	}
}

//evalIntegerInfixExpression evaluates infix expressions where both operands are integers
func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Float{Value: float64(leftVal) / float64(rightVal)}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

//evalFloatInfixExpression evaluates infix expressions where both operands are floats
func evalFloatInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value
	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		return &object.Float{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

//evalFloatInfixExpression evaluates infix expressions where both operands are floats
func evalFloatIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	var leftVal float64
	var rightVal float64
	if left.Type() == object.INTEGER_OBJ {
		leftVal = float64(left.(*object.Integer).Value)
	} else {
		leftVal = left.(*object.Float).Value
	}
	if right.Type() == object.INTEGER_OBJ {
		rightVal = float64(right.(*object.Integer).Value)
	} else {
		rightVal = right.(*object.Float).Value
	}
	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		return &object.Float{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

//evalStringInfixExpression evaluates infix operations involving strings
func evalStringInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}

}

//evalIfExpression evaluates a conditional if-else-statement
func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

//evalIfExpression evaluates a conditional if-else-statement
func evalWhileExpression(we *ast.WhileExpression, env *object.Environment) object.Object {
	var result object.Object
	result = NULL
	for {
		condition := Eval(we.Condition, env)

		if isError(condition) {
			return condition
		}
		if isTruthy(condition) {
			result = Eval(we.Consequence, env)
			if result == nil {
				result = NULL
			}
		} else {
			return result
		}
	}
}

// Helpers

//isTruthy checks whether an object can be considered true
func isTruthy(obj object.Object) bool {
	//all values are true except NULL and FALSE
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}
