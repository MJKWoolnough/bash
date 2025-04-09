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

func (c Command) printSource(w io.Writer, v bool) {
	if len(c.Vars) > 0 {
		c.Vars[0].printSource(w, v)

		for _, vr := range c.Vars[1:] {
			io.WriteString(w, " ")
			vr.printSource(w, v)
		}
	}

	if len(c.Words) > 0 {
		if len(c.Vars) > 0 {
			io.WriteString(w, " ")
		}

		c.Words[0].printSource(w, v)

		for _, wd := range c.Words[1:] {
			io.WriteString(w, " ")
			wd.printSource(w, v)
		}
	}

	if len(c.Redirections) > 0 {
		if len(c.Vars) > 0 || len(c.Words) > 0 {
			io.WriteString(w, " ")
		}

		c.Redirections[0].printSource(w, v)

		for _, r := range c.Redirections[1:] {
			r.printSource(w, v)
		}
	}
}

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

func (p ParameterExpansion) printSource(w io.Writer, v bool) {
	io.WriteString(w, "${")

	if p.Indirect {
		io.WriteString(w, "!")
	} else if p.Type == ParameterLength {
		io.WriteString(w, "#")
	}

	p.Parameter.printSource(w, v)

	if p.Word != nil {
		switch p.Type {
		case ParameterSubstitution:
			io.WriteString(w, ":=")
			p.Word.printSource(w, v)
		case ParameterAssignment:
			io.WriteString(w, ":?")
			p.Word.printSource(w, v)
		case ParameterMessage:
			io.WriteString(w, ":+")
			p.Word.printSource(w, v)
		case ParameterSetAssign:
			io.WriteString(w, ":-")
			p.Word.printSource(w, v)
		case ParameterRemoveStartShortest:
			io.WriteString(w, "#")
			p.Word.printSource(w, v)
		case ParameterRemoveStartLongest:
			io.WriteString(w, "##")
			p.Word.printSource(w, v)
		case ParameterRemoveEndShortest:
			io.WriteString(w, "%")
			p.Word.printSource(w, v)
		case ParameterRemoveEndLongest:
			io.WriteString(w, "%%")
			p.Word.printSource(w, v)
		}
	} else if p.Pattern != nil {
		isReplacement := false

		switch p.Type {
		case ParameterReplace:
			isReplacement = true

			io.WriteString(w, "/")
			io.WriteString(w, p.Pattern.Data)
		case ParameterReplaceAll:
			isReplacement = true

			io.WriteString(w, "//")
			io.WriteString(w, p.Pattern.Data)
		case ParameterReplaceStart:
			isReplacement = true

			io.WriteString(w, "/#")
			io.WriteString(w, p.Pattern.Data)
		case ParameterReplaceEnd:
			isReplacement = true

			io.WriteString(w, "/%")
			io.WriteString(w, p.Pattern.Data)
		case ParameterLowercaseFirstMatch:
			io.WriteString(w, "^")
			io.WriteString(w, p.Pattern.Data)
		case ParameterLowercaseAllMatches:
			io.WriteString(w, "^^")
			io.WriteString(w, p.Pattern.Data)
		case ParameterUppercaseFirstMatch:
			io.WriteString(w, ",")
			io.WriteString(w, p.Pattern.Data)
		case ParameterUppercaseAllMatches:
			io.WriteString(w, ",,")
			io.WriteString(w, p.Pattern.Data)
		}

		if isReplacement && p.String != nil {
			io.WriteString(w, "/")
			p.String.printSource(w, v)
		}
	} else if !p.Indirect {
		switch p.Type {
		case ParameterPrefix:
			io.WriteString(w, "*")
		case ParameterPrefixSeperate:
			io.WriteString(w, "@")
		}
	} else {
		switch p.Type {
		case ParameterSubstring:
			if p.SubstringStart != nil {
				io.WriteString(w, ":")
				io.WriteString(w, p.SubstringStart.Data)

				if p.SubstringEnd != nil {
					io.WriteString(w, ":")
					io.WriteString(w, p.SubstringEnd.Data)
				}
			}
		case ParameterUppercase:
			io.WriteString(w, "@U")
		case ParameterUppercaseFirst:
			io.WriteString(w, "@u")
		case ParameterLowercase:
			io.WriteString(w, "@L")
		case ParameterQuoted:
			io.WriteString(w, "@Q")
		case ParameterEscaped:
			io.WriteString(w, "@E")
		case ParameterPrompt:
			io.WriteString(w, "@P")
		case ParameterDeclare:
			io.WriteString(w, "@A")
		case ParameterQuotedArrays:
			io.WriteString(w, "@K")
		case ParameterQuotedArraysSeperate:
			io.WriteString(w, "@a")
		case ParameterAttributes:
			io.WriteString(w, "@k")
		}
	}

	io.WriteString(w, "}")
}

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

func (p Pipeline) printSource(w io.Writer, v bool) {
	p.PipelineTime.printSource(w, v)

	if p.Not {
		io.WriteString(w, "! ")
	}

	p.Command.printSource(w, v)

	if p.Pipeline != nil {
		io.WriteString(w, " | ")
		p.Pipeline.printSource(w, v)
	}
}

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
