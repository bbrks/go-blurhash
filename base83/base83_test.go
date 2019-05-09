package base83_test

import (
	"strings"
	"testing"

	"github.com/matryer/is"

	"github.com/bbrks/go-blurhash/base83"
)

var tests = []struct {
	str string
	val int
}{
	{"3", 3},
	{"A", 10},
	{":", 70},
	{"~", 82},
	{"01", 1}, // leading zeros are "trimmed"
	{"11", 84},
	{"33", 252},
	{"~$", 6869},
	{"%%%%%%", 255172974336},
}

func TestDecodeEncode(t *testing.T) {
	for _, test := range tests {
		t.Run(test.str, func(t *testing.T) {
			is := is.NewRelaxed(t)

			val, err := base83.Decode(test.str)
			is.NoErr(err)           // Decode returned unexpected error
			is.Equal(val, test.val) // Decode got unexpected result
		})
	}
}

func TestEncode(t *testing.T) {
	for _, test := range tests {
		t.Run(test.str, func(t *testing.T) {
			is := is.NewRelaxed(t)

			str, err := base83.Encode(test.val, len(test.str))
			is.NoErr(err)           // Encode returned unexpected error
			is.Equal(str, test.str) // Encode got unexpected result
		})
	}
}

func TestDecodeInvalidInput(t *testing.T) {
	tests := []struct {
		str string
		val int
		err error
	}{
		{"&", 0, base83.ErrInvalidInput},
	}

	for _, test := range tests {
		t.Run(test.str, func(t *testing.T) {
			is := is.NewRelaxed(t)

			val, err := base83.Decode(test.str)
			is.True(err != nil)                                      // Decode should've returned error for invalid input
			is.True(strings.Contains(err.Error(), test.err.Error())) // Decode returned wrong error
			is.Equal(val, test.val)                                  // Decode got unexpected result
		})
	}
}

func TestEncodeInvalidLength(t *testing.T) {
	tests := []struct {
		val    int
		length int
		str    string
	}{
		{255172974336, 3, "%%%"},
		{255172974336, 6, "%%%%%%"},
		{255172974336, 9, "000%%%%%%"},
	}

	for _, test := range tests {
		t.Run(test.str, func(t *testing.T) {
			is := is.NewRelaxed(t)

			output, err := base83.Encode(test.val, test.length)
			is.NoErr(err)              // Encode should've returned error for invalid input
			is.Equal(output, test.str) // Encode got unexpected result
		})
	}
}

func BenchmarkDecode(b *testing.B) {
	for _, test := range tests {
		b.Run(test.str, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, _ = base83.Decode("~$")
			}
		})
	}
}

func BenchmarkEncode(b *testing.B) {
	for _, test := range tests {
		b.Run(test.str, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = base83.Encode(6869, 2)
			}
		})
	}
}
