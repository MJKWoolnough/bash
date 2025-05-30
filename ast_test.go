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

	f, err = Parse(makeTokeniser(parser.NewStringTokeniser("||")))
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

func TestStatement(t *testing.T) {
	doTests(t, []sourceFn{
		{"a", func(t *test, tk Tokens) { // 1
			t.Output = Statement{
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
			}
		}},
		{"a||b", func(t *test, tk Tokens) { // 2
			t.Output = Statement{
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
				Tokens: tk[:3],
			}
		}},
		{"a && b", func(t *test, tk Tokens) { // 3
			t.Output = Statement{
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
				Tokens: tk[:5],
			}
		}},
		{"a||b;", func(t *test, tk Tokens) { // 4
			t.Output = Statement{
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
				Tokens: tk[:4],
			}
		}},
		{"a||b &", func(t *test, tk Tokens) { // 5
			t.Output = Statement{
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
						Parsing: "CommandOrCompound",
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
							Parsing: "CommandOrCompound",
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
			}
		}},
		{"time a", func(t *test, tk Tokens) { // 2
			t.Output = Pipeline{
				PipelineTime: PipelineTimeBash,
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
				Tokens: tk[:3],
			}
		}},
		{"time -p a", func(t *test, tk Tokens) { // 3
			t.Output = Pipeline{
				PipelineTime: PipelineTimePosix,
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
				Tokens: tk[:5],
			}
		}},
		{"! a", func(t *test, tk Tokens) { // 4
			t.Output = Pipeline{
				Not: true,
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
				Tokens: tk[:3],
			}
		}},
		{"coproc a", func(t *test, tk Tokens) { // 5
			t.Output = Pipeline{
				Coproc: true,
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
				Tokens: tk[:3],
			}
		}},
		{"coproc a if b; then c;fi", func(t *test, tk Tokens) { // 6
			t.Output = Pipeline{
				Coproc:           true,
				CoprocIdentifier: &tk[2],
				CommandOrCompound: CommandOrCompound{
					Compound: &Compound{
						IfCompound: &IfCompound{
							If: TestConsequence{
								Test: Statement{
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
									Tokens: tk[6:8],
								},
								Consequence: File{
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
																					Part:   &tk[11],
																					Tokens: tk[11:12],
																				},
																			},
																			Tokens: tk[11:12],
																		},
																		Tokens: tk[11:12],
																	},
																},
																Tokens: tk[11:12],
															},
															Tokens: tk[11:12],
														},
														Tokens: tk[11:12],
													},
													Tokens: tk[11:13],
												},
											},
											Tokens: tk[11:13],
										},
									},
									Tokens: tk[11:13],
								},
								Tokens: tk[6:13],
							},
							Tokens: tk[4:14],
						},
						Tokens: tk[4:14],
					},
					Tokens: tk[4:14],
				},
				Tokens: tk[:14],
			}
		}},
		{"a|b", func(t *test, tk Tokens) { // 7
			t.Output = Pipeline{
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
				Tokens: tk[:3],
			}
		}},
		{"a | b", func(t *test, tk Tokens) { // 8
			t.Output = Pipeline{
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
							Tokens: tk[4:5],
						},
						Tokens: tk[4:5],
					},
					Tokens: tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
		{"||", func(t *test, tk Tokens) { // 9
			t.Err = Error{
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
			}
		}},
		{"a | ||", func(t *test, tk Tokens) { // 10
			t.Err = Error{
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

func TestCommandOrCompound(t *testing.T) {
	doTests(t, []sourceFn{
		{"a", func(t *test, tk Tokens) { // 1
			t.Output = CommandOrCompound{
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
			}
		}},
		{"if a; then b; fi", func(t *test, tk Tokens) { // 2
			t.Output = CommandOrCompound{
				Compound: &Compound{
					IfCompound: &IfCompound{
						If: TestConsequence{
							Test: Statement{
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
							Consequence: File{
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
							Tokens: tk[2:9],
						},
						Tokens: tk[:11],
					},
					Tokens: tk[:11],
				},
				Tokens: tk[:11],
			}
		}},
		{"case a in b)c\nesac", func(t *test, tk Tokens) { // 3
			t.Output = CommandOrCompound{
				Compound: &Compound{
					CaseCompound: &CaseCompound{
						Word: Word{
							Parts: []WordPart{
								{
									Part:   &tk[2],
									Tokens: tk[2:3],
								},
							},
							Tokens: tk[2:3],
						},
						Matches: []PatternLines{
							{
								Patterns: []Word{
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
								Lines: File{
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
																					Part:   &tk[8],
																					Tokens: tk[8:9],
																				},
																			},
																			Tokens: tk[8:9],
																		},
																		Tokens: tk[8:9],
																	},
																},
																Tokens: tk[8:9],
															},
															Tokens: tk[8:9],
														},
														Tokens: tk[8:9],
													},
													Tokens: tk[8:9],
												},
											},
											Tokens: tk[8:9],
										},
									},
									Tokens: tk[8:9],
								},
								Tokens: tk[6:9],
							},
						},
						Tokens: tk[:11],
					},
					Tokens: tk[:11],
				},
				Tokens: tk[:11],
			}
		}},
		{"while a\ndo\nb\ndone", func(t *test, tk Tokens) { // 4
			t.Output = CommandOrCompound{
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
							Tokens: tk[2:3],
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
									Tokens: tk[6:7],
								},
							},
							Tokens: tk[6:7],
						},
						Tokens: tk[:9],
					},
					Tokens: tk[:9],
				},
				Tokens: tk[:9],
			}
		}},
		{"until a; do b; done", func(t *test, tk Tokens) { // 5
			t.Output = CommandOrCompound{
				Compound: &Compound{
					LoopCompound: &LoopCompound{
						Until: true,
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
					Tokens: tk[:11],
				},
				Tokens: tk[:11],
			}
		}},
		{"for a; do b;done", func(t *test, tk Tokens) { // 6
			t.Output = CommandOrCompound{
				Compound: &Compound{
					ForCompound: &ForCompound{
						Identifier: &tk[2],
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
						Tokens: tk[:10],
					},
					Tokens: tk[:10],
				},
				Tokens: tk[:10],
			}
		}},
		{"select a; do b;done", func(t *test, tk Tokens) { // 7
			t.Output = CommandOrCompound{
				Compound: &Compound{
					SelectCompound: &SelectCompound{
						Identifier: &tk[2],
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
						Tokens: tk[:10],
					},
					Tokens: tk[:10],
				},
				Tokens: tk[:10],
			}
		}},
		{"[[ a = b ]]", func(t *test, tk Tokens) { // 8
			t.Output = CommandOrCompound{
				Compound: &Compound{
					TestCompound: &TestCompound{
						Tests: Tests{
							Test: TestOperatorStringsEqual,
							Word: &Word{
								Parts: []WordPart{
									{
										Part:   &tk[2],
										Tokens: tk[2:3],
									},
								},
								Tokens: tk[2:3],
							},
							Pattern: &Pattern{
								Parts: []WordPart{
									{
										Part:   &tk[6],
										Tokens: tk[6:7],
									},
								},
								Tokens: tk[6:7],
							},
							Tokens: tk[2:7],
						},
						Tokens: tk[:9],
					},
					Tokens: tk[:9],
				},
				Tokens: tk[:9],
			}
		}},
		{"(a)", func(t *test, tk Tokens) { // 9
			t.Output = CommandOrCompound{
				Compound: &Compound{
					GroupingCompound: &GroupingCompound{
						SubShell: true,
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
																			Part:   &tk[1],
																			Tokens: tk[1:2],
																		},
																	},
																	Tokens: tk[1:2],
																},
																Tokens: tk[1:2],
															},
														},
														Tokens: tk[1:2],
													},
													Tokens: tk[1:2],
												},
												Tokens: tk[1:2],
											},
											Tokens: tk[1:2],
										},
									},
									Tokens: tk[1:2],
								},
							},
							Tokens: tk[1:2],
						},
						Tokens: tk[:3],
					},
					Tokens: tk[:3],
				},
				Tokens: tk[:3],
			}
		}},
		{"{\na\n}", func(t *test, tk Tokens) { // 10
			t.Output = CommandOrCompound{
				Compound: &Compound{
					GroupingCompound: &GroupingCompound{
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
							Tokens: tk[2:3],
						},
						Tokens: tk[:5],
					},
					Tokens: tk[:5],
				},
				Tokens: tk[:5],
			}
		}},
		{"function a() { b; }", func(t *test, tk Tokens) { // 11
			t.Output = CommandOrCompound{
				Compound: &Compound{
					FunctionCompound: &FunctionCompound{
						HasKeyword: true,
						Identifier: &tk[2],
						Body: Compound{
							GroupingCompound: &GroupingCompound{
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
																					Part:   &tk[8],
																					Tokens: tk[8:9],
																				},
																			},
																			Tokens: tk[8:9],
																		},
																		Tokens: tk[8:9],
																	},
																},
																Tokens: tk[8:9],
															},
															Tokens: tk[8:9],
														},
														Tokens: tk[8:9],
													},
													Tokens: tk[8:10],
												},
											},
											Tokens: tk[8:10],
										},
									},
									Tokens: tk[8:10],
								},
								Tokens: tk[6:12],
							},
							Tokens: tk[6:12],
						},
						Tokens: tk[:12],
					},
					Tokens: tk[:12],
				},
				Tokens: tk[:12],
			}
		}},
		{"a() { b; }", func(t *test, tk Tokens) { // 12
			t.Output = CommandOrCompound{
				Compound: &Compound{
					FunctionCompound: &FunctionCompound{
						Identifier: &tk[0],
						Body: Compound{
							GroupingCompound: &GroupingCompound{
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
													Tokens: tk[6:8],
												},
											},
											Tokens: tk[6:8],
										},
									},
									Tokens: tk[6:8],
								},
								Tokens: tk[4:10],
							},
							Tokens: tk[4:10],
						},
						Tokens: tk[:10],
					},
					Tokens: tk[:10],
				},
				Tokens: tk[:10],
			}
		}},
		{"(( a ))", func(t *test, tk Tokens) { // 13
			t.Output = CommandOrCompound{
				Compound: &Compound{
					ArithmeticCompound: &ArithmeticExpansion{
						Expression: true,
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
					},
					Tokens: tk[:5],
				},
				Tokens: tk[:5],
			}
		}},
		{"$(||)", func(t *test, tk Tokens) { // 14
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
															Token:   tk[1],
														},
														Parsing: "CommandOrCompound",
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
						Parsing: "AssignmentOrWord",
						Token:   tk[0],
					},
					Parsing: "Command",
					Token:   tk[0],
				},
				Parsing: "CommandOrCompound",
				Token:   tk[0],
			}
		}},
		{"case $(||) in b)c;esac", func(t *test, tk Tokens) { // 15
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
						Parsing: "CaseCompound",
						Token:   tk[2],
					},
					Parsing: "Compound",
					Token:   tk[0],
				},
				Parsing: "CommandOrCompound",
				Token:   tk[0],
			}
		}},
	}, func(t *test) (Type, error) {
		var c CommandOrCompound

		err := c.parse(t.Parser, true)

		return c, err
	})
}

func TestCompounds(t *testing.T) {
	doTests(t, []sourceFn{
		{"if a; then b; fi", func(t *test, tk Tokens) { // 1
			t.Output = Compound{
				IfCompound: &IfCompound{
					If: TestConsequence{
						Test: Statement{
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
						Consequence: File{
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
						Tokens: tk[2:9],
					},
					Tokens: tk[:11],
				},
				Tokens: tk[:11],
			}
		}},
		{"case a in b)c\nesac", func(t *test, tk Tokens) { // 2
			t.Output = Compound{
				CaseCompound: &CaseCompound{
					Word: Word{
						Parts: []WordPart{
							{
								Part:   &tk[2],
								Tokens: tk[2:3],
							},
						},
						Tokens: tk[2:3],
					},
					Matches: []PatternLines{
						{
							Patterns: []Word{
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
							Lines: File{
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
																				Part:   &tk[8],
																				Tokens: tk[8:9],
																			},
																		},
																		Tokens: tk[8:9],
																	},
																	Tokens: tk[8:9],
																},
															},
															Tokens: tk[8:9],
														},
														Tokens: tk[8:9],
													},
													Tokens: tk[8:9],
												},
												Tokens: tk[8:9],
											},
										},
										Tokens: tk[8:9],
									},
								},
								Tokens: tk[8:9],
							},
							Tokens: tk[6:9],
						},
					},
					Tokens: tk[:11],
				},
				Tokens: tk[:11],
			}
		}},
		{"while a\ndo\nb\ndone", func(t *test, tk Tokens) { // 3
			t.Output = Compound{
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
						Tokens: tk[2:3],
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
								Tokens: tk[6:7],
							},
						},
						Tokens: tk[6:7],
					},
					Tokens: tk[:9],
				},
				Tokens: tk[:9],
			}
		}},
		{"until a; do b; done", func(t *test, tk Tokens) { // 4
			t.Output = Compound{
				LoopCompound: &LoopCompound{
					Until: true,
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
				Tokens: tk[:11],
			}
		}},
		{"for a; do b;done", func(t *test, tk Tokens) { // 5
			t.Output = Compound{
				ForCompound: &ForCompound{
					Identifier: &tk[2],
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
					Tokens: tk[:10],
				},
				Tokens: tk[:10],
			}
		}},
		{"select a; do b;done", func(t *test, tk Tokens) { // 6
			t.Output = Compound{
				SelectCompound: &SelectCompound{
					Identifier: &tk[2],
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
					Tokens: tk[:10],
				},
				Tokens: tk[:10],
			}
		}},
		{"[[ a = b ]]", func(t *test, tk Tokens) { // 7
			t.Output = Compound{
				TestCompound: &TestCompound{
					Tests: Tests{
						Test: TestOperatorStringsEqual,
						Word: &Word{
							Parts: []WordPart{
								{
									Part:   &tk[2],
									Tokens: tk[2:3],
								},
							},
							Tokens: tk[2:3],
						},
						Pattern: &Pattern{
							Parts: []WordPart{
								{
									Part:   &tk[6],
									Tokens: tk[6:7],
								},
							},
							Tokens: tk[6:7],
						},
						Tokens: tk[2:7],
					},
					Tokens: tk[:9],
				},
				Tokens: tk[:9],
			}
		}},
		{"(a)", func(t *test, tk Tokens) { // 8
			t.Output = Compound{
				GroupingCompound: &GroupingCompound{
					SubShell: true,
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
																		Part:   &tk[1],
																		Tokens: tk[1:2],
																	},
																},
																Tokens: tk[1:2],
															},
															Tokens: tk[1:2],
														},
													},
													Tokens: tk[1:2],
												},
												Tokens: tk[1:2],
											},
											Tokens: tk[1:2],
										},
										Tokens: tk[1:2],
									},
								},
								Tokens: tk[1:2],
							},
						},
						Tokens: tk[1:2],
					},
					Tokens: tk[:3],
				},
				Tokens: tk[:3],
			}
		}},
		{"{\na\n}", func(t *test, tk Tokens) { // 9
			t.Output = Compound{
				GroupingCompound: &GroupingCompound{
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
						Tokens: tk[2:3],
					},
					Tokens: tk[:5],
				},
				Tokens: tk[:5],
			}
		}},
		{"function a() { b; }", func(t *test, tk Tokens) { // 10
			t.Output = Compound{
				FunctionCompound: &FunctionCompound{
					HasKeyword: true,
					Identifier: &tk[2],
					Body: Compound{
						GroupingCompound: &GroupingCompound{
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
																				Part:   &tk[8],
																				Tokens: tk[8:9],
																			},
																		},
																		Tokens: tk[8:9],
																	},
																	Tokens: tk[8:9],
																},
															},
															Tokens: tk[8:9],
														},
														Tokens: tk[8:9],
													},
													Tokens: tk[8:9],
												},
												Tokens: tk[8:10],
											},
										},
										Tokens: tk[8:10],
									},
								},
								Tokens: tk[8:10],
							},
							Tokens: tk[6:12],
						},
						Tokens: tk[6:12],
					},
					Tokens: tk[:12],
				},
				Tokens: tk[:12],
			}
		}},
		{"a() { b; }", func(t *test, tk Tokens) { // 11
			t.Output = Compound{
				FunctionCompound: &FunctionCompound{
					Identifier: &tk[0],
					Body: Compound{
						GroupingCompound: &GroupingCompound{
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
												Tokens: tk[6:8],
											},
										},
										Tokens: tk[6:8],
									},
								},
								Tokens: tk[6:8],
							},
							Tokens: tk[4:10],
						},
						Tokens: tk[4:10],
					},
					Tokens: tk[:10],
				},
				Tokens: tk[:10],
			}
		}},
		{"(( a ))", func(t *test, tk Tokens) { // 12
			t.Output = Compound{
				ArithmeticCompound: &ArithmeticExpansion{
					Expression: true,
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
				},
				Tokens: tk[:5],
			}
		}},
		{"[[ a = b ]] >c", func(t *test, tk Tokens) { // 13
			t.Output = Compound{
				TestCompound: &TestCompound{
					Tests: Tests{
						Test: TestOperatorStringsEqual,
						Word: &Word{
							Parts: []WordPart{
								{
									Part:   &tk[2],
									Tokens: tk[2:3],
								},
							},
							Tokens: tk[2:3],
						},
						Pattern: &Pattern{
							Parts: []WordPart{
								{
									Part:   &tk[6],
									Tokens: tk[6:7],
								},
							},
							Tokens: tk[6:7],
						},
						Tokens: tk[2:7],
					},
					Tokens: tk[:9],
				},
				Redirections: []Redirection{
					{
						Redirector: &tk[10],
						Output: Word{
							Parts: []WordPart{
								{
									Part:   &tk[11],
									Tokens: tk[11:12],
								},
							},
							Tokens: tk[11:12],
						},
						Tokens: tk[10:12],
					},
				},
				Tokens: tk[:12],
			}
		}},
		{"[[ a = b ]] >a 2>&1", func(t *test, tk Tokens) { // 14
			t.Output = Compound{
				TestCompound: &TestCompound{
					Tests: Tests{
						Test: TestOperatorStringsEqual,
						Word: &Word{
							Parts: []WordPart{
								{
									Part:   &tk[2],
									Tokens: tk[2:3],
								},
							},
							Tokens: tk[2:3],
						},
						Pattern: &Pattern{
							Parts: []WordPart{
								{
									Part:   &tk[6],
									Tokens: tk[6:7],
								},
							},
							Tokens: tk[6:7],
						},
						Tokens: tk[2:7],
					},
					Tokens: tk[:9],
				},
				Redirections: []Redirection{
					{
						Redirector: &tk[10],
						Output: Word{
							Parts: []WordPart{
								{
									Part:   &tk[11],
									Tokens: tk[11:12],
								},
							},
							Tokens: tk[11:12],
						},
						Tokens: tk[10:12],
					},
					{
						Input:      &tk[13],
						Redirector: &tk[14],
						Output: Word{
							Parts: []WordPart{
								{
									Part:   &tk[15],
									Tokens: tk[15:16],
								},
							},
							Tokens: tk[15:16],
						},
						Tokens: tk[13:16],
					},
				},
				Tokens: tk[:16],
			}
		}},
		{"if ||;then b;fi", func(t *test, tk Tokens) { // 15
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
									Parsing: "CommandOrCompound",
									Token:   tk[2],
								},
								Parsing: "Pipeline",
								Token:   tk[2],
							},
							Parsing: "Statement",
							Token:   tk[2],
						},
						Parsing: "TestConsequence",
						Token:   tk[2],
					},
					Parsing: "IfCompound",
					Token:   tk[2],
				},
				Parsing: "Compound",
				Token:   tk[0],
			}
		}},
		{"case $(||) in b)c;esac", func(t *test, tk Tokens) { // 16
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
					Parsing: "CaseCompound",
					Token:   tk[2],
				},
				Parsing: "Compound",
				Token:   tk[0],
			}
		}},
		{"while ||; do b; done", func(t *test, tk Tokens) { // 17
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err: Error{
									Err:     ErrMissingWord,
									Parsing: "Command",
									Token:   tk[2],
								},
								Parsing: "CommandOrCompound",
								Token:   tk[2],
							},
							Parsing: "Pipeline",
							Token:   tk[2],
						},
						Parsing: "Statement",
						Token:   tk[2],
					},
					Parsing: "LoopCompound",
					Token:   tk[2],
				},
				Parsing: "Compound",
				Token:   tk[0],
			}
		}},
		{"until a; do ||; done", func(t *test, tk Tokens) { // 18
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
											Token:   tk[7],
										},
										Parsing: "CommandOrCompound",
										Token:   tk[7],
									},
									Parsing: "Pipeline",
									Token:   tk[7],
								},
								Parsing: "Statement",
								Token:   tk[7],
							},
							Parsing: "Line",
							Token:   tk[7],
						},
						Parsing: "File",
						Token:   tk[7],
					},
					Parsing: "LoopCompound",
					Token:   tk[7],
				},
				Parsing: "Compound",
				Token:   tk[0],
			}
		}},
		{"for a in $(||); do b;done", func(t *test, tk Tokens) { // 19
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
														Token:   tk[7],
													},
													Parsing: "CommandOrCompound",
													Token:   tk[7],
												},
												Parsing: "Pipeline",
												Token:   tk[7],
											},
											Parsing: "Statement",
											Token:   tk[7],
										},
										Parsing: "Line",
										Token:   tk[7],
									},
									Parsing: "File",
									Token:   tk[7],
								},
								Parsing: "CommandSubstitution",
								Token:   tk[7],
							},
							Parsing: "WordPart",
							Token:   tk[6],
						},
						Parsing: "Word",
						Token:   tk[6],
					},
					Parsing: "ForCompound",
					Token:   tk[6],
				},
				Parsing: "Compound",
				Token:   tk[0],
			}
		}},
		{"select a in $(||); do b;done", func(t *test, tk Tokens) { // 20
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
														Token:   tk[7],
													},
													Parsing: "CommandOrCompound",
													Token:   tk[7],
												},
												Parsing: "Pipeline",
												Token:   tk[7],
											},
											Parsing: "Statement",
											Token:   tk[7],
										},
										Parsing: "Line",
										Token:   tk[7],
									},
									Parsing: "File",
									Token:   tk[7],
								},
								Parsing: "CommandSubstitution",
								Token:   tk[7],
							},
							Parsing: "WordPart",
							Token:   tk[6],
						},
						Parsing: "Word",
						Token:   tk[6],
					},
					Parsing: "SelectCompound",
					Token:   tk[6],
				},
				Parsing: "Compound",
				Token:   tk[0],
			}
		}},
		{"[[ -a $(||) ]]", func(t *test, tk Tokens) { // 21
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
										},
										Parsing: "File",
										Token:   tk[5],
									},
									Parsing: "CommandSubstitution",
									Token:   tk[5],
								},
								Parsing: "WordPart",
								Token:   tk[4],
							},
							Parsing: "Word",
							Token:   tk[4],
						},
						Parsing: "Tests",
						Token:   tk[4],
					},
					Parsing: "TestCompound",
					Token:   tk[2],
				},
				Parsing: "Compound",
				Token:   tk[0],
			}
		}},
		{"(||)", func(t *test, tk Tokens) { // 22
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
										Parsing: "CommandOrCompound",
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
					Parsing: "GroupingCompound",
					Token:   tk[1],
				},
				Parsing: "Compound",
				Token:   tk[0],
			}
		}},
		{"function a() { || }", func(t *test, tk Tokens) { // 23
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
													Token:   tk[8],
												},
												Parsing: "CommandOrCompound",
												Token:   tk[8],
											},
											Parsing: "Pipeline",
											Token:   tk[8],
										},
										Parsing: "Statement",
										Token:   tk[8],
									},
									Parsing: "Line",
									Token:   tk[8],
								},
								Parsing: "File",
								Token:   tk[8],
							},
							Parsing: "GroupingCompound",
							Token:   tk[8],
						},
						Parsing: "Compound",
						Token:   tk[6],
					},
					Parsing: "FunctionCompound",
					Token:   tk[6],
				},
				Parsing: "Compound",
				Token:   tk[0],
			}
		}},
		{"(($(||)))", func(t *test, tk Tokens) { // 24
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
														Parsing: "CommandOrCompound",
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
				Parsing: "Compound",
				Token:   tk[0],
			}
		}},
		{"[[ a = b ]] >$(||)", func(t *test, tk Tokens) { // 25
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
														Token:   tk[12],
													},
													Parsing: "CommandOrCompound",
													Token:   tk[12],
												},
												Parsing: "Pipeline",
												Token:   tk[12],
											},
											Parsing: "Statement",
											Token:   tk[12],
										},
										Parsing: "Line",
										Token:   tk[12],
									},
									Parsing: "File",
									Token:   tk[12],
								},
								Parsing: "CommandSubstitution",
								Token:   tk[12],
							},
							Parsing: "WordPart",
							Token:   tk[11],
						},
						Parsing: "Word",
						Token:   tk[11],
					},
					Parsing: "Redirection",
					Token:   tk[11],
				},
				Parsing: "Compound",
				Token:   tk[10],
			}
		}},
	}, func(t *test) (Type, error) {
		var c Compound

		err := c.parse(t.Parser)

		return c, err
	})
}

func TestIfCompound(t *testing.T) {
	doTests(t, []sourceFn{
		{"if a; then b; fi", func(t *test, tk Tokens) { // 1
			t.Output = IfCompound{
				If: TestConsequence{
					Test: Statement{
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
					Consequence: File{
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
					Tokens: tk[2:9],
				},
				Tokens: tk[:11],
			}
		}},
		{"if\na\nthen\nb\nfi", func(t *test, tk Tokens) { // 2
			t.Output = IfCompound{
				If: TestConsequence{
					Test: Statement{
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
					Consequence: File{
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
								Tokens: tk[6:7],
							},
						},
						Tokens: tk[6:7],
					},
					Tokens: tk[2:7],
				},
				Tokens: tk[:9],
			}
		}},
		{"if a; then b; else c; fi", func(t *test, tk Tokens) { // 3
			t.Output = IfCompound{
				If: TestConsequence{
					Test: Statement{
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
					Consequence: File{
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
					Tokens: tk[2:9],
				},
				Else: &File{
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
																	Part:   &tk[12],
																	Tokens: tk[12:13],
																},
															},
															Tokens: tk[12:13],
														},
														Tokens: tk[12:13],
													},
												},
												Tokens: tk[12:13],
											},
											Tokens: tk[12:13],
										},
										Tokens: tk[12:13],
									},
									Tokens: tk[12:14],
								},
							},
							Tokens: tk[12:14],
						},
					},
					Tokens: tk[12:14],
				},
				Tokens: tk[:16],
			}
		}},
		{"if a; then b; elif c; then d; fi", func(t *test, tk Tokens) { // 4
			t.Output = IfCompound{
				If: TestConsequence{
					Test: Statement{
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
					Consequence: File{
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
					Tokens: tk[2:9],
				},
				ElIf: []TestConsequence{
					{
						Test: Statement{
							Pipeline: Pipeline{
								CommandOrCompound: CommandOrCompound{
									Command: &Command{
										AssignmentsOrWords: []AssignmentOrWord{
											{
												Word: &Word{
													Parts: []WordPart{
														{
															Part:   &tk[12],
															Tokens: tk[12:13],
														},
													},
													Tokens: tk[12:13],
												},
												Tokens: tk[12:13],
											},
										},
										Tokens: tk[12:13],
									},
									Tokens: tk[12:13],
								},
								Tokens: tk[12:13],
							},
							Tokens: tk[12:14],
						},
						Consequence: File{
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
																			Part:   &tk[17],
																			Tokens: tk[17:18],
																		},
																	},
																	Tokens: tk[17:18],
																},
																Tokens: tk[17:18],
															},
														},
														Tokens: tk[17:18],
													},
													Tokens: tk[17:18],
												},
												Tokens: tk[17:18],
											},
											Tokens: tk[17:19],
										},
									},
									Tokens: tk[17:19],
								},
							},
							Tokens: tk[17:19],
						},
						Tokens: tk[12:19],
					},
				},
				Tokens: tk[:21],
			}
		}},
		{"if a; then b; else\n# comment\nc; fi", func(t *test, tk Tokens) { // 3
			t.Output = IfCompound{
				If: TestConsequence{
					Test: Statement{
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
					Consequence: File{
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
					Tokens: tk[2:9],
				},
				Else: &File{
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
																	Part:   &tk[14],
																	Tokens: tk[14:15],
																},
															},
															Tokens: tk[14:15],
														},
														Tokens: tk[14:15],
													},
												},
												Tokens: tk[14:15],
											},
											Tokens: tk[14:15],
										},
										Tokens: tk[14:15],
									},
									Tokens: tk[14:16],
								},
							},
							Comments: [2]Comments{{tk[12]}},
							Tokens:   tk[12:16],
						},
					},
					Tokens: tk[12:16],
				},
				Tokens: tk[:18],
			}
		}},
		{"if a; then b; else # comment\nc; fi", func(t *test, tk Tokens) { // 3
			t.Output = IfCompound{
				If: TestConsequence{
					Test: Statement{
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
					Consequence: File{
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
					Tokens: tk[2:9],
				},
				Else: &File{
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
																	Part:   &tk[14],
																	Tokens: tk[14:15],
																},
															},
															Tokens: tk[14:15],
														},
														Tokens: tk[14:15],
													},
												},
												Tokens: tk[14:15],
											},
											Tokens: tk[14:15],
										},
										Tokens: tk[14:15],
									},
									Tokens: tk[14:16],
								},
							},
							Tokens: tk[14:16],
						},
					},
					Comments: [2]Comments{{tk[12]}},
					Tokens:   tk[12:16],
				},
				Tokens: tk[:18],
			}
		}},
		{"if ||;then b;fi", func(t *test, tk Tokens) { // 5
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err: Error{
									Err:     ErrMissingWord,
									Parsing: "Command",
									Token:   tk[2],
								},
								Parsing: "CommandOrCompound",
								Token:   tk[2],
							},
							Parsing: "Pipeline",
							Token:   tk[2],
						},
						Parsing: "Statement",
						Token:   tk[2],
					},
					Parsing: "TestConsequence",
					Token:   tk[2],
				},
				Parsing: "IfCompound",
				Token:   tk[2],
			}
		}},
		{"if a;then b;elif ||;then d;fi", func(t *test, tk Tokens) { // 6
			t.Err = Error{
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
					Parsing: "TestConsequence",
					Token:   tk[10],
				},
				Parsing: "IfCompound",
				Token:   tk[10],
			}
		}},
		{"if a;then b;else ||;fi", func(t *test, tk Tokens) { // 7
			t.Err = Error{
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
				Parsing: "IfCompound",
				Token:   tk[10],
			}
		}},
	}, func(t *test) (Type, error) {
		var i IfCompound

		err := i.parse(t.Parser)

		return i, err
	})
}

func TestTestConsequence(t *testing.T) {
	doTests(t, []sourceFn{
		{"if a; then b;fi", func(t *test, tk Tokens) { // 1
			t.Output = TestConsequence{
				Test: Statement{
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
				Consequence: File{
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
				Tokens: tk[2:9],
			}
		}},
		{"if a\nthen\nb\nc\nfi", func(t *test, tk Tokens) { // 2
			t.Output = TestConsequence{
				Test: Statement{
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
				Consequence: File{
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
							Tokens: tk[6:7],
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
																	Part:   &tk[8],
																	Tokens: tk[8:9],
																},
															},
															Tokens: tk[8:9],
														},
														Tokens: tk[8:9],
													},
												},
												Tokens: tk[8:9],
											},
											Tokens: tk[8:9],
										},
										Tokens: tk[8:9],
									},
									Tokens: tk[8:9],
								},
							},
							Tokens: tk[8:9],
						},
					},
					Tokens: tk[6:9],
				},
				Tokens: tk[2:9],
			}
		}},
		{"if a; #comment\nthen b;fi", func(t *test, tk Tokens) { // 3
			t.Output = TestConsequence{
				Test: Statement{
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
				Consequence: File{
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
																	Part:   &tk[9],
																	Tokens: tk[9:10],
																},
															},
															Tokens: tk[9:10],
														},
														Tokens: tk[9:10],
													},
												},
												Tokens: tk[9:10],
											},
											Tokens: tk[9:10],
										},
										Tokens: tk[9:10],
									},
									Tokens: tk[9:11],
								},
							},
							Tokens: tk[9:11],
						},
					},
					Tokens: tk[9:11],
				},
				Comments: Comments{tk[5]},
				Tokens:   tk[2:11],
			}
		}},
		{"if a\n#comment\nthen b;fi", func(t *test, tk Tokens) { // 4
			t.Output = TestConsequence{
				Test: Statement{
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
				Consequence: File{
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
																	Part:   &tk[8],
																	Tokens: tk[8:9],
																},
															},
															Tokens: tk[8:9],
														},
														Tokens: tk[8:9],
													},
												},
												Tokens: tk[8:9],
											},
											Tokens: tk[8:9],
										},
										Tokens: tk[8:9],
									},
									Tokens: tk[8:10],
								},
							},
							Tokens: tk[8:10],
						},
					},
					Tokens: tk[8:10],
				},
				Comments: Comments{tk[4]},
				Tokens:   tk[2:10],
			}
		}},
		{"if ||; then b;fi", func(t *test, tk Tokens) { // 5
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err:     ErrMissingWord,
								Parsing: "Command",
								Token:   tk[2],
							},
							Parsing: "CommandOrCompound",
							Token:   tk[2],
						},
						Parsing: "Pipeline",
						Token:   tk[2],
					},
					Parsing: "Statement",
					Token:   tk[2],
				},
				Parsing: "TestConsequence",
				Token:   tk[2],
			}
		}},
		{"if a; then ||;fi", func(t *test, tk Tokens) { // 6
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err: Error{
									Err: Error{
										Err:     ErrMissingWord,
										Parsing: "Command",
										Token:   tk[7],
									},
									Parsing: "CommandOrCompound",
									Token:   tk[7],
								},
								Parsing: "Pipeline",
								Token:   tk[7],
							},
							Parsing: "Statement",
							Token:   tk[7],
						},
						Parsing: "Line",
						Token:   tk[7],
					},
					Parsing: "File",
					Token:   tk[7],
				},
				Parsing: "TestConsequence",
				Token:   tk[7],
			}
		}},
	}, func(t *test) (Type, error) {
		var tc TestConsequence

		t.Parser.Tokens = t.Parser.Tokens[2:2]
		err := tc.parse(t.Parser)

		return tc, err
	})
}

func TestCaseCompound(t *testing.T) {
	doTests(t, []sourceFn{
		{"case a in b)c\nesac", func(t *test, tk Tokens) { // 1
			t.Output = CaseCompound{
				Word: Word{
					Parts: []WordPart{
						{
							Part:   &tk[2],
							Tokens: tk[2:3],
						},
					},
					Tokens: tk[2:3],
				},
				Matches: []PatternLines{
					{
						Patterns: []Word{
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
						Lines: File{
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
																			Part:   &tk[8],
																			Tokens: tk[8:9],
																		},
																	},
																	Tokens: tk[8:9],
																},
																Tokens: tk[8:9],
															},
														},
														Tokens: tk[8:9],
													},
													Tokens: tk[8:9],
												},
												Tokens: tk[8:9],
											},
											Tokens: tk[8:9],
										},
									},
									Tokens: tk[8:9],
								},
							},
							Tokens: tk[8:9],
						},
						Tokens: tk[6:9],
					},
				},
				Tokens: tk[:11],
			}
		}},
		{"case a in b)c;esac", func(t *test, tk Tokens) { // 2
			t.Output = CaseCompound{
				Word: Word{
					Parts: []WordPart{
						{
							Part:   &tk[2],
							Tokens: tk[2:3],
						},
					},
					Tokens: tk[2:3],
				},
				Matches: []PatternLines{
					{
						Patterns: []Word{
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
						Lines: File{
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
																			Part:   &tk[8],
																			Tokens: tk[8:9],
																		},
																	},
																	Tokens: tk[8:9],
																},
																Tokens: tk[8:9],
															},
														},
														Tokens: tk[8:9],
													},
													Tokens: tk[8:9],
												},
												Tokens: tk[8:9],
											},
											Tokens: tk[8:10],
										},
									},
									Tokens: tk[8:10],
								},
							},
							Tokens: tk[8:10],
						},
						Tokens: tk[6:10],
					},
				},
				Tokens: tk[:11],
			}
		}},
		{"case a in b)c;;esac", func(t *test, tk Tokens) { // 3
			t.Output = CaseCompound{
				Word: Word{
					Parts: []WordPart{
						{
							Part:   &tk[2],
							Tokens: tk[2:3],
						},
					},
					Tokens: tk[2:3],
				},
				Matches: []PatternLines{
					{
						Patterns: []Word{
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
						Lines: File{
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
																			Part:   &tk[8],
																			Tokens: tk[8:9],
																		},
																	},
																	Tokens: tk[8:9],
																},
																Tokens: tk[8:9],
															},
														},
														Tokens: tk[8:9],
													},
													Tokens: tk[8:9],
												},
												Tokens: tk[8:9],
											},
											Tokens: tk[8:9],
										},
									},
									Tokens: tk[8:9],
								},
							},
							Tokens: tk[8:9],
						},
						CaseTerminationType: CaseTerminationEnd,
						Tokens:              tk[6:10],
					},
				},
				Tokens: tk[:11],
			}
		}},
		{"case a in b)c;;\nd|e)f;;esac", func(t *test, tk Tokens) { // 4
			t.Output = CaseCompound{
				Word: Word{
					Parts: []WordPart{
						{
							Part:   &tk[2],
							Tokens: tk[2:3],
						},
					},
					Tokens: tk[2:3],
				},
				Matches: []PatternLines{
					{
						Patterns: []Word{
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
						Lines: File{
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
																			Part:   &tk[8],
																			Tokens: tk[8:9],
																		},
																	},
																	Tokens: tk[8:9],
																},
																Tokens: tk[8:9],
															},
														},
														Tokens: tk[8:9],
													},
													Tokens: tk[8:9],
												},
												Tokens: tk[8:9],
											},
											Tokens: tk[8:9],
										},
									},
									Tokens: tk[8:9],
								},
							},
							Tokens: tk[8:9],
						},
						CaseTerminationType: CaseTerminationEnd,
						Tokens:              tk[6:10],
					},
					{
						Patterns: []Word{
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
								},
								Tokens: tk[13:14],
							},
						},
						Lines: File{
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
																			Part:   &tk[15],
																			Tokens: tk[15:16],
																		},
																	},
																	Tokens: tk[15:16],
																},
																Tokens: tk[15:16],
															},
														},
														Tokens: tk[15:16],
													},
													Tokens: tk[15:16],
												},
												Tokens: tk[15:16],
											},
											Tokens: tk[15:16],
										},
									},
									Tokens: tk[15:16],
								},
							},
							Tokens: tk[15:16],
						},
						CaseTerminationType: CaseTerminationEnd,
						Tokens:              tk[11:17],
					},
				},
				Tokens: tk[:18],
			}
		}},
		{"case $(||) in b)c;esac", func(t *test, tk Tokens) { // 5
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
				Parsing: "CaseCompound",
				Token:   tk[2],
			}
		}},
		{"case a in b)||;esac", func(t *test, tk Tokens) { // 6
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
											Token:   tk[8],
										},
										Parsing: "CommandOrCompound",
										Token:   tk[8],
									},
									Parsing: "Pipeline",
									Token:   tk[8],
								},
								Parsing: "Statement",
								Token:   tk[8],
							},
							Parsing: "Line",
							Token:   tk[8],
						},
						Parsing: "File",
						Token:   tk[8],
					},
					Parsing: "PatternLines",
					Token:   tk[8],
				},
				Parsing: "CaseCompound",
				Token:   tk[6],
			}
		}},
	}, func(t *test) (Type, error) {
		var cc CaseCompound

		err := cc.parse(t.Parser)

		return cc, err
	})
}

func TestPatternLines(t *testing.T) {
	doTests(t, []sourceFn{
		{"case a in a)b\nesac", func(t *test, tk Tokens) { // 1
			t.Output = PatternLines{
				Patterns: []Word{
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
				Lines: File{
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
																	Part:   &tk[8],
																	Tokens: tk[8:9],
																},
															},
															Tokens: tk[8:9],
														},
														Tokens: tk[8:9],
													},
												},
												Tokens: tk[8:9],
											},
											Tokens: tk[8:9],
										},
										Tokens: tk[8:9],
									},
									Tokens: tk[8:9],
								},
							},
							Tokens: tk[8:9],
						},
					},
					Tokens: tk[8:9],
				},
				Tokens: tk[6:9],
			}
		}},
		{"case a in a)b;esac", func(t *test, tk Tokens) { // 2
			t.Output = PatternLines{
				Patterns: []Word{
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
				Lines: File{
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
																	Part:   &tk[8],
																	Tokens: tk[8:9],
																},
															},
															Tokens: tk[8:9],
														},
														Tokens: tk[8:9],
													},
												},
												Tokens: tk[8:9],
											},
											Tokens: tk[8:9],
										},
										Tokens: tk[8:9],
									},
									Tokens: tk[8:10],
								},
							},
							Tokens: tk[8:10],
						},
					},
					Tokens: tk[8:10],
				},
				Tokens: tk[6:10],
			}
		}},
		{"case a in a)b;;esac", func(t *test, tk Tokens) { // 3
			t.Output = PatternLines{
				Patterns: []Word{
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
				Lines: File{
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
																	Part:   &tk[8],
																	Tokens: tk[8:9],
																},
															},
															Tokens: tk[8:9],
														},
														Tokens: tk[8:9],
													},
												},
												Tokens: tk[8:9],
											},
											Tokens: tk[8:9],
										},
										Tokens: tk[8:9],
									},
									Tokens: tk[8:9],
								},
							},
							Tokens: tk[8:9],
						},
					},
					Tokens: tk[8:9],
				},
				CaseTerminationType: CaseTerminationEnd,
				Tokens:              tk[6:10],
			}
		}},
		{"case a in a)b;&esac", func(t *test, tk Tokens) { // 4
			t.Output = PatternLines{
				Patterns: []Word{
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
				Lines: File{
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
																	Part:   &tk[8],
																	Tokens: tk[8:9],
																},
															},
															Tokens: tk[8:9],
														},
														Tokens: tk[8:9],
													},
												},
												Tokens: tk[8:9],
											},
											Tokens: tk[8:9],
										},
										Tokens: tk[8:9],
									},
									Tokens: tk[8:9],
								},
							},
							Tokens: tk[8:9],
						},
					},
					Tokens: tk[8:9],
				},
				CaseTerminationType: CaseTerminationContinue,
				Tokens:              tk[6:10],
			}
		}},
		{"case a in a)b;;&esac", func(t *test, tk Tokens) { // 5
			t.Output = PatternLines{
				Patterns: []Word{
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
				Lines: File{
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
																	Part:   &tk[8],
																	Tokens: tk[8:9],
																},
															},
															Tokens: tk[8:9],
														},
														Tokens: tk[8:9],
													},
												},
												Tokens: tk[8:9],
											},
											Tokens: tk[8:9],
										},
										Tokens: tk[8:9],
									},
									Tokens: tk[8:9],
								},
							},
							Tokens: tk[8:9],
						},
					},
					Tokens: tk[8:9],
				},
				CaseTerminationType: CaseTerminationFallthrough,
				Tokens:              tk[6:10],
			}
		}},
		{"case a in a|b|c)d;;esac", func(t *test, tk Tokens) { // 6
			t.Output = PatternLines{
				Patterns: []Word{
					{
						Parts: []WordPart{
							{
								Part:   &tk[6],
								Tokens: tk[6:7],
							},
						},
						Tokens: tk[6:7],
					},
					{
						Parts: []WordPart{
							{
								Part:   &tk[8],
								Tokens: tk[8:9],
							},
						},
						Tokens: tk[8:9],
					},
					{
						Parts: []WordPart{
							{
								Part:   &tk[10],
								Tokens: tk[10:11],
							},
						},
						Tokens: tk[10:11],
					},
				},
				Lines: File{
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
																	Part:   &tk[12],
																	Tokens: tk[12:13],
																},
															},
															Tokens: tk[12:13],
														},
														Tokens: tk[12:13],
													},
												},
												Tokens: tk[12:13],
											},
											Tokens: tk[12:13],
										},
										Tokens: tk[12:13],
									},
									Tokens: tk[12:13],
								},
							},
							Tokens: tk[12:13],
						},
					},
					Tokens: tk[12:13],
				},
				CaseTerminationType: CaseTerminationEnd,
				Tokens:              tk[6:14],
			}
		}},
		{"case a in $(||))d;;esac", func(t *test, tk Tokens) { // 7
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
													Token:   tk[7],
												},
												Parsing: "CommandOrCompound",
												Token:   tk[7],
											},
											Parsing: "Pipeline",
											Token:   tk[7],
										},
										Parsing: "Statement",
										Token:   tk[7],
									},
									Parsing: "Line",
									Token:   tk[7],
								},
								Parsing: "File",
								Token:   tk[7],
							},
							Parsing: "CommandSubstitution",
							Token:   tk[7],
						},
						Parsing: "WordPart",
						Token:   tk[6],
					},
					Parsing: "Word",
					Token:   tk[6],
				},
				Parsing: "PatternLines",
				Token:   tk[6],
			}
		}},
		{"case a in a|\nb)c;;esac", func(t *test, tk Tokens) { // 8
			t.Err = Error{
				Err:     ErrMissingClosingPattern,
				Parsing: "PatternLines",
				Token:   tk[8],
			}
		}},
		{"case a in a)||;;esac", func(t *test, tk Tokens) { // 9
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err: Error{
									Err: Error{
										Err:     ErrMissingWord,
										Parsing: "Command",
										Token:   tk[8],
									},
									Parsing: "CommandOrCompound",
									Token:   tk[8],
								},
								Parsing: "Pipeline",
								Token:   tk[8],
							},
							Parsing: "Statement",
							Token:   tk[8],
						},
						Parsing: "Line",
						Token:   tk[8],
					},
					Parsing: "File",
					Token:   tk[8],
				},
				Parsing: "PatternLines",
				Token:   tk[8],
			}
		}},
	}, func(t *test) (Type, error) {
		var p PatternLines

		t.Parser.Tokens = t.Parser.Tokens[6:6]
		err := p.parse(t.Parser)

		return p, err
	})
}

func TestLoopCompound(t *testing.T) {
	doTests(t, []sourceFn{
		{"while a\ndo\nb\ndone", func(t *test, tk Tokens) { // 1
			t.Output = LoopCompound{
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
					Tokens: tk[2:3],
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
							Tokens: tk[6:7],
						},
					},
					Tokens: tk[6:7],
				},
				Tokens: tk[:9],
			}
		}},
		{"until a; do b; done", func(t *test, tk Tokens) { // 2
			t.Output = LoopCompound{
				Until: true,
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
			}
		}},
		{"while ||; do b; done", func(t *test, tk Tokens) { // 3
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err:     ErrMissingWord,
								Parsing: "Command",
								Token:   tk[2],
							},
							Parsing: "CommandOrCompound",
							Token:   tk[2],
						},
						Parsing: "Pipeline",
						Token:   tk[2],
					},
					Parsing: "Statement",
					Token:   tk[2],
				},
				Parsing: "LoopCompound",
				Token:   tk[2],
			}
		}},
		{"until a; do ||; done", func(t *test, tk Tokens) { // 4
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err: Error{
									Err: Error{
										Err:     ErrMissingWord,
										Parsing: "Command",
										Token:   tk[7],
									},
									Parsing: "CommandOrCompound",
									Token:   tk[7],
								},
								Parsing: "Pipeline",
								Token:   tk[7],
							},
							Parsing: "Statement",
							Token:   tk[7],
						},
						Parsing: "Line",
						Token:   tk[7],
					},
					Parsing: "File",
					Token:   tk[7],
				},
				Parsing: "LoopCompound",
				Token:   tk[7],
			}
		}},
	}, func(t *test) (Type, error) {
		var l LoopCompound

		err := l.parse(t.Parser)

		return l, err
	})
}

func TestForCompound(t *testing.T) {
	doTests(t, []sourceFn{
		{"for a; do b;done", func(t *test, tk Tokens) { // 1
			t.Output = ForCompound{
				Identifier: &tk[2],
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
				Tokens: tk[:10],
			}
		}},
		{"for a\ndo b\ndone", func(t *test, tk Tokens) { // 2
			t.Output = ForCompound{
				Identifier: &tk[2],
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
							Tokens: tk[6:7],
						},
					},
					Tokens: tk[6:7],
				},
				Tokens: tk[:9],
			}
		}},
		{"for a in; do b;done", func(t *test, tk Tokens) { // 3
			t.Output = ForCompound{
				Identifier: &tk[2],
				Words:      []Word{},
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
																	Part:   &tk[9],
																	Tokens: tk[9:10],
																},
															},
															Tokens: tk[9:10],
														},
														Tokens: tk[9:10],
													},
												},
												Tokens: tk[9:10],
											},
											Tokens: tk[9:10],
										},
										Tokens: tk[9:10],
									},
									Tokens: tk[9:11],
								},
							},
							Tokens: tk[9:11],
						},
					},
					Tokens: tk[9:11],
				},
				Tokens: tk[:12],
			}
		}},
		{"for a in b; do c;done", func(t *test, tk Tokens) { // 4
			t.Output = ForCompound{
				Identifier: &tk[2],
				Words: []Word{
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
																	Part:   &tk[11],
																	Tokens: tk[11:12],
																},
															},
															Tokens: tk[11:12],
														},
														Tokens: tk[11:12],
													},
												},
												Tokens: tk[11:12],
											},
											Tokens: tk[11:12],
										},
										Tokens: tk[11:12],
									},
									Tokens: tk[11:13],
								},
							},
							Tokens: tk[11:13],
						},
					},
					Tokens: tk[11:13],
				},
				Tokens: tk[:14],
			}
		}},
		{"for a in b c; do d;done", func(t *test, tk Tokens) { // 5
			t.Output = ForCompound{
				Identifier: &tk[2],
				Words: []Word{
					{
						Parts: []WordPart{
							{
								Part:   &tk[6],
								Tokens: tk[6:7],
							},
						},
						Tokens: tk[6:7],
					},
					{
						Parts: []WordPart{
							{
								Part:   &tk[8],
								Tokens: tk[8:9],
							},
						},
						Tokens: tk[8:9],
					},
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
																	Part:   &tk[13],
																	Tokens: tk[13:14],
																},
															},
															Tokens: tk[13:14],
														},
														Tokens: tk[13:14],
													},
												},
												Tokens: tk[13:14],
											},
											Tokens: tk[13:14],
										},
										Tokens: tk[13:14],
									},
									Tokens: tk[13:15],
								},
							},
							Tokens: tk[13:15],
						},
					},
					Tokens: tk[13:15],
				},
				Tokens: tk[:16],
			}
		}},
		{"for (( a=1; a<2; a++ )); do b;done", func(t *test, tk Tokens) { // 6
			t.Output = ForCompound{
				ArithmeticExpansion: &ArithmeticExpansion{
					Expression: true,
					WordsAndOperators: []WordOrOperator{
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
							Operator: &tk[5],
							Tokens:   tk[5:6],
						},
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
						{
							Operator: &tk[7],
							Tokens:   tk[7:8],
						},
						{
							Word: &Word{
								Parts: []WordPart{
									{
										Part:   &tk[9],
										Tokens: tk[9:10],
									},
								},
								Tokens: tk[9:10],
							},
							Tokens: tk[9:10],
						},
						{
							Operator: &tk[10],
							Tokens:   tk[10:11],
						},
						{
							Word: &Word{
								Parts: []WordPart{
									{
										Part:   &tk[11],
										Tokens: tk[11:12],
									},
								},
								Tokens: tk[11:12],
							},
							Tokens: tk[11:12],
						},
						{
							Operator: &tk[12],
							Tokens:   tk[12:13],
						},
						{
							Word: &Word{
								Parts: []WordPart{
									{
										Part:   &tk[14],
										Tokens: tk[14:15],
									},
								},
								Tokens: tk[14:15],
							},
							Tokens: tk[14:15],
						},
						{
							Operator: &tk[15],
							Tokens:   tk[15:16],
						},
					},
					Tokens: tk[2:18],
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
																	Part:   &tk[22],
																	Tokens: tk[22:23],
																},
															},
															Tokens: tk[22:23],
														},
														Tokens: tk[22:23],
													},
												},
												Tokens: tk[22:23],
											},
											Tokens: tk[22:23],
										},
										Tokens: tk[22:23],
									},
									Tokens: tk[22:24],
								},
							},
							Tokens: tk[22:24],
						},
					},
					Tokens: tk[22:24],
				},
				Tokens: tk[:25],
			}
		}},
		{"for a in $(||); do b;done", func(t *test, tk Tokens) { // 7
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
													Token:   tk[7],
												},
												Parsing: "CommandOrCompound",
												Token:   tk[7],
											},
											Parsing: "Pipeline",
											Token:   tk[7],
										},
										Parsing: "Statement",
										Token:   tk[7],
									},
									Parsing: "Line",
									Token:   tk[7],
								},
								Parsing: "File",
								Token:   tk[7],
							},
							Parsing: "CommandSubstitution",
							Token:   tk[7],
						},
						Parsing: "WordPart",
						Token:   tk[6],
					},
					Parsing: "Word",
					Token:   tk[6],
				},
				Parsing: "ForCompound",
				Token:   tk[6],
			}
		}},
		{"for (( $(||) )); do b;done", func(t *test, tk Tokens) { // 8
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
										},
										Parsing: "File",
										Token:   tk[5],
									},
									Parsing: "CommandSubstitution",
									Token:   tk[5],
								},
								Parsing: "WordPart",
								Token:   tk[4],
							},
							Parsing: "Word",
							Token:   tk[4],
						},
						Parsing: "WordOrOperator",
						Token:   tk[4],
					},
					Parsing: "ArithmeticExpansion",
					Token:   tk[4],
				},
				Parsing: "ForCompound",
				Token:   tk[2],
			}
		}},
		{"for a; do ||;done", func(t *test, tk Tokens) { // 9
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err: Error{
									Err: Error{
										Err:     ErrMissingWord,
										Parsing: "Command",
										Token:   tk[7],
									},
									Parsing: "CommandOrCompound",
									Token:   tk[7],
								},
								Parsing: "Pipeline",
								Token:   tk[7],
							},
							Parsing: "Statement",
							Token:   tk[7],
						},
						Parsing: "Line",
						Token:   tk[7],
					},
					Parsing: "File",
					Token:   tk[7],
				},
				Parsing: "ForCompound",
				Token:   tk[7],
			}
		}},
	}, func(t *test) (Type, error) {
		var f ForCompound

		err := f.parse(t.Parser)

		return f, err
	})
}

func TestSelectCompound(t *testing.T) {
	doTests(t, []sourceFn{
		{"select a; do b;done", func(t *test, tk Tokens) { // 1
			t.Output = SelectCompound{
				Identifier: &tk[2],
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
				Tokens: tk[:10],
			}
		}},
		{"select a\ndo b\ndone", func(t *test, tk Tokens) { // 2
			t.Output = SelectCompound{
				Identifier: &tk[2],
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
							Tokens: tk[6:7],
						},
					},
					Tokens: tk[6:7],
				},
				Tokens: tk[:9],
			}
		}},
		{"select a in; do b;done", func(t *test, tk Tokens) { // 3
			t.Output = SelectCompound{
				Identifier: &tk[2],
				Words:      []Word{},
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
																	Part:   &tk[9],
																	Tokens: tk[9:10],
																},
															},
															Tokens: tk[9:10],
														},
														Tokens: tk[9:10],
													},
												},
												Tokens: tk[9:10],
											},
											Tokens: tk[9:10],
										},
										Tokens: tk[9:10],
									},
									Tokens: tk[9:11],
								},
							},
							Tokens: tk[9:11],
						},
					},
					Tokens: tk[9:11],
				},
				Tokens: tk[:12],
			}
		}},
		{"select a in b; do c;done", func(t *test, tk Tokens) { // 4
			t.Output = SelectCompound{
				Identifier: &tk[2],
				Words: []Word{
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
																	Part:   &tk[11],
																	Tokens: tk[11:12],
																},
															},
															Tokens: tk[11:12],
														},
														Tokens: tk[11:12],
													},
												},
												Tokens: tk[11:12],
											},
											Tokens: tk[11:12],
										},
										Tokens: tk[11:12],
									},
									Tokens: tk[11:13],
								},
							},
							Tokens: tk[11:13],
						},
					},
					Tokens: tk[11:13],
				},
				Tokens: tk[:14],
			}
		}},
		{"select a in b c; do d;done", func(t *test, tk Tokens) { // 5
			t.Output = SelectCompound{
				Identifier: &tk[2],
				Words: []Word{
					{
						Parts: []WordPart{
							{
								Part:   &tk[6],
								Tokens: tk[6:7],
							},
						},
						Tokens: tk[6:7],
					},
					{
						Parts: []WordPart{
							{
								Part:   &tk[8],
								Tokens: tk[8:9],
							},
						},
						Tokens: tk[8:9],
					},
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
																	Part:   &tk[13],
																	Tokens: tk[13:14],
																},
															},
															Tokens: tk[13:14],
														},
														Tokens: tk[13:14],
													},
												},
												Tokens: tk[13:14],
											},
											Tokens: tk[13:14],
										},
										Tokens: tk[13:14],
									},
									Tokens: tk[13:15],
								},
							},
							Tokens: tk[13:15],
						},
					},
					Tokens: tk[13:15],
				},
				Tokens: tk[:16],
			}
		}},
		{"select a in $(||); do b;done", func(t *test, tk Tokens) { // 6
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
													Token:   tk[7],
												},
												Parsing: "CommandOrCompound",
												Token:   tk[7],
											},
											Parsing: "Pipeline",
											Token:   tk[7],
										},
										Parsing: "Statement",
										Token:   tk[7],
									},
									Parsing: "Line",
									Token:   tk[7],
								},
								Parsing: "File",
								Token:   tk[7],
							},
							Parsing: "CommandSubstitution",
							Token:   tk[7],
						},
						Parsing: "WordPart",
						Token:   tk[6],
					},
					Parsing: "Word",
					Token:   tk[6],
				},
				Parsing: "SelectCompound",
				Token:   tk[6],
			}
		}},
		{"select a; do ||;done", func(t *test, tk Tokens) { // 7
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err: Error{
									Err: Error{
										Err:     ErrMissingWord,
										Parsing: "Command",
										Token:   tk[7],
									},
									Parsing: "CommandOrCompound",
									Token:   tk[7],
								},
								Parsing: "Pipeline",
								Token:   tk[7],
							},
							Parsing: "Statement",
							Token:   tk[7],
						},
						Parsing: "Line",
						Token:   tk[7],
					},
					Parsing: "File",
					Token:   tk[7],
				},
				Parsing: "SelectCompound",
				Token:   tk[7],
			}
		}},
	}, func(t *test) (Type, error) {
		var s SelectCompound

		err := s.parse(t.Parser)

		return s, err
	})
}

func TestTestCompound(t *testing.T) {
	doTests(t, []sourceFn{
		{"[[ a = b ]]", func(t *test, tk Tokens) { // 1
			t.Output = TestCompound{
				Tests: Tests{
					Test: TestOperatorStringsEqual,
					Word: &Word{
						Parts: []WordPart{
							{
								Part:   &tk[2],
								Tokens: tk[2:3],
							},
						},
						Tokens: tk[2:3],
					},
					Pattern: &Pattern{
						Parts: []WordPart{
							{
								Part:   &tk[6],
								Tokens: tk[6:7],
							},
						},
						Tokens: tk[6:7],
					},
					Tokens: tk[2:7],
				},
				Tokens: tk[:9],
			}
		}},
		{"[[ -a $(||) ]]", func(t *test, tk Tokens) { // 2
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
									},
									Parsing: "File",
									Token:   tk[5],
								},
								Parsing: "CommandSubstitution",
								Token:   tk[5],
							},
							Parsing: "WordPart",
							Token:   tk[4],
						},
						Parsing: "Word",
						Token:   tk[4],
					},
					Parsing: "Tests",
					Token:   tk[4],
				},
				Parsing: "TestCompound",
				Token:   tk[2],
			}
		}},
	}, func(t *test) (Type, error) {
		var tc TestCompound

		err := tc.parse(t.Parser)

		return tc, err
	})
}

func TestTests(t *testing.T) {
	doTests(t, []sourceFn{
		{"[[ a = b ]]", func(t *test, tk Tokens) { // 1
			t.Output = Tests{
				Test: TestOperatorStringsEqual,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[2],
							Tokens: tk[2:3],
						},
					},
					Tokens: tk[2:3],
				},
				Pattern: &Pattern{
					Parts: []WordPart{
						{
							Part:   &tk[6],
							Tokens: tk[6:7],
						},
					},
					Tokens: tk[6:7],
				},
				Tokens: tk[2:7],
			}
		}},
		{"[[ a == b ]]", func(t *test, tk Tokens) { // 2
			t.Output = Tests{
				Test: TestOperatorStringsEqual,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[2],
							Tokens: tk[2:3],
						},
					},
					Tokens: tk[2:3],
				},
				Pattern: &Pattern{
					Parts: []WordPart{
						{
							Part:   &tk[6],
							Tokens: tk[6:7],
						},
					},
					Tokens: tk[6:7],
				},
				Tokens: tk[2:7],
			}
		}},
		{"[[ a != b ]]", func(t *test, tk Tokens) { // 3
			t.Output = Tests{
				Test: TestOperatorStringsNotEqual,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[2],
							Tokens: tk[2:3],
						},
					},
					Tokens: tk[2:3],
				},
				Pattern: &Pattern{
					Parts: []WordPart{
						{
							Part:   &tk[6],
							Tokens: tk[6:7],
						},
					},
					Tokens: tk[6:7],
				},
				Tokens: tk[2:7],
			}
		}},
		{"[[ a =~ b ]]", func(t *test, tk Tokens) { // 4
			t.Output = Tests{
				Test: TestOperatorStringsMatch,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[2],
							Tokens: tk[2:3],
						},
					},
					Tokens: tk[2:3],
				},
				Pattern: &Pattern{
					Parts: []WordPart{
						{
							Part:   &tk[6],
							Tokens: tk[6:7],
						},
					},
					Tokens: tk[6:7],
				},
				Tokens: tk[2:7],
			}
		}},
		{"[[ a < b ]]", func(t *test, tk Tokens) { // 5
			t.Output = Tests{
				Test: TestOperatorStringBefore,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[2],
							Tokens: tk[2:3],
						},
					},
					Tokens: tk[2:3],
				},
				Pattern: &Pattern{
					Parts: []WordPart{
						{
							Part:   &tk[6],
							Tokens: tk[6:7],
						},
					},
					Tokens: tk[6:7],
				},
				Tokens: tk[2:7],
			}
		}},
		{"[[ a > b ]]", func(t *test, tk Tokens) { // 6
			t.Output = Tests{
				Test: TestOperatorStringAfter,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[2],
							Tokens: tk[2:3],
						},
					},
					Tokens: tk[2:3],
				},
				Pattern: &Pattern{
					Parts: []WordPart{
						{
							Part:   &tk[6],
							Tokens: tk[6:7],
						},
					},
					Tokens: tk[6:7],
				},
				Tokens: tk[2:7],
			}
		}},
		{"[[ a -ef b ]]", func(t *test, tk Tokens) { // 7
			t.Output = Tests{
				Test: TestOperatorFilesAreSameInode,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[2],
							Tokens: tk[2:3],
						},
					},
					Tokens: tk[2:3],
				},
				Pattern: &Pattern{
					Parts: []WordPart{
						{
							Part:   &tk[6],
							Tokens: tk[6:7],
						},
					},
					Tokens: tk[6:7],
				},
				Tokens: tk[2:7],
			}
		}},
		{"[[ a -nt b ]]", func(t *test, tk Tokens) { // 8
			t.Output = Tests{
				Test: TestOperatorFileIsNewerThan,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[2],
							Tokens: tk[2:3],
						},
					},
					Tokens: tk[2:3],
				},
				Pattern: &Pattern{
					Parts: []WordPart{
						{
							Part:   &tk[6],
							Tokens: tk[6:7],
						},
					},
					Tokens: tk[6:7],
				},
				Tokens: tk[2:7],
			}
		}},
		{"[[ a -ot b ]]", func(t *test, tk Tokens) { // 9
			t.Output = Tests{
				Test: TestOperatorFileIsOlderThan,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[2],
							Tokens: tk[2:3],
						},
					},
					Tokens: tk[2:3],
				},
				Pattern: &Pattern{
					Parts: []WordPart{
						{
							Part:   &tk[6],
							Tokens: tk[6:7],
						},
					},
					Tokens: tk[6:7],
				},
				Tokens: tk[2:7],
			}
		}},
		{"[[ a -eq b ]]", func(t *test, tk Tokens) { // 10
			t.Output = Tests{
				Test: TestOperatorEqual,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[2],
							Tokens: tk[2:3],
						},
					},
					Tokens: tk[2:3],
				},
				Pattern: &Pattern{
					Parts: []WordPart{
						{
							Part:   &tk[6],
							Tokens: tk[6:7],
						},
					},
					Tokens: tk[6:7],
				},
				Tokens: tk[2:7],
			}
		}},
		{"[[ a -ne b ]]", func(t *test, tk Tokens) { // 11
			t.Output = Tests{
				Test: TestOperatorNotEqual,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[2],
							Tokens: tk[2:3],
						},
					},
					Tokens: tk[2:3],
				},
				Pattern: &Pattern{
					Parts: []WordPart{
						{
							Part:   &tk[6],
							Tokens: tk[6:7],
						},
					},
					Tokens: tk[6:7],
				},
				Tokens: tk[2:7],
			}
		}},
		{"[[ a -le b ]]", func(t *test, tk Tokens) { // 12
			t.Output = Tests{
				Test: TestOperatorLessThanEqual,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[2],
							Tokens: tk[2:3],
						},
					},
					Tokens: tk[2:3],
				},
				Pattern: &Pattern{
					Parts: []WordPart{
						{
							Part:   &tk[6],
							Tokens: tk[6:7],
						},
					},
					Tokens: tk[6:7],
				},
				Tokens: tk[2:7],
			}
		}},
		{"[[ a -lt b ]]", func(t *test, tk Tokens) { // 13
			t.Output = Tests{
				Test: TestOperatorLessThan,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[2],
							Tokens: tk[2:3],
						},
					},
					Tokens: tk[2:3],
				},
				Pattern: &Pattern{
					Parts: []WordPart{
						{
							Part:   &tk[6],
							Tokens: tk[6:7],
						},
					},
					Tokens: tk[6:7],
				},
				Tokens: tk[2:7],
			}
		}},
		{"[[ a -ge b ]]", func(t *test, tk Tokens) { // 14
			t.Output = Tests{
				Test: TestOperatorGreaterThanEqual,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[2],
							Tokens: tk[2:3],
						},
					},
					Tokens: tk[2:3],
				},
				Pattern: &Pattern{
					Parts: []WordPart{
						{
							Part:   &tk[6],
							Tokens: tk[6:7],
						},
					},
					Tokens: tk[6:7],
				},
				Tokens: tk[2:7],
			}
		}},
		{"[[ a -gt b ]]", func(t *test, tk Tokens) { // 15
			t.Output = Tests{
				Test: TestOperatorGreaterThan,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[2],
							Tokens: tk[2:3],
						},
					},
					Tokens: tk[2:3],
				},
				Pattern: &Pattern{
					Parts: []WordPart{
						{
							Part:   &tk[6],
							Tokens: tk[6:7],
						},
					},
					Tokens: tk[6:7],
				},
				Tokens: tk[2:7],
			}
		}},
		{"[[ -a a ]]", func(t *test, tk Tokens) { // 16
			t.Output = Tests{
				Test: TestOperatorFileExists,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[4],
							Tokens: tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				Tokens: tk[2:5],
			}
		}},
		{"[[ -e a ]]", func(t *test, tk Tokens) { // 17
			t.Output = Tests{
				Test: TestOperatorFileExists,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[4],
							Tokens: tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				Tokens: tk[2:5],
			}
		}},
		{"[[ -b a ]]", func(t *test, tk Tokens) { // 18
			t.Output = Tests{
				Test: TestOperatorFileIsBlock,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[4],
							Tokens: tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				Tokens: tk[2:5],
			}
		}},
		{"[[ -c a ]]", func(t *test, tk Tokens) { // 19
			t.Output = Tests{
				Test: TestOperatorFileIsCharacter,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[4],
							Tokens: tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				Tokens: tk[2:5],
			}
		}},
		{"[[ -d a ]]", func(t *test, tk Tokens) { // 20
			t.Output = Tests{
				Test: TestOperatorDirectoryExists,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[4],
							Tokens: tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				Tokens: tk[2:5],
			}
		}},
		{"[[ -f a ]]", func(t *test, tk Tokens) { // 21
			t.Output = Tests{
				Test: TestOperatorFileIsRegular,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[4],
							Tokens: tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				Tokens: tk[2:5],
			}
		}},
		{"[[ -g a ]]", func(t *test, tk Tokens) { // 22
			t.Output = Tests{
				Test: TestOperatorFileHasSetGroupID,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[4],
							Tokens: tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				Tokens: tk[2:5],
			}
		}},
		{"[[ -h a ]]", func(t *test, tk Tokens) { // 23
			t.Output = Tests{
				Test: TestOperatorFileIsSymbolic,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[4],
							Tokens: tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				Tokens: tk[2:5],
			}
		}},
		{"[[ -L a ]]", func(t *test, tk Tokens) { // 24
			t.Output = Tests{
				Test: TestOperatorFileIsSymbolic,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[4],
							Tokens: tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				Tokens: tk[2:5],
			}
		}},
		{"[[ -k a ]]", func(t *test, tk Tokens) { // 25
			t.Output = Tests{
				Test: TestOperatorFileHasStickyBit,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[4],
							Tokens: tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				Tokens: tk[2:5],
			}
		}},
		{"[[ -p a ]]", func(t *test, tk Tokens) { // 26
			t.Output = Tests{
				Test: TestOperatorFileIsPipe,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[4],
							Tokens: tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				Tokens: tk[2:5],
			}
		}},
		{"[[ -r a ]]", func(t *test, tk Tokens) { // 27
			t.Output = Tests{
				Test: TestOperatorFileIsReadable,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[4],
							Tokens: tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				Tokens: tk[2:5],
			}
		}},
		{"[[ -s a ]]", func(t *test, tk Tokens) { // 28
			t.Output = Tests{
				Test: TestOperatorFileIsNonZero,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[4],
							Tokens: tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				Tokens: tk[2:5],
			}
		}},
		{"[[ -t a ]]", func(t *test, tk Tokens) { // 29
			t.Output = Tests{
				Test: TestOperatorFileIsTerminal,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[4],
							Tokens: tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				Tokens: tk[2:5],
			}
		}},
		{"[[ -u a ]]", func(t *test, tk Tokens) { // 30
			t.Output = Tests{
				Test: TestOperatorFileHasSetUserID,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[4],
							Tokens: tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				Tokens: tk[2:5],
			}
		}},
		{"[[ -w a ]]", func(t *test, tk Tokens) { // 31
			t.Output = Tests{
				Test: TestOperatorFileIsWritable,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[4],
							Tokens: tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				Tokens: tk[2:5],
			}
		}},
		{"[[ -x a ]]", func(t *test, tk Tokens) { // 32
			t.Output = Tests{
				Test: TestOperatorFileIsExecutable,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[4],
							Tokens: tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				Tokens: tk[2:5],
			}
		}},
		{"[[ -G a ]]", func(t *test, tk Tokens) { // 33
			t.Output = Tests{
				Test: TestOperatorFileIsOwnedByEffectiveGroup,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[4],
							Tokens: tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				Tokens: tk[2:5],
			}
		}},
		{"[[ -N a ]]", func(t *test, tk Tokens) { // 34
			t.Output = Tests{
				Test: TestOperatorFileWasModifiedSinceLastRead,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[4],
							Tokens: tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				Tokens: tk[2:5],
			}
		}},
		{"[[ -O a ]]", func(t *test, tk Tokens) { // 35
			t.Output = Tests{
				Test: TestOperatorFileIsOwnedByEffectiveUser,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[4],
							Tokens: tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				Tokens: tk[2:5],
			}
		}},
		{"[[ -S a ]]", func(t *test, tk Tokens) { // 36
			t.Output = Tests{
				Test: TestOperatorFileIsSocket,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[4],
							Tokens: tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				Tokens: tk[2:5],
			}
		}},
		{"[[ -o a ]]", func(t *test, tk Tokens) { // 37
			t.Output = Tests{
				Test: TestOperatorOptNameIsEnabled,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[4],
							Tokens: tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				Tokens: tk[2:5],
			}
		}},
		{"[[ -v a ]]", func(t *test, tk Tokens) { // 38
			t.Output = Tests{
				Test: TestOperatorVarNameIsSet,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[4],
							Tokens: tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				Tokens: tk[2:5],
			}
		}},
		{"[[ -R a ]]", func(t *test, tk Tokens) { // 39
			t.Output = Tests{
				Test: TestOperatorVarnameIsRef,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[4],
							Tokens: tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				Tokens: tk[2:5],
			}
		}},
		{"[[ -z a ]]", func(t *test, tk Tokens) { // 40
			t.Output = Tests{
				Test: TestOperatorStringIsZero,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[4],
							Tokens: tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				Tokens: tk[2:5],
			}
		}},
		{"[[ -n a ]]", func(t *test, tk Tokens) { // 41
			t.Output = Tests{
				Test: TestOperatorStringIsNonZero,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[4],
							Tokens: tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				Tokens: tk[2:5],
			}
		}},
		{"[[ ! -z a ]]", func(t *test, tk Tokens) { // 42
			t.Output = Tests{
				Not:  true,
				Test: TestOperatorStringIsZero,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[6],
							Tokens: tk[6:7],
						},
					},
					Tokens: tk[6:7],
				},
				Tokens: tk[2:7],
			}
		}},
		{"[[ (a = b) ]]", func(t *test, tk Tokens) { // 43
			t.Output = Tests{
				Parens: &Tests{
					Test: TestOperatorStringsEqual,
					Word: &Word{
						Parts: []WordPart{
							{
								Part:   &tk[3],
								Tokens: tk[3:4],
							},
						},
						Tokens: tk[3:4],
					},
					Pattern: &Pattern{
						Parts: []WordPart{
							{
								Part:   &tk[7],
								Tokens: tk[7:8],
							},
						},
						Tokens: tk[7:8],
					},
					Tokens: tk[3:8],
				},
				Tokens: tk[2:9],
			}
		}},
		{"[[ a ]]", func(t *test, tk Tokens) { // 44
			t.Output = Tests{
				Test: TestOperatorNone,
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
		{"[[ a || b ]]", func(t *test, tk Tokens) { // 45
			t.Output = Tests{
				Test: TestOperatorNone,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[2],
							Tokens: tk[2:3],
						},
					},
					Tokens: tk[2:3],
				},
				LogicalOperator: LogicalOperatorOr,
				Tests: &Tests{
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
				Tokens: tk[2:7],
			}
		}},
		{"[[ a && b ]]", func(t *test, tk Tokens) { // 46
			t.Output = Tests{
				Test: TestOperatorNone,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[2],
							Tokens: tk[2:3],
						},
					},
					Tokens: tk[2:3],
				},
				LogicalOperator: LogicalOperatorAnd,
				Tests: &Tests{
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
				Tokens: tk[2:7],
			}
		}},
		{"[[ -a $(||) ]]", func(t *test, tk Tokens) { // 47
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
								},
								Parsing: "File",
								Token:   tk[5],
							},
							Parsing: "CommandSubstitution",
							Token:   tk[5],
						},
						Parsing: "WordPart",
						Token:   tk[4],
					},
					Parsing: "Word",
					Token:   tk[4],
				},
				Parsing: "Tests",
				Token:   tk[4],
			}
		}},
		{"[[ $(||) = a ]]", func(t *test, tk Tokens) { // 48
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
				Parsing: "Tests",
				Token:   tk[2],
			}
		}},
		{"[[ a = $(||) ]]", func(t *test, tk Tokens) { // 49
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
													Token:   tk[7],
												},
												Parsing: "CommandOrCompound",
												Token:   tk[7],
											},
											Parsing: "Pipeline",
											Token:   tk[7],
										},
										Parsing: "Statement",
										Token:   tk[7],
									},
									Parsing: "Line",
									Token:   tk[7],
								},
								Parsing: "File",
								Token:   tk[7],
							},
							Parsing: "CommandSubstitution",
							Token:   tk[7],
						},
						Parsing: "WordPart",
						Token:   tk[6],
					},
					Parsing: "Pattern",
					Token:   tk[6],
				},
				Parsing: "Tests",
				Token:   tk[6],
			}
		}},
		{"[[ a -eq $(||) ]]", func(t *test, tk Tokens) { // 50
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
													Token:   tk[7],
												},
												Parsing: "CommandOrCompound",
												Token:   tk[7],
											},
											Parsing: "Pipeline",
											Token:   tk[7],
										},
										Parsing: "Statement",
										Token:   tk[7],
									},
									Parsing: "Line",
									Token:   tk[7],
								},
								Parsing: "File",
								Token:   tk[7],
							},
							Parsing: "CommandSubstitution",
							Token:   tk[7],
						},
						Parsing: "WordPart",
						Token:   tk[6],
					},
					Parsing: "Pattern",
					Token:   tk[6],
				},
				Parsing: "Tests",
				Token:   tk[6],
			}
		}},
		{"[[ ( $(||) ) ]]", func(t *test, tk Tokens) { // 51
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
									},
									Parsing: "File",
									Token:   tk[5],
								},
								Parsing: "CommandSubstitution",
								Token:   tk[5],
							},
							Parsing: "WordPart",
							Token:   tk[4],
						},
						Parsing: "Word",
						Token:   tk[4],
					},
					Parsing: "Tests",
					Token:   tk[4],
				},
				Parsing: "Tests",
				Token:   tk[4],
			}
		}},
		{"[[ ( a b ) ]]", func(t *test, tk Tokens) { // 52
			t.Err = Error{
				Err:     ErrMissingClosingParen,
				Parsing: "Tests",
				Token:   tk[6],
			}
		}},
		{"[[ a || $(||) ]]", func(t *test, tk Tokens) { // 53
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
														Token:   tk[7],
													},
													Parsing: "CommandOrCompound",
													Token:   tk[7],
												},
												Parsing: "Pipeline",
												Token:   tk[7],
											},
											Parsing: "Statement",
											Token:   tk[7],
										},
										Parsing: "Line",
										Token:   tk[7],
									},
									Parsing: "File",
									Token:   tk[7],
								},
								Parsing: "CommandSubstitution",
								Token:   tk[7],
							},
							Parsing: "WordPart",
							Token:   tk[6],
						},
						Parsing: "Word",
						Token:   tk[6],
					},
					Parsing: "Tests",
					Token:   tk[6],
				},
				Parsing: "Tests",
				Token:   tk[6],
			}
		}},
	}, func(t *test) (Type, error) {
		var ts Tests

		t.Parser.Tokens = t.Parser.Tokens[2:2]
		err := ts.parse(t.Parser)

		return ts, err
	})
}

func TestPattern(t *testing.T) {
	doTests(t, []sourceFn{
		{"[[ z = a ]]", func(t *test, tk Tokens) { // 1
			t.Output = Pattern{
				Parts: []WordPart{
					{
						Part:   &tk[6],
						Tokens: tk[6:7],
					},
				},
				Tokens: tk[6:7],
			}
		}},
		{"[[ z = a$b ]]", func(t *test, tk Tokens) { // 2
			t.Output = Pattern{
				Parts: []WordPart{
					{
						Part:   &tk[6],
						Tokens: tk[6:7],
					},
					{
						Part:   &tk[7],
						Tokens: tk[7:8],
					},
				},
				Tokens: tk[6:8],
			}
		}},
		{"[[ z = $(||) ]]", func(t *test, tk Tokens) { // 3
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
												Token:   tk[7],
											},
											Parsing: "CommandOrCompound",
											Token:   tk[7],
										},
										Parsing: "Pipeline",
										Token:   tk[7],
									},
									Parsing: "Statement",
									Token:   tk[7],
								},
								Parsing: "Line",
								Token:   tk[7],
							},
							Parsing: "File",
							Token:   tk[7],
						},
						Parsing: "CommandSubstitution",
						Token:   tk[7],
					},
					Parsing: "WordPart",
					Token:   tk[6],
				},
				Parsing: "Pattern",
				Token:   tk[6],
			}
		}},
	}, func(t *test) (Type, error) {
		var p Pattern

		t.Parser.Tokens = t.Parser.Tokens[6:6]
		err := p.parse(t.Parser)

		return p, err
	})
}

func TestGroupingCompound(t *testing.T) {
	doTests(t, []sourceFn{
		{"(a)", func(t *test, tk Tokens) { // 1
			t.Output = GroupingCompound{
				SubShell: true,
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
																	Part:   &tk[1],
																	Tokens: tk[1:2],
																},
															},
															Tokens: tk[1:2],
														},
														Tokens: tk[1:2],
													},
												},
												Tokens: tk[1:2],
											},
											Tokens: tk[1:2],
										},
										Tokens: tk[1:2],
									},
									Tokens: tk[1:2],
								},
							},
							Tokens: tk[1:2],
						},
					},
					Tokens: tk[1:2],
				},
				Tokens: tk[:3],
			}
		}},
		{"( a )", func(t *test, tk Tokens) { // 2
			t.Output = GroupingCompound{
				SubShell: true,
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
					Tokens: tk[2:3],
				},
				Tokens: tk[:5],
			}
		}},
		{"(\na\n)", func(t *test, tk Tokens) { // 3
			t.Output = GroupingCompound{
				SubShell: true,
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
					Tokens: tk[2:3],
				},
				Tokens: tk[:5],
			}
		}},
		{"(\na;\n)", func(t *test, tk Tokens) { // 4
			t.Output = GroupingCompound{
				SubShell: true,
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
							},
							Tokens: tk[2:4],
						},
					},
					Tokens: tk[2:4],
				},
				Tokens: tk[:6],
			}
		}},
		{"{\na\n}", func(t *test, tk Tokens) { // 5
			t.Output = GroupingCompound{
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
					Tokens: tk[2:3],
				},
				Tokens: tk[:5],
			}
		}},
		{"{ a; }", func(t *test, tk Tokens) { // 6
			t.Output = GroupingCompound{
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
							},
							Tokens: tk[2:4],
						},
					},
					Tokens: tk[2:4],
				},
				Tokens: tk[:6],
			}
		}},
		{"(||)", func(t *test, tk Tokens) { // 7
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
									Parsing: "CommandOrCompound",
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
				Parsing: "GroupingCompound",
				Token:   tk[1],
			}
		}},
	}, func(t *test) (Type, error) {
		var g GroupingCompound

		err := g.parse(t.Parser)

		return g, err
	})
}

func TestFunctionCompound(t *testing.T) {
	doTests(t, []sourceFn{
		{"function a() { b; }", func(t *test, tk Tokens) { // 1
			t.Output = FunctionCompound{
				HasKeyword: true,
				Identifier: &tk[2],
				Body: Compound{
					GroupingCompound: &GroupingCompound{
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
																			Part:   &tk[8],
																			Tokens: tk[8:9],
																		},
																	},
																	Tokens: tk[8:9],
																},
																Tokens: tk[8:9],
															},
														},
														Tokens: tk[8:9],
													},
													Tokens: tk[8:9],
												},
												Tokens: tk[8:9],
											},
											Tokens: tk[8:10],
										},
									},
									Tokens: tk[8:10],
								},
							},
							Tokens: tk[8:10],
						},
						Tokens: tk[6:12],
					},
					Tokens: tk[6:12],
				},
				Tokens: tk[:12],
			}
		}},
		{"function a ( ) { b; }", func(t *test, tk Tokens) { // 2
			t.Output = FunctionCompound{
				HasKeyword: true,
				Identifier: &tk[2],
				Body: Compound{
					GroupingCompound: &GroupingCompound{
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
																			Part:   &tk[10],
																			Tokens: tk[10:11],
																		},
																	},
																	Tokens: tk[10:11],
																},
																Tokens: tk[10:11],
															},
														},
														Tokens: tk[10:11],
													},
													Tokens: tk[10:11],
												},
												Tokens: tk[10:11],
											},
											Tokens: tk[10:12],
										},
									},
									Tokens: tk[10:12],
								},
							},
							Tokens: tk[10:12],
						},
						Tokens: tk[8:14],
					},
					Tokens: tk[8:14],
				},
				Tokens: tk[:14],
			}
		}},
		{"a() { b; }", func(t *test, tk Tokens) { // 3
			t.Output = FunctionCompound{
				Identifier: &tk[0],
				Body: Compound{
					GroupingCompound: &GroupingCompound{
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
											Tokens: tk[6:8],
										},
									},
									Tokens: tk[6:8],
								},
							},
							Tokens: tk[6:8],
						},
						Tokens: tk[4:10],
					},
					Tokens: tk[4:10],
				},
				Tokens: tk[:10],
			}
		}},
		{"a ( ) { b; }", func(t *test, tk Tokens) { // 4
			t.Output = FunctionCompound{
				Identifier: &tk[0],
				Body: Compound{
					GroupingCompound: &GroupingCompound{
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
																			Part:   &tk[8],
																			Tokens: tk[8:9],
																		},
																	},
																	Tokens: tk[8:9],
																},
																Tokens: tk[8:9],
															},
														},
														Tokens: tk[8:9],
													},
													Tokens: tk[8:9],
												},
												Tokens: tk[8:9],
											},
											Tokens: tk[8:10],
										},
									},
									Tokens: tk[8:10],
								},
							},
							Tokens: tk[8:10],
						},
						Tokens: tk[6:12],
					},
					Tokens: tk[6:12],
				},
				Tokens: tk[:12],
			}
		}},
		{"function a() { || }", func(t *test, tk Tokens) { // 5
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
												Token:   tk[8],
											},
											Parsing: "CommandOrCompound",
											Token:   tk[8],
										},
										Parsing: "Pipeline",
										Token:   tk[8],
									},
									Parsing: "Statement",
									Token:   tk[8],
								},
								Parsing: "Line",
								Token:   tk[8],
							},
							Parsing: "File",
							Token:   tk[8],
						},
						Parsing: "GroupingCompound",
						Token:   tk[8],
					},
					Parsing: "Compound",
					Token:   tk[6],
				},
				Parsing: "FunctionCompound",
				Token:   tk[6],
			}
		}},
	}, func(t *test) (Type, error) {
		var f FunctionCompound

		err := f.parse(t.Parser)

		return f, err
	})
}

func TestAssignmentOrWord(t *testing.T) {
	doTests(t, []sourceFn{
		{"a", func(t *test, tk Tokens) { // 1
			t.Output = AssignmentOrWord{
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
		{"a=b", func(t *test, tk Tokens) { // 2
			t.Output = AssignmentOrWord{
				Assignment: &Assignment{
					Identifier: ParameterAssign{
						Identifier: &tk[0],
						Tokens:     tk[:1],
					},
					Assignment: AssignmentAssign,
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
				Tokens: tk[:3],
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
												Parsing: "CommandOrCompound",
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
				Parsing: "AssignmentOrWord",
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
														Err: Error{
															Err:     ErrMissingWord,
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
				},
				Parsing: "AssignmentOrWord",
				Token:   tk[0],
			}
		}},
	}, func(t *test) (Type, error) {
		var a AssignmentOrWord

		err := a.parse(t.Parser)

		return a, err
	})
}

func TestCommand(t *testing.T) {
	doTests(t, []sourceFn{
		{"a", func(t *test, tk Tokens) { // 1
			t.Output = Command{
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
				AssignmentsOrWords: []AssignmentOrWord{
					{
						Word: &Word{
							Parts: []WordPart{
								{
									Part:   &tk[11],
									Tokens: tk[11:12],
								},
							},
							Tokens: tk[11:12],
						},
						Tokens: tk[11:12],
					},
					{
						Assignment: &Assignment{
							Identifier: ParameterAssign{
								Identifier: &tk[13],
								Tokens:     tk[13:14],
							},
							Assignment: AssignmentAssign,
							Value: Value{
								Word: &Word{
									Parts: []WordPart{
										{
											Part:   &tk[15],
											Tokens: tk[15:16],
										},
									},
									Tokens: tk[15:16],
								},
								Tokens: tk[15:16],
							},
							Tokens: tk[13:16],
						},
						Tokens: tk[13:16],
					},
				},
				Tokens: tk[:20],
			}
		}},
		{"declare a", func(t *test, tk Tokens) { // 6
			t.Output = Command{
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
				Tokens: tk[:3],
			}
		}},
		{"local -n a=b", func(t *test, tk Tokens) { // 7
			t.Output = Command{
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
						Assignment: &Assignment{
							Identifier: ParameterAssign{
								Identifier: &tk[4],
								Tokens:     tk[4:5],
							},
							Assignment: AssignmentAssign,
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
						Tokens: tk[4:7],
					},
				},
				Tokens: tk[:7],
			}
		}},
		{"readonly -a -p a=b c=d", func(t *test, tk Tokens) { // 8
			t.Output = Command{
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
						Assignment: &Assignment{
							Identifier: ParameterAssign{
								Identifier: &tk[6],
								Tokens:     tk[6:7],
							},
							Assignment: AssignmentAssign,
							Value: Value{
								Word: &Word{
									Parts: []WordPart{
										{
											Part:   &tk[8],
											Tokens: tk[8:9],
										},
									},
									Tokens: tk[8:9],
								},
								Tokens: tk[8:9],
							},
							Tokens: tk[6:9],
						},
						Tokens: tk[6:9],
					},
					{
						Assignment: &Assignment{
							Identifier: ParameterAssign{
								Identifier: &tk[10],
								Tokens:     tk[10:11],
							},
							Assignment: AssignmentAssign,
							Value: Value{
								Word: &Word{
									Parts: []WordPart{
										{
											Part:   &tk[12],
											Tokens: tk[12:13],
										},
									},
									Tokens: tk[12:13],
								},
								Tokens: tk[12:13],
							},
							Tokens: tk[10:13],
						},
						Tokens: tk[10:13],
					},
				},
				Tokens: tk[:13],
			}
		}},
		{"export -p a=b c=d", func(t *test, tk Tokens) { // 9
			t.Output = Command{
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
						Assignment: &Assignment{
							Identifier: ParameterAssign{
								Identifier: &tk[4],
								Tokens:     tk[4:5],
							},
							Assignment: AssignmentAssign,
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
						Tokens: tk[4:7],
					},
					{
						Assignment: &Assignment{
							Identifier: ParameterAssign{
								Identifier: &tk[8],
								Tokens:     tk[8:9],
							},
							Assignment: AssignmentAssign,
							Value: Value{
								Word: &Word{
									Parts: []WordPart{
										{
											Part:   &tk[10],
											Tokens: tk[10:11],
										},
									},
									Tokens: tk[10:11],
								},
								Tokens: tk[10:11],
							},
							Tokens: tk[8:11],
						},
						Tokens: tk[8:11],
					},
				},
				Tokens: tk[:11],
			}
		}},
		{"typeset -ag -pl a=b", func(t *test, tk Tokens) { // 10
			t.Output = Command{
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
						Assignment: &Assignment{
							Identifier: ParameterAssign{
								Identifier: &tk[6],
								Tokens:     tk[6:7],
							},
							Assignment: AssignmentAssign,
							Value: Value{
								Word: &Word{
									Parts: []WordPart{
										{
											Part:   &tk[8],
											Tokens: tk[8:9],
										},
									},
									Tokens: tk[8:9],
								},
								Tokens: tk[8:9],
							},
							Tokens: tk[6:9],
						},
						Tokens: tk[6:9],
					},
				},
				Tokens: tk[:9],
			}
		}},
		{"a[$(||)]=", func(t *test, tk Tokens) { // 11
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
		{">$(||)", func(t *test, tk Tokens) { // 12
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
													Parsing: "CommandOrCompound",
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
		{"$(||)", func(t *test, tk Tokens) { // 13
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
													Parsing: "CommandOrCompound",
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
					Parsing: "AssignmentOrWord",
					Token:   tk[0],
				},
				Parsing: "Command",
				Token:   tk[0],
			}
		}},
		{"a >$(||)", func(t *test, tk Tokens) { // 14
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
					Parsing: "Redirection",
					Token:   tk[3],
				},
				Parsing: "Command",
				Token:   tk[2],
			}
		}},
		{"export a=$(||)", func(t *test, tk Tokens) { // 15
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
																Err:     ErrMissingWord,
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
											},
											Parsing: "File",
											Token:   tk[5],
										},
										Parsing: "CommandSubstitution",
										Token:   tk[5],
									},
									Parsing: "WordPart",
									Token:   tk[4],
								},
								Parsing: "Word",
								Token:   tk[4],
							},
							Parsing: "Value",
							Token:   tk[4],
						},
						Parsing: "Assignment",
						Token:   tk[4],
					},
					Parsing: "AssignmentOrWord",
					Token:   tk[2],
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
											Parsing: "CommandOrCompound",
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

		err := w.parse(t.Parser, false)

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
														Parsing: "CommandOrCompound",
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
										Parsing: "CommandOrCompound",
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
		{"${a/b/$(||)}", func(t *test, tk Tokens) { // 51
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
						Parsing: "WordOrToken",
						Token:   tk[5],
					},
					Parsing: "String",
					Token:   tk[5],
				},
				Parsing: "ParameterExpansion",
				Token:   tk[5],
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
		{"${a[ 0 ]}", func(t *test, tk Tokens) { // 7
			t.Output = Parameter{
				Parameter: &tk[1],
				Array: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[4],
							Tokens: tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				Tokens: tk[1:7],
			}
		}},
		{"${a[@]}", func(t *test, tk Tokens) { // 8
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
		{"${a[*]}", func(t *test, tk Tokens) { // 9
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
		{"${a[$(||)]}", func(t *test, tk Tokens) { // 10
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
				Parsing: "Parameter",
				Token:   tk[3],
			}
		}},
		{"${a[b c]}", func(t *test, tk Tokens) { // 11
			t.Err = Error{
				Err:     ErrMissingClosingBracket,
				Parsing: "Parameter",
				Token:   tk[5],
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
													Parsing: "CommandOrCompound",
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
												Parsing: "CommandOrCompound",
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
										CommandOrCompound: CommandOrCompound{
											Command: &Command{
												AssignmentsOrWords: []AssignmentOrWord{
													{
														Word: &Word{
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
										CommandOrCompound: CommandOrCompound{
											Command: &Command{
												AssignmentsOrWords: []AssignmentOrWord{
													{
														Word: &Word{
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
									Parsing: "CommandOrCompound",
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
												Parsing: "CommandOrCompound",
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
		{"(( a ))", func(t *test, tk Tokens) { // 2
			t.Output = ArithmeticExpansion{
				Expression: true,
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
		{"$(( a ))", func(t *test, tk Tokens) { // 3
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
		{"$(( a$b ))", func(t *test, tk Tokens) { // 4
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
		{"$((a+b))", func(t *test, tk Tokens) { // 5
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
		{"$(($(||)))", func(t *test, tk Tokens) { // 6
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
													Parsing: "CommandOrCompound",
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
												Parsing: "CommandOrCompound",
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
