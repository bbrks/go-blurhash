package blurhash

import (
	"image"
	"image/draw"
	"math"
	"strings"

	"github.com/bbrks/go-blurhash/base83"
)

const (
	minComponents = 1
	maxComponents = 9
)

// Encode returns the blurhash for the given image.
func Encode(xComponents, yComponents int, img image.Image) (hash string, err error) {
	if xComponents < minComponents || xComponents > maxComponents ||
		yComponents < minComponents || yComponents > maxComponents {
		return "", ErrInvalidComponents
	}

	b := strings.Builder{}

	sizeFlag := (xComponents - 1) + (yComponents-1)*9
	sizeFlagEncoded, err := base83.Encode(sizeFlag, 1)
	if err != nil {
		return "", err
	}

	_, err = b.WriteString(sizeFlagEncoded)
	if err != nil {
		return "", err
	}

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Convert image to NRGBA for efficient direct pixel access (one-time cost)
	var nrgba *image.NRGBA
	if n, ok := img.(*image.NRGBA); ok {
		nrgba = n
	} else {
		nrgba = image.NewNRGBA(bounds)
		draw.Draw(nrgba, bounds, img, bounds.Min, draw.Src)
	}

	// Precompute cosine tables to avoid repeated math.Cos calls
	cosX := make([][]float64, xComponents)
	for i := 0; i < xComponents; i++ {
		cosX[i] = make([]float64, width)
		for x := 0; x < width; x++ {
			cosX[i][x] = math.Cos(math.Pi * float64(i) * float64(x) / float64(width))
		}
	}
	cosY := make([][]float64, yComponents)
	for j := 0; j < yComponents; j++ {
		cosY[j] = make([]float64, height)
		for y := 0; y < height; y++ {
			cosY[j][y] = math.Cos(math.Pi * float64(j) * float64(y) / float64(height))
		}
	}

	// vector of yComponents*xComponents*(RGB)
	factors := make([][][3]float64, yComponents)
	for y := 0; y < yComponents; y++ {
		factors[y] = make([][3]float64, xComponents)
		for x := 0; x < xComponents; x++ {
			factor := multiplyBasisFunction(x, y, nrgba, cosX[x], cosY[y])
			factors[y][x][0] = factor[0]
			factors[y][x][1] = factor[1]
			factors[y][x][2] = factor[2]
		}
	}

	maximumValue := 0.0
	if xComponents*yComponents-1 > 0 {
		actualMaximumValue := 0.0
		for y := 0; y < yComponents; y++ {
			for x := 0; x < xComponents; x++ {
				if y == 0 && x == 0 {
					continue
				}
				actualMaximumValue = math.Max(math.Abs(factors[y][x][0]), actualMaximumValue)
				actualMaximumValue = math.Max(math.Abs(factors[y][x][1]), actualMaximumValue)
				actualMaximumValue = math.Max(math.Abs(factors[y][x][2]), actualMaximumValue)
			}
		}

		quantisedMaximumValue := math.Max(0, math.Min(82, math.Floor(actualMaximumValue*166-0.5)))
		maximumValue = (quantisedMaximumValue + 1) / 166
		str, err := base83.Encode(int(quantisedMaximumValue), 1)
		if err != nil {
			return "", err
		}
		b.WriteString(str)
	} else {
		maximumValue = 1
		str, err := base83.Encode(0, 1)
		if err != nil {
			return "", err
		}
		b.WriteString(str)
	}

	dc := factors[0][0]
	str, err := base83.Encode(encodeDC(dc[0], dc[1], dc[2]), 4)
	if err != nil {
		return "", err
	}
	b.WriteString(str)

	for y := 0; y < yComponents; y++ {
		for x := 0; x < xComponents; x++ {
			if y == 0 && x == 0 {
				continue
			}
			str, err := base83.Encode(encodeAC(factors[y][x][0], factors[y][x][1], factors[y][x][2], maximumValue), 2)
			if err != nil {
				return "", err
			}
			b.WriteString(str)
		}
	}

	return b.String(), nil
}

func encodeDC(r, g, b float64) int {
	return (linearTosRGB(r) << 16) + (linearTosRGB(g) << 8) + linearTosRGB(b)
}

func encodeAC(r, g, b, maximumValue float64) int {
	quantR := math.Max(0, math.Min(18, math.Floor(signPow(r/maximumValue, 0.5)*9+9.5)))
	quantG := math.Max(0, math.Min(18, math.Floor(signPow(g/maximumValue, 0.5)*9+9.5)))
	quantB := math.Max(0, math.Min(18, math.Floor(signPow(b/maximumValue, 0.5)*9+9.5)))

	return int(quantR*19*19 + quantG*19 + quantB)
}

func multiplyBasisFunction(xComp, yComp int, img *image.NRGBA, cosX, cosY []float64) [3]float64 {
	var r, g, b float64
	width, height := len(cosX), len(cosY)

	normalisation := 2.0
	if xComp == 0 && yComp == 0 {
		normalisation = 1.0
	}

	// Direct pixel access - avoids interface calls and allocations
	pix := img.Pix
	stride := img.Stride

	for y := 0; y < height; y++ {
		rowOffset := y * stride
		basisY := cosY[y]
		for x := 0; x < width; x++ {
			i := rowOffset + x*4
			basis := cosX[x] * basisY
			r += basis * sRGBToLinear(int(pix[i]))
			g += basis * sRGBToLinear(int(pix[i+1]))
			b += basis * sRGBToLinear(int(pix[i+2]))
		}
	}

	scale := normalisation / float64(width*height)
	return [3]float64{
		r * scale,
		g * scale,
		b * scale,
	}
}
