package images

import (
	"path"
	"testing"
)

func testProportions(fA, fB string, isSimilar bool,
	t *testing.T) {
	p := path.Join("testdata", "proportions")
	imgA, err := Open(path.Join(p, fA))
	if err != nil {
		t.Error("Error opening image:", err)
	}
	imgB, err := Open(path.Join(p, fB))
	if err != nil {
		t.Error("Error opening image:", err)
	}
	_, sizeA := Icon(imgA)
	_, sizeB := Icon(imgB)

	if isSimilar == true {
		if !PropSimilar(sizeA, sizeB) {
			t.Errorf("Expecting similarity of %v to %v.", fA, fB)
		}
	}
	if isSimilar == false {
		if PropSimilar(sizeA, sizeB) {
			t.Errorf("Expecting non-similarity of %v to %v.", fA, fB)
		}
	}
}

func TestSimilarByProportions(t *testing.T) {
	testProportions("100x130.png", "100x124.png", true, t)
	testProportions("100x130.png", "100x122.png", false, t)
	testProportions("130x100.png", "260x200.png", true, t)
	testProportions("200x200.png", "260x200.png", false, t)
	testProportions("130x100.png", "124x100.png", true, t)
	testProportions("130x100.png", "122x100.png", false, t)
	testProportions("130x100.png", "130x100.png", true, t)
	testProportions("100x130.png", "130x100.png", false, t)
	testProportions("124x100.png", "260x200.png", true, t)
	testProportions("122x100.png", "260x200.png", false, t)
	testProportions("100x124.png", "100x130.png", true, t)
}

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
