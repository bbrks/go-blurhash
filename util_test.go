package blurhash

import (
	"fmt"
	"testing"
)

func TestSRGBLinearRoundtrip(t *testing.T) {
	// Roundtrip sRGB -> Linear -> sRGB to ensure same value
	for srgb := 0; srgb <= 255; srgb++ {
		t.Run(fmt.Sprintf("sRGB %d", srgb), func(t *testing.T) {
			linear := sRGBToLinear(srgb)
			back := linearToSRGB(linear)
			if srgb != back {
				t.Errorf("expecting sRGB value to roundtrip (srgb -> linear -> srgb): got %d, want %d", back, srgb)
			}
		})
	}
}
