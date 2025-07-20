package bash

// File automatically generated with format.sh.

import "fmt"

// Format implements the fmt.Formatter interface.
func (f ArithmeticExpansion) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = ArithmeticExpansion
		type ArithmeticExpansion X

		fmt.Fprintf(s, "%#v", ArithmeticExpansion(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface.
func (f ArrayWord) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = ArrayWord
		type ArrayWord X

		fmt.Fprintf(s, "%#v", ArrayWord(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface.
func (f Assignment) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = Assignment
		type Assignment X

		fmt.Fprintf(s, "%#v", Assignment(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface.
func (f AssignmentOrWord) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = AssignmentOrWord
		type AssignmentOrWord X

		fmt.Fprintf(s, "%#v", AssignmentOrWord(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface.
func (f BraceExpansion) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = BraceExpansion
		type BraceExpansion X

		fmt.Fprintf(s, "%#v", BraceExpansion(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface.
func (f BraceWord) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = BraceWord
		type BraceWord X

		fmt.Fprintf(s, "%#v", BraceWord(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface.
func (f CaseCompound) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = CaseCompound
		type CaseCompound X

		fmt.Fprintf(s, "%#v", CaseCompound(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface.
func (f Command) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = Command
		type Command X

		fmt.Fprintf(s, "%#v", Command(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface.
func (f CommandOrCompound) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = CommandOrCompound
		type CommandOrCompound X

		fmt.Fprintf(s, "%#v", CommandOrCompound(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface.
func (f CommandSubstitution) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = CommandSubstitution
		type CommandSubstitution X

		fmt.Fprintf(s, "%#v", CommandSubstitution(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface.
func (f Compound) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = Compound
		type Compound X

		fmt.Fprintf(s, "%#v", Compound(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface.
func (f File) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = File
		type File X

		fmt.Fprintf(s, "%#v", File(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface.
func (f ForCompound) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = ForCompound
		type ForCompound X

		fmt.Fprintf(s, "%#v", ForCompound(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface.
func (f FunctionCompound) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = FunctionCompound
		type FunctionCompound X

		fmt.Fprintf(s, "%#v", FunctionCompound(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface.
func (f GroupingCompound) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = GroupingCompound
		type GroupingCompound X

		fmt.Fprintf(s, "%#v", GroupingCompound(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface.
func (f Heredoc) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = Heredoc
		type Heredoc X

		fmt.Fprintf(s, "%#v", Heredoc(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface.
func (f HeredocPartOrWord) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = HeredocPartOrWord
		type HeredocPartOrWord X

		fmt.Fprintf(s, "%#v", HeredocPartOrWord(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface.
func (f IfCompound) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = IfCompound
		type IfCompound X

		fmt.Fprintf(s, "%#v", IfCompound(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface.
func (f Line) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = Line
		type Line X

		fmt.Fprintf(s, "%#v", Line(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface.
func (f LoopCompound) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = LoopCompound
		type LoopCompound X

		fmt.Fprintf(s, "%#v", LoopCompound(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface.
func (f Parameter) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = Parameter
		type Parameter X

		fmt.Fprintf(s, "%#v", Parameter(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface.
func (f ParameterAssign) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = ParameterAssign
		type ParameterAssign X

		fmt.Fprintf(s, "%#v", ParameterAssign(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface.
func (f ParameterExpansion) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = ParameterExpansion
		type ParameterExpansion X

		fmt.Fprintf(s, "%#v", ParameterExpansion(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface.
func (f Pattern) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = Pattern
		type Pattern X

		fmt.Fprintf(s, "%#v", Pattern(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface.
func (f PatternLines) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = PatternLines
		type PatternLines X

		fmt.Fprintf(s, "%#v", PatternLines(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface.
func (f Pipeline) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = Pipeline
		type Pipeline X

		fmt.Fprintf(s, "%#v", Pipeline(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface.
func (f Redirection) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = Redirection
		type Redirection X

		fmt.Fprintf(s, "%#v", Redirection(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface.
func (f SelectCompound) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = SelectCompound
		type SelectCompound X

		fmt.Fprintf(s, "%#v", SelectCompound(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface.
func (f Statement) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = Statement
		type Statement X

		fmt.Fprintf(s, "%#v", Statement(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface.
func (f String) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = String
		type String X

		fmt.Fprintf(s, "%#v", String(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface.
func (f TestCompound) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = TestCompound
		type TestCompound X

		fmt.Fprintf(s, "%#v", TestCompound(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface.
func (f TestConsequence) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = TestConsequence
		type TestConsequence X

		fmt.Fprintf(s, "%#v", TestConsequence(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface.
func (f Tests) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = Tests
		type Tests X

		fmt.Fprintf(s, "%#v", Tests(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface.
func (f Value) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = Value
		type Value X

		fmt.Fprintf(s, "%#v", Value(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface.
func (f Word) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = Word
		type Word X

		fmt.Fprintf(s, "%#v", Word(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface.
func (f WordOrOperator) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = WordOrOperator
		type WordOrOperator X

		fmt.Fprintf(s, "%#v", WordOrOperator(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface.
func (f WordOrToken) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = WordOrToken
		type WordOrToken X

		fmt.Fprintf(s, "%#v", WordOrToken(f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface.
func (f WordPart) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = WordPart
		type WordPart X

		fmt.Fprintf(s, "%#v", WordPart(f))
	} else {
		format(&f, s, v)
	}
}
