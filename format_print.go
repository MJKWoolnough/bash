package bash

import "io"

func (a ArithmeticExpansion) printSource(w io.Writer, v bool) {}

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
