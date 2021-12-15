package cvmath

import (
	"image"
	"log"
	"math"

	"github.com/digisan/go-generics/f64"
)

func Max(data ...float64) float64 {
	return f64.Max(data...)
}

func Min(data ...float64) float64 {
	return f64.Min(data...)
}

func MaxIdx(data ...float64) (float64, int) {
	return f64.MaxIdx(data...)
}

func MinIdx(data ...float64) (float64, int) {
	return f64.MinIdx(data...)
}

func MaxMinAbs(data ...float64) (float64, float64) {
	tempAbs := f64.FM(data, nil, func(i int, e float64) float64 {
		if e < 0 {
			return -e
		}
		return e
	})
	return Max(tempAbs...), Min(tempAbs...)
}

func Sum(data ...float64) (sum float64) {
	for _, v := range data {
		sum += v
	}
	return
}

func Average(data ...float64) float64 {
	return Sum(data...) / float64(len(data))
}

func StdDev(data ...float64) float64 {
	sum2 := 0.0
	ave := Average(data...)
	for _, v := range data {
		d2 := (v - ave) * (v - ave)
		sum2 += d2
	}
	return math.Sqrt(sum2 / float64(len(data)))
}

func DotProduct(v1, v2 []float64) (dp float64) {
	if len(v1) != len(v2) {
		log.Fatalf("DotProduct vector dimensions error")
	}
	for i := 0; i < len(v1); i++ {
		dp += v1[i] * v2[i]
	}
	return
}

func Derivative1(data ...float64) (ret []float64) {
	dp := []float64{86, -142, -193, -126, 126, 193, 142, -86}
	ret = make([]float64, len(data))
	for i := 4; i < len(data)-4; i++ {
		// ret[i] = (86*temp[i-4] - 142*temp[i-3] - 193*temp[i-2] - 126*temp[i-1] + 126*temp[i+1] + 193*temp[i+2] + 142*temp[i+3] - 86*temp[i+4]) / 1188
		temp := data[i-4 : i+5]
		ret[i] = DotProduct(dp, temp) / 1188
	}
	return
}

func Smooth9(pts []float64) (ret []float64) {
	sp9 := []float64{-21, 14, 39, 54, 59, 54, 39, 14, -21}
	ret = make([]float64, len(pts))
	copy(ret, pts[:4])
	copy(ret[len(ret)-4:], pts[len(pts)-4:])
	for i := 4; i < len(pts)-4; i++ {
		// ret[i] = ((-21)*pts[i-4] + 14*pts[i-3] + 39*pts[i-2] + 54*pts[i-1] + 59*pts[i] + 54*pts[i+1] + 39*pts[i+2] + 14*pts[i+3] - 21*pts[i+4]) / 231
		pts9 := pts[i-4 : i+5]
		ret[i] = DotProduct(sp9, pts9) / 231
	}
	for i := 0; i < len(ret); i++ {
		if ret[i] < 0 {
			ret[i] = 0
		}
	}
	return
}

func Dis(a, b float64) float64 {
	return math.Abs(a - b)
}

func DisInt(a, b int) int {
	fa, fb := float64(a), float64(b)
	return int(Dis(fa, fb))
}

func DisByte(a, b byte) byte {
	fa, fb := float64(a), float64(b)
	return byte(Dis(fa, fb))
}

func DisPt(pt1, pt2 image.Point) float64 {
	dx := float64(DisInt(pt1.X, pt2.X))
	dy := float64(DisInt(pt1.Y, pt2.Y))
	return math.Sqrt(dx*dx + dy*dy)
}

func DisPtX(pt1, pt2 image.Point) float64 {
	return float64(DisInt(pt1.X, pt2.X))
}

func DisPtY(pt1, pt2 image.Point) float64 {
	return float64(DisInt(pt1.Y, pt2.Y))
}
