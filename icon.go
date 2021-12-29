package images

import (
	"image"
	"image/color"
)

// Icon parameters.
const (
	// Side dimension of an intermediate icon. Image will be
	// resampled to this size. The value is odd to match later
	// box filter parameters with stride 2.
	iconLargeSize = 23
	// Side dimension (in pixels) of a downsample square
	// to reasonably well approximate color area of a full
	// size image.
	sampleSide       = 12
	invSamplePixels2 = 1 / float32(sampleSide*sampleSide)
	resizedImgSize   = iconLargeSize * sampleSide
	iconSmallSize    = iconLargeSize / 2
	numIcon2Pixels   = iconSmallSize * iconSmallSize
	oneNinth         = float32(1) / float32(9)
)

// Icon has square shape. Its array contains pixel values
// for 3 channels. Float is intentional to preserve color
// relationships from the image of original size.
// Color values range is [0.0, 255.0].
type IconT []float32

type Point struct{ X, Y int }

// Icon generates the signature (icon) and gets original
// image size. The icon data can then be stored in a database
// and used for comparisons. Image resolution of the icon
// is very small (11x11 pixels): original image details
// are lost in downsampling, except when source images have
// very low resolution (e.g. favicons or simple logos).
func Icon(img image.Image) (icon IconT, imgSize Point) {
	icon1, imgSizeX, imgSizeY := IconLarge(img)
	return IconSmall(icon1), Point{imgSizeX, imgSizeY}
}

// IconLarge resizes a source image to an icon approximating
// average color values. The function also returns the original
// image size. YCbCr space is used to take advantage of the luma
// component.
func IconLarge(img image.Image) (
	icon IconT, imgSizeX, imgSizeY int) {
	// Image is resampled to the icon size.
	resImg, imgSizeX, imgSizeY :=
		ResampleByNearest(img, resizedImgSize, resizedImgSize)
	icon = NewIcon(iconLargeSize)
	var r, g, b, sumR, sumG, sumB uint32
	// For each pixel of the icon.
	for x := 0; x < iconLargeSize; x++ {
		for y := 0; y < iconLargeSize; y++ {
			sumR, sumG, sumB = 0, 0, 0
			// Sum over pixels of resImg.
			for m := 0; m < sampleSide; m++ {
				for n := 0; n < sampleSide; n++ {
					r, g, b, _ =
						resImg.At(
							x*sampleSide+m, y*sampleSide+n).RGBA()
					sumR += r >> 8
					sumG += g >> 8
					sumB += b >> 8
				}
			}
			yc, cb, cr := yCbCr(
				float32(sumR)*invSamplePixels2,
				float32(sumG)*invSamplePixels2,
				float32(sumB)*invSamplePixels2)
			Set(icon, iconLargeSize, Point{x, y}, yc, cb, cr)
		}
	}
	return icon, imgSizeX, imgSizeY
}

// IconSmall returns a smaller icon with pixels representing
// box blur values from IconLarge.
func IconSmall(src IconT) (dst IconT) {
	dst = NewIcon(iconSmallSize)
	// Pixel positions in the destination icon.
	var xd, yd int
	var c1, c2, c3, s1, s2, s3 float32
	// For pixels of source icon with stride 2.
	for x := 1; x < iconLargeSize-1; x += 2 {
		xd = x / 2
		for y := 1; y < iconLargeSize-1; y += 2 {
			yd = y / 2
			// For each pixel of a 3x3 box.
			for n := -1; n <= 1; n++ {
				for m := -1; m <= 1; m++ {
					c1, c2, c3 =
						Get(src, iconLargeSize, Point{x + n, y + m})
					s1, s2, s3 = s1+c1, s2+c2, s3+c3
				}
			}
			Set(dst, iconSmallSize, Point{xd, yd},
				s1*oneNinth, s2*oneNinth, s3*oneNinth)
			s1, s2, s3 = 0, 0, 0
		}
	}
	return Normalize(dst, numIcon2Pixels)
}

// NewIcon is an icon maker.
func NewIcon(size int) IconT {
	arr := make([]float32, size*size*3) // 3 channels.
	return IconT(arr)
}

// ArrIndex gets a pixel position in 1D array from a point
// of 2D array. ch is color channel index (0 to 2).
func ArrIndex(p Point, size, ch int) (index int) {
	return size*(ch*size+p.Y) + p.X
}

// Set places pixel values in an icon at a point.
// c1, c2, c3 are color values for each channel
// (RGB for example).
func Set(icon IconT, size int, p Point, c1, c2, c3 float32) {
	icon[ArrIndex(p, size, 0)] = c1
	icon[ArrIndex(p, size, 1)] = c2
	icon[ArrIndex(p, size, 2)] = c3
}

// Get reads pixel values in an icon at a point.
// c1, c2, c3 are color values for each channel
// (RGB for example).
func Get(icon IconT, size int, p Point) (c1, c2, c3 float32) {
	c1 = icon[ArrIndex(p, size, 0)]
	c2 = icon[ArrIndex(p, size, 1)]
	c3 = icon[ArrIndex(p, size, 2)]
	return c1, c2, c3
}

// LumaValues returns luma values at sample pixels of the small icon
// (from IconSmall).
func LumaValues(icon IconT, sample []Point) (v []float64) {
	for i := range sample {
		c1, _, _ := Get(icon, iconSmallSize, sample[i])
		v = append(v, float64(c1))
	}
	return v
}

// Auxiliary function, used to evaluate icon-generation
// functions in a separate program, but not in tests, which
// could be brittle.
func ToRGBA(icon IconT, size int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			r, g, b := Get(icon, size, Point{x, y})
			img.Set(x, y,
				color.RGBA{uint8(r), uint8(g), uint8(b), 255})
		}
	}
	return img
}
