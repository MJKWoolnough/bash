package bash

import (
	"fmt"
	"io"
	"strings"
	"unsafe"
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
	return i.Write(unsafe.Slice(unsafe.StringData(s), len(s)))
}

func unwrapIndentPrinter(w io.Writer) io.Writer {
	for {
		switch ip := w.(type) {
		case *indentPrinter:
			w = ip.Writer
		default:
			return w
		}
	}
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
	case TokenFunctionIdentifier:
		typ = "TokenFunctionIdentifier"
	case TokenIdentifierAssign:
		typ = "TokenIdentifierAssign"
	case TokenLetIdentifierAssign:
		typ = "TokenLetIdentifierAssign"
	case TokenKeyword:
		typ = "TokenKeyword"
	case TokenBuiltin:
		typ = "TokenBuiltin"
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
	case TokenHeredocIndent:
		typ = "TokenHeredocIndent"
	case TokenHeredocEnd:
		typ = "TokenHeredocEnd"
	case TokenOpenBacktick:
		typ = "TokenOpenBacktick"
	case TokenCloseBacktick:
		typ = "TokenCloseBacktick"
	case TokenPattern:
		typ = "TokenPattern"
	case TokenOperator:
		typ = "TokenOperator"
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

func (c Comments) printType(w io.Writer, v bool) {
	Tokens(c).printType(w, v)
}

func (c Comments) printSource(w io.Writer, v bool) {
	if len(c) > 0 {
		printComment(w, c[0].Data)

		line := c[0].Line

		for _, c := range c[1:] {
			io.WriteString(w, "\n")

			line++

			if line < c.Line {
				io.WriteString(w, "\n")

				line = c.Line
			}

			printComment(w, c.Data)
		}

		if v {
			io.WriteString(w, "\n")
		}
	}
}

func printComment(w io.Writer, c string) {
	if !strings.HasPrefix(c, "#") {
		io.WriteString(w, "#")
	}

	io.WriteString(w, c)
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

func (c CaseTerminationType) String() string {
	switch c {
	case CaseTerminationNone:
		return "CaseTerminationNone"
	case CaseTerminationEnd:
		return "CaseTerminationEnd"
	case CaseTerminationContinue:
		return "CaseTerminationContinue"
	case CaseTerminationFallthrough:
		return "CaseTerminationFallthrough"
	default:
		return "Unknown"
	}
}

func (c CaseTerminationType) printType(w io.Writer, v bool) {
	io.WriteString(w, c.String())
}

func (c CaseTerminationType) printSource(w io.Writer, v bool) {
	switch c {
	case CaseTerminationNone, CaseTerminationEnd:
		io.WriteString(w, ";")
	case CaseTerminationContinue:
		io.WriteString(w, "&")
	case CaseTerminationFallthrough:
		io.WriteString(w, ";&")
	}
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

func (t TestOperator) printSource(w io.Writer, v bool) {
	switch t {
	case TestOperatorFileExists:
		io.WriteString(w, "-e")
	case TestOperatorFileIsBlock:
		io.WriteString(w, "-b")
	case TestOperatorFileIsCharacter:
		io.WriteString(w, "-c")
	case TestOperatorDirectoryExists:
		io.WriteString(w, "-d")
	case TestOperatorFileIsRegular:
		io.WriteString(w, "-f")
	case TestOperatorFileHasSetGroupID:
		io.WriteString(w, "-g")
	case TestOperatorFileIsSymbolic:
		io.WriteString(w, "-L")
	case TestOperatorFileHasStickyBit:
		io.WriteString(w, "-k")
	case TestOperatorFileIsPipe:
		io.WriteString(w, "-p")
	case TestOperatorFileIsReadable:
		io.WriteString(w, "-r")
	case TestOperatorFileIsNonZero:
		io.WriteString(w, "-s")
	case TestOperatorFileIsTerminal:
		io.WriteString(w, "-t")
	case TestOperatorFileHasSetUserID:
		io.WriteString(w, "-u")
	case TestOperatorFileIsWritable:
		io.WriteString(w, "-w")
	case TestOperatorFileIsExecutable:
		io.WriteString(w, "-x")
	case TestOperatorFileIsOwnedByEffectiveGroup:
		io.WriteString(w, "-G")
	case TestOperatorFileWasModifiedSinceLastRead:
		io.WriteString(w, "-N")
	case TestOperatorFileIsOwnedByEffectiveUser:
		io.WriteString(w, "-O")
	case TestOperatorFileIsSocket:
		io.WriteString(w, "-S")
	case TestOperatorFilesAreSameInode:
		io.WriteString(w, "-ef")
	case TestOperatorFileIsNewerThan:
		io.WriteString(w, "-nt")
	case TestOperatorFileIsOlderThan:
		io.WriteString(w, "-ot")
	case TestOperatorOptNameIsEnabled:
		io.WriteString(w, "-o")
	case TestOperatorVarNameIsSet:
		io.WriteString(w, "-v")
	case TestOperatorVarnameIsRef:
		io.WriteString(w, "-R")
	case TestOperatorStringIsZero:
		io.WriteString(w, "-z")
	case TestOperatorStringIsNonZero:
		io.WriteString(w, "-n")
	case TestOperatorStringsEqual:
		io.WriteString(w, "==")
	case TestOperatorStringsMatch:
		io.WriteString(w, "~=")
	case TestOperatorStringsNotEqual:
		io.WriteString(w, "!=")
	case TestOperatorStringBefore:
		io.WriteString(w, "<")
	case TestOperatorStringAfter:
		io.WriteString(w, ">")
	case TestOperatorEqual:
		io.WriteString(w, "-eq")
	case TestOperatorNotEqual:
		io.WriteString(w, "-ne")
	case TestOperatorLessThan:
		io.WriteString(w, "-lt")
	case TestOperatorLessThanEqual:
		io.WriteString(w, "-le")
	case TestOperatorGreaterThan:
		io.WriteString(w, "-gt")
	case TestOperatorGreaterThanEqual:
		io.WriteString(w, "-ge")
	}
}

func (t TestOperator) String() string {
	switch t {
	case TestOperatorNone:
		return "TestOperatorNone"
	case TestOperatorFileExists:
		return "TestOperatorFileExists"
	case TestOperatorFileIsBlock:
		return "TestOperatorFileIsBlock"
	case TestOperatorFileIsCharacter:
		return "TestOperatorFileIsCharacter"
	case TestOperatorDirectoryExists:
		return "TestOperatorDirectoryExists"
	case TestOperatorFileIsRegular:
		return "TestOperatorFileIsRegular"
	case TestOperatorFileHasSetGroupID:
		return "TestOperatorFileHasSetGroupID"
	case TestOperatorFileIsSymbolic:
		return "TestOperatorFileIsSymbolic"
	case TestOperatorFileHasStickyBit:
		return "TestOperatorFileHasStickyBit"
	case TestOperatorFileIsPipe:
		return "TestOperatorFileIsPipe"
	case TestOperatorFileIsReadable:
		return "TestOperatorFileIsReadable"
	case TestOperatorFileIsNonZero:
		return "TestOperatorFileIsNonZero"
	case TestOperatorFileIsTerminal:
		return "TestOperatorFileIsTerminal"
	case TestOperatorFileHasSetUserID:
		return "TestOperatorFileHasSetUserID"
	case TestOperatorFileIsWritable:
		return "TestOperatorFileIsWritable"
	case TestOperatorFileIsExecutable:
		return "TestOperatorFileIsExecutable"
	case TestOperatorFileIsOwnedByEffectiveGroup:
		return "TestOperatorFileIsOwnedByEffectiveGroup"
	case TestOperatorFileWasModifiedSinceLastRead:
		return "TestOperatorFileWasModifiedSinceLastRead"
	case TestOperatorFileIsOwnedByEffectiveUser:
		return "TestOperatorFileIsOwnedByEffectiveUser"
	case TestOperatorFileIsSocket:
		return "TestOperatorFileIsSocket"
	case TestOperatorFilesAreSameInode:
		return "TestOperatorFilesAreSameInode"
	case TestOperatorFileIsNewerThan:
		return "TestOperatorFileIsNewerThan"
	case TestOperatorFileIsOlderThan:
		return "TestOperatorFileIsOlderThan"
	case TestOperatorOptNameIsEnabled:
		return "TestOperatorOptNameIsEnabled"
	case TestOperatorVarNameIsSet:
		return "TestOperatorVarNameIsSet"
	case TestOperatorVarnameIsRef:
		return "TestOperatorVarnameIsRef"
	case TestOperatorStringIsZero:
		return "TestOperatorStringIsZero"
	case TestOperatorStringIsNonZero:
		return "TestOperatorStringIsNonZero"
	case TestOperatorStringsEqual:
		return "TestOperatorStringsEqual"
	case TestOperatorStringsMatch:
		return "TestOperatorStringsMatch"
	case TestOperatorStringsNotEqual:
		return "TestOperatorStringsNotEqual"
	case TestOperatorStringBefore:
		return "TestOperatorStringBefore"
	case TestOperatorStringAfter:
		return "TestOperatorStringAfter"
	case TestOperatorEqual:
		return "TestOperatorEqual"
	case TestOperatorNotEqual:
		return "TestOperatorNotEqual"
	case TestOperatorLessThan:
		return "TestOperatorLessThan"
	case TestOperatorLessThanEqual:
		return "TestOperatorLessThanEqual"
	case TestOperatorGreaterThan:
		return "TestOperatorGreaterThan"
	case TestOperatorGreaterThanEqual:
		return "TestOperatorGreaterThanEqual"
	default:
		return ""
	}
}

func (t TestOperator) printType(w io.Writer, v bool) {
	io.WriteString(w, t.String())
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
