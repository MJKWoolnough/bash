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
		{ // 7
			"a || b",
			nilRet,
			nil,
		},
		{ // 8
			"a || b",
			func(f *bash.File) bash.Type {
				return &f.Lines[0].Statements[0].Pipeline
			},
			[]string{"File", "Line", "Statement", "Pipeline"},
		},
		{ // 9
			"a || b",
			func(f *bash.File) bash.Type {
				return f.Lines[0].Statements[0].Statement
			},
			[]string{"File", "Line", "Statement", "Statement"},
		},
		{ // 10
			"a | b",
			nilRet,
			nil,
		},
		{ // 11
			"a | b",
			func(f *bash.File) bash.Type {
				return &f.Lines[0].Statements[0].Pipeline.CommandOrCompound
			},
			[]string{"File", "Line", "Statement", "Pipeline", "CommandOrCompound"},
		},
		{ // 12
			"a | b",
			func(f *bash.File) bash.Type {
				return f.Lines[0].Statements[0].Pipeline.Pipeline
			},
			[]string{"File", "Line", "Statement", "Pipeline", "Pipeline"},
		},
		{ // 13
			"a",
			nilRet,
			nil,
		},
		{ // 14
			"a",
			func(f *bash.File) bash.Type {
				return f.Lines[0].Statements[0].Pipeline.CommandOrCompound.Command
			},
			[]string{"File", "Line", "Statement", "Pipeline", "CommandOrCompound", "Command"},
		},
		{ // 15
			"(a)",
			func(f *bash.File) bash.Type {
				return f.Lines[0].Statements[0].Pipeline.CommandOrCompound.Compound
			},
			[]string{"File", "Line", "Statement", "Pipeline", "CommandOrCompound", "Compound"},
		},
		{ // 16
			"a=1 b=2 c d >e 2>f",
			nilRet,
			nil,
		},
		{ // 17
			"a=1 b=2 c d >e 2>f",
			func(f *bash.File) bash.Type {
				return &f.Lines[0].Statements[0].Pipeline.CommandOrCompound.Command.Vars[0]
			},
			[]string{"File", "Line", "Statement", "Pipeline", "CommandOrCompound", "Command", "Assignment"},
		},
		{ // 18
			"a=1 b=2 c d >e 2>f",
			func(f *bash.File) bash.Type {
				return &f.Lines[0].Statements[0].Pipeline.CommandOrCompound.Command.Vars[1]
			},
			[]string{"File", "Line", "Statement", "Pipeline", "CommandOrCompound", "Command", "Assignment"},
		},
		{ // 19
			"a=1 b=2 c d >e 2>f",
			func(f *bash.File) bash.Type {
				return &f.Lines[0].Statements[0].Pipeline.CommandOrCompound.Command.AssignmentsOrWords[0]
			},
			[]string{"File", "Line", "Statement", "Pipeline", "CommandOrCompound", "Command", "AssignmentOrWord"},
		},
		{ // 20
			"a=1 b=2 c d >e 2>f",
			func(f *bash.File) bash.Type {
				return &f.Lines[0].Statements[0].Pipeline.CommandOrCompound.Command.AssignmentsOrWords[1]
			},
			[]string{"File", "Line", "Statement", "Pipeline", "CommandOrCompound", "Command", "AssignmentOrWord"},
		},
		{ // 21
			"a=1 b=2 c d >e 2>f",
			func(f *bash.File) bash.Type {
				return &f.Lines[0].Statements[0].Pipeline.CommandOrCompound.Command.Redirections[0]
			},
			[]string{"File", "Line", "Statement", "Pipeline", "CommandOrCompound", "Command", "Redirection"},
		},
		{ // 22
			"a=1 b=2 c d >e 2>f",
			func(f *bash.File) bash.Type {
				return &f.Lines[0].Statements[0].Pipeline.CommandOrCompound.Command.Redirections[1]
			},
			[]string{"File", "Line", "Statement", "Pipeline", "CommandOrCompound", "Command", "Redirection"},
		},
		{ // 23
			"a=1",
			nilRet,
			nil,
		},
		{ // 24
			"a=1",
			func(f *bash.File) bash.Type {
				return &f.Lines[0].Statements[0].Pipeline.CommandOrCompound.Command.Vars[0].Identifier
			},
			[]string{"File", "Line", "Statement", "Pipeline", "CommandOrCompound", "Command", "Assignment", "ParameterAssign"},
		},
		{ // 25
			"a=1",
			func(f *bash.File) bash.Type {
				return f.Lines[0].Statements[0].Pipeline.CommandOrCompound.Command.Vars[0].Value
			},
			[]string{"File", "Line", "Statement", "Pipeline", "CommandOrCompound", "Command", "Assignment", "Value"},
		},
		{ // 26
			"let a=1 b=2",
			nilRet,
			nil,
		},
		{ // 27
			"let a=1 b=2",
			func(f *bash.File) bash.Type {
				return &f.Lines[0].Statements[0].Pipeline.CommandOrCompound.Command.AssignmentsOrWords[1].Assignment.Expression[0]
			},
			[]string{"File", "Line", "Statement", "Pipeline", "CommandOrCompound", "Command", "AssignmentOrWord", "Assignment", "WordOrOperator"},
		},
		{ // 28
			"let a=1+2",
			func(f *bash.File) bash.Type {
				return &f.Lines[0].Statements[0].Pipeline.CommandOrCompound.Command.AssignmentsOrWords[1].Assignment.Expression[2]
			},
			[]string{"File", "Line", "Statement", "Pipeline", "CommandOrCompound", "Command", "AssignmentOrWord", "Assignment", "WordOrOperator"},
		},
		{ // 29
			"let a=1",
			func(f *bash.File) bash.Type {
				return f.Lines[0].Statements[0].Pipeline.CommandOrCompound.Command.AssignmentsOrWords[0].Word
			},
			[]string{"File", "Line", "Statement", "Pipeline", "CommandOrCompound", "Command", "AssignmentOrWord", "Word"},
		},
		{ // 30
			"let a=1",
			func(f *bash.File) bash.Type {
				return f.Lines[0].Statements[0].Pipeline.CommandOrCompound.Command.AssignmentsOrWords[1].Assignment
			},
			[]string{"File", "Line", "Statement", "Pipeline", "CommandOrCompound", "Command", "AssignmentOrWord", "Assignment"},
		},
		{ // 31
			"a$b",
			nilRet,
			nil,
		},
		{ // 32
			"a$b",
			func(f *bash.File) bash.Type {
				return &f.Lines[0].Statements[0].Pipeline.CommandOrCompound.Command.AssignmentsOrWords[0].Word.Parts[0]
			},
			[]string{"File", "Line", "Statement", "Pipeline", "CommandOrCompound", "Command", "AssignmentOrWord", "Word", "WordPart"},
		},
		{ // 33
			"a",
			nilRet,
			nil,
		},
		{ // 34
			"${a}",
			func(f *bash.File) bash.Type {
				return f.Lines[0].Statements[0].Pipeline.CommandOrCompound.Command.AssignmentsOrWords[0].Word.Parts[0].ParameterExpansion
			},
			[]string{"File", "Line", "Statement", "Pipeline", "CommandOrCompound", "Command", "AssignmentOrWord", "Word", "WordPart", "ParameterExpansion"},
		},
		{ // 35
			"$((a))",
			func(f *bash.File) bash.Type {
				return f.Lines[0].Statements[0].Pipeline.CommandOrCompound.Command.AssignmentsOrWords[0].Word.Parts[0].ArithmeticExpansion
			},
			[]string{"File", "Line", "Statement", "Pipeline", "CommandOrCompound", "Command", "AssignmentOrWord", "Word", "WordPart", "ArithmeticExpansion"},
		},
		{ // 36
			"$(a)",
			func(f *bash.File) bash.Type {
				return f.Lines[0].Statements[0].Pipeline.CommandOrCompound.Command.AssignmentsOrWords[0].Word.Parts[0].CommandSubstitution
			},
			[]string{"File", "Line", "Statement", "Pipeline", "CommandOrCompound", "Command", "AssignmentOrWord", "Word", "WordPart", "CommandSubstitution"},
		},
		{ // 37
			"{a,b}",
			func(f *bash.File) bash.Type {
				return f.Lines[0].Statements[0].Pipeline.CommandOrCompound.Command.AssignmentsOrWords[0].Word.Parts[0].BraceExpansion
			},
			[]string{"File", "Line", "Statement", "Pipeline", "CommandOrCompound", "Command", "AssignmentOrWord", "Word", "WordPart", "BraceExpansion"},
		},
		{ // 38
			"${a}",
			nilRet,
			nil,
		},
		{ // 39
			"${a/b/c}",
			func(f *bash.File) bash.Type {
				return &f.Lines[0].Statements[0].Pipeline.CommandOrCompound.Command.AssignmentsOrWords[0].Word.Parts[0].ParameterExpansion.Parameter
			},
			[]string{"File", "Line", "Statement", "Pipeline", "CommandOrCompound", "Command", "AssignmentOrWord", "Word", "WordPart", "ParameterExpansion", "Parameter"},
		},
		{ // 40
			"${a#b}",
			func(f *bash.File) bash.Type {
				return f.Lines[0].Statements[0].Pipeline.CommandOrCompound.Command.AssignmentsOrWords[0].Word.Parts[0].ParameterExpansion.BraceWord
			},
			[]string{"File", "Line", "Statement", "Pipeline", "CommandOrCompound", "Command", "AssignmentOrWord", "Word", "WordPart", "ParameterExpansion", "BraceWord"},
		},
		{ // 41
			"${a/b/c}",
			func(f *bash.File) bash.Type {
				return f.Lines[0].Statements[0].Pipeline.CommandOrCompound.Command.AssignmentsOrWords[0].Word.Parts[0].ParameterExpansion.String
			},
			[]string{"File", "Line", "Statement", "Pipeline", "CommandOrCompound", "Command", "AssignmentOrWord", "Word", "WordPart", "ParameterExpansion", "String"},
		},
		{ // 42
			"${a[b c]}",
			nilRet,
			nil,
		},
		{ // 43
			"${a[b c]}",
			func(f *bash.File) bash.Type {
				return &f.Lines[0].Statements[0].Pipeline.CommandOrCompound.Command.AssignmentsOrWords[0].Word.Parts[0].ParameterExpansion.Parameter.Array[0]
			},
			[]string{"File", "Line", "Statement", "Pipeline", "CommandOrCompound", "Command", "AssignmentOrWord", "Word", "WordPart", "ParameterExpansion", "Parameter", "WordOrOperator"},
		},
		{ // 44
			"${a[b c]}",
			func(f *bash.File) bash.Type {
				return &f.Lines[0].Statements[0].Pipeline.CommandOrCompound.Command.AssignmentsOrWords[0].Word.Parts[0].ParameterExpansion.Parameter.Array[1]
			},
			[]string{"File", "Line", "Statement", "Pipeline", "CommandOrCompound", "Command", "AssignmentOrWord", "Word", "WordPart", "ParameterExpansion", "Parameter", "WordOrOperator"},
		},
		{ // 45
			"${a:-b c}",
			nilRet,
			nil,
		},
		{ // 46
			"${a:-b c}",
			func(f *bash.File) bash.Type {
				return &f.Lines[0].Statements[0].Pipeline.CommandOrCompound.Command.AssignmentsOrWords[0].Word.Parts[0].ParameterExpansion.BraceWord.Parts[0]
			},
			[]string{"File", "Line", "Statement", "Pipeline", "CommandOrCompound", "Command", "AssignmentOrWord", "Word", "WordPart", "ParameterExpansion", "BraceWord", "WordPart"},
		},
		{ // 47
			"${a:-b c}",
			func(f *bash.File) bash.Type {
				return &f.Lines[0].Statements[0].Pipeline.CommandOrCompound.Command.AssignmentsOrWords[0].Word.Parts[0].ParameterExpansion.BraceWord.Parts[1]
			},
			[]string{"File", "Line", "Statement", "Pipeline", "CommandOrCompound", "Command", "AssignmentOrWord", "Word", "WordPart", "ParameterExpansion", "BraceWord", "WordPart"},
		},
		{ // 48
			"${a/b/c d}",
			nilRet,
			nil,
		},
		{ // 49
			"${a/b/c d}",
			func(f *bash.File) bash.Type {
				return &f.Lines[0].Statements[0].Pipeline.CommandOrCompound.Command.AssignmentsOrWords[0].Word.Parts[0].ParameterExpansion.String.WordsOrTokens[0]
			},
			[]string{"File", "Line", "Statement", "Pipeline", "CommandOrCompound", "Command", "AssignmentOrWord", "Word", "WordPart", "ParameterExpansion", "String", "WordOrToken"},
		},
		{ // 50
			"${a/b/c d}",
			func(f *bash.File) bash.Type {
				return &f.Lines[0].Statements[0].Pipeline.CommandOrCompound.Command.AssignmentsOrWords[0].Word.Parts[0].ParameterExpansion.String.WordsOrTokens[1]
			},
			[]string{"File", "Line", "Statement", "Pipeline", "CommandOrCompound", "Command", "AssignmentOrWord", "Word", "WordPart", "ParameterExpansion", "String", "WordOrToken"},
		},
		{ // 51
			"$((a + b))",
			nilRet,
			nil,
		},
		{ // 52
			"$((a + b))",
			func(f *bash.File) bash.Type {
				return &f.Lines[0].Statements[0].Pipeline.CommandOrCompound.Command.AssignmentsOrWords[0].Word.Parts[0].ArithmeticExpansion.WordsAndOperators[0]
			},
			[]string{"File", "Line", "Statement", "Pipeline", "CommandOrCompound", "Command", "AssignmentOrWord", "Word", "WordPart", "ArithmeticExpansion", "WordOrOperator"},
		},
		{ // 53
			"$((a + b))",
			func(f *bash.File) bash.Type {
				return &f.Lines[0].Statements[0].Pipeline.CommandOrCompound.Command.AssignmentsOrWords[0].Word.Parts[0].ArithmeticExpansion.WordsAndOperators[1]
			},
			[]string{"File", "Line", "Statement", "Pipeline", "CommandOrCompound", "Command", "AssignmentOrWord", "Word", "WordPart", "ArithmeticExpansion", "WordOrOperator"},
		},
		{ // 54
			"$((a))",
			func(f *bash.File) bash.Type {
				return f.Lines[0].Statements[0].Pipeline.CommandOrCompound.Command.AssignmentsOrWords[0].Word.Parts[0].ArithmeticExpansion.WordsAndOperators[0].Word
			},
			[]string{"File", "Line", "Statement", "Pipeline", "CommandOrCompound", "Command", "AssignmentOrWord", "Word", "WordPart", "ArithmeticExpansion", "WordOrOperator", "Word"},
		},
		{ // 55
			"a[b c]=",
			func(f *bash.File) bash.Type {
				return &f.Lines[0].Statements[0].Pipeline.CommandOrCompound.Command.Vars[0].Identifier.Subscript[0]
			},
			[]string{"File", "Line", "Statement", "Pipeline", "CommandOrCompound", "Command", "Assignment", "ParameterAssign", "WordOrOperator"},
		},
		{ // 56
			"a[b c]=",
			func(f *bash.File) bash.Type {
				return &f.Lines[0].Statements[0].Pipeline.CommandOrCompound.Command.Vars[0].Identifier.Subscript[1]
			},
			[]string{"File", "Line", "Statement", "Pipeline", "CommandOrCompound", "Command", "Assignment", "ParameterAssign", "WordOrOperator"},
		},
		{ // 57
			"a=(b)",
			nilRet,
			nil,
		},
		{ // 58
			"a=b",
			func(f *bash.File) bash.Type {
				return f.Lines[0].Statements[0].Pipeline.CommandOrCompound.Command.Vars[0].Value.Word
			},
			[]string{"File", "Line", "Statement", "Pipeline", "CommandOrCompound", "Command", "Assignment", "Value", "Word"},
		},
		{ // 59
			"a=(b c)",
			func(f *bash.File) bash.Type {
				return &f.Lines[0].Statements[0].Pipeline.CommandOrCompound.Command.Vars[0].Value.Array[0]
			},
			[]string{"File", "Line", "Statement", "Pipeline", "CommandOrCompound", "Command", "Assignment", "Value", "ArrayWord"},
		},
		{ // 60
			"a=(b c)",
			func(f *bash.File) bash.Type {
				return &f.Lines[0].Statements[0].Pipeline.CommandOrCompound.Command.Vars[0].Value.Array[1]
			},
			[]string{"File", "Line", "Statement", "Pipeline", "CommandOrCompound", "Command", "Assignment", "Value", "ArrayWord"},
		},
		{ // 61
			"a=(b)",
			func(f *bash.File) bash.Type {
				return &f.Lines[0].Statements[0].Pipeline.CommandOrCompound.Command.Vars[0].Value.Array[0].Word
			},
			[]string{"File", "Line", "Statement", "Pipeline", "CommandOrCompound", "Command", "Assignment", "Value", "ArrayWord", "Word"},
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
