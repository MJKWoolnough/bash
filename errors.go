package bash

import "errors"

// Errors.
var (
	ErrInvalidCharacter          = errors.New("invalid character")
	ErrInvalidParameterExpansion = errors.New("invalid parameter expansion")
	ErrInvalidNumber             = errors.New("invalid number")
	ErrInvalidAssignment         = errors.New("invalid assignment")
	ErrMissingClosingBracket     = errors.New("missing closing bracket")
	ErrMissingClosingBrace       = errors.New("missing closing brace")
	ErrMissingClosingParen       = errors.New("missing closing paren")
	ErrMissingCloser             = errors.New("missing closer")
	ErrInvalidEndOfStatement     = errors.New("invalid end of statement")
	ErrIncorrectBacktick         = errors.New("incorrect backtick depth")
	ErrMissingWord               = errors.New("missing word")
	ErrMissingClosingIf          = errors.New("missing if closing")
	ErrMissingThen               = errors.New("missing then")
	ErrMissingIn                 = errors.New("missing in")
	ErrMissingDo                 = errors.New("missing do")
	ErrMissingClosingCase        = errors.New("missing case closing")
	ErrMissingClosingPattern     = errors.New("missing pattern closing")
	ErrInvalidKeyword            = errors.New("invalid keyword")
	ErrInvalidIdentifier         = errors.New("invalid identifier")
	ErrMissingOperator           = errors.New("missing operator")
	ErrInvalidOperator           = errors.New("invalid operator")
)
