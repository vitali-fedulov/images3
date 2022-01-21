# Comparing images in Go

Near duplicates and resized images can be found with the package. Function `Open` supports JPEG, PNG and GIF (Go image-package default). But other image types are possible through third-party libraries, because the input for processing is simply image.Image, as you can see in the examples below.

There is only one dependency: my package [hyper](https://github.com/vitali-fedulov/hyper).

**Demo**: [Similar image search and clustering](https://similar.pictures).

`Similar` function gives a verdict whether 2 images are similar or not, based on package-default thresholds. If instead you need similarity metrics and choose your own thresholds, use functions `PropMetric` and `EucMetric`.

If you are planning to process millions of images, comparing images with `Similar` directly may be slow and consume a lot of RAM to keep all icons in memory. To address the problem use a hash table as a preliminary filter (func `CentralHash` and `HashSet`) and read icons by their ids from hard drives. [More info](https://vitali-fedulov.github.io/algorithm-for-hashing-high-dimensional-float-vectors.html) on the hyperspace hashes used in the package.

The library also contains basic functions to open/save/resize images.


## Example of comparing 2 photos with func Similar


```go
package main

import (
	"fmt"
	"github.com/vitali-fedulov/images3"
)

func main() {
	
	// Open photos.
	imgA, _ := images.Open("photoA.jpg")
	imgB, _ := images.Open("photoB.jpg")
	
	// Calculate hashes and image sizes.
	hashA, imgSizeA := images.Hash(imgA)
	hashB, imgSizeB := images.Hash(imgB)
	
	// Image comparison.
	if images.Similar(hashA, hashB, imgSizeA, imgSizeB) {
		fmt.Println("Images are similar.")
	} else {
		fmt.Println("Images are distinct.")
	}
}
```

## Algorithm for image comparison

[Detailed explanation with illustrations](https://vitali-fedulov.github.io/algorithm-for-perceptual-image-comparison.html).

Summary: In the algorithm images are resized to small squares of fixed size.
A number of masks representing several sample pixels are run against the resized
images to calculate average color values. Then the values are compared to
give the similarity verdict. Also image proportions are used to avoid matching
images of distinct shape.
