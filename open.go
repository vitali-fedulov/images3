package images

import (
	"image"
	"os"
)

// Open opens and decodes an image file for a given path.
func Open(path string) (img image.Image, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err = image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, err
}
