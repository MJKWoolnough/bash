package bash

// File automatically generated with format.sh.

import "fmt"

// Type is an interface satisfied by all bash structural types.
type Type interface {
	fmt.Formatter
	bashType()
}

func (Tokens) bashType() {}

func (ArithmeticExpansion) bashType() {}

func (Assignment) bashType() {}

func (Command) bashType() {}

func (CommandOrControl) bashType() {}

func (CommandSubstitution) bashType() {}

func (Control) bashType() {}

func (File) bashType() {}

func (Line) bashType() {}

func (Parameter) bashType() {}

func (ParameterAssign) bashType() {}

func (ParameterExpansion) bashType() {}

func (Pipeline) bashType() {}

func (Redirection) bashType() {}

func (Statement) bashType() {}

func (String) bashType() {}

func (Value) bashType() {}

func (Word) bashType() {}

func (WordOrOperator) bashType() {}

func (WordOrToken) bashType() {}

func (WordPart) bashType() {}
