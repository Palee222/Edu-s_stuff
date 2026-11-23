package mask

import (
	"image"
	"image/color"
	"image/gif"
	"testing"
	"time"
)

func TestEncodeGIFFramesCompositesFramesAndHonorsDisposal(t *testing.T) {
	palette := color.Palette{
		color.RGBA{A: 0},
		color.RGBA{R: 255, A: 255},
		color.RGBA{G: 255, A: 255},
	}
	first := image.NewPaletted(image.Rect(0, 0, 2, 1), palette)
	first.SetColorIndex(0, 0, 1)
	first.SetColorIndex(1, 0, 1)
	second := image.NewPaletted(image.Rect(1, 0, 2, 1), palette)
	second.SetColorIndex(1, 0, 2)
	third := image.NewPaletted(image.Rect(0, 0, 1, 1), palette)
	third.SetColorIndex(0, 0, 0)

	animation := &gif.GIF{
		Image:    []*image.Paletted{first, second, third},
		Delay:    []int{5, 7, 9},
		Disposal: []byte{gif.DisposalNone, gif.DisposalPrevious, gif.DisposalNone},
		Config:   image.Config{ColorModel: palette, Width: 2, Height: 1},
	}

	frames, err := encodeGIFFrames(animation)
	if err != nil {
		t.Fatalf("encodeGIFFrames() error = %v", err)
	}
	if len(frames) != 3 {
		t.Fatalf("frame count = %d, want 3", len(frames))
	}
	if frames[0].delay != 50*time.Millisecond {
		t.Errorf("first frame delay = %v, want 50ms", frames[0].delay)
	}

	assertColumnColor(t, frames[0], 3, color.RGBA{R: 255, A: 255})
	assertColumnColor(t, frames[1], 3, color.RGBA{R: 255, A: 255})
	assertColumnColor(t, frames[1], 28, color.RGBA{G: 255, A: 255})
	assertColumnColor(t, frames[2], 3, color.RGBA{R: 255, A: 255})
	assertColumnColor(t, frames[2], 28, color.RGBA{R: 255, A: 255})
}

func TestEncodeGIFFramesRejectsNoFrames(t *testing.T) {
	if _, err := encodeGIFFrames(&gif.GIF{}); err == nil {
		t.Fatal("encodeGIFFrames() error = nil, want an error")
	}
}

func assertColumnColor(t *testing.T, frame encodedGIFFrame, column int, want color.RGBA) {
	t.Helper()
	got := frame.colorArray[column*3 : column*3+3]
	if got[0] != want.R || got[1] != want.G || got[2] != want.B {
		t.Errorf("column %d color = %v, want RGB(%d, %d, %d)", column, got, want.R, want.G, want.B)
	}
}
