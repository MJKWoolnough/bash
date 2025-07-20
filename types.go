package bash

// File automatically generated with format.sh.

import "fmt"

// Type is an interface satisfied by all bash structural types.
type Type interface {
	fmt.Formatter
	bashType()
}

func (ArithmeticExpansion) bashType() {}

func (ArrayWord) bashType() {}

func (Assignment) bashType() {}

func (AssignmentOrWord) bashType() {}

func (BraceExpansion) bashType() {}

func (BraceWord) bashType() {}

func (CaseCompound) bashType() {}

func (Command) bashType() {}

func (CommandOrCompound) bashType() {}

func (CommandSubstitution) bashType() {}

func (Compound) bashType() {}

func (File) bashType() {}

func (ForCompound) bashType() {}

func (FunctionCompound) bashType() {}

func (GroupingCompound) bashType() {}

func (Heredoc) bashType() {}

func (HeredocPartOrWord) bashType() {}

func (IfCompound) bashType() {}

func (Line) bashType() {}

func (LoopCompound) bashType() {}

func (Parameter) bashType() {}

func (ParameterAssign) bashType() {}

func (ParameterExpansion) bashType() {}

func (Pattern) bashType() {}

func (PatternLines) bashType() {}

func (Pipeline) bashType() {}

func (Redirection) bashType() {}

func (SelectCompound) bashType() {}

func (Statement) bashType() {}

func (String) bashType() {}

func (TestCompound) bashType() {}

func (TestConsequence) bashType() {}

func (Tests) bashType() {}

func (Value) bashType() {}

func (Word) bashType() {}

func (WordOrOperator) bashType() {}

func (WordOrToken) bashType() {}

func (WordPart) bashType() {}
