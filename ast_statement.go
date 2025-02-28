package bash

import "vimagination.zapto.org/parser"

type LogicalOperator uint8

const (
	LogicalOperatorNone LogicalOperator = iota
	LogicalOperatorAnd
	LogicalOperatorOr
)

type LogicalExpression struct {
	Pipeline          Pipeline
	LogicalOperator   LogicalOperator
	LogicalExpression *LogicalExpression
	Tokens
}

func (l *LogicalExpression) parse(b *bashParser) error {
	c := b.NewGoal()

	if err := l.Pipeline.parse(c); err != nil {
		return b.Error("LogicalExpression", err)
	}

	b.Score(c)

	c = b.NewGoal()

	c.AcceptRunWhitespace()

	if c.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "&&"}) {
		l.LogicalOperator = LogicalOperatorAnd
	} else if c.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "||"}) {
		l.LogicalOperator = LogicalOperatorOr
	}

	if l.LogicalOperator != LogicalOperatorNone {
		c.AcceptRunWhitespace()
		b.Score(c)

		c = b.NewGoal()
		l.LogicalExpression = new(LogicalExpression)

		if err := l.LogicalExpression.parse(c); err != nil {
			return b.Error("LogicalExpression", err)
		}

		b.Score(c)
	}

	l.Tokens = b.ToTokens()

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
	RedirectionsOrVars []RedirectionOrVar
	Statement          Statement
	Redirections       []Redirection
	Tokens             Tokens
}

func (r *Redirections) parse(b *bashParser) error {
	for b.Peek().Type == TokenIdentifierAssign || isRedirection(b) {
		var rv RedirectionOrVar

		c := b.NewGoal()

		if err := r.parse(c); err != nil {
			return b.Error("Redirections", err)
		}

		b.Score(c)
		b.AcceptRunWhitespace()

		r.RedirectionsOrVars = append(r.RedirectionsOrVars, rv)

	}

	c := b.NewGoal()

	if err := r.Statement.parse(c); err != nil {
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
	return false
}

type RedirectionOrVar struct{}

func (r *RedirectionOrVar) parse(b *bashParser) error {
	return nil
}

type Redirection struct{}

func (r *Redirection) parse(b *bashParser) error {
	return nil
}

type Statement struct{}

func (s *Statement) parse(b *bashParser) error {
	return nil
}
