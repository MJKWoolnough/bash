package bash

import (
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

		if output, err := fn(&ts); !reflect.DeepEqual(err, ts.Err) {
			t.Errorf("test %d: expecting error: %v, got %v", n+1, ts.Err, err)
		} else if ts.Output != nil && !reflect.DeepEqual(output, ts.Output) {
			t.Errorf("test %d: expecting \n%+v\n...got...\n%+v", n+1, ts.Output, output)
		}
	}
}

func TestFile(t *testing.T) {
	doTests(t, []sourceFn{
		{"a", func(t *test, tk Tokens) { // 1
			t.Output = File{
				Lines: []Line{
					{
						Statements: []Statement{
							{
								Pipeline: Pipeline{
									CommandOrControl: CommandOrControl{
										Command: &Command{
											Words: []Word{
												{
													Parts: []WordPart{
														{
															Part:   &tk[0],
															Tokens: tk[:1],
														},
													},
													Tokens: tk[:1],
												},
											},
											Tokens: tk[:1],
										},
										Tokens: tk[:1],
									},
									Tokens: tk[:1],
								},
								Tokens: tk[:1],
							},
						},
						Tokens: tk[:1],
					},
				},
				Tokens: tk[:1],
			}
		}},
		{"a\nb", func(t *test, tk Tokens) { // 2
			t.Output = File{
				Lines: []Line{
					{
						Statements: []Statement{
							{
								Pipeline: Pipeline{
									CommandOrControl: CommandOrControl{
										Command: &Command{
											Words: []Word{
												{
													Parts: []WordPart{
														{
															Part:   &tk[0],
															Tokens: tk[:1],
														},
													},
													Tokens: tk[:1],
												},
											},
											Tokens: tk[:1],
										},
										Tokens: tk[:1],
									},
									Tokens: tk[:1],
								},
								Tokens: tk[:1],
							},
						},
						Tokens: tk[:1],
					},
					{
						Statements: []Statement{
							{
								Pipeline: Pipeline{
									CommandOrControl: CommandOrControl{
										Command: &Command{
											Words: []Word{
												{
													Parts: []WordPart{
														{
															Part:   &tk[2],
															Tokens: tk[2:3],
														},
													},
													Tokens: tk[2:3],
												},
											},
											Tokens: tk[2:3],
										},
										Tokens: tk[2:3],
									},
									Tokens: tk[2:3],
								},
								Tokens: tk[2:3],
							},
						},
						Tokens: tk[2:3],
					},
				},
				Tokens: tk[:3],
			}
		}},
		{"a\n\nb\n", func(t *test, tk Tokens) { // 3
			t.Output = File{
				Lines: []Line{
					{
						Statements: []Statement{
							{
								Pipeline: Pipeline{
									CommandOrControl: CommandOrControl{
										Command: &Command{
											Words: []Word{
												{
													Parts: []WordPart{
														{
															Part:   &tk[0],
															Tokens: tk[:1],
														},
													},
													Tokens: tk[:1],
												},
											},
											Tokens: tk[:1],
										},
										Tokens: tk[:1],
									},
									Tokens: tk[:1],
								},
								Tokens: tk[:1],
							},
						},
						Tokens: tk[:1],
					},
					{
						Statements: []Statement{
							{
								Pipeline: Pipeline{
									CommandOrControl: CommandOrControl{
										Command: &Command{
											Words: []Word{
												{
													Parts: []WordPart{
														{
															Part:   &tk[2],
															Tokens: tk[2:3],
														},
													},
													Tokens: tk[2:3],
												},
											},
											Tokens: tk[2:3],
										},
										Tokens: tk[2:3],
									},
									Tokens: tk[2:3],
								},
								Tokens: tk[2:3],
							},
						},
						Tokens: tk[2:3],
					},
				},
				Tokens: tk[:3],
			}
		}},
		{"||", func(t *test, tk Tokens) { // 4
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err: Error{
									Err:     ErrMissingWord,
									Parsing: "Command",
									Token:   tk[0],
								},
								Parsing: "CommandOrControl",
								Token:   tk[0],
							},
							Parsing: "Pipeline",
							Token:   tk[0],
						},
						Parsing: "Statement",
						Token:   tk[0],
					},
					Parsing: "Line",
					Token:   tk[0],
				},
				Parsing: "File",
				Token:   tk[0],
			}
		}},
		{"<<a\nb\na", func(t *test, tk Tokens) { // 5
			t.Output = File{
				Lines: []Line{
					{
						Statements: []Statement{
							{
								Pipeline: Pipeline{
									CommandOrControl: CommandOrControl{
										Command: &Command{
											Redirections: []Redirection{
												{
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
													Heredoc: &Heredoc{
														HeredocPartsOrWords: []HeredocPartOrWord{
															{
																HeredocPart: &tk[3],
																Tokens:      tk[3:4],
															},
														},
														Tokens: tk[3:5],
													},
													Tokens: tk[:2],
												},
											},
											Tokens: tk[:2],
										},
										Tokens: tk[:2],
									},
									Tokens: tk[:2],
								},
								Tokens: tk[:2],
							},
						},
						Tokens: tk[:5],
					},
				},
				Tokens: tk[:5],
			}
		}},
		{"<<a;<<b\nc\na\nd\nb", func(t *test, tk Tokens) { // 6
			t.Output = File{
				Lines: []Line{
					{
						Statements: []Statement{
							{
								Pipeline: Pipeline{
									CommandOrControl: CommandOrControl{
										Command: &Command{
											Redirections: []Redirection{
												{
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
													Heredoc: &Heredoc{
														HeredocPartsOrWords: []HeredocPartOrWord{
															{
																HeredocPart: &tk[6],
																Tokens:      tk[6:7],
															},
														},
														Tokens: tk[6:8],
													},
													Tokens: tk[:2],
												},
											},
											Tokens: tk[:2],
										},
										Tokens: tk[:2],
									},
									Tokens: tk[:2],
								},
								Tokens: tk[:3],
							},
							{
								Pipeline: Pipeline{
									CommandOrControl: CommandOrControl{
										Command: &Command{
											Redirections: []Redirection{
												{
													Redirector: &tk[3],
													Output: Word{
														Parts: []WordPart{
															{
																Part:   &tk[4],
																Tokens: tk[4:5],
															},
														},
														Tokens: tk[4:5],
													},
													Heredoc: &Heredoc{
														HeredocPartsOrWords: []HeredocPartOrWord{
															{
																HeredocPart: &tk[9],
																Tokens:      tk[9:10],
															},
														},
														Tokens: tk[9:11],
													},
													Tokens: tk[3:5],
												},
											},
											Tokens: tk[3:5],
										},
										Tokens: tk[3:5],
									},
									Tokens: tk[3:5],
								},
								Tokens: tk[3:5],
							},
						},
						Tokens: tk[:11],
					},
				},
				Tokens: tk[:11],
			}
		}},
		{"<<a\n$(||)\na", func(t *test, tk Tokens) { // 7
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
																				Parsing: "CommandOrControl",
																				Token:   tk[4],
																			},
																			Parsing: "Pipeline",
																			Token:   tk[4],
																		},
																		Parsing: "Statement",
																		Token:   tk[4],
																	},
																	Parsing: "Line",
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
												Parsing: "HeredocPartOrWord",
												Token:   tk[3],
											},
											Parsing: "Heredoc",
											Token:   tk[3],
										},
										Parsing: "Redirection",
										Token:   tk[3],
									},
									Parsing: "Command",
									Token:   tk[3],
								},
								Parsing: "CommandOrControl",
								Token:   tk[3],
							},
							Parsing: "Pipeline",
							Token:   tk[3],
						},
						Parsing: "Statement",
						Token:   tk[3],
					},
					Parsing: "Line",
					Token:   tk[3],
				},
				Parsing: "File",
				Token:   tk[0],
			}
		}},
	}, func(t *test) (Type, error) {
		var f File

		err := f.parse(t.Parser)

		return f, err
	})
}

func TestLine(t *testing.T) {
	doTests(t, []sourceFn{
		{"a", func(t *test, tk Tokens) { // 1
			t.Output = Line{
				Statements: []Statement{
					{
						Pipeline: Pipeline{
							CommandOrControl: CommandOrControl{
								Command: &Command{
									Words: []Word{
										{
											Parts: []WordPart{
												{
													Part:   &tk[0],
													Tokens: tk[:1],
												},
											},
											Tokens: tk[:1],
										},
									},
									Tokens: tk[:1],
								},
								Tokens: tk[:1],
							},
							Tokens: tk[:1],
						},
						Tokens: tk[:1],
					},
				},
				Tokens: tk[:1],
			}
		}},
		{"a;b", func(t *test, tk Tokens) { // 2
			t.Output = Line{
				Statements: []Statement{
					{
						Pipeline: Pipeline{
							CommandOrControl: CommandOrControl{
								Command: &Command{
									Words: []Word{
										{
											Parts: []WordPart{
												{
													Part:   &tk[0],
													Tokens: tk[:1],
												},
											},
											Tokens: tk[:1],
										},
									},
									Tokens: tk[:1],
								},
								Tokens: tk[:1],
							},
							Tokens: tk[:1],
						},
						Tokens: tk[:2],
					},
					{
						Pipeline: Pipeline{
							CommandOrControl: CommandOrControl{
								Command: &Command{
									Words: []Word{
										{
											Parts: []WordPart{
												{
													Part:   &tk[2],
													Tokens: tk[2:3],
												},
											},
											Tokens: tk[2:3],
										},
									},
									Tokens: tk[2:3],
								},
								Tokens: tk[2:3],
							},
							Tokens: tk[2:3],
						},
						Tokens: tk[2:3],
					},
				},
				Tokens: tk[:3],
			}
		}},
		{"a & b;", func(t *test, tk Tokens) { // 3
			t.Output = Line{
				Statements: []Statement{
					{
						Pipeline: Pipeline{
							CommandOrControl: CommandOrControl{
								Command: &Command{
									Words: []Word{
										{
											Parts: []WordPart{
												{
													Part:   &tk[0],
													Tokens: tk[:1],
												},
											},
											Tokens: tk[:1],
										},
									},
									Tokens: tk[:1],
								},
								Tokens: tk[:1],
							},
							Tokens: tk[:1],
						},
						JobControl: JobControlBackground,
						Tokens:     tk[:3],
					},
					{
						Pipeline: Pipeline{
							CommandOrControl: CommandOrControl{
								Command: &Command{
									Words: []Word{
										{
											Parts: []WordPart{
												{
													Part:   &tk[4],
													Tokens: tk[4:5],
												},
											},
											Tokens: tk[4:5],
										},
									},
									Tokens: tk[4:5],
								},
								Tokens: tk[4:5],
							},
							Tokens: tk[4:5],
						},
						Tokens: tk[4:6],
					},
				},
				Tokens: tk[:6],
			}
		}},
		{"||", func(t *test, tk Tokens) { // 4
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err:     ErrMissingWord,
								Parsing: "Command",
								Token:   tk[0],
							},
							Parsing: "CommandOrControl",
							Token:   tk[0],
						},
						Parsing: "Pipeline",
						Token:   tk[0],
					},
					Parsing: "Statement",
					Token:   tk[0],
				},
				Parsing: "Line",
				Token:   tk[0],
			}
		}},
	}, func(t *test) (Type, error) {
		var l Line

		err := l.parse(t.Parser)

		return l, err
	})
}

func TestStatement(t *testing.T) {
	doTests(t, []sourceFn{
		{"a", func(t *test, tk Tokens) { // 1
			t.Output = Statement{
				Pipeline: Pipeline{
					CommandOrControl: CommandOrControl{
						Command: &Command{
							Words: []Word{
								{
									Parts: []WordPart{
										{
											Part:   &tk[0],
											Tokens: tk[:1],
										},
									},
									Tokens: tk[:1],
								},
							},
							Tokens: tk[:1],
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				Tokens: tk[:1],
			}
		}},
		{"a||b", func(t *test, tk Tokens) { // 2
			t.Output = Statement{
				Pipeline: Pipeline{
					CommandOrControl: CommandOrControl{
						Command: &Command{
							Words: []Word{
								{
									Parts: []WordPart{
										{
											Part:   &tk[0],
											Tokens: tk[:1],
										},
									},
									Tokens: tk[:1],
								},
							},
							Tokens: tk[:1],
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				LogicalOperator: LogicalOperatorOr,
				Statement: &Statement{
					Pipeline: Pipeline{
						CommandOrControl: CommandOrControl{
							Command: &Command{
								Words: []Word{
									{
										Parts: []WordPart{
											{
												Part:   &tk[2],
												Tokens: tk[2:3],
											},
										},
										Tokens: tk[2:3],
									},
								},
								Tokens: tk[2:3],
							},
							Tokens: tk[2:3],
						},
						Tokens: tk[2:3],
					},
					Tokens: tk[2:3],
				},
				Tokens: tk[:3],
			}
		}},
		{"a && b", func(t *test, tk Tokens) { // 3
			t.Output = Statement{
				Pipeline: Pipeline{
					CommandOrControl: CommandOrControl{
						Command: &Command{
							Words: []Word{
								{
									Parts: []WordPart{
										{
											Part:   &tk[0],
											Tokens: tk[:1],
										},
									},
									Tokens: tk[:1],
								},
							},
							Tokens: tk[:1],
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				LogicalOperator: LogicalOperatorAnd,
				Statement: &Statement{
					Pipeline: Pipeline{
						CommandOrControl: CommandOrControl{
							Command: &Command{
								Words: []Word{
									{
										Parts: []WordPart{
											{
												Part:   &tk[4],
												Tokens: tk[4:5],
											},
										},
										Tokens: tk[4:5],
									},
								},
								Tokens: tk[4:5],
							},
							Tokens: tk[4:5],
						},
						Tokens: tk[4:5],
					},
					Tokens: tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
		{"a||b;", func(t *test, tk Tokens) { // 4
			t.Output = Statement{
				Pipeline: Pipeline{
					CommandOrControl: CommandOrControl{
						Command: &Command{
							Words: []Word{
								{
									Parts: []WordPart{
										{
											Part:   &tk[0],
											Tokens: tk[:1],
										},
									},
									Tokens: tk[:1],
								},
							},
							Tokens: tk[:1],
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				LogicalOperator: LogicalOperatorOr,
				Statement: &Statement{
					Pipeline: Pipeline{
						CommandOrControl: CommandOrControl{
							Command: &Command{
								Words: []Word{
									{
										Parts: []WordPart{
											{
												Part:   &tk[2],
												Tokens: tk[2:3],
											},
										},
										Tokens: tk[2:3],
									},
								},
								Tokens: tk[2:3],
							},
							Tokens: tk[2:3],
						},
						Tokens: tk[2:3],
					},
					Tokens: tk[2:3],
				},
				Tokens: tk[:4],
			}
		}},
		{"a||b &", func(t *test, tk Tokens) { // 5
			t.Output = Statement{
				Pipeline: Pipeline{
					CommandOrControl: CommandOrControl{
						Command: &Command{
							Words: []Word{
								{
									Parts: []WordPart{
										{
											Part:   &tk[0],
											Tokens: tk[:1],
										},
									},
									Tokens: tk[:1],
								},
							},
							Tokens: tk[:1],
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				LogicalOperator: LogicalOperatorOr,
				Statement: &Statement{
					Pipeline: Pipeline{
						CommandOrControl: CommandOrControl{
							Command: &Command{
								Words: []Word{
									{
										Parts: []WordPart{
											{
												Part:   &tk[2],
												Tokens: tk[2:3],
											},
										},
										Tokens: tk[2:3],
									},
								},
								Tokens: tk[2:3],
							},
							Tokens: tk[2:3],
						},
						Tokens: tk[2:3],
					},
					Tokens: tk[2:3],
				},
				JobControl: JobControlBackground,
				Tokens:     tk[:5],
			}
		}},
		{"||", func(t *test, tk Tokens) { // 6
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err:     ErrMissingWord,
							Parsing: "Command",
							Token:   tk[0],
						},
						Parsing: "CommandOrControl",
						Token:   tk[0],
					},
					Parsing: "Pipeline",
					Token:   tk[0],
				},
				Parsing: "Statement",
				Token:   tk[0],
			}
		}},
		{"a || ||", func(t *test, tk Tokens) { // 7
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err:     ErrMissingWord,
								Parsing: "Command",
								Token:   tk[4],
							},
							Parsing: "CommandOrControl",
							Token:   tk[4],
						},
						Parsing: "Pipeline",
						Token:   tk[4],
					},
					Parsing: "Statement",
					Token:   tk[4],
				},
				Parsing: "Statement",
				Token:   tk[4],
			}
		}},
	}, func(t *test) (Type, error) {
		var s Statement

		err := s.parse(t.Parser, true)

		return s, err
	})
}

func TestPipeline(t *testing.T) {
	doTests(t, []sourceFn{
		{"a", func(t *test, tk Tokens) { // 1
			t.Output = Pipeline{
				CommandOrControl: CommandOrControl{
					Command: &Command{
						Words: []Word{
							{
								Parts: []WordPart{
									{
										Part:   &tk[0],
										Tokens: tk[:1],
									},
								},
								Tokens: tk[:1],
							},
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				Tokens: tk[:1],
			}
		}},
		{"time a", func(t *test, tk Tokens) { // 2
			t.Output = Pipeline{
				PipelineTime: PipelineTimeBash,
				CommandOrControl: CommandOrControl{
					Command: &Command{
						Words: []Word{
							{
								Parts: []WordPart{
									{
										Part:   &tk[2],
										Tokens: tk[2:3],
									},
								},
								Tokens: tk[2:3],
							},
						},
						Tokens: tk[2:3],
					},
					Tokens: tk[2:3],
				},
				Tokens: tk[:3],
			}
		}},
		{"time -p a", func(t *test, tk Tokens) { // 3
			t.Output = Pipeline{
				PipelineTime: PipelineTimePosix,
				CommandOrControl: CommandOrControl{
					Command: &Command{
						Words: []Word{
							{
								Parts: []WordPart{
									{
										Part:   &tk[4],
										Tokens: tk[4:5],
									},
								},
								Tokens: tk[4:5],
							},
						},
						Tokens: tk[4:5],
					},
					Tokens: tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
		{"! a", func(t *test, tk Tokens) { // 4
			t.Output = Pipeline{
				Not: true,
				CommandOrControl: CommandOrControl{
					Command: &Command{
						Words: []Word{
							{
								Parts: []WordPart{
									{
										Part:   &tk[2],
										Tokens: tk[2:3],
									},
								},
								Tokens: tk[2:3],
							},
						},
						Tokens: tk[2:3],
					},
					Tokens: tk[2:3],
				},
				Tokens: tk[:3],
			}
		}},
		{"a|b", func(t *test, tk Tokens) { // 5
			t.Output = Pipeline{
				CommandOrControl: CommandOrControl{
					Command: &Command{
						Words: []Word{
							{
								Parts: []WordPart{
									{
										Part:   &tk[0],
										Tokens: tk[:1],
									},
								},
								Tokens: tk[:1],
							},
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				Pipeline: &Pipeline{
					CommandOrControl: CommandOrControl{
						Command: &Command{
							Words: []Word{
								{
									Parts: []WordPart{
										{
											Part:   &tk[2],
											Tokens: tk[2:3],
										},
									},
									Tokens: tk[2:3],
								},
							},
							Tokens: tk[2:3],
						},
						Tokens: tk[2:3],
					},
					Tokens: tk[2:3],
				},
				Tokens: tk[:3],
			}
		}},
		{"a | b", func(t *test, tk Tokens) { // 6
			t.Output = Pipeline{
				CommandOrControl: CommandOrControl{
					Command: &Command{
						Words: []Word{
							{
								Parts: []WordPart{
									{
										Part:   &tk[0],
										Tokens: tk[:1],
									},
								},
								Tokens: tk[:1],
							},
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				Pipeline: &Pipeline{
					CommandOrControl: CommandOrControl{
						Command: &Command{
							Words: []Word{
								{
									Parts: []WordPart{
										{
											Part:   &tk[4],
											Tokens: tk[4:5],
										},
									},
									Tokens: tk[4:5],
								},
							},
							Tokens: tk[4:5],
						},
						Tokens: tk[4:5],
					},
					Tokens: tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
		{"||", func(t *test, tk Tokens) { // 7
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrMissingWord,
						Parsing: "Command",
						Token:   tk[0],
					},
					Parsing: "CommandOrControl",
					Token:   tk[0],
				},
				Parsing: "Pipeline",
				Token:   tk[0],
			}
		}},
		{"a | ||", func(t *test, tk Tokens) { // 8
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err:     ErrMissingWord,
							Parsing: "Command",
							Token:   tk[4],
						},
						Parsing: "CommandOrControl",
						Token:   tk[4],
					},
					Parsing: "Pipeline",
					Token:   tk[4],
				},
				Parsing: "Pipeline",
				Token:   tk[4],
			}
		}},
	}, func(t *test) (Type, error) {
		var p Pipeline

		err := p.parse(t.Parser, true)

		return p, err
	})
}

func TestCommand(t *testing.T) {
	doTests(t, []sourceFn{
		{"a", func(t *test, tk Tokens) { // 1
			t.Output = Command{
				Words: []Word{
					{
						Parts: []WordPart{
							{
								Part:   &tk[0],
								Tokens: tk[:1],
							},
						},
						Tokens: tk[:1],
					},
				},
				Tokens: tk[:1],
			}
		}},
		{"a=", func(t *test, tk Tokens) { // 2
			t.Output = Command{
				Vars: []Assignment{
					{
						Identifier: ParameterAssign{
							Identifier: &tk[0],
							Tokens:     tk[:1],
						},
						Value: Value{
							Word: &Word{
								Tokens: tk[2:2],
							},
							Tokens: tk[2:2],
						},
						Tokens: tk[:2],
					},
				},
				Tokens: tk[:2],
			}
		}},
		{"a=b c", func(t *test, tk Tokens) { // 3
			t.Output = Command{
				Vars: []Assignment{
					{
						Identifier: ParameterAssign{
							Identifier: &tk[0],
							Tokens:     tk[:1],
						},
						Value: Value{
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
						Tokens: tk[:3],
					},
				},
				Words: []Word{
					{
						Parts: []WordPart{
							{
								Part:   &tk[4],
								Tokens: tk[4:5],
							},
						},
						Tokens: tk[4:5],
					},
				},
				Tokens: tk[:5],
			}
		}},
		{"a=b >c d", func(t *test, tk Tokens) { // 4
			t.Output = Command{
				Vars: []Assignment{
					{
						Identifier: ParameterAssign{
							Identifier: &tk[0],
							Tokens:     tk[:1],
						},
						Value: Value{
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
						Tokens: tk[:3],
					},
				},
				Redirections: []Redirection{
					{
						Redirector: &tk[4],
						Output: Word{
							Parts: []WordPart{
								{
									Part:   &tk[5],
									Tokens: tk[5:6],
								},
							},
							Tokens: tk[5:6],
						},
						Tokens: tk[4:6],
					},
				},
				Words: []Word{
					{
						Parts: []WordPart{
							{
								Part:   &tk[7],
								Tokens: tk[7:8],
							},
						},
						Tokens: tk[7:8],
					},
				},
				Tokens: tk[:8],
			}
		}},
		{"a=b c=d >e f g=h 2>i", func(t *test, tk Tokens) { // 5
			t.Output = Command{
				Vars: []Assignment{
					{
						Identifier: ParameterAssign{
							Identifier: &tk[0],
							Tokens:     tk[:1],
						},
						Value: Value{
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
						Tokens: tk[:3],
					},
					{
						Identifier: ParameterAssign{
							Identifier: &tk[4],
							Tokens:     tk[4:5],
						},
						Value: Value{
							Word: &Word{
								Parts: []WordPart{
									{
										Part:   &tk[6],
										Tokens: tk[6:7],
									},
								},
								Tokens: tk[6:7],
							},
							Tokens: tk[6:7],
						},
						Tokens: tk[4:7],
					},
				},
				Redirections: []Redirection{
					{
						Redirector: &tk[8],
						Output: Word{
							Parts: []WordPart{
								{
									Part:   &tk[9],
									Tokens: tk[9:10],
								},
							},
							Tokens: tk[9:10],
						},
						Tokens: tk[8:10],
					},
					{
						Input:      &tk[17],
						Redirector: &tk[18],
						Output: Word{
							Parts: []WordPart{
								{
									Part:   &tk[19],
									Tokens: tk[19:20],
								},
							},
							Tokens: tk[19:20],
						},
						Tokens: tk[17:20],
					},
				},
				Words: []Word{
					{
						Parts: []WordPart{
							{
								Part:   &tk[11],
								Tokens: tk[11:12],
							},
						},
						Tokens: tk[11:12],
					},
					{
						Parts: []WordPart{
							{
								Part:   &tk[13],
								Tokens: tk[13:14],
							},
							{
								Part:   &tk[14],
								Tokens: tk[14:15],
							},
							{
								Part:   &tk[15],
								Tokens: tk[15:16],
							},
						},
						Tokens: tk[13:16],
					},
				},
				Tokens: tk[:20],
			}
		}},
		{"a[$(||)]=", func(t *test, tk Tokens) { // 6
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
													Err: Error{
														Err: Error{
															Err:     ErrMissingWord,
															Parsing: "Command",
															Token:   tk[3],
														},
														Parsing: "CommandOrControl",
														Token:   tk[3],
													},
													Parsing: "Pipeline",
													Token:   tk[3],
												},
												Parsing: "Statement",
												Token:   tk[3],
											},
											Parsing: "Line",
											Token:   tk[3],
										},
										Parsing: "File",
										Token:   tk[3],
									},
									Parsing: "CommandSubstitution",
									Token:   tk[3],
								},
								Parsing: "WordPart",
								Token:   tk[2],
							},
							Parsing: "Word",
							Token:   tk[2],
						},
						Parsing: "ParameterAssign",
						Token:   tk[2],
					},
					Parsing: "Assignment",
					Token:   tk[0],
				},
				Parsing: "Command",
				Token:   tk[0],
			}
		}},
		{">$(||)", func(t *test, tk Tokens) { // 7
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
													Err: Error{
														Err:     ErrMissingWord,
														Parsing: "Command",
														Token:   tk[2],
													},
													Parsing: "CommandOrControl",
													Token:   tk[2],
												},
												Parsing: "Pipeline",
												Token:   tk[2],
											},
											Parsing: "Statement",
											Token:   tk[2],
										},
										Parsing: "Line",
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
					Parsing: "Redirection",
					Token:   tk[1],
				},
				Parsing: "Command",
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
									Err: Error{
										Err: Error{
											Err: Error{
												Err: Error{
													Err:     ErrMissingWord,
													Parsing: "Command",
													Token:   tk[1],
												},
												Parsing: "CommandOrControl",
												Token:   tk[1],
											},
											Parsing: "Pipeline",
											Token:   tk[1],
										},
										Parsing: "Statement",
										Token:   tk[1],
									},
									Parsing: "Line",
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
				Parsing: "Command",
				Token:   tk[0],
			}
		}},
		{"a >$(||)", func(t *test, tk Tokens) { // 9
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
													Err: Error{
														Err:     ErrMissingWord,
														Parsing: "Command",
														Token:   tk[4],
													},
													Parsing: "CommandOrControl",
													Token:   tk[4],
												},
												Parsing: "Pipeline",
												Token:   tk[4],
											},
											Parsing: "Statement",
											Token:   tk[4],
										},
										Parsing: "Line",
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
					Parsing: "Redirection",
					Token:   tk[3],
				},
				Parsing: "Command",
				Token:   tk[2],
			}
		}},
	}, func(t *test) (Type, error) {
		var c Command

		err := c.parse(t.Parser, false)

		return c, err
	})
}

func TestAssignment(t *testing.T) {
	doTests(t, []sourceFn{
		{"a=", func(t *test, tk Tokens) { // 1
			t.Output = Assignment{
				Identifier: ParameterAssign{
					Identifier: &tk[0],
					Tokens:     tk[:1],
				},
				Value: Value{
					Word: &Word{
						Tokens: tk[2:2],
					},
					Tokens: tk[2:2],
				},
				Tokens: tk[:2],
			}
		}},
		{"a+=b", func(t *test, tk Tokens) { // 2
			t.Output = Assignment{
				Identifier: ParameterAssign{
					Identifier: &tk[0],
					Tokens:     tk[:1],
				},
				Assignment: AssignmentAppend,
				Value: Value{
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
				Tokens: tk[:3],
			}
		}},
		{"a[$(||)]=", func(t *test, tk Tokens) { // 3
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
													Err: Error{
														Err:     ErrMissingWord,
														Parsing: "Command",
														Token:   tk[3],
													},
													Parsing: "CommandOrControl",
													Token:   tk[3],
												},
												Parsing: "Pipeline",
												Token:   tk[3],
											},
											Parsing: "Statement",
											Token:   tk[3],
										},
										Parsing: "Line",
										Token:   tk[3],
									},
									Parsing: "File",
									Token:   tk[3],
								},
								Parsing: "CommandSubstitution",
								Token:   tk[3],
							},
							Parsing: "WordPart",
							Token:   tk[2],
						},
						Parsing: "Word",
						Token:   tk[2],
					},
					Parsing: "ParameterAssign",
					Token:   tk[2],
				},
				Parsing: "Assignment",
				Token:   tk[0],
			}
		}},
		{"a=$(||)", func(t *test, tk Tokens) { // 4
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
													Err: Error{
														Err:     ErrMissingWord,
														Parsing: "Command",
														Token:   tk[3],
													},
													Parsing: "CommandOrControl",
													Token:   tk[3],
												},
												Parsing: "Pipeline",
												Token:   tk[3],
											},
											Parsing: "Statement",
											Token:   tk[3],
										},
										Parsing: "Line",
										Token:   tk[3],
									},
									Parsing: "File",
									Token:   tk[3],
								},
								Parsing: "CommandSubstitution",
								Token:   tk[3],
							},
							Parsing: "WordPart",
							Token:   tk[2],
						},
						Parsing: "Word",
						Token:   tk[2],
					},
					Parsing: "Value",
					Token:   tk[2],
				},
				Parsing: "Assignment",
				Token:   tk[2],
			}
		}},
	}, func(t *test) (Type, error) {
		var a Assignment

		err := a.parse(t.Parser)

		return a, err
	})
}

func TestParameterAssign(t *testing.T) {
	doTests(t, []sourceFn{
		{"a=", func(t *test, tk Tokens) { // 1
			t.Output = ParameterAssign{
				Identifier: &tk[0],
				Tokens:     tk[:1],
			}
		}},
		{"a[0]=", func(t *test, tk Tokens) { // 2
			t.Output = ParameterAssign{
				Identifier: &tk[0],
				Subscript: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[2],
							Tokens: tk[2:3],
						},
					},
					Tokens: tk[2:3],
				},
				Tokens: tk[:4],
			}
		}},
		{"a[$a]=", func(t *test, tk Tokens) { // 3
			t.Output = ParameterAssign{
				Identifier: &tk[0],
				Subscript: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[2],
							Tokens: tk[2:3],
						},
					},
					Tokens: tk[2:3],
				},
				Tokens: tk[:4],
			}
		}},
		{"a[$(||)]=", func(t *test, tk Tokens) { // 4
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
													Token:   tk[3],
												},
												Parsing: "CommandOrControl",
												Token:   tk[3],
											},
											Parsing: "Pipeline",
											Token:   tk[3],
										},
										Parsing: "Statement",
										Token:   tk[3],
									},
									Parsing: "Line",
									Token:   tk[3],
								},
								Parsing: "File",
								Token:   tk[3],
							},
							Parsing: "CommandSubstitution",
							Token:   tk[3],
						},
						Parsing: "WordPart",
						Token:   tk[2],
					},
					Parsing: "Word",
					Token:   tk[2],
				},
				Parsing: "ParameterAssign",
				Token:   tk[2],
			}
		}},
		{"a[0 1]=", func(t *test, tk Tokens) { // 5
			t.Err = Error{
				Err:     ErrMissingClosingBracket,
				Parsing: "ParameterAssign",
				Token:   tk[4],
			}
		}},
	}, func(t *test) (Type, error) {
		var pa ParameterAssign

		err := pa.parse(t.Parser)

		return pa, err
	})
}

func TestValue(t *testing.T) {
	doTests(t, []sourceFn{
		{"a=b", func(t *test, tk Tokens) { // 1
			t.Output = Value{
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
			}
		}},
		{"a=()", func(t *test, tk Tokens) { // 2
			t.Output = Value{
				Array:  []Word{},
				Tokens: tk[2:4],
			}
		}},
		{"a=(b)", func(t *test, tk Tokens) { // 3
			t.Output = Value{
				Array: []Word{
					{
						Parts: []WordPart{
							{
								Part:   &tk[3],
								Tokens: tk[3:4],
							},
						},
						Tokens: tk[3:4],
					},
				},
				Tokens: tk[2:5],
			}
		}},
		{"a=( b )", func(t *test, tk Tokens) { // 4
			t.Output = Value{
				Array: []Word{
					{
						Parts: []WordPart{
							{
								Part:   &tk[4],
								Tokens: tk[4:5],
							},
						},
						Tokens: tk[4:5],
					},
				},
				Tokens: tk[2:7],
			}
		}},
		{"a=( b c )", func(t *test, tk Tokens) { // 5
			t.Output = Value{
				Array: []Word{
					{
						Parts: []WordPart{
							{
								Part:   &tk[4],
								Tokens: tk[4:5],
							},
						},
						Tokens: tk[4:5],
					},
					{
						Parts: []WordPart{
							{
								Part:   &tk[6],
								Tokens: tk[6:7],
							},
						},
						Tokens: tk[6:7],
					},
				},
				Tokens: tk[2:9],
			}
		}},
		{"a=$(||)", func(t *test, tk Tokens) { // 6
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
													Token:   tk[3],
												},
												Parsing: "CommandOrControl",
												Token:   tk[3],
											},
											Parsing: "Pipeline",
											Token:   tk[3],
										},
										Parsing: "Statement",
										Token:   tk[3],
									},
									Parsing: "Line",
									Token:   tk[3],
								},
								Parsing: "File",
								Token:   tk[3],
							},
							Parsing: "CommandSubstitution",
							Token:   tk[3],
						},
						Parsing: "WordPart",
						Token:   tk[2],
					},
					Parsing: "Word",
					Token:   tk[2],
				},
				Parsing: "Value",
				Token:   tk[2],
			}
		}},
		{"a=($(||))", func(t *test, tk Tokens) { // 7
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
												Parsing: "CommandOrControl",
												Token:   tk[4],
											},
											Parsing: "Pipeline",
											Token:   tk[4],
										},
										Parsing: "Statement",
										Token:   tk[4],
									},
									Parsing: "Line",
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
				Parsing: "Value",
				Token:   tk[3],
			}
		}},
	}, func(t *test) (Type, error) {
		var v Value

		t.Parser.Tokens = t.Parser.Tokens[2:2]

		err := v.parse(t.Parser)

		return v, err
	})
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
										Err: Error{
											Err: Error{
												Err:     ErrMissingWord,
												Parsing: "Command",
												Token:   tk[1],
											},
											Parsing: "CommandOrControl",
											Token:   tk[1],
										},
										Parsing: "Pipeline",
										Token:   tk[1],
									},
									Parsing: "Statement",
									Token:   tk[1],
								},
								Parsing: "Line",
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
													Err: Error{
														Err: Error{
															Err:     ErrMissingWord,
															Parsing: "Command",
															Token:   tk[4],
														},
														Parsing: "CommandOrControl",
														Token:   tk[4],
													},
													Parsing: "Pipeline",
													Token:   tk[4],
												},
												Parsing: "Statement",
												Token:   tk[4],
											},
											Parsing: "Line",
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
													Err: Error{
														Err: Error{
															Err:     ErrMissingWord,
															Parsing: "Command",
															Token:   tk[2],
														},
														Parsing: "CommandOrControl",
														Token:   tk[2],
													},
													Parsing: "Pipeline",
													Token:   tk[2],
												},
												Parsing: "Statement",
												Token:   tk[2],
											},
											Parsing: "Line",
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
									Err: Error{
										Err: Error{
											Err:     ErrMissingWord,
											Parsing: "Command",
											Token:   tk[1],
										},
										Parsing: "CommandOrControl",
										Token:   tk[1],
									},
									Parsing: "Pipeline",
									Token:   tk[1],
								},
								Parsing: "Statement",
								Token:   tk[1],
							},
							Parsing: "Line",
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

func TestParameterExpansion(t *testing.T) {
	doTests(t, []sourceFn{
		{"${a}", func(t *test, tk Tokens) { // 1
			t.Output = ParameterExpansion{
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Tokens: tk[:3],
			}
		}},
		{"${a[0]}", func(t *test, tk Tokens) { // 2
			t.Output = ParameterExpansion{
				Parameter: Parameter{
					Parameter: &tk[1],
					Array: &Word{
						Parts: []WordPart{
							{
								Part:   &tk[3],
								Tokens: tk[3:4],
							},
						},
						Tokens: tk[3:4],
					},
					Tokens: tk[1:5],
				},
				Tokens: tk[:6],
			}
		}},
		{"${@}", func(t *test, tk Tokens) { // 3
			t.Output = ParameterExpansion{
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Tokens: tk[:3],
			}
		}},
		{"${*}", func(t *test, tk Tokens) { // 4
			t.Output = ParameterExpansion{
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Tokens: tk[:3],
			}
		}},
		{"${!}", func(t *test, tk Tokens) { // 5
			t.Output = ParameterExpansion{
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Tokens: tk[:3],
			}
		}},
		{"${0}", func(t *test, tk Tokens) { // 6
			t.Output = ParameterExpansion{
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Tokens: tk[:3],
			}
		}},
		{"${99}", func(t *test, tk Tokens) { // 7
			t.Output = ParameterExpansion{
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Tokens: tk[:3],
			}
		}},
		{"${#a}", func(t *test, tk Tokens) { // 8
			t.Output = ParameterExpansion{
				Type: ParameterLength,
				Parameter: Parameter{
					Parameter: &tk[2],
					Tokens:    tk[2:3],
				},
				Tokens: tk[:4],
			}
		}},
		{"${!a}", func(t *test, tk Tokens) { // 9
			t.Output = ParameterExpansion{
				Indirect: true,
				Parameter: Parameter{
					Parameter: &tk[2],
					Tokens:    tk[2:3],
				},
				Tokens: tk[:4],
			}
		}},
		{"${a:=b}", func(t *test, tk Tokens) { // 10
			t.Output = ParameterExpansion{
				Type: ParameterSubstitution,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[3],
							Tokens: tk[3:4],
						},
					},
					Tokens: tk[3:4],
				},
				Tokens: tk[:5],
			}
		}},
		{"${a:?b}", func(t *test, tk Tokens) { // 11
			t.Output = ParameterExpansion{
				Type: ParameterAssignment,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[3],
							Tokens: tk[3:4],
						},
					},
					Tokens: tk[3:4],
				},
				Tokens: tk[:5],
			}
		}},
		{"${a:+b}", func(t *test, tk Tokens) { // 12
			t.Output = ParameterExpansion{
				Type: ParameterMessage,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[3],
							Tokens: tk[3:4],
						},
					},
					Tokens: tk[3:4],
				},
				Tokens: tk[:5],
			}
		}},
		{"${a:-b}", func(t *test, tk Tokens) { // 13
			t.Output = ParameterExpansion{
				Type: ParameterSetAssign,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[3],
							Tokens: tk[3:4],
						},
					},
					Tokens: tk[3:4],
				},
				Tokens: tk[:5],
			}
		}},
		{"${a:1}", func(t *test, tk Tokens) { // 14
			t.Output = ParameterExpansion{
				Type: ParameterSubstring,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				SubstringStart: &tk[3],
				Tokens:         tk[:5],
			}
		}},
		{"${a: 1}", func(t *test, tk Tokens) { // 15
			t.Output = ParameterExpansion{
				Type: ParameterSubstring,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				SubstringStart: &tk[4],
				Tokens:         tk[:6],
			}
		}},
		{"${a: -1}", func(t *test, tk Tokens) { // 16
			t.Output = ParameterExpansion{
				Type: ParameterSubstring,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				SubstringStart: &tk[4],
				Tokens:         tk[:6],
			}
		}},
		{"${a:1:2}", func(t *test, tk Tokens) { // 17
			t.Output = ParameterExpansion{
				Type: ParameterSubstring,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				SubstringStart: &tk[3],
				SubstringEnd:   &tk[5],
				Tokens:         tk[:7],
			}
		}},
		{"${a:1:-2}", func(t *test, tk Tokens) { // 18
			t.Output = ParameterExpansion{
				Type: ParameterSubstring,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				SubstringStart: &tk[3],
				SubstringEnd:   &tk[5],
				Tokens:         tk[:7],
			}
		}},
		{"${a:1: -2}", func(t *test, tk Tokens) { // 19
			t.Output = ParameterExpansion{
				Type: ParameterSubstring,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				SubstringStart: &tk[3],
				SubstringEnd:   &tk[6],
				Tokens:         tk[:8],
			}
		}},
		{"${a#b}", func(t *test, tk Tokens) { // 20
			t.Output = ParameterExpansion{
				Type: ParameterRemoveStartShortest,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[3],
							Tokens: tk[3:4],
						},
					},
					Tokens: tk[3:4],
				},
				Tokens: tk[:5],
			}
		}},
		{"${a##b}", func(t *test, tk Tokens) { // 21
			t.Output = ParameterExpansion{
				Type: ParameterRemoveStartLongest,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[3],
							Tokens: tk[3:4],
						},
					},
					Tokens: tk[3:4],
				},
				Tokens: tk[:5],
			}
		}},
		{"${a%b}", func(t *test, tk Tokens) { // 22
			t.Output = ParameterExpansion{
				Type: ParameterRemoveEndShortest,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[3],
							Tokens: tk[3:4],
						},
					},
					Tokens: tk[3:4],
				},
				Tokens: tk[:5],
			}
		}},
		{"${a%%b}", func(t *test, tk Tokens) { // 23
			t.Output = ParameterExpansion{
				Type: ParameterRemoveEndLongest,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[3],
							Tokens: tk[3:4],
						},
					},
					Tokens: tk[3:4],
				},
				Tokens: tk[:5],
			}
		}},
		{"${a/b}", func(t *test, tk Tokens) { // 24
			t.Output = ParameterExpansion{
				Type: ParameterReplace,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Pattern: &tk[3],
				Tokens:  tk[:5],
			}
		}},
		{"${a/b/c}", func(t *test, tk Tokens) { // 25
			t.Output = ParameterExpansion{
				Type: ParameterReplace,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Pattern: &tk[3],
				String: &String{
					WordsOrTokens: []WordOrToken{
						{
							Word: &Word{
								Parts: []WordPart{
									{
										Part:   &tk[5],
										Tokens: tk[5:6],
									},
								},
								Tokens: tk[5:6],
							},
							Tokens: tk[5:6],
						},
					},
					Tokens: tk[5:6],
				},
				Tokens: tk[:7],
			}
		}},
		{"${a//b}", func(t *test, tk Tokens) { // 26
			t.Output = ParameterExpansion{
				Type: ParameterReplaceAll,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Pattern: &tk[3],
				Tokens:  tk[:5],
			}
		}},
		{"${a//b/c}", func(t *test, tk Tokens) { // 27
			t.Output = ParameterExpansion{
				Type: ParameterReplaceAll,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Pattern: &tk[3],
				String: &String{
					WordsOrTokens: []WordOrToken{
						{
							Word: &Word{
								Parts: []WordPart{
									{
										Part:   &tk[5],
										Tokens: tk[5:6],
									},
								},
								Tokens: tk[5:6],
							},
							Tokens: tk[5:6],
						},
					},
					Tokens: tk[5:6],
				},
				Tokens: tk[:7],
			}
		}},
		{"${a/#b}", func(t *test, tk Tokens) { // 28
			t.Output = ParameterExpansion{
				Type: ParameterReplaceStart,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Pattern: &tk[3],
				Tokens:  tk[:5],
			}
		}},
		{"${a/#b/c}", func(t *test, tk Tokens) { // 29
			t.Output = ParameterExpansion{
				Type: ParameterReplaceStart,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Pattern: &tk[3],
				String: &String{
					WordsOrTokens: []WordOrToken{
						{
							Word: &Word{
								Parts: []WordPart{
									{
										Part:   &tk[5],
										Tokens: tk[5:6],
									},
								},
								Tokens: tk[5:6],
							},
							Tokens: tk[5:6],
						},
					},
					Tokens: tk[5:6],
				},
				Tokens: tk[:7],
			}
		}},
		{"${a/%b}", func(t *test, tk Tokens) { // 30
			t.Output = ParameterExpansion{
				Type: ParameterReplaceEnd,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Pattern: &tk[3],
				Tokens:  tk[:5],
			}
		}},
		{"${a/%b/c}", func(t *test, tk Tokens) { // 31
			t.Output = ParameterExpansion{
				Type: ParameterReplaceEnd,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Pattern: &tk[3],
				String: &String{
					WordsOrTokens: []WordOrToken{
						{
							Word: &Word{
								Parts: []WordPart{
									{
										Part:   &tk[5],
										Tokens: tk[5:6],
									},
								},
								Tokens: tk[5:6],
							},
							Tokens: tk[5:6],
						},
					},
					Tokens: tk[5:6],
				},
				Tokens: tk[:7],
			}
		}},
		{"${a^b}", func(t *test, tk Tokens) { // 32
			t.Output = ParameterExpansion{
				Type: ParameterUppercaseFirstMatch,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Pattern: &tk[3],
				Tokens:  tk[:5],
			}
		}},
		{"${a^^b}", func(t *test, tk Tokens) { // 33
			t.Output = ParameterExpansion{
				Type: ParameterUppercaseAllMatches,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Pattern: &tk[3],
				Tokens:  tk[:5],
			}
		}},
		{"${a,b}", func(t *test, tk Tokens) { // 34
			t.Output = ParameterExpansion{
				Type: ParameterLowercaseFirstMatch,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Pattern: &tk[3],
				Tokens:  tk[:5],
			}
		}},
		{"${a,,b}", func(t *test, tk Tokens) { // 35
			t.Output = ParameterExpansion{
				Type: ParameterLowercaseAllMatches,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Pattern: &tk[3],
				Tokens:  tk[:5],
			}
		}},
		{"${!a@}", func(t *test, tk Tokens) { // 36
			t.Output = ParameterExpansion{
				Type: ParameterPrefixSeperate,
				Parameter: Parameter{
					Parameter: &tk[2],
					Tokens:    tk[2:3],
				},
				Tokens: tk[:5],
			}
		}},
		{"${!a*}", func(t *test, tk Tokens) { // 37
			t.Output = ParameterExpansion{
				Type: ParameterPrefix,
				Parameter: Parameter{
					Parameter: &tk[2],
					Tokens:    tk[2:3],
				},
				Tokens: tk[:5],
			}
		}},
		{"${a@U}", func(t *test, tk Tokens) { // 38
			t.Output = ParameterExpansion{
				Type: ParameterUppercase,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Tokens: tk[:5],
			}
		}},
		{"${a@u}", func(t *test, tk Tokens) { // 39
			t.Output = ParameterExpansion{
				Type: ParameterUppercaseFirst,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Tokens: tk[:5],
			}
		}},
		{"${a@L}", func(t *test, tk Tokens) { // 40
			t.Output = ParameterExpansion{
				Type: ParameterLowercase,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Tokens: tk[:5],
			}
		}},
		{"${a@Q}", func(t *test, tk Tokens) { // 41
			t.Output = ParameterExpansion{
				Type: ParameterQuoted,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Tokens: tk[:5],
			}
		}},
		{"${a@E}", func(t *test, tk Tokens) { // 42
			t.Output = ParameterExpansion{
				Type: ParameterEscaped,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Tokens: tk[:5],
			}
		}},
		{"${a@P}", func(t *test, tk Tokens) { // 43
			t.Output = ParameterExpansion{
				Type: ParameterPrompt,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Tokens: tk[:5],
			}
		}},
		{"${a@A}", func(t *test, tk Tokens) { // 44
			t.Output = ParameterExpansion{
				Type: ParameterDeclare,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Tokens: tk[:5],
			}
		}},
		{"${a@K}", func(t *test, tk Tokens) { // 45
			t.Output = ParameterExpansion{
				Type: ParameterQuotedArrays,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Tokens: tk[:5],
			}
		}},
		{"${a@a}", func(t *test, tk Tokens) { // 46
			t.Output = ParameterExpansion{
				Type: ParameterAttributes,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Tokens: tk[:5],
			}
		}},
		{"${a@k}", func(t *test, tk Tokens) { // 47
			t.Output = ParameterExpansion{
				Type: ParameterQuotedArraysSeperate,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Tokens: tk[:5],
			}
		}},
		{"${a[$(||)]}", func(t *test, tk Tokens) { // 48
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
													Err: Error{
														Err:     ErrMissingWord,
														Parsing: "Command",
														Token:   tk[4],
													},
													Parsing: "CommandOrControl",
													Token:   tk[4],
												},
												Parsing: "Pipeline",
												Token:   tk[4],
											},
											Parsing: "Statement",
											Token:   tk[4],
										},
										Parsing: "Line",
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
			}
		}},
		{"${a:=$(||)}", func(t *test, tk Tokens) { // 49
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
												Parsing: "CommandOrControl",
												Token:   tk[4],
											},
											Parsing: "Pipeline",
											Token:   tk[4],
										},
										Parsing: "Statement",
										Token:   tk[4],
									},
									Parsing: "Line",
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
				Parsing: "ParameterExpansion",
				Token:   tk[3],
			}
		}},
		{"${a:1:2b}", func(t *test, tk Tokens) { // 50
			t.Err = Error{
				Err:     ErrMissingClosingBrace,
				Parsing: "ParameterExpansion",
				Token:   tk[6],
			}
		}},
	}, func(t *test) (Type, error) {
		var pe ParameterExpansion

		err := pe.parse(t.Parser)

		return pe, err
	})
}

func TestParameter(t *testing.T) {
	doTests(t, []sourceFn{
		{"${a}", func(t *test, tk Tokens) { // 1
			t.Output = Parameter{
				Parameter: &tk[1],
				Tokens:    tk[1:2],
			}
		}},
		{"${0}", func(t *test, tk Tokens) { // 2
			t.Output = Parameter{
				Parameter: &tk[1],
				Tokens:    tk[1:2],
			}
		}},
		{"${9}", func(t *test, tk Tokens) { // 3
			t.Output = Parameter{
				Parameter: &tk[1],
				Tokens:    tk[1:2],
			}
		}},
		{"${@}", func(t *test, tk Tokens) { // 4
			t.Output = Parameter{
				Parameter: &tk[1],
				Tokens:    tk[1:2],
			}
		}},
		{"${*}", func(t *test, tk Tokens) { // 5
			t.Output = Parameter{
				Parameter: &tk[1],
				Tokens:    tk[1:2],
			}
		}},
		{"${a[0]}", func(t *test, tk Tokens) { // 6
			t.Output = Parameter{
				Parameter: &tk[1],
				Array: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[3],
							Tokens: tk[3:4],
						},
					},
					Tokens: tk[3:4],
				},
				Tokens: tk[1:5],
			}
		}},
		{"${a[@]}", func(t *test, tk Tokens) { // 7
			t.Output = Parameter{
				Parameter: &tk[1],
				Array: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[3],
							Tokens: tk[3:4],
						},
					},
					Tokens: tk[3:4],
				},
				Tokens: tk[1:5],
			}
		}},
		{"${a[*]}", func(t *test, tk Tokens) { // 8
			t.Output = Parameter{
				Parameter: &tk[1],
				Array: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[3],
							Tokens: tk[3:4],
						},
					},
					Tokens: tk[3:4],
				},
				Tokens: tk[1:5],
			}
		}},
		{"${a[$(||)]}", func(t *test, tk Tokens) { // 9
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
												Parsing: "CommandOrControl",
												Token:   tk[4],
											},
											Parsing: "Pipeline",
											Token:   tk[4],
										},
										Parsing: "Statement",
										Token:   tk[4],
									},
									Parsing: "Line",
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
			}
		}},
	}, func(t *test) (Type, error) {
		var p Parameter

		t.Parser.Tokens = t.Parser.Tokens[1:1]

		err := p.parse(t.Parser)

		return p, err
	})
}

func TestString(t *testing.T) {
	doTests(t, []sourceFn{
		{"a", func(t *test, tk Tokens) { // 1
			t.Output = String{
				WordsOrTokens: []WordOrToken{
					{
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
					},
				},
				Tokens: tk[:1],
			}
		}},
		{"a\t or  b", func(t *test, tk Tokens) { // 2
			t.Output = String{
				WordsOrTokens: []WordOrToken{
					{
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
					},
					{
						Token:  &tk[1],
						Tokens: tk[1:2],
					},
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
					{
						Token:  &tk[3],
						Tokens: tk[3:4],
					},
					{
						Word: &Word{
							Parts: []WordPart{
								{
									Part:   &tk[4],
									Tokens: tk[4:5],
								},
							},
							Tokens: tk[4:5],
						},
						Tokens: tk[4:5],
					},
				},
				Tokens: tk[:5],
			}
		}},
		{"$a", func(t *test, tk Tokens) { // 3
			t.Output = String{
				WordsOrTokens: []WordOrToken{
					{
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
					},
				},
				Tokens: tk[:1],
			}
		}},
		{"$(||)", func(t *test, tk Tokens) { // 4
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
													Err: Error{
														Err:     ErrMissingWord,
														Parsing: "Command",
														Token:   tk[1],
													},
													Parsing: "CommandOrControl",
													Token:   tk[1],
												},
												Parsing: "Pipeline",
												Token:   tk[1],
											},
											Parsing: "Statement",
											Token:   tk[1],
										},
										Parsing: "Line",
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
				},
				Parsing: "String",
				Token:   tk[0],
			}
		}},
	}, func(t *test) (Type, error) {
		var s String

		err := s.parse(t.Parser)

		return s, err
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
											Err: Error{
												Err: Error{
													Err:     ErrMissingWord,
													Parsing: "Command",
													Token:   tk[1],
												},
												Parsing: "CommandOrControl",
												Token:   tk[1],
											},
											Parsing: "Pipeline",
											Token:   tk[1],
										},
										Parsing: "Statement",
										Token:   tk[1],
									},
									Parsing: "Line",
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
					Lines: []Line{
						{
							Statements: []Statement{
								{
									Pipeline: Pipeline{
										CommandOrControl: CommandOrControl{
											Command: &Command{
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
									Tokens: tk[1:3],
								},
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
					Lines: []Line{
						{
							Statements: []Statement{
								{
									Pipeline: Pipeline{
										CommandOrControl: CommandOrControl{
											Command: &Command{
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
									Tokens: tk[1:3],
								},
							},
							Tokens: tk[1:3],
						},
					},
					Tokens: tk[1:3],
				},
				Tokens: tk[:4],
			}
		}},
		{"$(||)", func(t *test, tk Tokens) { // 5
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
									Parsing: "CommandOrControl",
									Token:   tk[1],
								},
								Parsing: "Pipeline",
								Token:   tk[1],
							},
							Parsing: "Statement",
							Token:   tk[1],
						},
						Parsing: "Line",
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
										Err: Error{
											Err: Error{
												Err: Error{
													Err:     ErrMissingWord,
													Parsing: "Command",
													Token:   tk[2],
												},
												Parsing: "CommandOrControl",
												Token:   tk[2],
											},
											Parsing: "Pipeline",
											Token:   tk[2],
										},
										Parsing: "Statement",
										Token:   tk[2],
									},
									Parsing: "Line",
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
				Parsing: "Redirection",
				Token:   tk[1],
			}
		}},
	}, func(t *test) (Type, error) {
		var r Redirection

		err := r.parse(t.Parser)

		return r, err
	})
}

func TestHeredoc(t *testing.T) {
	doTests(t, []sourceFn{
		{"<<a\nb\na", func(t *test, tk Tokens) { // 1
			t.Output = Heredoc{
				HeredocPartsOrWords: []HeredocPartOrWord{
					{
						HeredocPart: &tk[3],
						Tokens:      tk[3:4],
					},
				},
				Tokens: tk[3:5],
			}
		}},
		{"<<a\n\tb\na", func(t *test, tk Tokens) { // 2
			t.Output = Heredoc{
				HeredocPartsOrWords: []HeredocPartOrWord{
					{
						HeredocPart: &tk[3],
						Tokens:      tk[3:4],
					},
				},
				Tokens: tk[3:5],
			}
		}},
		{"<<-a\n\tb\na", func(t *test, tk Tokens) { // 3
			t.Output = Heredoc{
				HeredocPartsOrWords: []HeredocPartOrWord{
					{
						HeredocPart: &tk[4],
						Tokens:      tk[4:5],
					},
				},
				Tokens: tk[3:6],
			}
		}},
		{"<<a\n$a\na", func(t *test, tk Tokens) { // 4
			t.Output = Heredoc{
				HeredocPartsOrWords: []HeredocPartOrWord{
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
					{
						HeredocPart: &tk[4],
						Tokens:      tk[4:5],
					},
				},
				Tokens: tk[3:6],
			}
		}},
		{"<<a\na$b\na", func(t *test, tk Tokens) { // 5
			t.Output = Heredoc{
				HeredocPartsOrWords: []HeredocPartOrWord{
					{
						HeredocPart: &tk[3],
						Tokens:      tk[3:4],
					},
					{
						Word: &Word{
							Parts: []WordPart{
								{
									Part:   &tk[4],
									Tokens: tk[4:5],
								},
							},
							Tokens: tk[4:5],
						},
						Tokens: tk[4:5],
					},
					{
						HeredocPart: &tk[5],
						Tokens:      tk[5:6],
					},
				},
				Tokens: tk[3:7],
			}
		}},
		{"<<a\n$(||)\na", func(t *test, tk Tokens) { // 6
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
													Err: Error{
														Err:     ErrMissingWord,
														Parsing: "Command",
														Token:   tk[4],
													},
													Parsing: "CommandOrControl",
													Token:   tk[4],
												},
												Parsing: "Pipeline",
												Token:   tk[4],
											},
											Parsing: "Statement",
											Token:   tk[4],
										},
										Parsing: "Line",
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
					Parsing: "HeredocPartOrWord",
					Token:   tk[3],
				},
				Parsing: "Heredoc",
				Token:   tk[3],
			}
		}},
	}, func(t *test) (Type, error) {
		var h Heredoc

		t.Parser.Tokens = t.Parser.Tokens[3:3]

		err := h.parse(t.Parser)

		return h, err
	})
}

func TestHeredocPartOrWord(t *testing.T) {
	doTests(t, []sourceFn{
		{"<<a\nb\na", func(t *test, tk Tokens) { // 1
			t.Output = HeredocPartOrWord{
				HeredocPart: &tk[3],
				Tokens:      tk[3:4],
			}
		}},
		{"<<a\n$b\na", func(t *test, tk Tokens) { // 2
			t.Output = HeredocPartOrWord{
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
			}
		}},
		{"<<a\n$(||)\na", func(t *test, tk Tokens) { // 3
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
												Parsing: "CommandOrControl",
												Token:   tk[4],
											},
											Parsing: "Pipeline",
											Token:   tk[4],
										},
										Parsing: "Statement",
										Token:   tk[4],
									},
									Parsing: "Line",
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
				Parsing: "HeredocPartOrWord",
				Token:   tk[3],
			}
		}},
	}, func(t *test) (Type, error) {
		var h HeredocPartOrWord

		t.Parser.Tokens = t.Parser.Tokens[3:3]

		err := h.parse(t.Parser)

		return h, err
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
												Err: Error{
													Err: Error{
														Err:     ErrMissingWord,
														Parsing: "Command",
														Token:   tk[2],
													},
													Parsing: "CommandOrControl",
													Token:   tk[2],
												},
												Parsing: "Pipeline",
												Token:   tk[2],
											},
											Parsing: "Statement",
											Token:   tk[2],
										},
										Parsing: "Line",
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
									Err: Error{
										Err: Error{
											Err: Error{
												Err: Error{
													Err:     ErrMissingWord,
													Parsing: "Command",
													Token:   tk[1],
												},
												Parsing: "CommandOrControl",
												Token:   tk[1],
											},
											Parsing: "Pipeline",
											Token:   tk[1],
										},
										Parsing: "Statement",
										Token:   tk[1],
									},
									Parsing: "Line",
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
				Parsing: "WordOrOperator",
				Token:   tk[0],
			}
		}},
	}, func(t *test) (Type, error) {
		var wo WordOrOperator

		err := wo.parse(t.Parser)

		return wo, err
	})
}
