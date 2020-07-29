# Monkey

The monkey interpreter from reading the book `writing an interpreter in go` by Thorsten Ball.
## Features
- First class Functions
- Recursion
- Strings
- HashMaps
- Lists (dynamic arrays)
- Integers
- Modules
- Classes and Objects

## Running the code
To run the repl with `go` installed, just run `go run main.go` and you're good to go


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