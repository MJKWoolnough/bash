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

func (a Assignment) printSource(w io.Writer, v bool) {
	if a.Assignment != AssignmentAssign && a.Assignment != AssignmentAppend {
		return
	}

	a.Identifier.printSource(w, v)
	a.Assignment.printSource(w, v)
	a.Value.printSource(w, v)
}

func (c Command) printSource(w io.Writer, v bool) {}

func (c CommandSubstitution) printSource(w io.Writer, v bool) {
	io.WriteString(w, "$(")
	c.Command.printSource(w, v)
	io.WriteString(w, ")")
}

func (f File) printSource(w io.Writer, v bool) {}

func (p ParameterAssign) printSource(w io.Writer, v bool) {
	if p.Identifier != nil {
		io.WriteString(w, p.Identifier.Data)

		if p.Subscript != nil {
			io.WriteString(w, "[")
			p.Subscript.printSource(w, v)
			io.WriteString(w, "]")
		}
	}
}

func (p ParameterExpansion) printSource(w io.Writer, v bool) {}

func (p Parameter) printSource(w io.Writer, v bool) {
	if p.Parameter != nil {
		io.WriteString(w, p.Parameter.Data)

		if p.Array != nil {
			io.WriteString(w, "[")
			p.Array.printSource(w, v)
			io.WriteString(w, "]")
		}
	}
}

func (p Pipeline) printSource(w io.Writer, v bool) {}

func (r Redirection) printSource(w io.Writer, v bool) {
	if r.Redirector == nil {
		return
	}

	if r.Input != nil {
		io.WriteString(w, r.Input.Data)
	}

	io.WriteString(w, r.Redirector.Data)
	r.Output.printSource(w, v)
}

func (s Statement) printSource(w io.Writer, v bool) {}

func (s String) printSource(w io.Writer, v bool) {
	for _, p := range s.WordsOrTokens {
		p.printSource(w, v)
	}
}

func (ve Value) printSource(w io.Writer, v bool) {
	if ve.Word != nil {
		ve.printSource(w, v)
	} else if ve.Array != nil {
		if len(ve.Array) == 0 {
			io.WriteString(w, "()")
		} else {

			if v {
				io.WriteString(w, "( ")
			} else {
				io.WriteString(w, "(")
			}

			ve.Array[0].printSource(w, v)

			for _, word := range ve.Array[1:] {
				io.WriteString(w, " ")
				word.printSource(w, v)
			}

			if v {
				io.WriteString(w, " )")
			} else {
				io.WriteString(w, ")")
			}
		}
	}
}

func (wo WordOrOperator) printSource(w io.Writer, v bool) {
	if wo.Operator != nil {
		io.WriteString(w, wo.Operator.Data)
	} else if wo.Word != nil {
		wo.Word.printSource(w, v)
	}
}

func (wt WordOrToken) printSource(w io.Writer, v bool) {
	if wt.Word != nil {
		wt.Word.printSource(w, v)
	} else if wt.Token != nil {
		io.WriteString(w, wt.Token.Data)
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
