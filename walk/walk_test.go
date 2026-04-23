package walk

import (
	"errors"
	"reflect"
	"testing"

	"vimagination.zapto.org/bash"
	"vimagination.zapto.org/parser"
)

var (
	sentinel = errors.New("")
	nilErr   = errors.New("nil received")
	nilRet   = func(_ *bash.File) bash.Type { return nil }
)

type walker struct {
	end   bash.Type
	level []string
}

func (w *walker) Handle(t bash.Type) error {
	if reflect.ValueOf(t).IsNil() {
		return nilErr
	}

	if t == w.end {
		w.level = append(w.level, reflect.TypeOf(t).Elem().Name())

		return sentinel
	}

	err := Walk(t, w)
	if err != nil {
		w.level = append(w.level, reflect.TypeOf(t).Elem().Name())
	}

	return err
}

func TestWalk(t *testing.T) {
	for n, test := range [...]struct {
		Input string
		End   func(m *bash.File) bash.Type
		Level []string
	}{
		{ // 1
			"",
			nilRet,
			nil,
		},
		{ // 2
			"a;\nb;",
			func(f *bash.File) bash.Type {
				return &f.Lines[0]
			},
			[]string{"File", "Line"},
		},
		{ // 3
			"a;\nb;",
			func(f *bash.File) bash.Type {
				return &f.Lines[1]
			},
			[]string{"File", "Line"},
		},
		{ // 4
			"a;b;",
			nilRet,
			nil,
		},
		{ // 5
			"a;b;",
			func(f *bash.File) bash.Type {
				return &f.Lines[0].Statements[0]
			},
			[]string{"File", "Line", "Statement"},
		},
		{ // 6
			"a;b;",
			func(f *bash.File) bash.Type {
				return &f.Lines[0].Statements[1]
			},
			[]string{"File", "Line", "Statement"},
		},
	} {
		tk := parser.NewStringTokeniser(test.Input)

		m, err := bash.Parse(&tk)
		if err != nil {
			t.Errorf("test %d: unexpected error parsing script: %s", n+1, err)
		} else {
			w := walker{end: test.End(m)}

			if err := w.Handle(m); err == nil && test.Level != nil {
				t.Errorf("test %d: expected to recieve sentinel error, but didn't", n+1)
			} else if err != nil && test.Level == nil {
				t.Errorf("test %d: expected no error, but recieved %v", n+1, err)
			} else if len(w.level) != len(test.Level) {
				t.Errorf("test %d: expected to have %d levels, got %d", n+1, len(test.Level), len(w.level))
			} else {
				for m, l := range w.level {
					if e := test.Level[len(test.Level)-m-1]; e != l {
						t.Errorf("test %d.%d: expected to read level %s, got %s", n+1, m+1, e, l)
					}
				}
			}
		}
	}
}
