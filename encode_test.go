package blurhash_test

import (
	"image"
	"image/draw"
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

func TestEncodeSubImage(t *testing.T) {
	is := is.New(t)

	// Load a test image
	f, err := os.Open(filepath.FromSlash("fixtures/test.png"))
	is.NoErr(err)
	defer f.Close()

	img, _, err := image.Decode(f)
	is.NoErr(err)

	// Create a sub-image with non-zero Min bounds
	bounds := img.Bounds()
	subRect := image.Rect(
		bounds.Min.X+10, bounds.Min.Y+10,
		bounds.Max.X-10, bounds.Max.Y-10,
	)

	type subImager interface {
		SubImage(r image.Rectangle) image.Image
	}
	subImg := img.(subImager).SubImage(subRect)

	// Verify sub-image has non-zero Min (this is the bug trigger)
	is.True(subImg.Bounds().Min.X != 0 || subImg.Bounds().Min.Y != 0)

	// Encode the sub-image
	subHash, err := blurhash.Encode(4, 3, subImg)
	is.NoErr(err)
	t.Logf("sub-image hash: %s", subHash)

	// Create a copy of the sub-image with (0,0) origin
	normalImg := image.NewNRGBA(image.Rect(0, 0, subRect.Dx(), subRect.Dy()))
	draw.Draw(normalImg, normalImg.Bounds(), subImg, subRect.Min, draw.Src)

	// Encode the (0,0)-origin copy
	normalHash, err := blurhash.Encode(4, 3, normalImg)
	is.NoErr(err)

	// Both should produce the same hash
	is.Equal(subHash, normalHash) // sub-image should encode same as equivalent normal image
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

func BenchmarkEncoderReuse(b *testing.B) {
	for _, test := range testFixtures {
		// skip tests without files or hashes
		if test.file == "" || test.hash == "" {
			continue
		}

		b.Run(test.hash, func(b *testing.B) {
			f, err := os.Open(filepath.FromSlash(test.file))
			if err != nil {
				b.Fatalf("error opening file: %v", err)
			}
			defer f.Close()

			img, _, err := image.Decode(f)
			if err != nil {
				b.Fatalf("error decoding image: %v", err)
			}

			enc := blurhash.NewEncoder()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = enc.Encode(test.xComp, test.yComp, img)
			}
		})
	}
}
