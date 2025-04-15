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
	Lines  []Line
	Tokens Tokens
}

func (f *File) parse(b *bashParser) error {
	c := b.NewGoal()

	for {
		if tk := c.AcceptRunAllWhitespace(); tk == parser.TokenDone || tk == TokenCloseBacktick || tk == TokenCloseParen {
			break
		}

		b.AcceptRunAllWhitespace()

		c = b.NewGoal()

		var l Line

		if err := l.parse(c); err != nil {
			return b.Error("File", err)
		}

		f.Lines = append(f.Lines, l)

		b.Score(c)

		c = b.NewGoal()
	}

	b.AcceptRunWhitespace()

	f.Tokens = b.ToTokens()

	return nil
}

type Line struct {
	Statements []Statement
	Tokens     Tokens
}

func (l *Line) parse(b *bashParser) error {
	c := b.NewGoal()

	for !c.Accept(TokenComment, TokenLineTerminator, TokenCloseBacktick, TokenCloseParen) && c.Peek().Type != parser.TokenDone {
		c = b.NewGoal()

		var s Statement

		if err := s.parse(c, true); err != nil {
			return b.Error("Line", err)
		}

		l.Statements = append(l.Statements, s)

		c.AcceptRunWhitespace()
		b.Score(c)

		c = b.NewGoal()
	}

	for n := range l.Statements {
		c := b.NewGoal()

		c.Accept(TokenLineTerminator)

		d := c.NewGoal()

		if err := l.Statements[n].parseHeredocs(d); err != nil {
			return c.Error("Line", err)
		}

		if len(d.Tokens) == 0 {
			continue
		}

		c.Score(d)
		b.Score(c)
	}

	l.Tokens = b.ToTokens()

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
	Pipeline        Pipeline
	LogicalOperator LogicalOperator
	Statement       *Statement
	JobControl      JobControl
	Tokens
}

func (s *Statement) parse(b *bashParser, first bool) error {
	c := b.NewGoal()

	if err := s.Pipeline.parse(c, !first); err != nil {
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
		s.Statement = new(Statement)

		if err := s.Statement.parse(c, false); err != nil {
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
		}
	}

	s.Tokens = b.ToTokens()

	return nil
}

func (s *Statement) parseHeredocs(b *bashParser) error {
	c := b.NewGoal()

	if err := s.Pipeline.parseHeredocs(c); err != nil {
		return b.Error("Statement", err)
	}

	b.Score(c)

	if s.Statement != nil {
		c = b.NewGoal()

		if err := s.Statement.parseHeredocs(c); err != nil {
			return b.Error("Statement", err)
		}

		b.Score(c)
	}

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
	Not               bool
	CommandOrCompound CommandOrCompound
	Pipeline          *Pipeline
	Tokens            Tokens
}

func (p *Pipeline) parse(b *bashParser, required bool) error {
	if b.AcceptToken(parser.Token{Type: TokenKeyword, Data: "time"}) {
		b.AcceptRunWhitespace()

		if b.AcceptToken(parser.Token{Type: TokenWord, Data: "-p"}) {
			p.PipelineTime = PipelineTimePosix
		} else {
			p.PipelineTime = PipelineTimeBash
		}

		b.AcceptRunWhitespace()
	}

	if p.Not = b.AcceptToken(parser.Token{Type: TokenKeyword, Data: "!"}); p.Not {
		b.AcceptRunWhitespace()
	}

	c := b.NewGoal()

	if err := p.CommandOrCompound.parse(c, required); err != nil {
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

		if err := p.Pipeline.parse(c, true); err != nil {
			return b.Error("Pipeline", err)
		}

		b.Score(c)
	}

	p.Tokens = b.ToTokens()

	return nil
}

func (p *Pipeline) parseHeredocs(b *bashParser) error {
	c := b.NewGoal()

	if err := p.CommandOrCompound.parseHeredoc(c); err != nil {
		return b.Error("Pipeline", err)
	}

	b.Score(c)

	if p.Pipeline != nil {
		c = b.NewGoal()

		if err := p.Pipeline.parseHeredocs(c); err != nil {
			return b.Error("Pipeline", err)
		}

		b.Score(c)
	}

	return nil
}

type CommandOrCompound struct {
	Command  *Command
	Compound *Compound
	Tokens   Tokens
}

func (cc *CommandOrCompound) parse(b *bashParser, required bool) error {
	var err error

	c := b.NewGoal()

	if isCompoundNext(b) {
		cc.Compound = new(Compound)
		err = cc.Compound.parse(c)
	} else {
		cc.Command = new(Command)
		err = cc.Command.parse(c, required)
	}

	if err != nil {
		return b.Error("CommandOrCompound", err)
	}

	b.Score(c)

	cc.Tokens = b.ToTokens()

	return nil
}

func (cc *CommandOrCompound) parseHeredoc(b *bashParser) error {
	var err error

	c := b.NewGoal()

	if cc.Command != nil {
		err = cc.Command.parseHeredocs(c)
	} else {
		err = cc.Compound.parseHeredocs(c)
	}

	if err != nil {
		return b.Error("CommandOrCompound", err)
	}

	b.Score(c)

	return nil
}

func isCompoundNext(b *bashParser) bool {
	tk := b.Peek()

	return tk.Type == TokenKeyword && (tk.Data == "if" || tk.Data == "case" || tk.Data == "while" || tk.Data == "for" || tk.Data == "until" || tk.Data == "select" || tk.Data == "[[") || tk.Type == TokenPunctuator && (tk.Data == "((" || tk.Data == "(" || tk.Data == "{")
}

type Compound struct {
	IfCompound         *IfCompound
	CaseCompound       *CaseCompound
	LoopCompound       *LoopCompound
	ForCompound        *ForCompound
	SelectCompound     *SelectCompound
	GroupingCompound   *GroupingCompound
	TestCompound       *TestCompound
	ArthimeticCompound *ArithmeticExpansion
	Tokens             Tokens
}

func (cc *Compound) parse(b *bashParser) error {
	var err error

	c := b.NewGoal()

	switch c.Peek() {
	case parser.Token{Type: TokenKeyword, Data: "if"}:
		cc.IfCompound = new(IfCompound)

		err = cc.IfCompound.parse(c)
	case parser.Token{Type: TokenKeyword, Data: "case"}:
		cc.CaseCompound = new(CaseCompound)

		err = cc.CaseCompound.parse(c)
	case parser.Token{Type: TokenKeyword, Data: "while"}, parser.Token{Type: TokenKeyword, Data: "until"}:
		cc.LoopCompound = new(LoopCompound)

		err = cc.LoopCompound.parse(c)
	case parser.Token{Type: TokenKeyword, Data: "for"}:
		cc.ForCompound = new(ForCompound)

		err = cc.ForCompound.parse(c)
	case parser.Token{Type: TokenKeyword, Data: "select"}:
		cc.SelectCompound = new(SelectCompound)

		err = cc.SelectCompound.parse(c)
	case parser.Token{Type: TokenKeyword, Data: "[["}:
		cc.TestCompound = new(TestCompound)

		err = cc.TestCompound.parse(c)
	case parser.Token{Type: TokenPunctuator, Data: "(("}:
		cc.ArthimeticCompound = new(ArithmeticExpansion)

		err = cc.ArthimeticCompound.parse(c)
	case parser.Token{Type: TokenPunctuator, Data: "("}, parser.Token{Type: TokenPunctuator, Data: "{"}:
		cc.GroupingCompound = new(GroupingCompound)

		err = cc.GroupingCompound.parse(c)
	}

	if err != nil {
		return b.Error("Compound", err)
	}

	b.Score(c)

	cc.Tokens = b.ToTokens()

	return nil
}

func (cc *Compound) parseHeredocs(b *bashParser) error {
	return nil
}

type IfCompound struct {
	Tokens Tokens
}

func (i *IfCompound) parse(b *bashParser) error {
	return nil
}

type CaseCompound struct {
	Tokens Tokens
}

func (cc *CaseCompound) parse(b *bashParser) error {
	return nil
}

type LoopCompound struct {
	Tokens Tokens
}

func (l *LoopCompound) parse(b *bashParser) error {
	return nil
}

type ForCompound struct {
	Tokens Tokens
}

func (f *ForCompound) parse(b *bashParser) error {
	return nil
}

type SelectCompound struct {
	Tokens Tokens
}

func (s *SelectCompound) parse(b *bashParser) error {
	return nil
}

type TestCompound struct {
	Tokens Tokens
}

func (t *TestCompound) parse(b *bashParser) error {
	return nil
}

type GroupingCompound struct {
	Tokens Tokens
}

func (g *GroupingCompound) parse(b *bashParser) error {
	return nil
}

type Command struct {
	Vars         []Assignment
	Redirections []Redirection
	Words        []Word
	Tokens       Tokens
}

func (cc *Command) parse(b *bashParser, required bool) error {
	for {
		c := b.NewGoal()

		if b.Peek().Type == TokenIdentifierAssign {
			var a Assignment

			if err := a.parse(c); err != nil {
				return b.Error("Command", err)
			}

			cc.Vars = append(cc.Vars, a)
		} else if isRedirection(b) {
			var r Redirection

			if err := r.parse(c); err != nil {
				return b.Error("Command", err)
			}

			cc.Redirections = append(cc.Redirections, r)
		} else {
			break
		}

		b.Score(c)
		b.AcceptRunWhitespace()
	}

	c := b.NewGoal()

	for nextIsWordPart(c) || isRedirection(c) {
		b.Score(c)
		c = b.NewGoal()

		if isRedirection(b) {
			var r Redirection

			if err := r.parse(c); err != nil {
				return b.Error("Command", err)
			}

			cc.Redirections = append(cc.Redirections, r)
		} else {
			var w Word

			if err := w.parse(c); err != nil {
				return b.Error("Command", err)
			}

			cc.Words = append(cc.Words, w)
		}

		b.Score(c)

		c = b.NewGoal()

		c.AcceptRunWhitespace()
	}

	if len(cc.Words) == 0 && (required || len(cc.Redirections) == 0 && len(cc.Vars) == 0) {
		return b.Error("Command", ErrMissingWord)
	}

	cc.Tokens = b.ToTokens()

	return nil
}

func (cc *Command) parseHeredocs(b *bashParser) error {
	for n := range cc.Redirections {
		c := b.NewGoal()

		if err := cc.Redirections[n].parseHeredocs(c); err != nil {
			return b.Error("Command", err)
		}

		b.Score(c)
	}

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
	case TokenWhitespace, TokenLineTerminator, TokenComment, TokenCloseBacktick, TokenCloseParen, TokenHeredoc, TokenHeredocEnd, parser.TokenDone:
		return false
	case TokenPunctuator:
		switch tk.Data {
		case "$((", "$(", "${", "=":
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

			if err := p.Word.parse(c); err != nil {
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
	Heredoc    *Heredoc
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

func (r *Redirection) parseHeredocs(b *bashParser) error {
	if !r.isHeredoc() {
		return nil
	}

	b.AcceptRunWhitespace()

	c := b.NewGoal()
	r.Heredoc = new(Heredoc)

	if err := r.Heredoc.parse(c); err != nil {
		return b.Error("Redirection", err)
	}

	b.Score(c)

	return nil
}

type Heredoc struct {
	HeredocPartsOrWords []HeredocPartOrWord
	Tokens              Tokens
}

func (h *Heredoc) parse(b *bashParser) error {
	b.Accept(TokenHeredocIndent)

	for !b.Accept(TokenHeredocEnd) {
		c := b.NewGoal()

		var hw HeredocPartOrWord

		if err := hw.parse(c); err != nil {
			return b.Error("Heredoc", err)
		}

		h.HeredocPartsOrWords = append(h.HeredocPartsOrWords, hw)

		b.Score(c)
		b.Accept(TokenHeredocIndent)
	}

	h.Tokens = b.ToTokens()

	return nil
}

type HeredocPartOrWord struct {
	HeredocPart *Token
	Word        *Word
	Tokens      Tokens
}

func (h *HeredocPartOrWord) parse(b *bashParser) error {
	if b.Accept(TokenHeredoc) {
		h.HeredocPart = b.GetLastToken()
	} else {
		c := b.NewGoal()

		h.Word = new(Word)

		if err := h.Word.parse(c); err != nil {
			return b.Error("HeredocPartOrWord", err)
		}

		b.Score(c)
	}

	h.Tokens = b.ToTokens()

	return nil
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
