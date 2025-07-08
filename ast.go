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
