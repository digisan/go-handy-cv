package gocv

import (
	"image"
	"math"
)

func U8ToF64(vb ...byte) (vf []float64) {
	for _, b := range vb {
		vf = append(vf, float64(b))
	}
	return
}

func F64ToU8(vf ...float64) (vb []byte) {
	for _, f := range vf {
		vb = append(vb, byte(math.Round(f)))
	}
	return
}

func I64ToF64(vi ...int) (vf []float64) {
	for _, i := range vi {
		vf = append(vf, float64(i))
	}
	return
}

func F64ToI64(vf ...float64) (vi []int) {
	for _, f := range vf {
		vi = append(vi, int(math.Round(f)))
	}
	return
}

///////////////////////////////////////////////////////////////////////

func ToRGBA(img image.Image) *image.RGBA {
	rect := img.Bounds()
	return ROI4RGBA(img, rect.Min.X, rect.Min.Y, rect.Max.X, rect.Max.Y)
}

func ToCMYK(img image.Image) *image.CMYK {
	rect := img.Bounds()
	return ROI4CMYK(img, rect.Min.X, rect.Min.Y, rect.Max.X, rect.Max.Y)
}

func ToGray(img image.Image) *image.Gray {
	rect := img.Bounds()
	return ROI4GRAY(img, rect.Min.X, rect.Min.Y, rect.Max.X, rect.Max.Y)
}
