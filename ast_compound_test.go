package bash

import "testing"

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
		{"if a; then b; else\n# comment\nc; fi", func(t *test, tk Tokens) { // 5
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
		{"if a; then b; else # comment\nc; fi", func(t *test, tk Tokens) { // 6
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
		{"if ||;then b;fi", func(t *test, tk Tokens) { // 7
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
		{"if a;then b;elif ||;then d;fi", func(t *test, tk Tokens) { // 8
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
		{"if a;then b;else ||;fi", func(t *test, tk Tokens) { // 9
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
		{"if a; then # comment\nb;fi", func(t *test, tk Tokens) { // 5
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
					Comments: [2]Comments{{tk[7]}},
					Tokens:   tk[7:11],
				},
				Tokens: tk[2:11],
			}
		}},
		{"if a; then\n# comment\nb;fi", func(t *test, tk Tokens) { // 6
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
							Comments: [2]Comments{{tk[7]}},
							Tokens:   tk[7:11],
						},
					},
					Tokens: tk[7:11],
				},
				Tokens: tk[2:11],
			}
		}},
		{"if ||; then b;fi", func(t *test, tk Tokens) { // 7
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
		{"if a; then ||;fi", func(t *test, tk Tokens) { // 8
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
		{"case a # comment\nin b)c;;esac", func(t *test, tk Tokens) { // 5
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
										Part:   &tk[8],
										Tokens: tk[8:9],
									},
								},
								Tokens: tk[8:9],
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
											Tokens: tk[10:11],
										},
									},
									Tokens: tk[10:11],
								},
							},
							Tokens: tk[10:11],
						},
						CaseTerminationType: CaseTerminationEnd,
						Tokens:              tk[8:12],
					},
				},
				Comments: [3]Comments{{tk[4]}},
				Tokens:   tk[:13],
			}
		}},
		{"case a in # comment\nb)c;;esac", func(t *test, tk Tokens) { // 6
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
										Part:   &tk[8],
										Tokens: tk[8:9],
									},
								},
								Tokens: tk[8:9],
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
											Tokens: tk[10:11],
										},
									},
									Tokens: tk[10:11],
								},
							},
							Tokens: tk[10:11],
						},
						CaseTerminationType: CaseTerminationEnd,
						Tokens:              tk[8:12],
					},
				},
				Comments: [3]Comments{nil, {tk[6]}},
				Tokens:   tk[:13],
			}
		}},
		{"case a in\n# comment\nb)c;;esac", func(t *test, tk Tokens) { // 7
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
										Part:   &tk[8],
										Tokens: tk[8:9],
									},
								},
								Tokens: tk[8:9],
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
											Tokens: tk[10:11],
										},
									},
									Tokens: tk[10:11],
								},
							},
							Tokens: tk[10:11],
						},
						CaseTerminationType: CaseTerminationEnd,
						Comments:            Comments{tk[6]},
						Tokens:              tk[6:12],
					},
				},
				Tokens: tk[:13],
			}
		}},
		{"case a in b)c;; # comment\nesac", func(t *test, tk Tokens) { // 8
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
				Comments: [3]Comments{nil, nil, {tk[11]}},
				Tokens:   tk[:14],
			}
		}},
		{"case $(||) in b)c;esac", func(t *test, tk Tokens) { // 9
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
		{"case a in b)||;esac", func(t *test, tk Tokens) { // 10
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
		{"case a in a)# comment\nb;;# comment2\nesac", func(t *test, tk Tokens) { // 7
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
									Tokens: tk[10:11],
								},
							},
							Tokens: tk[10:11],
						},
					},
					Comments: [2]Comments{{tk[8]}},
					Tokens:   tk[8:11],
				},
				CaseTerminationType: CaseTerminationEnd,
				Tokens:              tk[6:12],
			}
		}},
		{"case a in a)# comment\nb;\n# comment2\nesac", func(t *test, tk Tokens) { // 8
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
					Comments: [2]Comments{{tk[8]}, {tk[13]}},
					Tokens:   tk[8:14],
				},
				Tokens: tk[6:14],
			}
		}},
		{"case a in\n# comment\na)b;;esac", func(t *test, tk Tokens) { // 9
			t.Output = PatternLines{
				Patterns: []Word{
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
									Tokens: tk[10:11],
								},
							},
							Tokens: tk[10:11],
						},
					},
					Tokens: tk[10:11],
				},
				CaseTerminationType: CaseTerminationEnd,
				Comments:            Comments{tk[6]},
				Tokens:              tk[6:12],
			}
		}},
		{"case a in $(||))d;;esac", func(t *test, tk Tokens) { // 10
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
		{"case a in a|\nb)c;;esac", func(t *test, tk Tokens) { // 11
			t.Err = Error{
				Err:     ErrMissingClosingPattern,
				Parsing: "PatternLines",
				Token:   tk[8],
			}
		}},
		{"case a in a)||;;esac", func(t *test, tk Tokens) { // 12
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
		{"while a # comment\ndo b; done", func(t *test, tk Tokens) { // 3
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
				Tokens:   tk[:12],
			}
		}},
		{"until a; do # comment\nb; done", func(t *test, tk Tokens) { // 4
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
					Comments: [2]Comments{{tk[7]}},
					Tokens:   tk[7:11],
				},
				Tokens: tk[:13],
			}
		}},
		{"while ||; do b; done", func(t *test, tk Tokens) { // 5
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
		{"until a; do ||; done", func(t *test, tk Tokens) { // 6
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
		{"for a #comment\ndo b;done", func(t *test, tk Tokens) { // 7
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
				Comments: [2]Comments{{tk[4]}},
				Tokens:   tk[:11],
			}
		}},
		{"for a # comment\nin b #comment 2\ndo c;done", func(t *test, tk Tokens) { // 8
			t.Output = ForCompound{
				Identifier: &tk[2],
				Words: []Word{
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
					Tokens: tk[14:16],
				},
				Comments: [2]Comments{{tk[4]}, {tk[10]}},
				Tokens:   tk[:17],
			}
		}},
		{"for a in b #comment\ndo c;done", func(t *test, tk Tokens) { // 9
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
				Comments: [2]Comments{nil, {tk[8]}},
				Tokens:   tk[:15],
			}
		}},
		{"for (( a=1; a<2; a++ )) #comment\ndo b;done", func(t *test, tk Tokens) { // 10
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
																	Part:   &tk[23],
																	Tokens: tk[23:24],
																},
															},
															Tokens: tk[23:24],
														},
														Tokens: tk[23:24],
													},
												},
												Tokens: tk[23:24],
											},
											Tokens: tk[23:24],
										},
										Tokens: tk[23:24],
									},
									Tokens: tk[23:25],
								},
							},
							Tokens: tk[23:25],
						},
					},
					Tokens: tk[23:25],
				},
				Comments: [2]Comments{nil, {tk[19]}},
				Tokens:   tk[:26],
			}
		}},
		{"for (( a=1; a<2; a++ )) ;#comment\ndo b;done", func(t *test, tk Tokens) { // 11
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
																	Part:   &tk[24],
																	Tokens: tk[24:25],
																},
															},
															Tokens: tk[24:25],
														},
														Tokens: tk[24:25],
													},
												},
												Tokens: tk[24:25],
											},
											Tokens: tk[24:25],
										},
										Tokens: tk[24:25],
									},
									Tokens: tk[24:26],
								},
							},
							Tokens: tk[24:26],
						},
					},
					Tokens: tk[24:26],
				},
				Comments: [2]Comments{nil, {tk[20]}},
				Tokens:   tk[:27],
			}
		}},
		{"for a in $(||); do b;done", func(t *test, tk Tokens) { // 12
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
		{"for (( $(||) )); do b;done", func(t *test, tk Tokens) { // 13
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
		{"for a; do ||;done", func(t *test, tk Tokens) { // 14
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
		{"select a #comment\ndo b;done", func(t *test, tk Tokens) { // 6
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
				Comments: [2]Comments{{tk[4]}},
				Tokens:   tk[:11],
			}
		}},
		{"select a # comment\nin b #comment 2\ndo c;done", func(t *test, tk Tokens) { // 7
			t.Output = SelectCompound{
				Identifier: &tk[2],
				Words: []Word{
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
					Tokens: tk[14:16],
				},
				Comments: [2]Comments{{tk[4]}, {tk[10]}},
				Tokens:   tk[:17],
			}
		}},
		{"select a in b #comment\ndo c;done", func(t *test, tk Tokens) { // 8
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
				Comments: [2]Comments{nil, {tk[8]}},
				Tokens:   tk[:15],
			}
		}},
		{"select a in $(||); do b;done", func(t *test, tk Tokens) { // 9
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
		{"select a; do ||;done", func(t *test, tk Tokens) { // 10
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
		{"[[ #comment A\n\n#comment B\na = b ]]", func(t *test, tk Tokens) { // 2
			t.Output = TestCompound{
				Tests: Tests{
					Test: TestOperatorStringsEqual,
					Word: &Word{
						Parts: []WordPart{
							{
								Part:   &tk[6],
								Tokens: tk[6:7],
							},
						},
						Tokens: tk[6:7],
					},
					Pattern: &Pattern{
						Parts: []WordPart{
							{
								Part:   &tk[10],
								Tokens: tk[10:11],
							},
						},
						Tokens: tk[10:11],
					},
					Comments: [5]Comments{{tk[4]}},
					Tokens:   tk[4:11],
				},
				Comments: [2]Comments{{tk[2]}},
				Tokens:   tk[:13],
			}
		}},
		{"[[ a = b # comment A\n\n#comment B\n]]", func(t *test, tk Tokens) { // 3
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
					Comments: [5]Comments{nil, nil, nil, nil, {tk[8]}},
					Tokens:   tk[2:9],
				},
				Comments: [2]Comments{nil, {tk[10]}},
				Tokens:   tk[:13],
			}
		}},
		{"[[ -a $(||) ]]", func(t *test, tk Tokens) { // 4
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
		{"[[ a>b ]]", func(t *test, tk Tokens) { // 6
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
							Part:   &tk[4],
							Tokens: tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				Tokens: tk[2:5],
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
		{"[[\n#comment\na = b ]]", func(t *test, tk Tokens) { // 47
			t.Output = Tests{
				Test: TestOperatorStringsEqual,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[4],
							Tokens: tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				Pattern: &Pattern{
					Parts: []WordPart{
						{
							Part:   &tk[8],
							Tokens: tk[8:9],
						},
					},
					Tokens: tk[8:9],
				},
				Comments: [5]Comments{{tk[2]}},
				Tokens:   tk[2:9],
			}
		}},
		{"[[ a = b # comment\n ]]", func(t *test, tk Tokens) { // 48
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
				Comments: [5]Comments{nil, nil, nil, nil, {tk[8]}},
				Tokens:   tk[2:9],
			}
		}},
		{"[[ #comment A\n! #comment B\n-z a ]]", func(t *test, tk Tokens) { // 49
			t.Output = Tests{
				Not:  true,
				Test: TestOperatorStringIsZero,
				Word: &Word{
					Parts: []WordPart{
						{
							Part:   &tk[10],
							Tokens: tk[10:11],
						},
					},
					Tokens: tk[10:11],
				},
				Comments: [5]Comments{{tk[2]}, {tk[6]}},
				Tokens:   tk[2:11],
			}
		}},
		{"[[ ( #comment A\n\n#comment B\na = b #comment C\n\n#comment D\n) ]]", func(t *test, tk Tokens) { // 50
			t.Output = Tests{
				Parens: &Tests{
					Test: TestOperatorStringsEqual,
					Word: &Word{
						Parts: []WordPart{
							{
								Part:   &tk[8],
								Tokens: tk[8:9],
							},
						},
						Tokens: tk[8:9],
					},
					Pattern: &Pattern{
						Parts: []WordPart{
							{
								Part:   &tk[12],
								Tokens: tk[12:13],
							},
						},
						Tokens: tk[12:13],
					},
					Comments: [5]Comments{{tk[6]}, nil, nil, nil, {tk[14]}},
					Tokens:   tk[6:15],
				},
				Comments: [5]Comments{nil, nil, {tk[4]}, {tk[16]}},
				Tokens:   tk[2:19],
			}
		}},
		{"[[ -a $(||) ]]", func(t *test, tk Tokens) { // 51
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
		{"[[ $(||) = a ]]", func(t *test, tk Tokens) { // 52
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
		{"[[ a = $(||) ]]", func(t *test, tk Tokens) { // 53
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
		{"[[ a -eq $(||) ]]", func(t *test, tk Tokens) { // 54
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
		{"[[ ( $(||) ) ]]", func(t *test, tk Tokens) { // 55
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
		{"[[ ( a b ) ]]", func(t *test, tk Tokens) { // 56
			t.Err = Error{
				Err:     ErrMissingClosingParen,
				Parsing: "Tests",
				Token:   tk[6],
			}
		}},
		{"[[ a || $(||) ]]", func(t *test, tk Tokens) { // 57
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
		{"( # comment A\n\n# comment B\na # comment C\n\n# comment D\n)", func(t *test, tk Tokens) { // 7
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
							Comments: [2]Comments{{tk[4]}, {tk[8]}},
							Tokens:   tk[4:9],
						},
					},
					Comments: [2]Comments{{tk[2]}, {tk[10]}},
					Tokens:   tk[2:11],
				},
				Tokens: tk[:13],
			}
		}},
		{"{ # comment A\n\n# comment B\na # comment C\n\n# comment D\n}", func(t *test, tk Tokens) { // 8
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
							Comments: [2]Comments{{tk[4]}, {tk[8]}},
							Tokens:   tk[4:9],
						},
					},
					Comments: [2]Comments{{tk[2]}, {tk[10]}},
					Tokens:   tk[2:11],
				},
				Tokens: tk[:13],
			}
		}},
		{"(||)", func(t *test, tk Tokens) { // 9
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
		{"function a # comment\n{ b; }", func(t *test, tk Tokens) { // 5
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
				Comments: Comments{tk[4]},
				Tokens:   tk[:12],
			}
		}},
		{"function a() # comment\n{ b; }", func(t *test, tk Tokens) { // 6
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
				Comments: Comments{tk[6]},
				Tokens:   tk[:14],
			}
		}},
		{"a() #comment\n{ b; }", func(t *test, tk Tokens) { // 7
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
				Comments: Comments{tk[4]},
				Tokens:   tk[:12],
			}
		}},
		{"function a() { || }", func(t *test, tk Tokens) { // 8
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
