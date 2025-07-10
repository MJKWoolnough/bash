package main

import (
	"fmt"
	"os"

	"vimagination.zapto.org/bash"
	"vimagination.zapto.org/parser"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func run() error {
	r := os.Stdin

	if len(os.Args) > 1 {
		f, err := os.Open(os.Args[1])
		if err != nil {
			return nil
		}

		r = f
	}

	tk := parser.NewReaderTokeniser(r)

	b, err := bash.Parse(bash.SetTokeniser(&tk))
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%+s", b)

	return nil
}
