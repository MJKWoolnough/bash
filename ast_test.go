package bash

import (
	"errors"
	"reflect"
	"testing"

	"vimagination.zapto.org/parser"
)

type sourceFn struct {
	Source string
	Fn     func(*test, Tokens)
}

type test struct {
	Parser *bashParser
	Output Type
	Err    error
}

func makeTokeniser(tk parser.Tokeniser) *parser.Tokeniser {
	return &tk
}

func doTests(t *testing.T, tests []sourceFn, fn func(*test) (Type, error)) {
	t.Helper()

	var err error

	for n, tt := range tests {
		var ts test

		if ts.Parser, err = newBashParser(makeTokeniser(parser.NewStringTokeniser(tt.Source))); err != nil {
			t.Errorf("test %d: unexpected error: %s", n+1, err)

			continue
		}

		tt.Fn(&ts, Tokens(ts.Parser.Tokens[:cap(ts.Parser.Tokens)]))

		if output, err := fn(&ts); !errors.Is(err, ts.Err) {
			t.Errorf("test %d: expecting error: %v, got %v", n+1, ts.Err, err)
		} else if ts.Output != nil && !reflect.DeepEqual(output, ts.Output) {
			t.Errorf("test %d: expecting \n%+v\n...got...\n%+v", n+1, ts.Output, output)
		}
	}
}

func TestWord(t *testing.T) {
	doTests(t, []sourceFn{
		{"a", func(t *test, tk Tokens) { // 1
			t.Output = Word{
				Parts: []WordPart{
					{
						Part:   &tk[0],
						Tokens: tk[:1],
					},
				},
				Tokens: tk[:1],
			}
		}},
		{"a b", func(t *test, tk Tokens) { // 2
			t.Output = Word{
				Parts: []WordPart{
					{
						Part:   &tk[0],
						Tokens: tk[:1],
					},
				},
				Tokens: tk[:1],
			}
		}},
		{"a\nb", func(t *test, tk Tokens) { // 3
			t.Output = Word{
				Parts: []WordPart{
					{
						Part:   &tk[0],
						Tokens: tk[:1],
					},
				},
				Tokens: tk[:1],
			}
		}},
		{"a|b", func(t *test, tk Tokens) { // 4
			t.Output = Word{
				Parts: []WordPart{
					{
						Part:   &tk[0],
						Tokens: tk[:1],
					},
				},
				Tokens: tk[:1],
			}
		}},
		{"a$b", func(t *test, tk Tokens) { // 5
			t.Output = Word{
				Parts: []WordPart{
					{
						Part:   &tk[0],
						Tokens: tk[:1],
					},
					{
						Part:   &tk[1],
						Tokens: tk[1:2],
					},
				},
				Tokens: tk[:2],
			}
		}},
		{"a$()", func(t *test, tk Tokens) { // 6
			t.Output = Word{
				Parts: []WordPart{
					{
						Part:   &tk[0],
						Tokens: tk[:1],
					},
					{
						CommandSubstitution: &CommandSubstitution{
							Command: File{
								Tokens: tk[2:2],
							},
							Tokens: tk[1:3],
						},
						Tokens: tk[1:3],
					},
				},
				Tokens: tk[:3],
			}
		}},
		{"${a}b", func(t *test, tk Tokens) { // 7
			t.Output = Word{
				Parts: []WordPart{
					{
						ParameterExpansion: &ParameterExpansion{
							Parameter: Parameter{
								Parameter: &tk[1],
								Tokens:    tk[1:2],
							},
							Tokens: tk[:3],
						},
						Tokens: tk[:3],
					},
					{
						Part:   &tk[3],
						Tokens: tk[3:4],
					},
				},
				Tokens: tk[:4],
			}
		}},
		{"a$(())", func(t *test, tk Tokens) { // 8
			t.Output = Word{
				Parts: []WordPart{
					{
						Part:   &tk[0],
						Tokens: tk[:1],
					},
					{
						ArithmeticExpansion: &ArithmeticExpansion{
							Tokens: tk[1:3],
						},
						Tokens: tk[1:3],
					},
				},
				Tokens: tk[:3],
			}
		}},
		{"$(||)", func(t *test, tk Tokens) { // 9
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err: Error{
									Err: Error{
										Err:     ErrMissingWord,
										Parsing: "Command",
										Token:   tk[1],
									},
									Parsing: "Pipeline",
									Token:   tk[1],
								},
								Parsing: "Statement",
								Token:   tk[1],
							},
							Parsing: "File",
							Token:   tk[1],
						},
						Parsing: "CommandSubstitution",
						Token:   tk[1],
					},
					Parsing: "WordPart",
					Token:   tk[0],
				},
				Parsing: "Word",
				Token:   tk[0],
			}
		}},
	}, func(t *test) (Type, error) {
		var w Word

		err := w.parse(t.Parser)

		return w, err
	})
}

func TestWordPart(t *testing.T) {
	doTests(t, []sourceFn{
		{"a", func(t *test, tk Tokens) { // 1
			t.Output = WordPart{
				Part:   &tk[0],
				Tokens: tk[:1],
			}
		}},
		{"${a}", func(t *test, tk Tokens) { // 2
			t.Output = WordPart{
				ParameterExpansion: &ParameterExpansion{
					Parameter: Parameter{
						Parameter: &tk[1],
						Tokens:    tk[1:2],
					},
					Tokens: tk[:3],
				},
				Tokens: tk[:3],
			}
		}},
		{"$()", func(t *test, tk Tokens) { // 3
			t.Output = WordPart{
				CommandSubstitution: &CommandSubstitution{
					Command: File{
						Tokens: tk[1:1],
					},
					Tokens: tk[:2],
				},
				Tokens: tk[:2],
			}
		}},
		{"``", func(t *test, tk Tokens) { // 4
			t.Output = WordPart{
				CommandSubstitution: &CommandSubstitution{
					SubstitutionType: SubstitutionBacktick,
					Command: File{
						Tokens: tk[1:1],
					},
					Tokens: tk[:2],
				},
				Tokens: tk[:2],
			}
		}},
		{"$(())", func(t *test, tk Tokens) { // 5
			t.Output = WordPart{
				ArithmeticExpansion: &ArithmeticExpansion{
					Tokens: tk[:2],
				},
				Tokens: tk[:2],
			}
		}},
		{"${a[$(||)]}", func(t *test, tk Tokens) { // 6
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err: Error{
									Err: Error{
										Err: Error{
											Err: Error{
												Err: Error{
													Err:     ErrMissingWord,
													Parsing: "Command",
													Token:   tk[4],
												},
												Parsing: "Pipeline",
												Token:   tk[4],
											},
											Parsing: "Statement",
											Token:   tk[4],
										},
										Parsing: "File",
										Token:   tk[4],
									},
									Parsing: "CommandSubstitution",
									Token:   tk[4],
								},
								Parsing: "WordPart",
								Token:   tk[3],
							},
							Parsing: "Word",
							Token:   tk[3],
						},
						Parsing: "Parameter",
						Token:   tk[3],
					},
					Parsing: "ParameterExpansion",
					Token:   tk[1],
				},
				Parsing: "WordPart",
				Token:   tk[0],
			}
		}},
		{"$(($(||)))", func(t *test, tk Tokens) { // 7
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err: Error{
									Err: Error{
										Err: Error{
											Err: Error{
												Err: Error{
													Err:     ErrMissingWord,
													Parsing: "Command",
													Token:   tk[2],
												},
												Parsing: "Pipeline",
												Token:   tk[2],
											},
											Parsing: "Statement",
											Token:   tk[2],
										},
										Parsing: "File",
										Token:   tk[2],
									},
									Parsing: "CommandSubstitution",
									Token:   tk[2],
								},
								Parsing: "WordPart",
								Token:   tk[1],
							},
							Parsing: "Word",
							Token:   tk[1],
						},
						Parsing: "WordOrOperator",
						Token:   tk[1],
					},
					Parsing: "ArithmeticExpansion",
					Token:   tk[1],
				},
				Parsing: "WordPart",
				Token:   tk[0],
			}
		}},
		{"$(||)", func(t *test, tk Tokens) { // 8
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err: Error{
									Err:     ErrMissingWord,
									Parsing: "Command",
									Token:   tk[1],
								},
								Parsing: "Pipeline",
								Token:   tk[1],
							},
							Parsing: "Statement",
							Token:   tk[1],
						},
						Parsing: "File",
						Token:   tk[1],
					},
					Parsing: "CommandSubstitution",
					Token:   tk[1],
				},
				Parsing: "WordPart",
				Token:   tk[0],
			}
		}},
	}, func(t *test) (Type, error) {
		var wp WordPart

		err := wp.parse(t.Parser)

		return wp, err
	})
}

func TestWordOrToken(t *testing.T) {
	doTests(t, []sourceFn{
		{"a", func(t *test, tk Tokens) { // 1
			t.Output = WordOrToken{
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[0],
							Tokens: tk[:1],
						},
					},
					Tokens: tk[:1],
				},
				Tokens: tk[:1],
			}
		}},
		{" ", func(t *test, tk Tokens) { // 2
			t.Output = WordOrToken{
				Token:  &tk[0],
				Tokens: tk[:1],
			}
		}},
		{"$(||)", func(t *test, tk Tokens) { // 3
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err: Error{
									Err: Error{
										Err: Error{
											Err:     ErrMissingWord,
											Parsing: "Command",
											Token:   tk[1],
										},
										Parsing: "Pipeline",
										Token:   tk[1],
									},
									Parsing: "Statement",
									Token:   tk[1],
								},
								Parsing: "File",
								Token:   tk[1],
							},
							Parsing: "CommandSubstitution",
							Token:   tk[1],
						},
						Parsing: "WordPart",
						Token:   tk[0],
					},
					Parsing: "Word",
					Token:   tk[0],
				},
				Parsing: "WordOrToken",
				Token:   tk[0],
			}
		}},
	}, func(t *test) (Type, error) {
		var w WordOrToken

		err := w.parse(t.Parser)

		return w, err
	})
}

func TestCommandSubstitution(t *testing.T) {
	doTests(t, []sourceFn{
		{"$()", func(t *test, tk Tokens) { // 1
			t.Output = CommandSubstitution{
				SubstitutionType: SubstitutionNew,
				Command: File{
					Tokens: tk[1:1],
				},
				Tokens: tk[:2],
			}
		}},
		{"``", func(t *test, tk Tokens) { // 2
			t.Output = CommandSubstitution{
				SubstitutionType: SubstitutionBacktick,
				Command: File{
					Tokens: tk[1:1],
				},
				Tokens: tk[:2],
			}
		}},
		{"$(``)", func(t *test, tk Tokens) { // 3
			t.Output = CommandSubstitution{
				SubstitutionType: SubstitutionNew,
				Command: File{
					Statements: []Statement{
						{
							Pipeline: Pipeline{
								Command: Command{
									Words: []Word{
										{
											Parts: []WordPart{
												{
													CommandSubstitution: &CommandSubstitution{
														SubstitutionType: SubstitutionBacktick,
														Command: File{
															Tokens: tk[2:2],
														},
														Tokens: tk[1:3],
													},
													Tokens: tk[1:3],
												},
											},
											Tokens: tk[1:3],
										},
									},
									Tokens: tk[1:3],
								},
								Tokens: tk[1:3],
							},
							Tokens: tk[1:3],
						},
					},
					Tokens: tk[1:3],
				},
				Tokens: tk[:4],
			}
		}},
		{"`\\`\\``", func(t *test, tk Tokens) { // 4
			t.Output = CommandSubstitution{
				SubstitutionType: SubstitutionBacktick,
				Command: File{
					Statements: []Statement{
						{
							Pipeline: Pipeline{
								Command: Command{
									Words: []Word{
										{
											Parts: []WordPart{
												{
													CommandSubstitution: &CommandSubstitution{
														SubstitutionType: SubstitutionBacktick,
														Command: File{
															Tokens: tk[2:2],
														},
														Tokens: tk[1:3],
													},
													Tokens: tk[1:3],
												},
											},
											Tokens: tk[1:3],
										},
									},
									Tokens: tk[1:3],
								},
								Tokens: tk[1:3],
							},
							Tokens: tk[1:3],
						},
					},
					Tokens: tk[1:3],
				},
				Tokens: tk[:4],
			}
		}},
		{"$(||)", func(t *test, tk Tokens) { // 8
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err:     ErrMissingWord,
								Parsing: "Command",
								Token:   tk[1],
							},
							Parsing: "Pipeline",
							Token:   tk[1],
						},
						Parsing: "Statement",
						Token:   tk[1],
					},
					Parsing: "File",
					Token:   tk[1],
				},
				Parsing: "CommandSubstitution",
				Token:   tk[1],
			}
		}},
	}, func(t *test) (Type, error) {
		var c CommandSubstitution

		err := c.parse(t.Parser)

		return c, err
	})
}

func TestRedirection(t *testing.T) {
	doTests(t, []sourceFn{
		{">a", func(t *test, tk Tokens) { // 1
			t.Output = Redirection{
				Redirector: &tk[0],
				Output: Word{
					Parts: []WordPart{
						{
							Part:   &tk[1],
							Tokens: tk[1:2],
						},
					},
					Tokens: tk[1:2],
				},
				Tokens: tk[:2],
			}
		}},
		{"> a", func(t *test, tk Tokens) { // 2
			t.Output = Redirection{
				Redirector: &tk[0],
				Output: Word{
					Parts: []WordPart{
						{
							Part:   &tk[2],
							Tokens: tk[2:3],
						},
					},
					Tokens: tk[2:3],
				},
				Tokens: tk[:3],
			}
		}},
		{"2>&1", func(t *test, tk Tokens) { // 3
			t.Output = Redirection{
				Input:      &tk[0],
				Redirector: &tk[1],
				Output: Word{
					Parts: []WordPart{
						{
							Part:   &tk[2],
							Tokens: tk[2:3],
						},
					},
					Tokens: tk[2:3],
				},
				Tokens: tk[:3],
			}
		}},
		{"<a", func(t *test, tk Tokens) { // 4
			t.Output = Redirection{
				Redirector: &tk[0],
				Output: Word{
					Parts: []WordPart{
						{
							Part:   &tk[1],
							Tokens: tk[1:2],
						},
					},
					Tokens: tk[1:2],
				},
				Tokens: tk[:2],
			}
		}},
		{"2< a", func(t *test, tk Tokens) { // 5
			t.Output = Redirection{
				Input:      &tk[0],
				Redirector: &tk[1],
				Output: Word{
					Parts: []WordPart{
						{
							Part:   &tk[3],
							Tokens: tk[3:4],
						},
					},
					Tokens: tk[3:4],
				},
				Tokens: tk[:4],
			}
		}},
		{">|a", func(t *test, tk Tokens) { // 6
			t.Output = Redirection{
				Redirector: &tk[0],
				Output: Word{
					Parts: []WordPart{
						{
							Part:   &tk[1],
							Tokens: tk[1:2],
						},
					},
					Tokens: tk[1:2],
				},
				Tokens: tk[:2],
			}
		}},
		{"3>|a", func(t *test, tk Tokens) { // 7
			t.Output = Redirection{
				Input:      &tk[0],
				Redirector: &tk[1],
				Output: Word{
					Parts: []WordPart{
						{
							Part:   &tk[2],
							Tokens: tk[2:3],
						},
					},
					Tokens: tk[2:3],
				},
				Tokens: tk[:3],
			}
		}},
		{">>a", func(t *test, tk Tokens) { // 8
			t.Output = Redirection{
				Redirector: &tk[0],
				Output: Word{
					Parts: []WordPart{
						{
							Part:   &tk[1],
							Tokens: tk[1:2],
						},
					},
					Tokens: tk[1:2],
				},
				Tokens: tk[:2],
			}
		}},
		{"1>>a", func(t *test, tk Tokens) { // 9
			t.Output = Redirection{
				Input:      &tk[0],
				Redirector: &tk[1],
				Output: Word{
					Parts: []WordPart{
						{
							Part:   &tk[2],
							Tokens: tk[2:3],
						},
					},
					Tokens: tk[2:3],
				},
				Tokens: tk[:3],
			}
		}},
		{"&>a", func(t *test, tk Tokens) { // 10
			t.Output = Redirection{
				Redirector: &tk[0],
				Output: Word{
					Parts: []WordPart{
						{
							Part:   &tk[1],
							Tokens: tk[1:2],
						},
					},
					Tokens: tk[1:2],
				},
				Tokens: tk[:2],
			}
		}},
		{">&a", func(t *test, tk Tokens) { // 11
			t.Output = Redirection{
				Redirector: &tk[0],
				Output: Word{
					Parts: []WordPart{
						{
							Part:   &tk[1],
							Tokens: tk[1:2],
						},
					},
					Tokens: tk[1:2],
				},
				Tokens: tk[:2],
			}
		}},
		{"&>>a", func(t *test, tk Tokens) { // 12
			t.Output = Redirection{
				Redirector: &tk[0],
				Output: Word{
					Parts: []WordPart{
						{
							Part:   &tk[1],
							Tokens: tk[1:2],
						},
					},
					Tokens: tk[1:2],
				},
				Tokens: tk[:2],
			}
		}},
		{"<<abc\nabc", func(t *test, tk Tokens) { // 13
			t.Output = Redirection{
				Redirector: &tk[0],
				Output: Word{
					Parts: []WordPart{
						{
							Part:   &tk[1],
							Tokens: tk[1:2],
						},
					},
					Tokens: tk[1:2],
				},
				Tokens: tk[:2],
			}
		}},
		{"2<<-abc\nabc", func(t *test, tk Tokens) { // 14
			t.Output = Redirection{
				Input:      &tk[0],
				Redirector: &tk[1],
				Output: Word{
					Parts: []WordPart{
						{
							Part:   &tk[2],
							Tokens: tk[2:3],
						},
					},
					Tokens: tk[2:3],
				},
				Tokens: tk[:3],
			}
		}},
		{"2<&3", func(t *test, tk Tokens) { // 15
			t.Output = Redirection{
				Input:      &tk[0],
				Redirector: &tk[1],
				Output: Word{
					Parts: []WordPart{
						{
							Part:   &tk[2],
							Tokens: tk[2:3],
						},
					},
					Tokens: tk[2:3],
				},
				Tokens: tk[:3],
			}
		}},
		{"<<<abc", func(t *test, tk Tokens) { // 16
			t.Output = Redirection{
				Redirector: &tk[0],
				Output: Word{
					Parts: []WordPart{
						{
							Part:   &tk[1],
							Tokens: tk[1:2],
						},
					},
					Tokens: tk[1:2],
				},
				Tokens: tk[:2],
			}
		}},
		{"2<&3-", func(t *test, tk Tokens) { // 17
			t.Output = Redirection{
				Input:      &tk[0],
				Redirector: &tk[1],
				Output: Word{
					Parts: []WordPart{
						{
							Part:   &tk[2],
							Tokens: tk[2:3],
						},
					},
					Tokens: tk[2:3],
				},
				Tokens: tk[:3],
			}
		}},
		{">&2-", func(t *test, tk Tokens) { // 18
			t.Output = Redirection{
				Redirector: &tk[0],
				Output: Word{
					Parts: []WordPart{
						{
							Part:   &tk[1],
							Tokens: tk[1:2],
						},
					},
					Tokens: tk[1:2],
				},
				Tokens: tk[:2],
			}
		}},
		{"<>abc", func(t *test, tk Tokens) { // 19
			t.Output = Redirection{
				Redirector: &tk[0],
				Output: Word{
					Parts: []WordPart{
						{
							Part:   &tk[1],
							Tokens: tk[1:2],
						},
					},
					Tokens: tk[1:2],
				},
				Tokens: tk[:2],
			}
		}},
		{">$(||)", func(t *test, tk Tokens) { // 20
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err: Error{
									Err: Error{
										Err:     ErrMissingWord,
										Parsing: "Command",
										Token:   tk[2],
									},
									Parsing: "Pipeline",
									Token:   tk[2],
								},
								Parsing: "Statement",
								Token:   tk[2],
							},
							Parsing: "File",
							Token:   tk[2],
						},
						Parsing: "CommandSubstitution",
						Token:   tk[2],
					},
					Parsing: "WordPart",
					Token:   tk[1],
				},
				Parsing: "Word",
				Token:   tk[1],
			}
		}},
	}, func(t *test) (Type, error) {
		var r Redirection

		err := r.parse(t.Parser)

		return r, err
	})
}

func TestArithmeticExpansion(t *testing.T) {
	doTests(t, []sourceFn{
		{"$((a))", func(t *test, tk Tokens) { // 1
			t.Output = ArithmeticExpansion{
				WordsAndOperators: []WordOrOperator{
					{
						Word: &Word{
							Parts: []WordPart{
								{
									Part:   &tk[1],
									Tokens: tk[1:2],
								},
							},
							Tokens: tk[1:2],
						},
						Tokens: tk[1:2],
					},
				},
				Tokens: tk[:3],
			}
		}},
		{"$(( a ))", func(t *test, tk Tokens) { // 2
			t.Output = ArithmeticExpansion{
				WordsAndOperators: []WordOrOperator{
					{
						Word: &Word{
							Parts: []WordPart{
								{
									Part:   &tk[2],
									Tokens: tk[2:3],
								},
							},
							Tokens: tk[2:3],
						},
						Tokens: tk[2:3],
					},
				},
				Tokens: tk[:5],
			}
		}},
		{"$(( a$b ))", func(t *test, tk Tokens) { // 3
			t.Output = ArithmeticExpansion{
				WordsAndOperators: []WordOrOperator{
					{
						Word: &Word{
							Parts: []WordPart{
								{
									Part:   &tk[2],
									Tokens: tk[2:3],
								},
								{
									Part:   &tk[3],
									Tokens: tk[3:4],
								},
							},
							Tokens: tk[2:4],
						},
						Tokens: tk[2:4],
					},
				},
				Tokens: tk[:6],
			}
		}},
		{"$((a+b))", func(t *test, tk Tokens) { // 4
			t.Output = ArithmeticExpansion{
				WordsAndOperators: []WordOrOperator{
					{
						Word: &Word{
							Parts: []WordPart{
								{
									Part:   &tk[1],
									Tokens: tk[1:2],
								},
							},
							Tokens: tk[1:2],
						},
						Tokens: tk[1:2],
					},
					{
						Operator: &tk[2],
						Tokens:   tk[2:3],
					},
					{
						Word: &Word{
							Parts: []WordPart{
								{
									Part:   &tk[3],
									Tokens: tk[3:4],
								},
							},
							Tokens: tk[3:4],
						},
						Tokens: tk[3:4],
					},
				},
				Tokens: tk[:5],
			}
		}},
		{"$(($(||)))", func(t *test, tk Tokens) { // 5
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err: Error{
									Err: Error{
										Err: Error{
											Err: Error{
												Err:     ErrMissingWord,
												Parsing: "Command",
												Token:   tk[2],
											},
											Parsing: "Pipeline",
											Token:   tk[2],
										},
										Parsing: "Statement",
										Token:   tk[2],
									},
									Parsing: "File",
									Token:   tk[2],
								},
								Parsing: "CommandSubstitution",
								Token:   tk[2],
							},
							Parsing: "WordPart",
							Token:   tk[1],
						},
						Parsing: "Word",
						Token:   tk[1],
					},
					Parsing: "WordOrOperator",
					Token:   tk[1],
				},
				Parsing: "ArithmeticExpansion",
				Token:   tk[1],
			}
		}},
	}, func(t *test) (Type, error) {
		var ae ArithmeticExpansion

		err := ae.parse(t.Parser)

		return ae, err
	})
}

func TestWordOrOperator(t *testing.T) {
	doTests(t, []sourceFn{
		{"a", func(t *test, tk Tokens) { // 1
			t.Output = WordOrOperator{
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[0],
							Tokens: tk[:1],
						},
					},
					Tokens: tk[:1],
				},
				Tokens: tk[:1],
			}
		}},
		{"$a", func(t *test, tk Tokens) { // 2
			t.Output = WordOrOperator{
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[0],
							Tokens: tk[:1],
						},
					},
					Tokens: tk[:1],
				},
				Tokens: tk[:1],
			}
		}},
		{">", func(t *test, tk Tokens) { // 3
			t.Output = WordOrOperator{
				Operator: &tk[0],
				Tokens:   tk[:1],
			}
		}},
		{"&", func(t *test, tk Tokens) { // 4
			t.Output = WordOrOperator{
				Operator: &tk[0],
				Tokens:   tk[:1],
			}
		}},
		{"$(||)", func(t *test, tk Tokens) { // 5
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err: Error{
									Err:     ErrMissingWord,
									Parsing: "Command",
									Token:   tk[1],
								},
								Parsing: "Pipeline",
								Token:   tk[1],
							},
							Parsing: "Statement",
							Token:   tk[1],
						},
						Parsing: "File",
						Token:   tk[1],
					},
					Parsing: "CommandSubstitution",
					Token:   tk[1],
				},
				Parsing: "WordPart",
				Token:   tk[0],
			}
		}},
	}, func(t *test) (Type, error) {
		var wo WordOrOperator

		err := wo.parse(t.Parser)

		return wo, err
	})
}
