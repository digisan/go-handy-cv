package main

import (
	"fmt"

	"github.com/digisan/tiny-gocv/blob"
)

func main() {
	f := func(x, y int, p byte) bool {
		return p == 0
	}
	blobs := blob.DetectBlob(width, height, step, arr, f)
	fmt.Println(blobs)
	fmt.Println()

	for _, blob := range blobs {
		fmt.Println("loc:", blob.Loc())
		fmt.Println("center:", blob.Center())
		fmt.Println("area:", blob.Area())
		fmt.Println()
	}
}
