package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/jamestrew/go-interpreter/monkey/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s --- Let's get monkey\n", user)
	repl.Start(os.Stdin, os.Stdout)
}
