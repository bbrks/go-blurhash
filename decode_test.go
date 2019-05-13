package blurhash_test

import (
	"image/png"
	"io/ioutil"
	"testing"

	"github.com/matryer/is"

	"github.com/bbrks/go-blurhash"
)

func TestDecode(t *testing.T) {
	for _, test := range testFixtures {
		// skip tests without hashes
		if test.hash == "" {
			continue
		}

		t.Run(test.hash, func(t *testing.T) {
			is := is.New(t)

			img, err := blurhash.Decode(test.hash, 32, 32, 1)
			is.NoErr(err)

			err = png.Encode(ioutil.Discard, img)
			is.NoErr(err)
		})
	}
}

func TestComponents(t *testing.T) {
	for _, test := range testFixtures {
		// skip tests without expected component values
		if test.hash == "" || test.xComp == 0 || test.yComp == 0 {
			continue
		}

		t.Run(test.hash, func(t *testing.T) {
			is := is.NewRelaxed(t)

			x, y, err := blurhash.Components(test.hash)
			is.NoErr(err)           // unexpected error getting components
			is.Equal(x, test.xComp) // blurhash component mismatch
			is.Equal(y, test.yComp) // blurhash component mismatch
		})
	}
}

func BenchmarkComponents(b *testing.B) {
	for _, test := range testFixtures {
		// skip tests without hashes
		if test.hash == "" {
			continue
		}

		b.Run(test.hash, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _, _ = blurhash.Components(test.hash)
			}
		})
	}
}

func BenchmarkDecode(b *testing.B) {
	for _, test := range testFixtures {
		// skip tests without hashes
		if test.hash == "" {
			continue
		}

		b.Run(test.hash, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = blurhash.Decode(test.hash, 32, 32, 1)
			}
		})
	}
}
