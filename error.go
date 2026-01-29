package blurhash

import (
	"errors"
)

var (
	// ErrInvalidComponents is returned when components passed to Encode are invalid.
	ErrInvalidComponents = errors.New("blurhash: must have between 1 and 9 components")
	// ErrInvalidHash is returned when the library encounters a hash it can't recognise.
	ErrInvalidHash = errors.New("blurhash: invalid hash")
	// ErrInvalidDimensions is returned when width or height is invalid.
	ErrInvalidDimensions = errors.New("blurhash: width and height must be positive")
)
