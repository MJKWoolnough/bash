// Package bash implements a bash tokeniser and AST.
package bash // import "vimagination.zapto.org/bash"

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

type File struct{}

func (f *File) parse(p *bashParser) error {
	return nil
}
