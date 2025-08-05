package walk

import "vimagination.zapto.org/bash"

type Handler interface {
	Handle(bash.Type) error
}

type HandlerFunc func(bash.Type) error

func (h HandlerFunc) Handle(t bash.Type) error {
	return h(t)
}

func Walk(t bash.Type, fn Handler) error {
	switch t := t.(type) {
	case bash.ArithmeticExpansion:
		return walkArithmeticExpansion(&t, fn)
	case *bash.ArithmeticExpansion:
		return walkArithmeticExpansion(t, fn)
	case bash.ArrayWord:
		return walkArrayWord(&t, fn)
	case *bash.ArrayWord:
		return walkArrayWord(t, fn)
	case bash.Assignment:
		return walkAssignment(&t, fn)
	case *bash.Assignment:
		return walkAssignment(t, fn)
	case bash.AssignmentOrWord:
		return walkAssignmentOrWord(&t, fn)
	case *bash.AssignmentOrWord:
		return walkAssignmentOrWord(t, fn)
	case bash.BraceExpansion:
		return walkBraceExpansion(&t, fn)
	case *bash.BraceExpansion:
		return walkBraceExpansion(t, fn)
	case bash.BraceWord:
		return walkBraceWord(&t, fn)
	case *bash.BraceWord:
		return walkBraceWord(t, fn)
	case bash.CaseCompound:
		return walkCaseCompound(&t, fn)
	case *bash.CaseCompound:
		return walkCaseCompound(t, fn)
	case bash.Command:
		return walkCommand(&t, fn)
	case *bash.Command:
		return walkCommand(t, fn)
	case bash.CommandOrCompound:
		return walkCommandOrCompound(&t, fn)
	case *bash.CommandOrCompound:
		return walkCommandOrCompound(t, fn)
	case bash.CommandSubstitution:
		return walkCommandSubstitution(&t, fn)
	case *bash.CommandSubstitution:
		return walkCommandSubstitution(t, fn)
	case bash.Compound:
		return walkCompound(&t, fn)
	case *bash.Compound:
		return walkCompound(t, fn)
	case bash.File:
		return walkFile(&t, fn)
	case *bash.File:
		return walkFile(t, fn)
	case bash.ForCompound:
		return walkForCompound(&t, fn)
	case *bash.ForCompound:
		return walkForCompound(t, fn)
	case bash.FunctionCompound:
		return walkFunctionCompound(&t, fn)
	case *bash.FunctionCompound:
		return walkFunctionCompound(t, fn)
	case bash.GroupingCompound:
		return walkGroupingCompound(&t, fn)
	case *bash.GroupingCompound:
		return walkGroupingCompound(t, fn)
	case bash.Heredoc:
		return walkHeredoc(&t, fn)
	case *bash.Heredoc:
		return walkHeredoc(t, fn)
	case bash.HeredocPartOrWord:
		return walkHeredocPartOrWord(&t, fn)
	case *bash.HeredocPartOrWord:
		return walkHeredocPartOrWord(t, fn)
	case bash.IfCompound:
		return walkIfCompound(&t, fn)
	case *bash.IfCompound:
		return walkIfCompound(t, fn)
	case bash.Line:
		return walkLine(&t, fn)
	case *bash.Line:
		return walkLine(t, fn)
	case bash.LoopCompound:
		return walkLoopCompound(&t, fn)
	case *bash.LoopCompound:
		return walkLoopCompound(t, fn)
	case bash.Parameter:
		return walkParameter(&t, fn)
	case *bash.Parameter:
		return walkParameter(t, fn)
	case bash.ParameterAssign:
		return walkParameterAssign(&t, fn)
	case *bash.ParameterAssign:
		return walkParameterAssign(t, fn)
	case bash.ParameterExpansion:
		return walkParameterExpansion(&t, fn)
	case *bash.ParameterExpansion:
		return walkParameterExpansion(t, fn)
	case bash.Pattern:
		return walkPattern(&t, fn)
	case *bash.Pattern:
		return walkPattern(t, fn)
	case bash.PatternLines:
		return walkPatternLines(&t, fn)
	case *bash.PatternLines:
		return walkPatternLines(t, fn)
	case bash.Pipeline:
		return walkPipeline(&t, fn)
	case *bash.Pipeline:
		return walkPipeline(t, fn)
	case bash.Redirection:
		return walkRedirection(&t, fn)
	case *bash.Redirection:
		return walkRedirection(t, fn)
	case bash.SelectCompound:
		return walkSelectCompound(&t, fn)
	case *bash.SelectCompound:
		return walkSelectCompound(t, fn)
	case bash.Statement:
		return walkStatement(&t, fn)
	case *bash.Statement:
		return walkStatement(t, fn)
	case bash.String:
		return walkString(&t, fn)
	case *bash.String:
		return walkString(t, fn)
	case bash.TestCompound:
		return walkTestCompound(&t, fn)
	case *bash.TestCompound:
		return walkTestCompound(t, fn)
	case bash.TestConsequence:
		return walkTestConsequence(&t, fn)
	case *bash.TestConsequence:
		return walkTestConsequence(t, fn)
	case bash.Tests:
		return walkTests(&t, fn)
	case *bash.Tests:
		return walkTests(t, fn)
	case bash.Value:
		return walkValue(&t, fn)
	case *bash.Value:
		return walkValue(t, fn)
	case bash.Word:
		return walkWord(&t, fn)
	case *bash.Word:
		return walkWord(t, fn)
	case bash.WordOrOperator:
		return walkWordOrOperator(&t, fn)
	case *bash.WordOrOperator:
		return walkWordOrOperator(t, fn)
	case bash.WordOrToken:
		return walkWordOrToken(&t, fn)
	case *bash.WordOrToken:
		return walkWordOrToken(t, fn)
	case bash.WordPart:
		return walkWordPart(&t, fn)
	case *bash.WordPart:
		return walkWordPart(t, fn)
	}

	return nil
}

func walkArithmeticExpansion(t *bash.ArithmeticExpansion, fn Handler) error {
	for n := range t.WordsAndOperators {
		if err := fn.Handle(&t.WordsAndOperators[n]); err != nil {
			return err
		}
	}

	return nil
}

func walkArrayWord(t *bash.ArrayWord, fn Handler) error {
	return fn.Handle(&t.Word)
}

func walkAssignment(t *bash.Assignment, fn Handler) error {
	if err := fn.Handle(&t.Identifier); err != nil {
		return err
	}

	for n := range t.Expression {
		if err := fn.Handle(&t.Expression[n]); err != nil {
			return err
		}
	}

	if t.Value != nil {
		return fn.Handle(t.Value)
	}

	return nil
}

func walkAssignmentOrWord(t *bash.AssignmentOrWord, fn Handler) error {
	if t.Assignment != nil {
		return fn.Handle(t.Assignment)
	}

	if t.Word != nil {
		return fn.Handle(t.Word)
	}

	return nil
}

func walkBraceExpansion(t *bash.BraceExpansion, fn Handler) error {
	for n := range t.Words {
		if err := fn.Handle(&t.Words[n]); err != nil {
			return err
		}
	}

	return nil
}

func walkBraceWord(t *bash.BraceWord, fn Handler) error {
	for n := range t.Parts {
		if err := fn.Handle(&t.Parts[n]); err != nil {
			return err
		}
	}

	return nil
}

func walkCaseCompound(t *bash.CaseCompound, fn Handler) error {
	if err := fn.Handle(&t.Word); err != nil {
		return err
	}

	for n := range t.Matches {
		if err := fn.Handle(&t.Matches[n]); err != nil {
			return err
		}
	}

	return nil
}

func walkCommand(t *bash.Command, fn Handler) error {
	for n := range t.Vars {
		if err := fn.Handle(&t.Vars[n]); err != nil {
			return err
		}
	}

	for n := range t.AssignmentsOrWords {
		if err := fn.Handle(&t.AssignmentsOrWords[n]); err != nil {
			return err
		}
	}

	for n := range t.Redirections {
		if err := fn.Handle(&t.Redirections[n]); err != nil {
			return err
		}
	}

	return nil
}

func walkCommandOrCompound(t *bash.CommandOrCompound, fn Handler) error {
	if t.Command != nil {
		return fn.Handle(t.Command)
	}

	if t.Compound != nil {
		return fn.Handle(t.Compound)
	}

	return nil
}

func walkCommandSubstitution(t *bash.CommandSubstitution, fn Handler) error {
	return fn.Handle(&t.Command)
}

func walkCompound(t *bash.Compound, fn Handler) error {
	if t.IfCompound != nil {
		if err := fn.Handle(t.IfCompound); err != nil {
			return err
		}
	} else if t.CaseCompound != nil {
		if err := fn.Handle(t.CaseCompound); err != nil {
			return err
		}
	} else if t.LoopCompound != nil {
		if err := fn.Handle(t.LoopCompound); err != nil {
			return err
		}
	} else if t.ForCompound != nil {
		if err := fn.Handle(t.ForCompound); err != nil {
			return err
		}
	} else if t.SelectCompound != nil {
		if err := fn.Handle(t.SelectCompound); err != nil {
			return err
		}
	} else if t.GroupingCompound != nil {
		if err := fn.Handle(t.GroupingCompound); err != nil {
			return err
		}
	} else if t.TestCompound != nil {
		if err := fn.Handle(t.TestCompound); err != nil {
			return err
		}
	} else if t.ArithmeticCompound != nil {
		if err := fn.Handle(t.ArithmeticCompound); err != nil {
			return err
		}
	} else if t.FunctionCompound != nil {
		if err := fn.Handle(t.FunctionCompound); err != nil {
			return err
		}
	}

	for n := range t.Redirections {
		if err := fn.Handle(&t.Redirections[n]); err != nil {
			return err
		}
	}

	return nil
}

func walkFile(t *bash.File, fn Handler) error {
	for n := range t.Lines {
		if err := fn.Handle(&t.Lines[n]); err != nil {
			return err
		}
	}

	return nil
}

func walkForCompound(t *bash.ForCompound, fn Handler) error {
	if t.Identifier != nil {
		for n := range t.Words {
			if err := fn.Handle(&t.Words[n]); err != nil {
				return err
			}
		}
	} else if t.ArithmeticExpansion != nil {
		if err := fn.Handle(t.ArithmeticExpansion); err != nil {
			return err
		}
	}

	return fn.Handle(&t.File)
}

func walkFunctionCompound(t *bash.FunctionCompound, fn Handler) error {
	return fn.Handle(&t.Body)
}

func walkGroupingCompound(t *bash.GroupingCompound, fn Handler) error {
	return fn.Handle(&t.File)
}

func walkHeredoc(t *bash.Heredoc, fn Handler) error {
	for n := range t.HeredocPartsOrWords {
		if err := fn.Handle(&t.HeredocPartsOrWords[n]); err != nil {
			return err
		}
	}

	return nil
}

func walkHeredocPartOrWord(t *bash.HeredocPartOrWord, fn Handler) error {
	if t.Word != nil {
		return fn.Handle(t.Word)
	}

	return nil
}

func walkIfCompound(t *bash.IfCompound, fn Handler) error {
	if err := fn.Handle(&t.If); err != nil {
		return err
	}

	for n := range t.ElIf {
		if err := fn.Handle(&t.ElIf[n]); err != nil {
			return err
		}
	}

	if t.Else != nil {
		return fn.Handle(t.Else)
	}

	return nil
}

func walkLine(t *bash.Line, fn Handler) error {
	for n := range t.Statements {
		if err := fn.Handle(&t.Statements[n]); err != nil {
			return err
		}
	}

	return nil
}

func walkLoopCompound(t *bash.LoopCompound, fn Handler) error {
	if err := fn.Handle(&t.Statement); err != nil {
		return err
	}

	return fn.Handle(&t.File)
}

func walkParameter(t *bash.Parameter, fn Handler) error {
	for n := range t.Array {
		if err := fn.Handle(&t.Array[n]); err != nil {
			return err
		}
	}

	return nil
}

func walkParameterAssign(t *bash.ParameterAssign, fn Handler) error {
	for n := range t.Subscript {
		if err := fn.Handle(&t.Subscript[n]); err != nil {
			return err
		}
	}

	return nil
}

func walkParameterExpansion(t *bash.ParameterExpansion, fn Handler) error {
	if err := fn.Handle(t.Parameter); err != nil {
		return err
	}

	if t.BraceWord != nil {
		return fn.Handle(t.BraceWord)
	} else if t.String != nil {
		return fn.Handle(t.String)
	}

	return nil
}

func walkPattern(t *bash.Pattern, fn Handler) error {
	for n := range t.Parts {
		if err := fn.Handle(&t.Parts[n]); err != nil {
			return err
		}
	}

	return nil
}

func walkPatternLines(t *bash.PatternLines, fn Handler) error {
	for n := range t.Patterns {
		if err := fn.Handle(&t.Patterns[n]); err != nil {
			return err
		}
	}

	return fn.Handle(t.Lines)
}

func walkPipeline(t *bash.Pipeline, fn Handler) error {
	if err := fn.Handle(&t.CommandOrCompound); err != nil {
		return err
	}

	if t.Pipeline != nil {
		return fn.Handle(t.Pipeline)
	}

	return nil
}

func walkRedirection(t *bash.Redirection, fn Handler) error {
	if err := fn.Handle(&t.Output); err != nil {
		return err
	}

	if t.Heredoc != nil {
		return fn.Handle(t.Heredoc)
	}

	return nil
}

func walkSelectCompound(t *bash.SelectCompound, fn Handler) error {
	for n := range t.Words {
		if err := fn.Handle(&t.Words[n]); err != nil {
			return err
		}
	}

	return fn.Handle(&t.File)
}

func walkStatement(t *bash.Statement, fn Handler) error {
	if err := fn.Handle(&t.Pipeline); err != nil {
		return err
	}

	if t.Statement != nil {
		return fn.Handle(t.Statement)
	}

	return nil
}

func walkString(t *bash.String, fn Handler) error {
	return nil
}

func walkTestCompound(t *bash.TestCompound, fn Handler) error {
	return nil
}

func walkTestConsequence(t *bash.TestConsequence, fn Handler) error {
	return nil
}

func walkTests(t *bash.Tests, fn Handler) error {
	return nil
}

func walkValue(t *bash.Value, fn Handler) error {
	return nil
}

func walkWord(t *bash.Word, fn Handler) error {
	return nil
}

func walkWordOrOperator(t *bash.WordOrOperator, fn Handler) error {
	return nil
}

func walkWordOrToken(t *bash.WordOrToken, fn Handler) error {
	return nil
}

func walkWordPart(t *bash.WordPart, fn Handler) error {
	return nil
}
