package mask

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"

	xdraw "golang.org/x/image/draw"
)

const (
	maskImageHeight  = 16
	maskImageMaxCols = 4998
	maskPixelCutoff  = 64
)

// SetImageFromFile decodes an image file and uploads it in the mask's image format.
func SetImageFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open image %q: %w", path, err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("decode image %q: %w", path, err)
	}

	bitmap, colorArray, err := EncodeImageForMask(img)
	if err != nil {
		return err
	}
	return InitUpload(bitmap, colorArray)
}

// EncodeImageForMask scales an image to 16 pixels high and encodes its pixels
// and per-column colors for upload to the mask.
func EncodeImageForMask(img image.Image) ([]byte, []byte, error) {
	if img == nil {
		return nil, nil, fmt.Errorf("image is nil")
	}
	bounds := img.Bounds()
	if bounds.Empty() {
		return nil, nil, fmt.Errorf("image has no pixels")
	}

	width := int(math.Round(float64(bounds.Dx()) * maskImageHeight / float64(bounds.Dy())))
	if width < 1 {
		width = 1
	}
	if width > maskImageMaxCols {
		return nil, nil, fmt.Errorf("image is too wide for the mask after scaling (%d columns; maximum %d)", width, maskImageMaxCols)
	}

	resized := image.NewRGBA(image.Rect(0, 0, width, maskImageHeight))
	xdraw.CatmullRom.Scale(resized, resized.Bounds(), img, bounds, xdraw.Over, nil)

	pixelMap := make([][]byte, width)
	colorArray := make([]byte, 0, width*3)
	for x := 0; x < width; x++ {
		column := make([]byte, maskImageHeight)
		var red, green, blue uint64
		litPixels := uint64(0)

		for y := 0; y < maskImageHeight; y++ {
			pixel := resized.RGBAAt(x, y)
			luminance := (299*uint32(pixel.R) + 587*uint32(pixel.G) + 114*uint32(pixel.B)) / 1000
			if luminance < maskPixelCutoff {
				continue
			}

			column[y] = 1
			red += uint64(pixel.R)
			green += uint64(pixel.G)
			blue += uint64(pixel.B)
			litPixels++
		}

		pixelMap[x] = column
		if litPixels == 0 {
			colorArray = append(colorArray, 0, 0, 0)
		} else {
			colorArray = append(colorArray,
				byte(red/litPixels),
				byte(green/litPixels),
				byte(blue/litPixels))
		}
	}

	bitmap, err := EncodeBitmapForMask(pixelMap)
	if err != nil {
		return nil, nil, err
	}
	return bitmap, colorArray, nil
}

// ClearImage uploads a single blank column to replace the currently uploaded image.
func ClearImage() error {
	return InitUpload([]byte{0, 0}, []byte{0, 0, 0})
}
