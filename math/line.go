package cvmath

import (
	"fmt"
	"image"
	"math"
	"sort"

	"github.com/digisan/go-generics/pt"
)

type Segment struct {
	endPt1, endPt2 image.Point
}

func NewSegment(pt1, pt2 image.Point) Segment {
	ptPair := []image.Point{pt1, pt2}
	sort.SliceStable(ptPair, func(i, j int) bool {
		return ptPair[i].X < ptPair[j].X
	})
	return Segment{endPt1: ptPair[0], endPt2: ptPair[1]}
}

func (seg *Segment) String() string {
	return fmt.Sprint(seg.endPt1, seg.endPt2)
}

func (seg *Segment) EndPtX() (x1, x2 int) {
	return seg.endPt1.X, seg.endPt2.X
}

func (seg *Segment) EndPtY() (y1, y2 int) {
	return seg.endPt1.Y, seg.endPt2.Y
}

func (seg *Segment) Has(pt image.Point) bool {
	d1 := DisPt(seg.endPt1, pt)
	d2 := DisPt(seg.endPt2, pt)
	d := DisPt(seg.endPt1, seg.endPt2)
	return d1+d2 < d+ErrChkOnSeg
}

/////////////////////////////////////////////////////////////////

func YaXb(pt1, pt2 image.Point) (a, b float64, vertical bool, vX float64, horizontal bool, hY float64) {
	x1, y1 := float64(pt1.X), float64(pt1.Y)
	x2, y2 := float64(pt2.X), float64(pt2.Y)

	// vertical, "x = n"
	if x1 == x2 {
		vertical = true
		vX = x1
		a, b = 1, 0
		return
	}

	a = (y1 - y2) / (x1 - x2)
	b = y1 - a*x1

	// horizontal, "y = n"
	if a == 0 {
		horizontal = true
		hY = b
		a = 0
	}

	return
}

func Intersection(s1, s2 Segment) (pt *image.Point, inter, coincide bool) {

	a1, b1, v1, vX1, h1, hY1 := YaXb(s1.endPt1, s1.endPt2)
	a2, b2, v2, vX2, h2, hY2 := YaXb(s2.endPt1, s2.endPt2)

	if v1 && v2 {
		if vX1 != vX2 {
			return nil, false, false
		}
		return nil, true, true
	}

	if h1 && h2 {
		if hY1 != hY2 {
			return nil, false, false
		}
		return nil, true, true
	}

	if a1 == a2 {
		if b1 != b2 {
			return nil, false, false
		}
		return nil, true, true
	}

	x := (b2 - b1) / (a1 - a2)
	y := a1*x + b1
	pt = &image.Point{X: int(math.Round(x)), Y: int(math.Round(y))}
	return pt, s1.Has(*pt) && s2.Has(*pt), false
}

func (seg *Segment) Interpolate(step float64) (pts []image.Point) {
	return interpolate(seg.endPt1, seg.endPt2, step)
}

func interpolate(pt1, pt2 image.Point, step float64) (pts []image.Point) {

	x1, x2 := pt1.X, pt2.X
	y1, y2 := pt1.Y, pt2.Y

	xMin := Min(float64(x1), float64(x2))
	xMax := Max(float64(x1), float64(x2))
	yMin := Min(float64(y1), float64(y2))
	yMax := Max(float64(y1), float64(y2))

	a, b, v, vx, h, hy := YaXb(pt1, pt2)

	switch {
	case !v && !h:
		for x := xMin; x <= xMax; x += step {
			y := a*x + b
			pts = append(pts, image.Pt(int(math.Round(x)), int(math.Round(y))))
		}
	case h:
		for x := xMin; x <= xMax; x += step {
			pts = append(pts, image.Pt(int(math.Round(x)), int(math.Round(hy))))
		}
	case v:
		for y := yMin; y <= yMax; y += step {
			pts = append(pts, image.Pt(int(math.Round(vx)), int(math.Round(y))))
		}
	}

	return pt.MkSet(pts...)
}
