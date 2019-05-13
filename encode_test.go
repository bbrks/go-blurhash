package blurhash_test

import (
	"image"
	"os"
	"path/filepath"
	"testing"

	"github.com/matryer/is"

	"github.com/bbrks/go-blurhash"
)

func TestEncode(t *testing.T) {
	for _, test := range testFixtures {
		if test.file == "" {
			// skip tests without files
			continue
		}

		t.Run(test.hash, func(t *testing.T) {
			is := is.New(t)

			f, err := os.Open(filepath.FromSlash(test.file))
			is.NoErr(err) // error opening test fixture file
			defer f.Close()

			is.True(f != nil) // file should not be nil

			img, _, err := image.Decode(f)
			is.NoErr(err)       // error decoding image from test fixture
			is.True(img != nil) // image should not be nil

			hash, err := blurhash.Encode(test.xComp, test.yComp, img)
			is.NoErr(err)             // error hashing test fixture image
			is.Equal(hash, test.hash) // blurhash mismatch
		})
	}
}

func BenchmarkEncode(b *testing.B) {
	for _, test := range testFixtures {
		if test.file == "" {
			// skip tests without files
			continue
		}

		b.Run(test.hash, func(b *testing.B) {
			is := is.New(b)

			f, err := os.Open(filepath.FromSlash(test.file))
			is.NoErr(err) // error opening test fixture file
			defer f.Close()

			is.True(f != nil) // file should not be nil

			img, _, err := image.Decode(f)
			is.NoErr(err)       // error decoding image from test fixture
			is.True(img != nil) // image should not be nil

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = blurhash.Encode(test.xComp, test.yComp, img)
			}
		})
	}
}
