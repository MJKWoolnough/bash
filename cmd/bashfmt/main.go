package main

import (
	"flag"
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
	var write, concise bool

	flag.BoolVar(&write, "w", false, "write formatted bash code to source file instead of stdout")
	flag.BoolVar(&concise, "c", false, "print concise bash")
	flag.Parse()

	file := flag.CommandLine.Arg(0)

	r := os.Stdin

	if file != "" {
		f, err := os.Open(file)
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

	out := os.Stdout

	if write && r != os.Stdin {
		f, err := os.Create(file)
		if err != nil {
			return nil
		}
		defer f.Close()

		out = f
	}

	if concise {
		fmt.Fprintf(out, "%s", b)
	} else {
		fmt.Fprintf(out, "%+s", b)
	}

	return nil
}
