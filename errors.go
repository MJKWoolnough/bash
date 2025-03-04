package bash

import "errors"

// Errors.
var (
	ErrInvalidBraceExpansion = errors.New("invalid brace expansion")
	ErrInvalidCharacter      = errors.New("invalid character")
	ErrInvalidNumber         = errors.New("invalid number")
	ErrInvalidAssignment     = errors.New("invalid assignment")
	ErrMissingClosingBracket = errors.New("missing closing bracket")
	ErrInvalidEndOfStatement = errors.New("invalid end of statement")
)
