// +build go1.13

package base83

import "fmt"

func invalidError(c rune, i int) error {
	return fmt.Errorf("illegal rune %v at index %d: %w", c, i, ErrInvalidInput)
}
