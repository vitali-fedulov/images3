package images

// Similar returns similarity verdict based on Euclidean
// and proportion similarity (both with default parameters).
// For large image collections comparing images only with this
// function will be slow. In such case use Fast function
// to get a slice of rough matches, then use Similar
// to get the final verdict.
func Similar(iconA, iconB IconT) bool {

	if PropSimilar(iconA, iconB) {

		if EucSimilar(iconA, iconB) {
			return true
		}
	}
	return false
}

// Fast finds preliminary rough matches before final
// comparison with Similar function. Fast returns a slice
// of ids for images similar to one query image (hashes of one
// image as the query slice).
// Hash table format is map[hash][]ids.
// Hashes for both the query and hash table are produced
// with the Hashes function. It is sufficient to store any
// single hash (uint64 number) from the Hashes output as a key
// in the hash table (not all of them). But the query must
// contain all of them.
func Fast(
	query []uint64, hashTable map[uint64][]uint64) []uint64 {

	for i := range query {
		if ids, ok := hashTable[query[i]]; ok {
			return ids
		}
	}
	return nil
}
