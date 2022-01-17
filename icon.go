package images

import (
	"image"
	"image/color"
)

// TODO
// Image resolution of the icon
// is very small (11x11 pixels): original image details
// are lost in downsampling, except when source images have
// very low resolution (e.g. favicons or simple logos).

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
	oneNinth         = float32(1) / float32(9)
)

// Icon has square shape. It pixels are float32 values
// for 3 channels. Float32 is intentional to preserve color
// relationships from the full-size image.
type IconT struct {
	Pixels  []float32
	ImgSize Point  // Original image size.
	Path    string // Original image path.
}

type Point image.Point

// Icon generates image signature (icon) with related info. TODO Make a custom icon function for specific size, then use it for Icon11.
// The icon data can then be stored in a database and used
// for comparisons. Icon wraps IconNonNorm function
// by stretching histograms of icons similar to Auto Levels
// in Photoshop. It is recommended to use this function
// instead of IconNonNorm, because it showed better
// results for image comparison, especially for images
// with large uniform backgrounds (for example pen drawings
// or images of blue sky with a small object flying).
func Icon(img image.Image, path string) IconT {
	icon := IconNonNorm(img, path)
	icon.normalize(iconSmallSize) // TODO: Test Icon vs IconNonNorm.
	return icon                   // TODO: Make sure this one is used in DEMO. Normalize must preserve icon image path and image size. Make tests for these.
}

// IconNonNorm is similar to function Icon, except
// it does not normalize icon histograms. If used instead
// of Icon, it could produce false positives for images
// with large uniform backgrounds.
func IconNonNorm(img image.Image, path string) IconT {

	// Resizing the source image to level-1 icon size approximating
	// average color values. YCbCr space is used to take advantage
	// of the luma component.

	resImg, imgSizeX, imgSizeY :=
		ResampleByNearest(img, resizedImgSize, resizedImgSize)
	largeIcon := NewIcon(iconLargeSize)
	var r, g, b, sumR, sumG, sumB uint32
	// For each pixel of the largeIcon.
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
			set(largeIcon, iconLargeSize,
				Point{x, y}, yc, cb, cr)
		}
	}

	// Box blur filter with resizing to level-2 icon size. TODO Make size constants dependent on only smallest (11 pixels).

	smallIcon := NewIcon(iconSmallSize)
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
						get(largeIcon, iconLargeSize,
							Point{x + n, y + m})
					s1, s2, s3 = s1+c1, s2+c2, s3+c3
				}
			}
			set(smallIcon, iconSmallSize, Point{xd, yd},
				s1*oneNinth, s2*oneNinth, s3*oneNinth)
			s1, s2, s3 = 0, 0, 0
		}
	}

	smallIcon.ImgSize = Point{imgSizeX, imgSizeY}
	smallIcon.Path = path

	return smallIcon
}

func NewIcon(iconSize int) (icon IconT) {
	icon.Pixels = make([]float32, iconSize*iconSize*3)
	return icon
}

// ArrIndex gets a pixel position in 1D array from a point
// of 2D array. ch is color channel index (0 to 2).
func arrIndex(p Point, size, ch int) (index int) {
	return size*(ch*size+p.Y) + p.X
}

// Set places pixel values in an icon at a point.
// c1, c2, c3 are color values for each channel
// (RGB for example). Size is icon size.
func set(icon IconT, size int, p Point, c1, c2, c3 float32) {
	icon.Pixels[arrIndex(p, size, 0)] = c1
	icon.Pixels[arrIndex(p, size, 1)] = c2
	icon.Pixels[arrIndex(p, size, 2)] = c3
}

// Get reads pixel values in an icon at a point.
// c1, c2, c3 are color values for each channel
// (RGB for example).
func get(icon IconT, size int, p Point) (c1, c2, c3 float32) {
	c1 = icon.Pixels[arrIndex(p, size, 0)]
	c2 = icon.Pixels[arrIndex(p, size, 1)]
	c3 = icon.Pixels[arrIndex(p, size, 2)]
	return c1, c2, c3
}

// yCbCr transforms RGB components to YCbCr. This is a high
// precision version different from the Golang image library
// operating on uint8.
func yCbCr(r, g, b float32) (yc, cb, cr float32) {
	yc = 0.299000*r + 0.587000*g + 0.114000*b
	cb = 128 - 0.168736*r - 0.331264*g + 0.500000*b
	cr = 128 + 0.500000*r - 0.418688*g - 0.081312*b
	return yc, cb, cr
}

// LumaValues returns luma values at sample pixels of the small icon
// (from IconSmall).
func LumaValues(icon IconT, sample []Point) (v []float64) {
	for i := range sample {
		c1, _, _ := get(icon, iconSmallSize, sample[i])
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
			r, g, b := get(icon, size, Point{x, y})
			img.Set(x, y,
				color.RGBA{uint8(r), uint8(g), uint8(b), 255})
		}
	}
	return img
}

// Normalize stretches histograms for the 3 channels of an icon, so that
// minimum and maximum values of each are 0 and 255 correspondingly.
// iconSmallSize is the icon side size in pixels.
func (src IconT) normalize(iconSmallSize int) {

	numPixels := iconSmallSize * iconSmallSize

	var c1Min, c2Min, c3Min, c1Max, c2Max, c3Max float32
	c1Min, c2Min, c3Min = 256, 256, 256
	c1Max, c2Max, c3Max = 0, 0, 0
	var n int

	// Looking for extreme values.
	for n = 0; n < numPixels; n++ {
		// Channel 1.
		if src.Pixels[n] > c1Max {
			c1Max = src.Pixels[n]
		}
		if src.Pixels[n] < c1Min {
			c1Min = src.Pixels[n]
		}
		// Channel 2.
		if src.Pixels[n+numPixels] > c2Max {
			c2Max = src.Pixels[n+numPixels]
		}
		if src.Pixels[n+numPixels] < c2Min {
			c2Min = src.Pixels[n+numPixels]
		}
		// Channel 3.
		if src.Pixels[n+2*numPixels] > c3Max {
			c3Max = src.Pixels[n+2*numPixels]
		}
		if src.Pixels[n+2*numPixels] < c3Min {
			c3Min = src.Pixels[n+2*numPixels]
		}
	}

	// Normalization.
	rCoeff := 255 / (c1Max - c1Min)
	gCoeff := 255 / (c2Max - c2Min)
	bCoeff := 255 / (c3Max - c3Min)
	for n = 0; n < numPixels; n++ {
		src.Pixels[n] =
			(src.Pixels[n] - c1Min) * rCoeff
		src.Pixels[n+numPixels] =
			(src.Pixels[n+numPixels] - c2Min) * gCoeff
		src.Pixels[n+2*numPixels] =
			(src.Pixels[n+2*numPixels] - c3Min) * bCoeff
	}

}
