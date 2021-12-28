package gocv

import (
	"image"
	"math"

	cm "github.com/digisan/tiny-gocv/math"
)

func GetMean(img *image.Gray) byte {
	return byte(cm.Mean(U8ToF64(img.Pix...)...)) // img.Stride is always harmony with Width, Tested

	// sum, n := 0.0, 0.0
	// rect := img.Bounds()
	// for y := 0; y < rect.Dy(); y++ {
	// 	pHead := img.Pix[y*img.Stride:]
	// 	for x := 0; x < rect.Dx(); x++ {
	// 		pxl := pHead[x]
	// 		sum += float64(pxl)
	// 		n++
	// 	}
	// }
	// return byte(math.Round(sum / n))
}

func GetMedian(img *image.Gray) byte {
	return byte(cm.Median(U8ToF64(img.Pix...)...))
}

func GetMode(img *image.Gray) byte {
	return byte(cm.Mode(U8ToF64(img.Pix...)...))
}

func GetStripeV(img *image.Gray, x int) (stripe []byte) {
	for y := 0; y < img.Rect.Dy(); y++ {
		offset := y * img.Stride
		pixel := img.Pix[offset+x]
		stripe = append(stripe, pixel)
	}
	return
}

func GetStripeH(img *image.Gray, y int) (stripe []byte) {
	offset := img.Stride * y
	line := img.Pix[offset:]
	for x := 0; x < img.Rect.Dx(); x++ {
		stripe = append(stripe, line[x])
	}
	return
}

//////////////////////////////////////////////////////////

func PixelMap(o image.Point, XY string, d float64, f func(float64) []float64) (pts []image.Point) {
	switch XY {
	case "X", "x":
		for _, dy := range f(d) {
			pts = append(pts, image.Point{X: o.X + int(d), Y: o.Y + int(math.Round(dy))})
		}
	case "Y", "y":
		for _, dx := range f(d) {
			pts = append(pts, image.Point{X: o.X + int(math.Round(dx)), Y: o.Y + int(d)})
		}
	}
	return
}

// 'rng' is close range!
func RangePixelMap(o image.Point, XY string, rng [2]int, nSubPxl float64, f func(float64) []float64) (pts []image.Point) {

	s, e := rng[0], rng[1]
	switch XY {
	case "X", "x":
		for x := float64(s); x <= float64(e); x += nSubPxl {
			dx := x - float64(o.X)
			pts = append(pts, PixelMap(o, XY, dx, f)...)
		}
	case "Y", "y":
		for y := float64(s); y <= float64(e); y += nSubPxl {
			dy := y - float64(o.Y)
			pts = append(pts, PixelMap(o, XY, dy, f)...)
		}
	default:
		panic("<XY> must be 'X' or 'Y'")
	}

	return
}

// 'rng' is close range!
func Range2DPixelMap(os []image.Point, XY string, rng [2]int, nSubPxl float64, f func(x float64) []float64) (ptsGrp [][]image.Point) {
	for _, o := range os {
		ptsGrp = append(ptsGrp, RangePixelMap(o, XY, rng, nSubPxl, f))
	}
	return
}

//
func GetPtsVal(gray *image.Gray, pts []image.Point) []byte {

	// --- slow --- //
	// data := gray.Pix
	// offset := gray.Stride
	// for y := 0; y < gray.Rect.Dy(); y++ {
	// 	line := data[offset*y:]
	// 	for x := 0; x < gray.Rect.Dx(); x++ {
	// 		p := line[x]
	// 		for _, pt := range pts {
	// 			if y == pt.Y && x == pt.X {
	// 				vs = append(vs, p)
	// 				break
	// 			}
	// 		}
	// 	}
	// }

	// --- fast --- //
	vs := make([]byte, len(pts))
	for i, pt := range pts {
		vs[i] = gray.GrayAt(pt.X, pt.Y).Y
	}

	return vs
}
