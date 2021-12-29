package images

import (
	"testing"
)

func TestYCbCr(t *testing.T) {
	var r, g, b float32 = 255, 255, 255
	var eY, eCb, eCr float32 = 255, 128, 128
	y, cb, cr := yCbCr(r, g, b)
	// Int values, so the test does not become brittle.
	if int(y) != int(eY) || int(cb) != int(eCb) || int(cr) != int(eCr) {
		t.Errorf("Expected (%v,%v,%v) got (%v,%v,%v).", int(eY), int(eCb),
			int(eCr), int(y), int(cb), int(cr))
	}
	r, g, b = 14, 22, 250
	// From the original external formula.
	eY, eCb, eCr = 45.6, 243.3, 105.5
	y, cb, cr = yCbCr(r, g, b)
	// Int values, so the test does not become brittle.
	if int(y) != int(eY) || int(cb) != int(eCb) || int(cr) != int(eCr) {
		t.Errorf("Expected (%v,%v,%v) got (%v,%v,%v).", int(eY), int(eCb),
			int(eCr), int(y), int(cb), int(cr))
	}
}
