// Package bash implements a bash tokeniser and AST.
package bash // import "vimagination.zapto.org/bash"

import (
	"strings"

	"vimagination.zapto.org/parser"
)

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

func (f *File) parse(p *bashParser) error {
	q := p.NewGoal()

	for q.AcceptRunAllWhitespace() != parser.TokenDone {
		p.AcceptRunAllWhitespace()

		var s Statement

		q = p.NewGoal()

		if err := s.parse(q, true); err != nil {
			return p.Error("File", err)
		}

		f.Statements = append(f.Statements, s)

		p.Score(q)

		q = p.NewGoal()
	}

	f.Tokens = p.ToTokens()

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
		} else if tk := c.Peek().Type; tk != TokenLineTerminator && tk != parser.TokenDone {
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

	if err := p.Command.parse(c); err != nil {
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

func (c *Command) parse(b *bashParser) error {
	for {
		if b.Peek().Type == TokenIdentifierAssign {
			var a Assignment

			d := b.NewGoal()

			if err := a.parse(d); err != nil {
				return b.Error("Command", err)
			}

			b.Score(d)

			c.Vars = append(c.Vars, a)

		} else if isRedirection(b) {
			var r Redirection

			d := b.NewGoal()

			if err := r.parse(d); err != nil {
				return b.Error("Command", err)
			}

			b.Score(d)

			c.Redirections = append(c.Redirections, r)
		} else {
			break
		}

		b.AcceptRunWhitespace()
	}

	for nextIsWordPart(b) {
		if isRedirection(b) {
			var r Redirection

			d := b.NewGoal()

			if err := r.parse(d); err != nil {
				return b.Error("Command", err)
			}

			b.Score(d)

			c.Redirections = append(c.Redirections, r)
		} else {
			d := b.NewGoal()

			var w Word

			if err := w.parse(d); err != nil {
				return b.Error("Command", err)
			}

			b.Score(d)

			c.Words = append(c.Words, w)
		}
	}

	c.Tokens = b.ToTokens()

	return nil
}

func isRedirection(b *bashParser) bool {
	c := b.NewGoal()

	if c.Accept(TokenWord) {
		for _, r := range c.GetLastToken().Data {
			if !strings.ContainsRune(decimalDigit, r) {
				return false
			}
		}

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
	Identifier Paramater
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

type Paramater struct {
	Identifier *Token
	Subscript  *Word
	Tokens     Tokens
}

func (p *Paramater) parse(b *bashParser) error {
	b.Next()

	p.Identifier = b.GetLastToken()

	if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "["}) {
		b.AcceptRunAllWhitespace()

		c := b.NewGoal()
		p.Subscript = new(Word)

		if err := p.Subscript.parse(c); err != nil {
			return b.Error("Parameter", err)
		}

		b.Score(c)
		b.AcceptRunAllWhitespace()

		if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "]"}) {
			return b.Error("Parameter", ErrMissingClosingBracket)
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

		for b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ")"}) {
			c := b.NewGoal()

			var w Word

			if err := w.parse(c); err != nil {
				return b.Error("Value", err)
			}

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
	}

	w.Tokens = b.ToTokens()

	return nil
}

func nextIsWordPart(b *bashParser) bool {
	switch tk := b.Peek(); tk.Type {
	case TokenWhitespace, TokenLineTerminator, TokenComment:
		return false
	case TokenPunctuator:
		switch tk.Data {
		case "$(", "${", "`":
			return true
		}

		return false
	}

	return true
}

type WordPart struct {
	Part                *Token
	Parameter           *Parameter
	CommandSubstitution *CommandSubstitution
	ArithmeticExpansion *ArithmeticExpansion
	Tokens              Tokens
}

func (w *WordPart) parse(b *bashParser) error {
	c := b.NewGoal()

	switch tk := b.Peek(); {
	case tk == parser.Token{Type: TokenPunctuator, Data: "${"}:
		w.Parameter = new(Parameter)

		if err := w.Parameter.parse(c); err != nil {
			return b.Error("WordPart", err)
		}
	case tk == parser.Token{Type: TokenPunctuator, Data: "$("}:
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

type Parameter struct{}

func (p *Parameter) parse(b *bashParser) error {
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
	end := parser.Token{Type: TokenPunctuator, Data: ")"}

	if tk := b.Next(); tk.Type != TokenOpenBacktick {
		cs.SubstitutionType = SubstitutionBacktick
		end = parser.Token{Type: TokenCloseBacktick, Data: tk.Data}
	}

	b.AcceptRunWhitespace()

	c := b.NewGoal()
	c.StopAt = &end

	if err := cs.Command.parse(c); err != nil {
		return err
	}

	b.Score(c)

	cs.Tokens = b.ToTokens()

	return nil
}

type Redirection struct{}

func (r *Redirection) parse(b *bashParser) error {
	return nil
}

type ArithmeticExpansion struct{}

func (a *ArithmeticExpansion) parse(b *bashParser) error {
	return nil
}
