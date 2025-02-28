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
	doubleStops        = "\\\n`$\""
	singleStops        = "\n'"
	word               = "\\\"'`(){}- \t\n"
	wordNoBracket      = "\\\"'`(){}- \t\n]"
	wordBreak          = " `\\\t\n|&;<>()"
	wordBreakNoBracket = wordBreak + "]"
	wordBreakNoBrace   = wordBreak + "}"
	braceBreak         = " `\\\t\n|&;<>()=},"
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
	TokenPunctuator
	TokenHeredoc
)

type bashTokeniser struct {
	tokenDepth []byte
	heredoc    string
}

// SetTokeniser sets the initial tokeniser state of a parser.Tokeniser.
//
// Used if you want to manually tokenise bash code.
func SetTokeniser(t *parser.Tokeniser) *parser.Tokeniser {
	p := new(bashTokeniser)

	t.TokeniserState(p.main)

	return t
}

func (b *bashTokeniser) lastTokenDepth() byte {
	if len(b.tokenDepth) == 0 {
		return 0
	}

	return b.tokenDepth[len(b.tokenDepth)-1]
}

func (b *bashTokeniser) pushTokenDepth(c byte) {
	b.tokenDepth = append(b.tokenDepth, c)
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
	}

	if td == '"' || td == '\'' {
		return b.string(t)
	}

	if t.Accept(whitespace) || t.AcceptWord(escapedNewline, false) != "" {
		for t.AcceptRun(whitespace) != -1 {
			if t.AcceptWord(escapedNewline, false) == "" {
				break
			}
		}

		return t.Return(TokenWhitespace, b.main)
	}

	if t.Accept(newline) {
		if b.heredoc != "" {
			return t.Return(TokenLineTerminator, b.heredocString)
		}

		t.AcceptRun(newline)

		return t.Return(TokenLineTerminator, b.main)
	}

	if t.Accept("#") {
		t.ExceptRun(newline)

		return t.Return(TokenComment, b.main)
	}

	if td == '>' {
		return b.arithmeticExpansion(t)
	}

	return b.operatorOrWord(t)
}

func (b *bashTokeniser) string(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	stops := singleStops

	if b.lastTokenDepth() == '"' {
		stops = doubleStops
	}

	for {
		switch t.ExceptRun(stops) {
		default:
			return t.ReturnError(io.ErrUnexpectedEOF)
		case '\n':
			return t.ReturnError(ErrInvalidCharacter)
		case '`':
			return t.Return(TokenString, b.backtick)
		case '$':
			return t.Return(TokenString, b.identifier)
		case '"', '\'':
			t.Next()
			b.popTokenDepth()

			return t.Return(TokenString, b.main)
		case '\\':
			t.Next()
			t.Next()
		}
	}
}

func (b *bashTokeniser) arithmeticExpansion(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	var early bool

	switch c := t.Peek(); c {
	case -1:
		return t.ReturnError(io.ErrUnexpectedEOF)
	case '"', '\'':
		return b.stringStart(t)
	case '$':
		return b.identifier(t)
	case '+', '-', '&', '|':
		early = true

		fallthrough
	case '<', '>':
		t.Next()

		if t.Peek() == c {
			t.Next()

			if early {
				break
			}
		}

		t.Accept("=")
	case '=', '!', '*', '/', '%', '^':
		t.Next()
		t.Accept("=")
	case '~', '?', ':', ',':
		t.Next()
	case ')':
		t.Next()

		if !t.Accept(")") {
			return t.ReturnError(ErrInvalidCharacter)
		}

		b.popTokenDepth()
	case '(':
		t.Next()

		if !t.Accept("(") {
			return t.ReturnError(ErrInvalidCharacter)
		}

		b.pushTokenDepth('>')
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
		return b.word(t)
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
		t.Accept(";&")
	case '"', '\'':
		return b.stringStart(t)
	case '(':
		t.Next()
		b.pushTokenDepth(')')
	case '{':
		t.Next()

		if !strings.ContainsRune(word, t.Peek()) {
			return b.braceExpansion(t)
		}

		b.pushTokenDepth('}')
	case '}', ')', ']':
		t.Next()

		if rune(b.lastTokenDepth()) != c {
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
	case '`':
		if b.lastTokenDepth() != '`' {
			return b.backtick(t)
		}

		b.popTokenDepth()
		t.Next()
	}

	return t.Return(TokenPunctuator, b.main)
}

func (b *bashTokeniser) startHeredoc(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if t.Peek() == -1 || t.Accept(newline) || t.Accept("#") {
		return t.ReturnError(io.ErrUnexpectedEOF)
	}

	if t.Accept(whitespace) || t.AcceptWord(escapedNewline, false) != "" {
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

	b.heredoc = unstring(tk.Data)

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
	for {
		if t.ExceptRun(newline) == -1 {
			return t.ReturnError(io.ErrUnexpectedEOF)
		}

		t.Next()

		state := t.State()

		if t.AcceptString(b.heredoc, false) == len(b.heredoc) && (t.Accept("\n") || t.Peek() == -1) {
			break
		}

		state.Reset()
	}

	b.heredoc = ""

	return t.Return(TokenHeredoc, b.main)
}

func (b *bashTokeniser) identifier(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	t.Next()

	if t.Accept("(") {
		if t.Accept("(") {
			b.pushTokenDepth('>')

			return t.Return(TokenPunctuator, b.main)
		}

		b.pushTokenDepth(')')

		return t.Return(TokenPunctuator, b.main)
	}

	if t.Accept("{") {
		b.pushTokenDepth('}')

		return t.Return(TokenPunctuator, b.word)
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

func (b *bashTokeniser) backtick(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	b.pushTokenDepth('`')
	t.Next()

	return t.Return(TokenPunctuator, b.main)
}

func (b *bashTokeniser) stringStart(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if rune(b.lastTokenDepth()) == t.Peek() {
		b.popTokenDepth()

		t.Next()

		return t.Return(TokenString, b.main)
	}

	b.pushTokenDepth(byte(t.Next()))

	return b.string(t)
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
		return b.word(t)
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

func (b *bashTokeniser) word(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if t.AcceptWord(keywords, false) != "" {
		return t.Return(TokenKeyword, b.main)
	}

	if t.Accept(identStart) {
		t.AcceptRun(identCont)

		state := t.State()

		if t.AcceptWord(assignment, false) != "" {
			state.Reset()

			return t.Return(TokenIdentifierAssign, b.main)
		} else if t.Peek() == rune(b.lastTokenDepth()) {
			return t.Return(TokenWord, b.main)
		} else if t.Peek() == '[' {
			return t.Return(TokenIdentifierAssign, b.startArrayAssign)
		}
	}

	var wb string

	switch b.lastTokenDepth() {
	case '}':
		wb = wordBreakNoBrace
	case ']':
		wb = wordBreakNoBracket
	default:
		wb = wordBreak
	}

	if t.Accept("\\") {
		t.Next()
	} else if t.Len() == 0 && t.Accept(wb) {
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
		case '\\':
			t.Next()
			t.Next()
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
				return t.ReturnError(ErrInvalidBraceExpansion)
			}

			if t.AcceptWord(dotdot, false) != "" {
				if !t.Accept(decimalDigit) {
					return t.ReturnError(ErrInvalidBraceExpansion)
				}

				t.AcceptRun(decimalDigit)
			}

			if !t.Accept("}") {
				return t.ReturnError(ErrInvalidBraceExpansion)
			}
		}
	} else if t.Accept(decimalDigit) {
		switch t.AcceptRun(decimalDigit) {
		default:
			return t.ReturnError(ErrInvalidBraceExpansion)
		case ',':
			return b.braceExpansionWord(t)
		case '.':
			if t.AcceptWord(dotdot, false) != "" {
				if !t.Accept(decimalDigit) {
					return t.ReturnError(ErrInvalidBraceExpansion)
				}

				t.AcceptRun(decimalDigit)

				if t.AcceptWord(dotdot, false) != "" {
					if !t.Accept(decimalDigit) {
						return t.ReturnError(ErrInvalidBraceExpansion)
					}

					t.AcceptRun(decimalDigit)
				}

				if !t.Accept("}") {
					return t.ReturnError(ErrInvalidBraceExpansion)
				}
			}
		}
	} else {
		switch t.ExceptRun(braceBreak) {
		case '\\':
			t.Next()
			t.Next()

			fallthrough
		case ',':
			return b.braceExpansionWord(t)
		default:
			return t.ReturnError(ErrInvalidBraceExpansion)
		}
	}

	return t.Return(TokenString, b.main)
}

func (b *bashTokeniser) braceExpansionWord(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	for {
		switch t.ExceptRun(braceWordBreak) {
		default:
			return t.ReturnError(ErrInvalidBraceExpansion)
		case '}':
			t.Next()

			return t.Return(TokenString, b.main)
		case '\\':
			t.Next()
			t.Next()
		case ',':
			t.Next()
		}
	}
}
