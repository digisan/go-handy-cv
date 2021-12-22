package gocv

import "image"

func IsWidthStrideAlignRGBA(rgba *image.RGBA) bool {
	rect := rgba.Bounds()
	return rect.Dx()*4 == rgba.Stride
}

func IsWidthStrideAlignGRAY(gray *image.Gray) bool {
	rect := gray.Bounds()
	return rect.Dx() == gray.Stride
}
