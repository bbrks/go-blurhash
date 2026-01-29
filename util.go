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

func linearTosRGB(val float64) int {
	v := math.Max(0, math.Min(1, val))
	if v <= 0.0031308 {
		return int(v*12.92*255 + 0.5)
	}
	return int((1.055*math.Pow(v, 1/2.4)-0.055)*255 + 0.5)
}
