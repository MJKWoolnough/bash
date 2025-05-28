package bash

import (
	"fmt"

	"vimagination.zapto.org/parser"
)

// Token represents a parser.Token combined with positioning information.
type Token struct {
	parser.Token
	Pos, Line, LinePos uint64
}

// Tokens represents a list ok tokens that have been parsed.
type Tokens []Token

// Comments is a collection of Comment Tokens.
type Comments []Token

type bashParser struct {
	Tokens
}

// Tokeniser represents the methods required by the bash tokeniser.
type Tokeniser interface {
	Iter(func(parser.Token) bool)
	TokeniserState(parser.TokenFunc)
	GetError() error
}

func newBashParser(t Tokeniser) (*bashParser, error) {
	b := new(bashTokeniser)

	t.TokeniserState(b.main)

	var (
		tokens             Tokens
		err                error
		pos, line, linePos uint64
	)

	for tk := range t.Iter {
		tokens = append(tokens, Token{Token: tk, Pos: pos, Line: line, LinePos: linePos})

		switch tk.Type {
		case parser.TokenDone:
		case parser.TokenError:
			err = Error{Err: t.GetError(), Parsing: "Tokens", Token: tokens[len(tokens)-1]}
		case TokenLineTerminator:
			line += uint64(len(tk.Data))
			linePos = 0
		default:
			for _, c := range tk.Data {
				if c == '\n' {
					line++
					linePos = 0
				} else {
					linePos++
				}
			}
		}

		pos += uint64(len(tk.Data))
	}

	return &bashParser{Tokens: tokens[0:0:len(tokens)]}, err
}

func (b bashParser) NewGoal() *bashParser {
	return &bashParser{
		Tokens: b.Tokens[len(b.Tokens):],
	}
}

func (b *bashParser) Score(k *bashParser) {
	b.Tokens = b.Tokens[:len(b.Tokens)+len(k.Tokens)]
}

func (b *bashParser) Next() Token {
	l := len(b.Tokens)
	b.Tokens = b.Tokens[:l+1]
	tk := b.Tokens[l]

	return tk
}

func (b *bashParser) backup() {
	b.Tokens = b.Tokens[:len(b.Tokens)-1]
}

func (b *bashParser) Peek() parser.Token {
	tk := b.Next().Token

	b.backup()

	return tk
}

func (b *bashParser) Accept(ts ...parser.TokenType) bool {
	tt := b.Next().Type

	for _, pt := range ts {
		if pt == tt {
			return true
		}
	}

	b.backup()

	return false
}

func (b *bashParser) AcceptRun(ts ...parser.TokenType) parser.TokenType {
Loop:
	for {
		tt := b.Next().Type

		for _, pt := range ts {
			if pt == tt {
				continue Loop
			}
		}

		b.backup()

		return tt
	}
}

func (b *bashParser) AcceptToken(tk parser.Token) bool {
	if b.Next().Token == tk {
		return true
	}

	b.backup()

	return false
}

func (b *bashParser) ToTokens() Tokens {
	return b.Tokens[:len(b.Tokens):len(b.Tokens)]
}

func (b *bashParser) GetLastToken() *Token {
	return &b.Tokens[len(b.Tokens)-1]
}

func (b *bashParser) AcceptRunWhitespace() parser.TokenType {
	return b.AcceptRun(TokenWhitespace)
}

func (b *bashParser) AcceptRunAllWhitespaceComments() Comments {
	var c Comments

	s := b.NewGoal()

	for s.AcceptRunAllWhitespaceNoComments() == TokenComment {
		c = append(c, s.Next())

		b.Score(s)

		s = b.NewGoal()
	}

	return c
}

func (b *bashParser) AcceptRunWhitespaceComments() Comments {
	var c Comments

	s := b.NewGoal()

	for s.AcceptRunWhitespace() == TokenComment {
		b.Score(s)

		c = append(c, b.Next())
		s = b.NewGoal()

		if s.Accept(TokenLineTerminator) {
			if len(s.GetLastToken().Data) > 1 {
				break
			}
		}
	}

	return c
}

func (b *bashParser) AcceptRunAllWhitespace() parser.TokenType {
	return b.AcceptRun(TokenWhitespace, TokenComment, TokenLineTerminator)
}

func (b *bashParser) AcceptRunAllWhitespaceNoComments() parser.TokenType {
	return b.AcceptRun(TokenWhitespace, TokenLineTerminator)
}

// Error represents a Bash parsing error.
type Error struct {
	Err     error
	Parsing string
	Token   Token
}

// Error implements the error interface.
func (e Error) Error() string {
	return fmt.Sprintf("%s: error at position %d (%d:%d):\n%s", e.Parsing, e.Token.Pos+1, e.Token.Line+1, e.Token.LinePos+1, e.Err)
}

// Unwrap returns the underlying error.
func (e Error) Unwrap() error {
	return e.Err
}

func (b *bashParser) Error(parsingFunc string, err error) error {
	tk := b.Next()

	b.backup()

	return Error{
		Err:     err,
		Parsing: parsingFunc,
		Token:   tk,
	}
}
