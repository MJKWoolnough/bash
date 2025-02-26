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

type Pipeline struct {
	Statement Statement
	Pipeline  *Pipeline
	Tokens    Tokens
}

func (p *Pipeline) parse(b *bashParser) error {
	c := b.NewGoal()

	if err := p.Statement.parse(c); err != nil {
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

type Statement struct{}

func (s *Statement) parse(b *bashParser) error {
	return nil
}
