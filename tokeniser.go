package bash

import (
	"io"
	"strings"

	"vimagination.zapto.org/parser"
)

var (
	keywords           = []string{"if", "then", "else", "elif", "fi", "case", "esac", "while", "for", "in", "do", "done", "time", "until", "coproc", "select", "function", "{", "}", "[[", "]]", "!", "break", "continue"}
	compoundStart      = []string{"if", "while", "until", "for", "select", "{", "("}
	dotdot             = []string{".."}
	escapedNewline     = []string{"\\\n"}
	assignment         = []string{"=", "+="}
	expansionOperators = [...]string{"#", "%", "^", ","}
)

const (
	whitespace          = " \t"
	newline             = "\n"
	whitespaceNewline   = whitespace + newline
	heredocsBreak       = whitespace + newline + "|&;()<>\\\"'"
	heredocStringBreak  = newline + "$"
	doubleStops         = "\\`$\""
	singleStops         = "'"
	ansiStops           = "'\\"
	word                = "\\\"'`(){}- \t\n"
	wordNoBracket       = "\\\"'`(){}- \t\n]"
	wordBreak           = " `\\\t\n$|&;<>(){"
	wordBreakNoBracket  = wordBreak + "]"
	wordBreakNoBrace    = wordBreak + "}"
	wordBreakArithmetic = "\\\"'`(){} \t\n$+-!~*/%<=>&^|?:,"
	braceWordBreak      = " `\\\t\n|&;<>()={},"
	hexDigit            = "0123456789ABCDEFabcdef"
	octalDigit          = "012345678"
	decimalDigit        = "0123456789"
	letters             = "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz"
	identStart          = letters + "_"
	identCont           = decimalDigit + identStart
	numberChars         = identCont + "@"
)

const (
	TokenWhitespace parser.TokenType = iota
	TokenLineTerminator
	TokenComment
	TokenIdentifier
	TokenFunctionIdentifier
	TokenIdentifierAssign
	TokenKeyword
	TokenWord
	TokenNumberLiteral
	TokenString
	TokenStringStart
	TokenStringMid
	TokenStringEnd
	TokenBraceExpansion
	TokenBraceWord
	TokenPunctuator
	TokenHeredoc
	TokenHeredocIndent
	TokenHeredocEnd
	TokenOpenBacktick
	TokenCloseBacktick
	TokenPattern
)

type heredocType struct {
	stripped bool
	expand   bool
	delim    string
}

type bashTokeniser struct {
	tokenDepth            []byte
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

func (b *bashTokeniser) lastTokenDepth() rune {
	if len(b.tokenDepth) == 0 {
		return 0
	}

	return rune(b.tokenDepth[len(b.tokenDepth)-1])
}

func (b *bashTokeniser) pushTokenDepth(c rune) {
	b.tokenDepth = append(b.tokenDepth, byte(c))
}

func (b *bashTokeniser) popTokenDepth() {
	b.tokenDepth = b.tokenDepth[:len(b.tokenDepth)-1]
}

func (b *bashTokeniser) isInCommand() bool {
	return b.lastTokenDepth() == 'X'
}

func (b *bashTokeniser) endCommand() {
	if b.isInCommand() {
		b.popTokenDepth()

		if b.lastTokenDepth() == 'x' {
			b.popTokenDepth()
		}
	}
}

func (b *bashTokeniser) setInCommand() {
	switch b.lastTokenDepth() {
	case ']', '[', 'X', 'h', '"', '>', '~', 'C', 'f', 't', 'T':
	default:
		b.pushTokenDepth('X')
	}
}

func (b *bashTokeniser) main(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	td := b.lastTokenDepth()

	if isKeywordSeperator(t) && td == 'C' {
		return b.caseIn(t)
	} else if t.Peek() == -1 {
		if b.isInCommand() {
			b.endCommand()

			td = b.lastTokenDepth()
		}

		if td == 'x' {
			b.popTokenDepth()

			td = b.lastTokenDepth()
		}

		if td == 0 {
			return t.Done()
		}

		return t.ReturnError(io.ErrUnexpectedEOF)
	} else if td == 'h' {
		b.popTokenDepth()

		return b.heredocString(t)
	} else if td == '"' || td == '\'' {
		return b.string(t, false)
	} else if parseWhitespace(t) {
		return t.Return(TokenWhitespace, b.main)
	} else if t.Accept(newline) {
		b.endCommand()

		if td = b.lastTokenDepth(); td == 'H' {
			return t.Return(TokenLineTerminator, b.heredocString)
		}

		t.AcceptRun(newline)

		if td == 'I' {
			return t.Return(TokenLineTerminator, b.ifThen)
		} else if td == 'L' {
			return t.Return(TokenLineTerminator, b.loopDo)
		}

		return t.Return(TokenLineTerminator, b.main)
	} else if td == 't' {
		return b.test(t)
	} else if td == 'T' {
		return b.testBinary(t)
	} else if t.Accept("#") {
		if td == '}' || td == '~' {
			return b.word(t)
		}

		t.ExceptRun(newline)

		return t.Return(TokenComment, b.main)
	} else if td == '>' || td == '/' || td == ':' || td == 'f' {
		return b.arithmeticExpansion(t)
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
	td := b.lastTokenDepth()
	tk := TokenStringMid

	if td == '"' {
		stops = doubleStops
	} else if td == '$' {
		stops = ansiStops
	}

	if start {
		tk = TokenStringStart
	}

	for {
		switch t.ExceptRun(stops) {
		default:
			return t.ReturnError(io.ErrUnexpectedEOF)
		case '\n':
			return t.ReturnError(ErrInvalidCharacter)
		case '`':
			return t.Return(tk, b.startBacktick)
		case '$':
			return t.Return(tk, b.identifier)
		case '"', '\'':
			t.Next()
			b.popTokenDepth()

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
	case -1:
		return t.ReturnError(io.ErrUnexpectedEOF)
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
	case '<', '>':
		t.Next()
		t.Accept("=")
	case '=', '!', '/', '%', '^':
		t.Next()
		t.Accept("=")
	case '*':
		t.Next()
		t.Accept("*=")
	case '~', ',':
		t.Next()
	case '?':
		t.Next()
		b.pushTokenDepth(':')
	case ':':
		t.Next()

		if b.lastTokenDepth() != ':' {
			return t.ReturnError(ErrInvalidCharacter)
		}

		b.popTokenDepth()
	case ')':
		t.Next()

		if td := b.lastTokenDepth(); (td != '>' && td != 'f' || !t.Accept(")")) && td != '/' {
			return t.ReturnError(ErrInvalidCharacter)
		}

		b.popTokenDepth()
	case '(':
		t.Next()
		b.pushTokenDepth('/')
	case '0':
		return b.zero(t)
	case ';':
		if b.lastTokenDepth() != 'f' {
			return t.ReturnError(ErrInvalidCharacter)
		}

		t.Next()
	default:
		return b.number(t)
	}

	return t.Return(TokenPunctuator, b.main)
}

func (b *bashTokeniser) operatorOrWord(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	switch c := t.Peek(); c {
	default:
		return b.keywordIdentOrWord(t)
	case '<':
		t.Next()

		if t.Accept("<") {
			if !t.Accept("<") {
				b.nextHeredocIsStripped = t.Accept("-")

				return t.Return(TokenPunctuator, b.startHeredoc)
			}
		} else {
			t.Accept("&>")
		}
	case '>':
		t.Next()
		t.Accept(">&|")
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
				if td := b.lastTokenDepth(); td == 'I' {
					return t.Return(TokenPunctuator, b.ifThen)
				} else if td == 'L' {
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
			if b.lastTokenDepth() != 'c' {
				return t.ReturnError(ErrInvalidCharacter)
			} else {
				b.popTokenDepth()
				b.pushTokenDepth('p')
			}
		}

		if td := b.lastTokenDepth(); td == 'I' {
			return t.Return(TokenPunctuator, b.ifThen)
		} else if td == 'L' {
			return t.Return(TokenPunctuator, b.loopDo)
		}
	case '"', '\'':
		b.setInCommand()

		return b.stringStart(t)
	case '(':
		t.Next()

		if t.Accept("(") {
			b.setInCommand()
			b.pushTokenDepth('>')
		} else {
			b.setInCommand()
			b.pushTokenDepth(')')
		}
	case '{':
		t.Next()

		if tk := t.Peek(); !strings.ContainsRune(word, tk) || tk == '-' {
			b.setInCommand()

			return b.braceExpansion(t)
		} else if strings.ContainsRune(whitespaceNewline, tk) && !b.isInCommand() {
			b.pushTokenDepth('}')
		}
	case ']':
		t.Next()

		if b.lastTokenDepth() == '[' {
			b.popTokenDepth()

			return t.Return(TokenPunctuator, b.parameterExpansionOperation)
		}

		if b.lastTokenDepth() == ']' {
			b.popTokenDepth()
		}
	case ')':
		t.Next()
		b.endCommand()

		if td := b.lastTokenDepth(); td == ')' {
			b.popTokenDepth()
		} else if b.lastTokenDepth() == 'p' {
			b.popTokenDepth()
			b.pushTokenDepth('c')
		} else {
			return t.ReturnError(ErrInvalidCharacter)
		}
	case '}':
		t.Next()

		if td := b.lastTokenDepth(); td == '}' || td == '~' {
			b.popTokenDepth()
		}
	case '+':
		t.Next()
		t.Accept("=")
	case '=':
		t.Next()
	case '$':
		b.setInCommand()

		return b.identifier(t)
	case '`':
		b.setInCommand()

		return b.startBacktick(t)
	}

	return t.Return(TokenPunctuator, b.main)
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

	if b.lastTokenDepth() == 'H' {
		b.heredoc[len(b.heredoc)-1] = append(b.heredoc[len(b.heredoc)-1], hdt)
	} else {
		b.pushTokenDepth('H')

		b.heredoc = append(b.heredoc, []heredocType{hdt})
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
				b.pushTokenDepth('h')

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

		b.popTokenDepth()
	}

	return t.Return(TokenHeredocEnd, b.main)
}

func (b *bashTokeniser) identifier(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	t.Next()

	if t.Accept(decimalDigit) {
		return t.Return(TokenIdentifier, b.main)
	} else if t.Accept("(") {
		if t.Accept("(") {
			b.pushTokenDepth('>')

			return t.Return(TokenPunctuator, b.main)
		}

		b.pushTokenDepth(')')

		return t.Return(TokenPunctuator, b.main)
	} else if t.Accept("{") {
		b.pushTokenDepth('~')

		return t.Return(TokenPunctuator, b.parameterExpansionIdentifierOrPreOperator)
	} else if td := b.lastTokenDepth(); td != '"' && td != 'h' && t.Accept("'\"") {
		t.Reset()

		return b.stringStart(t)
	}

	var wb string

	switch b.lastTokenDepth() {
	case ']':
		wb = wordNoBracket
	case '>':
		wb = wordBreakArithmetic
	default:
		wb = word
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

	b.pushTokenDepth('[')

	return t.Return(TokenPunctuator, b.main)
}

func (b *bashTokeniser) parameterExpansionOperation(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if t.Accept(":") {
		if t.Accept("-=#?+") {
			return t.Return(TokenPunctuator, b.main)
		}

		return t.Return(TokenPunctuator, b.parameterExpansionSubstringStart)
	} else if t.Accept("/") {
		t.Accept("/#%")

		return t.Return(TokenPunctuator, b.parameterExpansionPattern)
	} else if t.Accept("*") {
		return t.Return(TokenPunctuator, b.main)
	} else if t.Accept("@") {
		return t.Return(TokenPunctuator, b.parameterExpansionOperator)
	} else if t.Accept("}") {
		b.popTokenDepth()

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

			return t.ReturnError(io.ErrUnexpectedEOF)
		case '/':
			if parens == 0 {
				return t.Return(TokenPattern, b.parameterExpansionPatternEnd)
			}

			fallthrough
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
				return t.ReturnError(ErrInvalidParameterExpansion)
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
		b.popTokenDepth()

		return t.Return(TokenPunctuator, b.main)
	}

	if !t.Accept("UuLQEPAKak") {
		return t.ReturnError(ErrInvalidParameterExpansion)
	}

	return t.Return(TokenBraceWord, b.main)
}

func (b *bashTokeniser) stringStart(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if b.lastTokenDepth() == t.Peek() {
		b.popTokenDepth()
		t.Next()

		return t.Return(TokenString, b.main)
	} else if t.Accept("$") && t.Accept("'") {
		b.pushTokenDepth('$')
	} else {
		b.pushTokenDepth(t.Next())
	}

	return b.string(t, true)
}

func (b *bashTokeniser) zero(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	t.Next()

	if t.Accept("xX") {
		if !t.Accept(hexDigit) {
			return t.ReturnError(ErrInvalidNumber)
		}

		t.AcceptRun(hexDigit)
	} else {
		t.AcceptRun(octalDigit)
	}

	return t.Return(TokenNumberLiteral, b.main)
}

func (b *bashTokeniser) number(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if !t.Accept(decimalDigit) {
		return b.keywordIdentOrWord(t)
	}

	t.AcceptRun(decimalDigit)

	if t.Accept("#") {
		if !t.Accept(numberChars) {
			return t.ReturnError(ErrInvalidNumber)
		}

		t.AcceptRun(numberChars)
	}

	return t.Return(TokenNumberLiteral, b.main)
}

func (b *bashTokeniser) keywordIdentOrWord(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if !b.isInCommand() {
		if b.lastTokenDepth() != 't' {
			state := t.State()
			kw := t.AcceptWord(keywords, false)

			if !isKeywordSeperator(t) {
				if b.lastTokenDepth() == 'x' {
					return t.ReturnError(ErrInvalidKeyword)
				}

				state.Reset()
			} else if kw != "" {
				return b.keyword(t, kw)
			}
		}
	} else if b.lastTokenDepth() == 'x' {
		return t.ReturnError(ErrInvalidKeyword)
	} else {
		b.endCommand()

		if b.lastTokenDepth() == 'F' {
			state := t.State()
			kw := t.AcceptWord([]string{"in", "do"}, false)

			if !isKeywordSeperator(t) || kw == "" {
				state.Reset()
			} else if kw == "in" {
				b.popTokenDepth()
				b.pushTokenDepth('L')

				return t.Return(TokenKeyword, b.main)
			} else if kw == "do" {
				b.popTokenDepth()
				b.pushTokenDepth('l')

				return t.Return(TokenKeyword, b.main)
			}
		}

		b.setInCommand()
	}

	if t.Accept(identStart) {
		t.AcceptRun(identCont)

		if state := t.State(); t.AcceptWord(assignment, false) != "" {
			state.Reset()

			return t.Return(TokenIdentifierAssign, b.main)
		} else if t.Peek() == '[' {
			return t.Return(TokenIdentifierAssign, b.startArrayAssign)
		} else if td := b.lastTokenDepth(); t.Peek() == td || td == '~' {
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

	return b.word(t)
}

func isKeywordSeperator(t *parser.Tokeniser) bool {
	switch t.Peek() {
	case ' ', '\t', '\n', ';', -1:
		return true
	}

	return false
}

func (b *bashTokeniser) keyword(t *parser.Tokeniser, kw string) (parser.Token, parser.TokenFunc) {
	switch kw {
	case "time":
		if b.lastTokenDepth() == 'x' {
			return t.ReturnError(ErrInvalidKeyword)
		}

		return t.Return(TokenKeyword, b.time)
	case "if":
		return t.Return(TokenKeyword, b.ifStart)
	case "then", "in":
		return t.ReturnError(ErrInvalidKeyword)
	case "do":
		if b.lastTokenDepth() != 'L' {
			return t.ReturnError(ErrInvalidKeyword)
		}

		b.popTokenDepth()
		b.pushTokenDepth('l')

		return t.Return(TokenKeyword, b.main)
	case "elif":
		if b.lastTokenDepth() != 'i' {
			return t.ReturnError(ErrInvalidKeyword)
		}

		b.popTokenDepth()

		return t.Return(TokenKeyword, b.ifStart)
	case "else":
		if b.lastTokenDepth() != 'i' {
			return t.ReturnError(ErrInvalidKeyword)
		}

		return t.Return(TokenKeyword, b.main)
	case "fi":
		return b.endCompound(t, 'i')
	case "case":
		return t.Return(TokenKeyword, b.caseStart)
	case "esac":
		if td := b.lastTokenDepth(); td != 'c' && td != 'p' {
			return t.ReturnError(ErrInvalidKeyword)
		}

		b.popTokenDepth()

		return t.Return(TokenKeyword, b.main)
	case "while", "until":
		return t.Return(TokenKeyword, b.loopStart)
	case "done":
		return b.endCompound(t, 'l')
	case "for":
		return t.Return(TokenKeyword, b.forStart)
	case "select":
		return t.Return(TokenKeyword, b.selectStart)
	case "coproc":
		if b.lastTokenDepth() == 'x' {
			return t.ReturnError(ErrInvalidKeyword)
		}

		return t.Return(TokenKeyword, b.coproc)
	case "function":
		if b.lastTokenDepth() == 'x' {
			return t.ReturnError(ErrInvalidKeyword)
		}

		return t.Return(TokenKeyword, b.function)
	case "[[":
		b.pushTokenDepth('t')

		return t.Return(TokenKeyword, b.test)
	case "continue", "break":
		if td := b.lastTokenDepth(); td != 'l' {
			return t.ReturnError(ErrInvalidKeyword)
		}

		fallthrough
	default:
		b.setInCommand()

		return t.Return(TokenKeyword, b.main)
	}
}

func (b *bashTokeniser) endCompound(t *parser.Tokeniser, td rune) (parser.Token, parser.TokenFunc) {
	if b.lastTokenDepth() != td {
		return t.ReturnError(ErrInvalidKeyword)
	}

	b.popTokenDepth()

	return t.Return(TokenKeyword, b.main)
}

func (b *bashTokeniser) time(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if parseWhitespace(t) {
		return t.Return(TokenWhitespace, b.time)
	}

	state := t.State()

	if t.AcceptString("-p", false) == 2 && isKeywordSeperator(t) {
		return t.Return(TokenWord, b.main)
	}

	state.Reset()

	return b.main(t)
}

func (b *bashTokeniser) startCompound(t *parser.Tokeniser, fn parser.TokenFunc, td rune) (parser.Token, parser.TokenFunc) {
	if parseWhitespace(t) {
		return t.Return(TokenWhitespace, fn)
	} else if t.Accept(newline) {
		t.AcceptRun(newline)

		return t.Return(TokenLineTerminator, fn)
	}

	b.pushTokenDepth(td)

	return b.main(t)
}

func (b *bashTokeniser) middleCompound(t *parser.Tokeniser, fn parser.TokenFunc, kw string, td rune, missing error) (parser.Token, parser.TokenFunc) {
	if parseWhitespace(t) {
		return t.Return(TokenWhitespace, fn)
	} else if t.Accept(newline) {
		t.AcceptRun(newline)

		return t.Return(TokenLineTerminator, fn)
	}

	b.popTokenDepth()

	if t.AcceptString(kw, false) == len(kw) && isKeywordSeperator(t) {
		b.pushTokenDepth(td)

		return t.Return(TokenKeyword, b.main)
	}

	return t.ReturnError(missing)
}

func (b *bashTokeniser) ifStart(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	return b.startCompound(t, b.ifStart, 'I')
}

func (b *bashTokeniser) ifThen(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	return b.middleCompound(t, b.ifThen, "then", 'i', ErrMissingThen)
}

func (b *bashTokeniser) caseStart(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	return b.startCompound(t, b.caseStart, 'C')
}

func (b *bashTokeniser) caseIn(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	return b.middleCompound(t, b.caseIn, "in", 'p', ErrMissingIn)
}

func (b *bashTokeniser) loopStart(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	return b.startCompound(t, b.loopStart, 'L')
}

func (b *bashTokeniser) loopDo(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	return b.middleCompound(t, b.loopDo, "do", 'l', ErrMissingDo)
}

func (b *bashTokeniser) forStart(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if parseWhitespace(t) {
		return t.Return(TokenWhitespace, b.forStart)
	} else if t.Accept(newline) {
		t.AcceptRun(newline)

		return t.Return(TokenLineTerminator, b.forStart)
	}

	if t.Accept("(") {
		if !t.Accept("(") {
			return t.ReturnError(ErrInvalidCharacter)
		}

		b.pushTokenDepth('L')
		b.pushTokenDepth('f')
		b.setInCommand()

		return t.Return(TokenPunctuator, b.main)
	}

	if !t.Accept(identStart) {
		return t.ReturnError(ErrInvalidParameterExpansion)
	}

	t.AcceptRun(identCont)

	return t.Return(TokenIdentifier, b.forInDo)
}

func (b *bashTokeniser) selectStart(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if parseWhitespace(t) {
		return t.Return(TokenWhitespace, b.selectStart)
	} else if t.Accept(newline) {
		t.AcceptRun(newline)

		return t.Return(TokenLineTerminator, b.selectStart)
	}

	if !t.Accept(identStart) {
		return t.ReturnError(ErrInvalidParameterExpansion)
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
	}

	b.pushTokenDepth('L')

	state := t.State()

	if t.AcceptString("in", false) == 2 && isKeywordSeperator(t) {
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
		if isKeywordSeperator(t) {
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

			if t.AcceptWord(compoundStart, false) != "" && isKeywordSeperator(t) {
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

	b.pushTokenDepth('x')

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
	} else if t.Accept("!") {
		return t.Return(TokenPunctuator, b.test)
	}

	state := t.State()

	if t.Accept("-") && t.Accept("abcdefghknoprstuvwxzGLNORS") && isKeywordSeperator(t) {
		return t.Return(TokenPunctuator, b.testWord)
	}

	state.Reset()

	return b.testWord(t)
}

func (b *bashTokeniser) testWord(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if parseWhitespace(t) {
		return t.Return(TokenWhitespace, b.test)
	} else if t.Accept(newline) {
		t.AcceptRun(newline)

		return t.Return(TokenLineTerminator, b.test)
	}

	switch c := t.Peek(); c {
	case -1:
		return t.ReturnError(io.ErrUnexpectedEOF)
	case '(':
		t.Next()
		b.pushTokenDepth('t')
	case ')':
		b.popTokenDepth()

		if b.lastTokenDepth() != 't' {
			return t.ReturnError(ErrInvalidCharacter)
		}
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
		return b.identifier(t)
	case '"', '\'':
		return b.stringStart(t)
	case ']':
		state := t.State()

		t.Next()

		if t.Accept("]") {
			b.popTokenDepth()

			if b.lastTokenDepth() != 't' {
				return t.ReturnError(ErrInvalidCharacter)
			}

			return t.Return(TokenKeyword, b.main)
		}

		state.Reset()

		fallthrough
	default:
		b.pushTokenDepth('T')

		return b.keywordIdentOrWord(t)
	}

	return t.Return(TokenPunctuator, b.test)
}

func (b *bashTokeniser) testBinary(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	switch t.Peek() {
	case -1:
		return t.ReturnError(io.ErrUnexpectedEOF)
	case '=':
		t.Next()
		t.Accept("=~")
	case '!':
		if !t.Accept("=") {
			return t.ReturnError(ErrInvalidCharacter)
		}
	case '<', '>':
	case '-':
		t.Next()

		if t.Accept("e") && !t.Accept("q") || t.Accept("n") && !t.Accept("e") || t.Accept("gl") && !t.Accept("et") {
			return t.ReturnError(ErrInvalidCharacter)
		}
	default:
		return t.ReturnError(ErrInvalidCharacter)
	}

	return t.Return(TokenPunctuator, b.testPattern)
}

func (b *bashTokeniser) testPattern(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	return t.Error()
}

func (b *bashTokeniser) word(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	var wb string

	switch b.lastTokenDepth() {
	case '~':
		wb = wordBreakNoBrace
	case ']', '[':
		wb = wordBreakNoBracket
	case '>', '/', ':', 'f':
		wb = wordBreakArithmetic
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
		case -1:
			if t.Len() == 0 {
				if b.lastTokenDepth() == 0 {
					return t.Done()
				}

				return t.ReturnError(io.ErrUnexpectedEOF)
			}

			fallthrough
		default:
			return t.Return(TokenWord, b.main)
		case '{':
			state := t.State()

			t.Next()

			if t.Accept(whitespace) || t.Accept(newline) || t.Peek() == -1 {
				state.Reset()
			} else {
				tk, _ := b.braceExpansion(t.SubTokeniser())

				state.Reset()

				if tk.Type == TokenBraceExpansion {
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
	b.pushTokenDepth(']')

	return t.Return(TokenPunctuator, b.main)
}

func (b *bashTokeniser) braceExpansion(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if t.Accept(letters) {
		if t.AcceptWord(dotdot, false) != "" {
			if !t.Accept(letters) {
				return b.word(t)
			}

			if t.AcceptWord(dotdot, false) != "" {
				t.Accept("-")

				if !t.Accept(decimalDigit) {
					return b.word(t)
				}

				t.AcceptRun(decimalDigit)
			}

			if !t.Accept("}") {
				return b.word(t)
			}

			return t.Return(TokenBraceExpansion, b.main)
		}

		return b.braceWord(t)
	} else if t.Accept("_") {
		return b.braceWord(t)
	} else {
		t.Accept("-")

		if t.Accept(decimalDigit) {
			switch t.AcceptRun(decimalDigit) {
			default:
				return b.word(t)
			case ',':
				return b.braceExpansionWord(t)
			case '.':
				if t.AcceptWord(dotdot, false) != "" {
					t.Accept("-")

					if !t.Accept(decimalDigit) {
						return b.word(t)
					}

					t.AcceptRun(decimalDigit)

					if t.AcceptWord(dotdot, false) != "" {
						t.Accept("-")

						if !t.Accept(decimalDigit) {
							return b.word(t)
						}

						t.AcceptRun(decimalDigit)
					}

					if !t.Accept("}") {
						return b.word(t)
					}

					return t.Return(TokenBraceExpansion, b.main)
				}

			}
		}
	}

	return b.braceExpansionWord(t)
}

func (b *bashTokeniser) braceWord(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	t.AcceptRun(identCont)

	if !t.Accept("}") {
		return b.braceExpansionWord(t)
	}

	return t.Return(TokenBraceWord, b.main)
}

func (b *bashTokeniser) braceExpansionWord(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	var hasComma bool

	for {
		switch t.ExceptRun(braceWordBreak) {
		case '}':
			if hasComma {
				t.Next()

				return t.Return(TokenBraceExpansion, b.main)
			}

			fallthrough
		default:
			return b.word(t)
		case '\\':
			t.Next()
			t.Next()
		case ',':
			t.Next()

			hasComma = true
		}
	}
}
