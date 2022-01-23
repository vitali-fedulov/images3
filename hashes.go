package images3

import (
	"math"

	// N-dimensional space discretization and hashing module.
	"github.com/vitali-fedulov/hyper"
)

// CentralHash generates a central hash for a given icon
// by sampling luma values at well-distributed icon points
// (hyperPoints, HyperPoints10) and later using package "hyper".
// This hash can then be used for a record or a query.
// When used for a record, you will need a hash set made
// with func HashSet for a query. And vice versa.
// To better understand CentralHash, read the following doc:
// https://vitali-fedulov.github.io/algorithm-for-hashing-high-dimensional-float-vectors.html
func CentralHash(icon IconT, hyperPoints []Point,
	epsPercent float64, numBuckets int) uint64 {

	vector := lumaVector(icon, hyperPoints)
	cube := hyper.CentralCube(vector,
		hyper.Params{
			Min:        0,
			Max:        255,
			EpsPercent: epsPercent,
			NumBuckets: numBuckets})

	if numBuckets > 10 || len(hyperPoints) > 19 {
		return cube.FNV1aHash()
	}

	return cube.DecimalHash()
}

// HashSet generates a hash set for a given icon by sampling
// luma values of well-distributed icon points (hyperPoints,
// HyperPoints10) and later using package "hyper".
// This hash set can then be used for records or a query.
// When used for a query, you will need a hash made with
// func CentralHash as a record. And vice versa.
// To better understand HashSet, read the following doc:
// https://vitali-fedulov.github.io/algorithm-for-hashing-high-dimensional-float-vectors.html
func HashSet(icon IconT, hyperPoints []Point,
	epsPercent float64, numBuckets int) []uint64 {

	vector := lumaVector(icon, hyperPoints)
	cubeSet := hyper.CubeSet(vector,
		hyper.Params{
			Min:        0,
			Max:        255,
			EpsPercent: epsPercent,
			NumBuckets: numBuckets})

	if numBuckets > 10 || len(hyperPoints) > 19 {
		return cubeSet.HashSet((hyper.Cube).FNV1aHash)
	}

	return cubeSet.HashSet((hyper.Cube).DecimalHash)
}

// HyperPoints10 is a convenience 10-point predefined set with
// coordinates of icon values to become 10 dimensions needed
// for hash generation with package "hyper".
// The 10 points are the only pixels from an icon to be used
// for hash generation (unless you define your own set of hyper
// points with CustomPoints function, or manually.
// The 10 points have been modified manually a little to avoid
// texture-like symmetries.
var HyperPoints10 = []Point{
	{2, 5}, {3, 3}, {3, 8}, {4, 6}, {5, 2},
	{6, 4}, {6, 7}, {8, 2}, {8, 5}, {8, 8}}

// CustomPoints is a utility function to create hyper points similar
// to HyperPoints10. It is needed if you are planning to use
// the package with billions of images, and might need higher number
// of sample points (more dimensions). You may also decide to reduce
// number of dimensions in order to reduce number of hashes per image.
// In both cases CustomPoints will help generate point sets similar
// to HyperPoints10.
// The function chooses a set of points (pixels from Icon) placed apart
// as far as possible from each other to increase variable independence.
// Number of chosen points corresponds to the number of dimensions n.
// Brightness values at those points represent one coordinate each
// in n-dimensional space for hash generation with package "hyper".
// Final point patterns are somewhat irregular, which is good to avoid
// occasional mutual pixel dependence of textures in images.
// For cases of low n, to avoid texture-like symmetries and visible
// patterns, it is recommended to slightly modify point positions
// manually, and with that distribute points irregularly across the Icon.
func CustomPoints(n int) map[Point]bool {
	// margin is a number of pixels near icon border to be left unused,
	// as images tend to contain noisy information there.
	margin := 2
	if n > 11 {
		margin = 1
	}
	pts := make(map[Point]bool)
	// First point to be in the upper left corner.
	pts[Point{margin, margin}] = true
	sumDist := 0.0 // Sum of distances for pairs of points.
	// Initializing point positions.
	for len(pts) < n {
		// Sums of distances for each x, y to already pts points.
		d := make(map[Point]float64)
		for x := margin; x < iconSize-margin; x++ {
			for y := margin; y < iconSize-margin; y++ {
				if _, ok := pts[Point{x, y}]; ok {
					continue
				}
				for p := range pts {
					sumDist += 1 / distance(p, Point{x, y})
				}
				d[Point{x, y}] = 1 / sumDist
				sumDist = 0
			}
		}
		// Find the max distance point.
		maxPoint := maxKey(d)
		pts[maxPoint] = true
	}

	// Moving a point to a space with a larger distance
	// to a nearest point.
	for i := 0; i < 50; i++ {
		for p0 := range pts {
			pts0 := exclude(p0, pts)
			d := make(map[Point]float64)
			for x := margin; x < iconSize-margin; x++ {
				for y := margin; y < iconSize-margin; y++ {
					if _, ok := pts0[Point{x, y}]; ok {
						continue
					}
					n := nearest(pts0, Point{x, y})
					d[Point{x, y}] = distance(Point{x, y}, n)
				}
			}
			newPoint := maxKey(d)
			delete(pts, p0)
			pts[newPoint] = true
		}
	}
	return pts
}

// distance calculates distance between 2 points.
func distance(p1, p2 Point) float64 {
	return math.Sqrt(
		float64((p1.X-p2.X)*(p1.X-p2.X)) +
			float64((p1.Y-p2.Y)*(p1.Y-p2.Y)))
}

// minKey finds key for smallest value of a map.
func minKey(m map[Point]float64) (key Point) {
	minVal := math.Inf(1)
	for k, v := range m {
		if v < minVal {
			key = k
			minVal = v
		}
	}
	return key
}

// maxKey finds key for largest value of a map.
func maxKey(m map[Point]float64) (key Point) {
	maxVal := math.Inf(-1)
	for k, v := range m {
		if v > maxVal {
			key = k
			maxVal = v
		}
	}
	return key
}

// exclude returns a copy of point set s with removed point p.
func exclude(p Point, s map[Point]bool) map[Point]bool {
	e := make(map[Point]bool)
	for k := range s {
		// Skip point p.
		if k == p {
			continue
		}
		e[k] = true
	}
	return e
}

// nearest finds a nearest point from a set of points s to
// another point o.
func nearest(s map[Point]bool, o Point) Point {
	// Distance from point o to point i.
	d := make(map[Point]float64)
	e := exclude(o, s)
	for k := range e {
		d[k] = distance(o, k)
	}
	return minKey(d)
}
