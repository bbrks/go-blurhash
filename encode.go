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

// Encoder is a reusable blurhash encoder that minimizes allocations
// by reusing internal buffers across encode operations.
//
// An Encoder is safe for sequential use but not for concurrent use.
// For concurrent workloads, use a sync.Pool of Encoders.
//
// The zero value is ready to use.
type Encoder struct {
	cosX, cosY []float64
	factors    [][3]float64
	nrgba      *image.NRGBA
	builder    strings.Builder
}

// NewEncoder creates a new reusable Encoder.
// Buffers are allocated lazily on first use and grown as needed.
func NewEncoder() *Encoder {
	return &Encoder{}
}

// Encode returns the blurhash for the given image.
// Internal buffers are reused across calls when possible.
func (e *Encoder) Encode(xComponents, yComponents int, img image.Image) (string, error) {
	if xComponents < minComponents || xComponents > maxComponents ||
		yComponents < minComponents || yComponents > maxComponents {
		return "", ErrInvalidComponents
	}

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Ensure buffers are large enough
	e.maybeGrowBuffers(width, height, xComponents, yComponents)

	// Reset builder for new encode
	e.builder.Reset()

	sizeFlag := (xComponents - 1) + (yComponents-1)*9
	sizeFlagEncoded, err := base83.Encode(sizeFlag, 1)
	if err != nil {
		return "", err
	}
	e.builder.WriteString(sizeFlagEncoded)

	// Get direct 4-byte per pixel data - fast path [N]RGBA
	var pix []uint8
	var stride int
	switch src := img.(type) {
	case *image.NRGBA:
		pix = src.Pix
		stride = src.Stride
	case *image.RGBA:
		pix = src.Pix
		stride = src.Stride
	default:
		// Reuse NRGBA buffer if large enough
		if e.nrgba == nil || e.nrgba.Bounds().Dx() < width || e.nrgba.Bounds().Dy() < height {
			e.nrgba = image.NewNRGBA(bounds)
		} else {
			e.nrgba.Rect = bounds
		}
		draw.Draw(e.nrgba, bounds, img, bounds.Min, draw.Src)
		pix = e.nrgba.Pix
		stride = e.nrgba.Stride
	}

	// Compute cosine tables into reusable buffers
	for i := 0; i < xComponents; i++ {
		for x := 0; x < width; x++ {
			e.cosX[i*width+x] = math.Cos(math.Pi * float64(i) * float64(x) / float64(width))
		}
	}
	for j := 0; j < yComponents; j++ {
		for y := 0; y < height; y++ {
			e.cosY[j*height+y] = math.Cos(math.Pi * float64(j) * float64(y) / float64(height))
		}
	}

	// Compute DCT factors
	for j := 0; j < yComponents; j++ {
		for i := 0; i < xComponents; i++ {
			cosXSlice := e.cosX[i*width : i*width+width]
			cosYSlice := e.cosY[j*height : j*height+height]
			factor := multiplyBasisFunction(i, j, pix, stride, cosXSlice, cosYSlice)
			e.factors[j*xComponents+i] = factor
		}
	}

	maximumValue := 0.0
	if xComponents*yComponents-1 > 0 {
		actualMaximumValue := 0.0
		for j := 0; j < yComponents; j++ {
			for i := 0; i < xComponents; i++ {
				if j == 0 && i == 0 {
					continue
				}
				f := e.factors[j*xComponents+i]
				actualMaximumValue = math.Max(math.Abs(f[0]), actualMaximumValue)
				actualMaximumValue = math.Max(math.Abs(f[1]), actualMaximumValue)
				actualMaximumValue = math.Max(math.Abs(f[2]), actualMaximumValue)
			}
		}

		quantisedMaximumValue := math.Max(0, math.Min(82, math.Floor(actualMaximumValue*166-0.5)))
		maximumValue = (quantisedMaximumValue + 1) / 166
		str, err := base83.Encode(int(quantisedMaximumValue), 1)
		if err != nil {
			return "", err
		}
		e.builder.WriteString(str)
	} else {
		maximumValue = 1
		str, err := base83.Encode(0, 1)
		if err != nil {
			return "", err
		}
		e.builder.WriteString(str)
	}

	dc := e.factors[0]
	str, err := base83.Encode(encodeDC(dc[0], dc[1], dc[2]), 4)
	if err != nil {
		return "", err
	}
	e.builder.WriteString(str)

	for j := 0; j < yComponents; j++ {
		for i := 0; i < xComponents; i++ {
			if j == 0 && i == 0 {
				continue
			}
			f := e.factors[j*xComponents+i]
			str, err := base83.Encode(encodeAC(f[0], f[1], f[2], maximumValue), 2)
			if err != nil {
				return "", err
			}
			e.builder.WriteString(str)
		}
	}

	return e.builder.String(), nil
}

func (e *Encoder) maybeGrowBuffers(width, height, xComponents, yComponents int) {
	e.cosX = growTo(e.cosX, xComponents*width)
	e.cosY = growTo(e.cosY, yComponents*height)
	e.factors = growTo(e.factors, xComponents*yComponents)
}

// Encode returns the blurhash for the given image.
func Encode(xComponents, yComponents int, img image.Image) (string, error) {
	var e Encoder
	return e.Encode(xComponents, yComponents, img)
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

func multiplyBasisFunction(xComp, yComp int, pix []uint8, stride int, cosX, cosY []float64) [3]float64 {
	var r, g, b float64
	width, height := len(cosX), len(cosY)

	normalisation := 2.0
	if xComp == 0 && yComp == 0 {
		normalisation = 1.0
	}

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
