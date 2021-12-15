package gocv

import (
	"image"
)

type PolPt struct {
	ValPt
	R, Phi float64
}

func (ppt *PolPt) GetVPt() *ValPt {
	return &ValPt{
		Point: ppt.Point,
		Val:   ppt.Val,
		Red:   ppt.Red,
		Green: ppt.Green,
		Blue:  ppt.Blue,
		Alpha: ppt.Alpha,
	}
}

func VPts2PPts(vPts []ValPt, centre image.Point) []PolPt {
	pPts := make([]PolPt, len(vPts))
	for i, pt := range vPts {
		pPts[i] = *pt.MkPolar(centre)
	}
	return pPts
}

func Pts2PPts(pts []image.Point, gray, red, green, blue, alpha *image.Gray, centre image.Point) []PolPt {
	return VPts2PPts(Pts2VPts(pts, gray, red, green, blue, alpha), centre)
}

func PPts2Pts(pPts []PolPt) []image.Point {
	pts := make([]image.Point, len(pPts))
	for i, pt := range pPts {
		pts[i] = pt.Point
	}
	return pts
}

func PPts2VPts(pPts []PolPt) []ValPt {
	vPts := make([]ValPt, len(pPts))
	for i, pt := range pPts {
		vPts[i] = *pt.GetVPt()
	}
	return vPts
}
