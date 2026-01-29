package blurhash

import "math"

// growTo returns a slice with length size, reusing the existing
// slice's backing array if it has sufficient capacity.
func growTo[T any](s []T, size int) []T {
	if cap(s) < size {
		return make([]T, size)
	}
	return s[:size]
}

// signSqrt returns sign(val) * sqrt(|val|)
func signSqrt(val float64) float64 {
	return math.Copysign(math.Sqrt(math.Abs(val)), val)
}

func sRGBToLinear(val int) float64 {
	return sRGBToLinearLUT[val]
}

// linearToSRGBLUTSize is the number of entries in the linearToSRGB lookup table.
// 4096 entries provides sufficient precision for 8-bit output.
const linearToSRGBLUTSize = 4096

// linearToSRGBLUT maps linear values [0,1] to sRGB [0,255].
// Index i corresponds to linear value i/(LUTSize-1).
var linearToSRGBLUT [linearToSRGBLUTSize]uint8

func init() {
	for i := 0; i < linearToSRGBLUTSize; i++ {
		v := float64(i) / float64(linearToSRGBLUTSize-1)
		var srgb float64
		if v <= 0.0031308 {
			srgb = v * 12.92
		} else {
			srgb = 1.055*math.Pow(v, 1/2.4) - 0.055
		}
		linearToSRGBLUT[i] = uint8(srgb*255 + 0.5)
	}
}

func linearToSRGB(val float64) int {
	if val <= 0 {
		return 0
	}
	if val >= 1 {
		return 255
	}
	return int(linearToSRGBLUT[int(val*float64(linearToSRGBLUTSize-1)+0.5)])
}
