package cvmath

import (
	"fmt"
	"image"
	"reflect"
	"testing"
)

func TestSegment_String(t *testing.T) {
	seg := NewSegment(image.Pt(2, 3), image.Pt(0, 2), true)
	fmt.Println(seg)
	seg = NewSegment(image.Pt(2, 3), image.Pt(0, 2), false)
	fmt.Println(seg)
}

func TestOnSeg(t *testing.T) {
	seg := NewSegment(image.Pt(200, 198), image.Pt(0, 0), true)
	fmt.Println(seg.Has(image.Pt(19, 19)))
	seg = NewSegment(image.Pt(200, 198), image.Pt(0, 0), false)
	fmt.Println(seg.Has(image.Pt(19, 19)))
}

func TestIntersection(t *testing.T) {
	seg1 := NewSegment(image.Pt(200, 200), image.Pt(0, 0), true)
	seg2 := NewSegment(image.Pt(9, 300), image.Pt(150, -20), true)
	fmt.Println(Intersection(seg1, seg2))

	seg1 = NewSegment(image.Pt(200, 200), image.Pt(0, 0), false)
	seg2 = NewSegment(image.Pt(9, 300), image.Pt(150, -20), false)
	fmt.Println(Intersection(seg1, seg2))
}

func TestSegInterpolate(t *testing.T) {
	seg := NewSegment(image.Pt(150, -20), image.Pt(9, 300), true)
	pts := seg.Interpolate(2)
	for _, pt := range pts {
		fmt.Println(pt, seg.Has(pt))
	}

	fmt.Println("----------------------------------------------")

	seg = NewSegment(image.Pt(150, -20), image.Pt(9, 300), false)
	pts = seg.Interpolate(2)
	for _, pt := range pts {
		fmt.Println(pt, seg.Has(pt))
	}
}

func TestSegments_Len(t *testing.T) {

	s, _ := NewSegments(image.Pt(0, 0), image.Pt(100, 100), image.Pt(100, 0))

	type fields struct {
		segs []Segment
	}
	tests := []struct {
		name       string
		fields     fields
		wantLength float64
	}{
		// TODO: Add test cases.
		{
			name:       "",
			fields:     fields(s),
			wantLength: 241.4213562373095,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			segs := &Segments{
				segs: tt.fields.segs,
			}
			if gotLength := segs.Len(); gotLength != tt.wantLength {
				t.Errorf("Segments.Len() = %v, want %v", gotLength, tt.wantLength)
			}
		})
	}
}

func TestSegments_Reverse(t *testing.T) {

	s, _ := NewSegments(image.Pt(0, 0), image.Pt(100, 100), image.Pt(100, 0), image.Pt(-30, 20))
	rs, _ := NewSegments(image.Pt(-30, 20), image.Pt(100, 0), image.Pt(100, 100), image.Pt(0, 0))

	fmt.Print(s)
	fmt.Println(s.Len())
	fmt.Print(rs)
	fmt.Println(rs.Len())

	fmt.Println()
	fmt.Println(s.DivideBy(image.Pt(29, 30)))
	fmt.Println()

	fmt.Println()
	fmt.Println(rs.DivideBy(image.Pt(29, 30)))
	fmt.Println()

	type fields struct {
		segs []Segment
	}
	tests := []struct {
		name   string
		fields fields
		want   Segments
	}{
		// TODO: Add test cases.
		{
			name:   "",
			fields: fields(s),
			want:   rs,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Segments{
				segs: tt.fields.segs,
			}
			if got := s.Reverse(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Segments.Reverse() = %v, want %v", got, tt.want)
			}
		})
	}
}
