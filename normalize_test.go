package images

import (
	"reflect"
	"testing"
)

func testNormalize(src, want IconT, numPix int, t *testing.T) {
	dst := Normalize(src, numPix)
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("Want %v, got %v.", want, dst)
	}
}

func TestNormalize(t *testing.T) {

	// 2x2 icon.
	src := IconT{0.5, 89, 14, 211,
		9, 193, 20, 14,
		97, 31, 7, 67.9}
	want := IconT{0, 107.20902, 16.35392, 255,
		0, 255, 15.244565, 6.929348,
		255, 68, 0, 172.55}
	testNormalize(src, want, 4, t)

	// 2x2 icon.
	src = IconT{111, 111, 22, 77,
		99, 99, 255, 33,
		88, 0, 222, 33}
	want = IconT{255, 255, 0, 157.58427,
		75.810814, 75.810814, 255, 0,
		101.08108, 0, 255, 37.905407}
	testNormalize(src, want, 4, t)

}
