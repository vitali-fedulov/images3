package images

import (
	"path"
	"testing"
)

func testEuclidean(fA, fB string, isSimilar bool,
	t *testing.T) {
	p := path.Join("testdata", "euclidean")
	imgA, err := Open(path.Join(p, fA))
	if err != nil {
		t.Error("Error opening image:", err)
	}
	iconA, _ := Icon(imgA)
	imgB, err := Open(path.Join(p, fB))
	if err != nil {
		t.Error("Error opening image:", err)
	}
	iconB, _ := Icon(imgB)
	if isSimilar == true {
		if !EucSimilar(iconA, iconB) {
			t.Errorf("Expecting similarity of %v to %v.", fA, fB)
		}
	}
	if isSimilar == false {
		if EucSimilar(iconA, iconB) {
			t.Errorf("Expecting non-similarity of %v to %v.", fA, fB)
		}
	}
}

func TestSimilarByEuclidean(t *testing.T) {
	testEuclidean("large.jpg", "distorted.jpg", true, t)
	testEuclidean("large.jpg", "flipped.jpg", false, t)
	testEuclidean("large.jpg", "small.jpg", true, t)
}
