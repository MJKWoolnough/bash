package bash

import "vimagination.zapto.org/parser"

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
	TestOperatorFilesAreSameInode
	TestOperatorFileIsNewerThan
	TestOperatorFileIsOlderThan
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

			c = b.NewGoal()
			t.Pattern = new(Pattern)

			if err := t.Pattern.parse(c); err != nil {
				return b.Error("Tests", err)
			}

			b.Score(c)
		} else if tk.Type == TokenBinaryOperator {
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

			c = b.NewGoal()
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
		c.AcceptRunAllWhitespaceNoComments()
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
