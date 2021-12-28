package cvmath

import (
	"fmt"
	"image"
	"testing"
)

func TestSegment_String(t *testing.T) {
	seg := NewSegment(image.Pt(2, 3), image.Pt(0, 2))
	fmt.Println(seg)
}

func TestOnSeg(t *testing.T) {
	seg := NewSegment(image.Pt(200, 198), image.Pt(0, 0))
	fmt.Println(seg.Has(image.Pt(19, 19)))
}

func TestIntersection(t *testing.T) {
	seg1 := NewSegment(image.Pt(200, 200), image.Pt(0, 0))
	seg2 := NewSegment(image.Pt(9, 300), image.Pt(150, -20))
	fmt.Println(Intersection(seg1, seg2))
}

func TestSegInterpolate(t *testing.T) {
	seg := NewSegment(image.Pt(9, 300), image.Pt(150, -20))
	pts := seg.Interpolate(2)
	for _, pt := range pts {
		fmt.Println(pt, seg.Has(pt))
	}
}
