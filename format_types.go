package bash

// File automatically generated with format.sh.

func (f *ArithmeticExpansion) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("ArithmeticExpansion {")

	if f.Expression || v {
		pp.Printf("\nExpression: %v", f.Expression)
	}

	if f.WordsAndOperators == nil {
		pp.WriteString("\nWordsAndOperators: nil")
	} else if len(f.WordsAndOperators) > 0 {
		pp.WriteString("\nWordsAndOperators: [")

		ipp := pp.Indent()

		for n, e := range f.WordsAndOperators {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nWordsAndOperators: []")
	}

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *ArrayWord) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("ArrayWord {")

	pp.WriteString("\nWord: ")
	f.Word.printType(pp, v)

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *Assignment) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("Assignment {")

	pp.WriteString("\nIdentifier: ")
	f.Identifier.printType(pp, v)

	pp.WriteString("\nAssignment: ")
	f.Assignment.printType(pp, v)

	if f.Expression == nil {
		pp.WriteString("\nExpression: nil")
	} else if len(f.Expression) > 0 {
		pp.WriteString("\nExpression: [")

		ipp := pp.Indent()

		for n, e := range f.Expression {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nExpression: []")
	}

	if f.Value != nil {
		pp.WriteString("\nValue: ")
		f.Value.printType(pp, v)
	} else if v {
		pp.WriteString("\nValue: nil")
	}

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *AssignmentOrWord) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("AssignmentOrWord {")

	if f.Assignment != nil {
		pp.WriteString("\nAssignment: ")
		f.Assignment.printType(pp, v)
	} else if v {
		pp.WriteString("\nAssignment: nil")
	}

	if f.Word != nil {
		pp.WriteString("\nWord: ")
		f.Word.printType(pp, v)
	} else if v {
		pp.WriteString("\nWord: nil")
	}

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *BraceExpansion) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("BraceExpansion {")

	pp.WriteString("\nBraceExpansionType: ")
	f.BraceExpansionType.printType(pp, v)

	if f.Words == nil {
		pp.WriteString("\nWords: nil")
	} else if len(f.Words) > 0 {
		pp.WriteString("\nWords: [")

		ipp := pp.Indent()

		for n, e := range f.Words {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nWords: []")
	}

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *BraceWord) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("BraceWord {")

	if f.Parts == nil {
		pp.WriteString("\nParts: nil")
	} else if len(f.Parts) > 0 {
		pp.WriteString("\nParts: [")

		ipp := pp.Indent()

		for n, e := range f.Parts {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nParts: []")
	}

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *CaseCompound) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("CaseCompound {")

	pp.WriteString("\nWord: ")
	f.Word.printType(pp, v)

	if f.Matches == nil {
		pp.WriteString("\nMatches: nil")
	} else if len(f.Matches) > 0 {
		pp.WriteString("\nMatches: [")

		ipp := pp.Indent()

		for n, e := range f.Matches {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nMatches: []")
	}

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *Command) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("Command {")

	if f.Vars == nil {
		pp.WriteString("\nVars: nil")
	} else if len(f.Vars) > 0 {
		pp.WriteString("\nVars: [")

		ipp := pp.Indent()

		for n, e := range f.Vars {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nVars: []")
	}

	if f.Redirections == nil {
		pp.WriteString("\nRedirections: nil")
	} else if len(f.Redirections) > 0 {
		pp.WriteString("\nRedirections: [")

		ipp := pp.Indent()

		for n, e := range f.Redirections {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nRedirections: []")
	}

	if f.AssignmentsOrWords == nil {
		pp.WriteString("\nAssignmentsOrWords: nil")
	} else if len(f.AssignmentsOrWords) > 0 {
		pp.WriteString("\nAssignmentsOrWords: [")

		ipp := pp.Indent()

		for n, e := range f.AssignmentsOrWords {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nAssignmentsOrWords: []")
	}

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *CommandOrCompound) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("CommandOrCompound {")

	if f.Command != nil {
		pp.WriteString("\nCommand: ")
		f.Command.printType(pp, v)
	} else if v {
		pp.WriteString("\nCommand: nil")
	}

	if f.Compound != nil {
		pp.WriteString("\nCompound: ")
		f.Compound.printType(pp, v)
	} else if v {
		pp.WriteString("\nCompound: nil")
	}

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *CommandSubstitution) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("CommandSubstitution {")

	pp.WriteString("\nSubstitutionType: ")
	f.SubstitutionType.printType(pp, v)

	if f.Backtick != nil {
		pp.WriteString("\nBacktick: ")
		f.Backtick.printType(pp, v)
	} else if v {
		pp.WriteString("\nBacktick: nil")
	}

	pp.WriteString("\nCommand: ")
	f.Command.printType(pp, v)

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *Compound) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("Compound {")

	if f.IfCompound != nil {
		pp.WriteString("\nIfCompound: ")
		f.IfCompound.printType(pp, v)
	} else if v {
		pp.WriteString("\nIfCompound: nil")
	}

	if f.CaseCompound != nil {
		pp.WriteString("\nCaseCompound: ")
		f.CaseCompound.printType(pp, v)
	} else if v {
		pp.WriteString("\nCaseCompound: nil")
	}

	if f.LoopCompound != nil {
		pp.WriteString("\nLoopCompound: ")
		f.LoopCompound.printType(pp, v)
	} else if v {
		pp.WriteString("\nLoopCompound: nil")
	}

	if f.ForCompound != nil {
		pp.WriteString("\nForCompound: ")
		f.ForCompound.printType(pp, v)
	} else if v {
		pp.WriteString("\nForCompound: nil")
	}

	if f.SelectCompound != nil {
		pp.WriteString("\nSelectCompound: ")
		f.SelectCompound.printType(pp, v)
	} else if v {
		pp.WriteString("\nSelectCompound: nil")
	}

	if f.GroupingCompound != nil {
		pp.WriteString("\nGroupingCompound: ")
		f.GroupingCompound.printType(pp, v)
	} else if v {
		pp.WriteString("\nGroupingCompound: nil")
	}

	if f.TestCompound != nil {
		pp.WriteString("\nTestCompound: ")
		f.TestCompound.printType(pp, v)
	} else if v {
		pp.WriteString("\nTestCompound: nil")
	}

	if f.ArithmeticCompound != nil {
		pp.WriteString("\nArithmeticCompound: ")
		f.ArithmeticCompound.printType(pp, v)
	} else if v {
		pp.WriteString("\nArithmeticCompound: nil")
	}

	if f.FunctionCompound != nil {
		pp.WriteString("\nFunctionCompound: ")
		f.FunctionCompound.printType(pp, v)
	} else if v {
		pp.WriteString("\nFunctionCompound: nil")
	}

	if f.Redirections == nil {
		pp.WriteString("\nRedirections: nil")
	} else if len(f.Redirections) > 0 {
		pp.WriteString("\nRedirections: [")

		ipp := pp.Indent()

		for n, e := range f.Redirections {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nRedirections: []")
	}

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *File) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("File {")

	if f.Lines == nil {
		pp.WriteString("\nLines: nil")
	} else if len(f.Lines) > 0 {
		pp.WriteString("\nLines: [")

		ipp := pp.Indent()

		for n, e := range f.Lines {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nLines: []")
	}

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *ForCompound) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("ForCompound {")

	if f.Identifier != nil {
		pp.WriteString("\nIdentifier: ")
		f.Identifier.printType(pp, v)
	} else if v {
		pp.WriteString("\nIdentifier: nil")
	}

	if f.Words == nil {
		pp.WriteString("\nWords: nil")
	} else if len(f.Words) > 0 {
		pp.WriteString("\nWords: [")

		ipp := pp.Indent()

		for n, e := range f.Words {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nWords: []")
	}

	if f.ArithmeticExpansion != nil {
		pp.WriteString("\nArithmeticExpansion: ")
		f.ArithmeticExpansion.printType(pp, v)
	} else if v {
		pp.WriteString("\nArithmeticExpansion: nil")
	}

	pp.WriteString("\nFile: ")
	f.File.printType(pp, v)

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *FunctionCompound) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("FunctionCompound {")

	if f.HasKeyword || v {
		pp.Printf("\nHasKeyword: %v", f.HasKeyword)
	}

	if f.Identifier != nil {
		pp.WriteString("\nIdentifier: ")
		f.Identifier.printType(pp, v)
	} else if v {
		pp.WriteString("\nIdentifier: nil")
	}

	pp.WriteString("\nBody: ")
	f.Body.printType(pp, v)

	pp.WriteString("\nComments: ")
	f.Comments.printType(pp, v)

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *GroupingCompound) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("GroupingCompound {")

	if f.SubShell || v {
		pp.Printf("\nSubShell: %v", f.SubShell)
	}

	pp.WriteString("\nFile: ")
	f.File.printType(pp, v)

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *Heredoc) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("Heredoc {")

	if f.HeredocPartsOrWords == nil {
		pp.WriteString("\nHeredocPartsOrWords: nil")
	} else if len(f.HeredocPartsOrWords) > 0 {
		pp.WriteString("\nHeredocPartsOrWords: [")

		ipp := pp.Indent()

		for n, e := range f.HeredocPartsOrWords {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nHeredocPartsOrWords: []")
	}

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *HeredocPartOrWord) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("HeredocPartOrWord {")

	if f.HeredocPart != nil {
		pp.WriteString("\nHeredocPart: ")
		f.HeredocPart.printType(pp, v)
	} else if v {
		pp.WriteString("\nHeredocPart: nil")
	}

	if f.Word != nil {
		pp.WriteString("\nWord: ")
		f.Word.printType(pp, v)
	} else if v {
		pp.WriteString("\nWord: nil")
	}

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *IfCompound) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("IfCompound {")

	pp.WriteString("\nIf: ")
	f.If.printType(pp, v)

	if f.ElIf == nil {
		pp.WriteString("\nElIf: nil")
	} else if len(f.ElIf) > 0 {
		pp.WriteString("\nElIf: [")

		ipp := pp.Indent()

		for n, e := range f.ElIf {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nElIf: []")
	}

	if f.Else != nil {
		pp.WriteString("\nElse: ")
		f.Else.printType(pp, v)
	} else if v {
		pp.WriteString("\nElse: nil")
	}

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *Line) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("Line {")

	if f.Statements == nil {
		pp.WriteString("\nStatements: nil")
	} else if len(f.Statements) > 0 {
		pp.WriteString("\nStatements: [")

		ipp := pp.Indent()

		for n, e := range f.Statements {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nStatements: []")
	}

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *LoopCompound) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("LoopCompound {")

	if f.Until || v {
		pp.Printf("\nUntil: %v", f.Until)
	}

	pp.WriteString("\nStatement: ")
	f.Statement.printType(pp, v)

	pp.WriteString("\nFile: ")
	f.File.printType(pp, v)

	pp.WriteString("\nComments: ")
	f.Comments.printType(pp, v)

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *Parameter) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("Parameter {")

	if f.Parameter != nil {
		pp.WriteString("\nParameter: ")
		f.Parameter.printType(pp, v)
	} else if v {
		pp.WriteString("\nParameter: nil")
	}

	if f.Array == nil {
		pp.WriteString("\nArray: nil")
	} else if len(f.Array) > 0 {
		pp.WriteString("\nArray: [")

		ipp := pp.Indent()

		for n, e := range f.Array {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nArray: []")
	}

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *ParameterAssign) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("ParameterAssign {")

	if f.Identifier != nil {
		pp.WriteString("\nIdentifier: ")
		f.Identifier.printType(pp, v)
	} else if v {
		pp.WriteString("\nIdentifier: nil")
	}

	if f.Subscript == nil {
		pp.WriteString("\nSubscript: nil")
	} else if len(f.Subscript) > 0 {
		pp.WriteString("\nSubscript: [")

		ipp := pp.Indent()

		for n, e := range f.Subscript {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nSubscript: []")
	}

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *ParameterExpansion) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("ParameterExpansion {")

	if f.Indirect || v {
		pp.Printf("\nIndirect: %v", f.Indirect)
	}

	pp.WriteString("\nParameter: ")
	f.Parameter.printType(pp, v)

	pp.WriteString("\nType: ")
	f.Type.printType(pp, v)

	if f.SubstringStart != nil {
		pp.WriteString("\nSubstringStart: ")
		f.SubstringStart.printType(pp, v)
	} else if v {
		pp.WriteString("\nSubstringStart: nil")
	}

	if f.SubstringEnd != nil {
		pp.WriteString("\nSubstringEnd: ")
		f.SubstringEnd.printType(pp, v)
	} else if v {
		pp.WriteString("\nSubstringEnd: nil")
	}

	if f.BraceWord != nil {
		pp.WriteString("\nBraceWord: ")
		f.BraceWord.printType(pp, v)
	} else if v {
		pp.WriteString("\nBraceWord: nil")
	}

	if f.Pattern != nil {
		pp.WriteString("\nPattern: ")
		f.Pattern.printType(pp, v)
	} else if v {
		pp.WriteString("\nPattern: nil")
	}

	if f.String != nil {
		pp.WriteString("\nString: ")
		f.String.printType(pp, v)
	} else if v {
		pp.WriteString("\nString: nil")
	}

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *Pattern) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("Pattern {")

	if f.Parts == nil {
		pp.WriteString("\nParts: nil")
	} else if len(f.Parts) > 0 {
		pp.WriteString("\nParts: [")

		ipp := pp.Indent()

		for n, e := range f.Parts {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nParts: []")
	}

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *PatternLines) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("PatternLines {")

	if f.Patterns == nil {
		pp.WriteString("\nPatterns: nil")
	} else if len(f.Patterns) > 0 {
		pp.WriteString("\nPatterns: [")

		ipp := pp.Indent()

		for n, e := range f.Patterns {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nPatterns: []")
	}

	pp.WriteString("\nLines: ")
	f.Lines.printType(pp, v)

	pp.WriteString("\nCaseTerminationType: ")
	f.CaseTerminationType.printType(pp, v)

	pp.WriteString("\nComments: ")
	f.Comments.printType(pp, v)

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *Pipeline) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("Pipeline {")

	pp.WriteString("\nPipelineTime: ")
	f.PipelineTime.printType(pp, v)

	if f.Not || v {
		pp.Printf("\nNot: %v", f.Not)
	}

	if f.Coproc || v {
		pp.Printf("\nCoproc: %v", f.Coproc)
	}

	if f.CoprocIdentifier != nil {
		pp.WriteString("\nCoprocIdentifier: ")
		f.CoprocIdentifier.printType(pp, v)
	} else if v {
		pp.WriteString("\nCoprocIdentifier: nil")
	}

	pp.WriteString("\nCommandOrCompound: ")
	f.CommandOrCompound.printType(pp, v)

	if f.Pipeline != nil {
		pp.WriteString("\nPipeline: ")
		f.Pipeline.printType(pp, v)
	} else if v {
		pp.WriteString("\nPipeline: nil")
	}

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *Redirection) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("Redirection {")

	if f.Input != nil {
		pp.WriteString("\nInput: ")
		f.Input.printType(pp, v)
	} else if v {
		pp.WriteString("\nInput: nil")
	}

	if f.Redirector != nil {
		pp.WriteString("\nRedirector: ")
		f.Redirector.printType(pp, v)
	} else if v {
		pp.WriteString("\nRedirector: nil")
	}

	pp.WriteString("\nOutput: ")
	f.Output.printType(pp, v)

	if f.Heredoc != nil {
		pp.WriteString("\nHeredoc: ")
		f.Heredoc.printType(pp, v)
	} else if v {
		pp.WriteString("\nHeredoc: nil")
	}

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *SelectCompound) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("SelectCompound {")

	if f.Identifier != nil {
		pp.WriteString("\nIdentifier: ")
		f.Identifier.printType(pp, v)
	} else if v {
		pp.WriteString("\nIdentifier: nil")
	}

	if f.Words == nil {
		pp.WriteString("\nWords: nil")
	} else if len(f.Words) > 0 {
		pp.WriteString("\nWords: [")

		ipp := pp.Indent()

		for n, e := range f.Words {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nWords: []")
	}

	pp.WriteString("\nFile: ")
	f.File.printType(pp, v)

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *Statement) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("Statement {")

	pp.WriteString("\nPipeline: ")
	f.Pipeline.printType(pp, v)

	pp.WriteString("\nLogicalOperator: ")
	f.LogicalOperator.printType(pp, v)

	if f.Statement != nil {
		pp.WriteString("\nStatement: ")
		f.Statement.printType(pp, v)
	} else if v {
		pp.WriteString("\nStatement: nil")
	}

	pp.WriteString("\nJobControl: ")
	f.JobControl.printType(pp, v)

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *String) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("String {")

	if f.WordsOrTokens == nil {
		pp.WriteString("\nWordsOrTokens: nil")
	} else if len(f.WordsOrTokens) > 0 {
		pp.WriteString("\nWordsOrTokens: [")

		ipp := pp.Indent()

		for n, e := range f.WordsOrTokens {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nWordsOrTokens: []")
	}

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *TestCompound) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("TestCompound {")

	pp.WriteString("\nTests: ")
	f.Tests.printType(pp, v)

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *TestConsequence) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("TestConsequence {")

	pp.WriteString("\nTest: ")
	f.Test.printType(pp, v)

	pp.WriteString("\nConsequence: ")
	f.Consequence.printType(pp, v)

	pp.WriteString("\nComments: ")
	f.Comments.printType(pp, v)

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *Tests) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("Tests {")

	if f.Not || v {
		pp.Printf("\nNot: %v", f.Not)
	}

	pp.WriteString("\nTest: ")
	f.Test.printType(pp, v)

	if f.Word != nil {
		pp.WriteString("\nWord: ")
		f.Word.printType(pp, v)
	} else if v {
		pp.WriteString("\nWord: nil")
	}

	if f.Pattern != nil {
		pp.WriteString("\nPattern: ")
		f.Pattern.printType(pp, v)
	} else if v {
		pp.WriteString("\nPattern: nil")
	}

	if f.Parens != nil {
		pp.WriteString("\nParens: ")
		f.Parens.printType(pp, v)
	} else if v {
		pp.WriteString("\nParens: nil")
	}

	pp.WriteString("\nLogicalOperator: ")
	f.LogicalOperator.printType(pp, v)

	if f.Tests != nil {
		pp.WriteString("\nTests: ")
		f.Tests.printType(pp, v)
	} else if v {
		pp.WriteString("\nTests: nil")
	}

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *Value) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("Value {")

	if f.Word != nil {
		pp.WriteString("\nWord: ")
		f.Word.printType(pp, v)
	} else if v {
		pp.WriteString("\nWord: nil")
	}

	if f.Array == nil {
		pp.WriteString("\nArray: nil")
	} else if len(f.Array) > 0 {
		pp.WriteString("\nArray: [")

		ipp := pp.Indent()

		for n, e := range f.Array {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nArray: []")
	}

	pp.WriteString("\nComments: [")

	ipp := pp.Indent()

	for n, e := range f.Comments {
		ipp.Printf("\n%d: ", n)
		e.printType(ipp, v)
	}

	pp.WriteString("\n]")

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *Word) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("Word {")

	if f.Parts == nil {
		pp.WriteString("\nParts: nil")
	} else if len(f.Parts) > 0 {
		pp.WriteString("\nParts: [")

		ipp := pp.Indent()

		for n, e := range f.Parts {
			ipp.Printf("\n%d: ", n)
			e.printType(ipp, v)
		}

		pp.WriteString("\n]")
	} else if v {
		pp.WriteString("\nParts: []")
	}

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *WordOrOperator) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("WordOrOperator {")

	if f.Word != nil {
		pp.WriteString("\nWord: ")
		f.Word.printType(pp, v)
	} else if v {
		pp.WriteString("\nWord: nil")
	}

	if f.Operator != nil {
		pp.WriteString("\nOperator: ")
		f.Operator.printType(pp, v)
	} else if v {
		pp.WriteString("\nOperator: nil")
	}

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *WordOrToken) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("WordOrToken {")

	if f.Token != nil {
		pp.WriteString("\nToken: ")
		f.Token.printType(pp, v)
	} else if v {
		pp.WriteString("\nToken: nil")
	}

	if f.Word != nil {
		pp.WriteString("\nWord: ")
		f.Word.printType(pp, v)
	} else if v {
		pp.WriteString("\nWord: nil")
	}

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}

func (f *WordPart) printType(w writer, v bool) {
	pp := w.Indent()

	pp.WriteString("WordPart {")

	if f.Part != nil {
		pp.WriteString("\nPart: ")
		f.Part.printType(pp, v)
	} else if v {
		pp.WriteString("\nPart: nil")
	}

	if f.ParameterExpansion != nil {
		pp.WriteString("\nParameterExpansion: ")
		f.ParameterExpansion.printType(pp, v)
	} else if v {
		pp.WriteString("\nParameterExpansion: nil")
	}

	if f.CommandSubstitution != nil {
		pp.WriteString("\nCommandSubstitution: ")
		f.CommandSubstitution.printType(pp, v)
	} else if v {
		pp.WriteString("\nCommandSubstitution: nil")
	}

	if f.ArithmeticExpansion != nil {
		pp.WriteString("\nArithmeticExpansion: ")
		f.ArithmeticExpansion.printType(pp, v)
	} else if v {
		pp.WriteString("\nArithmeticExpansion: nil")
	}

	if f.BraceExpansion != nil {
		pp.WriteString("\nBraceExpansion: ")
		f.BraceExpansion.printType(pp, v)
	} else if v {
		pp.WriteString("\nBraceExpansion: nil")
	}

	pp.WriteString("\nTokens: ")
	f.Tokens.printType(pp, v)

	w.WriteString("\n}")
}
