// +build go1.13

package blurhash

import "fmt"

func lengthError(expectedLength, actualLength int) error {
	return fmt.Errorf("expected hash length %d but got %d: %w", expectedLength, actualLength, ErrInvalidHash)
}
