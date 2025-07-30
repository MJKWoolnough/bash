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
		if err := Walk(&t.WordsAndOperators[n], fn); err != nil {
			return err
		}
	}

	return nil
}

func walkArrayWord(t *bash.ArrayWord, fn Handler) error {
	return Walk(&t.Word, fn)
}

func walkAssignment(t *bash.Assignment, fn Handler) error {
	if err := Walk(&t.Identifier, fn); err != nil {
		return err
	}

	for n := range t.Expression {
		if err := Walk(&t.Expression[n], fn); err != nil {
			return err
		}
	}

	if t.Value != nil {
		return Walk(t.Value, fn)
	}

	return nil
}

func walkAssignmentOrWord(t *bash.AssignmentOrWord, fn Handler) error {
	if t.Assignment != nil {
		return Walk(t.Assignment, fn)
	}

	if t.Word != nil {
		return Walk(t.Word, fn)
	}

	return nil
}

func walkBraceExpansion(t *bash.BraceExpansion, fn Handler) error {
	for n := range t.Words {
		if err := Walk(&t.Words[n], fn); err != nil {
			return err
		}
	}

	return nil
}

func walkBraceWord(t *bash.BraceWord, fn Handler) error {
	return nil
}

func walkCaseCompound(t *bash.CaseCompound, fn Handler) error {
	return nil
}

func walkCommand(t *bash.Command, fn Handler) error {
	return nil
}

func walkCommandOrCompound(t *bash.CommandOrCompound, fn Handler) error {
	return nil
}

func walkCommandSubstitution(t *bash.CommandSubstitution, fn Handler) error {
	return nil
}

func walkCompound(t *bash.Compound, fn Handler) error {
	return nil
}

func walkFile(t *bash.File, fn Handler) error {
	return nil
}

func walkForCompound(t *bash.ForCompound, fn Handler) error {
	return nil
}

func walkFunctionCompound(t *bash.FunctionCompound, fn Handler) error {
	return nil
}

func walkGroupingCompound(t *bash.GroupingCompound, fn Handler) error {
	return nil
}

func walkHeredoc(t *bash.Heredoc, fn Handler) error {
	return nil
}

func walkHeredocPartOrWord(t *bash.HeredocPartOrWord, fn Handler) error {
	return nil
}

func walkIfCompound(t *bash.IfCompound, fn Handler) error {
	return nil
}

func walkLine(t *bash.Line, fn Handler) error {
	return nil
}

func walkLoopCompound(t *bash.LoopCompound, fn Handler) error {
	return nil
}

func walkParameter(t *bash.Parameter, fn Handler) error {
	return nil
}

func walkParameterAssign(t *bash.ParameterAssign, fn Handler) error {
	return nil
}

func walkParameterExpansion(t *bash.ParameterExpansion, fn Handler) error {
	return nil
}

func walkPattern(t *bash.Pattern, fn Handler) error {
	return nil
}

func walkPatternLines(t *bash.PatternLines, fn Handler) error {
	return nil
}

func walkPipeline(t *bash.Pipeline, fn Handler) error {
	return nil
}

func walkRedirection(t *bash.Redirection, fn Handler) error {
	return nil
}

func walkSelectCompound(t *bash.SelectCompound, fn Handler) error {
	return nil
}

func walkStatement(t *bash.Statement, fn Handler) error {
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
