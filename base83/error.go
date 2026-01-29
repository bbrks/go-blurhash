package base83

import (
	"errors"
)

var (
	// ErrInvalidInput is returned when the given input to decode is not valid base83.
	ErrInvalidInput = errors.New("base83: invalid input")
)
