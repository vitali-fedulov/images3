package images

// Default Euclidean similarity parameters.
const (
	// Cutoff value for color distance.
	colorDiff = 50
	// Cutoff coefficient for Euclidean distance (squared).
	euclCoeff = 0.2
	// Coefficient of sensitivity for Cb/Cr channels vs Y.
	chanCoeff = 2

	// Euclidean distance threshold (squared) for Y-channel.
	euclDist2Y = float32(numIcon2Pixels) *
		float32(colorDiff*colorDiff) * euclCoeff
	// Euclidean distance threshold (squared) for Cb and Cr
	// channels.
	euclDist2CbCr = euclDist2Y * chanCoeff
)

// Default proportion similarity threshold.
const threshold = 0.05

// PropSimilar gives a similarity verdict for image A and B based on
// their height and width. When proportions are similar, it returns
// true. The function uses recommended threshold constant and wraps
// customizable PropSimilarCustom function.
func PropSimilar(imgSizeA, imgSizeB Point) bool {
	return PropSimilarCustom(imgSizeA, imgSizeB, threshold)
}

// PropSimilarCustom gives a similarity verdict for image A and B
// based on their height, width, and cuttoff threshold. When proportions
// are similar, it returns true.
func PropSimilarCustom(imgSizeA, imgSizeB Point, threshold float64) bool {

	// Filtering is based on rescaling a narrower side of images to 1,
	// then cutting off at threshold of a longer image vs shorter image.
	xA, yA := float64(imgSizeA.X), float64(imgSizeA.Y)
	xB, yB := float64(imgSizeB.X), float64(imgSizeB.Y)
	var delta float64
	if xA <= yA { // x to 1.
		yA = yA / xA
		yB = yB / xB
		if yA > yB {
			delta = (yA - yB) / yA
		} else {
			delta = (yB - yA) / yB
		}
	} else { // y to 1.
		xA = xA / yA
		xB = xB / yB
		if xA > xB {
			delta = (xA - xB) / xA
		} else {
			delta = (xB - xA) / xB
		}
	}
	if delta > threshold {
		return false
	}
	return true
}

// EucSimilar gives a similarity verdict for image A and B based
// on Euclidean distance between pixel values of their icons.
// When the distance is small, the function returns true.
// iconA and iconB are generated with the Icon function.
// EucSimilar is recommended for general use, because it uses
// thorougly tested thresholds. It wraps EucMetrics with
// those thresholds.
func EucSimilar(iconA, iconB IconT) bool {
	m1, m2, m3 := EucMetrics(iconA, iconB)
	return m1 < euclDist2Y && m2 < euclDist2CbCr && m3 < euclDist2CbCr
}

// EucMetrics returns Euclidean distances between 2 icons.
// These are 3 metrics corresponding to each color channel.
// The distances are squared to avoid square root calculations.
func EucMetrics(iconA, iconB IconT) (m1, m2, m3 float32) {

	var cA, cB float32
	for i := 0; i < numIcon2Pixels; i++ {
		// Channel 1.
		cA = iconA[i]
		cB = iconB[i]
		m1 += (cA - cB) * (cA - cB)
		// Channel 2.
		cA = iconA[i+numIcon2Pixels]
		cB = iconB[i+numIcon2Pixels]
		m2 += (cA - cB) * (cA - cB)
		// Channel 3.
		cA = iconA[i+2*numIcon2Pixels]
		cB = iconB[i+2*numIcon2Pixels]
		m3 += (cA - cB) * (cA - cB)
	}
	return m1, m2, m3
}
