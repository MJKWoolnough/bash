package bash

import (
	"errors"
	"io"
	"slices"
	"strings"

	"vimagination.zapto.org/parser"
)

var keywords = []string{"if", "then", "else", "elif", "fi", "case", "esac", "while", "for", "in", "do", "done", "time", "until", "coproc", "select", "function", "{", "}", "[[", "]]", "!"}

const (
	whitespace   = " \t"
	newline      = "\n"
	doubleStops  = "\\\n`$\""
	singleStops  = "\n'"
	word         = "\\\"'`(){}- \t\n"
	wordBreak    = " `\\\t\n|&;<>()={}"
	hexDigit     = "0123456789ABCDEFabcdef"
	octalDigit   = "012345678"
	decimalDigit = "0123456789"
	numberChars  = decimalDigit + "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz@_"
)

const (
	TokenWhitespace parser.TokenType = iota
	TokenLineTerminator
	TokenComment
	TokenIdentifier
	TokenKeyword
	TokenWord
	TokenNumberLiteral
	TokenString
	TokenPunctuator
)

type bashTokeniser struct {
	tokenDepth []byte
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

	if t.Accept(whitespace) {
		t.AcceptRun(whitespace)

		return t.Return(TokenWhitespace, b.main)
	}

	if t.Accept(newline) {
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
		t.Accept("<&>")
	case '>':
		t.Next()
		t.Accept(">&|")
	case '|':
		t.Next()
		t.Accept("&|")
	case '&':
		t.Next()
		t.Accept("&")
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

		if strings.ContainsRune(word, t.Peek()) {
			b.pushTokenDepth('}')

			return t.Return(TokenPunctuator, b.main)
		}

		return b.braceExpansion(t)
	case '}', ')':
		if rune(b.lastTokenDepth()) != c {
			return t.ReturnError(ErrInvalidCharacter)
		}

		b.popTokenDepth()
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

	t.ExceptRun(word)

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
		t.AcceptRun(t.AcceptRun(octalDigit))
	}

	return t.Return(TokenNumberLiteral, b.main)
}

func (b *bashTokeniser) number(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if !t.Accept(decimalDigit) {
		return word(t)
	}

	t.AcceptRun(decimalDigit)

	if t.Accept("#") {
		if !t.Accept(numberChars) {
			return t.ReturnError(ErrInvalidNumber)
		}

		t.AcceptRun(numberChars)
	}

	return tk.Return(TokenNumberLiteral, b.main)
}

func (b *bashTokeniser) word(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	var hasEscape bool

	for {
		switch t.ExceptRun(wordBreak) {
		case -1:
			if t.Len() == 0 {
				if b.lastTokenDepth() == 0 {
					return t.Done()
				}

				return t.ReturnError(io.ErrUnexpectedEOF)
			}

			fallthrough
		default:
			data := t.Get()

			if slices.Contains(keywords, data) {
				return parser.Token{Type: TokenKeyword, Data: data}, b.main
			}

			if !hasEscape && t.Peek() == '=' {
				return parser.Token{Type: TokenIdentifier, Data: data}, b.main
			}

			return parser.Token{Type: TokenWord, Data: data}, b.main
		case '\\':
			t.Next()
			t.Next()

			hasEscape = true
		}
	}
}

func (b *bashTokeniser) braceExpansion(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
}

var (
	ErrInvalidCharacter = errors.New("invalid character")
	ErrInvalidNumber    = errors.New("invalid number")
)
