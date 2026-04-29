package walk_test

import (
	"fmt"

	"vimagination.zapto.org/bash"
	"vimagination.zapto.org/bash/walk"
	"vimagination.zapto.org/parser"
)

func Example() {
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
