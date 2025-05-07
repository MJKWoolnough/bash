package bash

// File automatically generated with format.sh.

import "io"

func (f *ArithmeticExpansion) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("ArithmeticExpansion {")

	if f.Expression || v {
		pp.Printf("\nExpression: %v", f.Expression)
	}

	if f.WordsAndOperators == nil {
		pp.Print("\nWordsAndOperators: nil")
	} else if len(f.WordsAndOperators) > 0 {
		pp.Print("\nWordsAndOperators: [")

		ipp := indentPrinter{&pp}

		for n, e := range f.WordsAndOperators {
			ipp.Printf("\n%d: ", n)
			e.printType(&ipp, v)
		}

		pp.Print("\n]")
	} else if v {
		pp.Print("\nWordsAndOperators: []")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *Assignment) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("Assignment {")

	pp.Print("\nIdentifier: ")
	f.Identifier.printType(&pp, v)

	pp.Print("\nAssignment: ")
	f.Assignment.printType(&pp, v)

	pp.Print("\nValue: ")
	f.Value.printType(&pp, v)

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *CaseCompound) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("CaseCompound {")

	pp.Print("\nWord: ")
	f.Word.printType(&pp, v)

	if f.Matches == nil {
		pp.Print("\nMatches: nil")
	} else if len(f.Matches) > 0 {
		pp.Print("\nMatches: [")

		ipp := indentPrinter{&pp}

		for n, e := range f.Matches {
			ipp.Printf("\n%d: ", n)
			e.printType(&ipp, v)
		}

		pp.Print("\n]")
	} else if v {
		pp.Print("\nMatches: []")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *Command) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("Command {")

	if f.Vars == nil {
		pp.Print("\nVars: nil")
	} else if len(f.Vars) > 0 {
		pp.Print("\nVars: [")

		ipp := indentPrinter{&pp}

		for n, e := range f.Vars {
			ipp.Printf("\n%d: ", n)
			e.printType(&ipp, v)
		}

		pp.Print("\n]")
	} else if v {
		pp.Print("\nVars: []")
	}

	if f.Redirections == nil {
		pp.Print("\nRedirections: nil")
	} else if len(f.Redirections) > 0 {
		pp.Print("\nRedirections: [")

		ipp := indentPrinter{&pp}

		for n, e := range f.Redirections {
			ipp.Printf("\n%d: ", n)
			e.printType(&ipp, v)
		}

		pp.Print("\n]")
	} else if v {
		pp.Print("\nRedirections: []")
	}

	if f.Words == nil {
		pp.Print("\nWords: nil")
	} else if len(f.Words) > 0 {
		pp.Print("\nWords: [")

		ipp := indentPrinter{&pp}

		for n, e := range f.Words {
			ipp.Printf("\n%d: ", n)
			e.printType(&ipp, v)
		}

		pp.Print("\n]")
	} else if v {
		pp.Print("\nWords: []")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *CommandOrCompound) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("CommandOrCompound {")

	if f.Command != nil {
		pp.Print("\nCommand: ")
		f.Command.printType(&pp, v)
	} else if v {
		pp.Print("\nCommand: nil")
	}

	if f.Compound != nil {
		pp.Print("\nCompound: ")
		f.Compound.printType(&pp, v)
	} else if v {
		pp.Print("\nCompound: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *CommandSubstitution) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("CommandSubstitution {")

	pp.Print("\nSubstitutionType: ")
	f.SubstitutionType.printType(&pp, v)

	pp.Print("\nCommand: ")
	f.Command.printType(&pp, v)

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *Compound) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("Compound {")

	if f.IfCompound != nil {
		pp.Print("\nIfCompound: ")
		f.IfCompound.printType(&pp, v)
	} else if v {
		pp.Print("\nIfCompound: nil")
	}

	if f.CaseCompound != nil {
		pp.Print("\nCaseCompound: ")
		f.CaseCompound.printType(&pp, v)
	} else if v {
		pp.Print("\nCaseCompound: nil")
	}

	if f.LoopCompound != nil {
		pp.Print("\nLoopCompound: ")
		f.LoopCompound.printType(&pp, v)
	} else if v {
		pp.Print("\nLoopCompound: nil")
	}

	if f.ForCompound != nil {
		pp.Print("\nForCompound: ")
		f.ForCompound.printType(&pp, v)
	} else if v {
		pp.Print("\nForCompound: nil")
	}

	if f.SelectCompound != nil {
		pp.Print("\nSelectCompound: ")
		f.SelectCompound.printType(&pp, v)
	} else if v {
		pp.Print("\nSelectCompound: nil")
	}

	if f.GroupingCompound != nil {
		pp.Print("\nGroupingCompound: ")
		f.GroupingCompound.printType(&pp, v)
	} else if v {
		pp.Print("\nGroupingCompound: nil")
	}

	if f.TestCompound != nil {
		pp.Print("\nTestCompound: ")
		f.TestCompound.printType(&pp, v)
	} else if v {
		pp.Print("\nTestCompound: nil")
	}

	if f.ArthimeticCompound != nil {
		pp.Print("\nArthimeticCompound: ")
		f.ArthimeticCompound.printType(&pp, v)
	} else if v {
		pp.Print("\nArthimeticCompound: nil")
	}

	if f.FunctionCompound != nil {
		pp.Print("\nFunctionCompound: ")
		f.FunctionCompound.printType(&pp, v)
	} else if v {
		pp.Print("\nFunctionCompound: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *File) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("File {")

	if f.Lines == nil {
		pp.Print("\nLines: nil")
	} else if len(f.Lines) > 0 {
		pp.Print("\nLines: [")

		ipp := indentPrinter{&pp}

		for n, e := range f.Lines {
			ipp.Printf("\n%d: ", n)
			e.printType(&ipp, v)
		}

		pp.Print("\n]")
	} else if v {
		pp.Print("\nLines: []")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *ForCompound) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("ForCompound {")

	if f.Identifier != nil {
		pp.Print("\nIdentifier: ")
		f.Identifier.printType(&pp, v)
	} else if v {
		pp.Print("\nIdentifier: nil")
	}

	if f.Words == nil {
		pp.Print("\nWords: nil")
	} else if len(f.Words) > 0 {
		pp.Print("\nWords: [")

		ipp := indentPrinter{&pp}

		for n, e := range f.Words {
			ipp.Printf("\n%d: ", n)
			e.printType(&ipp, v)
		}

		pp.Print("\n]")
	} else if v {
		pp.Print("\nWords: []")
	}

	if f.ArithmeticExpansion != nil {
		pp.Print("\nArithmeticExpansion: ")
		f.ArithmeticExpansion.printType(&pp, v)
	} else if v {
		pp.Print("\nArithmeticExpansion: nil")
	}

	pp.Print("\nFile: ")
	f.File.printType(&pp, v)

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *FunctionCompound) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("FunctionCompound {")

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *GroupingCompound) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("GroupingCompound {")

	if f.SubShell || v {
		pp.Printf("\nSubShell: %v", f.SubShell)
	}

	pp.Print("\nFile: ")
	f.File.printType(&pp, v)

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *Heredoc) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("Heredoc {")

	if f.HeredocPartsOrWords == nil {
		pp.Print("\nHeredocPartsOrWords: nil")
	} else if len(f.HeredocPartsOrWords) > 0 {
		pp.Print("\nHeredocPartsOrWords: [")

		ipp := indentPrinter{&pp}

		for n, e := range f.HeredocPartsOrWords {
			ipp.Printf("\n%d: ", n)
			e.printType(&ipp, v)
		}

		pp.Print("\n]")
	} else if v {
		pp.Print("\nHeredocPartsOrWords: []")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *HeredocPartOrWord) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("HeredocPartOrWord {")

	if f.HeredocPart != nil {
		pp.Print("\nHeredocPart: ")
		f.HeredocPart.printType(&pp, v)
	} else if v {
		pp.Print("\nHeredocPart: nil")
	}

	if f.Word != nil {
		pp.Print("\nWord: ")
		f.Word.printType(&pp, v)
	} else if v {
		pp.Print("\nWord: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *IfCompound) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("IfCompound {")

	pp.Print("\nIf: ")
	f.If.printType(&pp, v)

	if f.ElIf == nil {
		pp.Print("\nElIf: nil")
	} else if len(f.ElIf) > 0 {
		pp.Print("\nElIf: [")

		ipp := indentPrinter{&pp}

		for n, e := range f.ElIf {
			ipp.Printf("\n%d: ", n)
			e.printType(&ipp, v)
		}

		pp.Print("\n]")
	} else if v {
		pp.Print("\nElIf: []")
	}

	if f.Else != nil {
		pp.Print("\nElse: ")
		f.Else.printType(&pp, v)
	} else if v {
		pp.Print("\nElse: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *Line) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("Line {")

	if f.Statements == nil {
		pp.Print("\nStatements: nil")
	} else if len(f.Statements) > 0 {
		pp.Print("\nStatements: [")

		ipp := indentPrinter{&pp}

		for n, e := range f.Statements {
			ipp.Printf("\n%d: ", n)
			e.printType(&ipp, v)
		}

		pp.Print("\n]")
	} else if v {
		pp.Print("\nStatements: []")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *LoopCompound) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("LoopCompound {")

	if f.Until || v {
		pp.Printf("\nUntil: %v", f.Until)
	}

	pp.Print("\nStatement: ")
	f.Statement.printType(&pp, v)

	pp.Print("\nFile: ")
	f.File.printType(&pp, v)

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *Parameter) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("Parameter {")

	if f.Parameter != nil {
		pp.Print("\nParameter: ")
		f.Parameter.printType(&pp, v)
	} else if v {
		pp.Print("\nParameter: nil")
	}

	if f.Array != nil {
		pp.Print("\nArray: ")
		f.Array.printType(&pp, v)
	} else if v {
		pp.Print("\nArray: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *ParameterAssign) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("ParameterAssign {")

	if f.Identifier != nil {
		pp.Print("\nIdentifier: ")
		f.Identifier.printType(&pp, v)
	} else if v {
		pp.Print("\nIdentifier: nil")
	}

	if f.Subscript != nil {
		pp.Print("\nSubscript: ")
		f.Subscript.printType(&pp, v)
	} else if v {
		pp.Print("\nSubscript: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *ParameterExpansion) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("ParameterExpansion {")

	if f.Indirect || v {
		pp.Printf("\nIndirect: %v", f.Indirect)
	}

	pp.Print("\nParameter: ")
	f.Parameter.printType(&pp, v)

	pp.Print("\nType: ")
	f.Type.printType(&pp, v)

	if f.SubstringStart != nil {
		pp.Print("\nSubstringStart: ")
		f.SubstringStart.printType(&pp, v)
	} else if v {
		pp.Print("\nSubstringStart: nil")
	}

	if f.SubstringEnd != nil {
		pp.Print("\nSubstringEnd: ")
		f.SubstringEnd.printType(&pp, v)
	} else if v {
		pp.Print("\nSubstringEnd: nil")
	}

	if f.Word != nil {
		pp.Print("\nWord: ")
		f.Word.printType(&pp, v)
	} else if v {
		pp.Print("\nWord: nil")
	}

	if f.Pattern != nil {
		pp.Print("\nPattern: ")
		f.Pattern.printType(&pp, v)
	} else if v {
		pp.Print("\nPattern: nil")
	}

	if f.String != nil {
		pp.Print("\nString: ")
		f.String.printType(&pp, v)
	} else if v {
		pp.Print("\nString: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *PatternLines) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("PatternLines {")

	if f.Patterns == nil {
		pp.Print("\nPatterns: nil")
	} else if len(f.Patterns) > 0 {
		pp.Print("\nPatterns: [")

		ipp := indentPrinter{&pp}

		for n, e := range f.Patterns {
			ipp.Printf("\n%d: ", n)
			e.printType(&ipp, v)
		}

		pp.Print("\n]")
	} else if v {
		pp.Print("\nPatterns: []")
	}

	pp.Print("\nLines: ")
	f.Lines.printType(&pp, v)

	pp.Print("\nCaseTerminationType: ")
	f.CaseTerminationType.printType(&pp, v)

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *Pipeline) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("Pipeline {")

	pp.Print("\nPipelineTime: ")
	f.PipelineTime.printType(&pp, v)

	if f.Not || v {
		pp.Printf("\nNot: %v", f.Not)
	}

	if f.Coproc || v {
		pp.Printf("\nCoproc: %v", f.Coproc)
	}

	if f.CoprocIdentifier != nil {
		pp.Print("\nCoprocIdentifier: ")
		f.CoprocIdentifier.printType(&pp, v)
	} else if v {
		pp.Print("\nCoprocIdentifier: nil")
	}

	pp.Print("\nCommandOrCompound: ")
	f.CommandOrCompound.printType(&pp, v)

	if f.Pipeline != nil {
		pp.Print("\nPipeline: ")
		f.Pipeline.printType(&pp, v)
	} else if v {
		pp.Print("\nPipeline: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *Redirection) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("Redirection {")

	if f.Input != nil {
		pp.Print("\nInput: ")
		f.Input.printType(&pp, v)
	} else if v {
		pp.Print("\nInput: nil")
	}

	if f.Redirector != nil {
		pp.Print("\nRedirector: ")
		f.Redirector.printType(&pp, v)
	} else if v {
		pp.Print("\nRedirector: nil")
	}

	pp.Print("\nOutput: ")
	f.Output.printType(&pp, v)

	if f.Heredoc != nil {
		pp.Print("\nHeredoc: ")
		f.Heredoc.printType(&pp, v)
	} else if v {
		pp.Print("\nHeredoc: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *SelectCompound) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("SelectCompound {")

	if f.Identifier != nil {
		pp.Print("\nIdentifier: ")
		f.Identifier.printType(&pp, v)
	} else if v {
		pp.Print("\nIdentifier: nil")
	}

	if f.Words == nil {
		pp.Print("\nWords: nil")
	} else if len(f.Words) > 0 {
		pp.Print("\nWords: [")

		ipp := indentPrinter{&pp}

		for n, e := range f.Words {
			ipp.Printf("\n%d: ", n)
			e.printType(&ipp, v)
		}

		pp.Print("\n]")
	} else if v {
		pp.Print("\nWords: []")
	}

	pp.Print("\nFile: ")
	f.File.printType(&pp, v)

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *Statement) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("Statement {")

	pp.Print("\nPipeline: ")
	f.Pipeline.printType(&pp, v)

	pp.Print("\nLogicalOperator: ")
	f.LogicalOperator.printType(&pp, v)

	if f.Statement != nil {
		pp.Print("\nStatement: ")
		f.Statement.printType(&pp, v)
	} else if v {
		pp.Print("\nStatement: nil")
	}

	pp.Print("\nJobControl: ")
	f.JobControl.printType(&pp, v)

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *String) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("String {")

	if f.WordsOrTokens == nil {
		pp.Print("\nWordsOrTokens: nil")
	} else if len(f.WordsOrTokens) > 0 {
		pp.Print("\nWordsOrTokens: [")

		ipp := indentPrinter{&pp}

		for n, e := range f.WordsOrTokens {
			ipp.Printf("\n%d: ", n)
			e.printType(&ipp, v)
		}

		pp.Print("\n]")
	} else if v {
		pp.Print("\nWordsOrTokens: []")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *TestCompound) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("TestCompound {")

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *TestConsequence) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("TestConsequence {")

	pp.Print("\nTest: ")
	f.Test.printType(&pp, v)

	pp.Print("\nConsequence: ")
	f.Consequence.printType(&pp, v)

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *Value) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("Value {")

	if f.Word != nil {
		pp.Print("\nWord: ")
		f.Word.printType(&pp, v)
	} else if v {
		pp.Print("\nWord: nil")
	}

	if f.Array == nil {
		pp.Print("\nArray: nil")
	} else if len(f.Array) > 0 {
		pp.Print("\nArray: [")

		ipp := indentPrinter{&pp}

		for n, e := range f.Array {
			ipp.Printf("\n%d: ", n)
			e.printType(&ipp, v)
		}

		pp.Print("\n]")
	} else if v {
		pp.Print("\nArray: []")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *Word) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("Word {")

	if f.Parts == nil {
		pp.Print("\nParts: nil")
	} else if len(f.Parts) > 0 {
		pp.Print("\nParts: [")

		ipp := indentPrinter{&pp}

		for n, e := range f.Parts {
			ipp.Printf("\n%d: ", n)
			e.printType(&ipp, v)
		}

		pp.Print("\n]")
	} else if v {
		pp.Print("\nParts: []")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *WordOrOperator) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("WordOrOperator {")

	if f.Word != nil {
		pp.Print("\nWord: ")
		f.Word.printType(&pp, v)
	} else if v {
		pp.Print("\nWord: nil")
	}

	if f.Operator != nil {
		pp.Print("\nOperator: ")
		f.Operator.printType(&pp, v)
	} else if v {
		pp.Print("\nOperator: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *WordOrToken) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("WordOrToken {")

	if f.Token != nil {
		pp.Print("\nToken: ")
		f.Token.printType(&pp, v)
	} else if v {
		pp.Print("\nToken: nil")
	}

	if f.Word != nil {
		pp.Print("\nWord: ")
		f.Word.printType(&pp, v)
	} else if v {
		pp.Print("\nWord: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}

func (f *WordPart) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("WordPart {")

	if f.Part != nil {
		pp.Print("\nPart: ")
		f.Part.printType(&pp, v)
	} else if v {
		pp.Print("\nPart: nil")
	}

	if f.ParameterExpansion != nil {
		pp.Print("\nParameterExpansion: ")
		f.ParameterExpansion.printType(&pp, v)
	} else if v {
		pp.Print("\nParameterExpansion: nil")
	}

	if f.CommandSubstitution != nil {
		pp.Print("\nCommandSubstitution: ")
		f.CommandSubstitution.printType(&pp, v)
	} else if v {
		pp.Print("\nCommandSubstitution: nil")
	}

	if f.ArithmeticExpansion != nil {
		pp.Print("\nArithmeticExpansion: ")
		f.ArithmeticExpansion.printType(&pp, v)
	} else if v {
		pp.Print("\nArithmeticExpansion: nil")
	}

	pp.Print("\nTokens: ")
	f.Tokens.printType(&pp, v)

	io.WriteString(w, "\n}")
}
