package images

import (
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
)

// Png encodes and saves image.RGBA to a file.
func Png(img *image.RGBA, path string) {
	if destFile, err := os.Create(path); err != nil {
		log.Println("Cannot create file: ", path, err)
	} else {
		defer destFile.Close()
		png.Encode(destFile, img)
	}
}

// Jpg encodes and saves image.RGBA to a file.
func Jpg(img *image.RGBA, path string, quality int) {
	if destFile, err := os.Create(path); err != nil {
		log.Println("Cannot create file: ", path, err)
	} else {
		defer destFile.Close()
		jpeg.Encode(destFile, img, &jpeg.Options{Quality: quality})
	}
}
