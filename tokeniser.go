package bash

import (
	"io"
	"strings"

	"vimagination.zapto.org/parser"
)

var (
	keywords       = []string{"if", "then", "else", "elif", "fi", "case", "esac", "while", "for", "in", "do", "done", "time", "until", "coproc", "select", "function", "{", "}", "[[", "]]", "!"}
	dotdot         = []string{".."}
	escapedNewline = []string{"\\\n"}
	assignment     = []string{"=", "+="}
)

const (
	whitespace         = " \t"
	newline            = "\n"
	metachars          = whitespace + newline + "|&;()<>"
	heredocsBreak      = metachars + "\\\"'"
	heredocStringBreak = newline + "$"
	doubleStops        = "\\`$\""
	singleStops        = "'"
	ansiStops          = "'\\"
	word               = "\\\"'`(){}- \t\n"
	wordNoBracket      = "\\\"'`(){}- \t\n]"
	wordBreak          = " `\\\t\n$|&;<>(){"
	wordBreakNoBracket = wordBreak + "]"
	wordBreakNoBrace   = wordBreak + "}"
	braceWordBreak     = " `\\\t\n|&;<>()={},"
	hexDigit           = "0123456789ABCDEFabcdef"
	octalDigit         = "012345678"
	decimalDigit       = "0123456789"
	letters            = "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz"
	identStart         = letters + "_"
	identCont          = decimalDigit + identStart
	numberChars        = identCont + "@"
)

const (
	TokenWhitespace parser.TokenType = iota
	TokenLineTerminator
	TokenComment
	TokenIdentifier
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
	TokenHeredocEnd
	TokenOpenBacktick
	TokenCloseBacktick
)

type bashTokeniser struct {
	tokenDepth []byte
	heredoc    [][]string
	backticks  int
}

// SetTokeniser sets the initial tokeniser state of a parser.Tokeniser.
//
// Used if you want to manually tokenise bash code.
func SetTokeniser(t *parser.Tokeniser) *parser.Tokeniser {
	p := new(bashTokeniser)

	t.TokeniserState(p.main)

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

func (b *bashTokeniser) main(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	td := b.lastTokenDepth()

	if t.Peek() == -1 {
		if td == 0 {
			return t.Done()
		}

		return t.ReturnError(io.ErrUnexpectedEOF)
	} else if td == 'h' {
		b.popTokenDepth()

		return b.heredocString(t)
	} else if td == '"' || td == '\'' {
		return b.string(t, false)
	} else if t.Accept(whitespace) || t.AcceptWord(escapedNewline, false) != "" {
		for t.AcceptRun(whitespace) != -1 {
			if t.AcceptWord(escapedNewline, false) == "" {
				break
			}
		}

		return t.Return(TokenWhitespace, b.main)
	} else if t.Accept(newline) {
		if td == 'H' {
			return t.Return(TokenLineTerminator, b.heredocString)
		}

		t.AcceptRun(newline)

		return t.Return(TokenLineTerminator, b.main)
	} else if t.Accept("#") {
		t.ExceptRun(newline)

		return t.Return(TokenComment, b.main)
	} else if td == '>' || td == '/' || td == ':' {
		return b.arithmeticExpansion(t)
	}

	return b.operatorOrWord(t)
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
			return t.Return(tk, b.backtickOrIdentOrWord)
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
			if td == '"' && b.isBacktick(t, true) != backtickNone {
				return t.Return(tk, b.backtickOrIdentOrWord)
			}

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

		if td := b.lastTokenDepth(); (td != '>' || !t.Accept(")")) && td != '/' {
			return t.ReturnError(ErrInvalidCharacter)
		}

		b.popTokenDepth()
	case '(':
		t.Next()
		b.pushTokenDepth('/')
	case '0':
		return b.zero(t)
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
			if !t.Accept("<-") {
				t.Accept("-")

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
	case '&':
		t.Next()

		if t.Accept(">") {
			t.Accept(">")
		} else {
			t.Accept("&")
		}
	case ';':
		t.Next()
		t.Accept(";")
		t.Accept("&")
	case '"', '\'':
		return b.stringStart(t)
	case '(':
		t.Next()
		b.pushTokenDepth(')')
	case '{':
		t.Next()

		if !strings.ContainsRune(word, t.Peek()) || t.Peek() == '-' {
			return b.braceExpansion(t)
		}

		b.pushTokenDepth('}')
	case '}', ')', ']':
		t.Next()

		if b.lastTokenDepth() != c {
			return t.ReturnError(ErrInvalidCharacter)
		}

		b.popTokenDepth()
	case '+':
		t.Next()

		if !t.Accept("=") {
			return t.ReturnError(ErrInvalidCharacter)
		}
	case '=':
		t.Next()
	case '$':
		return b.identifier(t)
	case '`', '\\':
		return b.backtickOrIdentOrWord(t)
	}

	return t.Return(TokenPunctuator, b.main)
}

func (b *bashTokeniser) backtickOrIdentOrWord(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	switch b.isBacktick(t, false) {
	case backtickNone:
		return b.keywordIdentOrWord(t)
	case backtickOpen:
		b.pushTokenDepth('`')

		b.backticks++

		return t.Return(TokenOpenBacktick, b.main)
	case backtickClose:
		b.popTokenDepth()

		b.backticks--

		return t.Return(TokenCloseBacktick, b.main)
	default:
		return t.ReturnError(ErrIncorrectBacktick)
	}
}

type backtick uint8

const (
	backtickNone backtick = iota
	backtickOpen
	backtickClose
	backtickInvalid
)

func (b *bashTokeniser) isBacktick(t *parser.Tokeniser, reset bool) backtick {
	if reset {
		defer t.State().Reset()
	}

	pos := t.Len()

	t.AcceptRun("\\")

	slashes := t.Len() - pos

	if slashes == 0 {
		if b.backticks == 0 && t.Accept("`") {
			return backtickOpen
		} else if b.backticks == 1 && b.lastTokenDepth() == '`' && t.Accept("`") {
			return backtickClose
		}
	} else if (b.backticks<<1)-1 == slashes && t.Accept("`") {
		return backtickOpen
	} else if ((b.backticks-1)<<1)-1 == slashes && t.Accept("`") {
		return backtickClose
	}

	if (b.backticks<<1)-1 > slashes && t.Accept("`") {
		return backtickInvalid
	}

	return backtickNone
}

func (b *bashTokeniser) startHeredoc(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if t.Peek() == -1 || t.Accept(newline) || t.Accept("#") {
		return t.ReturnError(io.ErrUnexpectedEOF)
	} else if t.Accept(whitespace) || t.AcceptWord(escapedNewline, false) != "" {
		for t.AcceptRun(whitespace) != -1 {
			if t.AcceptWord(escapedNewline, false) == "" {
				break
			}
		}

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

	if b.lastTokenDepth() == 'H' {
		b.heredoc[len(b.heredoc)-1] = append(b.heredoc[len(b.heredoc)-1], unstring(tk.Data))
	} else {
		b.pushTokenDepth('H')
		b.heredoc = append(b.heredoc, []string{unstring(tk.Data)})
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

	for {
		state := t.State()

		if t.AcceptString(heredoc, false) == len(heredoc) && (t.Peek() == '\n' || t.Peek() == -1) {
			state.Reset()

			return t.Return(TokenHeredoc, b.heredocEnd)
		}

		switch t.ExceptRun(heredocStringBreak) {
		case -1:
			return t.ReturnError(io.ErrUnexpectedEOF)
		case '$':
			state := t.State()

			t.Next()

			if t.Accept(decimalDigit) || t.Accept(identStart) || t.Accept("({") {
				state.Reset()
				b.pushTokenDepth('h')

				return t.Return(TokenHeredoc, b.identifier)
			}

			continue
		}

		t.Next()
	}
}

func (b *bashTokeniser) heredocEnd(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	last := len(b.heredoc) - 1
	heredoc := b.heredoc[last][0]
	b.heredoc[last] = b.heredoc[last][1:]

	t.AcceptString(heredoc, false)

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
		b.pushTokenDepth('}')

		return t.Return(TokenPunctuator, b.keywordIdentOrWord)
	} else if td := b.lastTokenDepth(); td != '"' && td != 'h' && t.Accept("'\"") {
		t.Reset()

		return b.stringStart(t)
	}

	var wb string

	switch b.lastTokenDepth() {
	case ']':
		wb = wordNoBracket
	default:
		wb = word
	}

	t.ExceptRun(wb)

	return t.Return(TokenIdentifier, b.main)
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
	if t.AcceptWord(keywords, false) != "" {
		return t.Return(TokenKeyword, b.main)
	}

	if t.Accept(identStart) {
		t.AcceptRun(identCont)

		if state := t.State(); t.AcceptWord(assignment, false) != "" {
			state.Reset()

			return t.Return(TokenIdentifierAssign, b.main)
		} else if t.Peek() == b.lastTokenDepth() {
			return t.Return(TokenWord, b.main)
		} else if t.Peek() == '[' {
			return t.Return(TokenIdentifierAssign, b.startArrayAssign)
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

func (b *bashTokeniser) word(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	var wb string

	switch b.lastTokenDepth() {
	case '}':
		wb = wordBreakNoBrace
	case ']':
		wb = wordBreakNoBracket
	default:
		wb = wordBreak
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
			if b.isBacktick(t, true) != backtickNone {
				return t.Return(TokenWord, b.backtickOrIdentOrWord)
			}

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
