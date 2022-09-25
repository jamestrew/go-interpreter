package repl

import (
	"bufio"
	"fmt"
	"github.com/jamestrew/go-interpreter/monkey/lexer"
	"github.com/jamestrew/go-interpreter/monkey/token"
	"io"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		if !scanner.Scan() {
			return
		}

		line := scanner.Text()
		lexer := lexer.New(line)
		for tok := lexer.NextToken(); tok.Type != token.EOF; tok = lexer.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
