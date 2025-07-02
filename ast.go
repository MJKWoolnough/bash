// Package bash implements a bash tokeniser and AST.
package bash // import "vimagination.zapto.org/bash"

import (
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
	Lines    []Line
	Comments [2]Comments
	Tokens   Tokens
}

func (f *File) parse(b *bashParser) error {
	if b.HasTopComment() {
		f.Comments[0] = b.AcceptRunWhitespaceComments()
	}

	c := b.NewGoal()

	for {
		c.AcceptRunAllWhitespace()

		if tk := c.Peek(); isEnd(tk) {
			break
		}

		b.AcceptRunAllWhitespaceNoComments()

		c = b.NewGoal()

		var l Line

		if err := l.parse(c); err != nil {
			return b.Error("File", err)
		}

		f.Lines = append(f.Lines, l)

		b.Score(c)

		c = b.NewGoal()
	}

	f.Comments[1] = b.AcceptRunAllWhitespaceComments()
	f.Tokens = b.ToTokens()

	return nil
}

func (f *File) isMultiline(v bool) bool {
	if len(f.Lines) > 1 || len(f.Comments[0]) > 0 || len(f.Comments[1]) > 0 {
		return true
	}

	for _, l := range f.Lines {
		if l.isMultiline(v) {
			return true
		}
	}

	return false
}

func isEnd(tk parser.Token) bool {
	return tk.Type == parser.TokenDone || tk.Type == TokenCloseBacktick || tk.Type == TokenKeyword && (tk.Data == "then" || tk.Data == "elif" || tk.Data == "else" || tk.Data == "fi" || tk.Data == "esac" || tk.Data == "done") || tk.Type == TokenPunctuator && (tk.Data == ";;" || tk.Data == ";&" || tk.Data == ";;&" || tk.Data == ")" || tk.Data == "}")
}

type Line struct {
	Statements []Statement
	Comments   [2]Comments
	Tokens     Tokens
}

func (l *Line) parse(b *bashParser) error {
	l.Comments[0] = b.AcceptRunAllWhitespaceComments()

	b.AcceptRunAllWhitespace()

	c := b.NewGoal()

	for {
		if tk := c.Peek(); tk.Type == TokenComment || tk.Type == TokenLineTerminator || isEnd(tk) {
			break
		}

		b.AcceptRunWhitespace()

		c = b.NewGoal()

		var s Statement

		if err := s.parse(c, true); err != nil {
			return b.Error("Line", err)
		}

		l.Statements = append(l.Statements, s)

		b.Score(c)

		c = b.NewGoal()
		c.AcceptRunWhitespace()
	}

	l.Comments[1] = b.AcceptRunWhitespaceComments()

	if err := l.parseHeredocs(b); err != nil {
		return err
	}

	l.Tokens = b.ToTokens()

	return nil
}

func (l *Line) isMultiline(v bool) bool {
	if len(l.Comments[0]) > 0 || len(l.Comments[1]) > 0 || v && len(l.Statements) > 1 {
		return true
	}

	for _, s := range l.Statements {
		if s.isMultiline(v) {
			return true
		}
	}

	return false
}

func (l *Line) parseHeredocs(b *bashParser) error {
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

	if b.AcceptToken(parser.Token{Type: TokenKeyword, Data: "coproc"}) {
		p.Coproc = true

		b.AcceptRunWhitespace()

		if b.Accept(TokenIdentifier) {
			p.CoprocIdentifier = b.GetLastToken()

			b.AcceptRunWhitespace()
		}
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

type Compound struct {
	IfCompound         *IfCompound
	CaseCompound       *CaseCompound
	LoopCompound       *LoopCompound
	ForCompound        *ForCompound
	SelectCompound     *SelectCompound
	GroupingCompound   *GroupingCompound
	TestCompound       *TestCompound
	ArithmeticCompound *ArithmeticExpansion
	FunctionCompound   *FunctionCompound
	Redirections       []Redirection
	Tokens             Tokens
}

func (cc *Compound) parse(b *bashParser) error {
	var err error

	c := b.NewGoal()

	if c.Peek().Type == TokenFunctionIdentifier {
		cc.FunctionCompound = new(FunctionCompound)

		err = cc.FunctionCompound.parse(c)
	} else {
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
		case parser.Token{Type: TokenKeyword, Data: "function"}:
			cc.FunctionCompound = new(FunctionCompound)

			err = cc.FunctionCompound.parse(c)
		case parser.Token{Type: TokenKeyword, Data: "[["}:
			cc.TestCompound = new(TestCompound)

			err = cc.TestCompound.parse(c)
		case parser.Token{Type: TokenPunctuator, Data: "(("}:
			cc.ArithmeticCompound = new(ArithmeticExpansion)

			err = cc.ArithmeticCompound.parse(c)
		case parser.Token{Type: TokenPunctuator, Data: "("}, parser.Token{Type: TokenPunctuator, Data: "{"}:
			cc.GroupingCompound = new(GroupingCompound)

			err = cc.GroupingCompound.parse(c)
		}
	}

	if err != nil {
		return b.Error("Compound", err)
	}

	b.Score(c)

	c = b.NewGoal()

	c.AcceptRunWhitespace()

	for isRedirection(c) {
		b.Score(c)

		c = b.NewGoal()

		var r Redirection

		if err := r.parse(c); err != nil {
			return b.Error("Compound", err)
		}

		cc.Redirections = append(cc.Redirections, r)

		b.Score(c)

		c = b.NewGoal()

		c.AcceptRunWhitespace()
	}

	cc.Tokens = b.ToTokens()

	return nil
}

func (cc *Compound) isMultiline(v bool) bool {
	if cc.IfCompound != nil || cc.CaseCompound != nil || cc.LoopCompound != nil || cc.ForCompound != nil || cc.SelectCompound != nil {
		return true
	} else if cc.GroupingCompound != nil && cc.GroupingCompound.isMultiline(v) {
		return true
	} else if cc.TestCompound != nil && cc.TestCompound.isMultiline(v) {
		return true
	} else if cc.ArithmeticCompound != nil && cc.ArithmeticCompound.isMultiline(v) {
		return true
	} else if cc.FunctionCompound != nil && cc.FunctionCompound.isMultiline(v) {
		return true
	}

	for _, r := range cc.Redirections {
		if r.isMultiline(v) {
			return true
		}
	}

	return false
}

func (cc *Compound) parseHeredocs(b *bashParser) error {
	for n := range cc.Redirections {
		c := b.NewGoal()

		if err := cc.Redirections[n].parseHeredocs(c); err != nil {
			return b.Error("Compound", err)
		}

		b.Score(c)
	}

	return nil
}

type IfCompound struct {
	If     TestConsequence
	ElIf   []TestConsequence
	Else   *File
	Tokens Tokens
}

func (i *IfCompound) parse(b *bashParser) error {
	b.AcceptToken(parser.Token{Type: TokenKeyword, Data: "if"})
	b.AcceptRunAllWhitespace()

	c := b.NewGoal()

	if err := i.If.parse(c); err != nil {
		return b.Error("IfCompound", err)
	}

	b.Score(c)
	b.AcceptRunAllWhitespace()

	for b.AcceptToken(parser.Token{Type: TokenKeyword, Data: "elif"}) {
		b.AcceptRunAllWhitespace()

		c := b.NewGoal()

		var tc TestConsequence

		if err := tc.parse(c); err != nil {
			return b.Error("IfCompound", err)
		}

		i.ElIf = append(i.ElIf, tc)

		b.Score(c)
		b.AcceptRunAllWhitespace()
	}

	if b.AcceptToken(parser.Token{Type: TokenKeyword, Data: "else"}) {
		c := b.NewFileGoal()
		i.Else = new(File)

		if err := i.Else.parse(c); err != nil {
			return b.Error("IfCompound", err)
		}

		b.Score(c)
		b.AcceptRunAllWhitespace()
	}

	b.AcceptToken(parser.Token{Type: TokenKeyword, Data: "fi"})

	i.Tokens = b.ToTokens()

	return nil
}

type TestConsequence struct {
	Test        Statement
	Consequence File
	Comments    Comments
	Tokens
}

func (t *TestConsequence) parse(b *bashParser) error {
	c := b.NewGoal()

	if err := t.Test.parse(c, true); err != nil {
		return b.Error("TestConsequence", err)
	}

	b.Score(c)

	t.Comments = b.AcceptRunAllWhitespaceComments()
	b.AcceptRunAllWhitespace()

	b.AcceptToken(parser.Token{Type: TokenKeyword, Data: "then"})
	c = b.NewFileGoal()

	if err := t.Consequence.parse(c); err != nil {
		return b.Error("TestConsequence", err)
	}

	b.Score(c)

	t.Tokens = b.ToTokens()

	return nil
}

type CaseCompound struct {
	Word     Word
	Matches  []PatternLines
	Comments [3]Comments
	Tokens   Tokens
}

func (cc *CaseCompound) parse(b *bashParser) error {
	b.AcceptToken(parser.Token{Type: TokenKeyword, Data: "case"})
	b.AcceptRunAllWhitespace()

	c := b.NewGoal()

	if err := cc.Word.parse(c, false); err != nil {
		return b.Error("CaseCompound", err)
	}

	b.Score(c)

	cc.Comments[0] = b.AcceptRunAllWhitespaceComments()

	b.AcceptRunAllWhitespaceNoComments()

	b.AcceptToken(parser.Token{Type: TokenKeyword, Data: "in"})

	cc.Comments[1] = b.AcceptRunWhitespaceComments()

	b.AcceptRunAllWhitespaceNoComments()

	for {
		c := b.NewGoal()

		cc.Comments[2] = c.AcceptRunAllWhitespaceComments()

		c.AcceptRunAllWhitespaceNoComments()

		if c.AcceptToken(parser.Token{Type: TokenKeyword, Data: "esac"}) {
			b.Score(c)

			break
		}

		b.AcceptRunAllWhitespaceNoComments()

		c = b.NewGoal()

		var pl PatternLines

		if err := pl.parse(c); err != nil {
			return b.Error("CaseCompound", err)
		}

		cc.Matches = append(cc.Matches, pl)

		b.Score(c)
	}

	cc.Tokens = b.ToTokens()

	return nil
}

type CaseTerminationType uint8

const (
	CaseTerminationNone CaseTerminationType = iota
	CaseTerminationEnd
	CaseTerminationContinue
	CaseTerminationFallthrough
)

type PatternLines struct {
	Patterns []Word
	Lines    File
	CaseTerminationType
	Comments Comments
	Tokens   Tokens
}

func (pl *PatternLines) parse(b *bashParser) error {
	pl.Comments = b.AcceptRunWhitespaceComments()

	b.AcceptRunAllWhitespaceNoComments()

	for {
		c := b.NewGoal()

		var w Word

		if err := w.parse(c, false); err != nil {
			return b.Error("PatternLines", err)
		}

		pl.Patterns = append(pl.Patterns, w)

		b.Score(c)
		b.AcceptRunWhitespace()

		if !b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "|"}) {
			break
		}
	}

	if !b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ")"}) {
		return b.Error("PatternLines", ErrMissingClosingPattern)
	}

	c := b.NewFileGoal()

	if err := pl.Lines.parse(c); err != nil {
		return b.Error("PatternLines", err)
	}

	b.Score(c)

	c = b.NewGoal()

	c.AcceptRunAllWhitespace()

	if c.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ";;"}) {
		b.Score(c)

		pl.CaseTerminationType = CaseTerminationEnd
	} else if c.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ";&"}) {
		b.Score(c)

		pl.CaseTerminationType = CaseTerminationContinue
	} else if c.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ";;&"}) {
		b.Score(c)

		pl.CaseTerminationType = CaseTerminationFallthrough
	}

	pl.Tokens = b.ToTokens()

	return nil
}

type LoopCompound struct {
	Until     bool
	Statement Statement
	File      File
	Comments  Comments
	Tokens    Tokens
}

func (l *LoopCompound) parse(b *bashParser) error {
	if !b.AcceptToken(parser.Token{Type: TokenKeyword, Data: "while"}) {
		b.AcceptToken(parser.Token{Type: TokenKeyword, Data: "until"})

		l.Until = true
	}

	b.AcceptRunWhitespace()

	c := b.NewGoal()

	if err := l.Statement.parse(c, true); err != nil {
		return b.Error("LoopCompound", err)
	}

	b.Score(c)

	l.Comments = b.AcceptRunAllWhitespaceComments()

	b.AcceptRunAllWhitespaceNoComments()
	b.AcceptToken(parser.Token{Type: TokenKeyword, Data: "do"})

	c = b.NewFileGoal()

	if err := l.File.parse(c); err != nil {
		return b.Error("LoopCompound", err)
	}

	b.Score(c)
	b.AcceptRunAllWhitespace()
	b.AcceptToken(parser.Token{Type: TokenKeyword, Data: "done"})

	l.Tokens = b.ToTokens()

	return nil
}

type ForCompound struct {
	Identifier          *Token
	Words               []Word
	ArithmeticExpansion *ArithmeticExpansion
	File                File
	Comments            [2]Comments
	Tokens              Tokens
}

func (f *ForCompound) parse(b *bashParser) error {
	b.AcceptToken(parser.Token{Type: TokenKeyword, Data: "for"})
	b.AcceptRunWhitespace()

	if b.Accept(TokenIdentifier) {
		f.Identifier = b.GetLastToken()

		f.Comments[0] = b.AcceptRunAllWhitespaceComments()

		b.AcceptRunAllWhitespaceNoComments()

		if b.AcceptToken(parser.Token{Type: TokenKeyword, Data: "in"}) {
			b.AcceptRunWhitespace()

			f.Words = []Word{}

			for {
				if tk := b.Peek(); tk == (parser.Token{Type: TokenPunctuator, Data: ";"}) || tk.Type == TokenLineTerminator || tk.Type == TokenComment {
					break
				}

				c := b.NewGoal()

				var w Word

				if err := w.parse(c, false); err != nil {
					return b.Error("ForCompound", err)
				}

				f.Words = append(f.Words, w)

				b.Score(c)
				b.AcceptRunWhitespace()
			}
		}
	} else {
		c := b.NewGoal()
		f.ArithmeticExpansion = new(ArithmeticExpansion)

		if err := f.ArithmeticExpansion.parse(c); err != nil {
			return b.Error("ForCompound", err)
		}

		b.Score(c)
	}

	if f.Comments[1] = b.AcceptRunAllWhitespaceComments(); len(f.Comments[1]) == 0 {
		b.AcceptRunWhitespace()
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ";"})

		f.Comments[1] = b.AcceptRunAllWhitespaceComments()
	}

	b.AcceptRunAllWhitespaceNoComments()
	b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ";"})
	b.AcceptRunAllWhitespace()
	b.AcceptToken(parser.Token{Type: TokenKeyword, Data: "do"})
	b.AcceptRunAllWhitespace()

	c := b.NewGoal()

	if err := f.File.parse(c); err != nil {
		return b.Error("ForCompound", err)
	}

	b.Score(c)
	b.AcceptRunAllWhitespace()
	b.AcceptToken(parser.Token{Type: TokenKeyword, Data: "done"})

	f.Tokens = b.ToTokens()

	return nil
}

type SelectCompound struct {
	Identifier *Token
	Words      []Word
	File       File
	Comments   [2]Comments
	Tokens     Tokens
}

func (s *SelectCompound) parse(b *bashParser) error {
	b.AcceptToken(parser.Token{Type: TokenKeyword, Data: "select"})
	b.AcceptRunWhitespace()
	b.Accept(TokenIdentifier)

	s.Identifier = b.GetLastToken()

	s.Comments[0] = b.AcceptRunAllWhitespaceComments()

	b.AcceptRunAllWhitespaceNoComments()

	if b.AcceptToken(parser.Token{Type: TokenKeyword, Data: "in"}) {
		b.AcceptRunWhitespace()

		s.Words = []Word{}

		for {
			if tk := b.Peek(); tk == (parser.Token{Type: TokenPunctuator, Data: ";"}) || tk.Type == TokenLineTerminator || tk.Type == TokenComment {
				break
			}

			c := b.NewGoal()

			var w Word

			if err := w.parse(c, false); err != nil {
				return b.Error("SelectCompound", err)
			}

			s.Words = append(s.Words, w)

			b.Score(c)
			b.AcceptRunWhitespace()
		}
	}

	b.AcceptRunWhitespace()
	b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ";"})

	s.Comments[1] = b.AcceptRunAllWhitespaceComments()

	b.AcceptRunAllWhitespaceNoComments()
	b.AcceptToken(parser.Token{Type: TokenKeyword, Data: "do"})

	c := b.NewFileGoal()

	if err := s.File.parse(c); err != nil {
		return b.Error("SelectCompound", err)
	}

	b.Score(c)
	b.AcceptRunAllWhitespace()
	b.AcceptToken(parser.Token{Type: TokenKeyword, Data: "done"})

	s.Tokens = b.ToTokens()

	return nil
}

type TestCompound struct {
	Tests    Tests
	Comments [2]Comments
	Tokens   Tokens
}

func (t *TestCompound) parse(b *bashParser) error {
	b.AcceptToken(parser.Token{Type: TokenKeyword, Data: "[["})

	t.Comments[0] = b.AcceptRunWhitespaceComments()

	b.AcceptRunAllWhitespaceNoComments()

	c := b.NewGoal()

	if err := t.Tests.parse(c); err != nil {
		return b.Error("TestCompound", err)
	}

	b.Score(c)

	t.Comments[1] = b.AcceptRunAllWhitespaceComments()

	b.AcceptRunAllWhitespaceNoComments()
	b.AcceptToken(parser.Token{Type: TokenKeyword, Data: "]]"})

	t.Tokens = b.ToTokens()

	return nil
}

func (t *TestCompound) isMultiline(v bool) bool {
	return len(t.Comments[0]) > 0 || len(t.Comments[1]) > 0 || t.Tests.isMultiline(v)
}

type TestOperator uint8

const (
	TestOperatorNone TestOperator = iota
	TestOperatorFileExists
	TestOperatorFileIsBlock
	TestOperatorFileIsCharacter
	TestOperatorDirectoryExists
	TestOperatorFileIsRegular
	TestOperatorFileHasSetGroupID
	TestOperatorFileIsSymbolic
	TestOperatorFileHasStickyBit
	TestOperatorFileIsPipe
	TestOperatorFileIsReadable
	TestOperatorFileIsNonZero
	TestOperatorFileIsTerminal
	TestOperatorFileHasSetUserID
	TestOperatorFileIsWritable
	TestOperatorFileIsExecutable
	TestOperatorFileIsOwnedByEffectiveGroup
	TestOperatorFileWasModifiedSinceLastRead
	TestOperatorFileIsOwnedByEffectiveUser
	TestOperatorFileIsSocket
	TestOperatorFilesAreSameInode
	TestOperatorFileIsNewerThan
	TestOperatorFileIsOlderThan
	TestOperatorOptNameIsEnabled
	TestOperatorVarNameIsSet
	TestOperatorVarnameIsRef
	TestOperatorStringIsZero
	TestOperatorStringIsNonZero
	TestOperatorStringsEqual
	TestOperatorStringsMatch
	TestOperatorStringsNotEqual
	TestOperatorStringBefore
	TestOperatorStringAfter
	TestOperatorEqual
	TestOperatorNotEqual
	TestOperatorLessThan
	TestOperatorLessThanEqual
	TestOperatorGreaterThan
	TestOperatorGreaterThanEqual
)

type Tests struct {
	Not             bool
	Test            TestOperator
	Word            *Word
	Pattern         *Pattern
	Parens          *Tests
	LogicalOperator LogicalOperator
	Tests           *Tests
	Comments        [5]Comments
	Tokens          Tokens
}

func (t *Tests) parse(b *bashParser) error {
	t.Comments[0] = b.AcceptRunAllWhitespaceComments()
	b.AcceptRunAllWhitespaceNoComments()

	if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "!"}) {
		t.Not = true

		t.Comments[1] = b.AcceptRunAllWhitespaceComments()

		b.AcceptRunAllWhitespaceNoComments()
	}

	if tk := b.Peek(); tk.Type == TokenKeyword {
		switch tk.Data {
		case "-a", "-e":
			t.Test = TestOperatorFileExists
		case "-b":
			t.Test = TestOperatorFileIsBlock
		case "-c":
			t.Test = TestOperatorFileIsCharacter
		case "-d":
			t.Test = TestOperatorDirectoryExists
		case "-f":
			t.Test = TestOperatorFileIsRegular
		case "-g":
			t.Test = TestOperatorFileHasSetGroupID
		case "-h", "-L":
			t.Test = TestOperatorFileIsSymbolic
		case "-k":
			t.Test = TestOperatorFileHasStickyBit
		case "-p":
			t.Test = TestOperatorFileIsPipe
		case "-r":
			t.Test = TestOperatorFileIsReadable
		case "-s":
			t.Test = TestOperatorFileIsNonZero
		case "-t":
			t.Test = TestOperatorFileIsTerminal
		case "-u":
			t.Test = TestOperatorFileHasSetUserID
		case "-w":
			t.Test = TestOperatorFileIsWritable
		case "-x":
			t.Test = TestOperatorFileIsExecutable
		case "-G":
			t.Test = TestOperatorFileIsOwnedByEffectiveGroup
		case "-N":
			t.Test = TestOperatorFileWasModifiedSinceLastRead
		case "-O":
			t.Test = TestOperatorFileIsOwnedByEffectiveUser
		case "-S":
			t.Test = TestOperatorFileIsSocket
		case "-o":
			t.Test = TestOperatorOptNameIsEnabled
		case "-v":
			t.Test = TestOperatorVarNameIsSet
		case "-R":
			t.Test = TestOperatorVarnameIsRef
		case "-z":
			t.Test = TestOperatorStringIsZero
		case "-n":
			t.Test = TestOperatorStringIsNonZero
		}

		b.Next()
		b.AcceptRunWhitespace()

		c := b.NewGoal()
		t.Word = new(Word)

		if err := t.Word.parse(c, false); err != nil {
			return b.Error("Tests", err)
		}

		b.Score(c)
	} else if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "("}) {
		t.Comments[2] = b.AcceptRunWhitespaceComments()

		b.AcceptRunAllWhitespaceNoComments()

		c := b.NewGoal()
		t.Parens = new(Tests)

		if err := t.Parens.parse(c); err != nil {
			return b.Error("Tests", err)
		}

		b.Score(c)

		t.Comments[3] = b.AcceptRunAllWhitespaceComments()

		b.AcceptRunAllWhitespaceNoComments()

		if !b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ")"}) {
			return b.Error("Tests", ErrMissingClosingParen)
		}
	} else {
		c := b.NewGoal()
		t.Word = new(Word)

		if err := t.Word.parse(c, false); err != nil {
			return b.Error("Tests", err)
		}

		b.Score(c)

		c = b.NewGoal()
		c.AcceptRunWhitespace()

		if tk := c.Peek(); tk.Type == TokenKeyword && tk.Data != "]]" {
			b.Score(c)

			switch tk.Data {
			case "-ef":
				t.Test = TestOperatorFilesAreSameInode
			case "-nt":
				t.Test = TestOperatorFileIsNewerThan
			case "-ot":
				t.Test = TestOperatorFileIsOlderThan
			case "-eq":
				t.Test = TestOperatorEqual
			case "-ne":
				t.Test = TestOperatorNotEqual
			case "-lt":
				t.Test = TestOperatorLessThan
			case "-le":
				t.Test = TestOperatorLessThanEqual
			case "-gt":
				t.Test = TestOperatorGreaterThan
			case "-ge":
				t.Test = TestOperatorGreaterThanEqual
			}

			b.Next()

			b.AcceptRunWhitespace()

			c := b.NewGoal()
			t.Pattern = new(Pattern)

			if err := t.Pattern.parse(c); err != nil {
				return b.Error("Tests", err)
			}

			b.Score(c)
		} else if tk.Type == TokenOperator {
			b.Score(c)

			switch tk.Data {
			case "=", "==":
				t.Test = TestOperatorStringsEqual
			case "!=":
				t.Test = TestOperatorStringsNotEqual
			case "=~":
				t.Test = TestOperatorStringsMatch
			case "<":
				t.Test = TestOperatorStringBefore
			case ">":
				t.Test = TestOperatorStringAfter
			}

			b.Next()

			b.AcceptRunWhitespace()

			c := b.NewGoal()
			t.Pattern = new(Pattern)

			if err := t.Pattern.parse(c); err != nil {
				return b.Error("Tests", err)
			}

			b.Score(c)
		}
	}

	c := b.NewGoal()

	t.Comments[4] = c.AcceptRunAllWhitespaceComments()

	c.AcceptRunAllWhitespaceNoComments()

	if c.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "||"}) {
		t.LogicalOperator = LogicalOperatorOr
	} else if c.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "&&"}) {
		t.LogicalOperator = LogicalOperatorAnd
	}

	if t.LogicalOperator != LogicalOperatorNone {
		c.AcceptRunAllWhitespace()
		b.Score(c)

		c = b.NewGoal()
		t.Tests = new(Tests)

		if err := t.Tests.parse(c); err != nil {
			return b.Error("Tests", err)
		}

		b.Score(c)
	} else {
		t.Comments[4] = b.AcceptRunWhitespaceComments()
	}

	t.Tokens = b.ToTokens()

	return nil
}

func (t *Tests) isMultiline(v bool) bool {
	if len(t.Comments[0]) > 0 || len(t.Comments[4]) > 0 ||
		t.Not && len(t.Comments[1]) > 0 ||
		t.Parens != nil && (len(t.Comments[2]) > 0 || len(t.Comments[3]) > 0) ||
		len(t.Comments[2]) > 0 && t.Word != nil && (t.Pattern != nil && t.Test >= TestOperatorStringsEqual || t.Test >= TestOperatorFileExists && t.Test <= TestOperatorVarnameIsRef) {
		return true
	}

	if t.Parens != nil && t.Parens.isMultiline(v) {
		return true
	}

	if t.Word != nil && t.Word.isMultiline(v) {
		return true
	}

	if t.Pattern != nil && t.Pattern.isMultiline(v) {
		return true
	}

	if t.Tests != nil {
		return t.Tests.isMultiline(v)
	}

	return false
}

type Pattern struct {
	Parts  []WordPart
	Tokens Tokens
}

func (p *Pattern) parse(b *bashParser) error {
	for nextIsPatternPart(b) {
		c := b.NewGoal()

		var pp WordPart

		if err := pp.parse(c); err != nil {
			return b.Error("Pattern", err)
		}

		p.Parts = append(p.Parts, pp)

		b.Score(c)
	}

	p.Tokens = b.ToTokens()

	return nil
}

func (p *Pattern) isMultiline(v bool) bool {
	for _, pt := range p.Parts {
		if pt.isMultiline(v) {
			return true
		}
	}

	return false
}

func nextIsPatternPart(b *bashParser) bool {
	switch tk := b.Peek(); tk.Type {
	case TokenWhitespace, TokenLineTerminator, TokenComment, TokenKeyword:
		return false
	case TokenPunctuator:
		switch tk.Data {
		case ")":
			return false
		}
	}

	return true
}

type GroupingCompound struct {
	SubShell bool
	File
	Tokens Tokens
}

func (g *GroupingCompound) parse(b *bashParser) error {
	if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "("}) {
		g.SubShell = true
	} else {
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "{"})
	}

	c := b.NewFileGoal()

	if err := g.File.parse(c); err != nil {
		return b.Error("GroupingCompound", err)
	}

	b.Score(c)
	b.AcceptRunAllWhitespace()
	b.Next()

	g.Tokens = b.ToTokens()

	return nil
}

type FunctionCompound struct {
	HasKeyword bool
	Identifier *Token
	Body       Compound
	Comments   Comments
	Tokens     Tokens
}

func (f *FunctionCompound) parse(b *bashParser) error {
	if b.AcceptToken(parser.Token{Type: TokenKeyword, Data: "function"}) {
		f.HasKeyword = true

		b.AcceptRunWhitespace()
	}

	b.Accept(TokenFunctionIdentifier)

	f.Identifier = b.GetLastToken()

	b.AcceptRunWhitespace()

	if b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "("}) {
		b.AcceptRunWhitespace()
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ")"})
		b.AcceptRunWhitespace()
	}

	f.Comments = b.AcceptRunAllWhitespaceComments()

	b.AcceptRunAllWhitespaceNoComments()

	c := b.NewGoal()

	if err := f.Body.parse(c); err != nil {
		return b.Error("FunctionCompound", err)
	}

	b.Score(c)

	f.Tokens = b.ToTokens()

	return nil
}

func (f *FunctionCompound) isMultiline(v bool) bool {
	return len(f.Comments) > 0 || f.Body.isMultiline(v)
}

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

type Command struct {
	Vars               []Assignment
	Redirections       []Redirection
	AssignmentsOrWords []AssignmentOrWord
	Tokens             Tokens
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

	if len(cc.AssignmentsOrWords) == 0 && (required || len(cc.Redirections) == 0 && len(cc.Vars) == 0) {
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
	case TokenWhitespace, TokenLineTerminator, TokenComment, TokenCloseBacktick, TokenHeredoc, TokenHeredocEnd, parser.TokenDone:
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
	if w.Word != nil {
		return w.Word.isMultiline(v)
	}

	return false
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

func (cs *CommandSubstitution) isMultiline(v bool) bool {
	return cs.Command.isMultiline(v)
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

type ArithmeticExpansion struct {
	Expression        bool
	WordsAndOperators []WordOrOperator
	Tokens            Tokens
}

func (a *ArithmeticExpansion) parse(b *bashParser) error {
	if !b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "$(("}) {
		b.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "(("})

		a.Expression = true
	}

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

func (a *ArithmeticExpansion) isMultiline(v bool) bool {
	for _, w := range a.WordsAndOperators {
		if w.isMultiline(v) {
			return true
		}
	}

	return false
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
	if w.Word != nil {
		return w.Word.isMultiline(v)
	}

	return false
}

func (w *WordOrOperator) operatorIsToken(tk parser.Token) bool {
	if w.Operator != nil {
		return w.Operator.Token == tk
	}

	return false
}
