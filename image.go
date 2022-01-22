package images3

import (
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
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

// ResizeByNearest resizes an image by the nearest neighbour method to the
// output size outX, outY. It also returns the size inX, inY of the input image.
func ResizeByNearest(src image.Image, dstX, dstY int) (dst image.RGBA,
	srcX, srcY int) {
	// Original image size.
	xMax, xMin := src.Bounds().Max.X, src.Bounds().Min.X
	yMax, yMin := src.Bounds().Max.Y, src.Bounds().Min.Y
	srcX = xMax - xMin
	srcY = yMax - yMin

	// Destination rectangle.
	outRect := image.Rectangle{image.Point{0, 0}, image.Point{dstX, dstY}}
	// Color model of uint8 per color.
	dst = *image.NewRGBA(outRect)
	var (
		r, g, b, a uint32
	)
	for y := 0; y < dstY; y++ {
		for x := 0; x < dstX; x++ {
			r, g, b, a = src.At(
				x*srcX/dstX+xMin,
				y*srcY/dstY+yMin).RGBA()
			dst.Set(x, y, color.RGBA{
				uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)})
		}
	}
	return dst, srcX, srcY
}

// SaveToPNG encodes and saves image.RGBA to a file.
func SaveToPNG(img *image.RGBA, path string) {
	if destFile, err := os.Create(path); err != nil {
		log.Println("Cannot create file: ", path, err)
	} else {
		defer destFile.Close()
		png.Encode(destFile, img)
	}
}

// SaveToJPG encodes and saves image.RGBA to a file.
func SaveToJPG(img *image.RGBA, path string, quality int) {
	if destFile, err := os.Create(path); err != nil {
		log.Println("Cannot create file: ", path, err)
	} else {
		defer destFile.Close()
		jpeg.Encode(destFile, img, &jpeg.Options{Quality: quality})
	}
}
