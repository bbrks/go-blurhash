package blurhash_test

import (
	"image"
	"image/png"
	"io/ioutil"
	"testing"

	"github.com/matryer/is"

	"github.com/bbrks/go-blurhash"
	"github.com/bbrks/go-blurhash/base83"
)

func TestDecodeRGBA(t *testing.T) {
	for _, test := range testFixtures {
		// skip tests without hashes
		if test.hash == "" {
			continue
		}

		t.Run(test.hash, func(t *testing.T) {
			is := is.New(t)

			img := image.NewRGBA(image.Rect(0, 0, 32, 32))

			err := blurhash.DecodeDraw(img, test.hash, 1)
			is.NoErr(err)

			err = png.Encode(ioutil.Discard, img)
			is.NoErr(err)
		})
	}
}

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

func TestComponentsInvalidHash(t *testing.T) {
	t.Run("too short", func(t *testing.T) {
		// Hashes shorter than 6 characters should return ErrInvalidHash
		shortHashes := []string{"", "A", "ABCDE"}
		for _, hash := range shortHashes {
			is := is.New(t)
			_, _, err := blurhash.Components(hash)
			is.Equal(err, blurhash.ErrInvalidHash) // short hash should return ErrInvalidHash
		}
	})

	t.Run("invalid base83 character", func(t *testing.T) {
		is := is.New(t)
		// '&' is not a valid base83 character
		_, _, err := blurhash.Components("&BCDEF")
		is.True(err != nil)                   // invalid character should return error
		is.Equal(err, base83.ErrInvalidInput) // expected invalid base83 error
	})

	t.Run("wrong length for components", func(t *testing.T) {
		is := is.New(t)
		// '9' encodes 1x2 components (sizeFlag=9), expecting 4+2*1*2=8 chars
		// but we provide only 6 chars
		_, _, err := blurhash.Components("900000")
		is.Equal(err, blurhash.ErrInvalidHash) // wrong length should return ErrInvalidHash

		// Valid 1x1 hash is 6 chars, but provide 8
		_, _, err = blurhash.Components("00000000")
		is.Equal(err, blurhash.ErrInvalidHash) // wrong length should return ErrInvalidHash
	})
}

func TestDecodeInvalidDimensions(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
	}{
		{"zero width", 0, 32},
		{"zero height", 32, 0},
		{"both zero", 0, 0},
		{"negative width", -1, 32},
		{"negative height", 32, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			_, err := blurhash.Decode(testFixtures[0].hash, tt.width, tt.height, 1)
			is.Equal(err, blurhash.ErrInvalidDimensions) // invalid dimensions should return error
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

func BenchmarkDecodeDraw(b *testing.B) {
	for _, test := range testFixtures {
		// skip tests without hashes
		if test.hash == "" {
			continue
		}

		b.Run(test.hash, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				dst := image.NewRGBA(image.Rect(0, 0, 32, 32))
				_ = blurhash.DecodeDraw(dst, test.hash, 1)
			}
		})
	}
}
