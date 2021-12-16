package gocv

import (
	"image"
	"image/draw"
)

//// RGBA

func ROI4RGBA(img image.Image, left, top, right, bottom int) *image.RGBA {
	rect := image.Rect(0, 0, right-left, bottom-top)
	rgba := image.NewRGBA(rect)
	draw.Draw(rgba, rect, img, image.Point{left, top}, draw.Src)
	return rgba
}

func ROI4RGBAv2(img image.Image, cx, cy, sRadius int) *image.RGBA {
	left := cx - sRadius
	top := cy - sRadius
	right := cx + sRadius
	bottom := cy + sRadius
	return ROI4RGBA(img, left, top, right, bottom)
}

//// CMYK

func ROI4CMYK(img image.Image, left, top, right, bottom int) *image.CMYK {
	rect := image.Rect(0, 0, right-left, bottom-top)
	cmyk := image.NewCMYK(rect)
	draw.Draw(cmyk, rect, img, image.Point{left, top}, draw.Src)
	return cmyk
}

func ROI4CMYKv2(img image.Image, cx, cy, sRadius int) *image.CMYK {
	left := cx - sRadius
	top := cy - sRadius
	right := cx + sRadius
	bottom := cy + sRadius
	return ROI4CMYK(img, left, top, right, bottom)
}

//// GRAY

func ROI4GRAY(img image.Image, left, top, right, bottom int) *image.Gray {
	rect := image.Rect(0, 0, right-left, bottom-top)
	gray := image.NewGray(rect)
	draw.Draw(gray, rect, img, image.Point{left, top}, draw.Src)
	return gray
}

func ROI4GRAYv2(img image.Image, cx, cy, sRadius int) *image.Gray {
	left := cx - sRadius
	top := cy - sRadius
	right := cx + sRadius
	bottom := cy + sRadius
	return ROI4GRAY(img, left, top, right, bottom)
}
