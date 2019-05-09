package base83

import "errors"

// ErrInvalidInput is returned when the given input to decode is not valid base83.
var ErrInvalidInput = errors.New("base83: invalid input")

func invalidError(r rune, i int) error {
	// No stdlib support for wrapped errors, so return as-is pre-1.13
	return ErrInvalidInput
}
