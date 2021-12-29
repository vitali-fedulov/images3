package images

import (
	"image"
	"path"
	"reflect"
	"testing"
)

const (
	testDir1 = "testdata"
	testDir2 = "resample"
)

func TestResampleByNearest(t *testing.T) {
	testDir := path.Join(testDir1, testDir2)
	tables := []struct {
		inFile     string
		srcX, srcY int
		outFile    string
		dstX, dstY int
	}{
		{"original.png", 533, 400,
			"nearest100x100.png", 100, 100},
		{"nearest100x100.png", 100, 100,
			"nearest533x400.png", 533, 400},
	}

	for _, table := range tables {
		inImg, err := Open(path.Join(testDir, table.inFile))
		if err != nil {
			t.Error("Cannot decode", path.Join(testDir, table.inFile))
		}
		outImg, err := Open(path.Join(testDir, table.outFile))
		if err != nil {
			t.Error("Cannot decode", path.Join(testDir, table.outFile))
		}
		resampled, srcX, srcY := ResampleByNearest(inImg,
			table.dstX, table.dstY)
		if !reflect.DeepEqual(
			outImg.(*image.RGBA), &resampled) || table.srcX != srcX ||
			table.srcY != srcY {
			t.Errorf(
				"Resample data do not match for %s and %s.",
				table.inFile, table.outFile)
		}
	}
}
