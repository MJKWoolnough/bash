package bash

import (
	"strings"

	"vimagination.zapto.org/parser"
)

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
	Not          bool
	Redirections Redirections
	Pipeline     *Pipeline
	Tokens       Tokens
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

	if err := p.Redirections.parse(c); err != nil {
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

type Redirections struct {
	RedirectionsOrVars []RedirectionOrAssignment
	Command            Command
	Redirections       []Redirection
	Tokens             Tokens
}

func (r *Redirections) parse(b *bashParser) error {
	for b.Peek().Type == TokenIdentifierAssign || isRedirection(b) {
		var rv RedirectionOrAssignment

		c := b.NewGoal()

		if err := r.parse(c); err != nil {
			return b.Error("Redirections", err)
		}

		b.Score(c)
		b.AcceptRunWhitespace()

		r.RedirectionsOrVars = append(r.RedirectionsOrVars, rv)

	}

	c := b.NewGoal()

	if err := r.Command.parse(c); err != nil {
		return b.Error("Redirections", err)
	}

	b.Score(c)

	c = b.NewGoal()

	c.AcceptRunWhitespace()

	for isRedirection(c) {
		b.Score(c)

		c = b.NewGoal()

		var rv Redirection

		if err := r.parse(c); err != nil {
			return b.Error("Redirections", err)
		}

		b.Score(c)

		r.Redirections = append(r.Redirections, rv)
		c = b.NewGoal()

		c.AcceptRunWhitespace()
	}

	r.Tokens = b.ToTokens()

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

type RedirectionOrAssignment struct {
	Redirection *Redirection
	Assignment  *Assignment
	Tokens      Tokens
}

func (r *RedirectionOrAssignment) parse(b *bashParser) error {
	c := b.NewGoal()

	var err error

	if b.Peek().Type == TokenIdentifierAssign {
		r.Assignment = new(Assignment)
		err = r.Assignment.parse(c)
	} else {
		r.Redirection = new(Redirection)
		err = r.Redirection.parse(c)
	}

	if err != nil {
		return b.Error("RedirectionOrAssignment", err)
	}

	b.Score(c)

	r.Tokens = b.ToTokens()

	return nil
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
	Tokens              Tokens
}

func (w *WordPart) parse(b *bashParser) error {
	c := b.NewGoal()
	switch b.Peek() {
	case parser.Token{Type: TokenPunctuator, Data: "${"}:
		w.Parameter = new(Parameter)

		if err := w.Parameter.parse(c); err != nil {
			return b.Error("WordPart", err)
		}
	case parser.Token{Type: TokenPunctuator, Data: "$("}, parser.Token{Type: TokenPunctuator, Data: "`"}:
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
	SubstitutionBacktick SubstitutionType = iota
	SubstitutionNew
)

type CommandSubstitution struct {
	SubstitutionType SubstitutionType
	Command          File
	Tokens           Tokens
}

func (cs *CommandSubstitution) parse(b *bashParser) error {
	end := "`"

	if b.Next().Data != "`" {
		cs.SubstitutionType = SubstitutionNew
		end = ")"
	}

	b.AcceptRunWhitespace()

	c := b.NewGoal()
	c.StopAt = &parser.Token{Type: TokenPunctuator, Data: end}

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

type Command struct{}

func (s *Command) parse(b *bashParser) error {
	return nil
}
