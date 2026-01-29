package blurhash

import "math"

func signPow(val, exp float64) float64 {
	sign := 1.0
	if val < 0 {
		sign = -1
	}
	return sign * math.Pow(math.Abs(val), exp)
}

//go:generate go run srgb_lut_generator.go
func sRGBToLinear(val int) float64 {
	return sRGBToLinearLUT[val]
}

// linearTosRGBLUTSize is the number of entries in the linearTosRGB lookup table.
// 4096 entries provides sufficient precision for 8-bit output.
const linearTosRGBLUTSize = 4096

// linearTosRGBLUT maps linear values [0,1] to sRGB [0,255].
// Index i corresponds to linear value i/(LUTSize-1).
var linearTosRGBLUT [linearTosRGBLUTSize]uint8

func init() {
	for i := 0; i < linearTosRGBLUTSize; i++ {
		v := float64(i) / float64(linearTosRGBLUTSize-1)
		var srgb float64
		if v <= 0.0031308 {
			srgb = v * 12.92
		} else {
			srgb = 1.055*math.Pow(v, 1/2.4) - 0.055
		}
		linearTosRGBLUT[i] = uint8(srgb*255 + 0.5)
	}
}

func linearTosRGB(val float64) int {
	if val <= 0 {
		return 0
	}
	if val >= 1 {
		return 255
	}
	return int(linearTosRGBLUT[int(val*float64(linearTosRGBLUTSize-1)+0.5)])
}
