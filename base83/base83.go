package base83

const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz#$%*+,-.:;=?@[]^_{|}~"

// lookup table for ASCII->base83 mapping
var charLookup = [256]int{
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, 62, 63, 64, -1, -1, -1, -1, 65, 66, 67, 68, 69, -1,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 70, 71, -1, 72, -1, 73,
	74, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24,
	25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 75, -1, 76, 77, 78,
	-1, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
	51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 79, 80, 81, 82, -1,
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
}

// Decode decodes a base83 string into an integer value.
func Decode(str string) (val int, err error) {
	for i := 0; i < len(str); i++ {
		idx := charLookup[str[i]]
		if idx == -1 {
			return 0, invalidError(rune(str[i]), i)
		}
		val = val*83 + idx
	}
	return val, nil
}

// Encode encodes an integer value into a base83 string of the given length.
func Encode(val, length int) (str string, err error) {
	buf := make([]byte, length)

	for i := length - 1; i >= 0; i-- {
		buf[i] = chars[val%83]
		val /= 83
	}

	return string(buf), nil
}
