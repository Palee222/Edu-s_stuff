package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"os"
)

func main() {
	const width, height, frameCount = 64, 16, 8
	palette := color.Palette{
		color.RGBA{A: 255},
		color.RGBA{R: 255, A: 255},
		color.RGBA{R: 255, G: 255, B: 255, A: 255},
	}

	animation := &gif.GIF{
		Image:     make([]*image.Paletted, frameCount),
		Delay:     make([]int, frameCount),
		LoopCount: 0,
	}
	for frameIndex := 0; frameIndex < frameCount; frameIndex++ {
		frame := image.NewPaletted(image.Rect(0, 0, width, height), palette)
		for bar := 0; bar < 4; bar++ {
			x := (frameIndex*4 + bar*16) % width
			for y := 2; y < height-2; y++ {
				for dx := 0; dx < 4; dx++ {
					frame.SetColorIndex((x+dx)%width, y, uint8(1+(bar%2)))
				}
			}
		}
		animation.Image[frameIndex] = frame
		animation.Delay[frameIndex] = 8
	}

	const path = "photosforuse/mask-test-moving-bars-valid.gif"
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	if err := gif.EncodeAll(file, animation); err != nil {
		file.Close()
		panic(err)
	}
	if err := file.Close(); err != nil {
		panic(err)
	}

	check, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	decoded, decodeErr := gif.DecodeAll(check)
	closeErr := check.Close()
	if decodeErr != nil {
		panic(decodeErr)
	}
	if closeErr != nil {
		panic(closeErr)
	}
	fmt.Printf("Wrote %s (%d frames)\n", path, len(decoded.Image))
}
