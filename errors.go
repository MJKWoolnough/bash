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
	ErrMissingCloser             = errors.New("missing closer")
	ErrInvalidEndOfStatement     = errors.New("invalid end of statement")
	ErrIncorrectBacktick         = errors.New("incorrect backtick depth")
	ErrMissingWord               = errors.New("missing word")
	ErrMissingClosingIf          = errors.New("missing if closing")
)
