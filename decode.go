package blurhash

import (
	"image"
	"image/color"
	"math"

	"github.com/bbrks/go-blurhash/base83"
)

// Components returns the X and Y components of a blurhash.
func Components(hash string) (x, y int, err error) {
	sizeFlag, err := base83.Decode(string(hash[0]))
	if err != nil {
		return 0, 0, err
	}

	x = (sizeFlag % 9) + 1
	y = (sizeFlag / 9) + 1

	expectedLength := 4 + 2*x*y
	actualLength := len(hash)
	if expectedLength != actualLength {
		return 0, 0, lengthError(expectedLength, actualLength)
	}

	return x, y, nil
}

// Decode returns an image of the given hash with the given size.
func Decode(hash string, width, height int, punch int) (img image.Image, err error) {
	numX, numY, err := Components(hash)
	if err != nil {
		return nil, err
	}

	fPunch := float64(punch)

	quantisedMaximumValue, err := base83.Decode(string(hash[1]))
	if err != nil {
		return nil, err
	}
	maximumValue := float64(quantisedMaximumValue+1) / 166

	// for each component
	colors := make([][3]float64, numX*numY)
	for i := range colors {
		if i == 0 {
			val, err := base83.Decode(hash[2:6])
			if err != nil {
				return nil, err
			}
			colors[i] = decodeDC(val)
		} else {
			val, err := base83.Decode(hash[4+i*2 : 6+i*2])
			if err != nil {
				return nil, err
			}
			colors[i] = decodeAC(float64(val), maximumValue*fPunch)
		}
	}

	newImg := image.NewNRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var r, g, b float64

			for j := 0; j < numY; j++ {
				for i := 0; i < numX; i++ {
					basis := math.Cos(math.Pi*float64(x)*float64(i)/float64(width)) * math.Cos(math.Pi*float64(y)*float64(j)/float64(height))
					compColor := colors[i+j*numX]
					r += compColor[0] * basis
					g += compColor[1] * basis
					b += compColor[2] * basis
				}
			}

			newImg.SetNRGBA(x, y, color.NRGBA{
				R: uint8(linearTosRGB(r)),
				G: uint8(linearTosRGB(g)),
				B: uint8(linearTosRGB(b)),
				A: 255,
			})
		}
	}

	return newImg, nil
}

func decodeDC(val int) (c [3]float64) {
	c[0] = sRGBToLinear(val >> 16)
	c[1] = sRGBToLinear(val >> 8 & 255)
	c[2] = sRGBToLinear(val & 255)
	return c
}

func decodeAC(val, maximumValue float64) (c [3]float64) {
	quantR := math.Floor(val / (19 * 19))
	quantG := math.Mod(math.Floor(val/19), 19)
	quantB := math.Mod(val, 19)
	c[0] = signPow((quantR-9)/9, 2.0) * maximumValue
	c[1] = signPow((quantG-9)/9, 2.0) * maximumValue
	c[2] = signPow((quantB-9)/9, 2.0) * maximumValue
	return c
}
