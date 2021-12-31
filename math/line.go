package cvmath

import (
	"fmt"
	"image"
	"math"
	"sort"

	"github.com/digisan/go-generics/pt"
)

type Segment struct {
	ept1, ept2 image.Point
}

func NewSegment(pt1, pt2 image.Point, sortX bool) Segment {
	ptPair := []image.Point{pt1, pt2}
	if sortX {
		sort.SliceStable(ptPair, func(i, j int) bool {
			return ptPair[i].X < ptPair[j].X
		})
	}
	return Segment{ept1: ptPair[0], ept2: ptPair[1]}
}

func (s Segment) String() string {
	return fmt.Sprint(s.ept1, s.ept2)
}

func (s *Segment) Reverse() Segment {
	return NewSegment(s.ept2, s.ept1, false)
}

func (s *Segment) EndPtX() (x1, x2 int) {
	return s.ept1.X, s.ept2.X
}

func (s *Segment) EndPtY() (y1, y2 int) {
	return s.ept1.Y, s.ept2.Y
}

func (s *Segment) Has(pt image.Point) bool {
	d1 := DisPt(s.ept1, pt)
	d2 := DisPt(s.ept2, pt)
	d := DisPt(s.ept1, s.ept2)
	return d1+d2 < d+ErrChkOnSeg
}

func (s *Segment) Len() float64 {
	return DisPt(s.ept1, s.ept2)
}

func (s *Segment) DivideBy(pt image.Point) (float64, float64) {
	return DisPt(s.ept1, pt), DisPt(s.ept2, pt)
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

func Intersection(s1, s2 Segment) (pt *image.Point, isInter, isCoincide bool) {

	a1, b1, v1, vX1, h1, hY1 := YaXb(s1.ept1, s1.ept2)
	a2, b2, v2, vX2, h2, hY2 := YaXb(s2.ept1, s2.ept2)

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

func (s *Segment) Interpolate(step float64) (pts []image.Point) {
	return interpolate(s.ept1, s.ept2, step)
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

/////////////////////////////////////////////////////////////////

type Segments struct {
	segs []Segment
}

func NewSegments(pts ...image.Point) (s Segments, e error) {
	if len(pts) < 2 {
		e = fmt.Errorf("points count must >= 2")
		return s, e
	}
	for i := 0; i < len(pts)-1; i++ {
		s.segs = append(s.segs, NewSegment(pts[i], pts[i+1], false))
	}
	return
}

func SetSegments(segs ...Segment) Segments {
	return Segments{segs: segs}
}

func (s Segments) String() (info string) {
	for _, seg := range s.segs {
		info += fmt.Sprintln(seg)
	}
	return
}

func (s *Segments) HasPt(pt image.Point) bool {
	for _, seg := range s.segs {
		if seg.Has(pt) {
			return true
		}
	}
	return false
}

func (s *Segments) HasSeg(seg Segment) bool {
	for _, s := range s.segs {
		if s == seg || s.Reverse() == seg {
			return true
		}
	}
	return false
}

func (s *Segments) Len() (length float64) {
	for _, s := range s.segs {
		length += s.Len()
	}
	return
}

func (s *Segments) Points() (pts []image.Point) {
	for i, seg := range s.segs {
		if i == 0 {
			pts = append(pts, seg.ept1, seg.ept2)
			continue
		}
		pts = append(pts, seg.ept2)
	}
	return
}

func (s *Segments) Reverse() Segments {
	pts := s.Points()
	for i := 0; i < len(pts)/2; i++ {
		head, tail := i, len(pts)-1-i
		pts[head], pts[tail] = pts[tail], pts[head]
	}
	segs, err := NewSegments(pts...)
	if err != nil {
		panic(err)
	}
	return segs
}

func (s *Segments) DivideBy(pt image.Point) (float64, float64) {
	k := -1
	l1, l2 := 0.0, 0.0
	for i, seg := range s.segs {
		if seg.Has(pt) {
			k = i
			l1, l2 = seg.DivideBy(pt)
			break
		}
	}
	if k == -1 {
		return 0, 0
	}
	al1, al2 := 0.0, 0.0
	for i := 0; i < k; i++ {
		al1 += s.segs[i].Len()
	}
	for i := k + 1; i < len(s.segs); i++ {
		al2 += s.segs[i].Len()
	}
	return al1 + l1, al2 + l2
}

func IntersectionSegs(segs1, segs2 Segments) (pts []*image.Point, isInter, isCoincide []bool) {
	for _, s1 := range segs1.segs {
		for _, s2 := range segs2.segs {
			pt, inter, coincide := Intersection(s1, s2)
			pts = append(pts, pt)
			isInter = append(isInter, inter)
			isCoincide = append(isCoincide, coincide)
		}
	}
	return
}
