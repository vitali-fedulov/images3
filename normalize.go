package images

// Normalize stretches histograms for the 3 channels of an icon, so that
// minimum and maximum values of each are 0 and 255 correspondingly.
// numPixels is number of pixels in an icon.
func Normalize(src IconT, numPixels int) IconT {

	dst := make([]float32, numPixels*3)
	var c1Min, c2Min, c3Min, c1Max, c2Max, c3Max float32
	c1Min, c2Min, c3Min = 256, 256, 256
	c1Max, c2Max, c3Max = 0, 0, 0
	var n int

	// Looking for extreme values.
	for n = 0; n < numPixels; n++ {
		// Channel 1.
		if src[n] > c1Max {
			c1Max = src[n]
		}
		if src[n] < c1Min {
			c1Min = src[n]
		}
		// Channel 2.
		if src[n+numPixels] > c2Max {
			c2Max = src[n+numPixels]
		}
		if src[n+numPixels] < c2Min {
			c2Min = src[n+numPixels]
		}
		// Channel 3.
		if src[n+2*numPixels] > c3Max {
			c3Max = src[n+2*numPixels]
		}
		if src[n+2*numPixels] < c3Min {
			c3Min = src[n+2*numPixels]
		}
	}

	// Normalization.
	rCoeff := 255 / (c1Max - c1Min)
	gCoeff := 255 / (c2Max - c2Min)
	bCoeff := 255 / (c3Max - c3Min)
	for n = 0; n < numPixels; n++ {
		dst[n] = (src[n] - c1Min) * rCoeff
		dst[n+numPixels] = (src[n+numPixels] - c2Min) * gCoeff
		dst[n+2*numPixels] = (src[n+2*numPixels] - c3Min) * bCoeff
	}

	return dst
}
