package blurhash_test

import (
	"errors"
	"image"
	"image/png"
	"io"
	"testing"

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
			img := image.NewRGBA(image.Rect(0, 0, 32, 32))

			err := blurhash.DecodeDraw(img, test.hash, 1)
			if err != nil {
				t.Fatalf("error decoding: %v", err)
			}

			err = png.Encode(io.Discard, img)
			if err != nil {
				t.Fatalf("error encoding png: %v", err)
			}
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
			img, err := blurhash.Decode(test.hash, 32, 32, 1)
			if err != nil {
				t.Fatalf("error decoding: %v", err)
			}

			err = png.Encode(io.Discard, img)
			if err != nil {
				t.Fatalf("error encoding png: %v", err)
			}
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
			x, y, err := blurhash.Components(test.hash)
			if err != nil {
				t.Fatalf("unexpected error getting components: %v", err)
			}
			if x != test.xComp {
				t.Errorf("x component mismatch: got %d, want %d", x, test.xComp)
			}
			if y != test.yComp {
				t.Errorf("y component mismatch: got %d, want %d", y, test.yComp)
			}
		})
	}
}

func TestComponentsInvalidHash(t *testing.T) {
	t.Run("too short", func(t *testing.T) {
		// Hashes shorter than 6 characters should return ErrInvalidHash
		shortHashes := []string{"", "A", "ABCDE"}
		for _, hash := range shortHashes {
			_, _, err := blurhash.Components(hash)
			if !errors.Is(err, blurhash.ErrInvalidHash) {
				t.Errorf("short hash %q should return ErrInvalidHash, got %v", hash, err)
			}
		}
	})

	t.Run("invalid base83 character", func(t *testing.T) {
		// '&' is not a valid base83 character
		_, _, err := blurhash.Components("&BCDEF")
		if err == nil {
			t.Fatal("invalid character should return error")
		}
		if !errors.Is(err, base83.ErrInvalidInput) {
			t.Errorf("expected invalid base83 error, got %v", err)
		}
	})

	t.Run("wrong length for components", func(t *testing.T) {
		// '9' encodes 1x2 components (sizeFlag=9), expecting 4+2*1*2=8 chars
		// but we provide only 6 chars
		_, _, err := blurhash.Components("900000")
		if !errors.Is(err, blurhash.ErrInvalidHash) {
			t.Errorf("wrong length should return ErrInvalidHash, got %v", err)
		}

		// Valid 1x1 hash is 6 chars, but provide 8
		_, _, err = blurhash.Components("00000000")
		if !errors.Is(err, blurhash.ErrInvalidHash) {
			t.Errorf("wrong length should return ErrInvalidHash, got %v", err)
		}
	})
}

func TestDecodeDrawSubImage(t *testing.T) {
	// Create a larger image and get a sub-image from it
	parent := image.NewNRGBA(image.Rect(0, 0, 100, 100))
	subRect := image.Rect(10, 20, 42, 52) // 32x32 sub-image at offset (10, 20)
	subImg := parent.SubImage(subRect).(*image.NRGBA)

	// Decode into the sub-image
	err := blurhash.DecodeDraw(subImg, testFixtures[0].hash, 1)
	if err != nil {
		t.Fatalf("error decoding: %v", err)
	}

	// Verify pixels were written to the correct location
	// The sub-image should have non-zero pixels
	hasNonZero := false
	for y := subRect.Min.Y; y < subRect.Max.Y; y++ {
		for x := subRect.Min.X; x < subRect.Max.X; x++ {
			c := parent.NRGBAAt(x, y)
			if c.R != 0 || c.G != 0 || c.B != 0 {
				hasNonZero = true
				break
			}
		}
	}
	if !hasNonZero {
		t.Error("sub-image should have non-zero pixels")
	}

	// Verify pixels outside the sub-image are still zero
	// Check a few pixels outside the sub-rect
	outside := parent.NRGBAAt(0, 0)
	if outside.R != 0 {
		t.Errorf("pixel outside sub-image R should be zero, got %d", outside.R)
	}
	if outside.G != 0 {
		t.Errorf("pixel outside sub-image G should be zero, got %d", outside.G)
	}
	if outside.B != 0 {
		t.Errorf("pixel outside sub-image B should be zero, got %d", outside.B)
	}
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
			_, err := blurhash.Decode(testFixtures[0].hash, tt.width, tt.height, 1)
			if !errors.Is(err, blurhash.ErrInvalidDimensions) {
				t.Errorf("invalid dimensions should return ErrInvalidDimensions, got %v", err)
			}
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
			dst := image.NewRGBA(image.Rect(0, 0, 32, 32))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = blurhash.DecodeDraw(dst, test.hash, 1)
			}
		})
	}
}

func BenchmarkDecoderReuse(b *testing.B) {
	for _, test := range testFixtures {
		if test.hash == "" {
			continue
		}

		b.Run(test.hash, func(b *testing.B) {
			dec := blurhash.NewDecoder()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = dec.Decode(test.hash, 32, 32, 1)
			}
		})
	}
}

func BenchmarkDecoderDrawReuse(b *testing.B) {
	for _, test := range testFixtures {
		if test.hash == "" {
			continue
		}

		b.Run(test.hash, func(b *testing.B) {
			dec := blurhash.NewDecoder()
			dst := image.NewRGBA(image.Rect(0, 0, 32, 32))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = dec.DecodeDraw(dst, test.hash, 1)
			}
		})
	}
}
