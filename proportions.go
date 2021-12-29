package images

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
