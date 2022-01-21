package images3

const (

	// Euclidean similarity parameters.

	// Cutoff value for color distance.
	colorDiff = 50
	// Cutoff coefficient for Euclidean distance (squared).
	euclCoeff = 0.2
	// Coefficient of sensitivity for Cb/Cr channels vs Y.
	chanCoeff = 2
	// Euclidean distance threshold (squared) for Y-channel.
	thY = float32(iconSize*iconSize) * float32(colorDiff*colorDiff) * euclCoeff
	// Euclidean distance threshold (squared) for Cb and Cr channels.
	thCbCr = thY * chanCoeff

	// Proportion similarity threshold 5%.
	thProp = 0.05
)

// Similar returns similarity verdict based on Euclidean
// and proportion similarity.
func Similar(iconA, iconB IconT) bool {

	if PropSimilar(iconA, iconB) {

		if EucSimilar(iconA, iconB) {
			return true
		}
	}
	return false
}

// PropSimilar gives a similarity verdict for image A and B based on
// their height and width. When proportions are similar, it returns
// true.
func PropSimilar(iconA, iconB IconT) bool {
	return PropMetric(iconA, iconB) < thProp
}

// PropMetric gives image proportion similarity metric for image A
// and B. The smaller the metric the more similar are images by their
// x-y size.
func PropMetric(iconA, iconB IconT) (m float64) {

	// Filtering is based on rescaling a narrower side of images to 1,
	// then cutting off at threshold of a longer image vs shorter image.
	xA, yA := float64(iconA.ImgSize.X), float64(iconA.ImgSize.Y)
	xB, yB := float64(iconB.ImgSize.X), float64(iconB.ImgSize.Y)

	if xA <= yA { // x to 1.
		yA = yA / xA
		yB = yB / xB
		if yA > yB {
			m = (yA - yB) / yA
		} else {
			m = (yB - yA) / yB
		}
	} else { // y to 1.
		xA = xA / yA
		xB = xB / yB
		if xA > xB {
			m = (xA - xB) / xA
		} else {
			m = (xB - xA) / xB
		}
	}
	return m
}

// EucSimilar gives a similarity verdict for image A and B based
// on Euclidean distance between pixel values of their icons.
// When the distance is small, the function returns true.
// iconA and iconB are generated with the Icon function.
// EucSimilar wraps EucMetric with well-tested thresholds.
func EucSimilar(iconA, iconB IconT) bool {

	m1, m2, m3 := EucMetric(iconA, iconB)
	return m1 < thY && m2 < thCbCr && m3 < thCbCr
}

// EucMetric returns Euclidean distances between 2 icons.
// These are 3 metrics corresponding to each color channel.
// The distances are squared to avoid square root calculations.
// Note that the channels are not RGB, but YCbCr, thus their
// importance for similarity might be not the same.
func EucMetric(iconA, iconB IconT) (m1, m2, m3 float32) {

	numIconPixels := iconSize * iconSize
	var cA, cB float32
	for i := 0; i < numIconPixels; i++ {
		// Channel 1.
		cA = iconA.Pixels[i]
		cB = iconB.Pixels[i]
		m1 += (cA - cB) * (cA - cB)
		// Channel 2.
		cA = iconA.Pixels[i+numIconPixels]
		cB = iconB.Pixels[i+numIconPixels]
		m2 += (cA - cB) * (cA - cB)
		// Channel 3.
		cA = iconA.Pixels[i+2*numIconPixels]
		cB = iconB.Pixels[i+2*numIconPixels]
		m3 += (cA - cB) * (cA - cB)
	}

	return m1, m2, m3
}
