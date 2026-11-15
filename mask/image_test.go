package mask

import (
	"image"
	"image/color"
	"testing"
)

func TestEncodeImageForMask(t *testing.T) {
	tests := []struct {
		name       string
		pixel      color.RGBA
		wantBitmap []byte
		wantColor  []byte
	}{
		{
			name:       "white pixels are on",
			pixel:      color.RGBA{R: 255, G: 255, B: 255, A: 255},
			wantBitmap: []byte{0xff, 0xff},
			wantColor:  []byte{255, 255, 255},
		},
		{
			name:       "black pixels are off",
			pixel:      color.RGBA{A: 255},
			wantBitmap: []byte{0, 0},
			wantColor:  []byte{0, 0, 0},
		},
		{
			name:       "bright red pixels retain their color",
			pixel:      color.RGBA{R: 255, A: 255},
			wantBitmap: []byte{0xff, 0xff},
			wantColor:  []byte{255, 0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img := image.NewRGBA(image.Rect(0, 0, 1, maskImageHeight))
			for y := 0; y < maskImageHeight; y++ {
				img.SetRGBA(0, y, tt.pixel)
			}

			bitmap, colors, err := EncodeImageForMask(img)
			if err != nil {
				t.Fatalf("EncodeImageForMask() error = %v", err)
			}
			if len(bitmap) != len(tt.wantBitmap) || string(bitmap) != string(tt.wantBitmap) {
				t.Errorf("bitmap = %v, want %v", bitmap, tt.wantBitmap)
			}
			if len(colors) != len(tt.wantColor) || string(colors) != string(tt.wantColor) {
				t.Errorf("colors = %v, want %v", colors, tt.wantColor)
			}
		})
	}
}

func TestEncodeImageForMaskPreservesAspectRatio(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 2, 1))

	bitmap, colors, err := EncodeImageForMask(img)
	if err != nil {
		t.Fatalf("EncodeImageForMask() error = %v", err)
	}
	if got, want := len(bitmap), 2*maskImageHeight*2; got != want {
		t.Errorf("bitmap length = %d, want %d", got, want)
	}
	if got, want := len(colors), 2*maskImageHeight*3; got != want {
		t.Errorf("color array length = %d, want %d", got, want)
	}
}

func TestEncodeImageForMaskRejectsEmptyImage(t *testing.T) {
	_, _, err := EncodeImageForMask(image.NewRGBA(image.Rectangle{}))
	if err == nil {
		t.Fatal("EncodeImageForMask() error = nil, want an error")
	}
}
