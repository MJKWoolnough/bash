package bash

// File automatically generated with format.sh.

import "io"

func (f *ArithmeticExpansion) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("ArithmeticExpansion {")

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

func (f *File) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("File {")

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

	if f.Index != nil {
		pp.Print("\nIndex: ")
		f.Index.printType(&pp, v)
	} else if v {
		pp.Print("\nIndex: nil")
	}

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

func (f *Pipeline) printType(w io.Writer, v bool) {
	pp := indentPrinter{w}

	pp.Print("Pipeline {")

	pp.Print("\nPipelineTime: ")
	f.PipelineTime.printType(&pp, v)

	if f.Not || v {
		pp.Printf("\nNot: %v", f.Not)
	}

	pp.Print("\nCommand: ")
	f.Command.printType(&pp, v)

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

	if f.LogicalExpression != nil {
		pp.Print("\nLogicalExpression: ")
		f.LogicalExpression.printType(&pp, v)
	} else if v {
		pp.Print("\nLogicalExpression: nil")
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
