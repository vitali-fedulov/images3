package images

import (
	// N-dimensional space discretization and hashing module.
	"github.com/vitali-fedulov/hyper"
)

// CentralHash generates a central hash for a given icon
// by sampling luma values of well-distributed icon points
// (see var HyperPoints10 to understand how) and later using
// hyper package functions. This hash can then be used for
// a record or a query. When used for a record, you will
// need a hash set made with HashSet function for a query.
// And vice versa.
func CentralHash(
	icon IconT, hyperPoints []Point, epsPercent float64,
	numBuckets int) uint64 {
	return hyper.Decimal(
		hyper.CentralCube(
			LumaValues(icon, hyperPoints),
			0, 255, epsPercent, numBuckets))
}

// HashSet generates a hash set for a given icon
// by sampling luma values of well-distributed icon points
// (see var HyperPoints10 to understand how) and later using
// hyper package functions. The hash set can then be used for
// records or a query. When used for a query, you will
// need a hash made with HashSet function as a record.
// And vice versa.
func HashSet(
	icon IconT, hyperPoints []Point, epsPercent float64,
	numBuckets int) []uint64 {
	return hyper.HashSet(
		hyper.CubeSet(
			LumaValues(icon, HyperPoints10),
			0, 255, epsPercent, numBuckets),
		hyper.Decimal)
}
