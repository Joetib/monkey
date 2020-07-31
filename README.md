# Monkey

The monkey interpreter from reading the book `writing an interpreter in go` by Thorsten Ball.
## Features
- First class Functions
- Recursion
- Strings
- HashMaps
- Lists (dynamic arrays)
- Integers (int64)
- Floats (float64)
- Modules and import mechanism
- Classes and Objects
- Multiple inheritance
- while loop

# TODO
- Keyword arguments
- Default values for arguments
- More informative error messages with positions

## Running the code
----------------------------
To run the repl with `go` installed, just run `go run main.go` and you're good to go

or build an executable by running `go build .`
Then you can use the executable to run monkey applications
```
monkey filename.monkey
``` 
replacing `filename` with the name of the file


## examples

```
lef factorial = fn (n) {
    if (n < 2){
        return 1;
    }
    return n * factorial(n-1);
}
puts(factorial(2))
```

# Tutorials
## Basic Data Type
Data types in Monkey are not declared when writing code.
Variables are also dynamic, hence a variable that initially holds a float, may at another point in your program hold a string, or a function, or an object.

Variables are declared as 

```
let VariableName = value;
```
* An `int` is just a number without decimal points. These are int64 values

* A `float` is a number with a decimal point. these are  float64 values
* A `string` is a sequence of characters enclosed in quotations
eg. `let myString = "Kofi is a boy";`
