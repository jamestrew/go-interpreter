package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"

	"github.com/jamestrew/go-interpreter/monkey/interpreter"
	"github.com/jamestrew/go-interpreter/monkey/repl"
)

func startRepl() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s --- Let's get monkey\n", user.Username)
	repl.Start(os.Stdin, os.Stdout)
}

func execFile(filePath string) {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	interpreter.Start(f, os.Stdout)
}

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		startRepl()
	} else {
		execFile(flag.Arg(0))
	}
}
