package bash

import (
	"fmt"
	"io"
)

var indent = []byte{'\t'}

type indentPrinter struct {
	io.Writer
}

func (i *indentPrinter) Write(p []byte) (int, error) {
	var (
		total int
		last  int
	)

	for n, c := range p {
		if c == '\n' {
			m, err := i.Writer.Write(p[last : n+1])
			total += m

			if err != nil {
				return total, err
			}

			_, err = i.Writer.Write(indent)
			if err != nil {
				return total, err
			}

			last = n + 1
		}
	}

	if last != len(p) {
		m, err := i.Writer.Write(p[last:])
		total += m

		if err != nil {
			return total, err
		}
	}

	return total, nil
}

func (i *indentPrinter) Print(args ...interface{}) {
	fmt.Fprint(i, args...)
}

func (i *indentPrinter) Printf(format string, args ...interface{}) {
	fmt.Fprintf(i, format, args...)
}

func (i *indentPrinter) WriteString(s string) (int, error) {
	return i.Write([]byte(s))
}

func (t Token) printType(w io.Writer, v bool) {
	var typ string

	switch t.Type {
	default:
		typ = "Unknown"
	}

	fmt.Fprintf(w, "Type: %s - Data: %q", typ, t.Data)

	if v {
		fmt.Fprintf(w, " - Position: %d (%d: %d)", t.Pos, t.Line, t.LinePos)
	}
}

func (t Tokens) printType(w io.Writer, v bool) {
	if t == nil {
		io.WriteString(w, "nil")

		return
	}

	if len(t) == 0 {
		io.WriteString(w, "[]")

		return
	}

	io.WriteString(w, "[")

	ipp := indentPrinter{w}

	for n, t := range t {
		ipp.Printf("\n%d: ", n)
		t.printType(w, v)
	}

	io.WriteString(w, "\n]")
}

type formatter interface {
	printType(io.Writer, bool)
	printSource(io.Writer, bool)
}

func format(f formatter, s fmt.State, v rune) {
	switch v {
	case 'v':
		f.printType(s, s.Flag('+'))
	case 's':
		f.printSource(s, s.Flag('+'))
	}
}
