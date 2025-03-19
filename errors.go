package bash

import "errors"

// Errors.
var (
	ErrInvalidCharacter          = errors.New("invalid character")
	ErrInvalidParameterExpansion = errors.New("invalid parameter expansion")
	ErrInvalidNumber             = errors.New("invalid number")
	ErrInvalidAssignment         = errors.New("invalid assignment")
	ErrMissingClosingBracket     = errors.New("missing closing bracket")
	ErrInvalidEndOfStatement     = errors.New("invalid end of statement")
	ErrIncorrectBacktick         = errors.New("incorrect backtick depth")
)
