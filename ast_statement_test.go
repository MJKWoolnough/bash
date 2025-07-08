package bash

import "testing"

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
						Value: &Value{
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
					{
						Identifier: ParameterAssign{
							Identifier: &tk[4],
							Tokens:     tk[4:5],
						},
						Value: &Value{
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
							Value: &Value{
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
							Value: &Value{
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
							Value: &Value{
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
							Value: &Value{
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
							Value: &Value{
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
							Value: &Value{
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
							Value: &Value{
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
		{"let a=b c=(d)", func(t *test, tk Tokens) { // 11
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
						Assignment: &Assignment{
							Identifier: ParameterAssign{
								Identifier: &tk[2],
								Tokens:     tk[2:3],
							},
							Assignment: AssignmentAssign,
							Expression: []WordOrOperator{
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
							Tokens: tk[2:5],
						},
						Tokens: tk[2:5],
					},
					{
						Assignment: &Assignment{
							Identifier: ParameterAssign{
								Identifier: &tk[6],
								Tokens:     tk[6:7],
							},
							Assignment: AssignmentAssign,
							Expression: []WordOrOperator{
								{
									Operator: &tk[8],
									Tokens:   tk[8:9],
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
							},
							Tokens: tk[6:11],
						},
						Tokens: tk[6:11],
					},
				},
				Tokens: tk[:11],
			}
		}},
		{"let a=(b ? c : d);", func(t *test, tk Tokens) { // 12
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
						Assignment: &Assignment{
							Identifier: ParameterAssign{
								Identifier: &tk[2],
								Tokens:     tk[2:3],
							},
							Assignment: AssignmentAssign,
							Expression: []WordOrOperator{
								{
									Operator: &tk[4],
									Tokens:   tk[4:5],
								},
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
									Operator: &tk[11],
									Tokens:   tk[11:12],
								},
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
								{
									Operator: &tk[14],
									Tokens:   tk[14:15],
								},
							},
							Tokens: tk[2:15],
						},
						Tokens: tk[2:15],
					},
				},
				Tokens: tk[:15],
			}
		}},
		{"a[$(||)]=", func(t *test, tk Tokens) { // 13
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
							Parsing: "WordOrOperator",
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
		{">$(||)", func(t *test, tk Tokens) { // 14
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
		{"$(||)", func(t *test, tk Tokens) { // 15
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
		{"a >$(||)", func(t *test, tk Tokens) { // 16
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
		{"export a=$(||)", func(t *test, tk Tokens) { // 17
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
		{"let a=$(||)", func(t *test, tk Tokens) { // 18
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
							Parsing: "WordOrOperator",
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
				Value: &Value{
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
						Parsing: "WordOrOperator",
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
				Subscript: []WordOrOperator{
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
				Tokens: tk[:4],
			}
		}},
		{"a[$a]=", func(t *test, tk Tokens) { // 3
			t.Output = ParameterAssign{
				Identifier: &tk[0],
				Subscript: []WordOrOperator{
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
					Parsing: "WordOrOperator",
					Token:   tk[2],
				},
				Parsing: "ParameterAssign",
				Token:   tk[2],
			}
		}},
	}, func(t *test) (Type, error) {
		var pa ParameterAssign

		err := pa.parse(t.Parser)

		return pa, err
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
