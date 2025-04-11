package bash

// File automatically generated with format.sh.

import "fmt"

// Format implements the fmt.Formatter interface
func (f ArithmeticExpansion) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = ArithmeticExpansion
		type ArithmeticExpansion X

		fmt.Fprintf(s, "%#v", (f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f Assignment) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = Assignment
		type Assignment X

		fmt.Fprintf(s, "%#v", (f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f Command) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = Command
		type Command X

		fmt.Fprintf(s, "%#v", (f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f CommandOrControl) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = CommandOrControl
		type CommandOrControl X

		fmt.Fprintf(s, "%#v", (f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f CommandSubstitution) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = CommandSubstitution
		type CommandSubstitution X

		fmt.Fprintf(s, "%#v", (f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f Control) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = Control
		type Control X

		fmt.Fprintf(s, "%#v", (f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f File) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = File
		type File X

		fmt.Fprintf(s, "%#v", (f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f Heredoc) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = Heredoc
		type Heredoc X

		fmt.Fprintf(s, "%#v", (f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f HeredocPartOrWord) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = HeredocPartOrWord
		type HeredocPartOrWord X

		fmt.Fprintf(s, "%#v", (f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f Line) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = Line
		type Line X

		fmt.Fprintf(s, "%#v", (f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f Parameter) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = Parameter
		type Parameter X

		fmt.Fprintf(s, "%#v", (f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ParameterAssign) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = ParameterAssign
		type ParameterAssign X

		fmt.Fprintf(s, "%#v", (f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f ParameterExpansion) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = ParameterExpansion
		type ParameterExpansion X

		fmt.Fprintf(s, "%#v", (f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f Pipeline) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = Pipeline
		type Pipeline X

		fmt.Fprintf(s, "%#v", (f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f Redirection) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = Redirection
		type Redirection X

		fmt.Fprintf(s, "%#v", (f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f Statement) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = Statement
		type Statement X

		fmt.Fprintf(s, "%#v", (f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f String) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = String
		type String X

		fmt.Fprintf(s, "%#v", (f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f Value) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = Value
		type Value X

		fmt.Fprintf(s, "%#v", (f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f Word) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = Word
		type Word X

		fmt.Fprintf(s, "%#v", (f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f WordOrOperator) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = WordOrOperator
		type WordOrOperator X

		fmt.Fprintf(s, "%#v", (f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f WordOrToken) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = WordOrToken
		type WordOrToken X

		fmt.Fprintf(s, "%#v", (f))
	} else {
		format(&f, s, v)
	}
}

// Format implements the fmt.Formatter interface
func (f WordPart) Format(s fmt.State, v rune) {
	if v == 'v' && s.Flag('#') {
		type X = WordPart
		type WordPart X

		fmt.Fprintf(s, "%#v", (f))
	} else {
		format(&f, s, v)
	}
}
