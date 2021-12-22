package gocv

import (
	"testing"
)

func TestIsWidthStrideAlignRGBA(t *testing.T) {
	img := LoadImage("./data/sample.jpg")

	for i := 100; i < 1000; i++ {

		// rgba := ToRGBA(img)
		rgba := ROI4RGBA(img, 10, 10, i, i)
		// rgba := ROI4RGBAv2(img, 300, 300, 34)
		ok := IsWidthStrideAlignRGBA(rgba)
		if !ok {
			panic("Stride RGBA")
		}

		// gray := ToGray(img)
		gray := ROI4GRAY(img, 10, 10, i, i)
		// gray := ROI4GRAYv2(img, 200, 200, 33)
		ok = IsWidthStrideAlignGRAY(gray)
		if !ok {
			panic("Stride GRAY")
		}
	}
}
