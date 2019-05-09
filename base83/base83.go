package base83

import (
	"strings"
)

const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz#$%*+,-.:;=?@[]^_{|}~"

// Decode decodes a base83 string into an integer value.
func Decode(str string) (val int, err error) {
	for i, r := range str {
		idx := strings.IndexRune(chars, r)
		if idx == -1 {
			return 0, invalidError(r, i)
		}

		val = val*len(chars) + idx
	}
	return val, nil
}

// Encode encodes an integer value into a base83 string of the given length.
func Encode(val, length int) (str string, err error) {

	divisor := 1
	for i := 0; i < length-1; i++ {
		divisor *= len(chars)
	}

	for i := 0; i < length; i++ {
		idx := val / divisor % len(chars)
		divisor /= len(chars)
		str += string(chars[idx])
	}

	return str, nil
}
