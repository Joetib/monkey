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