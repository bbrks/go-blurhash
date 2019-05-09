// +build !go1.13

package base83

import "errors"

var ErrInvalidInput = errors.New("base83: invalid input")
var ErrInvalidLength = errors.New("base83: invalid length")

// invalidError returns ErrInvalidInput for the given rune and index
func invalidError(r rune, i int) error {
	// No stdlib support for wrapped errors, so return as-is pre-1.13
	return ErrInvalidInput
}
