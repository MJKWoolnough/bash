// Package bash implements a bash tokeniser and AST.
package bash // import "vimagination.zapto.org/bash"

import "vimagination.zapto.org/parser"

// Parse parses Bash input into AST.
func Parse(t Tokeniser) (*File, error) {
	p, err := newBashParser(t)
	if err != nil {
		return nil, err
	}

	f := new(File)
	if err = f.parse(p); err != nil {
		return nil, err
	}

	return f, nil
}

// Parse parses Bash input into AST.
type File struct {
	Statements []Statement
	Tokens     Tokens
}

func (f *File) parse(b *bashParser) error {
	c := b.NewGoal()

	for {
		if tk := c.AcceptRunAllWhitespace(); tk == parser.TokenDone || tk == TokenCloseBacktick || tk == TokenCloseParen {
			break
		}

		b.AcceptRunAllWhitespace()

		c = b.NewGoal()

		var s Statement

		if err := s.parse(c, true); err != nil {
			return b.Error("File", err)
		}

		f.Statements = append(f.Statements, s)

		b.Score(c)

		c = b.NewGoal()
	}

	f.Tokens = b.ToTokens()

	return nil
}

type LogicalOperator uint8

const (
	LogicalOperatorNone LogicalOperator = iota
	LogicalOperatorAnd
	LogicalOperatorOr
)

type JobControl uint8

const (
	JobControlForeground JobControl = iota
	JobControlBackground
)

type Statement struct {
	Pipeline          Pipeline
	LogicalOperator   LogicalOperator
	LogicalExpression *Statement
	JobControl        JobControl
	Tokens
}

func (s *Statement) parse(b *bashParser, first bool) error {
	c := b.NewGoal()

	if err := s.Pipeline.parse(c); err != nil {
		return b.Error("Statement", err)
	}

	b.Score(c)

	c = b.NewGoal()

	c.AcceptRunWhitespace()

	if c.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "&&"}) {
		s.LogicalOperator = LogicalOperatorAnd
	} else if c.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "||"}) {
		s.LogicalOperator = LogicalOperatorOr
	}

	if s.LogicalOperator != LogicalOperatorNone {
		c.AcceptRunWhitespace()
		b.Score(c)

		c = b.NewGoal()
		s.LogicalExpression = new(Statement)

		if err := s.LogicalExpression.parse(c, false); err != nil {
			return b.Error("Statement", err)
		}

		b.Score(c)
	}

	if first {
		c = b.NewGoal()

		c.AcceptRunWhitespace()

		if c.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "&"}) {
			s.JobControl = JobControlBackground

			b.Score(c)
		} else if c.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ";"}) {
			b.Score(c)
		} else if tk := c.Peek().Type; tk != TokenLineTerminator && tk != parser.TokenDone && tk != TokenCloseBacktick && tk != TokenCloseParen {
			return c.Error("Statement", ErrInvalidEndOfStatement)
		}
	}

	s.Tokens = b.ToTokens()

	return nil
}

type PipelineTime uint8

const (
	PipelineTimeNone PipelineTime = iota
	PipelineTimeBash
	PipelineTimePosix
)

type Pipeline struct {
	PipelineTime
	Not      bool
	Command  Command
	Pipeline *Pipeline
	Tokens   Tokens
}

func (p *Pipeline) parse(b *bashParser) error {
	if b.AcceptToken(parser.Token{Type: TokenWord, Data: "time"}) {
		b.AcceptRunWhitespace()

		if b.AcceptToken(parser.Token{Type: TokenWord, Data: "-p"}) {
			p.PipelineTime = PipelineTimePosix
		} else {
			p.PipelineTime = PipelineTimeBash
		}

		b.AcceptRunWhitespace()
	}

	if p.Not = b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "!"}); p.Not {
		b.AcceptRunWhitespace()
	}

	c := b.NewGoal()

	if err := p.Command.parse(c, true); err != nil {
		return b.Error("Pipeline", err)
	}

	b.Score(c)

	c = b.NewGoal()

	c.AcceptRunWhitespace()

	if c.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "|"}) {
		c.AcceptRunWhitespace()
		b.Score(c)

		c = b.NewGoal()
		p.Pipeline = new(Pipeline)

		if err := p.Pipeline.parse(c); err != nil {
			return b.Error("Pipeline", err)
		}

		b.Score(c)
	}

	p.Tokens = b.ToTokens()

	return nil
}

type Command struct {
	Vars         []Assignment
	Redirections []Redirection
	Words        []Word
	Tokens       Tokens
}

func (c *Command) parse(b *bashParser, required bool) error {
	for {
		d := b.NewGoal()

		if b.Peek().Type == TokenIdentifierAssign {
			var a Assignment

			if err := a.parse(d); err != nil {
				return b.Error("Command", err)
			}

			c.Vars = append(c.Vars, a)
		} else if isRedirection(b) {
			var r Redirection

			if err := r.parse(d); err != nil {
				return b.Error("Command", err)
			}

			c.Redirections = append(c.Redirections, r)
		} else {
			break
		}

		b.Score(d)
		b.AcceptRunWhitespace()
	}

	d := b.NewGoal()

	for nextIsWordPart(d) {
		b.Score(d)
		d = b.NewGoal()

		if isRedirection(b) {
			var r Redirection

			if err := r.parse(d); err != nil {
				return b.Error("Command", err)
			}

			c.Redirections = append(c.Redirections, r)
		} else {
			var w Word

			if err := w.parse(d); err != nil {
				return b.Error("Command", err)
			}

			c.Words = append(c.Words, w)
		}

		b.Score(d)

		d = b.NewGoal()

		d.AcceptRunWhitespace()
	}

	if required && len(c.Words) == 0 {
		return b.Error("Command", ErrMissingWord)
	}

	c.Tokens = b.ToTokens()

	return nil
}

func isRedirection(b *bashParser) bool {
	c := b.NewGoal()

	if c.Accept(TokenNumberLiteral, TokenBraceWord) {
		if c.Accept(TokenPunctuator) {
			switch c.GetLastToken().Data {
			case "<", ">", ">|", ">>", "<<", "<<-", "<<<", "<&", ">&", "<>":
				return true
			}
		}
	} else if c.Accept(TokenPunctuator) {
		switch c.GetLastToken().Data {
		case "<", ">", ">|", ">>", "<<", "<<-", "<<<", "<&", ">&", "<>", "&>", "&>>":
			return true
		}
	}

	return false
}

type AssignmentType uint8

const (
	AssignmentAssign AssignmentType = iota
	AssignmentAppend
)

type Assignment struct {
	Identifier ParameterAssign
	Assignment AssignmentType
	Value      Value
	Tokens     Tokens
}

func (a *Assignment) parse(b *bashParser) error {
	c := b.NewGoal()

	if err := a.Identifier.parse(c); err != nil {
		return b.Error("Assignment", err)
	}

	b.Score(c)

	if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "="}) {
		a.Assignment = AssignmentAssign
	} else if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "+="}) {
		a.Assignment = AssignmentAppend
	} else {
		return b.Error("Assignment", ErrInvalidAssignment)
	}

	c = b.NewGoal()

	if err := a.Value.parse(c); err != nil {
		return b.Error("Assignment", err)
	}

	b.Score(c)

	a.Tokens = b.ToTokens()

	return nil
}

type ParameterAssign struct {
	Identifier *Token
	Subscript  *Word
	Tokens     Tokens
}

func (p *ParameterAssign) parse(b *bashParser) error {
	b.Next()

	p.Identifier = b.GetLastToken()

	if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "["}) {
		b.AcceptRunAllWhitespace()

		c := b.NewGoal()
		p.Subscript = new(Word)

		if err := p.Subscript.parse(c); err != nil {
			return b.Error("ParameterAssign", err)
		}

		b.Score(c)
		b.AcceptRunAllWhitespace()

		if !b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "]"}) {
			return b.Error("ParameterAssign", ErrMissingClosingBracket)
		}
	}

	p.Tokens = b.ToTokens()

	return nil
}

type Value struct {
	Word   *Word
	Array  []Word
	Tokens Tokens
}

func (v *Value) parse(b *bashParser) error {
	if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "("}) {
		b.AcceptRunAllWhitespace()

		v.Array = []Word{}

		for !b.AcceptToken(parser.Token{Type: TokenCloseParen, Data: ")"}) {
			c := b.NewGoal()

			var w Word

			if err := w.parse(c); err != nil {
				return b.Error("Value", err)
			}

			v.Array = append(v.Array, w)

			b.Score(c)
			b.AcceptRunAllWhitespace()
		}
	} else {
		c := b.NewGoal()
		v.Word = new(Word)

		if err := v.Word.parse(c); err != nil {
			return b.Error("Value", err)
		}

		b.Score(c)
	}

	v.Tokens = b.ToTokens()

	return nil
}

type Word struct {
	Parts  []WordPart
	Tokens Tokens
}

func (w *Word) parse(b *bashParser) error {
	for nextIsWordPart(b) {
		c := b.NewGoal()

		var wp WordPart

		if err := wp.parse(c); err != nil {
			return b.Error("Word", err)
		}

		w.Parts = append(w.Parts, wp)

		b.Score(c)
	}

	w.Tokens = b.ToTokens()

	return nil
}

func nextIsWordPart(b *bashParser) bool {
	switch tk := b.Peek(); tk.Type {
	case TokenWhitespace, TokenLineTerminator, TokenComment, TokenCloseBacktick, TokenCloseParen, parser.TokenDone:
		return false
	case TokenPunctuator:
		switch tk.Data {
		case "$((", "$(", "${":
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
	case tk == parser.Token{Type: TokenPunctuator, Data: "$("}, tk.Type == TokenOpenBacktick:
		w.CommandSubstitution = new(CommandSubstitution)

		if err := w.CommandSubstitution.parse(c); err != nil {
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

type ParameterType uint8

const (
	ParameterValue ParameterType = iota
	ParameterLength
	ParameterSubstitution
	ParameterAssignment
	ParameterMessage
	ParameterSetAssign
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
	Index          *Word
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
		} else if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ":"}) {
			p.Type = ParameterSubstring

			b.AcceptRunWhitespace()

			if !b.Accept(TokenNumberLiteral) {
				return b.Error("ParameterExpansion", ErrInvalidParameterExpansion)
			}

			p.SubstringStart = b.GetLastToken()

			if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ":"}) {
				b.AcceptRunWhitespace()

				if !b.Accept(TokenNumberLiteral) {
					return b.Error("ParameterExpansion", ErrInvalidParameterExpansion)
				}

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

			if err := p.Word.parse(c); err != nil {
				return b.Error("ParameterExpasion", err)
			}

			b.Score(c)
		} else if parsePattern || parseReplace {
			b.Accept(TokenPattern)

			p.Pattern = b.GetLastToken()

			if parseReplace && b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "/"}) {
				c := b.NewGoal()
				p.String = new(String)

				if err := p.String.parse(c); err != nil {
					return b.Error("ParameterExpasion", err)
				}

				b.Score(b)
			}
		}
	}

	if !b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "}"}) {
		return b.Error("ParameterExpasion", ErrMissingClosingBrace)
	}

	p.Tokens = b.ToTokens()

	return nil
}

type Parameter struct {
	Parameter *Token
	Array     *Word
	Tokens    Tokens
}

func (p *Parameter) parse(b *bashParser) error {
	if b.Accept(TokenIdentifier) {
		p.Parameter = b.GetLastToken()

		if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "["}) {
			c := b.NewGoal()
			p.Array = new(Word)

			if err := p.Array.parse(c); err != nil {
				return b.Error("Parameter", err)
			}

			b.Score(c)

			if !b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "]"}) {
				return b.Error("Parameter", ErrMissingClosingBracket)
			}
		}
	} else {
		b.Next()

		p.Parameter = b.GetLastToken()
	}

	p.Tokens = b.ToTokens()

	return nil
}

type String struct {
	WordsOrTokens []WordOrToken
	Tokens        Tokens
}

func (s *String) parse(b *bashParser) error {
	for b.Peek().Type != parser.TokenDone && !b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "}"}) {
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

type WordOrToken struct {
	Token  *Token
	Word   *Word
	Tokens Tokens
}

func (w *WordOrToken) parse(b *bashParser) error {
	if nextIsWordPart(b) {
		c := b.NewGoal()
		w.Word = new(Word)

		if err := w.Word.parse(c); err != nil {
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

type SubstitutionType uint8

const (
	SubstitutionNew SubstitutionType = iota
	SubstitutionBacktick
)

type CommandSubstitution struct {
	SubstitutionType SubstitutionType
	Command          File
	Tokens           Tokens
}

func (cs *CommandSubstitution) parse(b *bashParser) error {
	end := parser.Token{Type: TokenCloseParen, Data: ")"}

	if tk := b.Next(); tk.Type == TokenOpenBacktick {
		cs.SubstitutionType = SubstitutionBacktick
		end = parser.Token{Type: TokenCloseBacktick, Data: tk.Data}
	}

	b.AcceptRunAllWhitespace()

	c := b.NewGoal()

	if err := cs.Command.parse(c); err != nil {
		return b.Error("CommandSubstitution", err)
	}

	b.Score(c)
	b.AcceptRunAllWhitespace()
	b.AcceptToken(end)

	cs.Tokens = b.ToTokens()

	return nil
}

type Redirection struct {
	Input      *Token
	Redirector *Token
	Output     Word
	Tokens     Tokens
}

func (r *Redirection) parse(b *bashParser) error {
	if b.Accept(TokenNumberLiteral, TokenBraceWord) {
		r.Input = b.GetLastToken()
	}

	b.Accept(TokenPunctuator)

	r.Redirector = b.GetLastToken()

	b.AcceptRunWhitespace()

	c := b.NewGoal()

	if err := r.Output.parse(c); err != nil {
		return b.Error("Redirection", err)
	}

	b.Score(c)

	r.Tokens = b.ToTokens()

	return nil
}

func (r *Redirection) isHeredoc() bool {
	return r.Redirector != nil && (r.Redirector.Data == "<<" || r.Redirector.Data == "<<-")
}

type ArithmeticExpansion struct {
	WordsAndOperators []WordOrOperator
	Tokens            Tokens
}

func (a *ArithmeticExpansion) parse(b *bashParser) error {
	b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "$(("})
	b.AcceptRunAllWhitespace()

	for !b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "))"}) {
		c := b.NewGoal()

		var w WordOrOperator

		if err := w.parse(c); err != nil {
			return b.Error("ArithmeticExpansion", err)
		}

		a.WordsAndOperators = append(a.WordsAndOperators, w)

		b.Score(c)
		b.AcceptRunAllWhitespace()

	}

	a.Tokens = b.ToTokens()

	return nil
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
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ","}) {
		w.Operator = b.GetLastToken()
	} else {
		c := b.NewGoal()
		w.Word = new(Word)

		if err := w.Word.parse(c); err != nil {
			return b.Error("WordOrOperator", err)
		}

		b.Score(c)
	}

	w.Tokens = b.ToTokens()

	return nil
}
