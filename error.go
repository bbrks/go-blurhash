package blurhash

import "errors"

// ErrInvalidComponents is returned when components passed to Encode are invalid.
var ErrInvalidComponents = errors.New("blurhash: must have between 1 and 9 components")

// ErrInvalidHash is returned when the library encounters a hash it can't recognise.
var ErrInvalidHash = errors.New("blurhash: invalid hash")

func lengthError(expectedLength, actualLength int) error {
	// No stdlib support for wrapped errors, so return as-is pre-1.13
	return ErrInvalidHash
}
