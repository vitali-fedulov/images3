package images

// Normalize stretches histograms for the 3 channels of an icon, so that
// minimum and maximum values of each are 0 and 255 correspondingly.
// numPixels is number of pixels in an icon.
func Normalize(src IconT, numPixels int) IconT {

	dest := make([]float32, numPixels*3)
	var c1Min, c2Min, c3Min, c1Max, c2Max, c3Max float32
	c1Min, c2Min, c3Min = 256, 256, 256
	c1Max, c2Max, c3Max = 0, 0, 0
	var n int

	// Looking for extreme values.
	for n = 0; n < numPixels; n++ {
		if src[n] > c1Max {
			c1Max = src[n]
		}
		if src[n] < c1Min {
			c1Min = src[n]
		}
	}
	for n = numPixels; n < 2*numPixels; n++ {
		if src[n] > c2Max {
			c2Max = src[n]
		}
		if src[n] < c2Min {
			c2Min = src[n]
		}
	}
	for n = 2 * numPixels; n < 3*numPixels; n++ {
		if src[n] > c3Max {
			c3Max = src[n]
		}
		if src[n] < c3Min {
			c3Min = src[n]
		}
	}

	// Normalization.
	rMM := c1Max - c1Min
	gMM := c2Max - c2Min
	bMM := c3Max - c3Min
	for n = 0; n < numPixels; n++ {
		dest[n] = (src[n] - c1Min) * 255 / rMM
	}
	for n = numPixels; n < 2*numPixels; n++ {
		dest[n] = (src[n] - c2Min) * 255 / gMM
	}
	for n = 2 * numPixels; n < 3*numPixels; n++ {
		dest[n] = (src[n] - c3Min) * 255 / bMM
	}

	return dest
}
