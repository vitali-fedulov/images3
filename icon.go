package images

import (
	"image"
	"image/color"
)

// Icon parameters.
const (
	iconSize = 11 // Image resolution of the icon
	// is very small (11x11 pixels), therefore original
	// image details are lost in downsampling, except
	// when source images have very low resolution
	// (e.g. favicons or simple logos). This is useful
	// from the privacy perspective if you are to use
	// generated icons in a large searchable database.
	samples          = 12
	largeIconSize    = iconSize*2 + 1
	resizedImgSize   = largeIconSize * samples
	invSamplePixels2 = 1 / float32(samples*samples)
	oneNinth         = float32(1) / float32(9)
)

// Icon has square shape. Its pixels are float32 values
// for 3 channels. Float32 is intentional to preserve color
// relationships from the full-size image.
type IconT struct {
	Pixels  []float32
	ImgSize Point  // Original image size.
	Path    string // Original image path.
}

type Point image.Point

// Icon generates image signature (icon) with related info.
// The icon data can then be stored in a database and used
// for comparisons.
func Icon(img image.Image, path string) IconT {

	// Resizing to a large icon approximating average color
	// values of the source image. YCbCr space is used instead
	// of RGB for better results in image comparison.
	resImg, imgSizeX, imgSizeY :=
		ResampleByNearest(img, resizedImgSize, resizedImgSize)
	largeIcon := sizedIcon(largeIconSize)
	var r, g, b, sumR, sumG, sumB uint32
	// For each pixel of the largeIcon.
	for x := 0; x < largeIconSize; x++ {
		for y := 0; y < largeIconSize; y++ {
			sumR, sumG, sumB = 0, 0, 0
			// Sum over pixels of resImg.
			for m := 0; m < samples; m++ {
				for n := 0; n < samples; n++ {
					r, g, b, _ =
						resImg.At(
							x*samples+m, y*samples+n).RGBA()
					sumR += r >> 8
					sumG += g >> 8
					sumB += b >> 8
				}
			}
			yc, cb, cr := yCbCr(
				float32(sumR)*invSamplePixels2,
				float32(sumG)*invSamplePixels2,
				float32(sumB)*invSamplePixels2)
			set(largeIcon, largeIconSize,
				Point{x, y}, yc, cb, cr)
		}
	}

	// Box blur filter with resizing to the final icon of smaller size.

	icon := sizedIcon(iconSize)
	// Pixel positions in the final icon.
	var xd, yd int
	var c1, c2, c3, s1, s2, s3 float32

	// For pixels of source largeIcon with stride 2.
	for x := 1; x < largeIconSize-1; x += 2 {
		xd = x / 2
		for y := 1; y < largeIconSize-1; y += 2 {
			yd = y / 2
			// For each pixel of a 3x3 box.
			for n := -1; n <= 1; n++ {
				for m := -1; m <= 1; m++ {
					c1, c2, c3 =
						get(largeIcon, largeIconSize,
							Point{x + n, y + m})
					s1, s2, s3 = s1+c1, s2+c2, s3+c3
				}
			}
			set(icon, iconSize, Point{xd, yd},
				s1*oneNinth, s2*oneNinth, s3*oneNinth)
			s1, s2, s3 = 0, 0, 0
		}
	}

	icon.ImgSize = Point{imgSizeX, imgSizeY}
	icon.Path = path
	icon.normalize(iconSize)

	return icon
}

// EmptyIcon is an icon constructor in case you need an icon
// with nil values, for example for convenient error handling.
// Then you can use icon.Pixels == nil condition.
func EmptyIcon() (icon IconT) {
	icon = sizedIcon(iconSize)
	icon.Pixels = nil
	return icon
}

func sizedIcon(size int) (icon IconT) {
	icon.Pixels = make([]float32, size*size*3)
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

// lumaValues returns luma values at sample pixels of the icon.
func lumaValues(icon IconT, sample []Point) (v []float64) {
	for i := range sample {
		c1, _, _ := get(icon, iconSize, sample[i])
		v = append(v, float64(c1))
	}
	return v
}

// ToRGBA transforms a sized icon to image.RGBA. This is
// an auxiliary function, used to visually evaluate an icon
// in a separate program (but not in tests, which could be brittle).
func (icon IconT) ToRGBA(size int) *image.RGBA {
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
func (src IconT) normalize(size int) {

	numPixels := size * size
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
