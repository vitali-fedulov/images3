# Comparing images in Go

Near duplicates and resized images can be found with the package. Function `Open` supports JPEG, PNG and GIF (Go image-package default). But other image types are possible through third-party libraries, because the input for func `Icon` is simply image.Image.

**Demo**: [Similar image search and clustering](https://similar.pictures).

There is only one dependency: my package [hyper](https://github.com/vitali-fedulov/hyper).

`Similar` function gives a verdict whether 2 images are similar or not, based on package-default thresholds. If instead you need similarity metrics and choose your own thresholds, use functions `PropMetric` and `EucMetric`.

If you are planning to process millions of images, comparing images with `Similar` directly may be slow and consume a lot of RAM to keep all icons in memory. To address the problem use a hash table as a preliminary filter (func `CentralHash` and `HashSet`) and read icons by their ids from the hard drive. [More info](https://vitali-fedulov.github.io/algorithm-for-hashing-high-dimensional-float-vectors.html) on the hyperspace hashes.

The library also contains basic functions to open/save/resize images.


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

Summary: Images are resized to small squares of fixed size (here called "icon"). A number of masks representing several sample pixels are run against the resized images to calculate average color values. Then Euclidean distance between the icons is used to give the similarity verdict. Also image proportions are used to avoid matching images of distinct shape.


## Example of comparing 2 photos using hashes instead of Euclidean distance

Hash-based comparison provides rough approximation of image similarity. After that use func `Similar` to get the final verdict. The demo shows only the hash-based similarity testing.

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
