package bash

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unsafe"
)

var (
	indent = []byte{'\t'}
	space  = []byte{' '}
)

type writer interface {
	io.Writer
	io.StringWriter
	Pos() int
	Underlying() writer
}

type indentPrinter struct {
	writer
	hadNewline bool
}

func (i *indentPrinter) Write(p []byte) (int, error) {
	var (
		total int
		last  int
	)

	for n, c := range p {
		if c == '\n' {
			if last != n {
				if err := i.printIndent(); err != nil {
					return total, err
				}
			}

			m, err := i.writer.Write(p[last : n+1])
			total += m

			if err != nil {
				return total, err
			}

			i.hadNewline = true
			last = n + 1
		}
	}

	if last != len(p) {
		if err := i.printIndent(); err != nil {
			return total, err
		}

		m, err := i.writer.Write(p[last:])
		total += m

		if err != nil {
			return total, err
		}
	}

	return total, nil
}

func (i *indentPrinter) printIndent() error {
	if i.hadNewline {
		if _, err := i.writer.Write(indent); err != nil {
			return err
		}

		i.hadNewline = false
	}

	return nil
}

func (i *indentPrinter) Printf(format string, args ...interface{}) {
	fmt.Fprintf(i, format, args...)
}

func (i *indentPrinter) WriteString(s string) (int, error) {
	return i.Write(unsafe.Slice(unsafe.StringData(s), len(s)))
}

type countPrinter struct {
	io.Writer
	pos int
}

func (c *countPrinter) Write(p []byte) (int, error) {
	for _, b := range p {
		if b == '\n' {
			c.pos = 0
		} else if b != '\t' || c.pos > 0 {
			c.pos++
		}
	}

	return c.Writer.Write(p)
}

func (c *countPrinter) WriteString(s string) (int, error) {
	return c.Write(unsafe.Slice(unsafe.StringData(s), len(s)))
}

func (c *countPrinter) Pos() int {
	return c.pos
}

func (c *countPrinter) Underlying() writer {
	return c
}

func (t Token) printType(w writer, v bool) {
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
	case TokenBinaryOperator:
		typ = "TokenBinaryOperator"
	default:
		typ = "Unknown"
	}

	fmt.Fprintf(w, "Type: %s - Data: %q", typ, t.Data)

	if v {
		fmt.Fprintf(w, " - Position: %d (%d: %d)", t.Pos, t.Line, t.LinePos)
	}
}

func (t Tokens) printType(w writer, v bool) {
	if t == nil {
		w.WriteString("nil")

		return
	}

	if len(t) == 0 {
		w.WriteString("[]")

		return
	}

	w.WriteString("[")

	ipp := indentPrinter{writer: &countPrinter{Writer: w}}

	for n, t := range t {
		ipp.Printf("\n%d: ", n)
		t.printType(w, v)
	}

	w.WriteString("\n]")
}

func (c Comments) printType(w writer, v bool) {
	Tokens(c).printType(w, v)
}

func (c Comments) printSource(w writer, v bool) {
	if len(c) > 0 {
		pos := w.Pos()
		line := c[0].Line

		printComment(w, c[0].Data, 0)

		for _, c := range c[1:] {
			w.WriteString("\n")

			line++

			if line < c.Line {
				w.WriteString("\n")

				line = c.Line
			}

			printComment(w, c.Data, pos)
		}

		if v {
			w.WriteString("\n")
		}
	}
}

func printComment(w writer, c string, indent int) {
	w.Write(bytes.Repeat(space, indent))

	if !strings.HasPrefix(c, "#") {
		w.WriteString("#")
	}

	w.WriteString(c)
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

func (a AssignmentType) printSource(w writer, _ bool) {
	switch a {
	case AssignmentAssign:
		w.WriteString("=")
	case AssignmentAppend:
		w.WriteString("+=")
	}
}

func (a AssignmentType) printType(w writer, _ bool) {
	w.WriteString(a.String())
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

func (c CaseTerminationType) printType(w writer, _ bool) {
	w.WriteString(c.String())
}

func (c CaseTerminationType) printSource(w writer, _ bool) {
	switch c {
	case CaseTerminationNone, CaseTerminationEnd:
		w.WriteString(";")
	case CaseTerminationContinue:
		w.WriteString("&")
	case CaseTerminationFallthrough:
		w.WriteString(";&")
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

func (s SubstitutionType) printType(w writer, _ bool) {
	w.WriteString(s.String())
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

func (p PipelineTime) printSource(w writer, _ bool) {
	switch p {
	case PipelineTimeBash:
		w.WriteString("time ")
	case PipelineTimePosix:
		w.WriteString("time -p ")
	}
}

func (p PipelineTime) printType(w writer, _ bool) {
	w.WriteString(p.String())
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

func (l LogicalOperator) printSource(w writer, _ bool) {
	switch l {
	case LogicalOperatorAnd:
		w.WriteString("&&")
	case LogicalOperatorOr:
		w.WriteString("||")
	}
}

func (l LogicalOperator) printType(w writer, _ bool) {
	w.WriteString(l.String())
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

func (j JobControl) printType(w writer, _ bool) {
	w.WriteString(j.String())
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
	case ParameterUnsetSubstitution:
		return "ParameterUnsetSubstitution"
	case ParameterUnsetAssignment:
		return "ParameterUnsetAssignment"
	case ParameterUnsetMessage:
		return "ParameterUnsetMessage"
	case ParameterUnsetSetAssign:
		return "ParameterUnsetSetAssign"
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

func (p ParameterType) printType(w writer, _ bool) {
	w.WriteString(p.String())
}

func (t TestOperator) printSource(w writer, _ bool) {
	switch t {
	case TestOperatorFileExists:
		w.WriteString("-e")
	case TestOperatorFileIsBlock:
		w.WriteString("-b")
	case TestOperatorFileIsCharacter:
		w.WriteString("-c")
	case TestOperatorDirectoryExists:
		w.WriteString("-d")
	case TestOperatorFileIsRegular:
		w.WriteString("-f")
	case TestOperatorFileHasSetGroupID:
		w.WriteString("-g")
	case TestOperatorFileIsSymbolic:
		w.WriteString("-L")
	case TestOperatorFileHasStickyBit:
		w.WriteString("-k")
	case TestOperatorFileIsPipe:
		w.WriteString("-p")
	case TestOperatorFileIsReadable:
		w.WriteString("-r")
	case TestOperatorFileIsNonZero:
		w.WriteString("-s")
	case TestOperatorFileIsTerminal:
		w.WriteString("-t")
	case TestOperatorFileHasSetUserID:
		w.WriteString("-u")
	case TestOperatorFileIsWritable:
		w.WriteString("-w")
	case TestOperatorFileIsExecutable:
		w.WriteString("-x")
	case TestOperatorFileIsOwnedByEffectiveGroup:
		w.WriteString("-G")
	case TestOperatorFileWasModifiedSinceLastRead:
		w.WriteString("-N")
	case TestOperatorFileIsOwnedByEffectiveUser:
		w.WriteString("-O")
	case TestOperatorFileIsSocket:
		w.WriteString("-S")
	case TestOperatorOptNameIsEnabled:
		w.WriteString("-o")
	case TestOperatorVarNameIsSet:
		w.WriteString("-v")
	case TestOperatorVarnameIsRef:
		w.WriteString("-R")
	case TestOperatorStringIsZero:
		w.WriteString("-z")
	case TestOperatorStringIsNonZero:
		w.WriteString("-n")
	case TestOperatorStringsEqual:
		w.WriteString("==")
	case TestOperatorStringsMatch:
		w.WriteString("=~")
	case TestOperatorStringsNotEqual:
		w.WriteString("!=")
	case TestOperatorStringBefore:
		w.WriteString("<")
	case TestOperatorStringAfter:
		w.WriteString(">")
	case TestOperatorEqual:
		w.WriteString("-eq")
	case TestOperatorNotEqual:
		w.WriteString("-ne")
	case TestOperatorLessThan:
		w.WriteString("-lt")
	case TestOperatorLessThanEqual:
		w.WriteString("-le")
	case TestOperatorGreaterThan:
		w.WriteString("-gt")
	case TestOperatorGreaterThanEqual:
		w.WriteString("-ge")
	case TestOperatorFilesAreSameInode:
		w.WriteString("-ef")
	case TestOperatorFileIsNewerThan:
		w.WriteString("-nt")
	case TestOperatorFileIsOlderThan:
		w.WriteString("-ot")
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

func (t TestOperator) printType(w writer, _ bool) {
	w.WriteString(t.String())
}

type formatter interface {
	printType(writer, bool)
	printSource(writer, bool)
}

func format(f formatter, s fmt.State, v rune) {
	switch v {
	case 'v':
		f.printType(&countPrinter{Writer: s}, s.Flag('+'))
	case 's':
		f.printSource(&countPrinter{Writer: s}, s.Flag('+'))
	}
}
