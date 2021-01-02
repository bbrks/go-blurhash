package blurhash

import (
	"image"
	"image/color"
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

	// vector of yComponents*xComponents*(RGB)
	factors := make([][][3]float64, yComponents)
	for y := 0; y < yComponents; y++ {
		factors[y] = make([][3]float64, xComponents)
		for x := 0; x < xComponents; x++ {
			factor := multiplyBasisFunction(x, y, img)
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

func multiplyBasisFunction(xComponents, yComponents int, img image.Image) [3]float64 {
	var r, g, b float64
	width, height := float64(img.Bounds().Dx()), float64(img.Bounds().Dy())

	normalisation := 2.0
	if xComponents == 0 && yComponents == 0 {
		normalisation = 1.0
	}

	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Max.Y; y++ {
			//cR, cG, cB, _ := img.At(x, y).RGBA()
			c, ok := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
			if !ok {
				panic("not color.NRGBA")
			}
			basis := math.Cos(math.Pi*float64(xComponents)*float64(x)/width) *
				math.Cos(math.Pi*float64(yComponents)*float64(y)/height)
			r += basis * sRGBToLinear(int(c.R))
			g += basis * sRGBToLinear(int(c.G))
			b += basis * sRGBToLinear(int(c.B))
		}
	}

	scale := normalisation / (width * height)
	return [3]float64{
		r * scale,
		g * scale,
		b * scale,
	}
}
