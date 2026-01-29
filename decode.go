package blurhash

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/bbrks/go-blurhash/base83"
)

// Components returns the X and Y components of a blurhash.
func Components(hash string) (x, y int, err error) {
	if len(hash) < 6 {
		return 0, 0, ErrInvalidHash
	}

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

// Decode returns an NRGBA image of the given hash with the given size.
func Decode(hash string, width, height int, punch int) (image.Image, error) {
	if width <= 0 || height <= 0 {
		return nil, ErrInvalidDimensions
	}
	newImg := image.NewNRGBA(image.Rect(0, 0, width, height))
	if err := DecodeDraw(newImg, hash, float64(punch)); err != nil {
		return nil, err
	}
	return newImg, nil
}

// DecodeDraw decodes the given hash into the given image.
func DecodeDraw(dst draw.Image, hash string, punch float64) error {
	numX, numY, err := Components(hash)
	if err != nil {
		return err
	}

	quantisedMaximumValue, err := base83.Decode(string(hash[1]))
	if err != nil {
		return err
	}
	maximumValue := float64(quantisedMaximumValue+1) / 166

	// for each component
	colors := make([][3]float64, numX*numY)
	for i := range colors {
		if i == 0 {
			val, err := base83.Decode(hash[2:6])
			if err != nil {
				return err
			}
			colors[i] = decodeDC(val)
		} else {
			val, err := base83.Decode(hash[4+i*2 : 6+i*2])
			if err != nil {
				return err
			}
			colors[i] = decodeAC(float64(val), maximumValue*punch)
		}
	}

	bounds := dst.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Precompute cosine tables
	cosX := make([][]float64, numX)
	for i := 0; i < numX; i++ {
		cosX[i] = make([]float64, width)
		for x := 0; x < width; x++ {
			cosX[i][x] = math.Cos(math.Pi * float64(i) * float64(x) / float64(width))
		}
	}
	cosY := make([][]float64, numY)
	for j := 0; j < numY; j++ {
		cosY[j] = make([]float64, height)
		for y := 0; y < height; y++ {
			cosY[j][y] = math.Cos(math.Pi * float64(j) * float64(y) / float64(height))
		}
	}

	// Get direct pixel access if available (NRGBA and RGBA have same layout)
	var pix []uint8
	var stride int
	switch img := dst.(type) {
	case *image.NRGBA:
		pix = img.Pix
		stride = img.Stride
	case *image.RGBA:
		pix = img.Pix
		stride = img.Stride
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var r, g, b float64
			for j := 0; j < numY; j++ {
				basisY := cosY[j][y]
				for i := 0; i < numX; i++ {
					basis := cosX[i][x] * basisY
					compColor := colors[i+j*numX]
					r += compColor[0] * basis
					g += compColor[1] * basis
					b += compColor[2] * basis
				}
			}

			if pix != nil {
				idx := y*stride + x*4
				pix[idx] = uint8(linearTosRGB(r))
				pix[idx+1] = uint8(linearTosRGB(g))
				pix[idx+2] = uint8(linearTosRGB(b))
				pix[idx+3] = 255
			} else {
				dst.Set(x, y, color.NRGBA{
					uint8(linearTosRGB(r)),
					uint8(linearTosRGB(g)),
					uint8(linearTosRGB(b)),
					255,
				})
			}
		}
	}

	return nil
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
