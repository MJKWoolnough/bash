package bash

import "io"

func (a ArithmeticExpansion) printSource(w io.Writer, v bool) {
	if a.Expression {
		io.WriteString(w, "((")
	} else {
		io.WriteString(w, "$((")
	}

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

func (a ArrayWord) printSource(w io.Writer, v bool) {
	if len(a.Comments[0]) > 0 {
		io.WriteString(w, " ")
		a.Comments[0].printSource(w, v)
	}

	a.Word.printSource(w, v)

	if len(a.Comments[1]) > 0 {
		io.WriteString(w, " ")
		a.Comments[1].printSource(w, v)
	}
}

func (a Assignment) printSource(w io.Writer, v bool) {
	if a.Assignment != AssignmentAssign && a.Assignment != AssignmentAppend {
		return
	}

	a.Identifier.printSource(w, v)
	a.Assignment.printSource(w, v)
	a.Value.printSource(w, v)
}

func (a AssignmentOrWord) printSource(w io.Writer, v bool) {
	if a.Assignment != nil {
		a.Assignment.printSource(w, v)
	} else if a.Word != nil {
		a.Word.printSource(w, v)
	}
}

func (c CaseCompound) printSource(w io.Writer, v bool) {
	io.WriteString(w, "case ")
	c.Word.printSource(w, v)

	if len(c.Comments[0]) > 0 {
		io.WriteString(w, " ")
		c.Comments[0].printSource(w, false)
		io.WriteString(w, "\nin ")
	} else {
		io.WriteString(w, " in ")
	}

	c.Comments[1].printSource(w, true)

	for _, m := range c.Matches {
		io.WriteString(w, "\n")
		m.printSource(w, v)
	}

	if len(c.Comments[2]) > 0 {
		io.WriteString(w, "\n")
		c.Comments[2].printSource(w, false)
	}

	io.WriteString(w, "\nesac")
}

func (c Command) printSource(w io.Writer, v bool) {
	if len(c.Vars) > 0 {
		c.Vars[0].printSource(w, v)

		for _, vr := range c.Vars[1:] {
			io.WriteString(w, " ")
			vr.printSource(w, v)
		}
	}

	if len(c.AssignmentsOrWords) > 0 {
		if len(c.Vars) > 0 {
			io.WriteString(w, " ")
		}

		c.AssignmentsOrWords[0].printSource(w, v)

		for _, wd := range c.AssignmentsOrWords[1:] {
			io.WriteString(w, " ")
			wd.printSource(w, v)
		}
	}

	if len(c.Redirections) > 0 {
		if len(c.Vars) > 0 || len(c.AssignmentsOrWords) > 0 {
			io.WriteString(w, " ")
		}

		c.Redirections[0].printSource(w, v)

		for _, r := range c.Redirections[1:] {
			r.printSource(w, v)
		}
	}
}

func (c Command) printHeredoc(w io.Writer, v bool) {
	for _, r := range c.Redirections {
		r.printHeredoc(w, v)
	}
}

func (c CommandOrCompound) printSource(w io.Writer, v bool) {
	if c.Command != nil {
		c.Command.printSource(w, v)
	} else if c.Compound != nil {
		c.Compound.printSource(w, v)
	}
}

func (c CommandOrCompound) printHeredoc(w io.Writer, v bool) {
	if c.Command != nil {
		c.Command.printHeredoc(w, v)
	} else if c.Compound != nil {
		c.Compound.printHeredoc(w, v)
	}
}

func (c CommandSubstitution) printSource(w io.Writer, v bool) {
	io.WriteString(w, "$(")
	c.Command.printSource(w, v)
	io.WriteString(w, ")")
}

func (c Compound) printSource(w io.Writer, v bool) {
	if c.IfCompound != nil {
		c.IfCompound.printSource(w, v)
	} else if c.CaseCompound != nil {
		c.CaseCompound.printSource(w, v)
	} else if c.LoopCompound != nil {
		c.LoopCompound.printSource(w, v)
	} else if c.ForCompound != nil {
		c.ForCompound.printSource(w, v)
	} else if c.SelectCompound != nil {
		c.SelectCompound.printSource(w, v)
	} else if c.GroupingCompound != nil {
		c.GroupingCompound.printSource(w, v)
	} else if c.TestCompound != nil {
		c.TestCompound.printSource(w, v)
	} else if c.ArithmeticCompound != nil {
		c.ArithmeticCompound.printSource(w, v)
	} else if c.FunctionCompound != nil {
		c.FunctionCompound.printSource(w, v)
	}

	for _, r := range c.Redirections {
		io.WriteString(w, " ")
		r.printSource(w, v)
	}
}

func (c Compound) printHeredoc(w io.Writer, v bool) {
	for _, r := range c.Redirections {
		r.printHeredoc(w, v)
	}
}

func (f File) printSource(w io.Writer, v bool) {
	f.Comments[0].printSource(w, true)

	if len(f.Lines) > 0 {
		f.Lines[0].printSource(w, v)

		for _, l := range f.Lines[1:] {
			io.WriteString(w, "\n")
			l.printSource(w, v)
		}
	}

	f.Comments[1].printSource(w, true)
}

func (f ForCompound) printSource(w io.Writer, v bool) {
	if f.ArithmeticExpansion == nil && f.Identifier == nil {
		return
	}

	io.WriteString(w, "for ")

	if f.ArithmeticExpansion != nil {
		f.ArithmeticExpansion.printSource(w, v)
	} else {
		io.WriteString(w, f.Identifier.Data)

		if f.Words != nil {
			io.WriteString(w, " ")
			f.Comments[0].printSource(w, true)
			io.WriteString(w, "in")

			for _, wd := range f.Words {
				io.WriteString(w, " ")
				wd.printSource(w, v)
			}
		}
	}

	ip := indentPrinter{Writer: w}

	io.WriteString(&ip, "; ")
	f.Comments[1].printSource(&ip, true)
	io.WriteString(&ip, "do\n")
	f.File.printSource(&ip, v)
	io.WriteString(&ip, "\ndone")
}

func (f FunctionCompound) printSource(w io.Writer, v bool) {
	if f.Identifier == nil {
		return
	}

	if f.HasKeyword {
		io.WriteString(w, "function ")
	}

	io.WriteString(w, f.Identifier.Data)
	io.WriteString(w, "() ")
	f.Comments.printSource(w, true)
	f.Body.printSource(w, v)
}

func (g GroupingCompound) printSource(w io.Writer, v bool) {
	ip := indentPrinter{Writer: w}

	if g.SubShell {
		io.WriteString(&ip, "(\n")
	} else {
		io.WriteString(&ip, "{\n")
	}

	g.File.printSource(&ip, v)

	if g.SubShell {
		io.WriteString(&ip, "\n)")
	} else {
		io.WriteString(&ip, "\n}")
	}
}

func (h Heredoc) printSource(w io.Writer, v bool) {
	io.WriteString(w, "\n")

	for _, p := range h.HeredocPartsOrWords {
		p.printSource(w, v)
	}
}

func (h HeredocPartOrWord) printSource(w io.Writer, v bool) {
	if h.HeredocPart != nil {
		io.WriteString(w, h.HeredocPart.Data)
	} else if h.Word != nil {
		h.Word.printSource(w, v)
	}
}

func (i IfCompound) printSource(w io.Writer, v bool) {
	io.WriteString(w, "if ")
	i.If.printSource(w, v)

	for _, e := range i.ElIf {
		io.WriteString(w, "\nelif ")
		e.printSource(w, v)
	}

	if i.Else != nil {
		ip := indentPrinter{Writer: w}

		io.WriteString(w, "\nelse")
		io.WriteString(&ip, "\n")
		i.Else.printSource(&ip, v)
	}

	io.WriteString(w, "fi")
}

func (l Line) printSource(w io.Writer, v bool) {
	if len(l.Statements) > 0 {
		l.Comments[0].printSource(w, true)
		l.Statements[0].printSource(w, v)

		for _, s := range l.Statements[1:] {
			io.WriteString(w, " ")
			s.printSource(w, v)
		}

		l.Comments[1].printSource(w, false)

		for _, s := range l.Statements {
			s.printHeredoc(w, v)
		}
	}
}

func (l LoopCompound) printSource(w io.Writer, v bool) {
	if l.Until {
		io.WriteString(w, "until ")
	} else {
		io.WriteString(w, "while ")
	}

	l.Statement.printSource(w, v)

	ip := indentPrinter{Writer: w}

	l.Comments.printSource(&ip, true)
	io.WriteString(&ip, " do\n")
	l.File.printSource(&ip, v)
	io.WriteString(w, "\ndone")
}

func (p ParameterAssign) printSource(w io.Writer, v bool) {
	if p.Identifier != nil {
		io.WriteString(w, p.Identifier.Data)

		if len(p.Subscript) > 0 {
			io.WriteString(w, "[")

			for _, s := range p.Subscript {
				s.printSource(w, v)
			}

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

			for _, a := range p.Array {
				a.printSource(w, v)
			}

			io.WriteString(w, "]")
		}
	}
}

func (p Pattern) printSource(w io.Writer, v bool) {
	for _, word := range p.Parts {
		word.printSource(w, v)
	}
}

func (p PatternLines) printSource(w io.Writer, v bool) {
	if len(p.Patterns) == 0 {
		return
	}

	p.Comments.printSource(w, true)

	p.Patterns[0].printSource(w, v)

	for _, pattern := range p.Patterns[1:] {
		io.WriteString(w, "|")
		pattern.printSource(w, v)
	}

	ip := indentPrinter{Writer: w}

	io.WriteString(w, ")")
	p.Lines.printSource(&ip, v)
	io.WriteString(w, "\n")
	p.CaseTerminationType.printSource(w, v)
}

func (p Pipeline) printSource(w io.Writer, v bool) {
	p.PipelineTime.printSource(w, v)

	if p.Not {
		io.WriteString(w, "! ")
	}

	if p.Coproc {
		io.WriteString(w, "coproc ")

		if p.CoprocIdentifier != nil {
			io.WriteString(w, p.CoprocIdentifier.Data)
			io.WriteString(w, " ")
		}
	}

	p.CommandOrCompound.printSource(w, v)

	if p.Pipeline != nil {
		io.WriteString(w, " | ")
		p.Pipeline.printSource(w, v)
	}
}

func (p Pipeline) printHeredoc(w io.Writer, v bool) {
	p.CommandOrCompound.printHeredoc(w, v)

	if p.Pipeline != nil {
		p.printHeredoc(w, v)
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

func (r Redirection) printHeredoc(w io.Writer, v bool) {
	if r.Redirector != nil && r.Heredoc != nil && (r.Redirector.Data == "<<" || r.Redirector.Data == "<<-") {
		if r.Redirector.Data == "<<" {
			w = unwrapIndentPrinter(w)
		}

		r.Heredoc.printSource(w, v)
		r.Output.printSource(w, v)
	}
}

func (s SelectCompound) printSource(w io.Writer, v bool) {
	if s.Identifier == nil {
		return
	}

	io.WriteString(w, "select ")
	io.WriteString(w, s.Identifier.Data)

	if s.Words != nil {
		io.WriteString(w, " ")
		s.Comments[0].printSource(w, true)
		io.WriteString(w, "in")

		for _, wd := range s.Words {
			io.WriteString(w, " ")
			wd.printSource(w, v)
		}
	}

	ip := indentPrinter{Writer: w}

	io.WriteString(&ip, "; ")
	s.Comments[1].printSource(w, true)
	io.WriteString(&ip, "do\n")
	s.File.printSource(&ip, v)
	io.WriteString(&ip, "\ndone")
}

func (s Statement) printSource(w io.Writer, v bool) {
	s.Pipeline.printSource(w, v)

	if (s.LogicalOperator == LogicalOperatorAnd || s.LogicalOperator == LogicalOperatorOr) && s.Statement != nil {
		s.LogicalOperator.printSource(w, v)
		s.Statement.printSource(w, v)
	}

	if s.JobControl == JobControlBackground {
		if v {
			io.WriteString(w, " &")
		} else {
			io.WriteString(w, "&")
		}
	} else {
		io.WriteString(w, ";")
	}
}

func (s Statement) printHeredoc(w io.Writer, v bool) {
	s.Pipeline.printHeredoc(w, v)

	if s.Statement != nil {
		s.Statement.printHeredoc(w, v)
	}
}

func (s String) printSource(w io.Writer, v bool) {
	for _, p := range s.WordsOrTokens {
		p.printSource(w, v)
	}
}

func (t TestCompound) printSource(w io.Writer, v bool) {
	io.WriteString(w, "[[ ")
	t.Comments[0].printSource(w, true)
	t.Tests.printSource(w, v)
	io.WriteString(w, " ")
	t.Comments[1].printSource(w, true)
	io.WriteString(w, "]]")
}

func (t Tests) printSource(w io.Writer, v bool) {
	t.Comments[0].printSource(w, true)

	if t.Not {
		io.WriteString(w, "! ")

		t.Comments[1].printSource(w, true)
	}

	if t.Parens != nil {
		io.WriteString(w, "( ")
		t.Comments[2].printSource(w, true)
		t.Parens.printSource(w, v)
		io.WriteString(w, " ")
		t.Comments[3].printSource(w, true)
		io.WriteString(w, ")")
	} else if t.Word != nil && t.Test == TestOperatorNone {
		t.Word.printSource(w, v)
	} else if t.Word != nil && t.Pattern != nil && t.Test >= TestOperatorStringsEqual {
		t.Word.printSource(w, v)
		io.WriteString(w, " ")
		t.Test.printSource(w, v)
		io.WriteString(w, " ")
		t.Comments[2].printSource(w, true)
		t.Pattern.printSource(w, v)
	} else if t.Word != nil && t.Test >= TestOperatorFileExists && t.Test <= TestOperatorVarnameIsRef {
		t.Test.printSource(w, v)
		io.WriteString(w, " ")
		t.Comments[2].printSource(w, true)
		t.Word.printSource(w, v)
	} else {
		return
	}

	if t.Tests != nil && (t.LogicalOperator == LogicalOperatorOr || t.LogicalOperator == LogicalOperatorAnd) {
		io.WriteString(w, " ")
		t.Comments[4].printSource(w, true)
		t.LogicalOperator.printSource(w, v)
		t.Tests.printSource(w, v)
	} else if len(t.Comments[4]) > 0 {
		io.WriteString(w, " ")
		t.Comments[4].printSource(w, true)
	}
}

func (t TestConsequence) printSource(w io.Writer, v bool) {
	t.Test.printSource(w, v)

	ip := indentPrinter{Writer: w}

	if len(t.Comments) > 0 {
		io.WriteString(w, " ")
		t.Comments.printSource(&ip, true)
		io.WriteString(&ip, "then\n")
	} else {
		io.WriteString(&ip, " then\n")
	}

	t.Consequence.printSource(&ip, v)
}

func (ve Value) printSource(w io.Writer, v bool) {
	if ve.Word != nil {
		ve.printSource(w, v)
	} else if ve.Array != nil {
		if v || len(ve.Comments[0]) > 0 {
			io.WriteString(w, "( ")
			ve.Comments[0].printSource(w, true)
		} else {
			io.WriteString(w, "(")
		}

		if len(ve.Array) > 0 {
			ve.Array[0].printSource(w, v)

			for _, word := range ve.Array[1:] {
				io.WriteString(w, " ")
				word.printSource(w, v)
			}
		}

		ve.Comments[1].printSource(w, true)

		if v && len(ve.Comments) == 0 {
			io.WriteString(w, " )")
		} else {
			io.WriteString(w, ")")
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
