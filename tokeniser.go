package bash

import (
	"io"
	"slices"
	"strings"

	"vimagination.zapto.org/parser"
)

var (
	keywords           = []string{"if", "then", "else", "elif", "fi", "case", "esac", "while", "for", "in", "do", "done", "time", "until", "coproc", "select", "function", "{", "}", "[[", "]]", "!", "break", "continue"}
	compoundStart      = []string{"if", "while", "until", "for", "select", "{", "("}
	builtins           = []string{"export", "readonly", "declare", "typeset", "local", "let"}
	dotdot             = []string{".."}
	escapedNewline     = []string{"\\\n"}
	assignment         = []string{"=", "+="}
	expansionOperators = [...]string{"#", "%", "^", ","}
	declareParams      = "IAapfnxutrligF"
	typesetParams      = declareParams[2:]
	exportParams       = declareParams[3:6]
	readonlyParams     = declareParams[1:5]
)

const (
	whitespace            = " \t"
	newline               = "\n"
	whitespaceNewline     = whitespace + newline
	heredocsBreak         = whitespace + newline + "|&;()<>\\\"'"
	heredocStringBreak    = newline + "$"
	doubleStops           = "\\`$\""
	singleStops           = "'"
	ansiStops             = "'\\"
	word                  = "\\\"'`(){}- \t\n"
	wordBreak             = "\\\"'`() \t\n$|&;<>{"
	wordBreakBrace        = "\\\"'`() \t\n$|&;,<>}"
	wordBreakArithmetic   = "\\\"'`(){} \t\n$+-!~*/%<=>&^|?:,;"
	wordBreakNoBrace      = wordBreak + "#}]"
	wordBreakIndex        = wordBreakArithmetic + "]"
	wordBreakCommandIndex = "\\\"'`(){} \t\n$+-!~*/%<=>&^|?:,]"
	testWordBreak         = " `\\\t\n\"'$|&;<>(){}!,"
	hexDigit              = "0123456789ABCDEFabcdef"
	octalDigit            = "012345678"
	decimalDigit          = "0123456789"
	letters               = "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz"
	identStart            = letters + "_"
	identCont             = decimalDigit + identStart
	numberChars           = identCont + "@"
)

const (
	TokenWhitespace parser.TokenType = iota
	TokenLineTerminator
	TokenComment
	TokenIdentifier
	TokenFunctionIdentifier
	TokenIdentifierAssign
	TokenLetIdentifierAssign
	TokenAssignment
	TokenKeyword
	TokenBuiltin
	TokenWord
	TokenNumberLiteral
	TokenString
	TokenStringStart
	TokenStringMid
	TokenStringEnd
	TokenBraceSequenceExpansion
	TokenBraceExpansion
	TokenBraceWord
	TokenPunctuator
	TokenHeredoc
	TokenHeredocIndent
	TokenHeredocEnd
	TokenOpenBacktick
	TokenCloseBacktick
	TokenPattern
	TokenOperator
	TokenBinaryOperator
)

type state uint8

const (
	stateNone state = iota
	stateArithmeticExpansion
	stateArithmeticParens
	stateArrayIndex
	stateBrace
	stateBraceExpansion
	stateBraceExpansionWord
	stateBraceExpansionArrayIndex
	stateBuiltinDeclare
	stateBuiltinExport
	stateBuiltinLet
	stateBuiltinLetExpression
	stateBuiltinLetParens
	stateBuiltinLetTernary
	stateBuiltinReadonly
	stateBuiltinTypeset
	stateCaseBody
	stateCaseEnd
	stateCaseParam
	stateCommandIndex
	stateForArithmetic
	stateFunctionBody
	stateHeredoc
	stateHeredocIdentifier
	stateIfBody
	stateIfTest
	stateInCommand
	stateLoopBody
	stateLoopCondition
	stateParens
	stateParensGroup
	stateStringDouble
	stateStringSingle
	stateStringSpecial
	stateTernary
	stateTest
	stateTestBinary
	stateTestPattern
	stateValue
)

type heredocType struct {
	stripped bool
	expand   bool
	delim    string
}

type bashTokeniser struct {
	state                 []state
	heredoc               [][]heredocType
	nextHeredocIsStripped bool
	child                 *parser.Tokeniser
}

// SetTokeniser sets the initial tokeniser state of a parser.Tokeniser.
//
// Used if you want to manually tokenise bash code.
func SetTokeniser(t *parser.Tokeniser) *parser.Tokeniser {
	t.TokeniserState(new(bashTokeniser).main)

	return t
}

func (b *bashTokeniser) lastState() state {
	if len(b.state) == 0 {
		return stateNone
	}

	return b.state[len(b.state)-1]
}

func (b *bashTokeniser) pushState(c state) {
	b.state = append(b.state, c)
}

func (b *bashTokeniser) popState() {
	b.state = b.state[:len(b.state)-1]
}

func (b *bashTokeniser) isInCommand() bool {
	return b.lastState() == stateInCommand
}

func (b *bashTokeniser) endCommand() {
	td := b.lastState()

	if td == stateBuiltinLetExpression {
		b.popState()

		td = b.lastState()
	}

	if td == stateBuiltinLet {
		b.popState()
	}

	if b.isInCommand() {
		b.popState()

		if b.lastState() == stateFunctionBody {
			b.popState()
		}
	}
}

func (b *bashTokeniser) setInCommand() {
	switch b.lastState() {
	case stateArrayIndex, stateBraceExpansionWord, stateBraceExpansionArrayIndex, stateInCommand, stateHeredocIdentifier, stateStringDouble, stateArithmeticExpansion, stateArithmeticParens, stateBraceExpansion, stateCaseParam, stateForArithmetic, stateTest, stateTestBinary, stateValue, stateCommandIndex, stateBuiltinLet, stateBuiltinLetExpression, stateBuiltinLetParens, stateBuiltinLetTernary:
	default:
		b.pushState(stateInCommand)
	}
}

func (b *bashTokeniser) main(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	td := b.lastState()

	if td == stateValue && (isWhitespace(t) || t.Peek() == ';') {
		b.popState()

		td = b.lastState()
	}

	if isWhitespace(t) && td == stateCaseParam {
		return b.caseIn(t)
	} else if td == stateTestPattern {
		b.popState()

		return b.testPattern(t)
	} else if t.Peek() == -1 {
		b.endCommand()

		td = b.lastState()

		if td == stateNone {
			return t.Done()
		}

		return t.ReturnError(io.ErrUnexpectedEOF)
	} else if td == stateHeredocIdentifier {
		b.popState()

		return b.heredocString(t)
	} else if td == stateStringDouble || td == stateStringSingle {
		return b.string(t, false)
	} else if td == stateTest {
		return b.testWord(t)
	} else if td == stateTestBinary {
		return b.testBinaryOperator(t)
	} else if parseWhitespace(t) {
		if td == stateArrayIndex || td == stateBraceExpansionArrayIndex {
			b.popState()

			if !b.isInCommand() {
				b.pushState(td)
			}
		} else if td == stateBuiltinLetExpression {
			b.popState()
		}

		if td == stateCommandIndex {
			return t.Return(TokenWord, b.main)
		}

		return t.Return(TokenWhitespace, b.main)
	} else if t.Accept(newline) {
		b.endCommand()

		if td = b.lastState(); td == stateHeredoc {
			return t.Return(TokenLineTerminator, b.heredocString)
		}

		b.endCommand()
		t.AcceptRun(newline)

		if td == stateIfTest {
			return t.Return(TokenLineTerminator, b.ifThen)
		} else if td == stateLoopCondition {
			return t.Return(TokenLineTerminator, b.loopDo)
		}

		if td == stateCommandIndex {
			return t.Return(TokenWord, b.main)
		}

		return t.Return(TokenLineTerminator, b.main)
	} else if t.Accept("#") {
		if td == stateBraceExpansion || td == stateCommandIndex {
			return b.word(t)
		} else if td == stateArithmeticExpansion || td == stateArithmeticParens || td == stateTernary || td == stateForArithmetic || td == stateArrayIndex || td == stateBuiltinLetExpression || td == stateBuiltinLetParens || td == stateBuiltinLetTernary {
			return t.ReturnError(ErrInvalidCharacter)
		}

		t.ExceptRun(newline)

		return t.Return(TokenComment, b.main)
	} else if td == stateArithmeticExpansion || td == stateArithmeticParens || td == stateTernary || td == stateForArithmetic || td == stateArrayIndex || td == stateCommandIndex || td == stateBuiltinLetExpression || td == stateBuiltinLetParens || td == stateBuiltinLetTernary {
		return b.arithmeticExpansion(t)
	} else if td == stateBuiltinLet {
		return b.letExpressionOrWord(t)
	}

	return b.operatorOrWord(t)
}

func parseWhitespace(t *parser.Tokeniser) bool {
	if t.Accept(whitespace) || t.AcceptWord(escapedNewline, false) != "" {
		for t.AcceptRun(whitespace) != -1 {
			if t.AcceptWord(escapedNewline, false) == "" {
				break
			}
		}

		return true
	}

	return false
}

func (b *bashTokeniser) string(t *parser.Tokeniser, start bool) (parser.Token, parser.TokenFunc) {
	stops := singleStops
	td := b.lastState()
	tk := TokenStringMid

	if td == stateStringDouble {
		stops = doubleStops
	} else if td == stateStringSpecial {
		stops = ansiStops
	}

	if start {
		tk = TokenStringStart
	}

	for {
		switch t.ExceptRun(stops) {
		default:
			return t.ReturnError(io.ErrUnexpectedEOF)
		case '`':
			return t.Return(tk, b.startBacktick)
		case '$':
			return t.Return(tk, b.identifier)
		case '"', '\'':
			t.Next()
			b.popState()

			tk = TokenStringEnd

			if start {
				tk = TokenString
			}

			return t.Return(tk, b.main)
		case '\\':
			t.Next()
			t.Next()
		}
	}
}

func (b *bashTokeniser) arithmeticExpansion(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	switch c := t.Peek(); c {
	case '"', '\'':
		return b.stringStart(t)
	case '$':
		return b.identifier(t)
	case '+', '-', '&', '|':
		t.Next()

		if t.Peek() == c {
			t.Next()
		} else {
			t.Accept("=")
		}
	case '<', '>', '=', '!', '/', '%', '^':
		t.Next()
		t.Accept("=")
	case '*':
		t.Next()
		t.Accept("*=")
	case '~', ',':
		t.Next()
	case '?':
		t.Next()

		if td := b.lastState(); td == stateBuiltinLetExpression || td == stateBuiltinLetParens || td == stateBuiltinLetTernary {
			b.pushState(stateBuiltinLetTernary)
		} else {
			b.pushState(stateTernary)
		}
	case ':':
		t.Next()

		if td := b.lastState(); td != stateTernary && td != stateBuiltinLetTernary {
			return t.ReturnError(ErrInvalidCharacter)
		}

		b.popState()
	case ']':
		t.Next()

		if td := b.lastState(); td == stateCommandIndex {
			return t.Return(TokenPunctuator, b.endCommandIndex)
		} else if td == stateArrayIndex {
			return t.Return(TokenPunctuator, b.startAssign)
		}

		return t.ReturnError(ErrInvalidCharacter)
	case ')':
		t.Next()

		if td := b.lastState(); (td != stateArithmeticExpansion && td != stateForArithmetic || !t.Accept(")")) && td != stateArithmeticParens && td != stateBuiltinLetParens {
			return t.ReturnError(ErrInvalidCharacter)
		}

		b.popState()
	case '(':
		t.Next()

		if td := b.lastState(); td == stateBuiltinLetExpression || td == stateBuiltinLetParens || td == stateBuiltinLetTernary {
			b.pushState(stateBuiltinLetParens)
		} else {
			b.pushState(stateArithmeticParens)
		}
	case '0':
		return b.zero(t)
	case ';':
		if td := b.lastState(); td == stateBuiltinLetExpression {
			b.endCommand()
		} else if td != stateForArithmetic {
			return t.ReturnError(ErrInvalidCharacter)
		}

		t.Next()
	case '{':
		if td := b.lastState(); td == stateBuiltinLetExpression || td == stateBuiltinLetParens || td == stateBuiltinLetTernary {
			t.Next()

			return b.braceExpansion(t)
		}

		fallthrough
	case '}':
		if b.lastState() == stateCommandIndex {
			t.Next()

			return t.Return(TokenPunctuator, b.main)
		}

		return t.ReturnError(ErrInvalidCharacter)
	default:
		return b.number(t)
	}

	return t.Return(TokenPunctuator, b.main)
}

func (b *bashTokeniser) operatorOrWord(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	switch c := t.Peek(); c {
	case '<':
		t.Next()

		if t.Accept("(") {
			b.pushState(stateParens)
		} else if t.Accept("<") {
			if !t.Accept("<") {
				b.nextHeredocIsStripped = t.Accept("-")

				return t.Return(TokenPunctuator, b.startHeredoc)
			}
		} else {
			t.Accept("&>")
		}
	case '>':
		t.Next()

		if t.Accept("(") {
			b.pushState(stateParens)
		} else {
			t.Accept(">&|")
		}
	case '|':
		t.Next()
		t.Accept("&|")
		b.endCommand()
	case '&':
		t.Next()

		if t.Accept(">") {
			t.Accept(">")
		} else {
			b.endCommand()

			if !t.Accept("&") {
				if td := b.lastState(); td == stateIfTest {
					return t.Return(TokenPunctuator, b.ifThen)
				} else if td == stateLoopCondition {
					return t.Return(TokenPunctuator, b.loopDo)
				}
			}
		}
	case ';':
		t.Next()

		l := t.Accept(";")

		if t.Accept("&") {
			l = true
		}

		b.endCommand()

		if l {
			if b.lastState() != stateCaseBody {
				return t.ReturnError(ErrInvalidCharacter)
			} else {
				b.popState()
				b.pushState(stateCaseEnd)
			}
		}

		if td := b.lastState(); td == stateIfTest {
			return t.Return(TokenPunctuator, b.ifThen)
		} else if td == stateLoopCondition {
			return t.Return(TokenPunctuator, b.loopDo)
		}
	case '"', '\'':
		b.setInCommand()

		return b.stringStart(t)
	case '(':
		if b.isInCommand() {
			return t.ReturnError(ErrInvalidCharacter)
		}

		t.Next()
		b.setInCommand()

		if t.Accept("(") {
			b.pushState(stateArithmeticExpansion)
		} else {
			b.setInCommand()
			b.pushState(stateParensGroup)
		}
	case '{':
		t.Next()

		if strings.ContainsRune(whitespaceNewline, t.Peek()) && !b.isInCommand() {
			b.setInCommand()
			b.pushState(stateBrace)
		} else {
			b.setInCommand()

			return b.braceExpansion(t)
		}
	case ')':
		t.Next()
		b.endCommand()

		if td := b.lastState(); td == stateParensGroup {
			b.popState()
			b.endCommand()

			return b.endGroup(t)
		} else if td == stateParens {
			b.popState()
		} else if td == stateCaseEnd {
			b.popState()
			b.pushState(stateCaseBody)
		} else {
			return t.ReturnError(ErrInvalidCharacter)
		}
	case '}':
		t.Next()

		if td := b.lastState(); td == stateBrace {
			b.popState()
			b.endCommand()

			return b.endGroup(t)
		} else if td == stateBraceExpansion {
			b.popState()
			b.endCommand()
		} else if td == stateBraceExpansionWord {
			b.popState()

			return t.Return(TokenBraceExpansion, b.main)
		}
	case '$':
		b.setInCommand()

		return b.identifier(t)
	case '`':
		b.setInCommand()

		return b.startBacktick(t)
	case ',':
		if b.lastState() != stateBraceExpansionWord {
			return b.keywordIdentOrWord(t)
		}

		t.Next()
	case ']':
		if b.lastState() == stateBraceExpansionArrayIndex {
			t.Next()
			b.popState()

			return t.Return(TokenPunctuator, b.parameterExpansionOperation)
		}

		fallthrough
	default:
		return b.keywordIdentOrWord(t)
	}

	return t.Return(TokenPunctuator, b.main)
}

func (b *bashTokeniser) endGroup(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	state := t.State()
	next := b.main

	t.AcceptRun(whitespace)

	switch b.lastState() {
	case stateIfTest:
		if t.AcceptString("then", false) == 4 && isWordSeperator(t) {
			next = b.ifThen
		}
	case stateLoopCondition:
		if t.AcceptString("do", false) == 2 && isWordSeperator(t) {
			next = b.loopDo
		}
	}

	state.Reset()

	return t.Return(TokenPunctuator, next)
}

func (b *bashTokeniser) startBacktick(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	t.Next()

	sub := parser.NewRuneReaderTokeniser(subTokeniser{t})
	b.child = SetTokeniser(&sub)

	return t.Return(TokenOpenBacktick, b.backtick)
}

type subTokeniser struct {
	*parser.Tokeniser
}

func (s subTokeniser) ReadRune() (rune, int, error) {
	if s.Peek() == '`' {
		return -1, 0, io.EOF
	}

	c := s.Next()

	if c == -1 {
		return -1, 0, io.ErrUnexpectedEOF
	} else if c == '\\' {
		switch s.Peek() {
		case -1:
			return -1, 0, io.ErrUnexpectedEOF
		case '\\', '`', '$':
			c = s.Next()
		}
	}

	return c, 1, nil
}

func (b *bashTokeniser) backtick(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	tk, err := b.child.GetToken()

	switch tk.Type {
	case parser.TokenDone:
		if !t.Accept("`") {
			return t.ReturnError(ErrIncorrectBacktick)
		}

		return t.Return(TokenCloseBacktick, b.main)
	case parser.TokenError:
		return t.ReturnError(err)
	}

	pos := t.Len()

	t.Reset()

	for _, c := range tk.Data {
		t.AcceptRun("\\")
		t.AcceptRune(c)
	}

	pos -= t.Len()
	tk.Data = t.Get()

	for t.Len() != pos {
		t.Next()
	}

	return tk, b.backtick
}

func (b *bashTokeniser) startHeredoc(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if t.Peek() == -1 || t.Accept(newline) || t.Accept("#") {
		return t.ReturnError(io.ErrUnexpectedEOF)
	} else if parseWhitespace(t) {
		return t.Return(TokenWhitespace, b.startHeredoc)
	}

	chars := heredocsBreak

Loop:
	for {
		switch t.ExceptRun(chars) {
		case -1:
			return t.ReturnError(io.ErrUnexpectedEOF)
		case '\\':
			t.Next()
			t.Next()
		case '\'':
			t.Next()

			if chars == heredocsBreak {
				chars = "'"
			} else {
				chars = heredocsBreak
			}
		case '"':
			t.Next()

			if chars == heredocsBreak {
				chars = "\\\""
			} else {
				chars = heredocsBreak
			}
		default:
			break Loop
		}
	}

	tk := parser.Token{
		Type: TokenWord,
		Data: t.Get(),
	}
	hdt := heredocType{
		stripped: b.nextHeredocIsStripped,
		delim:    unstring(tk.Data),
	}
	hdt.expand = hdt.delim == tk.Data
	inCommand := b.isInCommand()

	b.endCommand()

	if b.lastState() == stateHeredoc {
		b.heredoc[len(b.heredoc)-1] = append(b.heredoc[len(b.heredoc)-1], hdt)
	} else {
		b.pushState(stateHeredoc)

		b.heredoc = append(b.heredoc, []heredocType{hdt})
	}

	if inCommand {
		b.setInCommand()
	}

	return tk, b.main
}

func unstring(str string) string {
	var sb strings.Builder

	nextEscaped := false

	for _, c := range str {
		if nextEscaped {
			switch c {
			case 'n':
				c = '\n'
			case 't':
				c = '\t'
			}

			nextEscaped = false
		} else {
			switch c {
			case '"', '\'':
				continue
			case '\\':
				nextEscaped = true

				continue
			}
		}

		sb.WriteRune(c)
	}

	return sb.String()
}

func (b *bashTokeniser) heredocString(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	last := len(b.heredoc) - 1
	heredoc := b.heredoc[last][0]

	if heredoc.stripped && t.Accept("\t") {
		t.AcceptRun("\t")

		return t.Return(TokenHeredocIndent, b.heredocString)
	}

	charBreak := newline

	if heredoc.expand {
		charBreak = heredocStringBreak
	}

	for {
		state := t.State()

		if t.AcceptString(heredoc.delim, false) == len(heredoc.delim) && (t.Peek() == '\n' || t.Peek() == -1) {
			state.Reset()

			str := t.Get()

			if len(str) == 0 {
				return b.heredocEnd(t)
			}

			return parser.Token{Type: TokenHeredoc, Data: str}, b.heredocEnd
		}

		switch t.ExceptRun(charBreak) {
		case -1:
			return t.ReturnError(io.ErrUnexpectedEOF)
		case '$':
			state := t.State()

			t.Next()

			if t.Accept(decimalDigit) || t.Accept(identStart) || t.Accept("({$!?") {
				state.Reset()
				b.pushState(stateHeredocIdentifier)

				str := t.Get()

				if len(str) == 0 {
					return b.identifier(t)
				}

				return parser.Token{Type: TokenHeredoc, Data: str}, b.identifier
			}

			continue
		case '\n':
			t.Next()

			if heredoc.stripped && t.Peek() == '\t' {
				return t.Return(TokenHeredoc, b.heredocString)
			}
		}
	}
}

func (b *bashTokeniser) heredocEnd(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	last := len(b.heredoc) - 1
	heredoc := b.heredoc[last][0]
	b.heredoc[last] = b.heredoc[last][1:]

	t.AcceptString(heredoc.delim, false)

	if len(b.heredoc[last]) == 0 {
		b.heredoc = b.heredoc[:last]

		b.popState()
	}

	return t.Return(TokenHeredocEnd, b.main)
}

func (b *bashTokeniser) identifier(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	state := t.State()

	t.Next()

	if t.Accept(decimalDigit) {
		return t.Return(TokenIdentifier, b.main)
	} else if t.Accept("(") {
		if t.Accept("(") {
			b.pushState(stateArithmeticExpansion)

			return t.Return(TokenPunctuator, b.main)
		}

		b.pushState(stateParens)

		return t.Return(TokenPunctuator, b.main)
	} else if t.Accept("{") {
		b.pushState(stateBraceExpansion)

		return t.Return(TokenPunctuator, b.parameterExpansionIdentifierOrPreOperator)
	} else if td := b.lastState(); td != stateStringDouble && td != stateHeredocIdentifier && t.Accept("'\"") {
		state.Reset()

		return b.stringStart(t)
	} else if t.Accept("$!?@*") {
		return t.Return(TokenIdentifier, b.main)
	}

	var wb string

	switch b.lastState() {
	case stateBraceExpansion:
		wb = wordBreakNoBrace
	case stateArrayIndex, stateBraceExpansionArrayIndex:
		wb = wordBreakIndex
	case stateCommandIndex:
		wb = wordBreakCommandIndex
	case stateArithmeticExpansion, stateArithmeticParens, stateTernary, stateForArithmetic:
		wb = wordBreakArithmetic
	case stateTest, stateTestBinary:
		wb = testWordBreak
	default:
		wb = wordBreak
	}

	t.ExceptRun(wb)

	return t.Return(TokenIdentifier, b.main)
}

func (b *bashTokeniser) parameterExpansionIdentifierOrPreOperator(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if t.Accept("!#") {
		if t.Peek() != '}' {
			return t.Return(TokenPunctuator, b.parameterExpansionIdentifier)
		}

		return t.Return(TokenKeyword, b.main)
	}

	return b.parameterExpansionIdentifier(t)
}

func (b *bashTokeniser) parameterExpansionIdentifier(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if t.Accept("@*") {
		return t.Return(TokenKeyword, b.parameterExpansionOperation)
	}

	if t.Accept(decimalDigit) {
		t.AcceptRun(decimalDigit)

		return t.Return(TokenNumberLiteral, b.parameterExpansionOperation)
	}

	if !t.Accept(identStart) {
		return t.ReturnError(ErrInvalidParameterExpansion)
	}

	t.AcceptRun(identCont)

	return t.Return(TokenIdentifier, b.parameterExpansionArrayOrOperation)
}

func (b *bashTokeniser) parameterExpansionArrayOrOperation(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if !t.Accept("[") {
		return b.parameterExpansionOperation(t)
	}

	return t.Return(TokenPunctuator, b.parameterExpansionArraySpecial)
}

func (b *bashTokeniser) parameterExpansionArraySpecial(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if t.Accept("*@") {
		return t.Return(TokenWord, b.parameterExpansionArrayEnd)
	}

	b.pushState(stateBraceExpansionArrayIndex)

	return b.main(t)
}

func (b *bashTokeniser) parameterExpansionArrayEnd(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if !t.Accept("]") {
		return t.ReturnError(ErrInvalidCharacter)
	}

	return t.Return(TokenPunctuator, b.parameterExpansionOperation)
}

func (b *bashTokeniser) parameterExpansionOperation(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if t.Accept(":") {
		if t.Accept("-=?+") {
			return t.Return(TokenPunctuator, b.main)
		}

		return t.Return(TokenPunctuator, b.parameterExpansionSubstringStart)
	} else if t.Accept("-=?+") {
		return t.Return(TokenPunctuator, b.main)
	} else if t.Accept("/") {
		t.Accept("/#%")

		return t.Return(TokenPunctuator, b.parameterExpansionPattern)
	} else if t.Accept("*") {
		return t.Return(TokenPunctuator, b.main)
	} else if t.Accept("@") {
		return t.Return(TokenPunctuator, b.parameterExpansionOperator)
	} else if t.Accept("}") {
		b.popState()

		return t.Return(TokenPunctuator, b.main)
	}

	for _, c := range expansionOperators {
		if t.Accept(c) {
			t.Accept(c)

			return t.Return(TokenPunctuator, b.parameterExpansionPattern)
		}
	}

	return t.ReturnError(ErrInvalidParameterExpansion)
}

func (b *bashTokeniser) parameterExpansionSubstringStart(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if parseWhitespace(t) {
		return t.Return(TokenWhitespace, b.parameterExpansionSubstringStart)
	}

	t.Accept("-")

	if !t.Accept(decimalDigit) {
		return t.ReturnError(ErrInvalidParameterExpansion)
	}

	t.AcceptRun(decimalDigit)

	return t.Return(TokenNumberLiteral, b.parameterExpansionSubstringMid)
}

func (b *bashTokeniser) parameterExpansionSubstringMid(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if parseWhitespace(t) {
		return t.Return(TokenWhitespace, b.parameterExpansionSubstringMid)
	}

	if t.Accept(":") {
		return t.Return(TokenPunctuator, b.parameterExpansionSubstringEnd)
	}

	return b.main(t)
}

func (b *bashTokeniser) parameterExpansionSubstringEnd(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if parseWhitespace(t) {
		return t.Return(TokenWhitespace, b.parameterExpansionSubstringEnd)
	}

	t.Accept("-")

	if !t.Accept(decimalDigit) {
		return t.ReturnError(ErrInvalidParameterExpansion)
	}

	t.AcceptRun(decimalDigit)

	return t.Return(TokenNumberLiteral, b.main)
}

func (b *bashTokeniser) parameterExpansionPattern(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	parens := 0

	for {
		switch t.ExceptRun("\\()[/}") {
		case '}':
			if parens == 0 {
				return t.Return(TokenPattern, b.main)
			}

			return t.ReturnError(ErrInvalidCharacter)
		case '/':
			if parens == 0 {
				return t.Return(TokenPattern, b.parameterExpansionPatternEnd)
			}

			return t.ReturnError(ErrInvalidCharacter)
		case -1:
			return t.ReturnError(io.ErrUnexpectedEOF)
		case '\\':
			t.Next()
			t.Next()
		case '(':
			t.Next()

			parens++
		case ')':
			t.Next()

			if parens == 0 {
				return t.ReturnError(ErrInvalidCharacter)
			}

			parens--
		case '[':
			for !t.Accept("]") {
				switch t.ExceptRun("\\]") {
				case -1:
					return t.ReturnError(io.ErrUnexpectedEOF)
				case '\\':
					t.Next()
					t.Next()
				}
			}
		}
	}
}

func (b *bashTokeniser) parameterExpansionPatternEnd(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	t.Accept("/")

	return t.Return(TokenPunctuator, b.main)
}

func (b *bashTokeniser) parameterExpansionOperator(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if t.Accept("}") {
		b.popState()

		return t.Return(TokenPunctuator, b.main)
	}

	if !t.Accept("UuLQEPAKak") {
		return t.ReturnError(ErrInvalidParameterExpansion)
	}

	return t.Return(TokenBraceWord, b.main)
}

func (b *bashTokeniser) stringStart(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if t.Accept("$") && t.Accept("'") {
		b.pushState(stateStringSpecial)
	} else if t.Accept("'") {
		b.pushState(stateStringSingle)
	} else {
		t.Next()

		b.pushState(stateStringDouble)
	}

	return b.string(t, true)
}

func (b *bashTokeniser) zero(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	t.Next()

	if t.Accept("xX") {
		if !t.Accept(hexDigit) {
			return b.word(t)
		}

		t.AcceptRun(hexDigit)
	} else {
		t.AcceptRun(octalDigit)
	}

	return t.Return(TokenNumberLiteral, b.main)
}

func (b *bashTokeniser) number(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if !t.Accept(decimalDigit) {
		return b.word(t)
	}

	t.AcceptRun(decimalDigit)

	if t.Accept("#") {
		if !t.Accept(numberChars) {
			return b.word(t)
		}

		t.AcceptRun(numberChars)
	}

	return t.Return(TokenNumberLiteral, b.main)
}

func (b *bashTokeniser) keywordIdentOrWord(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if !b.isInCommand() {
		if td := b.lastState(); td != stateTest && td != stateTestBinary {
			state := t.State()
			kw := t.AcceptWord(keywords, false)

			if !isWordSeperator(t) {
				if b.lastState() == stateFunctionBody {
					return t.ReturnError(ErrInvalidKeyword)
				}

				state.Reset()
			} else if kw != "" {
				return b.keyword(t, kw)
			}

			bn := t.AcceptWord(builtins, false)

			if !isWordSeperator(t) {
				state.Reset()
			} else if bn == "let" {
				b.pushState(stateBuiltinLet)

				return t.Return(TokenBuiltin, b.main)
			} else if bn != "" {
				return b.builtin(t, bn)
			}
		}
	}

	return b.identOrWord(t)
}

func (b *bashTokeniser) identOrWord(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if td := b.lastState(); td != stateTest && td != stateTestBinary {
		if t.Accept(identStart) {
			t.AcceptRun(identCont)

			if state := t.State(); t.AcceptWord(assignment, false) != "" {
				state.Reset()

				return t.Return(TokenIdentifierAssign, b.startAssign)
			} else if !b.isInCommand() && b.isCommandIndex(t) {
				b.pushState(stateCommandIndex)

				return t.Return(TokenWord, b.startCommandIndex)
			} else if c := t.Peek(); !b.isInCommand() && c == '[' || b.isArrayStart(t) {
				return t.Return(TokenIdentifierAssign, b.startArrayAssign)
			} else if td := b.lastState(); c == '}' && td == stateBrace || c == ')' && (td == stateParens || td == stateParensGroup) || td == stateBraceExpansion {
				return t.Return(TokenWord, b.main)
			} else if !b.isInCommand() {
				t.AcceptRun(whitespace)

				isFunc := t.Accept("(")

				state.Reset()

				if isFunc {
					return t.Return(TokenFunctionIdentifier, b.functionOpenParen)
				}
			}
		} else if t.Accept(decimalDigit) {
			t.AcceptRun(decimalDigit)

			switch t.Peek() {
			case '<', '>':
				return t.Return(TokenNumberLiteral, b.main)
			}
		}
	}

	return b.word(t)
}

func isWhitespace(t *parser.Tokeniser) bool {
	switch t.Peek() {
	case ' ', '\t', '\n', -1:
		return true
	}

	return false
}

func isWordSeperator(t *parser.Tokeniser) bool {
	return isWhitespace(t) || t.Peek() == ';'
}

func (b *bashTokeniser) isCommandIndex(t *parser.Tokeniser) bool {
	state := t.State()
	defer state.Reset()

	if !t.Accept("[") {
		return false
	}

	return b.hasCompleteBracket(t, stateCommandIndex) && t.AcceptWord(assignment, false) == ""
}

func (b *bashTokeniser) isArrayStart(t *parser.Tokeniser) bool {
	state := t.State()
	defer state.Reset()

	if !t.Accept("[") || t.Accept("]") {
		return false
	}

	return b.hasCompleteBracket(t, stateArrayIndex) && t.AcceptWord(assignment, false) != ""
}

func (b *bashTokeniser) hasCompleteBracket(t *parser.Tokeniser, s state) bool {
	b.pushState(s)
	defer b.popState()

	sub := t.SubTokeniser()

	c := &bashTokeniser{state: b.state}

	sub.TokeniserState(c.main)

	for {
		tk, err := sub.GetToken()
		if err != nil {
			return false
		}

		if len(c.state) == len(b.state) && tk == (parser.Token{Type: TokenPunctuator, Data: "]"}) {
			return true
		} else if len(c.state) < len(b.state) {
			return false
		}
	}
}

func (b *bashTokeniser) startCommandIndex(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	t.Next()

	return t.Return(TokenPunctuator, b.main)
}

func (b *bashTokeniser) endCommandIndex(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	b.popState()
	b.setInCommand()

	return b.main(t)
}

func (b *bashTokeniser) keyword(t *parser.Tokeniser, kw string) (parser.Token, parser.TokenFunc) {
	switch kw {
	case "time":
		if b.lastState() == stateFunctionBody {
			return t.ReturnError(ErrInvalidKeyword)
		}

		return t.Return(TokenKeyword, b.time)
	case "if":
		return t.Return(TokenKeyword, b.ifStart)
	case "then", "in":
		return t.ReturnError(ErrInvalidKeyword)
	case "do":
		if b.lastState() != stateLoopCondition {
			return t.ReturnError(ErrInvalidKeyword)
		}

		b.popState()
		b.pushState(stateLoopBody)

		return t.Return(TokenKeyword, b.main)
	case "elif":
		if b.lastState() != stateIfBody {
			return t.ReturnError(ErrInvalidKeyword)
		}

		b.popState()

		return t.Return(TokenKeyword, b.ifStart)
	case "else":
		if b.lastState() != stateIfBody {
			return t.ReturnError(ErrInvalidKeyword)
		}

		return t.Return(TokenKeyword, b.main)
	case "fi":
		return b.endCompound(t, stateIfBody)
	case "case":
		return t.Return(TokenKeyword, b.caseStart)
	case "esac":
		if td := b.lastState(); td != stateCaseBody && td != stateCaseEnd {
			return t.ReturnError(ErrInvalidKeyword)
		}

		b.popState()

		return t.Return(TokenKeyword, b.main)
	case "while", "until":
		return t.Return(TokenKeyword, b.loopStart)
	case "done":
		return b.endCompound(t, stateLoopBody)
	case "for":
		return t.Return(TokenKeyword, b.forStart)
	case "select":
		return t.Return(TokenKeyword, b.selectStart)
	case "coproc":
		if b.lastState() == stateFunctionBody {
			return t.ReturnError(ErrInvalidKeyword)
		}

		return t.Return(TokenKeyword, b.coproc)
	case "function":
		if b.lastState() == stateFunctionBody {
			return t.ReturnError(ErrInvalidKeyword)
		}

		return t.Return(TokenKeyword, b.function)
	case "[[":
		b.setInCommand()
		b.pushState(stateTest)

		return t.Return(TokenKeyword, b.test)
	case "continue", "break":
		var inLoop bool

	Loop:
		for _, state := range slices.Backward(b.state) {
			switch state {
			case stateIfBody, stateCaseBody:
			case stateLoopBody:
				inLoop = true

				fallthrough
			default:
				break Loop
			}
		}

		if !inLoop {
			return t.ReturnError(ErrInvalidKeyword)
		}

		fallthrough
	default:
		b.setInCommand()

		return t.Return(TokenKeyword, b.main)
	}
}

func (b *bashTokeniser) endCompound(t *parser.Tokeniser, td state) (parser.Token, parser.TokenFunc) {
	if b.lastState() != td {
		return t.ReturnError(ErrInvalidKeyword)
	}

	b.popState()

	return t.Return(TokenKeyword, b.main)
}

func (b *bashTokeniser) time(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if parseWhitespace(t) {
		return t.Return(TokenWhitespace, b.time)
	}

	state := t.State()

	if t.AcceptString("-p", false) == 2 && isWordSeperator(t) {
		return t.Return(TokenWord, b.main)
	}

	state.Reset()

	return b.main(t)
}

func (b *bashTokeniser) startCompound(t *parser.Tokeniser, fn parser.TokenFunc, td state) (parser.Token, parser.TokenFunc) {
	if parseWhitespace(t) {
		return t.Return(TokenWhitespace, fn)
	} else if t.Accept(newline) {
		t.AcceptRun(newline)

		return t.Return(TokenLineTerminator, fn)
	}

	b.pushState(td)

	return b.main(t)
}

func (b *bashTokeniser) middleCompound(t *parser.Tokeniser, fn parser.TokenFunc, kw string, td state, missing error) (parser.Token, parser.TokenFunc) {
	if parseWhitespace(t) {
		return t.Return(TokenWhitespace, fn)
	} else if t.Accept(newline) {
		t.AcceptRun(newline)

		return t.Return(TokenLineTerminator, fn)
	} else if t.Accept("#") {
		t.ExceptRun("\n")

		return t.Return(TokenComment, fn)
	}

	b.popState()

	if t.AcceptString(kw, false) == len(kw) && isWhitespace(t) {
		b.pushState(td)

		return t.Return(TokenKeyword, b.main)
	}

	return t.ReturnError(missing)
}

func (b *bashTokeniser) ifStart(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	return b.startCompound(t, b.ifStart, stateIfTest)
}

func (b *bashTokeniser) ifThen(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	return b.middleCompound(t, b.ifThen, "then", stateIfBody, ErrMissingThen)
}

func (b *bashTokeniser) caseStart(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	return b.startCompound(t, b.caseStart, stateCaseParam)
}

func (b *bashTokeniser) caseIn(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	return b.middleCompound(t, b.caseIn, "in", stateCaseEnd, ErrMissingIn)
}

func (b *bashTokeniser) loopStart(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	return b.startCompound(t, b.loopStart, stateLoopCondition)
}

func (b *bashTokeniser) loopDo(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	return b.middleCompound(t, b.loopDo, "do", stateLoopBody, ErrMissingDo)
}

func (b *bashTokeniser) forStart(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if parseWhitespace(t) {
		return t.Return(TokenWhitespace, b.forStart)
	}

	if t.Accept("(") {
		if !t.Accept("(") {
			return t.ReturnError(ErrInvalidCharacter)
		}

		b.pushState(stateLoopCondition)
		b.pushState(stateForArithmetic)
		b.setInCommand()

		return t.Return(TokenPunctuator, b.main)
	}

	if !t.Accept(identStart) {
		return t.ReturnError(ErrInvalidIdentifier)
	}

	t.AcceptRun(identCont)

	return t.Return(TokenIdentifier, b.forInDo)
}

func (b *bashTokeniser) selectStart(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if parseWhitespace(t) {
		return t.Return(TokenWhitespace, b.selectStart)
	}

	if !t.Accept(identStart) {
		return t.ReturnError(ErrInvalidIdentifier)
	}

	t.AcceptRun(identCont)

	return t.Return(TokenIdentifier, b.forInDo)
}

func (b *bashTokeniser) forInDo(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if parseWhitespace(t) {
		return t.Return(TokenWhitespace, b.forInDo)
	} else if t.Accept(newline) {
		t.AcceptRun(newline)

		return t.Return(TokenLineTerminator, b.forInDo)
	} else if t.Accept("#") {
		t.ExceptRun("\n")

		return t.Return(TokenComment, b.forInDo)
	}

	b.pushState(stateLoopCondition)

	state := t.State()

	if t.AcceptString("in", false) == 2 && isWordSeperator(t) {
		b.setInCommand()

		return t.Return(TokenKeyword, b.main)
	}

	state.Reset()

	return b.main(t)
}

func (b *bashTokeniser) coproc(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if parseWhitespace(t) {
		return t.Return(TokenWhitespace, b.coproc)
	}

	state := t.State()

	if t.AcceptWord(keywords, false) != "" {
		if isWordSeperator(t) {
			state.Reset()

			return b.main(t)
		}

		state.Reset()
	}

	if t.Accept(identStart) {
		t.AcceptRun(identCont)

		nameEnd := t.State()

		if t.Accept(whitespace) {
			t.AcceptRun(whitespace)

			if t.AcceptWord(compoundStart, false) != "" && isWordSeperator(t) {
				nameEnd.Reset()

				return t.Return(TokenIdentifier, b.main)
			}
		}
	}

	state.Reset()

	return b.main(t)
}

func (b *bashTokeniser) function(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if parseWhitespace(t) {
		return t.Return(TokenWhitespace, b.function)
	}

	if !t.Accept(identStart) {
		return t.ReturnError(ErrInvalidIdentifier)
	}

	t.AcceptRun(identCont)

	return t.Return(TokenFunctionIdentifier, b.functionOpenParen)
}

func (b *bashTokeniser) functionOpenParen(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if parseWhitespace(t) {
		return t.Return(TokenWhitespace, b.functionOpenParen)
	}

	b.pushState(stateFunctionBody)

	if t.Accept("(") {
		return t.Return(TokenPunctuator, b.functionCloseParen)
	}

	return b.main(t)
}

func (b *bashTokeniser) functionCloseParen(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if parseWhitespace(t) {
		return t.Return(TokenWhitespace, b.functionCloseParen)
	}

	if !t.Accept(")") {
		return t.ReturnError(ErrMissingClosingParen)
	}

	return t.Return(TokenPunctuator, b.main)
}

func (b *bashTokeniser) test(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if parseWhitespace(t) {
		return t.Return(TokenWhitespace, b.test)
	} else if t.Accept(newline) {
		t.AcceptRun(newline)

		return t.Return(TokenLineTerminator, b.test)
	} else if t.Accept("#") {
		t.ExceptRun("\n")

		return t.Return(TokenComment, b.test)
	} else if t.Accept("!") {
		return t.Return(TokenPunctuator, b.test)
	}

	state := t.State()

	if t.Accept("-") && t.Accept("abcdefghknoprstuvwxzGLNORS") && isWhitespace(t) {
		return t.Return(TokenKeyword, b.testWordStart)
	}

	state.Reset()

	return b.testWordOrPunctuator(t)
}

func (b *bashTokeniser) testWordOrPunctuator(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if parseWhitespace(t) {
		return t.Return(TokenWhitespace, b.testWordOrPunctuator)
	} else if t.Accept(newline) {
		t.AcceptRun(newline)

		return t.Return(TokenLineTerminator, b.testWordOrPunctuator)
	} else if t.Accept("#") {
		t.ExceptRun("\n")

		return t.Return(TokenComment, b.testWordOrPunctuator)
	}

	switch c := t.Peek(); c {
	case -1:
		return t.ReturnError(io.ErrUnexpectedEOF)
	case '(':
		t.Next()
		b.pushState(stateTest)
	case ')':
		t.Next()
		b.popState()

		if b.lastState() != stateTest {
			return t.ReturnError(ErrInvalidCharacter)
		}

		return t.Return(TokenPunctuator, b.testWordOrPunctuator)
	case '|':
		t.Next()

		if !t.Accept("|") {
			return t.ReturnError(ErrInvalidCharacter)
		}
	case '&':
		t.Next()

		if !t.Accept("&") {
			return t.ReturnError(ErrInvalidCharacter)
		}
	case '$':
		b.pushState(stateTestBinary)

		return b.identifier(t)
	case '"', '\'':
		b.pushState(stateTestBinary)

		return b.stringStart(t)
	case ']':
		state := t.State()

		t.Next()

		if t.Accept("]") && isWordSeperator(t) {
			b.popState()

			if b.lastState() == stateTest {
				return t.ReturnError(ErrInvalidCharacter)
			}

			return t.Return(TokenKeyword, b.main)
		}

		state.Reset()

		fallthrough
	default:
		b.pushState(stateTestBinary)

		return b.keywordIdentOrWord(t)
	}

	return t.Return(TokenPunctuator, b.test)
}

func (b *bashTokeniser) testBinaryOperator(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if parseWhitespace(t) {
		return t.Return(TokenWhitespace, b.testBinaryOperator)
	} else if t.Accept(newline) {
		t.AcceptRun(newline)

		return t.Return(TokenLineTerminator, b.testBinaryOperator)
	} else if t.Accept("#") {
		t.ExceptRun(newline)

		return t.Return(TokenComment, b.testBinaryOperator)
	}

	b.popState()

	switch t.Peek() {
	case -1:
		return b.test(t)
	case '=':
		t.Next()
		t.Accept("=~")
	case '!':
		t.Next()

		if !t.Accept("=") {
			return t.ReturnError(ErrInvalidCharacter)
		}
	case '<', '>':
		t.Next()
	case '-':
		t.Next()

		if t.Accept("e") {
			if !t.Accept("qf") {
				return t.ReturnError(ErrInvalidCharacter)
			}
		} else if t.Accept("n") {
			if !t.Accept("et") {
				return t.ReturnError(ErrInvalidCharacter)
			}
		} else if t.Accept("gl") {
			if !t.Accept("et") {
				return t.ReturnError(ErrInvalidCharacter)
			}
		} else if t.Accept("o") {
			if !t.Accept("t") {
				return t.ReturnError(ErrInvalidCharacter)
			}
		} else {
			return t.ReturnError(ErrInvalidCharacter)
		}

		if !isWhitespace(t) {
			return t.ReturnError(ErrInvalidOperator)
		}

		return t.Return(TokenKeyword, b.testWordStart)
	default:
		return b.test(t)
	}

	return t.Return(TokenBinaryOperator, b.testPatternStart)
}

func (b *bashTokeniser) testWordStart(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if parseWhitespace(t) {
		return t.Return(TokenWhitespace, b.testWordStart)
	} else if t.Accept(newline) {
		t.AcceptRun(newline)

		return t.Return(TokenLineTerminator, b.testWordStart)
	} else if t.Accept("#") {
		return t.ReturnError(ErrInvalidCharacter)
	}

	return b.testWord(t)
}

func (b *bashTokeniser) testWord(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if c := t.Peek(); c == '$' {
		return b.identifier(t)
	} else if c == '"' || c == '\'' {
		return b.stringStart(t)
	} else if c == ' ' || c == '\n' {
		return b.test(t)
	} else if c == '`' {
		return b.startBacktick(t)
	}

	return b.testWordOrPunctuator(t)
}

func (b *bashTokeniser) testPatternStart(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if parseWhitespace(t) {
		return t.Return(TokenWhitespace, b.testPatternStart)
	} else if t.Accept(newline) {
		t.AcceptRun(newline)

		return t.Return(TokenLineTerminator, b.testPatternStart)
	} else if t.Accept("#") {
		return t.ReturnError(ErrInvalidCharacter)
	}

	return b.testPattern(t)
}

func (b *bashTokeniser) testPattern(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
Loop:
	for {
		switch t.ExceptRun("\\\"' \t\n$()") {
		default:
			break Loop
		case -1:
			return t.ReturnError(io.ErrUnexpectedEOF)
		case '\\':
			t.Next()
			t.Next()
		case '"', '\'':
			b.pushState(stateTestPattern)

			if t.Len() > 0 {
				return t.Return(TokenPattern, b.stringStart)
			}

			return b.stringStart(t)
		case '$':
			b.pushState(stateTestPattern)

			if t.Len() > 0 {
				return t.Return(TokenPattern, b.identifier)
			}

			return b.identifier(t)
		}
	}

	if t.Len() > 0 {
		return t.Return(TokenPattern, b.test)
	}

	return b.test(t)
}

func (b *bashTokeniser) letExpressionOrWord(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	tk, fn := b.operatorOrWord(t)

	if tk.Type == TokenIdentifierAssign {
		tk.Type = TokenLetIdentifierAssign

		b.pushState(stateBuiltinLetExpression)
	}

	return tk, fn
}

func (b *bashTokeniser) builtin(t *parser.Tokeniser, bn string) (parser.Token, parser.TokenFunc) {
	switch bn {
	case "export":
		b.pushState(stateBuiltinExport)
	case "readonly":
		b.pushState(stateBuiltinReadonly)
	case "typeset":
		b.pushState(stateBuiltinTypeset)
	default:
		b.pushState(stateBuiltinDeclare)
	}

	return t.Return(TokenBuiltin, b.builtinArgs)
}

func (b *bashTokeniser) builtinArgs(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if parseWhitespace(t) {
		return t.Return(TokenWhitespace, b.builtinArgs)
	} else if !t.Accept("-") {
		b.popState()

		return b.main(t)
	}

	params := declareParams

	switch b.lastState() {
	case stateBuiltinExport:
		params = exportParams
	case stateBuiltinReadonly:
		params = readonlyParams
	case stateBuiltinTypeset:
		params = typesetParams
	}

	if !t.Accept(params) {
		return t.ReturnError(ErrInvalidCharacter)
	}

	t.AcceptRun(params)

	return t.Return(TokenOperator, b.builtinArgs)
}

func (b *bashTokeniser) word(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	var wb string

	td := b.lastState()

	switch td {
	case stateBraceExpansion:
		wb = wordBreakNoBrace
	case stateBraceExpansionWord:
		wb = wordBreakBrace
	case stateArrayIndex, stateBraceExpansionArrayIndex:
		wb = wordBreakIndex
	case stateCommandIndex:
		wb = wordBreakCommandIndex
	case stateArithmeticExpansion, stateArithmeticParens, stateTernary, stateForArithmetic, stateBuiltinLetExpression, stateBuiltinLetParens, stateBuiltinLetTernary:
		wb = wordBreakArithmetic
	case stateTest, stateTestBinary:
		wb = testWordBreak
	default:
		wb = wordBreak
	}

	b.setInCommand()

	if t.Accept("\\") && t.Next() == -1 {
		return t.ReturnError(io.ErrUnexpectedEOF)
	}

	if t.Len() == 0 && t.Accept(wb) {
		return t.ReturnError(ErrInvalidCharacter)
	}

	for {
		switch t.ExceptRun(wb) {
		default:
			return t.Return(TokenWord, b.main)
		case '{':
			if td == stateArrayIndex || td == stateBraceExpansionArrayIndex {
				return t.Return(TokenWord, b.main)
			}

			state := t.State()

			t.Next()

			if t.Accept(whitespace) || t.Accept(newline) || t.Peek() == -1 {
				state.Reset()
			} else {
				tk, _ := b.braceExpansion(t.SubTokeniser())

				state.Reset()

				switch tk.Type {
				case TokenBraceExpansion:
					b.popState()

					fallthrough
				case TokenBraceSequenceExpansion:
					return t.Return(TokenWord, b.main)
				}
			}

			t.Next()
		case '\\':
			t.Next()
			t.Next()
		case '$':
			state := t.State()

			t.Next()

			if t.Accept(decimalDigit) || t.Accept(identStart) || t.Accept("({") {
				state.Reset()

				return t.Return(TokenWord, b.main)
			}
		}
	}
}

func (b *bashTokeniser) startArrayAssign(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	t.Accept("[")
	b.pushState(stateArrayIndex)

	return t.Return(TokenPunctuator, b.main)
}

func (b *bashTokeniser) startAssign(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	t.Accept("+")
	t.Accept("=")

	if b.lastState() == stateBuiltinLetExpression {
		return t.Return(TokenAssignment, b.main)
	}

	return t.Return(TokenAssignment, b.value)
}

func (b *bashTokeniser) value(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	isArray := b.lastState() == stateArrayIndex
	if isArray {
		b.popState()
	}

	switch t.Peek() {
	case '(':
		t.Next()

		if isArray || t.Accept("(") {
			return t.ReturnError(ErrInvalidCharacter)
		}

		b.pushState(stateParens)

		return t.Return(TokenPunctuator, b.main)
	case '$':
		return b.identifier(t)
	}

	b.pushState(stateValue)

	return b.main(t)
}

func (b *bashTokeniser) braceExpansion(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	state := t.State()

	if (t.Accept("-") && t.Accept(decimalDigit) || t.Accept(decimalDigit)) && t.AcceptRun(decimalDigit) == '.' && t.AcceptWord(dotdot, false) != "" && (t.Accept("-") && t.Accept(decimalDigit) || t.Accept(decimalDigit)) && (t.AcceptRun(decimalDigit) == '}' || (t.AcceptWord(dotdot, false) != "" && t.Accept("-") && t.Accept(decimalDigit) || t.Accept(decimalDigit) && t.AcceptRun(decimalDigit) == '}')) {
		state.Reset()

		return t.Return(TokenBraceSequenceExpansion, b.braceExpansionSequence)
	}

	state.Reset()

	if t.Accept(letters) && t.AcceptWord(dotdot, false) != "" && t.Accept(letters) && (t.Accept("}") || t.AcceptWord(dotdot, false) != "" && (t.Accept("-") && t.Accept(decimalDigit) || t.Accept(decimalDigit)) && t.AcceptRun(decimalDigit) == '}') {
		state.Reset()

		return t.Return(TokenBraceSequenceExpansion, b.braceExpansionSequence)
	}

	state.Reset()

	bew := b.isBraceExpansionWord(t)

	state.Reset()

	if bew {
		b.pushState(stateBraceExpansionWord)

		return t.Return(TokenBraceExpansion, b.main)
	} else if b.lastState() == stateBuiltinLetExpression {
		return t.ReturnError(ErrInvalidCharacter)
	}

	return b.word(t)
}

func (b *bashTokeniser) braceExpansionSequence(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	typ := TokenNumberLiteral

	if t.Accept(letters) {
		typ = TokenWord
	}

	if t.ExceptRun(".}") == '}' {
		b.pushState(stateBraceExpansionWord)

		return t.Return(typ, b.operatorOrWord)
	}

	return t.Return(typ, b.braceExpansionDelimiter)
}

func (b *bashTokeniser) braceExpansionDelimiter(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	t.Accept(".")
	t.Accept(".")

	return t.Return(TokenPunctuator, b.braceExpansionSequence)
}

func (b *bashTokeniser) isBraceExpansionWord(t *parser.Tokeniser) bool {
	b.pushState(stateBraceExpansionWord)
	defer b.popState()

	var hasComma bool

	sub := t.SubTokeniser()
	c := &bashTokeniser{state: b.state}

	sub.TokeniserState(c.main)

	for {
		tk, err := sub.GetToken()
		if err != nil {
			return false
		}

		if len(c.state) <= len(b.state) {
			switch tk {
			case parser.Token{Type: TokenBraceExpansion, Data: "}"}:
				return hasComma
			case parser.Token{Type: TokenPunctuator, Data: ","}:
				hasComma = true
			case parser.Token{Type: TokenPunctuator, Data: ";"}:
				return false
			default:
				switch tk.Type {
				case TokenWhitespace, TokenLineTerminator:
					return false
				}
			}
		}
	}
}
