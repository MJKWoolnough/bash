package bash

import (
	"vimagination.zapto.org/parser"
)

type AssignmentOrWord struct {
	Assignment *Assignment
	Word       *Word
	Tokens     Tokens
}

func (a *AssignmentOrWord) parse(b *bashParser) error {
	var err error

	c := b.NewGoal()

	if tk := b.Peek().Type; tk == TokenIdentifierAssign || tk == TokenLetIdentifierAssign {
		a.Assignment = new(Assignment)
		err = a.Assignment.parse(c)
	} else {
		a.Word = new(Word)
		err = a.Word.parse(c, false)
	}

	if err != nil {
		return b.Error("AssignmentOrWord", err)
	}

	b.Score(c)

	a.Tokens = b.ToTokens()

	return nil
}

func (a *AssignmentOrWord) isMultiline(v bool) bool {
	return a.Assignment != nil && a.Assignment.isMultiline(v) || a.Word != nil && a.Word.isMultiline(v)
}

type Value struct {
	Word     *Word
	Array    []ArrayWord
	Comments [2]Comments
	Tokens   Tokens
}

func (v *Value) parse(b *bashParser) error {
	if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "("}) {
		v.Comments[0] = b.AcceptRunWhitespaceComments()

		v.Array = []ArrayWord{}
		c := b.NewGoal()

		c.AcceptRunAllWhitespace()

		for !c.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ")"}) {
			b.AcceptRunAllWhitespaceNoComments()

			c = b.NewGoal()

			var a ArrayWord

			if err := a.parse(c); err != nil {
				return b.Error("Value", err)
			}

			v.Array = append(v.Array, a)

			b.Score(c)

			c = b.NewGoal()

			c.AcceptRunAllWhitespace()
		}

		v.Comments[1] = b.AcceptRunAllWhitespaceComments()

		b.AcceptRunAllWhitespaceNoComments()
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ")"})
	} else {
		c := b.NewGoal()
		v.Word = new(Word)

		if err := v.Word.parse(c, false); err != nil {
			return b.Error("Value", err)
		}

		b.Score(c)
	}

	v.Tokens = b.ToTokens()

	return nil
}

func (v *Value) isMultiline(vs bool) bool {
	if len(v.Comments[0]) > 0 || len(v.Comments[1]) > 0 {
		return true
	}

	if v.Word != nil {
		return v.Word.isMultiline(vs)
	}

	for _, ar := range v.Array {
		if ar.isMultiline(vs) {
			return true
		}
	}

	return false
}

type ArrayWord struct {
	Word     Word
	Comments [2]Comments
	Tokens   Tokens
}

func (a *ArrayWord) parse(b *bashParser) error {
	a.Comments[0] = b.AcceptRunAllWhitespaceComments()

	b.AcceptRunAllWhitespaceNoComments()

	c := b.NewGoal()

	if err := a.Word.parse(c, false); err != nil {
		return b.Error("ArrayWord", err)
	}

	if len(a.Word.Parts) == 0 {
		return b.Error("ArrayWord", ErrMissingWord)
	}

	b.Score(c)

	a.Comments[1] = b.AcceptRunWhitespaceComments()
	a.Tokens = b.ToTokens()

	return nil
}

func (a *ArrayWord) isMultiline(v bool) bool {
	if len(a.Comments[0]) > 0 || len(a.Comments[1]) > 0 {
		return true
	}

	return a.Word.isMultiline(v)
}

type Word struct {
	Parts  []WordPart
	Tokens Tokens
}

func (w *Word) parse(b *bashParser, splitAssign bool) error {
	for nextIsWordPart(b) {
		nextIsAssign := b.Peek().Type == TokenIdentifierAssign
		c := b.NewGoal()

		var wp WordPart

		if err := wp.parse(c); err != nil {
			return b.Error("Word", err)
		}

		w.Parts = append(w.Parts, wp)

		b.Score(c)

		if nextIsAssign && splitAssign {
			break
		}
	}

	w.Tokens = b.ToTokens()

	return nil
}

func (w *Word) isMultiline(v bool) bool {
	for _, p := range w.Parts {
		if p.isMultiline(v) {
			return true
		}
	}

	return false
}

func nextIsWordPart(b *bashParser) bool {
	switch tk := b.Peek(); tk.Type {
	case TokenWhitespace, TokenLineTerminator, TokenComment, TokenCloseBacktick, TokenHeredoc, TokenBinaryOperator, TokenHeredocEnd, parser.TokenDone:
		return false
	case TokenBraceExpansion:
		return tk.Data != "}"
	case TokenPunctuator:
		switch tk.Data {
		case "$((", "$(", "${", "<(", ">(":
			return true
		}

		return false
	}

	return true
}

type WordPart struct {
	Part                *Token
	ParameterExpansion  *ParameterExpansion
	CommandSubstitution *CommandSubstitution
	ArithmeticExpansion *ArithmeticExpansion
	BraceExpansion      *BraceExpansion
	Tokens              Tokens
}

func (w *WordPart) parse(b *bashParser) error {
	c := b.NewGoal()

	switch tk := b.Peek(); {
	case tk == parser.Token{Type: TokenPunctuator, Data: "${"}:
		w.ParameterExpansion = new(ParameterExpansion)

		if err := w.ParameterExpansion.parse(c); err != nil {
			return b.Error("WordPart", err)
		}
	case tk == parser.Token{Type: TokenPunctuator, Data: "$(("}:
		w.ArithmeticExpansion = new(ArithmeticExpansion)

		if err := w.ArithmeticExpansion.parse(c); err != nil {
			return b.Error("WordPart", err)
		}
	case tk == parser.Token{Type: TokenPunctuator, Data: "$("}, tk.Type == TokenOpenBacktick, tk == parser.Token{Type: TokenPunctuator, Data: "<("}, tk == parser.Token{Type: TokenPunctuator, Data: ">("}:
		w.CommandSubstitution = new(CommandSubstitution)

		if err := w.CommandSubstitution.parse(c); err != nil {
			return b.Error("WordPart", err)
		}
	case tk == parser.Token{Type: TokenBraceExpansion, Data: "{"}:
		w.BraceExpansion = new(BraceExpansion)

		if err := w.BraceExpansion.parse(c); err != nil {
			return b.Error("WordPart", err)
		}
	default:
		b.Next()

		w.Part = b.GetLastToken()
	}

	b.Score(c)

	w.Tokens = b.ToTokens()

	return nil
}

func (w *WordPart) isMultiline(v bool) bool {
	if w.ParameterExpansion != nil {
		return w.ParameterExpansion.isMultiline(v)
	} else if w.ArithmeticExpansion != nil {
		return w.ArithmeticExpansion.isMultiline(v)
	} else if w.CommandSubstitution != nil {
		return w.CommandSubstitution.isMultiline(v)
	}

	return false
}

type BraceExpansion struct {
	Words  []Word
	Tokens Tokens
}

func (be *BraceExpansion) parse(b *bashParser) error {
	b.Next()

	for !b.AcceptToken(parser.Token{Type: TokenBraceExpansion, Data: "}"}) {
		c := b.NewGoal()

		var w Word

		if err := w.parse(c, false); err != nil {
			return b.Error("BraceExpansion", err)
		}

		be.Words = append(be.Words, w)

		b.Score(c)
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ","})
	}

	be.Tokens = b.ToTokens()

	return nil
}

type ParameterType uint8

const (
	ParameterValue ParameterType = iota
	ParameterLength
	ParameterSubstitution
	ParameterAssignment
	ParameterMessage
	ParameterSetAssign
	ParameterUnsetSubstitution
	ParameterUnsetAssignment
	ParameterUnsetMessage
	ParameterUnsetSetAssign
	ParameterSubstring
	ParameterPrefix
	ParameterPrefixSeperate
	ParameterRemoveStartShortest
	ParameterRemoveStartLongest
	ParameterRemoveEndShortest
	ParameterRemoveEndLongest
	ParameterReplace
	ParameterReplaceAll
	ParameterReplaceStart
	ParameterReplaceEnd
	ParameterLowercaseFirstMatch
	ParameterLowercaseAllMatches
	ParameterUppercaseFirstMatch
	ParameterUppercaseAllMatches
	ParameterUppercase
	ParameterUppercaseFirst
	ParameterLowercase
	ParameterQuoted
	ParameterEscaped
	ParameterPrompt
	ParameterDeclare
	ParameterQuotedArrays
	ParameterQuotedArraysSeperate
	ParameterAttributes
)

type ParameterExpansion struct {
	Indirect       bool
	Parameter      Parameter
	Type           ParameterType
	SubstringStart *Token
	SubstringEnd   *Token
	Word           *Word
	Pattern        *Token
	String         *String
	Tokens         Tokens
}

func (p *ParameterExpansion) parse(b *bashParser) error {
	b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "${"})

	if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "#"}) {
		p.Type = ParameterLength
	} else {
		p.Indirect = b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "!"})
	}

	c := b.NewGoal()

	if err := p.Parameter.parse(c); err != nil {
		return b.Error("ParameterExpansion", err)
	}

	b.Score(c)

	if p.Type != ParameterLength {
		var parseWord, parseReplace, parsePattern bool

		if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ":="}) {
			p.Type = ParameterSubstitution
			parseWord = true
		} else if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ":?"}) {
			p.Type = ParameterAssignment
			parseWord = true
		} else if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ":+"}) {
			p.Type = ParameterMessage
			parseWord = true
		} else if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ":-"}) {
			p.Type = ParameterSetAssign
			parseWord = true
		} else if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "="}) {
			p.Type = ParameterUnsetSubstitution
			parseWord = true
		} else if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "?"}) {
			p.Type = ParameterUnsetAssignment
			parseWord = true
		} else if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "+"}) {
			p.Type = ParameterUnsetMessage
			parseWord = true
		} else if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "-"}) {
			p.Type = ParameterUnsetSetAssign
			parseWord = true
		} else if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ":"}) {
			p.Type = ParameterSubstring

			b.AcceptRunWhitespace()
			b.Accept(TokenNumberLiteral)

			p.SubstringStart = b.GetLastToken()

			if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ":"}) {
				b.AcceptRunWhitespace()
				b.Accept(TokenNumberLiteral)

				p.SubstringEnd = b.GetLastToken()
			}
		} else if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "#"}) {
			p.Type = ParameterRemoveStartShortest
			parseWord = true
		} else if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "##"}) {
			p.Type = ParameterRemoveStartLongest
			parseWord = true
		} else if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "%"}) {
			p.Type = ParameterRemoveEndShortest
			parseWord = true
		} else if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "%%"}) {
			p.Type = ParameterRemoveEndLongest
			parseWord = true
		} else if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "/"}) {
			p.Type = ParameterReplace
			parseReplace = true
		} else if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "//"}) {
			p.Type = ParameterReplaceAll
			parseReplace = true
		} else if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "/#"}) {
			p.Type = ParameterReplaceStart
			parseReplace = true
		} else if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "/%"}) {
			p.Type = ParameterReplaceEnd
			parseReplace = true
		} else if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "^"}) {
			p.Type = ParameterUppercaseFirstMatch
			parsePattern = true
		} else if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "^^"}) {
			p.Type = ParameterUppercaseAllMatches
			parsePattern = true
		} else if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ","}) {
			p.Type = ParameterLowercaseFirstMatch
			parsePattern = true
		} else if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ",,"}) {
			p.Type = ParameterLowercaseAllMatches
			parsePattern = true
		} else if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "@"}) {
			if p.Indirect && b.Peek() == (parser.Token{Type: TokenPunctuator, Data: "}"}) {
				p.Indirect = false
				p.Type = ParameterPrefixSeperate
			} else {
				b.Accept(TokenBraceWord)

				switch b.GetLastToken().Data {
				case "U":
					p.Type = ParameterUppercase
				case "u":
					p.Type = ParameterUppercaseFirst
				case "L":
					p.Type = ParameterLowercase
				case "Q":
					p.Type = ParameterQuoted
				case "E":
					p.Type = ParameterEscaped
				case "P":
					p.Type = ParameterPrompt
				case "A":
					p.Type = ParameterDeclare
				case "K":
					p.Type = ParameterQuotedArrays
				case "a":
					p.Type = ParameterAttributes
				case "k":
					p.Type = ParameterQuotedArraysSeperate
				}
			}
		} else if p.Indirect && b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "*"}) {
			p.Indirect = false
			p.Type = ParameterPrefix
		}

		if parseWord {
			c := b.NewGoal()
			p.Word = new(Word)

			if err := p.Word.parse(c, false); err != nil {
				return b.Error("ParameterExpansion", err)
			}

			b.Score(c)
		} else if parsePattern || parseReplace {
			b.Accept(TokenPattern)

			p.Pattern = b.GetLastToken()

			if parseReplace && b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "/"}) {
				c := b.NewGoal()
				p.String = new(String)

				if err := p.String.parse(c); err != nil {
					return b.Error("ParameterExpansion", err)
				}

				b.Score(c)
			}
		}
	}

	if !b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "}"}) {
		return b.Error("ParameterExpansion", ErrMissingClosingBrace)
	}

	p.Tokens = b.ToTokens()

	return nil
}

func (p *ParameterExpansion) isMultiline(v bool) bool {
	if p.Word != nil && p.Word.isMultiline(v) {
		return true
	} else if p.String != nil && p.String.isMultiline(v) {
		return true
	}

	return p.Parameter.isMultiline(v)
}

type Parameter struct {
	Parameter *Token
	Array     []WordOrOperator
	Tokens    Tokens
}

func (p *Parameter) parse(b *bashParser) error {
	if b.Accept(TokenIdentifier) {
		p.Parameter = b.GetLastToken()

		if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "["}) {
			b.AcceptRunWhitespace()

			for !b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "]"}) {
				c := b.NewGoal()

				var w WordOrOperator

				if err := w.parse(c); err != nil {
					return b.Error("Parameter", err)
				}

				p.Array = append(p.Array, w)

				b.Score(c)
				b.AcceptRunAllWhitespace()
			}
		}
	} else {
		b.Next()

		p.Parameter = b.GetLastToken()
	}

	p.Tokens = b.ToTokens()

	return nil
}

func (p *Parameter) isMultiline(v bool) bool {
	for _, wo := range p.Array {
		if wo.isMultiline(v) {
			return true
		}
	}

	return false
}

type String struct {
	WordsOrTokens []WordOrToken
	Tokens        Tokens
}

func (s *String) parse(b *bashParser) error {
	for b.Peek().Type != parser.TokenDone && b.Peek() != (parser.Token{Type: TokenPunctuator, Data: "}"}) {
		c := b.NewGoal()

		var wp WordOrToken

		if err := wp.parse(c); err != nil {
			return b.Error("String", err)
		}

		b.Score(c)

		s.WordsOrTokens = append(s.WordsOrTokens, wp)
	}

	s.Tokens = b.ToTokens()

	return nil
}

func (s *String) isMultiline(v bool) bool {
	for _, wt := range s.WordsOrTokens {
		if wt.isMultiline(v) {
			return true
		}
	}

	return false
}

type WordOrToken struct {
	Token  *Token
	Word   *Word
	Tokens Tokens
}

func (w *WordOrToken) parse(b *bashParser) error {
	if nextIsWordPart(b) {
		c := b.NewGoal()
		w.Word = new(Word)

		if err := w.Word.parse(c, false); err != nil {
			return b.Error("WordOrToken", err)
		}

		b.Score(c)
	} else {
		b.Next()

		w.Token = b.GetLastToken()
	}

	w.Tokens = b.ToTokens()

	return nil
}

func (w *WordOrToken) isMultiline(v bool) bool {
	return w.Word != nil && w.Word.isMultiline(v)
}

type WordOrOperator struct {
	Word     *Word
	Operator *Token
	Tokens   Tokens
}

func (w *WordOrOperator) parse(b *bashParser) error {
	if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "++"}) ||
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "--"}) ||
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "-"}) ||
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "+"}) ||
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "!"}) ||
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "~"}) ||
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "**"}) ||
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "*"}) ||
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "/"}) ||
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "%"}) ||
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "<<"}) ||
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ">>"}) ||
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "<="}) ||
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "=>"}) ||
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "<"}) ||
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ">"}) ||
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "?"}) ||
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "&"}) ||
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "^"}) ||
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "|"}) ||
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "&&"}) ||
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "||"}) ||
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "?"}) ||
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ":"}) ||
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "="}) ||
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "*="}) ||
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "/="}) ||
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "%="}) ||
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "+="}) ||
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "-="}) ||
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "<<="}) ||
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ">>="}) ||
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "&="}) ||
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "^="}) ||
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "!="}) ||
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "("}) ||
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ")"}) ||
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ";"}) ||
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ","}) {
		w.Operator = b.GetLastToken()
	} else {
		c := b.NewGoal()
		w.Word = new(Word)

		if err := w.Word.parse(c, true); err != nil {
			return b.Error("WordOrOperator", err)
		}

		b.Score(c)
	}

	w.Tokens = b.ToTokens()

	return nil
}

func (w *WordOrOperator) isMultiline(v bool) bool {
	return w.Word != nil && w.Word.isMultiline(v)
}

func (w *WordOrOperator) operatorIsToken(tk parser.Token) bool {
	if w.Operator != nil {
		return w.Operator.Token == tk
	}

	return false
}
