# bash

[![CI](https://github.com/MJKWoolnough/bash/actions/workflows/go-checks.yml/badge.svg)](https://github.com/MJKWoolnough/bash/actions)
[![Go Reference](https://pkg.go.dev/badge/vimagination.zapto.org/bash.svg)](https://pkg.go.dev/vimagination.zapto.org/bash)
[![Go Report Card](https://goreportcard.com/badge/vimagination.zapto.org/bash)](https://goreportcard.com/report/vimagination.zapto.org/bash)

--
    import "vimagination.zapto.org/bash"

Package bash implements a bash tokeniser and AST.

## Highlights

 - Parse Bash code into AST.
 - Modify parsed code.
 - Consistant bash formatting.

## Usage

```go
package main

import (
	"fmt"

	"vimagination.zapto.org/bash"
	"vimagination.zapto.org/parser"
)

func main() {
	src := `for name in Alice Bob Charlie; do echo "Hello, $name";done`

	tk := parser.NewStringTokeniser(src)

	ast, err := bash.Parse(&tk)
	if err != nil {
		fmt.Println(err)

		return
	}

	ast.Lines[0].Statements[0].Pipeline.CommandOrCompound.Compound.ForCompound.File.Lines[0].Statements[0].Pipeline.CommandOrCompound.Command.AssignmentsOrWords[1].Word.Parts[0].Part.Data = `"Hi, `

	fmt.Printf("%s", ast)

	// Output:
	// for name in Alice Bob Charlie; do
	// 	echo "Hi, $name";
	// done;
}
```

## Documentation

Full API docs can be found at:

https://pkg.go.dev/vimagination.zapto.org/bash
