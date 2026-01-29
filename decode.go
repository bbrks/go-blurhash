package blurhash

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/bbrks/go-blurhash/base83"
)

// Components returns the X and Y components of a blurhash.
func Components(hash string) (x, y int, err error) {
	if len(hash) < 6 {
		return 0, 0, fmt.Errorf("%w: hash too short", ErrInvalidHash)
	}

	sizeFlag, err := base83.Decode(string(hash[0]))
	if err != nil {
		return 0, 0, err
	}

	x = (sizeFlag % 9) + 1
	y = (sizeFlag / 9) + 1

	if x < minComponents || x > maxComponents || y < minComponents || y > maxComponents {
		return 0, 0, fmt.Errorf("%w: invalid components x=%d, y=%d", ErrInvalidHash, x, y)
	}

	expectedLength := 4 + 2*x*y
	actualLength := len(hash)
	if expectedLength != actualLength {
		return 0, 0, fmt.Errorf("%w: length mismatch: expected %d, got %d", ErrInvalidHash, expectedLength, actualLength)
	}

	return x, y, nil
}

// Decoder is a reusable blurhash decoder that minimizes allocations
// by reusing internal buffers across decode operations.
//
// A Decoder is safe for sequential use but not for concurrent use.
// For concurrent workloads, use a sync.Pool of Decoders.
//
// The zero value is ready to use.
type Decoder struct {
	cosX, cosY []float64
	colors     [][3]float64
}

// NewDecoder creates a new reusable Decoder.
// Buffers are allocated lazily on first use and grown as needed.
func NewDecoder() *Decoder {
	return &Decoder{}
}

// Decode decodes a blurhash to a new NRGBA image.
// Internal buffers are reused across calls when possible.
func (d *Decoder) Decode(hash string, width, height, punch int) (image.Image, error) {
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("%w: had width=%d, height=%d", ErrInvalidDimensions, width, height)
	}
	newImg := image.NewNRGBA(image.Rect(0, 0, width, height))
	if err := d.DecodeDraw(newImg, hash, float64(punch)); err != nil {
		return nil, err
	}
	return newImg, nil
}

// DecodeDraw decodes a blurhash into an existing image.
// Internal buffers are reused across calls when possible.
func (d *Decoder) DecodeDraw(dst draw.Image, hash string, punch float64) error {
	numX, numY, err := Components(hash)
	if err != nil {
		return err
	}

	quantisedMaximumValue, err := base83.Decode(string(hash[1]))
	if err != nil {
		return err
	}
	maximumValue := float64(quantisedMaximumValue+1) / 166

	bounds := dst.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Ensure buffers are large enough
	d.maybeGrowBuffers(width, height, numX, numY)

	// Decode colors into reusable buffer
	numColors := numX * numY
	for i := 0; i < numColors; i++ {
		if i == 0 {
			val, err := base83.Decode(hash[2:6])
			if err != nil {
				return err
			}
			d.colors[i] = decodeDC(val)
		} else {
			val, err := base83.Decode(hash[4+i*2 : 6+i*2])
			if err != nil {
				return err
			}
			d.colors[i] = decodeAC(val, maximumValue*punch)
		}
	}

	// Compute cosine tables into reusable buffers
	for i := 0; i < numX; i++ {
		for x := 0; x < width; x++ {
			d.cosX[i*width+x] = math.Cos(math.Pi * float64(i) * float64(x) / float64(width))
		}
	}
	for j := 0; j < numY; j++ {
		for y := 0; y < height; y++ {
			d.cosY[j*height+y] = math.Cos(math.Pi * float64(j) * float64(y) / float64(height))
		}
	}

	// Get direct pixel access if available
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

	// Account for sub-image offset
	minX, minY := bounds.Min.X, bounds.Min.Y

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var r, g, b float64
			for j := 0; j < numY; j++ {
				basisY := d.cosY[j*height+y]
				for i := 0; i < numX; i++ {
					basis := d.cosX[i*width+x] * basisY
					compColor := d.colors[i+j*numX]
					r += compColor[0] * basis
					g += compColor[1] * basis
					b += compColor[2] * basis
				}
			}

			if pix != nil {
				idx := (minY+y)*stride + (minX+x)*4
				pix[idx] = uint8(linearToSRGB(r))
				pix[idx+1] = uint8(linearToSRGB(g))
				pix[idx+2] = uint8(linearToSRGB(b))
				pix[idx+3] = 255
			} else {
				dst.Set(minX+x, minY+y, color.NRGBA{
					uint8(linearToSRGB(r)),
					uint8(linearToSRGB(g)),
					uint8(linearToSRGB(b)),
					255,
				})
			}
		}
	}

	return nil
}

func (d *Decoder) maybeGrowBuffers(width, height, numX, numY int) {
	d.cosX = growTo(d.cosX, numX*width)
	d.cosY = growTo(d.cosY, numY*height)
	d.colors = growTo(d.colors, numX*numY)
}

// Decode returns an NRGBA image of the given hash with the given size.
func Decode(hash string, width, height int, punch int) (image.Image, error) {
	var d Decoder
	return d.Decode(hash, width, height, punch)
}

// DecodeDraw decodes the given hash into the given image.
func DecodeDraw(dst draw.Image, hash string, punch float64) error {
	var d Decoder
	return d.DecodeDraw(dst, hash, punch)
}

func decodeDC(val int) (c [3]float64) {
	c[0] = sRGBToLinear(val >> 16 & 255)
	c[1] = sRGBToLinear(val >> 8 & 255)
	c[2] = sRGBToLinear(val & 255)
	return c
}

func decodeAC(val int, maximumValue float64) (c [3]float64) {
	quantR := val / 361 // 19*19
	quantG := (val / 19) % 19
	quantB := val % 19

	// signPow with exponent 2 is: sign(x) * x^2 = x * |x|
	rNorm := float64(quantR-9) / 9
	gNorm := float64(quantG-9) / 9
	bNorm := float64(quantB-9) / 9

	c[0] = rNorm * math.Abs(rNorm) * maximumValue
	c[1] = gNorm * math.Abs(gNorm) * maximumValue
	c[2] = bNorm * math.Abs(bNorm) * maximumValue
	return c
}
