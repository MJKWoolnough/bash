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
