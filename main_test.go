package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestRemovePhotoFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "photo.png")
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("create image: %v", err)
	}
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.White)
	if err := png.Encode(file, img); err != nil {
		file.Close()
		t.Fatalf("encode image: %v", err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("close image: %v", err)
	}

	if err := removePhotoFile(path); err != nil {
		t.Fatalf("removePhotoFile() error = %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("image still exists or stat failed unexpectedly: %v", err)
	}
}

func TestRemovePhotoFileRejectsNonImage(t *testing.T) {
	path := filepath.Join(t.TempDir(), "notes.txt")
	if err := os.WriteFile(path, []byte("not an image"), 0600); err != nil {
		t.Fatalf("create test file: %v", err)
	}

	if err := removePhotoFile(path); err == nil {
		t.Fatal("removePhotoFile() error = nil, want an error")
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("non-image file was removed: %v", err)
	}
}
