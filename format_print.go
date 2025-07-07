package bash

import (
	"strings"

	"vimagination.zapto.org/parser"
)

func (a ArithmeticExpansion) printSource(w writer, v bool) {
	if a.Expression {
		w.WriteString("((")
	} else {
		w.WriteString("$((")
	}

	if len(a.WordsAndOperators) > 0 {
		if v {
			w.WriteString(" ")
		}

		a.WordsAndOperators[0].printSource(w, v)

		for _, wo := range a.WordsAndOperators[1:] {
			if v && !wo.operatorIsToken(parser.Token{Type: TokenPunctuator, Data: ";"}) {
				w.WriteString(" ")
			}

			wo.printSource(w, v)
		}

		if v {
			w.WriteString(" ")
		}
	}

	w.WriteString("))")
}

func (a ArrayWord) printSource(w writer, v bool) {
	if len(a.Comments[0]) > 0 {
		a.Comments[0].printSource(w, true)
	}

	a.Word.printSource(w, v)

	if len(a.Comments[1]) > 0 {
		w.WriteString(" ")
		a.Comments[1].printSource(w, false)
	}
}

func (a Assignment) printSource(w writer, v bool) {
	if a.Assignment == AssignmentAssign || a.Assignment == AssignmentAppend {
		a.Identifier.printSource(w, v)
		a.Assignment.printSource(w, v)

		if a.Value != nil {
			a.Value.printSource(w, v)
		} else {
			parens := 0

			for _, e := range a.Expression {
				if parens > 0 {
					w.WriteString(" ")
				}

				e.printSource(w, v)

				if v && e.Operator != nil {
					if e.Operator.Token == (parser.Token{Type: TokenPunctuator, Data: "("}) {
						parens++
					} else if e.Operator.Token == (parser.Token{Type: TokenPunctuator, Data: ")"}) {
						parens--
					}
				}
			}
		}
	}
}

func (a AssignmentOrWord) printSource(w writer, v bool) {
	if a.Assignment != nil {
		a.Assignment.printSource(w, v)
	} else if a.Word != nil {
		a.Word.printSource(w, v)
	}
}

func (c CaseCompound) printSource(w writer, v bool) {
	w.WriteString("case ")
	c.Word.printSource(w, v)

	if len(c.Comments[0]) > 0 {
		w.WriteString(" ")
		c.Comments[0].printSource(w, false)
		w.WriteString("\nin")
	} else {
		w.WriteString(" in")
	}

	if len(c.Comments[1]) > 0 {
		w.WriteString(" ")
		c.Comments[1].printSource(w, false)
	}

	for _, m := range c.Matches {
		w.WriteString("\n")
		m.printSource(w, v)
	}

	if len(c.Comments[2]) > 0 {
		w.WriteString("\n")
		c.Comments[2].printSource(w, false)
	}

	w.WriteString("\nesac")
}

func (c Command) printSource(w writer, v bool) {
	if len(c.Vars) > 0 {
		c.Vars[0].printSource(w, v)

		for _, vr := range c.Vars[1:] {
			w.WriteString(" ")
			vr.printSource(w, v)
		}
	}

	if len(c.AssignmentsOrWords) > 0 {
		if len(c.Vars) > 0 {
			w.WriteString(" ")
		}

		c.AssignmentsOrWords[0].printSource(w, v)

		for _, wd := range c.AssignmentsOrWords[1:] {
			w.WriteString(" ")
			wd.printSource(w, v)
		}
	}

	if len(c.Redirections) > 0 {
		if len(c.Vars) > 0 || len(c.AssignmentsOrWords) > 0 {
			w.WriteString(" ")
		}

		c.Redirections[0].printSource(w, v)

		for _, r := range c.Redirections[1:] {
			w.WriteString(" ")
			r.printSource(w, v)
		}
	}
}

func (c Command) printHeredoc(w writer, v bool) {
	for _, r := range c.Redirections {
		r.printHeredoc(w, v)
	}
}

func (c CommandOrCompound) printSource(w writer, v bool) {
	if c.Command != nil {
		c.Command.printSource(w, v)
	} else if c.Compound != nil {
		c.Compound.printSource(w, v)
	}
}

func (c CommandOrCompound) printHeredoc(w writer, v bool) {
	if c.Command != nil {
		c.Command.printHeredoc(w, v)
	} else if c.Compound != nil {
		c.Compound.printHeredoc(w, v)
	}
}

func (c CommandSubstitution) printSource(w writer, v bool) {
	w.WriteString("$(")

	if c.Command.isMultiline(v) {
		ip := indentPrinter{writer: w}

		ip.WriteString("\n")
		c.Command.printSource(&ip, v)
		w.WriteString("\n")
	} else {
		c.Command.printSourceEnd(w, v, false)
	}

	w.WriteString(")")
}

func (c Compound) printSource(w writer, v bool) {
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
		w.WriteString(" ")
		r.printSource(w, v)
	}
}

func (c Compound) printHeredoc(w writer, v bool) {
	for _, r := range c.Redirections {
		r.printHeredoc(w, v)
	}
}

func (f File) printSource(w writer, v bool) {
	f.printSourceEnd(w, v, true)
}

func (f File) printSourceEnd(w writer, v, end bool) {
	f.Comments[0].printSource(w, true)

	if len(f.Lines) > 0 {
		f.Lines[0].printSourceEnd(w, v, end || len(f.Lines) > 1)

		lastLine := lastTokenPos(f.Lines[0].Tokens)

		for n, l := range f.Lines[1:] {
			if firstTokenPos(l.Tokens) > lastLine+1 {
				w.WriteString("\n")
			}

			w.WriteString("\n")
			l.printSourceEnd(w, v, end || len(f.Lines) > n+1)

			lastLine = lastTokenPos(l.Tokens)
		}
	}

	if len(f.Comments[1]) > 0 {
		w.WriteString("\n\n")
		f.Comments[1].printSource(w, false)
	}
}

func firstTokenPos(tk Tokens) (pos uint64) {
	if len(tk) > 0 {
		pos = tk[0].Line
	}

	return pos
}

func lastTokenPos(tk Tokens) (pos uint64) {
	if len(tk) > 0 {
		pos = tk[len(tk)-1].Line
	}

	return pos
}

func (f ForCompound) printSource(w writer, v bool) {
	if f.ArithmeticExpansion != nil || f.Identifier != nil {
		w.WriteString("for ")

		if f.ArithmeticExpansion != nil {
			f.ArithmeticExpansion.printSource(w, v)
		} else {
			w.WriteString(f.Identifier.Data)

			if f.Words != nil {
				w.WriteString(" ")
				f.Comments[0].printSource(w, true)
				w.WriteString("in")

				for _, wd := range f.Words {
					w.WriteString(" ")
					wd.printSource(w, v)
				}
			}
		}

		ip := indentPrinter{writer: w}

		ip.WriteString("; ")
		f.Comments[1].printSource(&ip, true)
		ip.WriteString("do\n")
		f.File.printSource(&ip, v)
		w.WriteString("\ndone")
	}
}

func (f FunctionCompound) printSource(w writer, v bool) {
	if f.Identifier != nil {
		if f.HasKeyword {
			w.WriteString("function ")
		}

		w.WriteString(f.Identifier.Data)
		w.WriteString("() ")
		f.Comments.printSource(w, true)
		f.Body.printSource(w, v)
	}
}

func (g GroupingCompound) printSource(w writer, v bool) {
	if g.SubShell {
		w.WriteString("(")
	} else {
		w.WriteString("{")
	}

	ip := indentPrinter{writer: w}

	multiline := v || g.File.isMultiline(v)

	if len(g.File.Comments[0]) > 0 || !multiline {
		w.WriteString(" ")
	} else {
		ip.WriteString("\n")
	}

	g.File.printSource(&ip, v)

	if multiline {
		w.WriteString("\n")
	} else {
		w.WriteString(" ")
	}

	if g.SubShell {
		w.WriteString(")")
	} else {
		w.WriteString("}")
	}
}

func (h Heredoc) printSource(w writer, v bool) {
	w.WriteString("\n")

	for _, p := range h.HeredocPartsOrWords {
		p.printSource(w, v)
	}
}

func (h HeredocPartOrWord) printSource(w writer, v bool) {
	if h.HeredocPart != nil {
		w.WriteString(h.HeredocPart.Data)
	} else if h.Word != nil {
		h.Word.printSource(w, v)
	}
}

func (i IfCompound) printSource(w writer, v bool) {
	w.WriteString("if ")
	i.If.printSource(w, v)

	for _, e := range i.ElIf {
		w.WriteString("\nelif ")
		e.printSource(w, v)
	}

	if i.Else != nil {
		ip := indentPrinter{writer: w}

		w.WriteString("\nelse")
		ip.WriteString("\n")
		i.Else.printSource(&ip, v)
	}

	w.WriteString("\nfi")
}

func (l Line) printSource(w writer, v bool) {
	l.printSourceEnd(w, v, true)
}

func (l Line) printSourceEnd(w writer, v, end bool) {
	if len(l.Statements) > 0 {
		l.Comments[0].printSource(w, true)
		l.Statements[0].printSourceEnd(w, v, end || len(l.Statements) > 1)

		eos := " "

		if v {
			eos = "\n"
		}

		for n, s := range l.Statements[1:] {
			w.WriteString(eos)
			s.printSourceEnd(w, v, end || len(l.Statements) > n+1)
		}

		if len(l.Comments[1]) > 0 {
			w.WriteString(" ")
			l.Comments[1].printSource(w, false)
		}

		for _, s := range l.Statements {
			s.printHeredoc(w, v)
		}
	}
}

func (l LoopCompound) printSource(w writer, v bool) {
	if l.Until {
		w.WriteString("until ")
	} else {
		w.WriteString("while ")
	}

	l.Statement.printSource(w, v)

	if len(l.Comments) > 0 {
		w.WriteString(" ")
		l.Comments.printSource(w, true)
		w.WriteString("do")
	} else {
		w.WriteString(" do")
	}

	ip := indentPrinter{writer: w}

	ip.WriteString("\n")
	l.File.printSource(&ip, v)
	w.WriteString("\ndone")
}

func (p ParameterAssign) printSource(w writer, v bool) {
	if p.Identifier != nil {
		w.WriteString(p.Identifier.Data)

		if len(p.Subscript) > 0 {
			w.WriteString("[")
			p.Subscript[0].printSource(w, v)

			for _, s := range p.Subscript[1:] {
				if v {
					w.WriteString(" ")
				}

				s.printSource(w, v)
			}

			w.WriteString("]")
		}
	}
}

func (p ParameterExpansion) printSource(w writer, v bool) {
	w.WriteString("${")

	if p.Indirect || p.Type == ParameterPrefix || p.Type == ParameterPrefixSeperate {
		w.WriteString("!")
	} else if p.Type == ParameterLength {
		w.WriteString("#")
	}

	p.Parameter.printSource(w, v)

	if p.Word != nil {
		switch p.Type {
		case ParameterSubstitution:
			w.WriteString(":=")
			p.Word.printSource(w, v)
		case ParameterAssignment:
			w.WriteString(":?")
			p.Word.printSource(w, v)
		case ParameterMessage:
			w.WriteString(":+")
			p.Word.printSource(w, v)
		case ParameterSetAssign:
			w.WriteString(":-")
			p.Word.printSource(w, v)
		case ParameterUnsetSubstitution:
			w.WriteString("=")
			p.Word.printSource(w, v)
		case ParameterUnsetAssignment:
			w.WriteString("?")
			p.Word.printSource(w, v)
		case ParameterUnsetMessage:
			w.WriteString("+")
			p.Word.printSource(w, v)
		case ParameterUnsetSetAssign:
			w.WriteString("-")
			p.Word.printSource(w, v)
		case ParameterRemoveStartShortest:
			w.WriteString("#")
			p.Word.printSource(w, v)
		case ParameterRemoveStartLongest:
			w.WriteString("##")
			p.Word.printSource(w, v)
		case ParameterRemoveEndShortest:
			w.WriteString("%")
			p.Word.printSource(w, v)
		case ParameterRemoveEndLongest:
			w.WriteString("%%")
			p.Word.printSource(w, v)
		}
	} else if p.Pattern != nil {
		isReplacement := false

		switch p.Type {
		case ParameterReplace:
			isReplacement = true

			w.WriteString("/")
			w.WriteString(p.Pattern.Data)
		case ParameterReplaceAll:
			isReplacement = true

			w.WriteString("//")
			w.WriteString(p.Pattern.Data)
		case ParameterReplaceStart:
			isReplacement = true

			w.WriteString("/#")
			w.WriteString(p.Pattern.Data)
		case ParameterReplaceEnd:
			isReplacement = true

			w.WriteString("/%")
			w.WriteString(p.Pattern.Data)
		case ParameterLowercaseFirstMatch:
			w.WriteString(",")
			w.WriteString(p.Pattern.Data)
		case ParameterLowercaseAllMatches:
			w.WriteString(",,")
			w.WriteString(p.Pattern.Data)
		case ParameterUppercaseFirstMatch:
			w.WriteString("^")
			w.WriteString(p.Pattern.Data)
		case ParameterUppercaseAllMatches:
			w.WriteString("^^")
			w.WriteString(p.Pattern.Data)
		}

		if isReplacement && p.String != nil {
			w.WriteString("/")
			p.String.printSource(w, v)
		}
	} else if !p.Indirect {
		switch p.Type {
		case ParameterPrefix:
			w.WriteString("*")
		case ParameterPrefixSeperate:
			w.WriteString("@")
		}
	}

	switch p.Type {
	case ParameterSubstring:
		if p.SubstringStart != nil {
			w.WriteString(":")

			if strings.HasPrefix(p.SubstringStart.Data, "-") {
				w.WriteString(" ")
			}

			w.WriteString(p.SubstringStart.Data)

			if p.SubstringEnd != nil {
				w.WriteString(":")
				w.WriteString(p.SubstringEnd.Data)
			}
		}
	case ParameterUppercase:
		w.WriteString("@U")
	case ParameterUppercaseFirst:
		w.WriteString("@u")
	case ParameterLowercase:
		w.WriteString("@L")
	case ParameterQuoted:
		w.WriteString("@Q")
	case ParameterEscaped:
		w.WriteString("@E")
	case ParameterPrompt:
		w.WriteString("@P")
	case ParameterDeclare:
		w.WriteString("@A")
	case ParameterQuotedArrays:
		w.WriteString("@K")
	case ParameterAttributes:
		w.WriteString("@a")
	case ParameterQuotedArraysSeperate:
		w.WriteString("@k")
	}

	w.WriteString("}")
}

func (p Parameter) printSource(w writer, v bool) {
	if p.Parameter != nil {
		w.WriteString(p.Parameter.Data)

		if p.Array != nil {
			w.WriteString("[")

			for _, a := range p.Array {
				a.printSource(w, v)
			}

			w.WriteString("]")
		}
	}
}

func (p Pattern) printSource(w writer, v bool) {
	for _, word := range p.Parts {
		word.printSource(w, v)
	}
}

func (p PatternLines) printSource(w writer, v bool) {
	if len(p.Patterns) > 0 {
		p.Comments.printSource(w, true)

		p.Patterns[0].printSource(w, v)

		for _, pattern := range p.Patterns[1:] {
			w.WriteString("|")
			pattern.printSource(w, v)
		}

		ip := indentPrinter{writer: w}

		w.WriteString(")")
		ip.WriteString("\n")

		if len(p.Lines.Lines) > 0 {
			p.Lines.printSource(&ip, v)
		} else {
			ip.WriteString(";")
		}

		p.CaseTerminationType.printSource(&ip, v)
	}
}

func (p Pipeline) printSource(w writer, v bool) {
	p.PipelineTime.printSource(w, v)

	if p.Not {
		w.WriteString("! ")
	}

	if p.Coproc {
		w.WriteString("coproc ")

		if p.CoprocIdentifier != nil {
			w.WriteString(p.CoprocIdentifier.Data)
			w.WriteString(" ")
		}
	}

	p.CommandOrCompound.printSource(w, v)

	if p.Pipeline != nil {
		w.WriteString(" | ")
		p.Pipeline.printSource(w, v)
	}
}

func (p Pipeline) endsWithGrouping() bool {
	return p.CommandOrCompound.Compound != nil && len(p.CommandOrCompound.Compound.Redirections) == 0 && (p.CommandOrCompound.Compound.GroupingCompound != nil || p.CommandOrCompound.Compound.FunctionCompound != nil && p.CommandOrCompound.Compound.FunctionCompound.Body.GroupingCompound != nil)
}

func (p Pipeline) printHeredoc(w writer, v bool) {
	p.CommandOrCompound.printHeredoc(w, v)

	if p.Pipeline != nil {
		p.Pipeline.printHeredoc(w, v)
	}
}

func (r Redirection) printSource(w writer, v bool) {
	if r.Redirector != nil {
		if r.Input != nil {
			w.WriteString(r.Input.Data)
		}

		w.WriteString(r.Redirector.Data)
		r.Output.printSource(w, v)
	}
}

func (r Redirection) printHeredoc(w writer, v bool) {
	if r.Redirector != nil && r.Heredoc != nil && (r.Redirector.Data == "<<" || r.Redirector.Data == "<<-") {
		if r.Redirector.Data == "<<" {
			w = w.Underlying()
		}

		r.Heredoc.printSource(w, v)
		r.Output.printSource(w, v)
	}
}

func (s SelectCompound) printSource(w writer, v bool) {
	if s.Identifier != nil {
		w.WriteString("select ")
		w.WriteString(s.Identifier.Data)

		if s.Words != nil {
			w.WriteString(" ")
			s.Comments[0].printSource(w, true)
			w.WriteString("in")

			for _, wd := range s.Words {
				w.WriteString(" ")
				wd.printSource(w, v)
			}
		}

		ip := indentPrinter{writer: w}

		ip.WriteString("; ")
		s.Comments[1].printSource(w, true)
		ip.WriteString("do\n")
		s.File.printSource(&ip, v)
		w.WriteString("\ndone")
	}
}

func (s Statement) printSource(w writer, v bool) {
	s.printSourceEnd(w, v, true)
}

func (s Statement) printSourceEnd(w writer, v, end bool) {
	s.Pipeline.printSource(w, v)

	if (s.LogicalOperator == LogicalOperatorAnd || s.LogicalOperator == LogicalOperatorOr) && s.Statement != nil {
		w.WriteString(" ")
		s.LogicalOperator.printSource(w, v)
		w.WriteString(" ")
		s.Statement.printSourceEnd(w, v, false)
	}

	if s.JobControl == JobControlBackground {
		if v {
			w.WriteString(" &")
		} else {
			w.WriteString("&")
		}
	} else if end && !s.endsWithGrouping() {
		w.WriteString(";")
	}
}

func (s Statement) endsWithGrouping() bool {
	if s.Statement != nil {
		return s.Statement.endsWithGrouping()
	}

	return s.Pipeline.endsWithGrouping()
}

func (s Statement) printHeredoc(w writer, v bool) {
	s.Pipeline.printHeredoc(w, v)

	if s.Statement != nil {
		s.Statement.printHeredoc(w, v)
	}
}

func (s String) printSource(w writer, v bool) {
	for _, p := range s.WordsOrTokens {
		p.printSource(w, v)
	}
}

func (t TestCompound) printSource(w writer, v bool) {
	w.WriteString("[[")

	iw := w
	multi := t.isMultiline(v)

	if multi {
		iw = &indentPrinter{writer: w}
		multi = true

		if len(t.Comments[0]) > 0 {
			w.WriteString(" ")
			t.Comments[0].printSource(w, len(t.Tests.Comments[0]) > 0)
		}

		iw.WriteString("\n")
	} else {
		w.WriteString(" ")
	}

	t.Tests.printSource(iw, v)

	if len(t.Comments[1]) > 0 {
		if t.Tests.lastIsComment() {
			w.WriteString("\n")
		}

		w.WriteString("\n")
		t.Comments[1].printSource(w, true)
	} else if multi {
		w.WriteString("\n")
	} else {
		w.WriteString(" ")
	}

	w.WriteString("]]")
}

func (t Tests) printSource(w writer, v bool) {
	t.Comments[0].printSource(w, true)

	if t.Not {
		w.WriteString("! ")

		t.Comments[1].printSource(w, true)
	}

	if t.Parens != nil {
		w.WriteString("(")

		multi := t.Parens.isMultiline(v) || len(t.Comments[2]) > 0 || len(t.Comments[3]) > 0
		iw := w

		if multi {
			iw = &indentPrinter{writer: w}

			if len(t.Comments[2]) > 0 {
				w.WriteString(" ")
				t.Comments[2].printSource(w, len(t.Parens.Comments[0]) > 0)
			}

			iw.WriteString("\n")
		} else {
			w.WriteString(" ")
		}

		t.Parens.printSource(iw, v)

		if multi {
			if len(t.Comments[3]) > 0 && len(t.Parens.Comments[4]) > 0 {
				w.WriteString("\n\n")
			} else {
				w.WriteString("\n")
			}
		} else {
			w.WriteString(" ")
		}

		t.Comments[3].printSource(w, true)
		w.WriteString(")")
	} else if t.Word != nil && t.Test == TestOperatorNone {
		t.Word.printSource(w, v)
	} else if t.Word != nil && t.Pattern != nil && t.Test >= TestOperatorStringsEqual {
		t.Word.printSource(w, v)
		w.WriteString(" ")
		t.Test.printSource(w, v)
		w.WriteString(" ")
		t.Comments[2].printSource(w, true)
		t.Pattern.printSource(w, v)
	} else if t.Word != nil && t.Test >= TestOperatorFileExists && t.Test <= TestOperatorVarnameIsRef {
		t.Test.printSource(w, v)
		w.WriteString(" ")
		t.Comments[2].printSource(w, true)
		t.Word.printSource(w, v)
	}

	if t.Tests != nil && (t.LogicalOperator == LogicalOperatorOr || t.LogicalOperator == LogicalOperatorAnd) {
		w.WriteString(" ")
		t.Comments[4].printSource(w, true)
		t.LogicalOperator.printSource(w, v)
		w.WriteString(" ")
		t.Tests.printSource(w, v)
	} else if len(t.Comments[4]) > 0 {
		w.WriteString(" ")
		t.Comments[4].printSource(w, false)
	}
}

func (t Tests) lastIsComment() bool {
	if t.Tests != nil {
		return t.Tests.lastIsComment()
	}

	return len(t.Comments[4]) > 0
}

func (t TestConsequence) printSource(w writer, v bool) {
	t.Test.printSource(w, v)

	ip := indentPrinter{writer: w}

	if len(t.Comments) > 0 {
		w.WriteString(" ")
		t.Comments.printSource(&ip, true)
		ip.WriteString("then\n")
	} else {
		ip.WriteString(" then\n")
	}

	t.Consequence.printSource(&ip, v)
}

func (ve Value) printSource(w writer, v bool) {
	if ve.Word != nil {
		ve.Word.printSource(w, v)
	} else if ve.Array != nil {
		iw := w
		ml := ve.isMultiline(v)
		lastHadComment := ml

		if ml {
			iw = &indentPrinter{writer: w}
		}

		if len(ve.Comments[0]) > 0 {
			w.WriteString("(")

			if len(ve.Comments[0]) > 0 {
				w.WriteString(" ")
				ve.Comments[0].printSource(w, len(ve.Array) > 0)

				lastHadComment = true
			}
		} else {
			w.WriteString("(")
		}

		if len(ve.Array) > 0 {
			lastHadComment = lastHadComment || len(ve.Array[0].Comments[0]) > 0

			for _, word := range ve.Array {
				if lastHadComment {
					iw.WriteString("\n")
				} else if v && len(word.Comments[0]) == 0 {
					iw.WriteString(" ")
				}

				word.printSource(iw, v)

				lastHadComment = len(word.Comments[1]) != 0
			}

			if v && !ml {
				iw.WriteString(" ")
			}
		}

		if len(ve.Comments[1]) != 0 {
			w.WriteString("\n")

			if lastHadComment {
				w.WriteString("\n")
			}

			ve.Comments[1].printSource(w, false)

			lastHadComment = true
		}

		if lastHadComment || ml {
			w.WriteString("\n")
		}

		w.WriteString(")")
	}
}

func (wo WordOrOperator) printSource(w writer, v bool) {
	if wo.Operator != nil {
		w.WriteString(wo.Operator.Data)
	} else if wo.Word != nil {
		wo.Word.printSource(w, v)
	}
}

func (wt WordOrToken) printSource(w writer, v bool) {
	if wt.Word != nil {
		wt.Word.printSource(w, v)
	} else if wt.Token != nil {
		w.WriteString(wt.Token.Data)
	}
}

func (wp WordPart) printSource(w writer, v bool) {
	if wp.Part != nil {
		w.WriteString(wp.Part.Data)
	} else if wp.ArithmeticExpansion != nil {
		wp.ArithmeticExpansion.printSource(w, v)
	} else if wp.CommandSubstitution != nil {
		wp.CommandSubstitution.printSource(w, v)
	} else if wp.ParameterExpansion != nil {
		wp.ParameterExpansion.printSource(w, v)
	}
}

func (wd Word) printSource(w writer, v bool) {
	for _, word := range wd.Parts {
		word.printSource(w, v)
	}
}
