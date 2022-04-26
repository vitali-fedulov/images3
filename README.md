# Comparing images in Go

Near duplicates and resized images can be found with the package.

**Demo**: [similar image clustering](https://similar.pictures) based on same algorithm.

**This is the latest major version** of [images](https://github.com/vitali-fedulov/images).

The only dependency is [hyper](https://github.com/vitali-fedulov/hyper) package, which in turn does not have any dependencies. If you are not using hashes, you can remove this dependency by deleting files [hashes.go, hashes_test.go] from your fork.

Func `Similar` gives a verdict whether 2 images are similar with well-tested default thresholds.

Use func `EucMetric`, when you need different precision or want to sort by similarity. Func `PropMetric` can be used for customization of image proportion threshold.

Func `Open` supports JPEG, PNG and GIF. But other image types are possible through third-party libraries, because func `Icon` input is `image.Image`.

For search in billions of images, use a hash table for preliminary filtering (see the 2nd example below).

[Go doc](https://pkg.go.dev/github.com/vitali-fedulov/images3) for code reference.

## Example of comparing 2 photos with func Similar

```go
package main

import (
	"fmt"
	"github.com/vitali-fedulov/images3"
)

func main() {

	// Photos to compare.
	path1 := "1.jpg"
	path2 := "2.jpg"

	// Open photos (ignoring errors here).
	img1, _ := images3.Open(path1)
	img2, _ := images3.Open(path2)

	// Icons are compact image representations.
	icon1 := images3.Icon(img1, path1)
	icon2 := images3.Icon(img2, path2)

	// Comparison.
	if images3.Similar(icon1, icon2) {
		fmt.Println("Images are similar.")
	} else {
		fmt.Println("Images are distinct.")
	}
}
```

## Algorithm used

[Detailed explanation](https://vitali-fedulov.github.io/algorithm-for-perceptual-image-comparison.html), also as a [PDF](https://github.com/vitali-fedulov/research/blob/main/Algorithm%20for%20perceptual%20image%20comparison.pdf).

Summary: Images are resized in a special way to squares of fixed size called "icons". Euclidean distance between the icons is used to give the similarity verdict. Also image proportions are used to avoid matching images of distinct shape.

## Customization suggestions

To increase precision you can either use your own thresholds in func `EucMetric` (and `PropMetric`) OR generate icons for image sub-regions and compare those icons.

To speedup file processing you may want to generate icons for available image thumbnails. Specifically, many JPEG images contain [EXIF thumbnails](https://www.similar.pictures/jpeg-thumbnail-reader.html), you could considerably speedup the reads by using decoded thumbnails to feed into func `Icon`. External packages to read thumbnails: [1](https://github.com/dsoprea/go-exif) and [2](https://github.com/rwcarlsen/goexif). A note of caution: in rare cases there could be [issues](https://security.stackexchange.com/questions/116552/the-history-of-thumbnails-or-just-a-previous-thumbnail-is-embedded-in-an-image/201785#201785) with thumbnails not matching image content. EXIF standard specification: [1](https://www.media.mit.edu/pia/Research/deepview/exif.html) and [2](https://www.exif.org/Exif2-2.PDF).


## Example of comparing 2 photos using hashes

Hash-based comparison provides fast and RAM-friendly rough approximation of image similarity, when you need to process millions of images. After matching hashes use func `Similar` to get the final verdict. The demo shows only the hash-based similarity testing in its simplified form (without using actual hash table).

```go
package main

import (
	"fmt"
	"github.com/vitali-fedulov/images3"
)

func main() {

	// Paths to photos.
	path1 := "1.jpg"
	path2 := "2.jpg"

	// Hyper space parameters.
	epsPct := 0.25
	numBuckets := 4
	
	// Open photos (skipping error handling for clarity).
	img1, _ := images3.Open(path1)
	img2, _ := images3.Open(path2)

	// Make icons. They are image representations for comparison.
	icon1 := images3.Icon(img1, path1)
	icon2 := images3.Icon(img2, path2)


	// Hash table values.

	// Value to save to the hash table as a key with corresponding
	// image ids. Table structure: map[centralHash][]imageId.
	// imageId is simply an image number in a directory tree.
	centralHash := images3.CentralHash(
		icon1, images3.HyperPoints10, epsPct, numBuckets)

	// Hash set to be used as a query to the hash table. Each hash from
	// the hashSet has to be checked against the hash table.
	// See more info in the package "hyper" README.
	hashSet := images3.HashSet(
		icon2, images3.HyperPoints10, epsPct, numBuckets)

	// Checking hash matches. In full implementation this will
	// be done on the hash table map[centralHash][]imageId.
	foundSimilarImage := false
	for _, hash := range hashSet {
		if centralHash == hash {
			foundSimilarImage = true
		}
	}

	// Image comparison result.
	if foundSimilarImage {
		fmt.Println("Images are similar.")
	} else {
		fmt.Println("Images are distinct.")
	}
}
```
