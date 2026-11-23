package mask

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"os"
	"time"
)

const (
	defaultGIFFrameDelay   = 100 * time.Millisecond
	animationUploadTimeout = 60 * time.Second
)

type encodedGIFFrame struct {
	bitmap     []byte
	colorArray []byte
	delay      time.Duration
}

// PlayGIFFromFile plays an animated GIF by uploading each frame in sequence.
// It blocks until playback finishes; GIFs configured to loop forever run until interrupted.
func PlayGIFFromFile(path string) error {
	if !IsConnected() {
		return fmt.Errorf("not connected")
	}

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open GIF %q: %w", path, err)
	}
	animation, decodeErr := gif.DecodeAll(file)
	closeErr := file.Close()
	if decodeErr != nil {
		return fmt.Errorf("decode GIF %q: %w", path, decodeErr)
	}
	if closeErr != nil {
		return fmt.Errorf("close GIF %q: %w", path, closeErr)
	}

	frames, err := encodeGIFFrames(animation)
	if err != nil {
		return fmt.Errorf("prepare GIF %q: %w", path, err)
	}

	infinite := animation.LoopCount == 0
	plays := 1
	if animation.LoopCount > 0 {
		plays = animation.LoopCount + 1
	}
	for play := 0; infinite || play < plays; play++ {
		for i, frame := range frames {
			started := time.Now()
			done := make(chan error, 1)
			if err := initUpload(frame.bitmap, frame.colorArray, done); err != nil {
				return fmt.Errorf("upload GIF frame %d: %w", i+1, err)
			}

			timer := time.NewTimer(animationUploadTimeout)
			select {
			case err := <-done:
				if !timer.Stop() {
					select {
					case <-timer.C:
					default:
					}
				}
				if err != nil {
					return fmt.Errorf("upload GIF frame %d: %w", i+1, err)
				}
			case <-timer.C:
				return fmt.Errorf("timed out uploading GIF frame %d", i+1)
			}

			if remaining := frame.delay - time.Since(started); remaining > 0 {
				time.Sleep(remaining)
			}
		}
	}
	return nil
}

func encodeGIFFrames(animation *gif.GIF) ([]encodedGIFFrame, error) {
	if animation == nil || len(animation.Image) == 0 {
		return nil, fmt.Errorf("GIF has no frames")
	}

	width, height := animation.Config.Width, animation.Config.Height
	if width <= 0 || height <= 0 {
		bounds := animation.Image[0].Bounds()
		width, height = bounds.Max.X, bounds.Max.Y
	}
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("GIF has invalid dimensions")
	}

	canvas := image.NewRGBA(image.Rect(0, 0, width, height))
	frames := make([]encodedGIFFrame, 0, len(animation.Image))
	var previousCanvas *image.RGBA
	var previousBounds image.Rectangle
	var previousDisposal byte

	for i, frame := range animation.Image {
		if frame == nil {
			return nil, fmt.Errorf("frame %d is empty", i+1)
		}
		if i > 0 {
			switch previousDisposal {
			case gif.DisposalBackground:
				background := color.Color(color.Transparent)
				if palette, ok := animation.Config.ColorModel.(color.Palette); ok &&
					int(animation.BackgroundIndex) < len(palette) {
					background = palette[animation.BackgroundIndex]
				}
				draw.Draw(canvas, previousBounds, image.NewUniform(background), image.Point{}, draw.Src)
			case gif.DisposalPrevious:
				if previousCanvas != nil {
					canvas = previousCanvas
				}
			}
		}

		disposal := byte(gif.DisposalNone)
		if i < len(animation.Disposal) {
			disposal = animation.Disposal[i]
		}
		if disposal == gif.DisposalPrevious {
			previousCanvas = cloneRGBA(canvas)
		} else {
			previousCanvas = nil
		}

		draw.Draw(canvas, frame.Bounds(), frame, frame.Bounds().Min, draw.Over)
		bitmap, colorArray, err := EncodeImageForMask(canvas)
		if err != nil {
			return nil, fmt.Errorf("encode frame %d: %w", i+1, err)
		}

		delay := defaultGIFFrameDelay
		if i < len(animation.Delay) && animation.Delay[i] > 0 {
			delay = time.Duration(animation.Delay[i]) * 10 * time.Millisecond
		}
		frames = append(frames, encodedGIFFrame{
			bitmap:     bitmap,
			colorArray: colorArray,
			delay:      delay,
		})
		previousBounds = frame.Bounds().Intersect(canvas.Bounds())
		previousDisposal = disposal
	}
	return frames, nil
}

func cloneRGBA(src *image.RGBA) *image.RGBA {
	dst := image.NewRGBA(src.Bounds())
	draw.Draw(dst, dst.Bounds(), src, src.Bounds().Min, draw.Src)
	return dst
}
