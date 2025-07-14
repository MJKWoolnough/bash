package bash

import "vimagination.zapto.org/parser"

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

func (s *Statement) isMultiline(v bool) bool {
	if s.Pipeline.isMultiline(v) {
		return true
	} else if s.Statement != nil {
		return s.Statement.isMultiline(v)
	}

	return false
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
	PipelineTime      PipelineTime
	Not               bool
	Coproc            bool
	CoprocIdentifier  *Token
	CommandOrCompound CommandOrCompound
	Pipeline          *Pipeline
	Tokens            Tokens
}

func (p *Pipeline) parse(b *bashParser) error {
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

	if b.AcceptToken(parser.Token{Type: TokenKeyword, Data: "coproc"}) {
		p.Coproc = true

		b.AcceptRunWhitespace()

		if b.Accept(TokenIdentifier) {
			p.CoprocIdentifier = b.GetLastToken()

			b.AcceptRunWhitespace()
		}
	}

	c := b.NewGoal()

	if err := p.CommandOrCompound.parse(c); err != nil {
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

func (p *Pipeline) isMultiline(v bool) bool {
	if p.CommandOrCompound.isMultiline(v) {
		return true
	} else if p.Pipeline != nil {
		return p.Pipeline.isMultiline(v)
	}

	return false
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

func (cc *CommandOrCompound) parse(b *bashParser) error {
	var err error

	c := b.NewGoal()

	if isCompoundNext(b) {
		cc.Compound = new(Compound)
		err = cc.Compound.parse(c)
	} else {
		cc.Command = new(Command)
		err = cc.Command.parse(c)
	}

	if err != nil {
		return b.Error("CommandOrCompound", err)
	}

	b.Score(c)

	cc.Tokens = b.ToTokens()

	return nil
}

func (cc *CommandOrCompound) isMultiline(v bool) bool {
	return cc.Command != nil && cc.Command.isMultiline(v) || cc.Compound != nil && cc.Compound.isMultiline(v)
}

func (cc *CommandOrCompound) parseHeredoc(b *bashParser) error {
	var err error

	c := b.NewGoal()

	if cc.Command != nil {
		err = cc.Command.parseHeredocs(c)
	} else if cc.Compound != nil {
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

	return tk.Type == TokenKeyword && (tk.Data == "function" || tk.Data == "if" || tk.Data == "case" || tk.Data == "while" || tk.Data == "for" || tk.Data == "until" || tk.Data == "select" || tk.Data == "[[") || tk.Type == TokenPunctuator && (tk.Data == "((" || tk.Data == "(" || tk.Data == "{") || tk.Type == TokenFunctionIdentifier
}

type Command struct {
	Vars               []Assignment
	Redirections       []Redirection
	AssignmentsOrWords []AssignmentOrWord
	Tokens             Tokens
}

func (cc *Command) parse(b *bashParser) error {
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
			var a AssignmentOrWord

			if err := a.parse(c); err != nil {
				return b.Error("Command", err)
			}

			cc.AssignmentsOrWords = append(cc.AssignmentsOrWords, a)
		}

		b.Score(c)

		c = b.NewGoal()

		c.AcceptRunWhitespace()
	}

	if len(cc.AssignmentsOrWords) == 0 && len(cc.Redirections) == 0 && len(cc.Vars) == 0 {
		return b.Error("Command", ErrMissingWord)
	}

	cc.Tokens = b.ToTokens()

	return nil
}

func (cc *Command) isMultiline(v bool) bool {
	for _, vs := range cc.Vars {
		if vs.isMultiline(v) {
			return true
		}
	}

	for _, r := range cc.Redirections {
		if r.isMultiline(v) {
			return true
		}
	}

	for _, a := range cc.AssignmentsOrWords {
		if a.isMultiline(v) {
			return true
		}
	}

	return false
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
	Expression []WordOrOperator
	Value      *Value
	Tokens     Tokens
}

func (a *Assignment) parse(b *bashParser) error {
	c := b.NewGoal()

	isLet := b.Peek().Type == TokenLetIdentifierAssign

	if err := a.Identifier.parse(c); err != nil {
		return b.Error("Assignment", err)
	}

	b.Score(c)

	if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "="}) {
		a.Assignment = AssignmentAssign
	} else if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "+="}) {
		a.Assignment = AssignmentAppend
	}

	if isLet {
		parens := 0

		for {
			if tk := b.Peek(); parens == 0 && (tk.Type == TokenWhitespace || tk.Type == TokenLineTerminator || tk.Type == TokenComment || tk == (parser.Token{Type: TokenPunctuator, Data: ";"}) || isEnd(tk)) {
				break
			} else if tk == (parser.Token{Type: TokenPunctuator, Data: "("}) {
				parens++
			} else if tk == (parser.Token{Type: TokenPunctuator, Data: ")"}) {
				parens--
			}

			c := b.NewGoal()

			var w WordOrOperator

			if err := w.parse(c); err != nil {
				return b.Error("Assignment", err)
			}

			a.Expression = append(a.Expression, w)

			b.Score(c)

			if parens > 0 {
				b.AcceptRunWhitespace()
			}
		}
	} else {
		c := b.NewGoal()

		a.Value = new(Value)
		if err := a.Value.parse(c); err != nil {
			return b.Error("Assignment", err)
		}

		b.Score(c)
	}

	a.Tokens = b.ToTokens()

	return nil
}

func (a *Assignment) isMultiline(v bool) bool {
	if a.Identifier.isMultiline(v) {
		return true
	}

	if a.Value != nil {
		return a.Value.isMultiline(v)
	}

	for _, wo := range a.Expression {
		if wo.isMultiline(v) {
			return true
		}
	}

	return false
}

type ParameterAssign struct {
	Identifier *Token
	Subscript  []WordOrOperator
	Tokens     Tokens
}

func (p *ParameterAssign) parse(b *bashParser) error {
	b.Next()

	p.Identifier = b.GetLastToken()

	if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "["}) {
		b.AcceptRunWhitespace()

		for !b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "]"}) {
			c := b.NewGoal()

			var w WordOrOperator

			if err := w.parse(c); err != nil {
				return b.Error("ParameterAssign", err)
			}

			p.Subscript = append(p.Subscript, w)

			b.Score(c)
			b.AcceptRunAllWhitespace()
		}
	}

	p.Tokens = b.ToTokens()

	return nil
}

func (p *ParameterAssign) isMultiline(v bool) bool {
	for _, wo := range p.Subscript {
		if wo.isMultiline(v) {
			return true
		}
	}

	return false
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

	if err := r.Output.parse(c, false); err != nil {
		return b.Error("Redirection", err)
	}

	b.Score(c)

	r.Tokens = b.ToTokens()

	return nil
}

func (r *Redirection) isMultiline(v bool) bool {
	return r.Heredoc != nil || r.Output.isMultiline(v)
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

		if err := h.Word.parse(c, false); err != nil {
			return b.Error("HeredocPartOrWord", err)
		}

		b.Score(c)
	}

	h.Tokens = b.ToTokens()

	return nil
}
