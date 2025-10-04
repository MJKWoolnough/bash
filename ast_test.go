package bash

import (
	"io"
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

func TestParse(t *testing.T) {
	tk := []Token{
		{
			Token: parser.Token{
				Type: TokenWord,
				Data: "(",
			},
		},
		{
			Token: parser.Token{
				Type: parser.TokenError,
				Data: "unexpected EOF",
			},
			Pos:     1,
			LinePos: 1,
		},
	}
	expectedErr := Error{
		Err:     io.ErrUnexpectedEOF,
		Parsing: "Tokens",
		Token:   tk[1],
	}
	_, err := Parse(makeTokeniser(parser.NewStringTokeniser("(")))
	if !reflect.DeepEqual(err, expectedErr) {
		t.Errorf("expecting error: %v, got %v", expectedErr, err)
	}

	tk = []Token{
		{
			Token: parser.Token{
				Type: TokenWord,
				Data: "a",
			},
		},
		{
			Token: parser.Token{
				Type: parser.TokenDone,
			},
			Pos:     1,
			LinePos: 1,
		},
	}
	expectedFile := &File{
		Lines: []Line{
			{
				Statements: []Statement{
					{
						Pipeline: Pipeline{
							CommandOrCompound: CommandOrCompound{
								Command: &Command{
									AssignmentsOrWords: []AssignmentOrWord{
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

	f, err := Parse(makeTokeniser(parser.NewStringTokeniser("a")))
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else if !reflect.DeepEqual(f, expectedFile) {
		t.Errorf("expecting \n%+v\n...got...\n%+v", expectedFile, f)
	}

	tk = []Token{
		{
			Token: parser.Token{
				Type: TokenPunctuator,
				Data: "||",
			},
		},
		{
			Token: parser.Token{
				Type: parser.TokenDone,
			},
			Pos:     1,
			LinePos: 1,
		},
	}
	expectedErr = Error{
		Err: Error{
			Err: Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err:     ErrMissingWord,
							Parsing: "Command",
							Token:   tk[0],
						},
						Parsing: "CommandOrCompound",
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

	_, err = Parse(makeTokeniser(parser.NewStringTokeniser("||")))
	if !reflect.DeepEqual(err, expectedErr) {
		t.Errorf("expecting error: %v, got %v", expectedErr, err)
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
									CommandOrCompound: CommandOrCompound{
										Command: &Command{
											AssignmentsOrWords: []AssignmentOrWord{
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
									CommandOrCompound: CommandOrCompound{
										Command: &Command{
											AssignmentsOrWords: []AssignmentOrWord{
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
									CommandOrCompound: CommandOrCompound{
										Command: &Command{
											AssignmentsOrWords: []AssignmentOrWord{
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
									CommandOrCompound: CommandOrCompound{
										Command: &Command{
											AssignmentsOrWords: []AssignmentOrWord{
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
									CommandOrCompound: CommandOrCompound{
										Command: &Command{
											AssignmentsOrWords: []AssignmentOrWord{
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
								Parsing: "CommandOrCompound",
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
									CommandOrCompound: CommandOrCompound{
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
									CommandOrCompound: CommandOrCompound{
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
									CommandOrCompound: CommandOrCompound{
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
		{"#!/bin/bash", func(t *test, tk Tokens) { // 7
			t.Output = File{
				Comments: [2]Comments{{tk[0]}},
				Tokens:   tk[:1],
			}
		}},
		{" #!/bin/bash", func(t *test, tk Tokens) { // 8
			t.Output = File{
				Comments: [2]Comments{nil, {tk[1]}},
				Tokens:   tk[:2],
			}
		}},
		{"\n#!/bin/bash", func(t *test, tk Tokens) { // 9
			t.Output = File{
				Comments: [2]Comments{nil, {tk[1]}},
				Tokens:   tk[:2],
			}
		}},
		{"#!/bin/bash\n# comment", func(t *test, tk Tokens) { // 10
			t.Output = File{
				Comments: [2]Comments{{tk[0], tk[2]}},
				Tokens:   tk[:3],
			}
		}},
		{"#!/bin/bash\n# comment\n\n# final\n# comment", func(t *test, tk Tokens) { // 11
			t.Output = File{
				Comments: [2]Comments{{tk[0], tk[2]}, {tk[4], tk[6]}},
				Tokens:   tk[:7],
			}
		}},
		{"#!/bin/bash\n# comment\n\n# pre-line comment\na #post-line comment\n# another post line comment\n\n# final\n# comment", func(t *test, tk Tokens) { // 12
			t.Output = File{
				Lines: []Line{
					{
						Statements: []Statement{
							{
								Pipeline: Pipeline{
									CommandOrCompound: CommandOrCompound{
										Command: &Command{
											AssignmentsOrWords: []AssignmentOrWord{
												{
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
											},
											Tokens: tk[6:7],
										},
										Tokens: tk[6:7],
									},
									Tokens: tk[6:7],
								},
								Tokens: tk[6:7],
							},
						},
						Comments: [2]Comments{{tk[4]}, {tk[8], tk[10]}},
						Tokens:   tk[4:11],
					},
				},
				Comments: [2]Comments{{tk[0], tk[2]}, {tk[12], tk[14]}},
				Tokens:   tk[:15],
			}
		}},
		{"<<a\n$(||)\na", func(t *test, tk Tokens) { // 13
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
																				Parsing: "CommandOrCompound",
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
								Parsing: "CommandOrCompound",
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
							CommandOrCompound: CommandOrCompound{
								Command: &Command{
									AssignmentsOrWords: []AssignmentOrWord{
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
							CommandOrCompound: CommandOrCompound{
								Command: &Command{
									AssignmentsOrWords: []AssignmentOrWord{
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
								},
								Tokens: tk[:1],
							},
							Tokens: tk[:1],
						},
						Tokens: tk[:2],
					},
					{
						Pipeline: Pipeline{
							CommandOrCompound: CommandOrCompound{
								Command: &Command{
									AssignmentsOrWords: []AssignmentOrWord{
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
							CommandOrCompound: CommandOrCompound{
								Command: &Command{
									AssignmentsOrWords: []AssignmentOrWord{
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
							CommandOrCompound: CommandOrCompound{
								Command: &Command{
									AssignmentsOrWords: []AssignmentOrWord{
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
		{"# comment\na", func(t *test, tk Tokens) { // 4
			t.Output = Line{
				Statements: []Statement{
					{
						Pipeline: Pipeline{
							CommandOrCompound: CommandOrCompound{
								Command: &Command{
									AssignmentsOrWords: []AssignmentOrWord{
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
									Tokens: tk[2:3],
								},
								Tokens: tk[2:3],
							},
							Tokens: tk[2:3],
						},
						Tokens: tk[2:3],
					},
				},
				Comments: [2]Comments{{tk[0]}},
				Tokens:   tk[:3],
			}
		}},
		{"# comment\n\n# another comment\na", func(t *test, tk Tokens) { // 5
			t.Output = Line{
				Statements: []Statement{
					{
						Pipeline: Pipeline{
							CommandOrCompound: CommandOrCompound{
								Command: &Command{
									AssignmentsOrWords: []AssignmentOrWord{
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
									Tokens: tk[4:5],
								},
								Tokens: tk[4:5],
							},
							Tokens: tk[4:5],
						},
						Tokens: tk[4:5],
					},
				},
				Comments: [2]Comments{{tk[0], tk[2]}},
				Tokens:   tk[:5],
			}
		}},
		{"a # comment", func(t *test, tk Tokens) { // 6
			t.Output = Line{
				Statements: []Statement{
					{
						Pipeline: Pipeline{
							CommandOrCompound: CommandOrCompound{
								Command: &Command{
									AssignmentsOrWords: []AssignmentOrWord{
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
								},
								Tokens: tk[:1],
							},
							Tokens: tk[:1],
						},
						Tokens: tk[:1],
					},
				},
				Comments: [2]Comments{nil, {tk[2]}},
				Tokens:   tk[:3],
			}
		}},
		{"a # comment\n# another comment", func(t *test, tk Tokens) { // 7
			t.Output = Line{
				Statements: []Statement{
					{
						Pipeline: Pipeline{
							CommandOrCompound: CommandOrCompound{
								Command: &Command{
									AssignmentsOrWords: []AssignmentOrWord{
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
								},
								Tokens: tk[:1],
							},
							Tokens: tk[:1],
						},
						Tokens: tk[:1],
					},
				},
				Comments: [2]Comments{nil, {tk[2], tk[4]}},
				Tokens:   tk[:5],
			}
		}},
		{"a # comment\n\n# another comment", func(t *test, tk Tokens) { // 8
			t.Output = Line{
				Statements: []Statement{
					{
						Pipeline: Pipeline{
							CommandOrCompound: CommandOrCompound{
								Command: &Command{
									AssignmentsOrWords: []AssignmentOrWord{
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
								},
								Tokens: tk[:1],
							},
							Tokens: tk[:1],
						},
						Tokens: tk[:1],
					},
				},
				Comments: [2]Comments{nil, {tk[2]}},
				Tokens:   tk[:3],
			}
		}},
		{"||", func(t *test, tk Tokens) { // 9
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err:     ErrMissingWord,
								Parsing: "Command",
								Token:   tk[0],
							},
							Parsing: "CommandOrCompound",
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

func TestLineHeredocs(t *testing.T) {
	doTests(t, []sourceFn{
		{"a >b", func(t *test, tk Tokens) { // 1
			t.Output = Line{
				Statements: []Statement{
					{
						Pipeline: Pipeline{
							CommandOrCompound: CommandOrCompound{
								Command: &Command{
									AssignmentsOrWords: []AssignmentOrWord{
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
									Redirections: []Redirection{
										{
											Redirector: &tk[2],
											Output: Word{
												Parts: []WordPart{
													{
														Part:   &tk[3],
														Tokens: tk[3:4],
													},
												},
												Tokens: tk[3:4],
											},
											Tokens: tk[2:4],
										},
									},
									Tokens: tk[:4],
								},
								Tokens: tk[:4],
							},
							Tokens: tk[:4],
						},
						Tokens: tk[:4],
					},
				},
				Tokens: tk[:4],
			}
		}},
		{"a <<b\nSome content\nb", func(t *test, tk Tokens) { // 2
			t.Output = Line{
				Statements: []Statement{
					{
						Pipeline: Pipeline{
							CommandOrCompound: CommandOrCompound{
								Command: &Command{
									AssignmentsOrWords: []AssignmentOrWord{
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
									Redirections: []Redirection{
										{
											Redirector: &tk[2],
											Output: Word{
												Parts: []WordPart{
													{
														Part:   &tk[3],
														Tokens: tk[3:4],
													},
												},
												Tokens: tk[3:4],
											},
											Heredoc: &Heredoc{
												HeredocPartsOrWords: []HeredocPartOrWord{
													{
														HeredocPart: &tk[5],
														Tokens:      tk[5:6],
													},
												},
												Tokens: tk[5:7],
											},
											Tokens: tk[2:4],
										},
									},
									Tokens: tk[:4],
								},
								Tokens: tk[:4],
							},
							Tokens: tk[:4],
						},
						Tokens: tk[:4],
					},
				},
				Tokens: tk[:7],
			}
		}},
		{"a <<b #comment\nSome content\nb", func(t *test, tk Tokens) { // 3
			t.Output = Line{
				Statements: []Statement{
					{
						Pipeline: Pipeline{
							CommandOrCompound: CommandOrCompound{
								Command: &Command{
									AssignmentsOrWords: []AssignmentOrWord{
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
									Redirections: []Redirection{
										{
											Redirector: &tk[2],
											Output: Word{
												Parts: []WordPart{
													{
														Part:   &tk[3],
														Tokens: tk[3:4],
													},
												},
												Tokens: tk[3:4],
											},
											Heredoc: &Heredoc{
												HeredocPartsOrWords: []HeredocPartOrWord{
													{
														HeredocPart: &tk[7],
														Tokens:      tk[7:8],
													},
												},
												Tokens: tk[7:9],
											},
											Tokens: tk[2:4],
										},
									},
									Tokens: tk[:4],
								},
								Tokens: tk[:4],
							},
							Tokens: tk[:4],
						},
						Tokens: tk[:4],
					},
				},
				Comments: [2]Comments{nil, {tk[5]}},
				Tokens:   tk[:9],
			}
		}},
		{"a | b <<c\nSome content\nc", func(t *test, tk Tokens) { // 4
			t.Output = Line{
				Statements: []Statement{
					{
						Pipeline: Pipeline{
							CommandOrCompound: CommandOrCompound{
								Command: &Command{
									AssignmentsOrWords: []AssignmentOrWord{
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
								},
								Tokens: tk[:1],
							},
							Pipeline: &Pipeline{
								CommandOrCompound: CommandOrCompound{
									Command: &Command{
										AssignmentsOrWords: []AssignmentOrWord{
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
										Redirections: []Redirection{
											{
												Redirector: &tk[6],
												Output: Word{
													Parts: []WordPart{
														{
															Part:   &tk[7],
															Tokens: tk[7:8],
														},
													},
													Tokens: tk[7:8],
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
												Tokens: tk[6:8],
											},
										},
										Tokens: tk[4:8],
									},
									Tokens: tk[4:8],
								},
								Tokens: tk[4:8],
							},
							Tokens: tk[:8],
						},
						Tokens: tk[:8],
					},
				},
				Tokens: tk[:11],
			}
		}},
		{"a <<b | c\nSome content\nb", func(t *test, tk Tokens) { // 5
			t.Output = Line{
				Statements: []Statement{
					{
						Pipeline: Pipeline{
							CommandOrCompound: CommandOrCompound{
								Command: &Command{
									AssignmentsOrWords: []AssignmentOrWord{
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
									Redirections: []Redirection{
										{
											Redirector: &tk[2],
											Output: Word{
												Parts: []WordPart{
													{
														Part:   &tk[3],
														Tokens: tk[3:4],
													},
												},
												Tokens: tk[3:4],
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
											Tokens: tk[2:4],
										},
									},
									Tokens: tk[:4],
								},
								Tokens: tk[:4],
							},
							Pipeline: &Pipeline{
								CommandOrCompound: CommandOrCompound{
									Command: &Command{
										AssignmentsOrWords: []AssignmentOrWord{
											{
												Word: &Word{
													Parts: []WordPart{
														{
															Part:   &tk[7],
															Tokens: tk[7:8],
														},
													},
													Tokens: tk[7:8],
												},
												Tokens: tk[7:8],
											},
										},
										Tokens: tk[7:8],
									},
									Tokens: tk[7:8],
								},
								Tokens: tk[7:8],
							},
							Tokens: tk[:8],
						},
						Tokens: tk[:8],
					},
				},
				Tokens: tk[:11],
			}
		}},
		{"a || b <<c\nSome content\nc", func(t *test, tk Tokens) { // 6
			t.Output = Line{
				Statements: []Statement{
					{
						Pipeline: Pipeline{
							CommandOrCompound: CommandOrCompound{
								Command: &Command{
									AssignmentsOrWords: []AssignmentOrWord{
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
								},
								Tokens: tk[:1],
							},
							Tokens: tk[:1],
						},
						LogicalOperator: LogicalOperatorOr,
						Statement: &Statement{
							Pipeline: Pipeline{
								CommandOrCompound: CommandOrCompound{
									Command: &Command{
										AssignmentsOrWords: []AssignmentOrWord{
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
										Redirections: []Redirection{
											{
												Redirector: &tk[6],
												Output: Word{
													Parts: []WordPart{
														{
															Part:   &tk[7],
															Tokens: tk[7:8],
														},
													},
													Tokens: tk[7:8],
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
												Tokens: tk[6:8],
											},
										},
										Tokens: tk[4:8],
									},
									Tokens: tk[4:8],
								},
								Tokens: tk[4:8],
							},
							Tokens: tk[4:8],
						},
						Tokens: tk[:8],
					},
				},
				Tokens: tk[:11],
			}
		}},
		{"a <<b && c\nSome content\nb", func(t *test, tk Tokens) { // 7
			t.Output = Line{
				Statements: []Statement{
					{
						Pipeline: Pipeline{
							CommandOrCompound: CommandOrCompound{
								Command: &Command{
									AssignmentsOrWords: []AssignmentOrWord{
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
									Redirections: []Redirection{
										{
											Redirector: &tk[2],
											Output: Word{
												Parts: []WordPart{
													{
														Part:   &tk[3],
														Tokens: tk[3:4],
													},
												},
												Tokens: tk[3:4],
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
											Tokens: tk[2:4],
										},
									},
									Tokens: tk[:4],
								},
								Tokens: tk[:4],
							},
							Tokens: tk[:4],
						},
						LogicalOperator: LogicalOperatorAnd,
						Statement: &Statement{
							Pipeline: Pipeline{
								CommandOrCompound: CommandOrCompound{
									Command: &Command{
										AssignmentsOrWords: []AssignmentOrWord{
											{
												Word: &Word{
													Parts: []WordPart{
														{
															Part:   &tk[7],
															Tokens: tk[7:8],
														},
													},
													Tokens: tk[7:8],
												},
												Tokens: tk[7:8],
											},
										},
										Tokens: tk[7:8],
									},
									Tokens: tk[7:8],
								},
								Tokens: tk[7:8],
							},
							Tokens: tk[7:8],
						},
						Tokens: tk[:8],
					},
				},
				Tokens: tk[:11],
			}
		}},
		{"while a; do b; done <<c\nSome content\nc", func(t *test, tk Tokens) { // 8
			t.Output = Line{
				Statements: []Statement{
					{
						Pipeline: Pipeline{
							CommandOrCompound: CommandOrCompound{
								Compound: &Compound{
									LoopCompound: &LoopCompound{
										Statement: Statement{
											Pipeline: Pipeline{
												CommandOrCompound: CommandOrCompound{
													Command: &Command{
														AssignmentsOrWords: []AssignmentOrWord{
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
														Tokens: tk[2:3],
													},
													Tokens: tk[2:3],
												},
												Tokens: tk[2:3],
											},
											Tokens: tk[2:4],
										},
										File: File{
											Lines: []Line{
												{
													Statements: []Statement{
														{
															Pipeline: Pipeline{
																CommandOrCompound: CommandOrCompound{
																	Command: &Command{
																		AssignmentsOrWords: []AssignmentOrWord{
																			{
																				Word: &Word{
																					Parts: []WordPart{
																						{
																							Part:   &tk[7],
																							Tokens: tk[7:8],
																						},
																					},
																					Tokens: tk[7:8],
																				},
																				Tokens: tk[7:8],
																			},
																		},
																		Tokens: tk[7:8],
																	},
																	Tokens: tk[7:8],
																},
																Tokens: tk[7:8],
															},
															Tokens: tk[7:9],
														},
													},
													Tokens: tk[7:9],
												},
											},
											Tokens: tk[7:9],
										},
										Tokens: tk[:11],
									},
									Redirections: []Redirection{
										{
											Redirector: &tk[12],
											Output: Word{
												Parts: []WordPart{
													{
														Part:   &tk[13],
														Tokens: tk[13:14],
													},
												},
												Tokens: tk[13:14],
											},
											Heredoc: &Heredoc{
												HeredocPartsOrWords: []HeredocPartOrWord{
													{
														HeredocPart: &tk[15],
														Tokens:      tk[15:16],
													},
												},
												Tokens: tk[15:17],
											},
											Tokens: tk[12:14],
										},
									},
									Tokens: tk[:14],
								},
								Tokens: tk[:14],
							},
							Tokens: tk[:14],
						},
						Tokens: tk[:14],
					},
				},
				Tokens: tk[:17],
			}
		}},
		{"a <<b\n$(||)\nb", func(t *test, tk Tokens) { // 9
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
																				Err:     ErrMissingWord,
																				Parsing: "Command",
																				Token:   tk[6],
																			},
																			Parsing: "CommandOrCompound",
																			Token:   tk[6],
																		},
																		Parsing: "Pipeline",
																		Token:   tk[6],
																	},
																	Parsing: "Statement",
																	Token:   tk[6],
																},
																Parsing: "Line",
																Token:   tk[6],
															},
															Parsing: "File",
															Token:   tk[6],
														},
														Parsing: "CommandSubstitution",
														Token:   tk[6],
													},
													Parsing: "WordPart",
													Token:   tk[5],
												},
												Parsing: "Word",
												Token:   tk[5],
											},
											Parsing: "HeredocPartOrWord",
											Token:   tk[5],
										},
										Parsing: "Heredoc",
										Token:   tk[5],
									},
									Parsing: "Redirection",
									Token:   tk[5],
								},
								Parsing: "Command",
								Token:   tk[5],
							},
							Parsing: "CommandOrCompound",
							Token:   tk[5],
						},
						Parsing: "Pipeline",
						Token:   tk[5],
					},
					Parsing: "Statement",
					Token:   tk[5],
				},
				Parsing: "Line",
				Token:   tk[5],
			}
		}},
		{"a | b <<c\n$(||)\nc", func(t *test, tk Tokens) { // 10
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
																					Token:   tk[10],
																				},
																				Parsing: "CommandOrCompound",
																				Token:   tk[10],
																			},
																			Parsing: "Pipeline",
																			Token:   tk[10],
																		},
																		Parsing: "Statement",
																		Token:   tk[10],
																	},
																	Parsing: "Line",
																	Token:   tk[10],
																},
																Parsing: "File",
																Token:   tk[10],
															},
															Parsing: "CommandSubstitution",
															Token:   tk[10],
														},
														Parsing: "WordPart",
														Token:   tk[9],
													},
													Parsing: "Word",
													Token:   tk[9],
												},
												Parsing: "HeredocPartOrWord",
												Token:   tk[9],
											},
											Parsing: "Heredoc",
											Token:   tk[9],
										},
										Parsing: "Redirection",
										Token:   tk[9],
									},
									Parsing: "Command",
									Token:   tk[9],
								},
								Parsing: "CommandOrCompound",
								Token:   tk[9],
							},
							Parsing: "Pipeline",
							Token:   tk[9],
						},
						Parsing: "Pipeline",
						Token:   tk[9],
					},
					Parsing: "Statement",
					Token:   tk[9],
				},
				Parsing: "Line",
				Token:   tk[9],
			}
		}},
		{"a && b <<c\n$(||)\nc", func(t *test, tk Tokens) { // 11
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
																					Token:   tk[10],
																				},
																				Parsing: "CommandOrCompound",
																				Token:   tk[10],
																			},
																			Parsing: "Pipeline",
																			Token:   tk[10],
																		},
																		Parsing: "Statement",
																		Token:   tk[10],
																	},
																	Parsing: "Line",
																	Token:   tk[10],
																},
																Parsing: "File",
																Token:   tk[10],
															},
															Parsing: "CommandSubstitution",
															Token:   tk[10],
														},
														Parsing: "WordPart",
														Token:   tk[9],
													},
													Parsing: "Word",
													Token:   tk[9],
												},
												Parsing: "HeredocPartOrWord",
												Token:   tk[9],
											},
											Parsing: "Heredoc",
											Token:   tk[9],
										},
										Parsing: "Redirection",
										Token:   tk[9],
									},
									Parsing: "Command",
									Token:   tk[9],
								},
								Parsing: "CommandOrCompound",
								Token:   tk[9],
							},
							Parsing: "Pipeline",
							Token:   tk[9],
						},
						Parsing: "Statement",
						Token:   tk[9],
					},
					Parsing: "Statement",
					Token:   tk[9],
				},
				Parsing: "Line",
				Token:   tk[9],
			}
		}},
		{"until a; do b; done <<b\n$(||)\nb", func(t *test, tk Tokens) { // 12
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
																				Err:     ErrMissingWord,
																				Parsing: "Command",
																				Token:   tk[16],
																			},
																			Parsing: "CommandOrCompound",
																			Token:   tk[16],
																		},
																		Parsing: "Pipeline",
																		Token:   tk[16],
																	},
																	Parsing: "Statement",
																	Token:   tk[16],
																},
																Parsing: "Line",
																Token:   tk[16],
															},
															Parsing: "File",
															Token:   tk[16],
														},
														Parsing: "CommandSubstitution",
														Token:   tk[16],
													},
													Parsing: "WordPart",
													Token:   tk[15],
												},
												Parsing: "Word",
												Token:   tk[15],
											},
											Parsing: "HeredocPartOrWord",
											Token:   tk[15],
										},
										Parsing: "Heredoc",
										Token:   tk[15],
									},
									Parsing: "Redirection",
									Token:   tk[15],
								},
								Parsing: "Compound",
								Token:   tk[15],
							},
							Parsing: "CommandOrCompound",
							Token:   tk[15],
						},
						Parsing: "Pipeline",
						Token:   tk[15],
					},
					Parsing: "Statement",
					Token:   tk[15],
				},
				Parsing: "Line",
				Token:   tk[15],
			}
		}},
	}, func(t *test) (Type, error) {
		var l Line

		err := l.parse(t.Parser)

		return l, err
	})
}
