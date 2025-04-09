package bash

import (
	"fmt"
	"io"
)

var indent = []byte{'\t'}

type indentPrinter struct {
	io.Writer
}

func (i *indentPrinter) Write(p []byte) (int, error) {
	var (
		total int
		last  int
	)

	for n, c := range p {
		if c == '\n' {
			m, err := i.Writer.Write(p[last : n+1])
			total += m

			if err != nil {
				return total, err
			}

			_, err = i.Writer.Write(indent)
			if err != nil {
				return total, err
			}

			last = n + 1
		}
	}

	if last != len(p) {
		m, err := i.Writer.Write(p[last:])
		total += m

		if err != nil {
			return total, err
		}
	}

	return total, nil
}

func (i *indentPrinter) Print(args ...interface{}) {
	fmt.Fprint(i, args...)
}

func (i *indentPrinter) Printf(format string, args ...interface{}) {
	fmt.Fprintf(i, format, args...)
}

func (i *indentPrinter) WriteString(s string) (int, error) {
	return i.Write([]byte(s))
}

func (t Token) printType(w io.Writer, v bool) {
	var typ string

	switch t.Type {
	case TokenWhitespace:
		typ = "TokenWhitespace"
	case TokenLineTerminator:
		typ = "TokenLineTerminator"
	case TokenComment:
		typ = "TokenComment"
	case TokenIdentifier:
		typ = "TokenIdentifier"
	case TokenIdentifierAssign:
		typ = "TokenIdentifierAssign"
	case TokenKeyword:
		typ = "TokenKeyword"
	case TokenWord:
		typ = "TokenWord"
	case TokenNumberLiteral:
		typ = "TokenNumberLiteral"
	case TokenString:
		typ = "TokenString"
	case TokenStringStart:
		typ = "TokenStringStart"
	case TokenStringMid:
		typ = "TokenStringMid"
	case TokenStringEnd:
		typ = "TokenStringEnd"
	case TokenBraceExpansion:
		typ = "TokenBraceExpansion"
	case TokenBraceWord:
		typ = "TokenBraceWord"
	case TokenPunctuator:
		typ = "TokenPunctuator"
	case TokenHeredoc:
		typ = "TokenHeredoc"
	case TokenHeredocEnd:
		typ = "TokenHeredocEnd"
	case TokenOpenBacktick:
		typ = "TokenOpenBacktick"
	case TokenCloseBacktick:
		typ = "TokenCloseBacktick"
	case TokenCloseParen:
		typ = "TokenCloseParen"
	case TokenPattern:
		typ = "TokenPattern"
	default:
		typ = "Unknown"
	}

	fmt.Fprintf(w, "Type: %s - Data: %q", typ, t.Data)

	if v {
		fmt.Fprintf(w, " - Position: %d (%d: %d)", t.Pos, t.Line, t.LinePos)
	}
}

func (t Tokens) printType(w io.Writer, v bool) {
	if t == nil {
		io.WriteString(w, "nil")

		return
	}

	if len(t) == 0 {
		io.WriteString(w, "[]")

		return
	}

	io.WriteString(w, "[")

	ipp := indentPrinter{w}

	for n, t := range t {
		ipp.Printf("\n%d: ", n)
		t.printType(w, v)
	}

	io.WriteString(w, "\n]")
}

func (a AssignmentType) String() string {
	switch a {
	case AssignmentAssign:
		return "AssignmentAssign"
	case AssignmentAppend:
		return "AssignmentAppend"
	default:
		return "Unknown"
	}
}

func (a AssignmentType) printSource(w io.Writer, v bool) {
	switch a {
	case AssignmentAssign:
		io.WriteString(w, "=")
	case AssignmentAppend:
		io.WriteString(w, "+=")
	}
}

func (a AssignmentType) printType(w io.Writer, v bool) {
	io.WriteString(w, a.String())
}

func (s SubstitutionType) String() string {
	switch s {
	case SubstitutionNew:
		return "SubstitutionNew"
	case SubstitutionBacktick:
		return "SubstitutionBacktick"
	default:
		return "Unknown"
	}
}

func (s SubstitutionType) printType(w io.Writer, v bool) {
	io.WriteString(w, s.String())
}

func (p PipelineTime) String() string {
	switch p {
	case PipelineTimeNone:
		return "PipelineTimeNone"
	case PipelineTimeBash:
		return "PipelineTimeBash"
	case PipelineTimePosix:
		return "PipelineTimePosix"
	default:
		return "Unknown"
	}
}

func (p PipelineTime) printSource(w io.Writer, v bool) {
	switch p {
	case PipelineTimeBash:
		io.WriteString(w, "time ")
	case PipelineTimePosix:
		io.WriteString(w, "time -p ")
	}
}

func (p PipelineTime) printType(w io.Writer, v bool) {
	io.WriteString(w, p.String())
}

func (l LogicalOperator) String() string {
	switch l {
	case LogicalOperatorNone:
		return "LogicalOperatorNone"
	case LogicalOperatorAnd:
		return "LogicalOperatorAnd"
	case LogicalOperatorOr:
		return "LogicalOperatorOr"
	default:
		return "Unknown"
	}
}

func (l LogicalOperator) printSource(w io.Writer, v bool) {
	switch l {
	case LogicalOperatorAnd:
		io.WriteString(w, " && ")
	case LogicalOperatorOr:
		io.WriteString(w, " || ")
	}
}

func (l LogicalOperator) printType(w io.Writer, v bool) {
	io.WriteString(w, l.String())
}

func (j JobControl) String() string {
	switch j {
	case JobControlForeground:
		return "JobControlForeground"
	case JobControlBackground:
		return "JobControlBackground"
	default:
		return "Unknown"
	}
}

func (j JobControl) printType(w io.Writer, v bool) {
	io.WriteString(w, j.String())
}

func (p ParameterType) String() string {
	switch p {
	case ParameterValue:
		return "ParameterValue"
	case ParameterLength:
		return "ParameterLength"
	case ParameterSubstitution:
		return "ParameterSubstitution"
	case ParameterAssignment:
		return "ParameterAssignment"
	case ParameterMessage:
		return "ParameterMessage"
	case ParameterSetAssign:
		return "ParameterSetAssign"
	case ParameterSubstring:
		return "ParameterSubstring"
	case ParameterPrefix:
		return "ParameterPrefix"
	case ParameterPrefixSeperate:
		return "ParameterPrefixSeperate"
	case ParameterRemoveStartShortest:
		return "ParameterRemoveStartShortest"
	case ParameterRemoveStartLongest:
		return "ParameterRemoveStartLongest"
	case ParameterRemoveEndShortest:
		return "ParameterRemoveEndShortest"
	case ParameterRemoveEndLongest:
		return "ParameterRemoveEndLongest"
	case ParameterReplace:
		return "ParameterReplace"
	case ParameterReplaceAll:
		return "ParameterReplaceAll"
	case ParameterReplaceStart:
		return "ParameterReplaceStart"
	case ParameterReplaceEnd:
		return "ParameterReplaceEnd"
	case ParameterLowercaseFirstMatch:
		return "ParameterLowercaseFirstMatch"
	case ParameterLowercaseAllMatches:
		return "ParameterLowercaseAllMatches"
	case ParameterUppercaseFirstMatch:
		return "ParameterUppercaseFirstMatch"
	case ParameterUppercaseAllMatches:
		return "ParameterUppercaseAllMatches"
	case ParameterUppercase:
		return "ParameterUppercase"
	case ParameterUppercaseFirst:
		return "ParameterUppercaseFirst"
	case ParameterLowercase:
		return "ParameterLowercase"
	case ParameterQuoted:
		return "ParameterQuoted"
	case ParameterEscaped:
		return "ParameterEscaped"
	case ParameterPrompt:
		return "ParameterPrompt"
	case ParameterDeclare:
		return "ParameterDeclare"
	case ParameterQuotedArrays:
		return "ParameterQuotedArrays"
	case ParameterQuotedArraysSeperate:
		return ""
	case ParameterAttributes:
		return "ParameterAttributes"
	default:
		return ""
	}
}

func (p ParameterType) printType(w io.Writer, v bool) {
	io.WriteString(w, p.String())
}

type formatter interface {
	printType(io.Writer, bool)
	printSource(io.Writer, bool)
}

func format(f formatter, s fmt.State, v rune) {
	switch v {
	case 'v':
		f.printType(s, s.Flag('+'))
	case 's':
		f.printSource(s, s.Flag('+'))
	}
}
