package main

import (
	"fmt"

	"github.com/digisan/go-handy-cv/blob"
)

func main() {
	f := func(p byte) bool {
		return p == 0
	}
	blobs := blob.DetectBlob(width, height, step, arr, f)
	fmt.Println(blobs)

	for _, blob := range blobs {
		fmt.Println("loc:", blob.Loc())
		fmt.Println("center:", blob.Center())
		fmt.Println("area:", blob.Area())
	}
}
