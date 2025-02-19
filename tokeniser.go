package bash

import (
	"io"

	"vimagination.zapto.org/parser"
)

const (
	whitespace = " \t"
	newline    = "\n"
)

const (
	TokenWhitespace parser.TokenType = iota
	TokenLineTerminator
	TokenComment
	TokenIdentifier
	TokenKeyword
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
}

func (b *bashTokeniser) arithmeticExpansion(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
}

func (b *bashTokeniser) operatorOrWord(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
}
