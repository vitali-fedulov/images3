package images

import (
	"path"
	"reflect"
	"testing"
)

func TestNewIcon(t *testing.T) {
	icon := NewIcon(4)
	expected := 4 * 4 * 3
	got := len(icon)
	if got != expected {
		t.Errorf(
			"Expected length %d, got %d.", expected, got)
	}
}

func TestArrIndex(t *testing.T) {
	x, y := 2, 3
	size := 4
	ch := 2
	got := ArrIndex(Point{x, y}, size, ch)
	expected := 46
	if got != expected {
		t.Errorf("Expected %d, got %d.", expected, got)
	}
	x, y = 1, 1
	ch = 1
	got = ArrIndex(Point{x, y}, size, ch)
	expected = 21
	if got != expected {
		t.Errorf("Expected %d, got %d.", expected, got)
	}
	x, y = 3, 3
	ch = 0
	got = ArrIndex(Point{x, y}, size, ch)
	expected = 15
	if got != expected {
		t.Errorf("Expected %d, got %d.", expected, got)
	}
}

func TestSet(t *testing.T) {
	icon := NewIcon(4)
	Set(icon, 4, Point{1, 1}, 13.5, 29.9, 95.9)
	expected := IconT{0, 0, 0, 0, 0, 13.5, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 29.9, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 95.9, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0}
	if !reflect.DeepEqual(expected, icon) {
		t.Errorf("Expected %v, got %v.", expected, icon)
	}
}

func TestGet(t *testing.T) {
	icon := IconT{0, 0, 0, 0, 0, 13.5, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 29.9, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 95.9, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0}
	c1, c2, c3 := Get(icon, 4, Point{1, 1})
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
	testDir := path.Join(testDir1, testDir2)
	img, err := Open(path.Join(testDir, imageName))
	if err != nil {
		t.Error(
			"Cannot decode", path.Join(testDir, imageName))
	}
	_, imgSize := Icon(img)
	if imgSize.X != 533 || imgSize.Y != 400 {
		t.Errorf(
			"Expected image size (533, 400), got (%d, %d).",
			imgSize.X, imgSize.Y)
	}
}

func TestLumaValues(t *testing.T) {
	icon := NewIcon(iconSmallSize)
	expectedColor1 := float32(13.1)
	expectedColor2 := float32(9.1)
	Set(icon, iconSmallSize,
		Point{1, 1}, expectedColor1, 29.9, 95.9)
	Set(icon, iconSmallSize,
		Point{9, 5}, expectedColor2, 11.0, 12.9)
	got := LumaValues(icon, []Point{{1, 1}, {9, 5}})
	if float32(got[0]) != expectedColor1 ||
		float32(got[1]) != expectedColor2 {
		t.Errorf(
			`Expected 2 color values %v and %v.
			 Got: %v and %v.`, expectedColor1, expectedColor2,
			float32(got[0]), float32(got[1]))
	}
}
