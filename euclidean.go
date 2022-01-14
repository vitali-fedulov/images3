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
