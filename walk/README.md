# walk

[![CI](https://github.com/MJKWoolnough/bash/actions/workflows/go-checks.yml/badge.svg)](https://github.com/MJKWoolnough/bash/actions)
[![Go Reference](https://pkg.go.dev/badge/vimagination.zapto.org/bash/walk.svg)](https://pkg.go.dev/vimagination.zapto.org/bash/walk)
[![Go Report Card](https://goreportcard.com/badge/vimagination.zapto.org/bash)](https://goreportcard.com/report/vimagination.zapto.org/bash)

--
    import "vimagination.zapto.org/bash/walk"

Package walk provides a Bash type walker.

## Highlights

 - Simple interface to allow control over walking through parsed Bash.
 - Allows modification to the tree as it's being walked.

## Usage

```go
package main

import (
	"fmt"

	"vimagination.zapto.org/bash"
	"vimagination.zapto.org/bash/walk"
	"vimagination.zapto.org/parser"
)

func main() {
	src := "a='Beep''Boop'\nprint() {\n\techo $a;\n}\n\nprint;"
	tk := parser.NewStringTokeniser(src)

	b, err := bash.Parse(&tk)
	if err != nil {
		fmt.Println(err)

		return
	}

	var walkFn walk.Handler

	walkFn = walk.HandlerFunc(func(t bash.Type) error {
		switch t := t.(type) {
		case *bash.WordPart:
			switch t.Part.Data {
			case "'Beep'":
				t.Part.Data = "'Hello'"
			case "'Boop'":
				t.Part.Data = "', world'"
			case "print":
				t.Part.Data = "do_print"
			}
		case *bash.FunctionCompound:
			t.Identifier.Data = "do_print"
		}

		return walk.Walk(t, walkFn)
	})

	walk.Walk(b, walkFn)

	fmt.Printf("%s", b)

	// Output:
	// a='Hello'', world';
	// do_print() { echo $a; }
	//
	// do_print;
}
```

## Documentation

Full API docs can be found at:

https://pkg.go.dev/vimagination.zapto.org/bash/walk
