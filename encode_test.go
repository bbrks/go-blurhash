package blurhash_test

import (
	"image"
	"image/draw"
	"os"
	"path/filepath"
	"testing"

	"github.com/bbrks/go-blurhash"
)

func TestEncode(t *testing.T) {
	for _, test := range testFixtures {
		if test.file == "" {
			// skip tests without files
			continue
		}

		t.Run(test.hash, func(t *testing.T) {
			f, err := os.Open(filepath.FromSlash(test.file))
			if err != nil {
				t.Fatalf("error opening test fixture file: %v", err)
			}
			defer f.Close() //nolint:errcheck

			if f == nil {
				t.Fatal("file should not be nil")
			}

			img, _, err := image.Decode(f)
			if err != nil {
				t.Fatalf("error decoding image from test fixture: %v", err)
			}
			if img == nil {
				t.Fatal("image should not be nil")
			}

			hash, err := blurhash.Encode(test.xComp, test.yComp, img)
			if err != nil {
				t.Fatalf("error hashing test fixture image: %v", err)
			}
			if hash != test.hash {
				t.Errorf("blurhash mismatch: got %q, want %q", hash, test.hash)
			}
		})
	}
}

func TestEncodeSubImage(t *testing.T) {
	// Load a test image
	f, err := os.Open(filepath.FromSlash("fixtures/test.png"))
	if err != nil {
		t.Fatalf("error opening file: %v", err)
	}
	defer f.Close() //nolint:errcheck

	img, _, err := image.Decode(f)
	if err != nil {
		t.Fatalf("error decoding image: %v", err)
	}

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
	if subImg.Bounds().Min.X == 0 && subImg.Bounds().Min.Y == 0 {
		t.Fatal("sub-image should have non-zero Min bounds")
	}

	// Encode the sub-image
	subHash, err := blurhash.Encode(4, 3, subImg)
	if err != nil {
		t.Fatalf("error encoding sub-image: %v", err)
	}
	t.Logf("sub-image hash: %s", subHash)

	// Create a copy of the sub-image with (0,0) origin
	normalImg := image.NewNRGBA(image.Rect(0, 0, subRect.Dx(), subRect.Dy()))
	draw.Draw(normalImg, normalImg.Bounds(), subImg, subRect.Min, draw.Src)

	// Encode the (0,0)-origin copy
	normalHash, err := blurhash.Encode(4, 3, normalImg)
	if err != nil {
		t.Fatalf("error encoding normal image: %v", err)
	}

	// Both should produce the same hash
	if subHash != normalHash {
		t.Errorf("sub-image should encode same as equivalent normal image: got %q, want %q", subHash, normalHash)
	}
}

func TestEncoderReuse(t *testing.T) {
	// Use a single encoder for all images to verify buffer reuse
	enc := blurhash.NewEncoder()

	// Run multiple iterations to catch any buffer corruption issues
	for iter := 0; iter < 3; iter++ {
		for _, test := range testFixtures {
			if test.file == "" {
				continue
			}

			t.Run(test.hash, func(t *testing.T) {
				f, err := os.Open(filepath.FromSlash(test.file))
				if err != nil {
					t.Fatalf("error opening file: %v", err)
				}
				defer f.Close() //nolint:errcheck

				img, _, err := image.Decode(f)
				if err != nil {
					t.Fatalf("error decoding image: %v", err)
				}

				hash, err := enc.Encode(test.xComp, test.yComp, img)
				if err != nil {
					t.Fatalf("encode error: %v", err)
				}
				if hash != test.hash {
					t.Errorf("hash mismatch: got %q, want %q", hash, test.hash)
				}
			})
		}
	}
}

func BenchmarkEncode(b *testing.B) {
	for _, test := range testFixtures {
		if test.file == "" {
			// skip tests without files
			continue
		}

		b.Run(test.hash, func(b *testing.B) {
			f, err := os.Open(filepath.FromSlash(test.file))
			if err != nil {
				b.Fatalf("error opening test fixture file: %v", err)
			}
			defer f.Close() //nolint:errcheck

			if f == nil {
				b.Fatal("file should not be nil")
			}

			img, _, err := image.Decode(f)
			if err != nil {
				b.Fatalf("error decoding image from test fixture: %v", err)
			}
			if img == nil {
				b.Fatal("image should not be nil")
			}

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
			defer f.Close() //nolint:errcheck

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
