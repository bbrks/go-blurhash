package blurhash

import (
	"fmt"
	"testing"

	"github.com/matryer/is"
)

func TestSRGBLinearRoundtrip(t *testing.T) {
	// Roundtrip sRGB -> Linear -> sRGB to ensure same value
	for srgb := 0; srgb <= 255; srgb++ {
		t.Run(fmt.Sprintf("sRGB %d", srgb), func(t *testing.T) {
			is := is.New(t)
			linear := sRGBToLinear(srgb)
			back := linearTosRGB(linear)
			is.Equal(srgb, back) // expecting sgrb value to roundtrip (srgb -> linear -> srgb)
		})
	}
}
