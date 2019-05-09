// +build !go1.13

package blurhash

import "errors"

// ErrInvalidHash is returned when the library encounters a hash it can't recognise.
var ErrInvalidHash = errors.New("blurhash: invalid hash")

func lengthError(expectedLength, actualLength int) error {
	return ErrInvalidHash
}
