# Comparing images in Go

**This is the LATEST major version** of old package [images](https://github.com/vitali-fedulov/images).

Near duplicates and resized images can be found with the package. Function `Open` supports JPEG, PNG and GIF (Go image-package default). But other image types are possible through third-party libraries, because the input for func `Icon` is simply image.Image. There is only one dependency: my [hyper](https://github.com/vitali-fedulov/hyper) package.

**Demo**: [Similar image search and clustering](https://similar.pictures).

`Similar` function gives a verdict whether 2 images are similar based on Euclidean distance between specially constructed signatures ("icons") and package-default thresholds. If you prefer your own thresholds or sort by similarity metrics, use functions `PropMetric` and `EucMetric` to get metric values.

If you work with millions of images, using func `Similar` directly could be slow and consume a lot of RAM to keep icons in memory. To address the problem use a hash table as a preliminary filter (see example 2 below). [More info](https://vitali-fedulov.github.io/algorithm-for-hashing-high-dimensional-float-vectors.html) on the hyperspace hashes.

The package also contains basic functions to open/save/resize images.

[Go doc](https://pkg.go.dev/github.com/vitali-fedulov/images3) for code reference.

How to improve precision: Since icon is a tiny compressed representation of the whole large image, in some cases you may want to increase precision by comparing image sub-parts. For that you can generate your own image.Image for image sub-parts and compare icons for those sub-parts.

Opening images is the most time-consuming operation, but since many JPEG images contain [EXIF thumbnails](https://www.similar.pictures/jpeg-thumbnail-reader.html), you could considerably speedup the reads by using decoded thumbnails to feed into func `Icon`. External packages to read thumbnails: [1](https://github.com/dsoprea/go-exif) and [2](https://github.com/rwcarlsen/goexif). A note of caution: in rare cases there could be [issues](https://security.stackexchange.com/questions/116552/the-history-of-thumbnails-or-just-a-previous-thumbnail-is-embedded-in-an-image/201785#201785) with thumbnails not matching image content. EXIF standard specification: [1](https://www.media.mit.edu/pia/Research/deepview/exif.html) and [2](https://www.exif.org/Exif2-2.PDF).


## Example of comparing 2 photos with func Similar

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

	// Open photos (skipping error handling for clarity).
	img1, _ := images3.Open(path1)
	img2, _ := images3.Open(path2)

	// Make icons. They are compact image representations.
	icon1 := images3.Icon(img1, path1)
	icon2 := images3.Icon(img2, path2)

	// Image comparison.
	if images3.Similar(icon1, icon2) {
		fmt.Println("Images are similar.")
	} else {
		fmt.Println("Images are distinct.")
	}
}
```

## Algorithm for image comparison

[Detailed explanation with illustrations](https://vitali-fedulov.github.io/algorithm-for-perceptual-image-comparison.html).

Summary: Images are resized to small squares of fixed size (here called "icon"). A box filter is run against the resized images to calculate average color values. Then Euclidean distance between the icons is used to give the similarity verdict. Also image proportions are used to avoid matching images of distinct shape.


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
