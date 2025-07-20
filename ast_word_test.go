package bash

import "testing"

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
					Value: &Value{
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
				Array:  []ArrayWord{},
				Tokens: tk[2:4],
			}
		}},
		{"a=(b)", func(t *test, tk Tokens) { // 3
			t.Output = Value{
				Array: []ArrayWord{
					{
						Word: Word{
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
				Tokens: tk[2:5],
			}
		}},
		{"a=( b )", func(t *test, tk Tokens) { // 4
			t.Output = Value{
				Array: []ArrayWord{
					{
						Word: Word{
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
				Tokens: tk[2:7],
			}
		}},
		{"a=( b c )", func(t *test, tk Tokens) { // 5
			t.Output = Value{
				Array: []ArrayWord{
					{
						Word: Word{
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
						Word: Word{
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
				Tokens: tk[2:9],
			}
		}},
		{"a=( # comment\n)", func(t *test, tk Tokens) { // 6
			t.Output = Value{
				Array:    []ArrayWord{},
				Comments: [2]Comments{{tk[4]}},
				Tokens:   tk[2:7],
			}
		}},
		{"a=( # comment A\n\n# b comment B\n)", func(t *test, tk Tokens) { // 7
			t.Output = Value{
				Array:    []ArrayWord{},
				Comments: [2]Comments{{tk[4]}, {tk[6]}},
				Tokens:   tk[2:9],
			}
		}},
		{"a=(\n# comment\n)", func(t *test, tk Tokens) { // 8
			t.Output = Value{
				Array:    []ArrayWord{},
				Comments: [2]Comments{nil, {tk[4]}},
				Tokens:   tk[2:7],
			}
		}},
		{"a=( # comment A\nb\n# comment B\n)", func(t *test, tk Tokens) { // 9
			t.Output = Value{
				Array: []ArrayWord{
					{
						Word: Word{
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
				Comments: [2]Comments{{tk[4]}, {tk[8]}},
				Tokens:   tk[2:11],
			}
		}},
		{"a=$(||)", func(t *test, tk Tokens) { // 10
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
		{"a=($(||))", func(t *test, tk Tokens) { // 11
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
					Parsing: "ArrayWord",
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

func TestArrayWord(t *testing.T) {
	doTests(t, []sourceFn{
		{"a", func(t *test, tk Tokens) { // 1
			t.Output = ArrayWord{
				Word: Word{
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
		{"# comment\na", func(t *test, tk Tokens) { // 2
			t.Output = ArrayWord{
				Word: Word{
					Parts: []WordPart{
						{
							Part:   &tk[2],
							Tokens: tk[2:3],
						},
					},
					Tokens: tk[2:3],
				},
				Comments: [2]Comments{{tk[0]}},
				Tokens:   tk[:3],
			}
		}},
		{"a # comment", func(t *test, tk Tokens) { // 3
			t.Output = ArrayWord{
				Word: Word{
					Parts: []WordPart{
						{
							Part:   &tk[0],
							Tokens: tk[:1],
						},
					},
					Tokens: tk[:1],
				},
				Comments: [2]Comments{nil, {tk[2]}},
				Tokens:   tk[:3],
			}
		}},
		{"# comment A\na # comment B", func(t *test, tk Tokens) { // 4
			t.Output = ArrayWord{
				Word: Word{
					Parts: []WordPart{
						{
							Part:   &tk[2],
							Tokens: tk[2:3],
						},
					},
					Tokens: tk[2:3],
				},
				Comments: [2]Comments{{tk[0]}, {tk[4]}},
				Tokens:   tk[:5],
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
				Parsing: "ArrayWord",
				Token:   tk[0],
			}
		}},
		{"", func(t *test, tk Tokens) { // 6
			t.Err = Error{
				Err:     ErrMissingWord,
				Parsing: "ArrayWord",
				Token:   tk[0],
			}
		}},
	}, func(t *test) (Type, error) {
		var a ArrayWord

		err := a.parse(t.Parser)

		return a, err
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
					Backtick:         &tk[0],
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
		{"{a,,\"b\"}", func(t *test, tk Tokens) { // 6
			t.Output = WordPart{
				BraceExpansion: &BraceExpansion{
					BraceExpansionType: BraceExpansionWords,
					Words: []Word{
						{
							Parts: []WordPart{
								{
									Part:   &tk[1],
									Tokens: tk[1:2],
								},
							},
							Tokens: tk[1:2],
						},
						{
							Tokens: tk[3:3],
						},
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
					Tokens: tk[:6],
				},
				Tokens: tk[:6],
			}
		}},
		{"{a..e}", func(t *test, tk Tokens) { // 7
			t.Output = WordPart{
				BraceExpansion: &BraceExpansion{
					BraceExpansionType: BraceExpansionSequence,
					Words: []Word{
						{
							Parts: []WordPart{
								{
									Part:   &tk[1],
									Tokens: tk[1:2],
								},
							},
							Tokens: tk[1:2],
						},
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
					Tokens: tk[:5],
				},
				Tokens: tk[:5],
			}
		}},
		{"${a[$(||)]}", func(t *test, tk Tokens) { // 8
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
							Parsing: "WordOrOperator",
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
		{"$(($(||)))", func(t *test, tk Tokens) { // 9
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
		{"$(||)", func(t *test, tk Tokens) { // 10
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
		{"{\"$(||)\",}", func(t *test, tk Tokens) { // 11
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
					Parsing: "BraceExpansion",
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

func TestBraceExpansion(t *testing.T) {
	doTests(t, []sourceFn{
		{"{a,b}", func(t *test, tk Tokens) { // 1
			t.Output = BraceExpansion{
				BraceExpansionType: BraceExpansionWords,
				Words: []Word{
					{
						Parts: []WordPart{
							{
								Part:   &tk[1],
								Tokens: tk[1:2],
							},
						},
						Tokens: tk[1:2],
					},
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
				Tokens: tk[:5],
			}
		}},
		{"{\"a\",bc,123}", func(t *test, tk Tokens) { // 2
			t.Output = BraceExpansion{
				BraceExpansionType: BraceExpansionWords,
				Words: []Word{
					{
						Parts: []WordPart{
							{
								Part:   &tk[1],
								Tokens: tk[1:2],
							},
						},
						Tokens: tk[1:2],
					},
					{
						Parts: []WordPart{
							{
								Part:   &tk[3],
								Tokens: tk[3:4],
							},
						},
						Tokens: tk[3:4],
					},
					{
						Parts: []WordPart{
							{
								Part:   &tk[5],
								Tokens: tk[5:6],
							},
						},
						Tokens: tk[5:6],
					},
				},
				Tokens: tk[:7],
			}
		}},
		{"{a..e}", func(t *test, tk Tokens) { // 3
			t.Output = BraceExpansion{
				BraceExpansionType: BraceExpansionSequence,
				Words: []Word{
					{
						Parts: []WordPart{
							{
								Part:   &tk[1],
								Tokens: tk[1:2],
							},
						},
						Tokens: tk[1:2],
					},
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
				Tokens: tk[:5],
			}
		}},
		{"{100..10..-1}", func(t *test, tk Tokens) { // 4
			t.Output = BraceExpansion{
				BraceExpansionType: BraceExpansionSequence,
				Words: []Word{
					{
						Parts: []WordPart{
							{
								Part:   &tk[1],
								Tokens: tk[1:2],
							},
						},
						Tokens: tk[1:2],
					},
					{
						Parts: []WordPart{
							{
								Part:   &tk[3],
								Tokens: tk[3:4],
							},
						},
						Tokens: tk[3:4],
					},
					{
						Parts: []WordPart{
							{
								Part:   &tk[5],
								Tokens: tk[5:6],
							},
						},
						Tokens: tk[5:6],
					},
				},
				Tokens: tk[:7],
			}
		}},
		{"{\"$(||)\",bc,123}", func(t *test, tk Tokens) { // 5
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
				Parsing: "BraceExpansion",
				Token:   tk[1],
			}
		}},
	}, func(t *test) (Type, error) {
		var b BraceExpansion

		err := b.parse(t.Parser)

		return b, err
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
					Array: []WordOrOperator{
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
				BraceWord: &BraceWord{
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
				BraceWord: &BraceWord{
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
				BraceWord: &BraceWord{
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
		{"${a:+b c}", func(t *test, tk Tokens) { // 13
			t.Output = ParameterExpansion{
				Type: ParameterMessage,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				BraceWord: &BraceWord{
					Parts: []WordPart{
						{
							Part:   &tk[3],
							Tokens: tk[3:4],
						},
						{
							Part:   &tk[4],
							Tokens: tk[4:5],
						},
						{
							Part:   &tk[5],
							Tokens: tk[5:6],
						},
					},
					Tokens: tk[3:6],
				},
				Tokens: tk[:7],
			}
		}},
		{"${a:-b}", func(t *test, tk Tokens) { // 14
			t.Output = ParameterExpansion{
				Type: ParameterSetAssign,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				BraceWord: &BraceWord{
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
		{"${a=b}", func(t *test, tk Tokens) { // 15
			t.Output = ParameterExpansion{
				Type: ParameterUnsetSubstitution,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				BraceWord: &BraceWord{
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
		{"${a?b}", func(t *test, tk Tokens) { // 16
			t.Output = ParameterExpansion{
				Type: ParameterUnsetAssignment,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				BraceWord: &BraceWord{
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
		{"${a+b}", func(t *test, tk Tokens) { // 17
			t.Output = ParameterExpansion{
				Type: ParameterUnsetMessage,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				BraceWord: &BraceWord{
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
		{"${a-b}", func(t *test, tk Tokens) { // 18
			t.Output = ParameterExpansion{
				Type: ParameterUnsetSetAssign,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				BraceWord: &BraceWord{
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
		{"${a:1}", func(t *test, tk Tokens) { // 19
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
		{"${a: 1}", func(t *test, tk Tokens) { // 20
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
		{"${a: -1}", func(t *test, tk Tokens) { // 21
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
		{"${a:1:2}", func(t *test, tk Tokens) { // 22
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
		{"${a:1:-2}", func(t *test, tk Tokens) { // 23
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
		{"${a:1: -2}", func(t *test, tk Tokens) { // 24
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
		{"${a#b}", func(t *test, tk Tokens) { // 25
			t.Output = ParameterExpansion{
				Type: ParameterRemoveStartShortest,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				BraceWord: &BraceWord{
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
		{"${a##b}", func(t *test, tk Tokens) { // 26
			t.Output = ParameterExpansion{
				Type: ParameterRemoveStartLongest,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				BraceWord: &BraceWord{
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
		{"${a%b}", func(t *test, tk Tokens) { // 27
			t.Output = ParameterExpansion{
				Type: ParameterRemoveEndShortest,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				BraceWord: &BraceWord{
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
		{"${a%%b}", func(t *test, tk Tokens) { // 28
			t.Output = ParameterExpansion{
				Type: ParameterRemoveEndLongest,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				BraceWord: &BraceWord{
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
		{"${a/b}", func(t *test, tk Tokens) { // 29
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
		{"${a/b/c}", func(t *test, tk Tokens) { // 30
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
		{"${a//b}", func(t *test, tk Tokens) { // 31
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
		{"${a//b/c}", func(t *test, tk Tokens) { // 32
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
		{"${a/#b}", func(t *test, tk Tokens) { // 33
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
		{"${a/#b/c}", func(t *test, tk Tokens) { // 34
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
		{"${a/%b}", func(t *test, tk Tokens) { // 35
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
		{"${a/%b/c}", func(t *test, tk Tokens) { // 36
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
		{"${a^b}", func(t *test, tk Tokens) { // 37
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
		{"${a^^b}", func(t *test, tk Tokens) { // 38
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
		{"${a,b}", func(t *test, tk Tokens) { // 39
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
		{"${a,,b}", func(t *test, tk Tokens) { // 40
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
		{"${!a@}", func(t *test, tk Tokens) { // 41
			t.Output = ParameterExpansion{
				Type: ParameterPrefixSeperate,
				Parameter: Parameter{
					Parameter: &tk[2],
					Tokens:    tk[2:3],
				},
				Tokens: tk[:5],
			}
		}},
		{"${!a*}", func(t *test, tk Tokens) { // 42
			t.Output = ParameterExpansion{
				Type: ParameterPrefix,
				Parameter: Parameter{
					Parameter: &tk[2],
					Tokens:    tk[2:3],
				},
				Tokens: tk[:5],
			}
		}},
		{"${a@U}", func(t *test, tk Tokens) { // 43
			t.Output = ParameterExpansion{
				Type: ParameterUppercase,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Tokens: tk[:5],
			}
		}},
		{"${a@u}", func(t *test, tk Tokens) { // 44
			t.Output = ParameterExpansion{
				Type: ParameterUppercaseFirst,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Tokens: tk[:5],
			}
		}},
		{"${a@L}", func(t *test, tk Tokens) { // 45
			t.Output = ParameterExpansion{
				Type: ParameterLowercase,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Tokens: tk[:5],
			}
		}},
		{"${a@Q}", func(t *test, tk Tokens) { // 46
			t.Output = ParameterExpansion{
				Type: ParameterQuoted,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Tokens: tk[:5],
			}
		}},
		{"${a@E}", func(t *test, tk Tokens) { // 47
			t.Output = ParameterExpansion{
				Type: ParameterEscaped,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Tokens: tk[:5],
			}
		}},
		{"${a@P}", func(t *test, tk Tokens) { // 48
			t.Output = ParameterExpansion{
				Type: ParameterPrompt,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Tokens: tk[:5],
			}
		}},
		{"${a@A}", func(t *test, tk Tokens) { // 49
			t.Output = ParameterExpansion{
				Type: ParameterDeclare,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Tokens: tk[:5],
			}
		}},
		{"${a@K}", func(t *test, tk Tokens) { // 50
			t.Output = ParameterExpansion{
				Type: ParameterQuotedArrays,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Tokens: tk[:5],
			}
		}},
		{"${a@a}", func(t *test, tk Tokens) { // 51
			t.Output = ParameterExpansion{
				Type: ParameterAttributes,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Tokens: tk[:5],
			}
		}},
		{"${a@k}", func(t *test, tk Tokens) { // 52
			t.Output = ParameterExpansion{
				Type: ParameterQuotedArraysSeperate,
				Parameter: Parameter{
					Parameter: &tk[1],
					Tokens:    tk[1:2],
				},
				Tokens: tk[:5],
			}
		}},
		{"${a[$(||)]}", func(t *test, tk Tokens) { // 53
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
						Parsing: "WordOrOperator",
						Token:   tk[3],
					},
					Parsing: "Parameter",
					Token:   tk[3],
				},
				Parsing: "ParameterExpansion",
				Token:   tk[1],
			}
		}},
		{"${a:=$(||)}", func(t *test, tk Tokens) { // 54
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
					Parsing: "BraceWord",
					Token:   tk[3],
				},
				Parsing: "ParameterExpansion",
				Token:   tk[3],
			}
		}},
		{"${a:1:2b}", func(t *test, tk Tokens) { // 55
			t.Err = Error{
				Err:     ErrMissingClosingBrace,
				Parsing: "ParameterExpansion",
				Token:   tk[6],
			}
		}},
		{"${a/b/$(||)}", func(t *test, tk Tokens) { // 56
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
				Array: []WordOrOperator{
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
				Tokens: tk[1:5],
			}
		}},
		{"${a[ 0 ]}", func(t *test, tk Tokens) { // 7
			t.Output = Parameter{
				Parameter: &tk[1],
				Array: []WordOrOperator{
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
				Tokens: tk[1:7],
			}
		}},
		{"${a[@]}", func(t *test, tk Tokens) { // 8
			t.Output = Parameter{
				Parameter: &tk[1],
				Array: []WordOrOperator{
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
				Tokens: tk[1:5],
			}
		}},
		{"${a[*]}", func(t *test, tk Tokens) { // 9
			t.Output = Parameter{
				Parameter: &tk[1],
				Array: []WordOrOperator{
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
					Parsing: "WordOrOperator",
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
		{"a=b", func(t *test, tk Tokens) { // 2
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
		{"$a", func(t *test, tk Tokens) { // 3
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
		{">", func(t *test, tk Tokens) { // 4
			t.Output = WordOrOperator{
				Operator: &tk[0],
				Tokens:   tk[:1],
			}
		}},
		{"&", func(t *test, tk Tokens) { // 5
			t.Output = WordOrOperator{
				Operator: &tk[0],
				Tokens:   tk[:1],
			}
		}},
		{"$(||)", func(t *test, tk Tokens) { // 6
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
