package gocv

import (
	"image"
	"math"

	cm "github.com/digisan/tiny-gocv/math"
)

type ValPt struct {
	image.Point
	Val, Red, Green, Blue, Alpha byte
}

func NewValPt(pt image.Point) *ValPt {
	return &ValPt{Point: pt}
}

func (pt *ValPt) SetVal(gray, red, green, blue, alpha *image.Gray) {
	if gray != nil {
		pt.Val = gray.GrayAt(pt.X, pt.Y).Y
	}
	if red != nil {
		pt.Red = red.GrayAt(pt.X, pt.Y).Y
	}
	if green != nil {
		pt.Green = green.GrayAt(pt.X, pt.Y).Y
	}
	if blue != nil {
		pt.Blue = blue.GrayAt(pt.X, pt.Y).Y
	}
	if alpha != nil {
		pt.Alpha = alpha.GrayAt(pt.X, pt.Y).Y
	}
}

func (vpt *ValPt) MkPolar(centre image.Point) *PolPt {
	pt := vpt.Point
	return &PolPt{
		R:     cm.DisPt(pt, centre),
		Phi:   math.Atan2(cm.DisPtY(pt, centre), cm.DisPtX(pt, centre)),
		ValPt: *vpt,
	}
}

func Pts2VPts(pts []image.Point, gray, red, green, blue, alpha *image.Gray) []ValPt {
	vPts := make([]ValPt, len(pts))
	for i, pt := range pts {
		vPts[i] = *NewValPt(pt)
		vPts[i].SetVal(gray, red, green, blue, alpha)
	}
	return vPts
}

func VPts2Pts(vPts []ValPt) []image.Point {
	pts := make([]image.Point, len(vPts))
	for i, pt := range vPts {
		pts[i] = pt.Point
	}
	return pts
}
