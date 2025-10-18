package bash_test

import (
	"fmt"

	"vimagination.zapto.org/bash"
	"vimagination.zapto.org/parser"
)

func Example() {
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
