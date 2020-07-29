package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Monkey programming language!\n",
		user.Username)
	fmt.Printf("Feel free to type in commands\n")
	if len(os.Args) > 1 {
		filename := os.Args[1]
		content, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Printf("Could not open file : %s", filename)
			panic(err)
		}
		env := object.NewEnvironment()

		l := lexer.New(string(content))
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(os.Stdout, p.Errors())
			return
		}

		evaluated := evaluator.Eval(program, env)
		Err, ok := evaluated.(*object.Error)
		if ok {
			io.WriteString(os.Stdout, repl.MONKEY_FACE)
			io.WriteString(os.Stdout, "\n\n>> Error Running program : ")
			io.WriteString(os.Stdout, "\n          >>   "+Err.Inspect()+"\n")
			return
		}
		if evaluated != nil {

			io.WriteString(os.Stdout, evaluated.Inspect())
			io.WriteString(os.Stdout, "\n")
		}

	} else {
		repl.Start(os.Stdin, os.Stdout)
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, ">>> Error >>> .................")
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
