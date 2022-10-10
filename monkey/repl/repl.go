package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/jamestrew/go-interpreter/monkey/evaluator"
	"github.com/jamestrew/go-interpreter/monkey/lexer"
	"github.com/jamestrew/go-interpreter/monkey/object"
	"github.com/jamestrew/go-interpreter/monkey/parser"
)

const PROMPT = ">> "

func printParseErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t")
		io.WriteString(out, msg)
		io.WriteString(out, "\n")
	}
}

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	eval := evaluator.New(env)

	for {
		fmt.Printf(PROMPT)
		if !scanner.Scan() {
			return
		}

		line := scanner.Text()
		p := parser.New(lexer.New(line))

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParseErrors(out, p.Errors())
			continue
		}

		evaluated := eval.Eval(program)
		io.WriteString(out, evaluated.Inspect())
		io.WriteString(out, "\n")
	}
}
