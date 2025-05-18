package bash

import (
	"errors"
	"testing"
)

func TestError(t *testing.T) {
	const expected = "Something: error at position 2 (4:3):\ninvalid character"
	err := Error{
		Err:     ErrInvalidCharacter,
		Parsing: "Something",
		Token:   Token{Pos: 1, LinePos: 2, Line: 3},
	}

	if errStr := err.Error(); errStr != expected {
		t.Errorf("expecting error string %q, got %q", expected, errStr)
	} else if under := errors.Unwrap(err); under != ErrInvalidCharacter {
		t.Errorf("expecting underlying error to be ErrInvalidCharacter, got %s", under)
	}
}
