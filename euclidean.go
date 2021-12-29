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
// thorougly tested constants. It wraps EucSimilarCustom with
// those constants.
func EucSimilar(iconA, iconB IconT) bool {
	return EucSimilarCustom(iconA, iconB, euclDist2Y, euclDist2CbCr)
}

// EucSimilarCustom is a customizable function, where parameters
// euclDist2Y and euclDist2CbCr can be changed to non-default
// distance coefficients.
func EucSimilarCustom(
	iconA, iconB IconT, euclDist2Y, euclDist2CbCr float32) bool {

	var cA, cB, s float32

	// Euclidean distance filter on Y-channel.
	for i := 0; i < numIcon2Pixels; i++ {
		cA = iconA[i]
		cB = iconB[i]
		s += (cA - cB) * (cA - cB)
	}
	if s > euclDist2Y {
		return false
	}
	s = 0

	// Euclidean distance filter on Cb-channel.
	for i := 0; i < numIcon2Pixels; i++ {
		cA = iconA[numIcon2Pixels+i]
		cB = iconB[numIcon2Pixels+i]
		s += (cA - cB) * (cA - cB)
	}
	if s > euclDist2CbCr {
		return false
	}
	s = 0

	// Euclidean distance filter on Cr-channel.
	for i := 0; i < numIcon2Pixels; i++ {
		cA = iconA[numIcon2Pixels+numIcon2Pixels+i]
		cB = iconB[numIcon2Pixels+numIcon2Pixels+i]
		s += (cA - cB) * (cA - cB)
	}
	if s > euclDist2CbCr {
		return false
	}

	return true
}
