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

	for i, blob := range blobs {
		fmt.Printf("---blob---%d\n", i)
		fmt.Println(blob.Y, blob.Idx)
		fmt.Println(blob.Tag())
		fmt.Println(blob.Loc())
		fmt.Println(blob.Area())
		fmt.Println()
	}
}
