package interpreter

import (
	"bufio"
	"io"

	"github.com/jamestrew/go-interpreter/monkey/evaluator"
	"github.com/jamestrew/go-interpreter/monkey/object"
	"github.com/jamestrew/go-interpreter/monkey/parser"
)

func PrintParseErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t")
		io.WriteString(out, msg)
		io.WriteString(out, "\n")
	}
}

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	eval := evaluator.New(object.NewEnvironment())

	for {
		if !scanner.Scan() {
			return
		}
		line := scanner.Text()
		program, p := parser.ParseInput(line)
		if len(p.Errors()) != 0 {
			PrintParseErrors(out, p.Errors())
			return
		}

		eval.Eval(program)
	}
}
