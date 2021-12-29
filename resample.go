package images

import (
	"image"
	"image/color"
)

// ResampleByNearest resizes an image by the nearest neighbour method to the
// output size outX, outY. It also returns the size inX, inY of the input image.
func ResampleByNearest(src image.Image, dstX, dstY int) (dst image.RGBA,
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
	for x := 0; x < dstX; x++ {
		for y := 0; y < dstY; y++ {
			r, g, b, a = src.At(
				x*srcX/dstX+xMin,
				y*srcY/dstY+yMin).RGBA()
			dst.Set(x, y, color.RGBA{
				uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)})
		}
	}
	return dst, srcX, srcY
}
