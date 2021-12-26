package gocv

import (
	"image"
	"math"

	"github.com/digisan/go-generics/pt"
)

func RoundPts(centre image.Point, r int, d float64) (pts []image.Point) {

	cx, cy := centre.X, centre.Y
	for i := 0; i <= r; i += 2 {

		pts1 := RangePixelMap(centre, "X", [2]int{cx - i, cx + i}, d, func(x float64) (ys []float64) {
			R := i
			y := math.Sqrt(float64(R)*float64(R) - x*x)
			ys = append(ys, y, -y)
			return
		})
		pts = append(pts, pts1...)

		pts2 := RangePixelMap(centre, "Y", [2]int{cy - i, cy + i}, d, func(y float64) (xs []float64) {
			R := i
			x := math.Sqrt(float64(R)*float64(R) - y*y)
			xs = append(xs, x, -x)
			return
		})
		pts = append(pts, pts2...)
	}

	pts = pt.MkSet(pts...)
	return pt.Filter(&pts, func(i int, e image.Point) bool { return e.X > 0 && e.Y > 0 })
}
