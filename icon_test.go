package images3

import (
	"math"
	"path"
	"reflect"
	"testing"
)

func TestSizedIcon(t *testing.T) {
	icon := sizedIcon(4)
	expected := 4 * 4 * 3
	got := len(icon.Pixels)
	if got != expected {
		t.Errorf(
			"Expected length %d, got %d.", expected, got)
	}
}

func TestEmptyIcon(t *testing.T) {
	icon1 := EmptyIcon()
	icon2 := IconT{nil, Point{0, 0}, ""}

	if !reflect.DeepEqual(icon1.Pixels, icon2.Pixels) {
		t.Errorf("Icons' Pixels mismatch. They must be equal: %v %v",
			icon1.Pixels, icon2.Pixels)
	}
	if !reflect.DeepEqual(icon1.ImgSize, icon2.ImgSize) {
		t.Errorf("Icons' ImgSize mismatch. They must be equal: %v %v",
			icon1.ImgSize, icon2.ImgSize)
	}
	if icon1.Path != icon2.Path {
		t.Errorf("Empty-icon Path must be equal to \"\", instead got %v",
			icon1.Path)
	}
}

func TestArrIndex(t *testing.T) {
	x, y := 2, 3
	size := 4
	ch := 2
	got := arrIndex(Point{x, y}, size, ch)
	expected := 46
	if got != expected {
		t.Errorf("Expected %d, got %d.", expected, got)
	}
	x, y = 1, 1
	ch = 1
	got = arrIndex(Point{x, y}, size, ch)
	expected = 21
	if got != expected {
		t.Errorf("Expected %d, got %d.", expected, got)
	}
	x, y = 3, 3
	ch = 0
	got = arrIndex(Point{x, y}, size, ch)
	expected = 15
	if got != expected {
		t.Errorf("Expected %d, got %d.", expected, got)
	}
}

func TestSet(t *testing.T) {
	icon := sizedIcon(4)
	set(icon, 4, Point{1, 1}, 13.5, 29.9, 95.9)
	expected := sizedIcon(4 * 4 * 3)
	expected.Pixels = []float32{0, 0, 0, 0, 0, 13.5, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 29.9, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 95.9, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0}
	if !reflect.DeepEqual(expected, icon) {
		t.Errorf("Expected %v, got %v.", expected, icon)
	}
}

func TestGet(t *testing.T) {
	icon := sizedIcon(4 * 4 * 3)
	icon.Pixels = []float32{0, 0, 0, 0, 0, 13.5, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 29.9, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 95.9, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0}
	c1, c2, c3 := get(icon, 4, Point{1, 1})
	if !(c1 == 13.5 && c2 == 29.9 && c3 == 95.9) {
		t.Errorf(
			"Expected 13.5, 29.9, 95.9, got %v, %v, %v.",
			c1, c2, c3)
	}
}

// Only checking that image size is correct.
func TestIcon(t *testing.T) {
	const (
		testDir1  = "testdata"
		testDir2  = "resample"
		imageName = "nearest533x400.png"
	)
	filePath := path.Join(testDir1, testDir2, imageName)
	img, err := Open(filePath)
	if err != nil {
		t.Error(
			"Cannot decode", filePath)
	}
	icon := Icon(img, filePath)
	if icon.ImgSize.X != 533 || icon.ImgSize.Y != 400 {
		t.Errorf(
			"Expected image size (533, 400), got (%d, %d).",
			icon.ImgSize.X, icon.ImgSize.Y)
	}
}

func TestYCbCr(t *testing.T) {
	var r, g, b float32 = 255, 255, 255
	var eY, eCb, eCr float32 = 255, 128, 128
	y, cb, cr := yCbCr(r, g, b)
	// Int values, so the test does not become brittle.
	if int(y) != int(eY) || int(cb) !=
		int(eCb) || int(cr) != int(eCr) {
		t.Errorf("Expected (%v,%v,%v) got (%v,%v,%v).",
			int(eY), int(eCb), int(eCr),
			int(y), int(cb), int(cr))
	}
	r, g, b = 14, 22, 250
	// From the original external formula.
	eY, eCb, eCr = 45.6, 243.3, 105.5
	y, cb, cr = yCbCr(r, g, b)
	// Int values, so the test does not become brittle.
	if int(y) != int(eY) || int(cb) !=
		int(eCb) || int(cr) != int(eCr) {
		t.Errorf("Expected (%v,%v,%v) got (%v,%v,%v).",
			int(eY), int(eCb), int(eCr),
			int(y), int(cb), int(cr))
	}
}

func TestLumaVector(t *testing.T) {
	iconSize := 11
	icon := sizedIcon(iconSize)
	expectedColor1 := float32(13.1)
	expectedColor2 := float32(9.1)
	set(icon, iconSize,
		Point{1, 1}, expectedColor1, 29.9, 95.9)
	set(icon, iconSize,
		Point{9, 5}, expectedColor2, 11.0, 12.9)
	got := lumaVector(icon, []Point{{1, 1}, {9, 5}})
	if float32(got[0]) != expectedColor1 ||
		float32(got[1]) != expectedColor2 {
		t.Errorf(
			`Expected 2 color values %v and %v.
			 Got: %v and %v.`, expectedColor1, expectedColor2,
			float32(got[0]), float32(got[1]))
	}
}

func testNormalize(src, want IconT, t *testing.T) {
	for i := range src.Pixels {
		if math.Round(float64(src.Pixels[i])) !=
			math.Round(float64(want.Pixels[i])) {
			t.Errorf("Want %v, got %v.", want, src)
			break
		}
	}
}

func TestNormalize(t *testing.T) {

	// 2x2 icon.
	src := sizedIcon(2)
	src.Pixels = []float32{
		0.5, 89, 14, 211,
		9, 193, 20, 14,
		97, 31, 7, 67.9}
	src.normalize(2)
	want := sizedIcon(2)
	want.Pixels = []float32{
		0, 107.20902, 16.35392, 255,
		0, 255, 15.244565, 6.929348,
		255, 68, 0, 172.55}
	testNormalize(src, want, t)

	// 2x2 icon.
	src.Pixels = []float32{
		111, 111, 22, 77,
		99, 99, 255, 33,
		88, 0, 222, 33}
	src.normalize(2)
	want.Pixels = []float32{
		255, 255, 0, 157.58427,
		75.810814, 75.810814, 255, 0,
		101.08108, 0, 255, 37.905407}
	testNormalize(src, want, t)

}
