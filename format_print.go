package bash

import "io"

func (a ArithmeticExpansion) printSource(w io.Writer, v bool) {
	io.WriteString(w, "$((")

	if len(a.WordsAndOperators) > 0 {
		if v {
			io.WriteString(w, " ")
		}

		a.WordsAndOperators[0].printSource(w, v)

		for _, wo := range a.WordsAndOperators[1:] {
			if v {
				io.WriteString(w, " ")
			}

			wo.printSource(w, v)
		}

		if v {
			io.WriteString(w, " ")
		}
	}

	io.WriteString(w, "))")
}

func (a Assignment) printSource(w io.Writer, v bool) {}

func (c Command) printSource(w io.Writer, v bool) {}

func (c CommandSubstitution) printSource(w io.Writer, v bool) {}

func (f File) printSource(w io.Writer, v bool) {}

func (p ParameterAssign) printSource(w io.Writer, v bool) {}

func (p ParameterExpansion) printSource(w io.Writer, v bool) {}

func (p Parameter) printSource(w io.Writer, v bool) {}

func (p Pipeline) printSource(w io.Writer, v bool) {}

func (r Redirection) printSource(w io.Writer, v bool) {}

func (s Statement) printSource(w io.Writer, v bool) {}

func (s String) printSource(w io.Writer, v bool) {}

func (ve Value) printSource(w io.Writer, v bool) {}

func (wo WordOrOperator) printSource(w io.Writer, v bool) {
	if wo.Operator != nil {
		io.WriteString(w, wo.Operator.Data)
	} else if wo.Word != nil {
		wo.Word.printSource(w, v)
	}
}

func (wp WordPart) printSource(w io.Writer, v bool) {
	if wp.Part != nil {
		io.WriteString(w, wp.Part.Data)
	} else if wp.ArithmeticExpansion != nil {
		wp.ArithmeticExpansion.printSource(w, v)
	} else if wp.CommandSubstitution != nil {
		wp.CommandSubstitution.printSource(w, v)
	} else if wp.ParameterExpansion != nil {
		wp.ParameterExpansion.printSource(w, v)
	}
}

func (wd Word) printSource(w io.Writer, v bool) {
	for _, word := range wd.Parts {
		word.printSource(w, v)
	}
}
